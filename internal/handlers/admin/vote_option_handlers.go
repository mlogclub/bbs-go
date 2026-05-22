package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/web"
)

func VoteOptionDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	t := services.VoteOptionService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("Not found, id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func VoteOptionList(ctx *gin.Context) {
	list, paging := services.VoteOptionService.FindPageByCnd(params.NewPagedSqlCnd(ctx,
		params.QueryFilter{
			ParamName: "id",
		},
	).Desc("id"))
	ginx.WriteJSON(ctx, &web.PageResult{Results: list, Page: paging})

}

func VoteOptionCreate(ctx *gin.Context) {
	t := &models.VoteOption{}
	if err := ginx.Bind(ctx, t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}

	if err := services.VoteOptionService.Create(t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func VoteOptionUpdate(ctx *gin.Context) {
	id, _ := params.GetInt64(ctx, "id")
	t := services.VoteOptionService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("entity not found"))
		return
	}

	if err := ginx.Bind(ctx, t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}

	if err := services.VoteOptionService.Update(t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func VoteOptionRemove(ctx *gin.Context) {
	ids := params.GetInt64Arr(ctx, "ids")
	if len(ids) == 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("delete ids is empty"))
		return
	}
	for _, id := range ids {
		services.VoteOptionService.Delete(id)
	}
	ginx.WriteJSON(ctx, nil)

}
