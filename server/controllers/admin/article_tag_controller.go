package admin

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
	"github.com/mlogclub/simple"
	"strconv"
)

type ArticleTagController struct {
	Ctx iris.Context
}

func (this *ArticleTagController) GetBy(id int64) *simple.JsonResult {
	t := services.ArticleTagService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *ArticleTagController) AnyList() *simple.JsonResult {
	list, paging := services.ArticleTagService.FindPageByParams(simple.NewQueryParams(this.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *ArticleTagController) PostCreate() *simple.JsonResult {
	t := &model.ArticleTag{}
	this.Ctx.ReadForm(t)

	err := services.ArticleTagService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *ArticleTagController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.ArticleTagService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	this.Ctx.ReadForm(t)

	err = services.ArticleTagService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
