package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/controllers/render"
	"bbs-go/pkg/common"
	"bbs-go/services"
)

type TopicController struct {
	Ctx iris.Context
}

func (c *TopicController) GetBy(id int64) *web.JsonResult {
	t := services.TopicService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *TopicController) AnyList() *web.JsonResult {
	list, paging := services.TopicService.FindPageByParams(params.NewQueryParams(c.Ctx).
		EqByReq("id").EqByReq("user_id").EqByReq("status").EqByReq("recommend").LikeByReq("title").PageByReq().Desc("id"))

	var results []map[string]interface{}
	for _, topic := range list {
		item := render.BuildSimpleTopic(&topic)
		builder := web.NewRspBuilder(item)
		builder.Put("status", topic.Status)
		results = append(results, builder.Build())
	}

	return web.JsonData(&web.PageResult{Results: results, Page: paging})
}

// 推荐
func (c *TopicController) PostRecommend() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	err = services.TopicService.SetRecommend(id, true)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonSuccess()
}

// 取消推荐
func (c *TopicController) DeleteRecommend() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	err = services.TopicService.SetRecommend(id, false)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonSuccess()
}

func (c *TopicController) PostDelete() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(common.ErrorNotLogin)
	}
	err = services.TopicService.Delete(id, user.Id, c.Ctx.Request())
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonSuccess()
}

func (c *TopicController) PostUndelete() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	err = services.TopicService.Undelete(id)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonSuccess()
}
