package api

import (
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/services"
	"bbs-go/services/cache"
)

type UserController struct {
	Ctx iris.Context
}

// 获取当前登录用户
func (c *UserController) GetCurrent() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user != nil {
		return simple.JsonData(render.BuildUser(user))
	}
	return simple.JsonSuccess()
}

// 用户详情
func (c *UserController) GetBy(userId int64) *simple.JsonResult {
	user := cache.UserCache.Get(userId)
	if user != nil && user.Status != model.StatusDeleted {
		return simple.JsonData(render.BuildUser(user))
	}
	return simple.JsonErrorMsg("用户不存在")
}

// 修改用户资料
func (c *UserController) PostEditBy(userId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	if user.Id != userId {
		return simple.JsonErrorMsg("无权限")
	}
	nickname := strings.TrimSpace(simple.FormValue(c.Ctx, "nickname"))
	avatar := strings.TrimSpace(simple.FormValue(c.Ctx, "avatar"))
	description := simple.FormValue(c.Ctx, "description")

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
func (c *UserController) PostSetUsername() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	username := strings.TrimSpace(simple.FormValue(c.Ctx, "username"))
	err := services.UserService.SetUsername(user.Id, username)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 设置邮箱
func (c *UserController) PostSetEmail() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	email := strings.TrimSpace(simple.FormValue(c.Ctx, "email"))
	err := services.UserService.SetEmail(user.Id, email)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 设置密码
func (c *UserController) PostSetPassword() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	password := simple.FormValue(c.Ctx, "password")
	rePassword := simple.FormValue(c.Ctx, "rePassword")
	err := services.UserService.SetPassword(user.Id, password, rePassword)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 修改密码
func (c *UserController) PostUpdatePassword() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	var (
		oldPassword = simple.FormValue(c.Ctx, "oldPassword")
		password    = simple.FormValue(c.Ctx, "password")
		rePassword  = simple.FormValue(c.Ctx, "rePassword")
	)
	err := services.UserService.UpdatePassword(user.Id, oldPassword, password, rePassword)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 用户收藏
func (c *UserController) GetFavorites() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	cursor := simple.FormValueInt64Default(c.Ctx, "cursor", 0)

	// 用户必须登录
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	// 查询列表
	var favorites []model.Favorite
	if cursor > 0 {
		favorites = services.FavoriteService.Find(simple.NewSqlCnd().Where("user_id = ? and id < ?",
			user.Id, cursor).Desc("id").Limit(20))
	} else {
		favorites = services.FavoriteService.Find(simple.NewSqlCnd().Where("user_id = ?", user.Id).Desc("id").Limit(20))
	}

	if len(favorites) > 0 {
		cursor = favorites[len(favorites)-1].Id
	}

	return simple.JsonCursorData(render.BuildFavorites(favorites), strconv.FormatInt(cursor, 10))
}

// 获取最近3条未读消息
func (c *UserController) GetMsgrecent() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	var count int64 = 0
	var messages []model.Message
	if user != nil {
		count = services.MessageService.GetUnReadCount(user.Id)
		messages = services.MessageService.Find(simple.NewSqlCnd().Eq("user_id", user.Id).Eq("status", model.MsgStatusUnread).Limit(3).Desc("id"))
	}
	return simple.NewEmptyRspBuilder().Put("count", count).Put("messages", render.BuildMessages(messages)).JsonResult()
}

// 用户消息
func (c *UserController) GetMessages() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	page := simple.FormValueIntDefault(c.Ctx, "page", 1)

	// 用户必须登录
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	messages, paging := services.MessageService.FindPageByCnd(simple.NewSqlCnd().
		Eq("user_id", user.Id).
		Page(page, 20).Desc("id"))

	// 全部标记为已读
	services.MessageService.MarkRead(user.Id)

	return simple.JsonPageData(render.BuildMessages(messages), paging)
}

// 最新用户
func (c *UserController) GetNewest() *simple.JsonResult {
	users := services.UserService.Find(simple.NewSqlCnd().Eq("type", model.UserTypeNormal).Desc("id").Limit(10))
	return simple.JsonData(render.BuildUsers(users))
}
