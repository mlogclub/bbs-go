package api

import (
	"bbs-go/model/constants"
	"bbs-go/package/validate"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"strconv"
	"strings"

	"bbs-go/cache"
	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/services"
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
	if user != nil && user.Status != constants.StatusDeleted {
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
	homePage := simple.FormValue(c.Ctx, "homePage")
	description := simple.FormValue(c.Ctx, "description")

	if len(nickname) == 0 {
		return simple.JsonErrorMsg("昵称不能为空")
	}
	if len(avatar) == 0 {
		return simple.JsonErrorMsg("头像不能为空")
	}

	if len(homePage) > 0 && validate.IsURL(homePage) != nil {
		return simple.JsonErrorMsg("个人主页地址错误")
	}

	err := services.UserService.Updates(user.Id, map[string]interface{}{
		"nickname":    nickname,
		"avatar":      avatar,
		"home_page":   homePage,
		"description": description,
	})
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 修改头像
func (c *UserController) PostUpdateAvatar() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	avatar := strings.TrimSpace(simple.FormValue(c.Ctx, "avatar"))
	if len(avatar) == 0 {
		return simple.JsonErrorMsg("头像不能为空")
	}
	err := services.UserService.UpdateAvatar(user.Id, avatar)
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
	if err := services.UserService.UpdatePassword(user.Id, oldPassword, password, rePassword); err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 设置背景图
func (c *UserController) PostSetBackgroundImage() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	backgroundImage := simple.FormValue(c.Ctx, "backgroundImage")
	if simple.IsBlank(backgroundImage) {
		return simple.JsonErrorMsg("请上传图片")
	}
	if err := services.UserService.UpdateBackgroundImage(user.Id, backgroundImage); err != nil {
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
		messages = services.MessageService.Find(simple.NewSqlCnd().Eq("user_id", user.Id).
			Eq("status", constants.MsgStatusUnread).Limit(3).Desc("id"))
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
	users := services.UserService.Find(simple.NewSqlCnd().Eq("type", constants.UserTypeNormal).Desc("id").Limit(10))
	return simple.JsonData(render.BuildUsers(users))
}

// 用户积分记录
func (c *UserController) GetScorelogs() *simple.JsonResult {
	page := simple.FormValueIntDefault(c.Ctx, "page", 1)
	user := services.UserTokenService.GetCurrent(c.Ctx)
	// 用户必须登录
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	logs, paging := services.UserScoreLogService.FindPageByCnd(simple.NewSqlCnd().
		Eq("user_id", user.Id).
		Page(page, 20).Desc("id"))

	return simple.JsonPageData(logs, paging)
}

// 积分排行
func (c *UserController) GetScoreRank() *simple.JsonResult {
	users := cache.UserCache.GetScoreRank()
	var results []*model.UserInfo
	for _, user := range users {
		results = append(results, render.BuildUser(&user))
	}
	return simple.JsonData(results)
}

// 禁言
func (c *UserController) PostForbidden() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	if !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return simple.JsonErrorMsg("无权限")
	}
	var (
		userId = simple.FormValueInt64Default(c.Ctx, "userId", 0)
		days   = simple.FormValueIntDefault(c.Ctx, "days", 0)
		reason = simple.FormValue(c.Ctx, "reason")
	)
	if userId < 0 {
		return simple.JsonErrorMsg("请传入：userId")
	}
	if days == -1 && !user.HasRole(constants.RoleOwner) {
		return simple.JsonErrorMsg("无永久禁言权限")
	}
	if days == 0 {
		services.UserService.RemoveForbidden(user.Id, userId, c.Ctx.Request())
	} else {
		if err := services.UserService.Forbidden(user.Id, userId, days, reason, c.Ctx.Request()); err != nil {
			return simple.JsonErrorMsg(err.Error())
		}
	}
	return simple.JsonSuccess()
}

// PostEmailVerify 请求邮箱验证邮件S
func (c *UserController) PostEmailVerify() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	if err := services.UserService.SendEmailVerifyEmail(user.Id); err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// GetEmailVerify 获取邮箱验证码
func (c *UserController) GetEmailVerify() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	token := simple.FormValue(c.Ctx, "token")
	if simple.IsBlank(token) {
		return simple.JsonErrorMsg("非法请求")
	}
	if err := services.UserService.VerifyEmail(user.Id, token); err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}
