package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/services"
)

type TweetController struct {
	Ctx iris.Context
}

func (c *TweetController) GetBy(id int64) *simple.JsonResult {
	t := services.TweetService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *TweetController) AnyList() *simple.JsonResult {
	list, paging := services.TweetService.FindPageByParams(simple.NewQueryParams(c.Ctx).
		EqByReq("id").
		EqByReq("user_id").
		EqByReq("status").
		PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: render.BuildTweets(list), Page: paging})
}

func (c *TweetController) PostDelete() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.TweetService.UpdateColumn(id, "status", model.StatusDeleted)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

func (c *TweetController) PostUndelete() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.TweetService.UpdateColumn(id, "status", model.StatusOk)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}
