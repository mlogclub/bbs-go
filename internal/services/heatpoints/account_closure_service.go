package heatpoints

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"fmt"
	"log/slog"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

// AccountClosureService 用户注销时的热度点资产处理
type AccountClosureService struct{}

var AccountClosure = &AccountClosureService{}

// CloseAccount 处理用户注销：强制赎回所有质押 → 余额归入公共奖池 → 记录审计日志
func (s *AccountClosureService) CloseAccount(userId int64, operatorId int64, reason string) error {
	slog.Info("处理用户注销热度点资产", "userId", userId, "operatorId", operatorId)

	// 事务内处理
	return sqls.DB().Transaction(func(tx *gorm.DB) error {

		// 1. 强制赎回该用户所有活跃质押（返还 OriginalPoints，不赚不赔）
		var stakes []models.TopicStake
		tx.Where("user_id = ? AND status = ?", userId, constants.StakeStatusActive).Find(&stakes)

		totalRedeemed := 0
		for _, stake := range stakes {
			returnAmount := stake.OriginalPoints
			if returnAmount > stake.HeatPoints {
				returnAmount = stake.HeatPoints
			}

			tx.Model(&models.TopicStake{}).Where("id = ?", stake.Id).Updates(map[string]interface{}{
				"status":      constants.StakeStatusAdminRedeemed,
				"heat_points": returnAmount,
				"update_time": dates.NowTimestamp(),
			})

			totalRedeemed += returnAmount

			tx.Create(&models.UserHeatLog{
				UserId:     userId,
				ChangeType: constants.HeatLogTypeAdminForceRedeem,
				Amount:     returnAmount,
				Balance:    0,
				RefId:      fmt.Sprintf("stake_%d_closure", stake.Id),
				Remark:     fmt.Sprintf("注销强制赎回，帖子ID=%d", stake.TopicId),
				CreateTime: dates.NowTimestamp(),
			})
		}

		// 2. 将用户当前余额（含刚赎回的部分）转入公共奖池
		var user models.User
		if err := tx.First(&user, userId).Error; err != nil {
			return err
		}

		if user.HeatPoints > 0 {
			// 获取奖池当前余额
			var latestPool models.HeatPublicPool
			tx.Order("id desc").First(&latestPool)
			newPoolBalance := latestPool.BalanceAfter + user.HeatPoints

			tx.Create(&models.HeatPublicPool{
				Source:       constants.HeatPoolSourceAccountClosure,
				Amount:       user.HeatPoints,
				BalanceAfter: newPoolBalance,
				RefId:        fmt.Sprintf("user_%d", userId),
				Remark:       fmt.Sprintf("用户注销，回收 %d 热度点", user.HeatPoints),
				CreateTime:   dates.NowTimestamp(),
			})

			// 记录流水
			tx.Create(&models.UserHeatLog{
				UserId:     userId,
				ChangeType: constants.HeatLogTypeAccountClosure,
				Amount:     -user.HeatPoints,
				Balance:    0,
				RefId:      "closure",
				Remark:     fmt.Sprintf("注销回收，%d 热度点归入公共奖池", user.HeatPoints),
				CreateTime: dates.NowTimestamp(),
			})

			// 清空用户余额
			tx.Model(&models.User{}).Where("id = ?", userId).
				Update("heat_points", 0)
		}

		// 3. 清零统计记录
		tx.Model(&models.UserHeatStats{}).Where("user_id = ?", userId).
			Updates(map[string]interface{}{
				"total_points":        0,
				"staked_in_window":    0,
				"cooldown_points":     0,
				"cooldown_until":      0,
				"decayed_accumulated": 0,
			})

		// 4. 失效缓存
		cache.UserCache.Invalidate(userId)

		// 5. 写入操作日志
		tx.Create(&models.OperateLog{
			UserId:      operatorId,
			OpType:      "account_closure",
			DataType:    "user",
			DataId:      userId,
			Description: fmt.Sprintf("用户注销热度点资产处理：强制赎回 %d 笔质押，回收 %d 热度点入池。原因：%s", len(stakes), user.HeatPoints, reason),
			CreateTime:  dates.NowTimestamp(),
		})

		slog.Info("用户注销热度点处理完成",
			"userId", userId,
			"stakes_redeemed", len(stakes),
			"total_redeemed", totalRedeemed,
			"balance_to_pool", user.HeatPoints,
		)
		return nil
	})
}
