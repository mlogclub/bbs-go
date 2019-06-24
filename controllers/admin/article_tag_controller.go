package admin

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/simple"
	"strconv"
)

type ArticleTagController struct {
	Ctx               iris.Context
	ArticleTagService *services.ArticleTagService
}

func (this *ArticleTagController) GetBy(id int64) *simple.JsonResult {
	t := this.ArticleTagService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *ArticleTagController) AnyList() *simple.JsonResult {
	list, paging := this.ArticleTagService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *ArticleTagController) PostCreate() *simple.JsonResult {
	t := &model.ArticleTag{}
	this.Ctx.ReadForm(t)

	err := this.ArticleTagService.Create(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *ArticleTagController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	t := this.ArticleTagService.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	this.Ctx.ReadForm(t)

	err = this.ArticleTagService.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
