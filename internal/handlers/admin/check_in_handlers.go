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

func CheckInDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	t := services.CheckInService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("Not found, id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func CheckInList(ctx *gin.Context) {
	list, paging := services.CheckInService.FindPageByParams(params.NewQueryParams(ctx).PageByReq().Desc("id"))
	ginx.WriteJSON(ctx, &web.PageResult{Results: list, Page: paging})

}

func CheckInCreate(ctx *gin.Context) {
	t := &models.CheckIn{}
	err := ginx.Bind(ctx, t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	err = services.CheckInService.Create(t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, t)

}

func CheckInUpdate(ctx *gin.Context) {
	id, err := params.FormValueInt64(ctx, "id")
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	t := services.CheckInService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("entity not found"))
		return
	}

	err = ginx.Bind(ctx, t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	err = services.CheckInService.Update(t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, t)

}
