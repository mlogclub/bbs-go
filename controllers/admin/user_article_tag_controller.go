package admin

import (
	"strconv"

	"github.com/kataras/iris"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
)

type UserArticleTagController struct {
	Ctx                   iris.Context
	UserArticleTagService *services.UserArticleTagService
}

func (this *UserArticleTagController) GetBy(id int64) *simple.JsonResult {
	t := this.UserArticleTagService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *UserArticleTagController) AnyList() *simple.JsonResult {
	list, paging := this.UserArticleTagService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *UserArticleTagController) PostCreate() *simple.JsonResult {
	t := &model.UserArticleTag{}
	this.Ctx.ReadForm(t)

	err := this.UserArticleTagService.Create(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *UserArticleTagController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	t := this.UserArticleTagService.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	this.Ctx.ReadForm(t)

	err = this.UserArticleTagService.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
