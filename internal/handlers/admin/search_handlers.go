package admin

import (
	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/services"

	"github.com/gin-gonic/gin"
)

func SearchReindex(ctx *gin.Context) {
	status, _ := services.SearchReindexService.Start()
	ginx.WriteJSON(ctx, status)
}

func SearchReindexStatus(ctx *gin.Context) {
	ginx.WriteJSON(ctx, services.SearchReindexService.Status())
}
