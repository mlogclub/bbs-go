package admin

import (
	"bbs-go/internal/models/constants"
	"strconv"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/web"

	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

func LinkDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	t := services.LinkService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("Not found, id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func LinkList(ctx *gin.Context) {
	list, paging := services.LinkService.FindPageByCnd(params.NewPagedSqlCnd(ctx,
		params.QueryFilter{
			ParamName: "status",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "title",
			Op:        params.Like,
		},
		params.QueryFilter{
			ParamName: "url",
			Op:        params.Like,
		},
	))
	ginx.WriteJSON(ctx, &web.PageResult{Results: list, Page: paging})

}

func LinkCreate(ctx *gin.Context) {
	t := &models.Link{}
	err := ginx.Bind(ctx, t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	t.CreateTime = dates.NowTimestamp()
	t.Status = constants.StatusOk

	err = services.LinkService.Create(t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, t)

}

func LinkRemove(ctx *gin.Context) {
	ids := params.GetInt64Arr(ctx, "ids")
	if len(ids) == 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("delete ids is empty"))
		return
	}
	for _, id := range ids {
		services.LinkService.Delete(id)
	}
	ginx.WriteJSON(ctx, nil)

}

func LinkUpdate(ctx *gin.Context) {
	id, err := params.FormValueInt64(ctx, "id")
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	t := services.LinkService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("entity not found"))
		return
	}

	err = ginx.Bind(ctx, t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	err = services.LinkService.Update(t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, t)

}
