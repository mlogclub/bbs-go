package api

import (
	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/pkg/es"
	"bbs-go/services"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
)

type SearchController struct {
	Ctx iris.Context
}

func (c *SearchController) AnyReindex() *simple.JsonResult {
	go services.TopicService.ScanDesc(func(topics []model.Topic) {
		for _, t := range topics {
			topic := services.TopicService.Get(t.Id)
			es.UpdateTopicIndex(topic)
		}
	})
	return simple.JsonSuccess()
}

func (c *SearchController) PostTopic() *simple.JsonResult {
	var (
		page      = simple.FormValueIntDefault(c.Ctx, "page", 1)
		keyword   = simple.FormValue(c.Ctx, "keyword")
		nodeId    = simple.FormValueInt64Default(c.Ctx, "nodeId", 0)
		timeRange = simple.FormValueIntDefault(c.Ctx, "timeRange", 0)
	)

	docs, paging, err := es.SearchTopic(keyword, nodeId, timeRange, page, 20)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	items := render.BuildSearchTopics(docs)
	return simple.JsonPageData(items, paging)
}
