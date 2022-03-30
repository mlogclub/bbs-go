package api

import (
	"bbs-go/cache"
	"bbs-go/controllers/render"
	"bbs-go/services"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/mvc"
)

type CheckinController struct {
	Ctx iris.Context
}

// PostCheckin 签到
func (c *CheckinController) PostCheckin() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return mvc.JsonError(err)
	}
	err := services.CheckInService.CheckIn(user.Id)
	if err == nil {
		return mvc.JsonSuccess()
	} else {
		return mvc.JsonErrorMsg(err.Error())
	}
}

// GetCheckin 获取签到信息
func (c *CheckinController) GetCheckin() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return mvc.JsonSuccess()
	}
	checkIn := services.CheckInService.GetByUserId(user.Id)
	if checkIn != nil {
		today := dates.GetDay(time.Now())
		return mvc.NewRspBuilder(checkIn).
			Put("checkIn", checkIn.LatestDayName == today). // 今日是否已签到
			JsonResult()
	}
	return mvc.JsonSuccess()
}

// GetRank 获取当天签到排行榜（最早签到的排在最前面）
func (c *CheckinController) GetRank() *mvc.JsonResult {
	list := cache.UserCache.GetCheckInRank()
	var itemList []map[string]interface{}
	for _, checkIn := range list {
		itemList = append(itemList, mvc.NewRspBuilder(checkIn).
			Put("user", render.BuildUserInfoDefaultIfNull(checkIn.UserId)).
			Build())
	}
	return mvc.JsonData(itemList)
}
