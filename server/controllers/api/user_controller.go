package api

import (
	"strconv"
	"strings"

	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/controllers/render"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
	"github.com/mlogclub/bbs-go/services/cache"
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

// 设置用户名
func (this *UserController) PostSetUsername() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	username := strings.TrimSpace(simple.FormValue(this.Ctx, "username"))
	err := services.UserService.SetUsername(user.Id, username)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 设置邮箱
func (this *UserController) PostSetEmail() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	email := strings.TrimSpace(simple.FormValue(this.Ctx, "email"))
	err := services.UserService.SetEmail(user.Id, email)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 设置密码
func (this *UserController) PostSetPassword() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	password := simple.FormValue(this.Ctx, "password")
	rePassword := simple.FormValue(this.Ctx, "rePassword")
	err := services.UserService.SetPassword(user.Id, password, rePassword)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 修改密码
func (this *UserController) PostUpdatePassword() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	var (
		oldPassword = simple.FormValue(this.Ctx, "oldPassword")
		password    = simple.FormValue(this.Ctx, "password")
		rePassword  = simple.FormValue(this.Ctx, "rePassword")
	)
	err := services.UserService.UpdatePassword(user.Id, oldPassword, password, rePassword)
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

// 用户收藏
func (this *UserController) GetFavorites() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	cursor := simple.FormValueInt64Default(this.Ctx, "cursor", 0)

	// 用户必须登录
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	// 查询列表
	var favorites []model.Favorite
	if cursor > 0 {
		favorites, _ = services.FavoriteService.QueryCnd(simple.NewQueryCnd("user_id = ? and id < ?",
			user.Id, cursor).Order("id desc").Size(20))
	} else {
		favorites, _ = services.FavoriteService.QueryCnd(simple.NewQueryCnd("user_id = ?",
			user.Id).Order("id desc").Size(20))
	}

	if len(favorites) > 0 {
		cursor = favorites[len(favorites)-1].Id
	}

	return simple.JsonCursorData(render.BuildFavorites(favorites), strconv.FormatInt(cursor, 10))
}

// 用户消息
func (this *UserController) GetMessages() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	cursor := simple.FormValueInt64Default(this.Ctx, "cursor", 0)

	// 用户必须登录
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	// 查询列表
	var messages []model.Message
	if cursor > 0 {
		messages, _ = services.MessageService.QueryCnd(simple.NewQueryCnd("user_id = ? and id < ?",
			user.Id, cursor).Order("id desc").Size(20))
	} else {
		messages, _ = services.MessageService.QueryCnd(simple.NewQueryCnd("user_id = ?",
			user.Id).Order("id desc").Size(20))
	}

	if len(messages) > 0 {
		cursor = messages[len(messages)-1].Id
	}

	// 全部标记为已读
	services.MessageService.MarkReadAll(user.Id)

	return simple.JsonCursorData(render.BuildMessages(messages), strconv.FormatInt(cursor, 10))
}
