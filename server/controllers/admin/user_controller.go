package admin

import (
	"bbs-go/model/constants"
	"strconv"
	"strings"

	"bbs-go/model"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/common/passwd"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/services"
)

type UserController struct {
	Ctx iris.Context
}

func (c *UserController) GetSynccount() *mvc.JsonResult {
	go func() {
		services.UserService.SyncUserCount()
	}()
	return mvc.JsonSuccess()
}

func (c *UserController) GetBy(id int64) *mvc.JsonResult {
	t := services.UserService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return mvc.JsonData(c.buildUserItem(t))
}

func (c *UserController) AnyList() *mvc.JsonResult {
	list, paging := services.UserService.FindPageByParams(params.NewQueryParams(c.Ctx).EqByReq("id").LikeByReq("nickname").EqByReq("username").PageByReq().Desc("id"))
	var itemList []map[string]interface{}
	for _, user := range list {
		itemList = append(itemList, c.buildUserItem(&user))
	}
	return mvc.JsonData(&sqls.PageResult{Results: itemList, Page: paging})
}

func (c *UserController) PostCreate() *mvc.JsonResult {
	username := params.FormValue(c.Ctx, "username")
	email := params.FormValue(c.Ctx, "email")
	nickname := params.FormValue(c.Ctx, "nickname")
	password := params.FormValue(c.Ctx, "password")

	user, err := services.UserService.SignUp(username, email, nickname, password, password)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonData(c.buildUserItem(user))
}

func (c *UserController) PostUpdate() *mvc.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	user := services.UserService.Get(id)
	if user == nil {
		return mvc.JsonErrorMsg("entity not found")
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
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonData(c.buildUserItem(user))
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
	if days == 0 {
		services.UserService.RemoveForbidden(user.Id, userId, c.Ctx.Request())
	} else {
		if err := services.UserService.Forbidden(user.Id, userId, days, reason, c.Ctx.Request()); err != nil {
			return mvc.JsonErrorMsg(err.Error())
		}
	}
	return mvc.JsonSuccess()
}

func (c *UserController) buildUserItem(user *model.User) map[string]interface{} {
	return mvc.NewRspBuilder(user).
		Put("roles", user.GetRoles()).
		Put("username", user.Username.String).
		Put("email", user.Email.String).
		Put("score", user.Score).
		Put("forbidden", user.IsForbidden()).
		Build()
}
