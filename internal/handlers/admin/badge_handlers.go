package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/web"
)

func BadgeDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	t := services.BadgeService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("Not found, id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func BadgeList(ctx *gin.Context) {
	list, paging := services.BadgeService.FindPageByCnd(params.NewPagedSqlCnd(ctx,
		params.QueryFilter{
			ParamName: "id",
		},
		params.QueryFilter{
			ParamName: "name",
			Op:        params.Like,
		},
		params.QueryFilter{
			ParamName: "title",
			Op:        params.Like,
		},
		params.QueryFilter{
			ParamName: "status",
			Op:        params.Eq,
		},
	).Asc("sort_no").Desc("id"))
	ginx.WriteJSON(ctx, &web.PageResult{Results: list, Page: paging})

}

func BadgeCreate(ctx *gin.Context) {
	t := &models.Badge{}
	if err := ginx.Bind(ctx, t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}

	now := dates.NowTimestamp()
	t.CreateTime = now
	t.UpdateTime = now
	t.Status = constants.StatusOk
	if err := services.BadgeService.Create(t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func BadgeUpdate(ctx *gin.Context) {
	id, _ := params.GetInt64(ctx, "id")
	t := services.BadgeService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("entity not found"))
		return
	}

	if err := ginx.Bind(ctx, t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}

	t.UpdateTime = dates.NowTimestamp()
	if err := services.BadgeService.Update(t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func BadgeRemove(ctx *gin.Context) {
	ids := params.GetInt64Arr(ctx, "ids")
	if len(ids) == 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("delete ids is empty"))
		return
	}
	now := dates.NowTimestamp()
	for _, id := range ids {
		services.BadgeService.Updates(id, map[string]interface{}{
			"status":      constants.StatusDeleted,
			"update_time": now,
		})
	}
	ginx.WriteJSON(ctx, nil)

}
