package api

import (
	"bbs-go/model/constants"
	"bbs-go/package/collect"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/services"
)

type SpiderController struct {
	Ctx iris.Context
}

// 微信采集发布接口
func (c *SpiderController) PostWxPublish() *simple.JsonResult {
	err := c.checkToken()
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	article := &collect.WxArticle{}
	err = c.Ctx.ReadJSON(article)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t, err := collect.NewWxbotApi().Publish(article)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("id", t.Id).JsonResult()
}

// 采集发布
func (c *SpiderController) PostArticlePublish() *simple.JsonResult {
	err := c.checkToken()
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	article := &collect.Article{}
	err = c.Ctx.ReadJSON(article)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	articleId, err := collect.NewSpiderApi().Publish(article)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("id", articleId).JsonResult()
}

func (c *SpiderController) PostCommentPublish() *simple.JsonResult {
	err := c.checkToken()
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	comment := &collect.Comment{}
	err = c.Ctx.ReadJSON(comment)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	commentId, err := collect.NewSpiderApi().PublishComment(comment)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("id", commentId).JsonResult()
}

func (c *SpiderController) PostProjectPublish() *simple.JsonResult {
	err := c.checkToken()
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	var (
		userIdStr   = c.Ctx.FormValue("userId")
		userId, _   = strconv.ParseInt(userIdStr, 10, 64)
		name        = c.Ctx.FormValue("name")
		title       = c.Ctx.FormValue("title")
		logo        = c.Ctx.FormValue("logo")
		url         = c.Ctx.FormValue("url")
		docUrl      = c.Ctx.FormValue("docUrl")
		downloadUrl = c.Ctx.FormValue("downloadUrl")
		content     = c.Ctx.FormValue("content")
		contentType = c.Ctx.FormValue("contentType")
	)

	if len(name) == 0 || len(title) == 0 || len(content) == 0 {
		return simple.JsonErrorMsg("数据不完善...")
	}

	temp := services.ProjectService.FindOne(simple.NewSqlCnd().Eq("name", name))
	if temp != nil {
		return simple.JsonErrorMsg("项目已经存在：" + name)
	}

	if len(contentType) == 0 {
		contentType = constants.ContentTypeHtml
	}

	p, err := services.ProjectService.Publish(userId, name, title, logo, url, docUrl, downloadUrl,
		contentType, content)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("id", p.Id).JsonResult()
}

func (c *SpiderController) checkToken() error {
	token := c.Ctx.FormValue("token")
	data, err := ioutil.ReadFile("/data/publish_token")
	if err != nil {
		return err
	}
	token2 := strings.TrimSpace(string(data))
	if token != token2 {
		return errors.New("token invalidate")
	}
	return nil
}
