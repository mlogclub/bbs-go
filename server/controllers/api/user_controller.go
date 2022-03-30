package api

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/msg"
	"bbs-go/pkg/validate"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/cache"
	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/services"
)

type UserController struct {
	Ctx iris.Context
}

// 获取当前登录用户
func (c *UserController) GetCurrent() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user != nil {
		return mvc.JsonData(render.BuildUserProfile(user))
	}
	return mvc.JsonSuccess()
}

// 用户详情
func (c *UserController) GetBy(userId int64) *mvc.JsonResult {
	user := cache.UserCache.Get(userId)
	if user != nil && user.Status != constants.StatusDeleted {
		return mvc.JsonData(render.BuildUserProfile(user))
	}
	return mvc.JsonErrorMsg("用户不存在")
}

// 修改用户资料
func (c *UserController) PostEditBy(userId int64) *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}
	if user.Id != userId {
		return mvc.JsonErrorMsg("无权限")
	}
	nickname := strings.TrimSpace(params.FormValue(c.Ctx, "nickname"))
	homePage := params.FormValue(c.Ctx, "homePage")
	description := params.FormValue(c.Ctx, "description")

	if len(nickname) == 0 {
		return mvc.JsonErrorMsg("昵称不能为空")
	}

	if len(homePage) > 0 && validate.IsURL(homePage) != nil {
		return mvc.JsonErrorMsg("个人主页地址错误")
	}

	err := services.UserService.Updates(user.Id, map[string]interface{}{
		"nickname":    nickname,
		"home_page":   homePage,
		"description": description,
	})
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}

// 修改头像
func (c *UserController) PostUpdateAvatar() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}
	avatar := strings.TrimSpace(params.FormValue(c.Ctx, "avatar"))
	if len(avatar) == 0 {
		return mvc.JsonErrorMsg("头像不能为空")
	}
	err := services.UserService.UpdateAvatar(user.Id, avatar)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}

// 设置用户名
func (c *UserController) PostSetUsername() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}
	username := strings.TrimSpace(params.FormValue(c.Ctx, "username"))
	err := services.UserService.SetUsername(user.Id, username)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}

// 设置邮箱
func (c *UserController) PostSetEmail() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}
	email := strings.TrimSpace(params.FormValue(c.Ctx, "email"))
	err := services.UserService.SetEmail(user.Id, email)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}

// 设置密码
func (c *UserController) PostSetPassword() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}
	password := params.FormValue(c.Ctx, "password")
	rePassword := params.FormValue(c.Ctx, "rePassword")
	err := services.UserService.SetPassword(user.Id, password, rePassword)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}

// 修改密码
func (c *UserController) PostUpdatePassword() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}
	var (
		oldPassword = params.FormValue(c.Ctx, "oldPassword")
		password    = params.FormValue(c.Ctx, "password")
		rePassword  = params.FormValue(c.Ctx, "rePassword")
	)
	if err := services.UserService.UpdatePassword(user.Id, oldPassword, password, rePassword); err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}

// 设置背景图
func (c *UserController) PostSetBackgroundImage() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}
	backgroundImage := params.FormValue(c.Ctx, "backgroundImage")
	if strs.IsBlank(backgroundImage) {
		return mvc.JsonErrorMsg("请上传图片")
	}
	if err := services.UserService.UpdateBackgroundImage(user.Id, backgroundImage); err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}

// 用户收藏
func (c *UserController) GetFavorites() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)

	// 用户必须登录
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}

	// 查询列表
	limit := 20
	var favorites []model.Favorite
	if cursor > 0 {
		favorites = services.FavoriteService.Find(sqls.NewSqlCnd().Where("user_id = ? and id < ?",
			user.Id, cursor).Desc("id").Limit(20))
	} else {
		favorites = services.FavoriteService.Find(sqls.NewSqlCnd().Where("user_id = ?", user.Id).Desc("id").Limit(limit))
	}

	hasMore := false
	if len(favorites) > 0 {
		cursor = favorites[len(favorites)-1].Id
		hasMore = len(favorites) >= limit
	}

	return mvc.JsonCursorData(render.BuildFavorites(favorites), strconv.FormatInt(cursor, 10), hasMore)
}

// 获取最近3条未读消息
func (c *UserController) GetMsgrecent() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	var count int64 = 0
	var messages []model.Message
	if user != nil {
		count = services.MessageService.GetUnReadCount(user.Id)
		messages = services.MessageService.Find(sqls.NewSqlCnd().Eq("user_id", user.Id).
			Eq("status", msg.StatusUnread).Limit(3).Desc("id"))
	}
	return mvc.NewEmptyRspBuilder().Put("count", count).Put("messages", render.BuildMessages(messages)).JsonResult()
}

// 用户消息
func (c *UserController) GetMessages() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	page := params.FormValueIntDefault(c.Ctx, "page", 1)

	// 用户必须登录
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}

	messages, paging := services.MessageService.FindPageByCnd(sqls.NewSqlCnd().
		Eq("user_id", user.Id).
		Page(page, 20).Desc("id"))

	// 全部标记为已读
	services.MessageService.MarkRead(user.Id)

	return mvc.JsonPageData(render.BuildMessages(messages), paging)
}

// 用户积分记录
func (c *UserController) GetScorelogs() *mvc.JsonResult {
	page := params.FormValueIntDefault(c.Ctx, "page", 1)
	user := services.UserTokenService.GetCurrent(c.Ctx)
	// 用户必须登录
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}

	logs, paging := services.UserScoreLogService.FindPageByCnd(sqls.NewSqlCnd().
		Eq("user_id", user.Id).
		Page(page, 20).Desc("id"))

	return mvc.JsonPageData(logs, paging)
}

// 积分排行
func (c *UserController) GetScoreRank() *mvc.JsonResult {
	users := cache.UserCache.GetScoreRank()
	var results []*model.UserInfo
	for _, user := range users {
		results = append(results, render.BuildUserInfo(&user))
	}
	return mvc.JsonData(results)
}

// 禁言
func (c *UserController) PostForbidden() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}
	if !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return mvc.JsonErrorMsg("无权限")
	}
	var (
		userId = params.FormValueInt64Default(c.Ctx, "userId", 0)
		days   = params.FormValueIntDefault(c.Ctx, "days", 0)
		reason = params.FormValue(c.Ctx, "reason")
	)
	if userId < 0 {
		return mvc.JsonErrorMsg("请传入：userId")
	}
	if days == -1 && !user.HasRole(constants.RoleOwner) {
		return mvc.JsonErrorMsg("无永久禁言权限")
	}
	if days == 0 {
		services.UserService.RemoveForbidden(user.Id, userId, c.Ctx.Request())
	} else {
		if err := services.UserService.Forbidden(user.Id, userId, days, reason, c.Ctx.Request()); err != nil {
			return mvc.JsonErrorMsg(err.Error())
		}
	}
	return mvc.JsonSuccess()
}

// PostEmailVerify 请求邮箱验证邮件
func (c *UserController) PostSend_verify_email() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}
	if err := services.UserService.SendEmailVerifyEmail(user.Id); err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}

// PostVerify_email 获取邮箱验证码
func (c *UserController) PostVerify_email() *mvc.JsonResult {
	token := params.FormValue(c.Ctx, "token")
	if strs.IsBlank(token) {
		return mvc.JsonErrorMsg("Illegal request")
	}
	var (
		email string
		err   error
	)
	if email, err = services.UserService.VerifyEmail(token); err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.NewEmptyRspBuilder().Put("email", email).JsonResult()
}
