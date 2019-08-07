package admin

import (
	"strconv"

	"github.com/mlogclub/mlog/controllers/render"

	"github.com/kataras/iris"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/services"
)

type CommentController struct {
	Ctx iris.Context
}

func (this *CommentController) GetBy(id int64) *simple.JsonResult {
	t := services.CommentService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *CommentController) AnyList() *simple.JsonResult {
	list, paging := services.CommentService.Query(simple.NewParamQueries(this.Ctx).EqAuto("status").PageAuto().Desc("id"))

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

func (this *CommentController) PostDeleteBy(id int64) *simple.JsonResult {
	err := services.CommentService.Delete(id)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	} else {
		return simple.JsonSuccess()
	}
}
