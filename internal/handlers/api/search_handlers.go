package api

import (
	"bbs-go/internal/handlers/render"
	"bbs-go/internal/pkg/search"
	"bbs-go/internal/services"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/spf13/cast"
)

func SearchTopic(ctx *gin.Context) {
	var (
		cursor     = params.FormValueIntDefault(ctx, "cursor", 1)
		keyword    = params.FormValue(ctx, "keyword")
		categoryId = params.FormValueInt64Default(ctx, "categoryId", 0)
		timeRange  = params.FormValueIntDefault(ctx, "timeRange", 0)
		limit      = 20
	)
	var categoryIds []int64
	if categoryId > 0 {
		categoryIds = services.CategoryService.GetCategoryIdsForList(categoryId)
	}
	list, _, err := search.SearchTopic(keyword, categoryId, categoryIds, timeRange, cursor, limit)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, ginx.CursorData(render.BuildSearchTopics(list), cast.ToString(cursor+1), len(list) >= limit))

}

func SearchArticle(ctx *gin.Context) {
	var (
		cursor    = params.FormValueIntDefault(ctx, "cursor", 1)
		keyword   = params.FormValue(ctx, "keyword")
		timeRange = params.FormValueIntDefault(ctx, "timeRange", 0)
		limit     = 20
	)
	list, _, err := search.SearchArticle(keyword, timeRange, cursor, limit)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, ginx.CursorData(render.BuildSearchArticles(list), cast.ToString(cursor+1), len(list) >= limit))
}

func SearchUser(ctx *gin.Context) {
	var (
		cursor  = params.FormValueIntDefault(ctx, "cursor", 1)
		keyword = params.FormValue(ctx, "keyword")
		limit   = 20
	)
	list, _, err := search.SearchUser(keyword, cursor, limit)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, ginx.CursorData(render.BuildSearchUsers(list), cast.ToString(cursor+1), len(list) >= limit))
}
