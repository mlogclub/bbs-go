package admin

import (
	"bytes"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
	"github.com/sirupsen/logrus"

	"bbs-go/model"
	"bbs-go/services"
)

type LinkController struct {
	Ctx iris.Context
}

func (c *LinkController) GetBy(id int64) *web.JsonResult {
	t := services.LinkService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *LinkController) AnyList() *web.JsonResult {
	list, paging := services.LinkService.FindPageByParams(params.NewQueryParams(c.Ctx).EqByReq("status").LikeByReq("title").LikeByReq("url").PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *LinkController) PostCreate() *web.JsonResult {
	t := &model.Link{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t.CreateTime = dates.NowTimestamp()

	err = services.LinkService.Create(t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *LinkController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t := services.LinkService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	err = services.LinkService.Update(t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *LinkController) GetDetect() *web.JsonResult {
	url := c.Ctx.FormValue("url")
	resp, err := resty.New().SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).R().Get(url)
	if err != nil {
		logrus.Error(err)
		return web.JsonSuccess()
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		logrus.Error(err)
		return web.JsonSuccess()
	}
	title := doc.Find("title").Text()
	description := doc.Find("meta[name=description]").AttrOr("content", "")
	return web.NewEmptyRspBuilder().Put("title", title).Put("description", description).JsonResult()
}
