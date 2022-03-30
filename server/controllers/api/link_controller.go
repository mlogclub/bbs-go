package api

import (
	"bbs-go/model/constants"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/model"
	"bbs-go/services"
)

type LinkController struct {
	Ctx iris.Context
}

func (c *LinkController) GetBy(id int64) *mvc.JsonResult {
	link := services.LinkService.Get(id)
	if link == nil || link.Status == constants.StatusDeleted {
		return mvc.JsonErrorMsg("数据不存在")
	}
	return mvc.JsonData(c.buildLink(*link))
}

// 列表
func (c *LinkController) GetLinks() *mvc.JsonResult {
	page := params.FormValueIntDefault(c.Ctx, "page", 1)

	links, paging := services.LinkService.FindPageByCnd(sqls.NewSqlCnd().
		Eq("status", constants.StatusOk).Page(page, 20).Asc("id"))

	var itemList []map[string]interface{}
	for _, v := range links {
		itemList = append(itemList, c.buildLink(v))
	}
	return mvc.JsonPageData(itemList, paging)
}

// 前10个链接
func (c *LinkController) GetToplinks() *mvc.JsonResult {
	links := services.LinkService.Find(sqls.NewSqlCnd().
		Eq("status", constants.StatusOk).Limit(10).Asc("id"))

	var itemList []map[string]interface{}
	for _, v := range links {
		itemList = append(itemList, c.buildLink(v))
	}
	return mvc.JsonData(itemList)
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
