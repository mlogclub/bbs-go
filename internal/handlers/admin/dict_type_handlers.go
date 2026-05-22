package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/common/dates"
)

func DictTypeDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	t := services.DictTypeService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("Not found, id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func DictTypeList(ctx *gin.Context) {
	list := services.DictTypeService.Find(params.NewSqlCnd(ctx,
		params.QueryFilter{
			ParamName: "id",
		},
	).Desc("id"))
	ginx.WriteJSON(ctx, list)

}

func DictTypeCreate(ctx *gin.Context) {
	t := &models.DictType{}
	if err := ginx.Bind(ctx, t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}

	t.CreateTime = dates.NowTimestamp()
	t.UpdateTime = dates.NowTimestamp()
	if err := services.DictTypeService.Create(t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func DictTypeUpdate(ctx *gin.Context) {
	id, err := params.FormValueInt64(ctx, "id")
	if err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	t := services.DictTypeService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("entity not found"))
		return
	}

	if err := ginx.Bind(ctx, t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}

	t.UpdateTime = dates.NowTimestamp()
	if err := services.DictTypeService.Update(t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func DictTypeRemove(ctx *gin.Context) {
	ids := params.GetInt64Arr(ctx, "ids")
	if len(ids) == 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("delete ids is empty"))
		return
	}
	for _, id := range ids {
		services.DictTypeService.Delete(id)
	}
	ginx.WriteJSON(ctx, nil)

}
