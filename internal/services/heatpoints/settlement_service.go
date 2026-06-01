package heatpoints

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/services"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

// SettlementService 结算服务
type SettlementService struct{}

var Settlement = &SettlementService{}

// stakeYield 单笔质押的收益计算结果
type stakeYield struct {
	StakeId   int64
	UserId    int64
	TopicId   int64
	Phase     int
	DailyRate float64
	Yield     int // 正=收益, 负=亏损
}

// SettleAll 执行每日结算（00:00 调用）
func (s *SettlementService) SettleAll() error {
	today := time.Now().Format("20060102")
	slog.Info("开始执行热度点每日结算", "date", today)

	// 0. 提升待生效的热度点参数（heat.* 前缀）
	//    必须在任何参数读取之前执行，确保本次结算使用最新参数
	if err := services.SysConfigService.PromotePendingParams(); err != nil {
		slog.Warn("参数提升失败，使用旧参数继续结算", "error", err)
	}

	// 1. 幂等性检查
	if s.isAlreadySettled(today) {
		slog.Info("今日已结算，跳过", "date", today)
		return nil
	}
	s.markSettlementStart(today)

	// 2. 检查昨日快照是否存在
	yesterday := time.Now().AddDate(0, 0, -1).Format("20060102")
	var circulationSnapshot models.HeatCirculationSnapshot
	if err := sqls.DB().Where("snapshot_date = ?", yesterday).First(&circulationSnapshot).Error; err != nil {
		slog.Error("昨日流通快照不存在，跳过结算", "date", yesterday)
		s.markSettlementFailed(today, "昨日流通快照不存在: "+yesterday)
		return nil
	}

	// 3. Phase 1: 分批扫描所有活跃质押，计算收益（不写入数据库）
	allYields, err := s.calculateAllYields(yesterday, &circulationSnapshot)
	if err != nil {
		slog.Error("收益计算失败", "error", err)
		s.markSettlementFailed(today, err.Error())
		return err
	}

	// 4. Phase 2: 分离正收益和负收益
	totalPositive := 0
	totalNegative := 0
	for _, y := range allYields {
		if y.Yield > 0 {
			totalPositive += y.Yield
		} else if y.Yield < 0 {
			totalNegative += -y.Yield // 转为正数
		}
	}

	slog.Info("收益计算完成", "stakes", len(allYields), "positive", totalPositive, "negative", totalNegative)

	// 5. 获取公共奖池当前余额
	poolBalance := s.getPoolBalance()

	// 6. 注入负收益到奖池（在正收益发放之前，避免低估支付能力）
	if totalNegative > 0 {
		poolBalance += totalNegative
		s.recordPoolFlow(constants.HeatPoolSourceStakeLoss, totalNegative, poolBalance,
			fmt.Sprintf("结算注入负收益，共 %d 笔", totalNegative))
	}

	// 7. 第一优先级：签到奖励（从奖池扣除，必须全额支付）
	if err := s.distributeCheckInRewards(&poolBalance, &circulationSnapshot); err != nil {
		slog.Error("签到发放失败", "error", err)
	}

	// 8. 第二优先级：正收益发放（带准备金约束）
	proRataRatio := 1.0
	if totalPositive > 0 && poolBalance < totalPositive {
		proRataRatio = float64(poolBalance) / float64(totalPositive)
		slog.Warn("准备金不足，等比削减正收益", "pool", poolBalance, "needed", totalPositive, "ratio", proRataRatio)
	}

	// 9. Phase 3: 写入质押记录和流水
	if err := s.applyYields(allYields, proRataRatio, yesterday); err != nil {
		slog.Error("收益应用失败", "error", err)
		s.markSettlementFailed(today, err.Error())
		return err
	}

	// 更新奖池余额（扣除正收益发放）
	actualPositivePaid := int(float64(totalPositive) * proRataRatio)
	if actualPositivePaid > 0 {
		poolBalance -= actualPositivePaid
		s.recordPoolFlow(constants.HeatPoolSourceSettlePayout, -actualPositivePaid, poolBalance,
			fmt.Sprintf("结算正收益发放，削减率=%.2f%%", (1-proRataRatio)*100))
	}

	// 10. Phase 4: 执行衰减
	truncatedTotal := 0
	if err := s.applyDecay(&truncatedTotal); err != nil {
		slog.Error("衰减执行失败", "error", err)
	}

	// 11. Phase 5: 标记 EverViral
	s.markEverViral(yesterday, &circulationSnapshot)

	// 12. Phase 6: 60 天完全归零
	s.applyFullForfeit(&poolBalance)

	// 13. Phase 7: 周利用率排名奖励（周一执行）
	s.distributeRankRewards(&poolBalance)

	// 14. 完成标记
	s.markSettlementCompleted(today, len(allYields), actualPositivePaid, totalNegative)

	slog.Info("结算完成",
		"stakes", len(allYields),
		"positive_paid", actualPositivePaid,
		"negative_injected", totalNegative,
		"pro_rata_ratio", proRataRatio,
		"truncated_decay", truncatedTotal,
		"pool_balance", poolBalance,
	)
	return nil
}

// calculateAllYields 分批计算所有活跃质押的收益（游标式扫描，不全量加载）
func (s *SettlementService) calculateAllYields(yesterday string, snapshot *models.HeatCirculationSnapshot) ([]stakeYield, error) {
	var allYields []stakeYield

	// 预加载所有达到 EverViral 的帖子 ID（一次查询，避免逐条查询）
	everViralMap := make(map[int64]bool)
	var viralIds []int64
	sqls.DB().Model(&models.Topic{}).Where("ever_viral = ?", true).Pluck("id", &viralIds)
	for _, id := range viralIds {
		everViralMap[id] = true
	}

	// 按 topic_id 范围分批（每批 500 条）
	var lastId int64 = 0
	for {
		var batch []models.TopicStake
		if err := sqls.DB().
			Where("status = ? AND id > ?", constants.StakeStatusActive, lastId).
			Order("id asc").
			Limit(500).
			Find(&batch).Error; err != nil {
			return nil, err
		}
		if len(batch) == 0 {
			break
		}

		for _, stake := range batch {
			lastId = stake.Id

			// 跳过已结算过的（幂等保护）
			if stake.LastSettleDay == yesterday {
				continue
			}

			// 获取帖子快照
			var topicSnapshot models.TopicInteractionSnapshot
			if err := sqls.DB().Where("topic_id = ? AND snapshot_date = ?", stake.TopicId, yesterday).
				First(&topicSnapshot).Error; err != nil {
				continue // 快照不存在，跳过
			}

			everViral := everViralMap[stake.TopicId]
			phase := s.calculatePhase(topicSnapshot.StakeTotal, snapshot.ActiveCirculation, everViral)
			dailyRate := s.calculateDailyRate(&topicSnapshot, yesterday)
			yield := s.calculateYield(stake.HeatPoints, dailyRate, phase)

			allYields = append(allYields, stakeYield{
				StakeId:   stake.Id,
				UserId:    stake.UserId,
				TopicId:   stake.TopicId,
				Phase:     phase,
				DailyRate: dailyRate,
				Yield:     yield,
			})
		}

		if len(batch) < 500 {
			break
		}
	}

	return allYields, nil
}

// applyYields 将计算好的收益写入质押记录和流水
func (s *SettlementService) applyYields(yields []stakeYield, proRataRatio float64, yesterday string) error {
	batchSize := 200
	for i := 0; i < len(yields); i += batchSize {
		end := i + batchSize
		if end > len(yields) {
			end = len(yields)
		}
		batch := yields[i:end]

		if err := sqls.DB().Transaction(func(tx *gorm.DB) error {
			for _, y := range batch {
				actualYield := y.Yield

				// 正收益按比例削减
				if y.Yield > 0 {
					actualYield = int(float64(y.Yield) * proRataRatio)
				}

				// 只读一次 stake（修复重复读取）
				var stake models.TopicStake
				if err := tx.First(&stake, y.StakeId).Error; err != nil {
					continue
				}

				// 亏损归零保护：不会超过本金
				if actualYield < 0 {
					if -actualYield > stake.HeatPoints {
						actualYield = -stake.HeatPoints
					}
				}

				newHeatPoints := stake.HeatPoints + actualYield

				if newHeatPoints <= 0 {
					// 本金归零，标记为失败
					tx.Model(&models.TopicStake{}).Where("id = ?", y.StakeId).Updates(map[string]interface{}{
						"status":        constants.StakeStatusFailed,
						"heat_points":   0,
						"update_time":   dates.NowTimestamp(),
						"last_settle_day": yesterday,
					})
				} else {
					// 更新质押金额（利息滚入本金）
					tx.Model(&models.TopicStake{}).Where("id = ?", y.StakeId).Updates(map[string]interface{}{
						"heat_points":     newHeatPoints,
						"update_time":     dates.NowTimestamp(),
						"last_settle_day": yesterday,
					})
				}

				// 获取用户实际余额（在事务内读取，避免竞态）
				var user models.User
				tx.First(&user, y.UserId)

				// 记录流水（Balance 使用用户实际余额，非质押金额）
				changeType := constants.HeatLogTypeSettleProfit
				if actualYield < 0 {
					changeType = constants.HeatLogTypeSettleLoss
				} else if actualYield < y.Yield {
					changeType = constants.HeatLogTypeSettleProfitPartial
				}

				remark := fmt.Sprintf("每日结算，阶段=%d, 利率=%.4f", y.Phase, y.DailyRate)
				if actualYield != y.Yield && y.Yield > 0 {
					remark += fmt.Sprintf(", 理论=%d, 实际=%d, 削减率=%.1f%%", y.Yield, actualYield, (1-proRataRatio)*100)
				}

				heatLog := models.UserHeatLog{
					UserId:     y.UserId,
					ChangeType: changeType,
					Amount:     actualYield,
					Balance:    user.HeatPoints, // 用户实际余额（事务内一致读取）
					RefId:      fmt.Sprintf("stake_%d", y.StakeId),
					Remark:     remark,
					CreateTime: dates.NowTimestamp(),
				}
				tx.Create(&heatLog)

				// 同步 UserHeatStats.TotalPoints（余额 + 活跃质押）
				s.syncTotalPoints(tx, y.UserId)
			}
			return nil
		}); err != nil {
			slog.Error("批次收益应用失败", "error", err, "batch_start", i)
		}
	}
	return nil
}

// calculatePhase 计算帖子所处阶段（考虑 EverViral 锁定）
func (s *SettlementService) calculatePhase(stakeTotal, activeCirculation int, everViral bool) int {
	if activeCirculation <= 0 {
		return 1
	}
	ratio := float64(stakeTotal) / float64(activeCirculation)

	var phase int
	if ratio < constants.HeatPhase1Threshold {
		phase = 1 // 冷帖期
	} else if ratio < constants.HeatPhase2Threshold {
		phase = 2 // 共识期
	} else {
		phase = 3 // 热门期（断崖）
	}

	// EverViral 帖子永远不能回到阶段一（防止套利循环）
	if everViral && phase < 2 {
		phase = 2
	}
	return phase
}

// calculateDailyRate 计算日利率（使用 log2 平滑 + 去重用户 + 生命周期系数）
func (s *SettlementService) calculateDailyRate(snapshot *models.TopicInteractionSnapshot, yesterday string) float64 {
	// 获取前日快照
	prevDate := time.Now().AddDate(0, 0, -2).Format("20060102")
	var prevSnapshot models.TopicInteractionSnapshot
	sqls.DB().Where("topic_id = ? AND snapshot_date = ?", snapshot.TopicId, prevDate).First(&prevSnapshot)

	// 使用去重后的互动得分：log2(1 + uniqueLikers*3 + validComments*2)
	// 对数平滑抑制小基数暴增，去重用户数代替原始计数防止刷量
	currScore := math.Log2(1 + float64(snapshot.UniqueLikers)*3 + float64(snapshot.ValidComments)*2)
	prevScore := math.Log2(1 + float64(prevSnapshot.UniqueLikers)*3 + float64(prevSnapshot.ValidComments)*2)

	// 互动增长率
	var interactionGrowth float64
	if prevScore > 0.01 {
		interactionGrowth = (currScore - prevScore) / prevScore
	} else if currScore > 0 {
		// 新帖首日：封顶处理
		interactionGrowth = math.Min(currScore/5.0, 1.0)
	}

	// 生命周期系数：发帖越久增长权重越低
	// 1/(1 + daysOld/30)，第 1 天 ≈ 0.97，第 30 天 = 0.5，第 90 天 ≈ 0.25
	daysOld := s.getTopicDaysOld(snapshot.TopicId)
	lifecycleCoeff := 1.0 / (1.0 + float64(daysOld)/30.0)

	interactionGrowth *= lifecycleCoeff

	// tanh 饱和，防止暴增
	interactionGrowth = math.Tanh(interactionGrowth)

	return interactionGrowth
}

// getTopicDaysOld 获取帖子发帖天数（最小为 1）
func (s *SettlementService) getTopicDaysOld(topicId int64) int {
	var topic models.Topic
	if err := sqls.DB().Select("create_time").First(&topic, topicId).Error; err != nil {
		return 1
	}
	days := int(time.Since(time.UnixMilli(topic.CreateTime)).Hours() / 24)
	if days < 1 {
		return 1
	}
	return days
}

// calculateYield 计算收益（风险系数只放大正收益）
func (s *SettlementService) calculateYield(principal int, dailyRate float64, phase int) int {
	var riskCoeff, rateCap, rateFloor float64

	switch phase {
	case 1:
		riskCoeff = 2.0
		rateCap = 0.50
		rateFloor = -0.30
	case 2:
		riskCoeff = 1.0
		rateCap = 0.20
		rateFloor = -0.20
	case 3:
		riskCoeff = 0.1
		rateCap = 0.02
		rateFloor = -0.30
	}

	// 限制利率范围
	if dailyRate > rateCap {
		dailyRate = rateCap
	} else if dailyRate < rateFloor {
		dailyRate = rateFloor
	}

	// 风险系数只放大正收益，亏损原值扣除
	var yield float64
	if dailyRate > 0 {
		yield = float64(principal) * dailyRate * riskCoeff
	} else {
		yield = float64(principal) * dailyRate
	}

	return int(yield)
}

// applyDecay 执行衰减（记录截断量）
func (s *SettlementService) applyDecay(truncatedTotal *int) error {
	var stats []models.UserHeatStats
	sqls.DB().Find(&stats)

	today := time.Now().Format("20060102")

	for _, stat := range stats {
		// 计算未使用量
		unused := stat.TotalPoints - stat.StakedInWindow
		if stat.StakedInWindow >= stat.TotalPoints {
			unused = 0
		}

		if unused <= 0 {
			continue
		}

		decayAmount := int(float64(unused) * constants.HeatPointsDecayRate)
		if decayAmount <= 0 {
			continue
		}

		// 获取用户余额
		var user models.User
		if err := sqls.DB().First(&user, stat.UserId).Error; err != nil {
			continue
		}

		// 截断：只扣余额，不产生负债
		actualDecay := decayAmount
		truncated := 0
		if user.HeatPoints < decayAmount {
			actualDecay = user.HeatPoints
			truncated = decayAmount - actualDecay
		}

		if actualDecay > 0 {
			sqls.DB().Model(&models.User{}).Where("id = ?", stat.UserId).
				UpdateColumn("heat_points", gorm.Expr("heat_points - ?", actualDecay))

			heatLog := models.UserHeatLog{
				UserId:     stat.UserId,
				ChangeType: constants.HeatLogTypeDecay,
				Amount:     -actualDecay,
				Balance:    user.HeatPoints - actualDecay,
				RefId:      "decay_" + today,
				Remark:     fmt.Sprintf("每日衰减，理论=%d, 实际=%d", decayAmount, actualDecay),
				CreateTime: dates.NowTimestamp(),
			}
			sqls.DB().Create(&heatLog)

			// 实际扣除部分流入公共奖池
			s.recordPoolFlow(constants.HeatPoolSourceDecayInflow, actualDecay, s.getPoolBalance()+actualDecay, "衰减回收")
		}

		// 记录截断量
		if truncated > 0 {
			*truncatedTotal += truncated

			heatLog := models.UserHeatLog{
				UserId:     stat.UserId,
				ChangeType: constants.HeatLogTypeDecayTruncated,
				Amount:     0, // 不发生实际变动
				Balance:    0, // 余额已为 0
				RefId:      "trunc_" + today,
				Remark:     fmt.Sprintf("衰减截断，丢弃=%d", truncated),
				CreateTime: dates.NowTimestamp(),
			}
			sqls.DB().Create(&heatLog)
		}

		// 更新统计表的衰减累计
		sqls.DB().Model(&models.UserHeatStats{}).Where("user_id = ?", stat.UserId).
			UpdateColumn("decayed_accumulated", gorm.Expr("decayed_accumulated + ?", actualDecay))

		// 同步 TotalPoints（余额因衰减减少，TotalPoints 需要同步）
		s.syncTotalPoints(sqls.DB(), stat.UserId)
	}

	// 更新流通快照的截断量
	if *truncatedTotal > 0 {
		sqls.DB().Model(&models.HeatCirculationSnapshot{}).
			Where("snapshot_date = ?", today).
			Update("daily_decay_truncated", *truncatedTotal)
	}

	return nil
}

// markEverViral 标记达到阶段三的帖子为 EverViral（不可逆）
func (s *SettlementService) markEverViral(yesterday string, snapshot *models.HeatCirculationSnapshot) {
	// 找出所有质押量 ≥ 阶段三阈值的帖子
	threshold := int(float64(snapshot.ActiveCirculation) * constants.HeatPhase2Threshold)

	var topicIds []int64
	sqls.DB().Model(&models.TopicStake{}).
		Where("status = ?", constants.StakeStatusActive).
		Group("topic_id").
		Having("SUM(heat_points) >= ?", threshold).
		Pluck("topic_id", &topicIds)

	if len(topicIds) > 0 {
		sqls.DB().Model(&models.Topic{}).
			Where("id IN ? AND ever_viral = ?", topicIds, false).
			Update("ever_viral", true)
		slog.Info("标记 EverViral", "count", len(topicIds))
	}
}

// distributeCheckInRewards 发放签到奖励（从奖池扣除，第一优先级）
func (s *SettlementService) distributeCheckInRewards(poolBalance *int, snapshot *models.HeatCirculationSnapshot) error {
	// 查找今日签到的用户
	today := dates.GetDay(time.Now())

	var checkIns []models.CheckIn
	sqls.DB().Where("latest_day_name = ?", today).Find(&checkIns)

	if len(checkIns) == 0 {
		return nil
	}

	// 直接使用已加载的快照，避免重复查询
	activeCirculation := snapshot.ActiveCirculation
	if activeCirculation <= 0 {
		activeCirculation = constants.HeatPointsDefaultCirculation
	}
	baseReward := int(float64(activeCirculation) * constants.HeatPointsDailyCheckinRate)
	if baseReward <= 0 {
		baseReward = 1 // 最低保底
	}

	totalNeeded := len(checkIns) * baseReward

	// 熔断检查：奖池余额不足时等比削减
	actualReward := baseReward
	if totalNeeded > *poolBalance {
		if *poolBalance <= 0 {
			slog.Warn("奖池余额为零，跳过签到发放")
			return nil
		}
		actualReward = *poolBalance / len(checkIns)
		if actualReward <= 0 {
			return nil
		}
		slog.Warn("签到熔断", "needed", totalNeeded, "pool", *poolBalance, "actual_per_user", actualReward)
	}

	actualTotal := 0
	for _, ci := range checkIns {
		if actualReward <= 0 {
			break
		}

		// 连续签到加成
		multiplier := 1.0
		if ci.ConsecutiveDays >= 30 {
			multiplier = 2.0
		} else if ci.ConsecutiveDays >= 7 {
			multiplier = 1.5
		} else if ci.ConsecutiveDays >= 3 {
			multiplier = 1.2
		}
		reward := int(float64(actualReward) * multiplier)

		if reward > *poolBalance-actualTotal {
			reward = *poolBalance - actualTotal
		}
		if reward <= 0 {
			break
		}

		// 发放到用户
		sqls.DB().Model(&models.User{}).Where("id = ?", ci.UserId).
			UpdateColumn("heat_points", gorm.Expr("heat_points + ?", reward))

		// 更新统计
		sqls.DB().Model(&models.UserHeatStats{}).Where("user_id = ?", ci.UserId).
			UpdateColumn("total_points", gorm.Expr("total_points + ?", reward))

		// 记录流水
		var user models.User
		sqls.DB().First(&user, ci.UserId)
		heatLog := models.UserHeatLog{
			UserId:     ci.UserId,
			ChangeType: constants.HeatLogTypeDailyCheckIn,
			Amount:     reward,
			Balance:    user.HeatPoints,
			RefId:      fmt.Sprintf("checkin_%d_%s", ci.Id, time.Now().Format("20060102")),
			Remark:     fmt.Sprintf("连续签到第%d天，系数=%.1f", ci.ConsecutiveDays, multiplier),
			CreateTime: dates.NowTimestamp(),
		}
		sqls.DB().Create(&heatLog)

		actualTotal += reward
	}

	if actualTotal > 0 {
		*poolBalance -= actualTotal
		s.recordPoolFlow(constants.HeatPoolSourceCheckInPayout, -actualTotal, *poolBalance,
			fmt.Sprintf("签到发放，%d人, 每人%d点", len(checkIns), actualReward))
	}

	return nil
}

// getPoolBalance 获取公共奖池当前余额
func (s *SettlementService) getPoolBalance() int {
	var latest models.HeatPublicPool
	if err := sqls.DB().Order("id desc").First(&latest).Error; err != nil {
		return 0
	}
	return latest.BalanceAfter
}

// recordPoolFlow 记录公共奖池流水
func (s *SettlementService) recordPoolFlow(source string, amount int, balanceAfter int, remark string) {
	pool := models.HeatPublicPool{
		Source:       source,
		Amount:       amount,
		BalanceAfter: balanceAfter,
		Remark:       remark,
		CreateTime:   dates.NowTimestamp(),
	}
	sqls.DB().Create(&pool)
}

// isAlreadySettled 幂等性检查
func (s *SettlementService) isAlreadySettled(today string) bool {
	var log models.SettlementTaskLog
	err := sqls.DB().Where("task_date = ? AND task_type = ?", today, constants.SettlementTaskTypeSettlement).First(&log).Error
	return err == nil && log.Status == constants.SettlementTaskStatusCompleted
}

func (s *SettlementService) markSettlementStart(today string) {
	log := models.SettlementTaskLog{
		TaskDate:  today,
		TaskType:  constants.SettlementTaskTypeSettlement,
		Status:    constants.SettlementTaskStatusRunning,
		StartedAt: dates.NowTimestamp(),
	}
	sqls.DB().Create(&log)
}

func (s *SettlementService) markSettlementCompleted(today string, totalProcessed, positive, negative int) {
	sqls.DB().Model(&models.SettlementTaskLog{}).
		Where("task_date = ? AND task_type = ?", today, constants.SettlementTaskTypeSettlement).
		Updates(map[string]interface{}{
			"status":          constants.SettlementTaskStatusCompleted,
			"finished_at":     dates.NowTimestamp(),
			"total_processed": totalProcessed,
			"batch_count":     positive + negative,
		})
}

func (s *SettlementService) markSettlementFailed(today string, errMsg string) {
	sqls.DB().Model(&models.SettlementTaskLog{}).
		Where("task_date = ? AND task_type = ?", today, constants.SettlementTaskTypeSettlement).
		Updates(map[string]interface{}{
			"status":      constants.SettlementTaskStatusFailed,
			"finished_at": dates.NowTimestamp(),
			"error_msg":   errMsg,
		})
}

// syncTotalPoints 同步 UserHeatStats.TotalPoints（余额 + 活跃质押）
func (s *SettlementService) syncTotalPoints(tx *gorm.DB, userId int64) {
	var user models.User
	if err := tx.First(&user, userId).Error; err != nil {
		return
	}
	var stakeTotal int
	tx.Model(&models.TopicStake{}).
		Where("user_id = ? AND status = ?", userId, constants.StakeStatusActive).
		Select("COALESCE(SUM(heat_points), 0)").Scan(&stakeTotal)

	newTotal := user.HeatPoints + stakeTotal
	tx.Model(&models.UserHeatStats{}).Where("user_id = ?", userId).
		Update("total_points", newTotal)
}

// applyFullForfeit 60 天完全未使用归零
func (s *SettlementService) applyFullForfeit(poolBalance *int) {
	today := time.Now().Format("20060102")
	cutoffMs := dates.Timestamp(time.Now().AddDate(0, 0, -constants.HeatPointsForfeitDays))

	// 找出连续 60 天未质押的用户（LastStakeTime > 0 且 < cutoff，排除从未质押的用户）
	var forfeitStats []models.UserHeatStats
	sqls.DB().Where("last_stake_time > 0 AND last_stake_time < ?", cutoffMs).
		Find(&forfeitStats)

	for _, stat := range forfeitStats {
		var user models.User
		if err := sqls.DB().First(&user, stat.UserId).Error; err != nil {
			continue
		}

		// 仅对"可用余额"归零，不动质押中的点数
		if user.HeatPoints <= 0 {
			continue
		}

		forfeitAmount := user.HeatPoints

		// 扣减余额
		sqls.DB().Model(&models.User{}).Where("id = ?", stat.UserId).
			UpdateColumn("heat_points", 0)

		// 流入公共奖池
		*poolBalance += forfeitAmount
		s.recordPoolFlow(constants.HeatPoolSourceFullForfeit, forfeitAmount, *poolBalance,
			fmt.Sprintf("60 天完全归零，用户=%d", stat.UserId))

		// 记录流水
		heatLog := models.UserHeatLog{
			UserId:     stat.UserId,
			ChangeType: constants.HeatLogTypeFullForfeit,
			Amount:     -forfeitAmount,
			Balance:    0,
			RefId:      "forfeit_" + today,
			Remark:     fmt.Sprintf("连续 %d 天未质押，余额归零入池", constants.HeatPointsForfeitDays),
			CreateTime: dates.NowTimestamp(),
		}
		sqls.DB().Create(&heatLog)

		// 同步统计
		s.syncTotalPoints(sqls.DB(), stat.UserId)
	}
}

// distributeRankRewards 每周利用率排名奖励（周一执行）
func (s *SettlementService) distributeRankRewards(poolBalance *int) {
	if time.Now().Weekday() != time.Monday {
		return
	}
	if *poolBalance <= 0 {
		return
	}

	// 奖励池 = 奖池余额的 50%
	rewardPool := *poolBalance / 2
	if rewardPool <= 0 {
		return
	}

	// 找出所有利用率 > 0 的用户
	type utilizationRow struct {
		UserId         int64
		StakedInWindow int
		TotalPoints    int
	}
	var rows []utilizationRow
	sqls.DB().Model(&models.UserHeatStats{}).
		Where("staked_in_window > 0 AND total_points > 0").
		Select("user_id, staked_in_window, total_points").
		Order("CAST(staked_in_window AS FLOAT) / CAST(total_points AS FLOAT) DESC").
		Find(&rows)

	if len(rows) == 0 {
		return
	}

	// 分段：前 5% 平分 50%，前 5-10% 平分 30%，前 10-20% 平分 20%
	perCap := 50 // 单人上限
	segments := []struct {
		startPct float64
		endPct   float64
		share    float64
	}{
		{0.0, 0.05, 0.50},
		{0.05, 0.10, 0.30},
		{0.10, 0.20, 0.20},
	}

	totalDistributed := 0
	today := time.Now().Format("20060102")

	for _, seg := range segments {
		startIdx := int(float64(len(rows)) * seg.startPct)
		endIdx := int(float64(len(rows)) * seg.endPct)
		if startIdx >= endIdx || startIdx >= len(rows) {
			continue
		}
		if endIdx > len(rows) {
			endIdx = len(rows)
		}

		segPool := int(float64(rewardPool) * seg.share)
		count := endIdx - startIdx
		if count <= 0 {
			continue
		}

		perPerson := segPool / count
		if perPerson > perCap {
			perPerson = perCap
		}
		if perPerson <= 0 {
			continue
		}

		for i := startIdx; i < endIdx; i++ {
			if totalDistributed+perPerson > rewardPool {
				break
			}

			userId := rows[i].UserId
			sqls.DB().Model(&models.User{}).Where("id = ?", userId).
				UpdateColumn("heat_points", gorm.Expr("heat_points + ?", perPerson))
			s.syncTotalPoints(sqls.DB(), userId)

			heatLog := models.UserHeatLog{
				UserId:     userId,
				ChangeType: constants.HeatLogTypeRankReward,
				Amount:     perPerson,
				Balance:    0, // 余额在 syncTotalPoints 中更新
				RefId:      "rank_" + today,
				Remark:     fmt.Sprintf("周利用率排名奖励，段=%.0f%%-%.0f%%", seg.startPct*100, seg.endPct*100),
				CreateTime: dates.NowTimestamp(),
			}
			sqls.DB().Create(&heatLog)

			totalDistributed += perPerson
		}
	}

	if totalDistributed > 0 {
		*poolBalance -= totalDistributed
		s.recordPoolFlow(constants.HeatPoolSourceRankReward, -totalDistributed, *poolBalance,
			fmt.Sprintf("周排名奖励发放，%d 人，共 %d 点", len(rows), totalDistributed))
	}
}
