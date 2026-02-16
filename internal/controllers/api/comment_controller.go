package api

import (
	"bbs-go/internal/models/req"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/spam"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/controllers/render"
	"bbs-go/internal/services"
)

type CommentController struct {
	Ctx iris.Context
}

func (c *CommentController) GetComments() *web.JsonResult {
	var (
		cursor, _     = params.GetInt64(c.Ctx, "cursor")
		entityType, _ = params.Get(c.Ctx, "entityType")
		entityId      = common.GetID(c.Ctx, "entityId")
		currentUser   = common.GetCurrentUser(c.Ctx)
	)
	comments, cursor, hasMore := services.CommentService.GetComments(entityType, entityId, cursor)
	return web.JsonCursorData(render.BuildComments(comments, currentUser, true, false), strconv.FormatInt(cursor, 10), hasMore)
}

func (c *CommentController) GetReplies() *web.JsonResult {
	var (
		cursor, _    = params.GetInt64(c.Ctx, "cursor")
		commentId, _ = params.GetInt64(c.Ctx, "commentId")
	)
	currentUser := common.GetCurrentUser(c.Ctx)
	comments, cursor, hasMore := services.CommentService.GetReplies(commentId, cursor, 10)
	return web.JsonCursorData(render.BuildComments(comments, currentUser, false, true), strconv.FormatInt(cursor, 10), hasMore)
}

func (c *CommentController) PostCreate() *web.JsonResult {
	user := common.GetCurrentUser(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}
	form := req.GetCreateCommentForm(c.Ctx)
	if err := spam.CheckComment(user, form); err != nil {
		return web.JsonError(err)
	}

	comment, err := services.CommentService.Publish(user.Id, form)
	if err != nil {
		return web.JsonError(err)
	}

	return web.JsonData(render.BuildComment(comment))
}
