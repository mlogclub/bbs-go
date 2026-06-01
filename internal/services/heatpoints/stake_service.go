package heatpoints

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"fmt"
	"sync"
	"time"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

// StakeService 质押服务
type StakeService struct{}

var Stake = &StakeService{}

// flameLevelCache 火焰等级缓存（用于滞回区间计算）
var flameLevelCache sync.Map

// StakeCreateRequest 创建质押请求
type StakeCreateRequest struct {
	TopicId    int64
	UserId     int64
	HeatPoints int
}

// StakeCreateResponse 创建质押响应
type StakeCreateResponse struct {
	StakeId        int64  `json:"stakeId"`
	RemainingQuota int    `json:"remainingQuota"`
	RiskLevel      string `json:"riskLevel"`
	RiskHint       string `json:"riskHint"`
	FlameLevel     int    `json:"flameLevel"`
}

// Create 创建质押
func (s *StakeService) Create(req *StakeCreateRequest) (*StakeCreateResponse, error) {
	// 获取用户
	var user models.User
	if err := sqls.DB().First(&user, req.UserId).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 验证质押额度
	if req.HeatPoints < constants.HeatPointsStakeMinAmount {
		return nil, fmt.Errorf("质押额度不能少于 %d 热度点", constants.HeatPointsStakeMinAmount)
	}

	// 检查余额（扣除冷却中的额度）
	cooldownPoints := s.getCooldownPoints(req.UserId)
	availableBalance := user.HeatPoints - cooldownPoints
	if availableBalance < req.HeatPoints {
		if user.HeatPoints >= req.HeatPoints {
			return nil, fmt.Errorf("可用热度点不足（含冷却中 %d 点）", cooldownPoints)
		}
		return nil, fmt.Errorf("热度点余额不足")
	}

	// 检查每日配额
	quotaUsed := s.GetTodayQuotaUsed(req.UserId)
	if quotaUsed >= constants.HeatPointsStakeQuotaDaily {
		return nil, fmt.Errorf("今日质押次数已达上限")
	}

	// 验证帖子存在且状态正常
	var topic models.Topic
	if err := sqls.DB().First(&topic, req.TopicId).Error; err != nil {
		return nil, fmt.Errorf("帖子不存在")
	}
	if topic.Status != constants.StatusOk {
		return nil, fmt.Errorf("该帖子已被删除或不可用")
	}

	// 检查单人单帖质押上限（阶梯式：保底 + 百分比）
	activeCirculation := HeatSnapshot.GetActiveCirculation()
	perTopicLimit := calcPerTopicLimit(activeCirculation)

	var userTopicStakeTotal int
	sqls.DB().Model(&models.TopicStake{}).
		Where("topic_id = ? AND user_id = ? AND status = ?", req.TopicId, req.UserId, constants.StakeStatusActive).
		Select("COALESCE(SUM(heat_points), 0)").Scan(&userTopicStakeTotal)

	if userTopicStakeTotal+req.HeatPoints > perTopicLimit {
		return nil, fmt.Errorf("对该帖的质押总额已达上限（当前 %d，上限 %d）", userTopicStakeTotal, perTopicLimit)
	}

	// 开启事务
	tx := sqls.DB().Begin()
	defer tx.Rollback()

	// 创建质押记录
	stake := models.TopicStake{
		TopicId:        req.TopicId,
		UserId:         req.UserId,
		HeatPoints:     req.HeatPoints,
		OriginalPoints: req.HeatPoints,
		StakeDay:       time.Now().Format("20060102"),
		Status:         constants.StakeStatusActive,
		CreateTime:     dates.NowTimestamp(),
		UpdateTime:     dates.NowTimestamp(),
	}
	if err := tx.Create(&stake).Error; err != nil {
		return nil, err
	}

	// 扣减余额
	if err := tx.Model(&models.User{}).Where("id = ?", req.UserId).
		UpdateColumn("heat_points", gorm.Expr("heat_points - ?", req.HeatPoints)).Error; err != nil {
		return nil, err
	}

	// 记录流水
	newBalance := user.HeatPoints - req.HeatPoints
	heatLog := models.UserHeatLog{
		UserId:     req.UserId,
		ChangeType: constants.HeatLogTypeStakeOut,
		Amount:     -req.HeatPoints,
		Balance:    newBalance,
		RefId:      fmt.Sprintf("stake_%d", stake.Id),
		Remark:     fmt.Sprintf("质押帖子 ID=%d", req.TopicId),
		CreateTime: dates.NowTimestamp(),
	}
	if err := tx.Create(&heatLog).Error; err != nil {
		return nil, err
	}

	// 更新统计
	var heatStats models.UserHeatStats
	if err := tx.Where("user_id = ?", req.UserId).First(&heatStats).Error; err != nil {
		// 不存在则创建
		heatStats = models.UserHeatStats{
			UserId:      req.UserId,
			TotalPoints: newBalance + req.HeatPoints,
			UpdateTime:  dates.NowTimestamp(),
		}
	}
	heatStats.StakedInWindow += req.HeatPoints
	heatStats.LastStakeTime = dates.NowTimestamp()
	heatStats.TotalPoints = newBalance + s.getActiveStakeTotal(tx, req.UserId)
	heatStats.UpdateTime = dates.NowTimestamp()
	tx.Save(&heatStats)

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 计算火焰等级
	flameLevel := s.CalculateFlameLevel(req.TopicId, activeCirculation)

	resp := &StakeCreateResponse{
		StakeId:        stake.Id,
		RemainingQuota: constants.HeatPointsStakeQuotaDaily - quotaUsed - 1,
		FlameLevel:     flameLevel,
	}

	if flameLevel <= 2 {
		resp.RiskLevel = "high"
		resp.RiskHint = "该帖处于冷帖期，收益波动较大"
	} else if flameLevel <= 3 {
		resp.RiskLevel = "medium"
		resp.RiskHint = "该帖热度正在形成中"
	} else {
		resp.RiskLevel = "low"
		resp.RiskHint = "该帖已达到较高热度"
	}

	return resp, nil
}

// Redeem 赎回质押
func (s *StakeService) Redeem(userId, stakeId int64) error {
	var stake models.TopicStake
	if err := sqls.DB().First(&stake, stakeId).Error; err != nil {
		return fmt.Errorf("质押记录不存在")
	}

	if stake.UserId != userId {
		return fmt.Errorf("无权操作他人的质押")
	}

	if stake.Status != constants.StakeStatusActive {
		return fmt.Errorf("该质押已赎回或已结算")
	}

	// 检查是否当日质押（不可赎回）
	if stake.StakeDay == time.Now().Format("20060102") {
		return fmt.Errorf("当日质押不可赎回，需等到次日结算后")
	}

	tx := sqls.DB().Begin()
	defer tx.Rollback()

	// 更新状态
	if err := tx.Model(&models.TopicStake{}).Where("id = ?", stakeId).
		Update("status", constants.StakeStatusRedeemed).Error; err != nil {
		return err
	}

	// 返还热度点（含利息滚入部分）
	if err := tx.Model(&models.User{}).Where("id = ?", userId).
		UpdateColumn("heat_points", gorm.Expr("heat_points + ?", stake.HeatPoints)).Error; err != nil {
		return err
	}

	// 设置冷却期
	cooldownUntil := dates.NowTimestamp() + int64(constants.HeatPointsCooldownHours)*3600*1000
	var heatStats models.UserHeatStats
	if err := tx.Where("user_id = ?", userId).First(&heatStats).Error; err == nil {
		heatStats.CooldownPoints += stake.HeatPoints
		heatStats.CooldownUntil = cooldownUntil
		heatStats.TotalPoints = s.getTotalPointsTx(tx, userId)
		heatStats.UpdateTime = dates.NowTimestamp()
		tx.Save(&heatStats)
	}

	// 在事务内读取用户余额（修复竞态条件：之前是在事务外读取）
	var user models.User
	if err := tx.First(&user, userId).Error; err != nil {
		return err
	}

	// 记录流水（Balance 使用事务内一致的余额）
	heatLog := models.UserHeatLog{
		UserId:     userId,
		ChangeType: constants.HeatLogTypeRedeem,
		Amount:     stake.HeatPoints,
		Balance:    user.HeatPoints, // 事务内读取，包含刚返还的点数
		RefId:      fmt.Sprintf("stake_%d", stakeId),
		Remark:     fmt.Sprintf("赎回质押，冷却至 %s", time.UnixMilli(cooldownUntil).Format("2006-01-02 15:04")),
		CreateTime: dates.NowTimestamp(),
	}
	if err := tx.Create(&heatLog).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

// GetCooldownPoints 获取冷却中的热度点数（已过期的自动清零）
func (s *StakeService) getCooldownPoints(userId int64) int {
	var stat models.UserHeatStats
	if err := sqls.DB().Where("user_id = ?", userId).First(&stat).Error; err != nil {
		return 0
	}
	if stat.CooldownUntil == 0 || stat.CooldownPoints == 0 {
		return 0
	}
	// 冷却期已过，清零
	if dates.NowTimestamp() >= stat.CooldownUntil {
		sqls.DB().Model(&models.UserHeatStats{}).Where("user_id = ?", userId).
			Updates(map[string]interface{}{
				"cooldown_points": 0,
				"cooldown_until":  0,
			})
		return 0
	}
	return stat.CooldownPoints
}

// GetUserStakes 获取用户质押记录
func (s *StakeService) GetUserStakes(userId int64, status int, limit int) ([]models.TopicStake, error) {
	var stakes []models.TopicStake
	query := sqls.DB().Where("user_id = ?", userId)
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	err := query.Order("id desc").Limit(limit).Find(&stakes).Error
	return stakes, err
}

// GetTodayQuotaUsed 获取今日已用质押次数（含所有状态，避免通过赎回再质押绕过）
func (s *StakeService) GetTodayQuotaUsed(userId int64) int {
	today := time.Now().Format("20060102")
	var count int64
	sqls.DB().Model(&models.TopicStake{}).
		Where("user_id = ? AND stake_day = ?", userId, today).
		Count(&count)
	return int(count)
}

// CalculateFlameLevel 计算帖子火焰等级（含 EverViral 锁定 + 日偏移量 + 滞回）
func (s *StakeService) CalculateFlameLevel(topicId int64, activeCirculation int) int {
	if activeCirculation <= 0 {
		return 0
	}

	var stakeTotal int
	sqls.DB().Model(&models.TopicStake{}).
		Where("topic_id = ? AND status = ?", topicId, constants.StakeStatusActive).
		Select("COALESCE(SUM(heat_points), 0)").Scan(&stakeTotal)

	if stakeTotal == 0 {
		return 0
	}

	ratio := float64(stakeTotal) / float64(activeCirculation)

	// 获取今日偏移量（不存在则使用 1.0）
	offset := s.getFlameOffset()

	// 应用偏移后的阈值
	flame2Thresh := constants.HeatFlameLevel2Threshold * offset.Flame2Offset
	flame3Thresh := constants.HeatFlameLevel3Threshold * offset.Flame3Offset
	flame4Thresh := constants.HeatFlameLevel4Threshold * offset.Flame4Offset
	flame5Thresh := constants.HeatFlameLevel5Threshold * offset.Flame5Offset

	// 计算原始等级
	rawLevel := 1
	if ratio >= flame5Thresh {
		rawLevel = 5
	} else if ratio >= flame4Thresh {
		rawLevel = 4
	} else if ratio >= flame3Thresh {
		rawLevel = 3
	} else if ratio >= flame2Thresh {
		rawLevel = 2
	}

	// 滞回区间：升档需超额 20%，降档允许回落 20%
	var lastLevel int
	s.getFlameLastLevel(topicId, &lastLevel)

	// 如果当前实际等级高于上次记录（升档）：需要超额 20% 才触发
	if rawLevel > lastLevel && lastLevel > 0 {
		// 检查是否真的达到上一级阈值 × 1.2
		var requiredThresh float64
		switch rawLevel {
		case 5:
			requiredThresh = flame5Thresh * 1.2
		case 4:
			requiredThresh = flame4Thresh * 1.2
		case 3:
			requiredThresh = flame3Thresh * 1.2
		case 2:
			requiredThresh = flame2Thresh * 1.2
		}
		if ratio < requiredThresh {
			return lastLevel // 未达到升档要求，保持原等级
		}
	}

	// 如果当前实际等级低于上次记录（降档）：允许 20% 回落空间
	if rawLevel < lastLevel && lastLevel > 0 {
		var holdThresh float64
		switch lastLevel {
		case 5:
			holdThresh = flame5Thresh * 0.8
		case 4:
			holdThresh = flame4Thresh * 0.8
		case 3:
			holdThresh = flame3Thresh * 0.8
		case 2:
			holdThresh = flame2Thresh * 0.8
		}
		if ratio >= holdThresh {
			return lastLevel // 仍在回落区间内，保持原等级
		}
	}

	// 更新缓存的火焰等级
	s.setFlameLastLevel(topicId, rawLevel)
	return rawLevel
}

// getFlameOffset 获取今日火焰偏移量
func (s *StakeService) getFlameOffset() models.DailyFlameOffset {
	today := time.Now().Format("20060102")
	var offset models.DailyFlameOffset
	if err := sqls.DB().Where("date = ?", today).First(&offset).Error; err != nil {
		return models.DailyFlameOffset{
			Flame2Offset: 1.0, Flame3Offset: 1.0,
			Flame4Offset: 1.0, Flame5Offset: 1.0,
		}
	}
	return offset
}

// getFlameLastLevel 从缓存获取上次计算的火焰等级
func (s *StakeService) getFlameLastLevel(topicId int64, level *int) {
	// 从 Topic 表的 flame_locked_level 字段获取（如果管理员手动锁定了则优先）
	var topic struct {
		FlameLockedLevel int
		EverViral        bool
	}
	sqls.DB().Model(&models.Topic{}).Select("flame_locked_level, ever_viral").Where("id = ?", topicId).First(&topic)

	if topic.FlameLockedLevel > 0 {
		*level = topic.FlameLockedLevel
		return
	}

	// 从内存缓存读取上次计算的火焰等级（修复：之前总是返回 0，导致滞回区间永不触发）
	if cached, ok := flameLevelCache.Load(topicId); ok {
		*level = cached.(int)
		return
	}
	*level = 0
}

// setFlameLastLevel 更新火焰等级缓存
func (s *StakeService) setFlameLastLevel(topicId int64, level int) {
	// 修复：之前直接丢弃了值，现在写入 sync.Map 缓存
	flameLevelCache.Store(topicId, level)
}

// getActiveStakeTotal 获取用户当前所有活跃质押的总额
func (s *StakeService) getActiveStakeTotal(tx *gorm.DB, userId int64) int {
	var total int
	tx.Model(&models.TopicStake{}).
		Where("user_id = ? AND status = ?", userId, constants.StakeStatusActive).
		Select("COALESCE(SUM(heat_points), 0)").Scan(&total)
	return total
}

// getTotalPointsTx 在事务中获取用户持有总量（余额 + 活跃质押）
func (s *StakeService) getTotalPointsTx(tx *gorm.DB, userId int64) int {
	var user models.User
	tx.First(&user, userId)
	stakeTotal := s.getActiveStakeTotal(tx, userId)
	return user.HeatPoints + stakeTotal
}

// calcPerTopicLimit 计算单人单帖质押上限（阶梯式：保底 + 百分比）
func calcPerTopicLimit(activeCirculation int) int {
	// 百分比部分
	ratio := int(float64(activeCirculation) * constants.HeatSingleTopicUserLimitRatio)

	// 阶梯保底：社区越小，百分比越低压不住，用固定保底
	switch {
	case activeCirculation < 100:
		return 5 // 极小社区：保底 5 点，够每天 3 次质押
	case activeCirculation < 500:
		return 10 // 小社区：保底 10 点
	case activeCirculation < 2000:
		// 过渡期：取 max(保底, 百分比)
		if ratio < 10 {
			return 10
		}
		return ratio
	default:
		// 大社区：纯百分比
		return ratio
	}
}
