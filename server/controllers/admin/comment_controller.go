package admin

import (
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
	list, paging := services.CommentService.FindPageByParams(simple.NewQueryParams(c.Ctx).EqByReq("status").PageByReq().Desc("id"))

	var results []map[string]interface{}
	for _, comment := range list {
		builder := simple.NewRspBuilderExcludes(comment, "content")

		// 用户
		builder = builder.Put("user", render.BuildUserDefaultIfNull(comment.UserId))

		// 简介
		mr := simple.NewMd().Run(comment.Content)
		builder.Put("content", mr.ContentHtml)

		results = append(results, builder.Build())
	}

	return simple.JsonData(&simple.PageResult{Results: results, Page: paging})
}

func (c *CommentController) PostDeleteBy(id int64) *simple.JsonResult {
	err := services.CommentService.Delete(id)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	} else {
		return simple.JsonSuccess()
	}
}
