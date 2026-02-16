package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/repositories"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

// TaskEngineService 负责处理用户事件驱动的任务进度与奖励发放
var TaskEngineService = newTaskEngineService()

func newTaskEngineService() *taskEngineService {
	return &taskEngineService{}
}

type taskEngineService struct{}

// GetPeriodKey 根据周期类型返回 periodKey（毫秒时间戳）
func (s *taskEngineService) GetPeriodKey(period constants.TaskPeriod, ms int64) string {
	if period == constants.TaskPeriodLifetime {
		return "0"
	}
	if ms <= 0 {
		ms = dates.NowTimestamp()
	}
	t := time.UnixMilli(ms).In(time.Local)

	switch period {
	case constants.TaskPeriodDaily:
		return t.Format("20060102")
	case constants.TaskPeriodWeekly:
		year, week := t.ISOWeek()
		return fmt.Sprintf("%04d%02d", year, week)
	case constants.TaskPeriodMonthly:
		return t.Format("200601")
	case constants.TaskPeriodYearly:
		return t.Format("2006")
	default:
		return "0"
	}
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	// Fallback for drivers that don't translate errors consistently.
	// - MySQL: "Duplicate entry"
	// - SQLite: "UNIQUE constraint failed"
	msg := err.Error()
	return strings.Contains(msg, "Duplicate entry") || strings.Contains(msg, "UNIQUE constraint failed")
}

// HandleUserEvent 处理用户行为事件（例如 topic.create/comment.create/checkin 等）
// - 更新 UserTaskEvent 进度
// - 产生 UserTaskLog
// - 自动发放 Score/Exp/Badge
func (s *taskEngineService) HandleUserEvent(userId int64, eventType string, eventTime int64) {
	if userId <= 0 || eventType == "" {
		return
	}
	if eventTime <= 0 {
		eventTime = dates.NowTimestamp()
	}

	taskConfigs := repositories.TaskConfigRepository.FindByEventType(eventType)
	if len(taskConfigs) == 0 {
		return
	}

	if err := sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		for i := range taskConfigs {
			taskConfig := &taskConfigs[i]
			userTaskLog, err := s.updateProgressAndCreateLog(ctx, userId, taskConfig, eventTime)
			if err != nil {
				return err
			}
			if userTaskLog != nil {
				if err := s.grantReward(ctx, userTaskLog); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		slog.Error("task progress update failed", slog.Any("userId", userId), slog.String("eventType", eventType), slog.Any("err", err))
	}
}

func (s *taskEngineService) updateProgressAndCreateLog(ctx *sqls.TxContext, userId int64, cfg *models.TaskConfig, eventTime int64) (*models.UserTaskLog, error) {
	now := eventTime
	if now <= 0 {
		now = dates.NowTimestamp()
	}

	if !s.isTaskActive(cfg, now) {
		return nil, nil
	}
	if cfg.EventCount <= 0 || cfg.MaxFinishCount <= 0 {
		return nil, nil
	}

	periodKey := 0
	if constants.TaskPeriod(cfg.Period) != constants.TaskPeriodLifetime {
		periodKeyStr := s.GetPeriodKey(constants.TaskPeriod(cfg.Period), now)
		if periodKeyStr != "" && periodKeyStr != "0" {
			if v, err := strconv.Atoi(periodKeyStr); err == nil {
				periodKey = v
			}
		}
	}

	userTaskEvent := repositories.UserTaskEventRepository.Take(ctx.Tx, "user_id = ? AND period_key = ? AND task_id = ?", userId, periodKey, cfg.Id)
	if userTaskEvent == nil {
		userTaskEvent = &models.UserTaskEvent{
			UserId:          userId,
			PeriodKey:       periodKey,
			TaskId:          cfg.Id,
			EventCount:      0,
			TaskFinishCount: 0,
			CreateTime:      now,
			UpdateTime:      now,
		}
		if err := repositories.UserTaskEventRepository.Create(ctx.Tx, userTaskEvent); err != nil {
			return nil, err
		}
	}

	remainingFinish := cfg.MaxFinishCount - userTaskEvent.TaskFinishCount
	if remainingFinish <= 0 {
		return nil, nil
	}

	// 当前实现下，单次事件最多完成一次任务
	userTaskEvent.EventCount++

	if userTaskEvent.EventCount < cfg.EventCount { // 任务没完成
		userTaskEvent.UpdateTime = now
		if err := repositories.UserTaskEventRepository.Update(ctx.Tx, userTaskEvent); err != nil {
			return nil, err
		}
		return nil, nil
	} else {
		userTaskEvent.TaskFinishCount++ // 任务完成次数加一
		userTaskEvent.EventCount = 0    // 任务完成一次后，剩余计数归零
		userTaskEvent.UpdateTime = now
		if err := repositories.UserTaskEventRepository.Update(ctx.Tx, userTaskEvent); err != nil {
			return nil, err
		}

		// 完成一次任务，记录一条日志
		log := &models.UserTaskLog{
			UserId:     userId,
			PeriodKey:  periodKey,
			TaskId:     cfg.Id,
			FinishNo:   userTaskEvent.TaskFinishCount,
			Score:      cfg.Score,
			Exp:        cfg.Exp,
			BadgeId:    cfg.BadgeId,
			CreateTime: now,
			UpdateTime: now,
		}
		if err := repositories.UserTaskLogRepository.Create(ctx.Tx, log); err != nil {
			// 已存在则跳过（避免重复发奖）
			if isDuplicateKeyError(err) {
				return nil, nil
			}
			return nil, err
		}
		return log, nil
	}
}

func (s *taskEngineService) grantReward(ctx *sqls.TxContext, logRow *models.UserTaskLog) error {
	var (
		sourceType = constants.EntityTask
		sourceId   = strconv.FormatInt(logRow.Id, 10)
	)

	// 发放积分
	if logRow.Score != 0 {
		UserService.addScore(ctx, logRow.UserId, logRow.Score, sourceType, sourceId, "task reward")
	}

	// 发放经验
	if logRow.Exp != 0 {
		if err := UserService.addExpTx(ctx, logRow.UserId, logRow.Exp, sourceType, sourceId, "task reward"); err != nil {
			return err
		}
	}

	// 发放徽章
	if logRow.BadgeId > 0 {
		if err := UserBadgeService.Give(ctx, logRow.UserId, logRow.BadgeId, sourceType, sourceId); err != nil {
			return err
		}
	}

	return nil
}

func (s *taskEngineService) isTaskActive(cfg *models.TaskConfig, now int64) bool {
	if cfg.Status != constants.StatusOk {
		return false
	}
	if cfg.StartTime > 0 && now < cfg.StartTime {
		return false
	}
	if cfg.EndTime > 0 && now > cfg.EndTime {
		return false
	}
	return true
}
