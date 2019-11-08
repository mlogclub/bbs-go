package admin

import (
	"bytes"
	"strconv"

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
	t := services.LinkService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *LinkController) AnyList() *simple.JsonResult {
	list, paging := services.LinkService.FindPageByParams(simple.NewQueryParams(this.Ctx).EqByReq("status").LikeByReq("title").LikeByReq("url").PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *LinkController) PostCreate() *simple.JsonResult {
	t := &model.Link{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.LinkService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *LinkController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.LinkService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err = this.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.LinkService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *LinkController) GetDetect() *simple.JsonResult {
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
