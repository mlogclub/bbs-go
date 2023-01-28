package api

import (
	"server/controllers/render"
	"server/model"
	"server/pkg/es"
	"server/services"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type SearchController struct {
	Ctx iris.Context
}

func (c *SearchController) AnyReindex() *web.JsonResult {
	go services.TopicService.ScanDesc(func(topics []model.Topic) {
		for _, t := range topics {
			topic := services.TopicService.Get(t.Id)
			es.UpdateTopicIndex(topic)
		}
	})
	return web.JsonSuccess()
}

func (c *SearchController) PostSearch() *web.JsonResult {
	var (
		page      = params.FormValueIntDefault(c.Ctx, "page", 1)
		keyword   = params.FormValue(c.Ctx, "keyword")
		nodeId    = params.FormValueInt64Default(c.Ctx, "nodeId", 0)
		timeRange = params.FormValueIntDefault(c.Ctx, "timeRange", 0)
	)

	key := keyword[0:2]
	Index := keyword[3:]
	switch key {
	case ":t":
		docs, paging, err := es.SearchTopic(Index, nodeId, timeRange, page, 20)
		if err != nil {
			return web.JsonError(err)
		}

		items := render.BuildSearchTopics(docs)
		return web.JsonPageData(items, paging)
	case ":a":
		docs, paging, err := es.SearchArticle(Index, timeRange, page, 20)
		if err != nil {
			return web.JsonError(err)
		}

		items := render.BuildSearchArticles(docs)
		return web.JsonPageData(items, paging)
	default:
		web.JsonErrorMsg("参数错误")
	}
	return web.JsonSuccess()
}
