package api

import (
	"bbs-go/internal/models/constants"
	"strconv"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/sqls"

	"bbs-go/internal/cache"
	"bbs-go/internal/handlers/render"
	"bbs-go/internal/services"
)

// 标签详情
// 标签列表
// 标签自动完成
func TagDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	tagId := id

	tag := cache.TagCache.Get(tagId)
	if tag == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("tag not found"))
		return
	}
	ginx.WriteJSON(ctx, render.BuildTag(tag))

}

func TagTags(ctx *gin.Context) {
	page := params.FormValueIntDefault(ctx, "page", 1)
	tags, paging := services.TagService.FindPageByCnd(sqls.NewCnd().
		Eq("status", constants.StatusOk).
		Page(page, 200).Desc("id"))

	ginx.WriteJSON(ctx, ginx.PageData(render.BuildTags(tags), paging))

}

func TagAutocompleteSubmit(ctx *gin.Context) {
	var req struct {
		Input string `json:"input" form:"input"`
	}
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	tags := services.TagService.Autocomplete(req.Input)
	ginx.WriteJSON(ctx, tags)

}
