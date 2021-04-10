package api

import (
	"bbs-go/controllers/render"
	"bbs-go/package/es"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
)

type SearchController struct {
	Ctx iris.Context
}

func (c *SearchController) GetTopic() *simple.JsonResult {
	var (
		page    = simple.FormValueIntDefault(c.Ctx, "page", 1)
		keyword = simple.FormValue(c.Ctx, "keyword")
	)

	docs, paging, err := es.SearchTopic(keyword, page, 20)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	items := render.BuildSearchTopics(docs)
	return simple.JsonPageData(items, paging)
}
