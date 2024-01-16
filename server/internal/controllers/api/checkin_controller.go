package api

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/services"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/web"
)

type CheckinController struct {
	Ctx iris.Context
}

// PostCheckin 签到
func (c *CheckinController) PostCheckin() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}
	err := services.CheckInService.CheckIn(user.Id)
	if err == nil {
		return web.JsonSuccess()
	} else {
		return web.JsonError(err)
	}
}

// GetCheckin 获取签到信息
func (c *CheckinController) GetCheckin() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonSuccess()
	}
	checkIn := services.CheckInService.GetByUserId(user.Id)
	if checkIn != nil {
		today := dates.GetDay(time.Now())
		return web.NewRspBuilder(checkIn).
			Put("checkIn", checkIn.LatestDayName == today). // 今日是否已签到
			JsonResult()
	}
	return web.JsonSuccess()
}

// GetRank 获取当天签到排行榜（最早签到的排在最前面）
func (c *CheckinController) GetRank() *web.JsonResult {
	list := cache.UserCache.GetCheckInRank()
	var itemList []map[string]interface{}
	for _, checkIn := range list {
		itemList = append(itemList, web.NewRspBuilder(checkIn).
			Put("user", render.BuildUserInfoDefaultIfNull(checkIn.UserId)).
			Build())
	}
	return web.JsonData(itemList)
}
