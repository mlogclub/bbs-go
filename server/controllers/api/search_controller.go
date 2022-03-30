package api

import (
	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/pkg/es"
	"bbs-go/services"

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

func (c *SearchController) PostTopic() *web.JsonResult {
	var (
		page      = params.FormValueIntDefault(c.Ctx, "page", 1)
		keyword   = params.FormValue(c.Ctx, "keyword")
		nodeId    = params.FormValueInt64Default(c.Ctx, "nodeId", 0)
		timeRange = params.FormValueIntDefault(c.Ctx, "timeRange", 0)
	)

	docs, paging, err := es.SearchTopic(keyword, nodeId, timeRange, page, 20)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	items := render.BuildSearchTopics(docs)
	return web.JsonPageData(items, paging)
}
