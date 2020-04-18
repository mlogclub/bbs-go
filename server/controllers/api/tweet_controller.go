package api

import (
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/controllers/render"
	"bbs-go/services"
)

type TweetController struct {
	Ctx iris.Context
}

func (c *TweetController) PostCreate() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	content := strings.TrimSpace(simple.FormValue(c.Ctx, "content"))
	imageList := simple.FormValue(c.Ctx, "imageList")
	tweets, err := services.TweetService.Publish(user.Id, content, imageList)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(render.BuildTweet(tweets))
}

func (c *TweetController) GetList() *simple.JsonResult {
	cursor := simple.FormValueInt64Default(c.Ctx, "cursor", 0)
	tweets, cursor := services.TweetService.GetTweets(cursor)
	return simple.JsonCursorData(render.BuildTweets(tweets), strconv.FormatInt(cursor, 10))
}
