package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/services"
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
	list, paging := services.TopicService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "id",
			Op:        params.Eq,
			ValueWrapper: func(origin string) string {
				if id := idcodec.Decode(origin); id > 0 {
					return strconv.FormatInt(id, 10)
				}
				return ""
			},
		},
		params.QueryFilter{
			ParamName: "userId",
			Op:        params.Eq,
			ValueWrapper: func(origin string) string {
				if id := idcodec.Decode(origin); id > 0 {
					return strconv.FormatInt(id, 10)
				}
				return ""
			},
		},
		params.QueryFilter{
			ParamName: "status",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "type",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "qaStatus",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "recommend",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "title",
			Op:        params.Like,
		},
	).Desc("id"))

	var results []map[string]interface{}
	for _, topic := range list {
		item := render.BuildSimpleTopic(&topic)
		builder := web.NewRspBuilder(item)
		builder.Put("id", topic.Id) // Admin 端保持使用明文 ID，避免后台管理接口参数不兼容
		builder.Put("idEncode", idcodec.Encode(topic.Id))
		builder.Put("status", topic.Status)
		if vote := services.VoteService.Get(topic.VoteId); vote != nil {
			builder.Put("vote", render.BuildVote(c.Ctx, vote))
		}
		results = append(results, builder.Build())
	}

	return web.JsonData(&web.PageResult{Results: results, Page: paging})
}

// 推荐
func (c *TopicController) PostRecommend() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	err = services.TopicService.SetRecommend(id, true)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 取消推荐
func (c *TopicController) DeleteRecommend() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	err = services.TopicService.SetRecommend(id, false)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *TopicController) PostDelete() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	user := common.GetCurrentUser(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin())
	}
	err = services.TopicService.Delete(id, user.Id, c.Ctx.Request())
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *TopicController) PostUndelete() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	err = services.TopicService.Undelete(id)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *TopicController) PostAudit() *web.JsonResult {
	id := c.Ctx.PostValueInt64Default("id", 0)
	if id <= 0 {
		return web.JsonErrorMsg("id is required")
	}
	err := services.TopicService.UpdateColumn(id, "status", constants.StatusOk)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *TopicController) PostAccept_answer() *web.JsonResult {
	topicId := c.Ctx.PostValueInt64Default("id", 0)
	commentId := c.Ctx.PostValueInt64Default("commentId", 0)
	if topicId <= 0 || commentId <= 0 {
		return web.JsonErrorMsg("id/commentId is required")
	}
	user := common.GetCurrentUser(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin())
	}
	if err := services.TopicService.AcceptAnswer(topicId, commentId, user.Id, true); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *TopicController) PostUnaccept_answer() *web.JsonResult {
	topicId := c.Ctx.PostValueInt64Default("id", 0)
	if topicId <= 0 {
		return web.JsonErrorMsg("id is required")
	}
	user := common.GetCurrentUser(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin())
	}
	if err := services.TopicService.UnacceptAnswer(topicId, user.Id, true); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *TopicController) PostMark_solved() *web.JsonResult {
	topicId := c.Ctx.PostValueInt64Default("id", 0)
	if topicId <= 0 {
		return web.JsonErrorMsg("id is required")
	}
	if err := services.TopicService.ForceSetQaStatus(topicId, constants.QaStatusSolved); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *TopicController) PostMark_unsolved() *web.JsonResult {
	topicId := c.Ctx.PostValueInt64Default("id", 0)
	if topicId <= 0 {
		return web.JsonErrorMsg("id is required")
	}
	if err := services.TopicService.ForceSetQaStatus(topicId, constants.QaStatusUnsolved); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
