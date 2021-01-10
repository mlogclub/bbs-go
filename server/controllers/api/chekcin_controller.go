package api

import (
	"bbs-go/services"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
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
		today := services.CheckInService.GetDayName(time.Now())
		return simple.NewRspBuilder(checkIn).
			Put("checkIn", checkIn.LatestDayName == today). // 今日是否已签到
			JsonResult()
	}
	return simple.JsonSuccess()
}
