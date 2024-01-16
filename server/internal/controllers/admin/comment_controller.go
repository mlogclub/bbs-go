package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/markdown"
	"strconv"

	"bbs-go/internal/controllers/render"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/services"
)

type CommentController struct {
	Ctx iris.Context
}

func (c *CommentController) GetBy(id int64) *web.JsonResult {
	t := services.CommentService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *CommentController) AnyList() *web.JsonResult {
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
		// return web.JsonErrorMsg("请输入必要的查询参数。")
		return web.JsonSuccess()
	}

	list, paging := services.CommentService.FindPageByParams(params)

	var results []map[string]interface{}
	for _, comment := range list {
		builder := web.NewRspBuilder(comment)

		// 用户
		builder = builder.Put("user", render.BuildUserInfoDefaultIfNull(comment.UserId))

		// 内容
		if comment.ContentType == constants.ContentTypeMarkdown {
			builder.Put("content", markdown.ToHTML(comment.Content))
		} else {
			builder.Put("content", comment.Content)
		}

		// 图片
		if strs.IsNotBlank(comment.ImageList) {
			var imageList []models.ImageDTO
			_ = jsons.Parse(comment.ImageList, &imageList)
			builder.Put("imageList", imageList)
		}

		results = append(results, builder.Build())
	}

	return web.JsonPageData(results, paging)
}

func (c *CommentController) PostDeleteBy(id int64) *web.JsonResult {
	if err := services.CommentService.Delete(id); err != nil {
		return web.JsonError(err)
	} else {
		return web.JsonSuccess()
	}
}
