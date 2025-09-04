package api

import (
	"bbs-go/internal/models/constants"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/pkg/simple/sqls"
	"bbs-go/internal/pkg/simple/web"

	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

type LinkController struct {
	Ctx iris.Context
}

// 列表
func (c *LinkController) GetList() *web.JsonResult {
	links := services.LinkService.Find(sqls.NewCnd().
		Eq("status", constants.StatusOk).Asc("id"))

	var itemList []map[string]any
	for _, v := range links {
		itemList = append(itemList, c.buildLink(v))
	}
	return web.JsonData(itemList)
}

// 前10个链接
func (c *LinkController) GetTop_links() *web.JsonResult {
	links := services.LinkService.Find(sqls.NewCnd().
		Eq("status", constants.StatusOk).Limit(10).Asc("id"))

	var itemList []map[string]any
	for _, v := range links {
		itemList = append(itemList, c.buildLink(v))
	}
	return web.JsonData(itemList)
}

func (c *LinkController) buildLink(link models.Link) map[string]any {
	return map[string]any{
		"id":         link.Id,
		"linkId":     link.Id,
		"url":        link.Url,
		"title":      link.Title,
		"summary":    link.Summary,
		"createTime": link.CreateTime,
	}
}
