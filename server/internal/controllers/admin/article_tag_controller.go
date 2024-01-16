package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

type ArticleTagController struct {
	Ctx iris.Context
}

func (c *ArticleTagController) GetBy(id int64) *web.JsonResult {
	t := services.ArticleTagService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *ArticleTagController) AnyList() *web.JsonResult {
	list, paging := services.ArticleTagService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *ArticleTagController) PostCreate() *web.JsonResult {
	t := &models.ArticleTag{}
	params.ReadForm(c.Ctx, t)

	err := services.ArticleTagService.Create(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *ArticleTagController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	t := services.ArticleTagService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	params.ReadForm(c.Ctx, t)

	err = services.ArticleTagService.Update(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}
