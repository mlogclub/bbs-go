package admin

import (
	"bbs-go/internal/handlers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
)

func DictDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	t := services.DictService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("Not found, id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, render.BuildDict(*t))

}

func DictList(ctx *gin.Context) {
	typeId, _ := params.GetInt64(ctx, "typeId")
	list := services.DictService.Find(sqls.NewCnd().Eq("type_id", typeId).Asc("sort_no").Desc("id"))
	ginx.WriteJSON(ctx, render.BuildDictTree(0, list))

}

func DictCreate(ctx *gin.Context) {
	t := &models.Dict{}
	if err := ginx.Bind(ctx, t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}

	t.SortNo = services.DictService.GetNextSortNo()
	t.CreateTime = dates.NowTimestamp()
	t.UpdateTime = dates.NowTimestamp()
	if err := services.DictService.Create(t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func DictUpdate(ctx *gin.Context) {
	id, err := params.FormValueInt64(ctx, "id")
	if err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	t := services.DictService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("entity not found"))
		return
	}

	if err := ginx.Bind(ctx, t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}

	t.UpdateTime = dates.NowTimestamp()
	if err := services.DictService.Update(t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func DictRemove(ctx *gin.Context) {
	ids := params.GetInt64Arr(ctx, "ids")
	if len(ids) == 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("delete ids is empty"))
		return
	}
	for _, id := range ids {
		services.DictService.Delete(id)
	}
	ginx.WriteJSON(ctx, nil)

}

func DictUpdateSort(ctx *gin.Context) {
	var ids []int64
	if err := ginx.BindJSON(ctx, &ids); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	if err := services.DictService.UpdateSort(ids); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func DictDicts(ctx *gin.Context) {
	var (
		typeId, _ = params.GetInt64(ctx, "typeId")
		code, _   = params.Get(ctx, "code")
	)
	if typeId <= 0 {
		if strs.IsNotBlank(code) {
			dictType := services.DictTypeService.GetByCode(code)
			if dictType != nil {
				typeId = dictType.Id
			}
		}
	}
	list := services.DictService.FindByTypeId(typeId)
	ginx.WriteJSON(ctx, render.BuildDictTree(0, list))

}
