package render

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/resp"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/repositories"
	"bbs-go/internal/services"
	"strconv"

	"github.com/mlogclub/simple/sqls"
)

func BuildTaskGroups() []resp.TaskGroupInfo {
	groups := []resp.TaskGroupInfo{
		{Key: constants.TaskGroupDaily},
		{Key: constants.TaskGroupNewbie},
		{Key: constants.TaskGroupAchievement},
	}

	for i := range groups {
		key := groups[i].Key
		groups[i].Name = locales.Get("task.group." + string(key))
		if groups[i].Name == "" {
			groups[i].Name = string(key)
		}
	}

	return groups
}

func BuildTask(cfg *models.TaskConfig, user *models.User, now int64) resp.TaskResponse {
	var progress *resp.TaskProgressResponse

	if user != nil {
		periodKey := 0
		if constants.TaskPeriod(cfg.Period) != constants.TaskPeriodLifetime {
			if pkStr := services.TaskEngineService.GetPeriodKey(constants.TaskPeriod(cfg.Period), now); pkStr != "" {
				if v, err := strconv.Atoi(pkStr); err == nil {
					periodKey = v
				}
			}
		}

		ute := repositories.UserTaskEventRepository.Take(sqls.DB(), "user_id = ? AND period_key = ? AND task_id = ?", user.Id, periodKey, cfg.Id)
		progress = &resp.TaskProgressResponse{
			PeriodKey:      periodKey,
			EventProgress:  0,
			EventTarget:    cfg.EventCount,
			FinishedCount:  0,
			MaxFinishCount: cfg.MaxFinishCount,
		}
		if ute != nil {
			progress.EventProgress = ute.EventCount
			progress.FinishedCount = ute.TaskFinishCount
		}
	}

	return resp.TaskResponse{
		Id:             cfg.Id,
		GroupName:      cfg.GroupName,
		Title:          cfg.Title,
		Description:    cfg.Description,
		EventType:      cfg.EventType,
		Period:         constants.TaskPeriod(cfg.Period),
		EventCount:     cfg.EventCount,
		MaxFinishCount: cfg.MaxFinishCount,
		Score:          cfg.Score,
		Exp:            cfg.Exp,
		BadgeId:        cfg.BadgeId,
		BtnName:        cfg.BtnName,
		ActionUrl:      cfg.ActionUrl,
		SortNo:         cfg.SortNo,
		StartTime:      cfg.StartTime,
		EndTime:        cfg.EndTime,
		Status:         cfg.Status,
		UserProgress:   progress,
	}
}
