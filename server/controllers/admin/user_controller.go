package admin

import (
	"strconv"
	"strings"

	"github.com/mlogclub/bbs-go/common"
	"github.com/mlogclub/bbs-go/model"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/services"
)

type UserController struct {
	Ctx iris.Context
}

func (this *UserController) GetBy(id int64) *simple.JsonResult {
	t := services.UserService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(this.buildUserItem(t))
}

func (this *UserController) AnyList() *simple.JsonResult {
	list, paging := services.UserService.FindPageByParams(simple.NewQueryParams(this.Ctx).EqByReq("id").LikeByReq("nickname").EqByReq("username").PageByReq().Desc("id"))
	var itemList []map[string]interface{}
	for _, user := range list {
		itemList = append(itemList, this.buildUserItem(&user))
	}
	return simple.JsonData(&simple.PageResult{Results: itemList, Page: paging})
}

func (this *UserController) PostCreate() *simple.JsonResult {
	username := simple.FormValue(this.Ctx, "username")
	email := simple.FormValue(this.Ctx, "email")
	nickname := simple.FormValue(this.Ctx, "nickname")
	password := simple.FormValue(this.Ctx, "password")

	user, err := services.UserService.SignUp(username, email, nickname, password, password)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(this.buildUserItem(user))
}

func (this *UserController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.UserService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	username := simple.FormValue(this.Ctx, "username")
	password := simple.FormValue(this.Ctx, "password")
	nickname := simple.FormValue(this.Ctx, "nickname")
	email := simple.FormValue(this.Ctx, "email")
	roles := simple.FormValueStringArray(this.Ctx, "roles")
	status := simple.FormValueIntDefault(this.Ctx, "status", -1)

	if len(username) > 0 {
		t.Username = simple.SqlNullString(username)
	}
	if len(password) > 0 {
		t.Password = simple.EncodePassword(t.Password)
	}
	if len(nickname) > 0 {
		t.Nickname = nickname
	}
	if len(email) > 0 {
		t.Email = simple.SqlNullString(email)
	}
	if status != -1 {
		t.Status = status
	}

	t.Roles = strings.Join(roles, ",")

	err = services.UserService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *UserController) buildUserItem(user *model.User) map[string]interface{} {
	return simple.NewRspBuilder(user).
		Put("roles", common.GetUserRoles(user.Roles)).
		Put("username", user.Username.String).
		Put("email", user.Email.String).
		Build()
}
