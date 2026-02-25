package admin

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/repositories"
	"strconv"

	"bbs-go/internal/models"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
	"github.com/spf13/cast"

	"bbs-go/internal/services"
)

type UserController struct {
	Ctx iris.Context
}

func (c *UserController) GetSynccount() *web.JsonResult {
	go func() {
		services.UserService.Scan(func(users []models.User) {
			for _, user := range users {
				topicCount := repositories.TopicRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).Eq("status", constants.StatusOk))
				commentCount := repositories.CommentRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).Eq("status", constants.StatusOk))
				_ = repositories.UserRepository.UpdateColumn(sqls.DB(), user.Id, "topic_count", topicCount)
				_ = repositories.UserRepository.UpdateColumn(sqls.DB(), user.Id, "comment_count", commentCount)
				cache.UserCache.Invalidate(user.Id)
			}
		})
	}()
	return web.JsonSuccess()
}

func (c *UserController) GetBy(id int64) *web.JsonResult {
	t := services.UserService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(c.buildUserItem(t, true))
}

func (c *UserController) AnyList() *web.JsonResult {
	list, paging := services.UserService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "id",
			Op:        params.Eq,
			ValueWrapper: func(origin string) string {
				if id := idcodec.Decode(origin); id > 0 {
					return cast.ToString(id)
				}
				return ""
			},
		},
		params.QueryFilter{
			ParamName: "nickname",
			Op:        params.Like,
		},
		params.QueryFilter{
			ParamName: "email",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "username",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "type",
			Op:        params.Eq,
		},
	).Desc("id"))
	var itemList []map[string]interface{}
	for _, user := range list {
		itemList = append(itemList, c.buildUserItem(&user, false))
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
		return web.JsonError(err)
	}
	return web.JsonData(c.buildUserItem(user, true))
}

func (c *UserController) PostUpdate() *web.JsonResult {
	var (
		id, _       = params.GetInt64(c.Ctx, "id")
		_type, _    = params.GetInt(c.Ctx, "type")
		username    = params.FormValue(c.Ctx, "username")
		email       = params.FormValue(c.Ctx, "email")
		nickname    = params.FormValue(c.Ctx, "nickname")
		avatar      = params.FormValue(c.Ctx, "avatar")
		gender      = params.FormValue(c.Ctx, "gender")
		homePage    = params.FormValue(c.Ctx, "homePage")
		description = params.FormValue(c.Ctx, "description")
		roleIds     = params.FormValueInt64Array(c.Ctx, "roleIds")
		status      = params.FormValueIntDefault(c.Ctx, "status", 0)
	)

	user := services.UserService.Get(id)
	if user == nil {
		return web.JsonErrorMsg("entity not found")
	}

	user.Type = _type
	user.Username = sqls.SqlNullString(username)
	user.Email = sqls.SqlNullString(email)
	user.Nickname = nickname
	user.Avatar = avatar
	user.Gender = constants.Gender(gender)
	user.HomePage = homePage
	user.Description = description
	user.Status = status

	if err := services.UserService.Update(user); err != nil {
		return web.JsonError(err)
	}
	if err := services.UserRoleService.UpdateUserRoles(user.Id, roleIds); err != nil {
		return web.JsonError(err)
	}
	user = services.UserService.Get(user.Id)
	return web.JsonData(c.buildUserItem(user, true))
}

// 禁言
func (c *UserController) PostForbidden() *web.JsonResult {
	user := common.GetCurrentUser(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin())
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
			return web.JsonError(err)
		}
	}
	return web.JsonSuccess()
}

// 修改自己的密码
func (c *UserController) PostUpdate_password() *web.JsonResult {
	user := common.GetCurrentUser(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin())
	}
	var (
		oldPassword = params.FormValue(c.Ctx, "oldPassword")
		password    = params.FormValue(c.Ctx, "password")
		rePassword  = params.FormValue(c.Ctx, "rePassword")
	)
	if err := services.UserService.UpdatePassword(user.Id, oldPassword, password, rePassword); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// PostResetPassword 重置密码
func (c *UserController) PostReset_password() *web.JsonResult {
	userId, _ := params.GetInt64(c.Ctx, "userId")

	if userId <= 0 {
		return web.JsonErrorMsg("invalid param: userId")
	}

	newPassword, err := services.UserService.ResetPassword(userId)
	if err != nil {
		return web.JsonError(err)
	}

	return web.JsonData(iris.Map{
		"password": newPassword,
	})
}

func (c *UserController) buildUserItem(user *models.User, buildRoleIds bool) map[string]interface{} {
	b := web.NewRspBuilder(user).
		Put("idEncode", idcodec.Encode(user.Id)).
		Put("roles", user.GetRoles()).
		Put("username", user.Username.String).
		Put("email", user.Email.String).
		Put("score", user.Score).
		Put("forbidden", user.IsForbidden())
	if buildRoleIds {
		b.Put("roleIds", services.UserRoleService.GetUserRoleIds(user.Id))
	}
	return b.Build()
}
