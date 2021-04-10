package admin

import (
	"bbs-go/package/markdown"
	"strconv"

	"bbs-go/controllers/render"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/services"
)

type CommentController struct {
	Ctx iris.Context
}

func (c *CommentController) GetBy(id int64) *simple.JsonResult {
	t := services.CommentService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *CommentController) AnyList() *simple.JsonResult {
	var (
		id         = simple.FormValueInt64Default(c.Ctx, "id", 0)
		userId     = simple.FormValueInt64Default(c.Ctx, "userId", 0)
		entityType = simple.FormValueDefault(c.Ctx, "entityType", "")
		entityId   = simple.FormValueInt64Default(c.Ctx, "entityId", 0)
	)
	params := simple.NewQueryParams(c.Ctx).
		EqByReq("status").
		PageByReq().Desc("id")

	if id > 0 {
		params.Eq("id", id)
	}
	if userId > 0 {
		params.Eq("user_id", userId)
	}
	if simple.IsNotBlank(entityType) && entityId > 0 {
		params.Eq("entity_type", entityType).Eq("entity_id", entityId)
	}

	if id <= 0 && userId <= 0 && (simple.IsBlank(entityType) || entityId <= 0) {
		// return simple.JsonErrorMsg("请输入必要的查询参数。")
		return simple.JsonSuccess()
	}

	list, paging := services.CommentService.FindPageByParams(params)

	var results []map[string]interface{}
	for _, comment := range list {
		builder := simple.NewRspBuilderExcludes(comment, "content")

		// 用户
		builder = builder.Put("user", render.BuildUserDefaultIfNull(comment.UserId))

		// 简介
		content := markdown.ToHTML(comment.Content)
		builder.Put("content", content)

		results = append(results, builder.Build())
	}

	return simple.JsonPageData(results, paging)
}

func (c *CommentController) PostDeleteBy(id int64) *simple.JsonResult {
	if err := services.CommentService.Delete(id); err != nil {
		return simple.JsonErrorMsg(err.Error())
	} else {
		return simple.JsonSuccess()
	}
}
