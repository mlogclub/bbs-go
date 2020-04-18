package api

import (
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/controllers/render"
	"bbs-go/services"
)

type TweetsController struct {
	Ctx iris.Context
}

func (c *TweetsController) PostCreate() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	content := strings.TrimSpace(simple.FormValue(c.Ctx, "content"))
	imageList := simple.FormValue(c.Ctx, "imageList")
	tweets, err := services.TweetsService.Publish(user.Id, content, imageList)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(render.BuildTweets(tweets))
}

func (c *TweetsController) GetList() *simple.JsonResult {
	cursor := simple.FormValueInt64Default(c.Ctx, "cursor", 0)
	tweets, cursor := services.TweetsService.GetTweets(cursor)
	return simple.JsonCursorData(render.BuildTweetsList(tweets), strconv.FormatInt(cursor, 10))
}
