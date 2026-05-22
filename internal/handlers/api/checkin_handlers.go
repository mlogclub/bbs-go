package api

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/handlers/render"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/services"
	"time"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/web"
)

// PostCheckin 签到
// GetCheckin 获取签到信息
// GetRank 获取当天签到排行榜（最早签到的排在最前面）
func CheckinSubmit(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	err := services.CheckInService.CheckIn(user.Id)
	if err == nil {
		ginx.WriteJSON(ctx, nil)
		return
	} else {
		ginx.WriteJSON(ctx, err)
		return
	}

}

func CheckinStatus(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, nil)
		return
	}
	checkIn := services.CheckInService.GetByUserId(user.Id)
	if checkIn != nil {
		today := dates.GetDay(time.Now())
		ginx.WriteJSON(ctx, web.NewRspBuilder(checkIn).
			Put("checkIn", checkIn.LatestDayName == today). // 今日是否已签到
			Build())
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func CheckinRank(ctx *gin.Context) {

	list := cache.UserCache.GetCheckInRank()
	var itemList []map[string]interface{}
	for _, checkIn := range list {
		itemList = append(itemList, web.NewRspBuilder(checkIn).
			Put("user", render.BuildUserInfoDefaultIfNull(checkIn.UserId)).
			Build())
	}
	ginx.WriteJSON(ctx, itemList)

}
