package api

import (
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/search"
	"bbs-go/internal/services"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
	"github.com/spf13/cast"
)

type SearchController struct {
	Ctx iris.Context
}

func (c *SearchController) AnyReindex() *web.JsonResult {
	go services.TopicService.ScanDesc(func(topics []models.Topic) {
		for _, topic := range topics {
			if topic.Status != constants.StatusDeleted {
				search.UpdateTopicIndex(&topic)
			}
		}
	})
	return web.JsonSuccess()
}

func (c *SearchController) GetTopic() *web.JsonResult {
	var (
		cursor    = params.FormValueIntDefault(c.Ctx, "cursor", 1)
		keyword   = params.FormValue(c.Ctx, "keyword")
		nodeId    = params.FormValueInt64Default(c.Ctx, "nodeId", 0)
		timeRange = params.FormValueIntDefault(c.Ctx, "timeRange", 0)
		limit     = 20
	)
	list, _, err := search.SearchTopic(keyword, nodeId, timeRange, cursor, limit)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonCursorData(render.BuildSearchTopics(list), cast.ToString(cursor+1), len(list) >= limit)
}
