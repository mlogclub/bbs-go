package api

import (
	"bbs-go/internal/models/constants"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"

	"github.com/mlogclub/simple/sqls"

	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

// 列表
// 前10个链接
func linkBuildLink(link models.Link) map[string]any {
	return map[string]any{
		"id":         link.Id,
		"linkId":     link.Id,
		"url":        link.Url,
		"title":      link.Title,
		"summary":    link.Summary,
		"createTime": link.CreateTime,
	}
}

func LinkList(ctx *gin.Context) {
	links := services.LinkService.Find(sqls.NewCnd().
		Eq("status", constants.StatusOk).Asc("sort_no").Desc("id"))

	var itemList []map[string]any
	for _, v := range links {
		itemList = append(itemList, linkBuildLink(v))
	}
	ginx.WriteJSON(ctx, itemList)

}

func LinkTopLinks(ctx *gin.Context) {
	links := services.LinkService.Find(sqls.NewCnd().
		Eq("status", constants.StatusOk).Limit(10).Asc("sort_no").Desc("id"))

	var itemList []map[string]any
	for _, v := range links {
		itemList = append(itemList, linkBuildLink(v))
	}
	ginx.WriteJSON(ctx, itemList)

}
