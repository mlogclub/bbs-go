package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/services"
)

type ArticleTagController struct {
	Ctx iris.Context
}

func (c *ArticleTagController) GetBy(id int64) *simple.JsonResult {
	t := services.ArticleTagService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *ArticleTagController) AnyList() *simple.JsonResult {
	list, paging := services.ArticleTagService.FindPageByParams(simple.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (c *ArticleTagController) PostCreate() *simple.JsonResult {
	t := &model.ArticleTag{}
	simple.ReadForm(c.Ctx, t)

	err := services.ArticleTagService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (c *ArticleTagController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.ArticleTagService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	simple.ReadForm(c.Ctx, t)

	err = services.ArticleTagService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
