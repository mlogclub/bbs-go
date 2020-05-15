package api

import (
	"bbs-go/model"
	"bbs-go/services"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"strconv"
)

type PollAnswerController struct {
	Ctx iris.Context
}

func (c *PollAnswerController) GetBy(id int64) *simple.JsonResult {
	t := services.PollAnswerService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

// Get is user voted or not
func (c *PollAnswerController) GetVoted() *simple.JsonResult {
	cnd := simple.NewSqlCnd("id").Eq("topic_id", c.Ctx.FormValue("topicId")).Eq("poll_user_id", c.Ctx.FormValue("pollUserId"))
	t := services.PollAnswerService.Count(cnd)
	return simple.JsonData(t != 0)
}

func (c *PollAnswerController) AnyList() *simple.JsonResult {
	list, paging := services.PollAnswerService.FindPageByParams(simple.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (c *PollAnswerController) PostCreate() *simple.JsonResult {
	t := &model.PollAnswer{}
	err := simple.ReadForm(c.Ctx, t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.PollAnswerService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (c *PollAnswerController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.PollAnswerService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err = simple.ReadForm(c.Ctx, t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.PollAnswerService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
