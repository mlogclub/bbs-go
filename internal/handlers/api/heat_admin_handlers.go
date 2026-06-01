package api

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/services/heatpoints"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

// AdminHeatForceRedeem 管理员强制赎回某帖子所有质押
func AdminHeatForceRedeem(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("未登录"))
		return
	}
	// 仅 Owner 角色可操作
	if !user.IsOwner() {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("无权限"))
		return
	}

	var req struct {
		TopicId int64  `json:"topicId" binding:"required"`
		Reason  string `json:"reason" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("参数错误：需要 topicId 和 reason"))
		return
	}
	if len(req.Reason) < 8 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("原因说明不少于 8 字"))
		return
	}

	// 查询该帖子所有活跃质押
	var stakes []models.TopicStake
	sqls.DB().Where("topic_id = ? AND status = ?", req.TopicId, constants.StakeStatusActive).Find(&stakes)

	if len(stakes) == 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("该帖子无活跃质押"))
		return
	}

	tx := sqls.DB().Begin()
	defer tx.Rollback()

	totalReturned := 0
	redeemedCount := 0

	for _, stake := range stakes {
		// 返还 OriginalPoints（撤销参与，不赚不赔）
		returnAmount := stake.OriginalPoints
		if returnAmount > stake.HeatPoints {
			returnAmount = stake.HeatPoints // 不会超过当前值
		}

		// 更新质押状态
		tx.Model(&models.TopicStake{}).Where("id = ?", stake.Id).Updates(map[string]interface{}{
			"status":      constants.StakeStatusAdminRedeemed,
			"heat_points": returnAmount,
			"update_time": dates.NowTimestamp(),
		})

		// 返还到用户余额
		tx.Model(&models.User{}).Where("id = ?", stake.UserId).
			UpdateColumn("heat_points", gorm.Expr("heat_points + ?", returnAmount))

		// 记录流水
		heatLog := models.UserHeatLog{
			UserId:     stake.UserId,
			ChangeType: constants.HeatLogTypeAdminForceRedeem,
			Amount:     returnAmount,
			Balance:    0,
			RefId:      fmt.Sprintf("stake_%d", stake.Id),
			Remark:     fmt.Sprintf("管理员强制赎回：帖子 ID=%d, 原因=%s", req.TopicId, req.Reason),
			CreateTime: dates.NowTimestamp(),
		}
		tx.Create(&heatLog)

		totalReturned += returnAmount
		redeemedCount++
	}

	// 写入操作日志
	sqls.DB().Create(&models.OperateLog{
		UserId:      user.Id,
		OpType:      "force_redeem",
		DataType:    "topic",
		DataId:      req.TopicId,
		Description: fmt.Sprintf("强制赎回质押 %d 笔，返还 %d 热度点。原因：%s", redeemedCount, totalReturned, req.Reason),
		CreateTime:  dates.NowTimestamp(),
	})

	if err := tx.Commit().Error; err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("事务提交失败：" + err.Error()))
		return
	}

	ginx.WriteJSON(ctx, map[string]interface{}{
		"redeemedCount":      redeemedCount,
		"totalPointsReturned": totalReturned,
	})
}

// AdminHeatFreezeFlame 管理员锁定/解锁火焰等级
func AdminHeatFreezeFlame(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("未登录"))
		return
	}
	if !user.IsOwner() {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("无权限"))
		return
	}

	var req struct {
		TopicId    int64 `json:"topicId" binding:"required"`
		FlameLevel int   `json:"flameLevel" binding:"required"` // 0=解锁，1-5=锁定
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("参数错误"))
		return
	}
	if req.FlameLevel < 0 || req.FlameLevel > 5 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("flameLevel 必须在 0-5 之间"))
		return
	}

	// 锁定数量上限检查（最多同时锁定 20 个帖子）
	if req.FlameLevel > 0 {
		var lockedCount int64
		sqls.DB().Model(&models.Topic{}).Where("flame_locked_level > 0").Count(&lockedCount)
		if lockedCount >= 20 {
			ginx.WriteJSON(ctx, ginx.ErrorMessage("同时锁定帖子数已达上限（20）"))
			return
		}
	}

	// 更新
	if err := sqls.DB().Model(&models.Topic{}).Where("id = ?", req.TopicId).
		Update("flame_locked_level", req.FlameLevel).Error; err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("更新失败：" + err.Error()))
		return
	}

	action := "锁定"
	if req.FlameLevel == 0 {
		action = "解锁"
	}

	// 写入操作日志
	sqls.DB().Create(&models.OperateLog{
		UserId:      user.Id,
		OpType:      "freeze_flame",
		DataType:    "topic",
		DataId:      req.TopicId,
		Description: fmt.Sprintf("%s火焰等级（%d）", action, req.FlameLevel),
		CreateTime:  dates.NowTimestamp(),
	})

	ginx.WriteJSON(ctx, map[string]interface{}{
		"topicId":    req.TopicId,
		"lockedLevel": req.FlameLevel,
	})
}

// AdminHeatFrozenTopics 列出所有被锁定火焰等级的帖子
func AdminHeatFrozenTopics(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("未登录"))
		return
	}
	if !user.IsOwner() {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("无权限"))
		return
	}

	type frozenRow struct {
		Id               int64
		Title            string
		FlameLockedLevel int
	}
	var rows []frozenRow
	sqls.DB().Model(&models.Topic{}).
		Where("flame_locked_level > 0").
		Select("id, title, flame_locked_level").
		Find(&rows)

	ginx.WriteJSON(ctx, rows)
}

// AdminHeatCirculationStatus 获取热度点流通状态
func AdminHeatCirculationStatus(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("未登录"))
		return
	}
	if !user.IsOwner() {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("无权限"))
		return
	}

	// 取最新快照
	var latest models.HeatCirculationSnapshot
	sqls.DB().Order("snapshot_date desc").First(&latest)

	// 取奖池余额
	var pool models.HeatPublicPool
	sqls.DB().Order("id desc").First(&pool)

	// 取活跃用户数
	var activeUsers int64
	sqls.DB().Model(&models.UserHeatStats{}).
		Where("staked_in_window > 0").Count(&activeUsers)

	ginx.WriteJSON(ctx, map[string]interface{}{
		"totalSupply":       latest.TotalSupply,
		"activeCirculation": latest.ActiveCirculation,
		"stakedTotal":       latest.StakedTotal,
		"activeUsers":       activeUsers,
		"inactiveRatio": func() float64 {
			if latest.TotalSupply == 0 {
				return 0
			}
			return float64(latest.TotalSupply-latest.ActiveCirculation) / float64(latest.TotalSupply)
		}(),
		"publicPoolBalance": pool.BalanceAfter,
		"snapshotDate":      latest.SnapshotDate,
		"decayTruncated":    latest.DailyDecayTruncated,
	})
}

// AdminHeatTriggerSettlement 手动触发结算（应急用）
func AdminHeatTriggerSettlement(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("未登录"))
		return
	}
	if !user.IsOwner() {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("无权限"))
		return
	}

	// 写入操作日志
	sqls.DB().Create(&models.OperateLog{
		UserId:      user.Id,
		OpType:      "trigger_settlement",
		DataType:    "system",
		DataId:      0,
		Description: fmt.Sprintf("手动触发结算，时间=%s", time.Now().Format("2006-01-02 15:04:05")),
		CreateTime:  dates.NowTimestamp(),
	})

	// 实际触发结算
	go func() {
		if err := heatpoints.Settlement.SettleAll(); err != nil {
			slog.Error("手动触发结算失败", "error", err)
		}
	}()

	ginx.WriteJSON(ctx, map[string]interface{}{
		"status":  "triggered",
		"message": "结算已触发，将在后台执行",
	})
}
