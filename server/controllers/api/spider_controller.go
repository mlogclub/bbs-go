package api

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
	"github.com/mlogclub/bbs-go/services/collect"
)

type SpiderController struct {
	Ctx iris.Context
}

// 微信采集发布接口
func (this *SpiderController) PostWxPublish() *simple.JsonResult {
	err := this.checkToken()
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	article := &collect.WxArticle{}
	err = this.Ctx.ReadJSON(article)
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
func (this *SpiderController) PostArticlePublish() *simple.JsonResult {
	err := this.checkToken()
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	article := &collect.SpiderArticle{}
	err = this.Ctx.ReadJSON(article)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	articleId, err := collect.NewSpiderApi().Publish(article)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("id", articleId).JsonResult()
}

func (this *SpiderController) PostCommentPublish() *simple.JsonResult {
	err := this.checkToken()
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	comment := &collect.SpiderComment{}
	err = this.Ctx.ReadJSON(comment)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	commentId, err := collect.NewSpiderApi().PublishComment(comment)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("id", commentId).JsonResult()
}

func (this *SpiderController) PostProjectPublish() *simple.JsonResult {
	err := this.checkToken()
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	var (
		userIdStr   = this.Ctx.FormValue("userId")
		userId, _   = strconv.ParseInt(userIdStr, 10, 64)
		name        = this.Ctx.FormValue("name")
		title       = this.Ctx.FormValue("title")
		logo        = this.Ctx.FormValue("logo")
		url         = this.Ctx.FormValue("url")
		docUrl      = this.Ctx.FormValue("docUrl")
		downloadUrl = this.Ctx.FormValue("downloadUrl")
		content     = this.Ctx.FormValue("content")
	)

	if len(name) == 0 || len(title) == 0 || len(content) == 0 {
		return simple.JsonErrorMsg("数据不完善...")
	}

	temp := services.ProjectService.FindOne(simple.NewSqlCnd().Eq("name", name))
	if temp != nil {
		return simple.JsonErrorMsg("项目已经存在：" + name)
	}

	p, err := services.ProjectService.Publish(userId, name, title, logo, url, docUrl, downloadUrl,
		model.ContentTypeHtml, content)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("id", p.Id).JsonResult()
}

func (this *SpiderController) checkToken() error {
	token := this.Ctx.FormValue("token")
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
