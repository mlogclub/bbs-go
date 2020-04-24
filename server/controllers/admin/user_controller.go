package admin

import (
	"strconv"
	"strings"

	"bbs-go/common"
	"bbs-go/model"
	"bbs-go/services/cache"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/services"
)

type UserController struct {
	Ctx iris.Context
}

func (c *UserController) GetSynccount() *simple.JsonResult {
	go func() {
		services.UserService.SyncUserCount()
	}()
	return simple.JsonSuccess()
}

func (c *UserController) GetBy(id int64) *simple.JsonResult {
	t := services.UserService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(c.buildUserItem(t))
}

func (c *UserController) AnyList() *simple.JsonResult {
	list, paging := services.UserService.FindPageByParams(simple.NewQueryParams(c.Ctx).EqByReq("id").LikeByReq("nickname").EqByReq("username").PageByReq().Desc("id"))
	var itemList []map[string]interface{}
	for _, user := range list {
		itemList = append(itemList, c.buildUserItem(&user))
	}
	return simple.JsonData(&simple.PageResult{Results: itemList, Page: paging})
}

func (c *UserController) PostCreate() *simple.JsonResult {
	username := simple.FormValue(c.Ctx, "username")
	email := simple.FormValue(c.Ctx, "email")
	nickname := simple.FormValue(c.Ctx, "nickname")
	password := simple.FormValue(c.Ctx, "password")

	user, err := services.UserService.SignUp(username, email, nickname, password, password)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(c.buildUserItem(user))
}

func (c *UserController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	user := services.UserService.Get(id)
	if user == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	username := simple.FormValue(c.Ctx, "username")
	password := simple.FormValue(c.Ctx, "password")
	nickname := simple.FormValue(c.Ctx, "nickname")
	email := simple.FormValue(c.Ctx, "email")
	roles := simple.FormValueStringArray(c.Ctx, "roles")
	status := simple.FormValueIntDefault(c.Ctx, "status", -1)

	user.Username = simple.SqlNullString(username)
	user.Nickname = nickname
	user.Email = simple.SqlNullString(email)
	user.Roles = strings.Join(roles, ",")
	user.Status = status

	if len(password) > 0 {
		user.Password = simple.EncodePassword(password)
	}

	err = services.UserService.Update(user)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(c.buildUserItem(user))
}

func (c *UserController) buildUserItem(user *model.User) map[string]interface{} {
	score := cache.UserCache.GetScore(user.Id)
	return simple.NewRspBuilder(user).
		Put("roles", common.GetUserRoles(user.Roles)).
		Put("username", user.Username.String).
		Put("email", user.Email.String).
		Put("score", score).
		Build()
}
