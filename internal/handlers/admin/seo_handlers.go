package admin

import (
	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/services"

	"github.com/gin-gonic/gin"
)

func SeoSitemapGenerate(ctx *gin.Context) {
	status, _ := services.SeoSitemapService.StartGenerate()
	ginx.WriteJSON(ctx, status)
}

func SeoSitemapStatus(ctx *gin.Context) {
	ginx.WriteJSON(ctx, services.SeoSitemapService.Status())
}
