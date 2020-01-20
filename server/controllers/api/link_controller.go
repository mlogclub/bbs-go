package api

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"bbs-go/common"
	"bbs-go/model"
	"bbs-go/services"
)

type LinkController struct {
	Ctx iris.Context
}

func (c *LinkController) GetBy(id int64) *simple.JsonResult {
	link := services.LinkService.Get(id)
	if link == nil || link.Status == model.StatusDeleted {
		return simple.JsonErrorMsg("数据不存在")
	}
	return simple.JsonData(c.buildLink(*link))
}

// 列表
func (c *LinkController) GetLinks() *simple.JsonResult {
	page := simple.FormValueIntDefault(c.Ctx, "page", 1)

	links, paging := services.LinkService.FindPageByCnd(simple.NewSqlCnd().
		Eq("status", model.StatusOk).Page(page, 20).Desc("id"))

	var itemList []map[string]interface{}
	for _, v := range links {
		itemList = append(itemList, c.buildLink(v))
	}
	return simple.JsonPageData(itemList, paging)
}

// 待审核
func (c *LinkController) GetPending() *simple.JsonResult {
	links := services.LinkService.Find(simple.NewSqlCnd().
		Eq("status", model.StatusPending).Limit(3).Desc("id"))

	var itemList []map[string]interface{}
	for _, v := range links {
		itemList = append(itemList, c.buildLink(v))
	}
	return simple.JsonData(itemList)
}

func (c *LinkController) PostCreate() *simple.JsonResult {
	var (
		title   = c.Ctx.FormValue("title")
		url     = c.Ctx.FormValue("url")
		summary = c.Ctx.FormValue("summary")
		logo    = c.Ctx.FormValue("logo")
	)
	if err := common.IsValidateUrl(url); err != nil {
		return simple.JsonErrorMsg("博客链接错误")
	}
	if len(title) == 0 {
		return simple.JsonErrorMsg("标题不能为空")
	}
	if len(summary) == 0 {
		return simple.JsonErrorMsg("描述不能为空")
	}

	userId := services.UserTokenService.GetCurrentUserId(c.Ctx)
	link := &model.Link{
		UserId:     userId,
		Url:        url,
		Title:      title,
		Summary:    summary,
		Logo:       logo,
		Status:     model.StatusPending,
		CreateTime: simple.NowTimestamp(),
	}
	err := services.LinkService.Create(link)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 通过网址检测标题和描述
func (c *LinkController) PostDetect() *simple.JsonResult {
	url := c.Ctx.FormValue("url")
	if err := common.IsValidateUrl(url); err != nil {
		logrus.Error(err.Error(), url)
		return simple.JsonSuccess()
	}
	resp, err := resty.New().SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).R().Get(url)
	if err != nil {
		logrus.Error(err)
		return simple.JsonSuccess()
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		logrus.Error(err)
		return simple.JsonSuccess()
	}
	title := doc.Find("title").Text()
	description := doc.Find("meta[name=description]").AttrOr("content", "")
	return simple.NewEmptyRspBuilder().Put("title", title).Put("description", description).JsonResult()
}

func (c *LinkController) buildLink(link model.Link) map[string]interface{} {
	return map[string]interface{}{
		"linkId":     link.Id,
		"url":        link.Url,
		"title":      link.Title,
		"summary":    link.Summary,
		"logo":       link.Logo,
		"category":   link.Category,
		"createTime": link.CreateTime,
	}
}
