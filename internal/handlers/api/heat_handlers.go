package api

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/services/heatpoints"

	"github.com/gin-gonic/gin"
)

// StakeCreate 创建质押
func StakeCreate(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("未登录"))
		return
	}

	var req struct {
		TopicId    string `json:"topicId"`
		HeatPoints int    `json:"heatPoints"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("参数错误："+err.Error()))
		return
	}

	topicId := idcodec.Decode(req.TopicId)
	if topicId <= 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("帖子ID无效"))
		return
	}

	if req.HeatPoints <= 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("质押额度必须大于 0"))
		return
	}

	resp, err := heatpoints.Stake.Create(&heatpoints.StakeCreateRequest{
		TopicId:    topicId,
		UserId:     user.Id,
		HeatPoints: req.HeatPoints,
	})

	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	ginx.WriteJSON(ctx, resp)
}

// StakeRedeem 赎回质押
func StakeRedeem(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("未登录"))
		return
	}

	stakeIdStr := ctx.Param("id")
	if stakeIdStr == "" {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("参数错误"))
		return
	}
	stakeId := idcodec.Decode(stakeIdStr)
	if stakeId <= 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("质押ID无效"))
		return
	}

	err := heatpoints.Stake.Redeem(user.Id, stakeId)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	ginx.WriteJSON(ctx, nil)
}

// StakeRecords 获取用户质押记录
func StakeRecords(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("未登录"))
		return
	}

	status := -1
	if s := ctx.Query("status"); s != "" {
		status = 0
	}
	limit := 20

	records, err := heatpoints.Stake.GetUserStakes(user.Id, status, limit)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	var items []map[string]interface{}
	for _, record := range records {
		items = append(items, map[string]interface{}{
			"id":             idcodec.Encode(record.Id),
			"topicId":        idcodec.Encode(record.TopicId),
			"heatPoints":     record.HeatPoints,
			"originalPoints": record.OriginalPoints,
			"stakeDay":       record.StakeDay,
			"status":         record.Status,
			"createTime":     record.CreateTime,
		})
	}

	ginx.WriteJSON(ctx, items)
}

// StakeQuota 获取用户配额
func StakeQuota(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("未登录"))
		return
	}

	quotaUsed := heatpoints.Stake.GetTodayQuotaUsed(user.Id)

	ginx.WriteJSON(ctx, map[string]interface{}{
		"remainingQuota": constants.HeatPointsStakeQuotaDaily - quotaUsed,
		"totalQuota":     constants.HeatPointsStakeQuotaDaily,
		"quotaUsed":      quotaUsed,
		"heatPoints":     user.HeatPoints,
	})
}

// StakeHeat 获取帖子火焰等级
func StakeHeat(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")
	if topicIdStr == "" {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("参数错误"))
		return
	}
	topicId := idcodec.Decode(topicIdStr)
	if topicId <= 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("帖子ID无效"))
		return
	}

	activeCirculation := heatpoints.HeatSnapshot.GetActiveCirculation()
	flameLevel := heatpoints.Stake.CalculateFlameLevel(topicId, activeCirculation)

	ginx.WriteJSON(ctx, map[string]interface{}{
		"flameLevel": flameLevel,
	})
}
