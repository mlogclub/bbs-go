package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/controllers/render"
	"bbs-go/services"
)

type TopicController struct {
	Ctx iris.Context
}

func (c *TopicController) GetBy(id int64) *mvc.JsonResult {
	t := services.TopicService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return mvc.JsonData(t)
}

func (c *TopicController) AnyList() *mvc.JsonResult {
	list, paging := services.TopicService.FindPageByParams(params.NewQueryParams(c.Ctx).
		EqByReq("id").EqByReq("user_id").EqByReq("status").EqByReq("recommend").LikeByReq("title").PageByReq().Desc("id"))

	var results []map[string]interface{}
	for _, topic := range list {
		item := render.BuildSimpleTopic(&topic)
		builder := mvc.NewRspBuilder(item)
		builder.Put("status", topic.Status)
		results = append(results, builder.Build())
	}

	return mvc.JsonData(&sqls.PageResult{Results: results, Page: paging})
}

// 推荐
func (c *TopicController) PostRecommend() *mvc.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	err = services.TopicService.SetRecommend(id, true)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}

// 取消推荐
func (c *TopicController) DeleteRecommend() *mvc.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	err = services.TopicService.SetRecommend(id, false)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}

func (c *TopicController) PostDelete() *mvc.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}
	err = services.TopicService.Delete(id, user.Id, c.Ctx.Request())
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}

func (c *TopicController) PostUndelete() *mvc.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	err = services.TopicService.Undelete(id)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}
