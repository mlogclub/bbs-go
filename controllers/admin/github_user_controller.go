package admin

import (
	"strconv"

	"github.com/kataras/iris"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
)

type GithubUserController struct {
	Ctx iris.Context
}

func (this *GithubUserController) GetBy(id int64) *simple.JsonResult {
	t := services.GithubUserService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *GithubUserController) AnyList() *simple.JsonResult {
	list, paging := services.GithubUserService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *GithubUserController) PostCreate() *simple.JsonResult {
	t := &model.GithubUser{}
	this.Ctx.ReadForm(t)

	err := services.GithubUserService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *GithubUserController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.GithubUserService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	this.Ctx.ReadForm(t)

	err = services.GithubUserService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
