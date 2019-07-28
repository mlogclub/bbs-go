package api

import (
	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/services/cache"
)

type UserController struct {
	Ctx context.Context
}

// 获取当前登录用户
func (this *UserController) GetCurrent() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user != nil {
		return simple.JsonData(render.BuildUser(user))
	}
	return simple.ErrorMsg("未登录")
}

// 用户详情
func (this *UserController) GetBy(userId int64) *simple.JsonResult {
	user := cache.UserCache.Get(userId)
	if user != nil {
		return simple.JsonData(render.BuildUser(user))
	}
	return simple.ErrorMsg("用户不存在")
}

// 未读消息数量
func (this *UserController) GetMsgcount() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	var count int64 = 0
	if user != nil {
		count = services.MessageService.GetUnReadCount(user.Id)
	}
	return simple.NewEmptyRspBuilder().Put("count", count).JsonResult()
}
