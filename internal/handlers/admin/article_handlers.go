package admin

import (
	"bbs-go/internal/models/constants"
	modelReq "bbs-go/internal/models/req"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/pkg/locales"
	"strconv"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/models"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/web"

	"bbs-go/internal/cache"
	"bbs-go/internal/handlers/render"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/services"
)

// 构建文章列表返回数据
func articleBuildArticles(articles []models.Article) []map[string]interface{} {
	var results []map[string]interface{}
	for _, article := range articles {
		builder := web.NewRspBuilderExcludes(article, "content")

		// 用户
		builder = builder.Put("user", render.BuildUserInfoDefaultIfNull(article.UserId))

		// 简介
		builder.Put("summary", common.GetSummary(article.ContentType, article.Content))

		// 标签
		tagIds := cache.ArticleTagCache.Get(article.Id)
		tags := cache.TagCache.GetList(tagIds)
		builder.Put("tags", render.BuildTags(tags))

		// 封面
		builder.Put("cover", render.BuildImage(article.Cover))

		results = append(results, builder.Build())
	}
	return results
}

func ArticleDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	t := services.ArticleService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("Not found, id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func ArticleList(ctx *gin.Context) {
	list, paging := services.ArticleService.FindPageByCnd(params.NewPagedSqlCnd(ctx,
		params.QueryFilter{
			ParamName: "id",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "userId",
			Op:        params.Eq,
			ValueWrapper: func(origin string) string {
				if id := idcodec.Decode(origin); id > 0 {
					return strconv.FormatInt(id, 10)
				}
				return ""
			},
		},
		params.QueryFilter{
			ParamName: "status",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "title",
			Op:        params.Like,
		},
	).Desc("id"))

	results := articleBuildArticles(list)
	ginx.WriteJSON(ctx, ginx.PageData(results, paging))

}

func ArticleUpdate(ctx *gin.Context) {
	id, _ := params.GetInt64(ctx, "id")
	if id <= 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("id is required"))
		return
	}
	t := services.ArticleService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("entity not found"))
		return
	}

	if err := ginx.Bind(ctx, t); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	// 数据校验
	if len(t.Title) == 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("article.title_required")))
		return
	}
	if len(t.Content) == 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("article.content_required")))
		return
	}
	if len(t.ContentType) == 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.content_type_required")))
		return
	}

	t.UpdateTime = dates.NowTimestamp()
	err := services.ArticleService.Update(t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	ginx.WriteJSON(ctx, t)

}

func ArticleTags(ctx *gin.Context) {
	var req modelReq.ArticleTagsReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	tags := services.ArticleService.GetArticleTags(req.ArticleId)
	ginx.WriteJSON(ctx, render.BuildTags(tags))

}

func ArticleSaveTags(ctx *gin.Context) {
	var req modelReq.ArticleTagsReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	services.ArticleService.PutTags(req.ArticleId, modelReq.SplitCommaStrings(req.Tags))
	ginx.WriteJSON(ctx, render.BuildTags(services.ArticleService.GetArticleTags(req.ArticleId)))

}

func ArticleRemove(ctx *gin.Context) {
	id, _ := params.GetInt64(ctx, "id")
	if id <= 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("id is required"))
		return
	}
	err := services.ArticleService.Delete(id)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func ArticleAudit(ctx *gin.Context) {
	id, _ := params.GetInt64(ctx, "id")
	if id <= 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("id is required"))
		return
	}
	err := services.ArticleService.UpdateColumn(id, "status", constants.StatusOk)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}
