package api

import (
	"bbs-go/cache"
	"bbs-go/controllers/render"
	"bbs-go/services"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/date"
	"time"
)

type CheckinController struct {
	Ctx iris.Context
}

// PostCheckin 签到
func (c *CheckinController) PostCheckin() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return simple.JsonError(err)
	}
	err := services.CheckInService.CheckIn(user.Id)
	if err == nil {
		return simple.JsonSuccess()
	} else {
		return simple.JsonErrorMsg(err.Error())
	}
}

// GetCheckin 获取签到信息
func (c *CheckinController) GetCheckin() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonSuccess()
	}
	checkIn := services.CheckInService.GetByUserId(user.Id)
	if checkIn != nil {
		today := date.GetDay(time.Now())
		return simple.NewRspBuilder(checkIn).
			Put("checkIn", checkIn.LatestDayName == today). // 今日是否已签到
			JsonResult()
	}
	return simple.JsonSuccess()
}

// GetRank 获取当天签到排行榜（最早签到的排在最前面）
func (c *CheckinController) GetRank() *simple.JsonResult {
	list := cache.UserCache.GetCheckInRank()
	var itemList []map[string]interface{}
	for _, checkIn := range list {
		itemList = append(itemList, simple.NewRspBuilder(checkIn).
			Put("user", render.BuildUserById(checkIn.UserId)).
			Build())
	}
	return simple.JsonData(itemList)
}
