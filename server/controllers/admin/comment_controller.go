package admin

import (
	"bbs-go/pkg/markdown"
	"strconv"

	"bbs-go/controllers/render"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"

	"bbs-go/services"
)

type CommentController struct {
	Ctx iris.Context
}

func (c *CommentController) GetBy(id int64) *mvc.JsonResult {
	t := services.CommentService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return mvc.JsonData(t)
}

func (c *CommentController) AnyList() *mvc.JsonResult {
	var (
		id         = params.FormValueInt64Default(c.Ctx, "id", 0)
		userId     = params.FormValueInt64Default(c.Ctx, "userId", 0)
		entityType = params.FormValueDefault(c.Ctx, "entityType", "")
		entityId   = params.FormValueInt64Default(c.Ctx, "entityId", 0)
	)
	params := params.NewQueryParams(c.Ctx).
		EqByReq("status").
		PageByReq().Desc("id")

	if id > 0 {
		params.Eq("id", id)
	}
	if userId > 0 {
		params.Eq("user_id", userId)
	}
	if strs.IsNotBlank(entityType) && entityId > 0 {
		params.Eq("entity_type", entityType).Eq("entity_id", entityId)
	}

	if id <= 0 && userId <= 0 && (strs.IsBlank(entityType) || entityId <= 0) {
		// return mvc.JsonErrorMsg("请输入必要的查询参数。")
		return mvc.JsonSuccess()
	}

	list, paging := services.CommentService.FindPageByParams(params)

	var results []map[string]interface{}
	for _, comment := range list {
		builder := mvc.NewRspBuilderExcludes(comment, "content")

		// 用户
		builder = builder.Put("user", render.BuildUserInfoDefaultIfNull(comment.UserId))

		// 简介
		content := markdown.ToHTML(comment.Content)
		builder.Put("content", content)

		results = append(results, builder.Build())
	}

	return mvc.JsonPageData(results, paging)
}

func (c *CommentController) PostDeleteBy(id int64) *mvc.JsonResult {
	if err := services.CommentService.Delete(id); err != nil {
		return mvc.JsonErrorMsg(err.Error())
	} else {
		return mvc.JsonSuccess()
	}
}
