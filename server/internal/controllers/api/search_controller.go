package api

import (
	"log/slog"

	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/search"
	"bbs-go/internal/services"

	"github.com/kataras/iris/v12"
	"github.com/spf13/cast"

	"bbs-go/internal/pkg/simple/web"
	"bbs-go/internal/pkg/simple/web/params"
)

type SearchController struct {
	Ctx iris.Context
}

func (c *SearchController) AnyReindex() *web.JsonResult {
	if config.Instance.MeiliSearch.Enabled {
		go func() {
			if err := search.ReindexAllTopicsMeili(); err != nil {
				slog.Error("Failed to reindex topics in MeiliSearch", slog.Any("error", err))
			}
		}()
	} else {
		go services.TopicService.ScanDesc(func(topics []models.Topic) {
			for _, topic := range topics {
				if topic.Status != constants.StatusDeleted {
					search.UpdateTopicIndex(&topic)
				}
			}
		})
	}
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
