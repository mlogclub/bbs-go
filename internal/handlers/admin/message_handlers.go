package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/web"

	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

func MessageDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	t := services.MessageService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("Not found, id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func MessageList(ctx *gin.Context) {
	list, paging := services.MessageService.FindPageByParams(params.NewQueryParams(ctx).PageByReq().Desc("id"))
	ginx.WriteJSON(ctx, &web.PageResult{Results: list, Page: paging})

}

func MessageCreate(ctx *gin.Context) {
	t := &models.Message{}
	if err := ginx.Bind(ctx, t); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	err := services.MessageService.Create(t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, t)

}

func MessageUpdate(ctx *gin.Context) {
	id, err := params.FormValueInt64(ctx, "id")
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	t := services.MessageService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("entity not found"))
		return
	}

	if err := ginx.Bind(ctx, t); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	err = services.MessageService.Update(t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, t)

}
