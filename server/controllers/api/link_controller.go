package api

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
)

type LinkController struct {
	Ctx iris.Context
}

func (this *LinkController) GetBy(id int64) *simple.JsonResult {
	link := services.LinkService.Get(id)
	if link == nil || link.Status == model.LinkStatusDeleted {
		return simple.JsonErrorMsg("数据不存在")
	}
	return simple.JsonData(this.buildLink(*link))
}

func (this *LinkController) GetLinks() *simple.JsonResult {
	page := simple.FormValueIntDefault(this.Ctx, "page", 1)

	links, paging := services.LinkService.FindPageByCnd(simple.NewSqlCnd().
		Eq("status", model.LinkStatusOk).Page(page, 20).Desc("id"))

	var itemList []map[string]interface{}
	for _, v := range links {
		itemList = append(itemList, this.buildLink(v))
	}
	return simple.JsonPageData(itemList, paging)
}

func (this *LinkController) PostCreate() *simple.JsonResult {
	var (
		title   = this.Ctx.FormValue("title")
		url     = this.Ctx.FormValue("url")
		summary = this.Ctx.FormValue("summary")
		logo    = this.Ctx.FormValue("logo")
	)
	if len(title) == 0 {
		return simple.JsonErrorMsg("标题不能为空")
	}
	if len(url) == 0 {
		return simple.JsonErrorMsg("链接不能为空")
	}
	if len(summary) == 0 {
		return simple.JsonErrorMsg("描述不能为空")
	}
	link := &model.Link{
		Url:        url,
		Title:      title,
		Summary:    summary,
		Logo:       logo,
		Status:     model.LinkStatusPending,
		CreateTime: simple.NowTimestamp(),
	}
	err := services.LinkService.Create(link)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 通过网址检测标题和描述
func (this *LinkController) PostDetect() *simple.JsonResult {
	url := this.Ctx.FormValue("url")
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

func (this *LinkController) buildLink(link model.Link) map[string]interface{} {
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
