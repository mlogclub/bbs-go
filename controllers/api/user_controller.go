package api

import (
	"strings"

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
	return simple.JsonSuccess()
}

// 用户详情
func (this *UserController) GetBy(userId int64) *simple.JsonResult {
	user := cache.UserCache.Get(userId)
	if user != nil {
		return simple.JsonData(render.BuildUser(user))
	}
	return simple.JsonErrorMsg("用户不存在")
}

// 修改用户资料
func (this *UserController) PostEditBy(userId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	if user.Id != userId {
		return simple.JsonErrorMsg("无权限")
	}
	nickname := strings.TrimSpace(simple.FormValue(this.Ctx, "nickname"))
	avatar := strings.TrimSpace(simple.FormValue(this.Ctx, "avatar"))
	description := simple.FormValue(this.Ctx, "description")

	if len(nickname) == 0 {
		return simple.JsonErrorMsg("昵称不能为空")
	}
	if len(avatar) == 0 {
		return simple.JsonErrorMsg("头像不能为空")
	}

	err := services.UserService.Updates(user.Id, map[string]interface{}{
		"nickname":    nickname,
		"avatar":      avatar,
		"description": description,
	})
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
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

// 活跃用户
func (this *UserController) GetActive() *simple.JsonResult {
	users := cache.UserCache.GetActiveUsers()
	return simple.JsonData(render.BuildUsers(users))
}
