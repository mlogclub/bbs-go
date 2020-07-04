package api

import (
	"bbs-go/model/constants"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/services"
)

type LinkController struct {
	Ctx iris.Context
}

func (c *LinkController) GetBy(id int64) *simple.JsonResult {
	link := services.LinkService.Get(id)
	if link == nil || link.Status == constants.StatusDeleted {
		return simple.JsonErrorMsg("数据不存在")
	}
	return simple.JsonData(c.buildLink(*link))
}

// 列表
func (c *LinkController) GetLinks() *simple.JsonResult {
	page := simple.FormValueIntDefault(c.Ctx, "page", 1)

	links, paging := services.LinkService.FindPageByCnd(simple.NewSqlCnd().
		Eq("status", constants.StatusOk).Page(page, 20).Asc("id"))

	var itemList []map[string]interface{}
	for _, v := range links {
		itemList = append(itemList, c.buildLink(v))
	}
	return simple.JsonPageData(itemList, paging)
}

// 前10个链接
func (c *LinkController) GetToplinks() *simple.JsonResult {
	links := services.LinkService.Find(simple.NewSqlCnd().
		Eq("status", constants.StatusOk).Limit(10).Asc("id"))

	var itemList []map[string]interface{}
	for _, v := range links {
		itemList = append(itemList, c.buildLink(v))
	}
	return simple.JsonData(itemList)
}

func (c *LinkController) buildLink(link model.Link) map[string]interface{} {
	return map[string]interface{}{
		"linkId":     link.Id,
		"url":        link.Url,
		"title":      link.Title,
		"summary":    link.Summary,
		"logo":       link.Logo,
		"createTime": link.CreateTime,
	}
}
