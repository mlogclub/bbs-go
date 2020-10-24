package admin

import (
	"bytes"
	"github.com/mlogclub/simple/date"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"bbs-go/model"
	"bbs-go/services"
)

type LinkController struct {
	Ctx iris.Context
}

func (c *LinkController) GetBy(id int64) *simple.JsonResult {
	t := services.LinkService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *LinkController) AnyList() *simple.JsonResult {
	list, paging := services.LinkService.FindPageByParams(simple.NewQueryParams(c.Ctx).EqByReq("status").LikeByReq("title").LikeByReq("url").PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (c *LinkController) PostCreate() *simple.JsonResult {
	t := &model.Link{}
	err := simple.ReadForm(c.Ctx, t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t.CreateTime = date.NowTimestamp()

	err = services.LinkService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (c *LinkController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.LinkService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err = simple.ReadForm(c.Ctx, t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.LinkService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (c *LinkController) GetDetect() *simple.JsonResult {
	url := c.Ctx.FormValue("url")
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
