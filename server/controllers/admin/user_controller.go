package admin

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/common"
	"strconv"
	"strings"

	"bbs-go/model"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/passwd"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/services"
)

type UserController struct {
	Ctx iris.Context
}

func (c *UserController) GetSynccount() *web.JsonResult {
	go func() {
		services.UserService.SyncUserCount()
	}()
	return web.JsonSuccess()
}

func (c *UserController) GetBy(id int64) *web.JsonResult {
	t := services.UserService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(c.buildUserItem(t))
}

func (c *UserController) AnyList() *web.JsonResult {
	list, paging := services.UserService.FindPageByParams(params.NewQueryParams(c.Ctx).EqByReq("id").LikeByReq("nickname").EqByReq("username").PageByReq().Desc("id"))
	var itemList []map[string]interface{}
	for _, user := range list {
		itemList = append(itemList, c.buildUserItem(&user))
	}
	return web.JsonData(&web.PageResult{Results: itemList, Page: paging})
}

func (c *UserController) PostCreate() *web.JsonResult {
	username := params.FormValue(c.Ctx, "username")
	email := params.FormValue(c.Ctx, "email")
	nickname := params.FormValue(c.Ctx, "nickname")
	password := params.FormValue(c.Ctx, "password")

	user, err := services.UserService.SignUp(username, email, nickname, password, password)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(c.buildUserItem(user))
}

func (c *UserController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	user := services.UserService.Get(id)
	if user == nil {
		return web.JsonErrorMsg("entity not found")
	}

	username := params.FormValue(c.Ctx, "username")
	password := params.FormValue(c.Ctx, "password")
	nickname := params.FormValue(c.Ctx, "nickname")
	email := params.FormValue(c.Ctx, "email")
	roles := params.FormValueStringArray(c.Ctx, "roles")
	status := params.FormValueIntDefault(c.Ctx, "status", -1)

	user.Username = sqls.SqlNullString(username)
	user.Nickname = nickname
	user.Email = sqls.SqlNullString(email)
	user.Roles = strings.Join(roles, ",")
	user.Status = status

	if len(password) > 0 {
		user.Password = passwd.EncodePassword(password)
	}

	err = services.UserService.Update(user)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(c.buildUserItem(user))
}

// 禁言
func (c *UserController) PostForbidden() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(common.ErrorNotLogin)
	}
	if !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return web.JsonErrorMsg("无权限")
	}
	var (
		userId = params.FormValueInt64Default(c.Ctx, "userId", 0)
		days   = params.FormValueIntDefault(c.Ctx, "days", 0)
		reason = params.FormValue(c.Ctx, "reason")
	)
	if userId < 0 {
		return web.JsonErrorMsg("请传入：userId")
	}
	if days == 0 {
		services.UserService.RemoveForbidden(user.Id, userId, c.Ctx.Request())
	} else {
		if err := services.UserService.Forbidden(user.Id, userId, days, reason, c.Ctx.Request()); err != nil {
			return web.JsonErrorMsg(err.Error())
		}
	}
	return web.JsonSuccess()
}

func (c *UserController) buildUserItem(user *model.User) map[string]interface{} {
	return web.NewRspBuilder(user).
		Put("roles", user.GetRoles()).
		Put("username", user.Username.String).
		Put("email", user.Email.String).
		Put("score", user.Score).
		Put("forbidden", user.IsForbidden()).
		Build()
}
