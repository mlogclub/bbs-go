package api

import (
	"bbs-go/internal/models/req"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/spam"
	"strconv"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"bbs-go/internal/handlers/render"
	"bbs-go/internal/services"
	"github.com/mlogclub/simple/web"
)

func CommentComments(ctx *gin.Context) {
	var (
		cursor, _     = params.GetInt64(ctx, "cursor")
		entityType, _ = params.Get(ctx, "entityType")
		entityId      = common.GetID(ctx, "entityId")
		currentUser   = common.GetCurrentUser(ctx)
	)
	comments, cursor, hasMore := services.CommentService.GetComments(entityType, entityId, cursor)
	ginx.WriteJSON(ctx, ginx.CursorData(render.BuildComments(comments, currentUser, true, false), strconv.FormatInt(cursor, 10), hasMore))

}

func CommentReplies(ctx *gin.Context) {
	var (
		cursor, _    = params.GetInt64(ctx, "cursor")
		commentId, _ = params.GetInt64(ctx, "commentId")
	)
	currentUser := common.GetCurrentUser(ctx)
	comments, cursor, hasMore := services.CommentService.GetReplies(commentId, cursor, 10)
	ginx.WriteJSON(ctx, ginx.CursorData(render.BuildComments(comments, currentUser, false, true), strconv.FormatInt(cursor, 10), hasMore))

}

func CommentCreate(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	var body req.CreateCommentReq
	if err := ginx.Bind(ctx, &body); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	body.UserAgent = web.GetUserAgent(ctx.Request)
	body.Ip = web.GetRequestIP(ctx.Request)
	if err := spam.CheckComment(user, body); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	comment, err := services.CommentService.Publish(user.Id, body)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	ginx.WriteJSON(ctx, render.BuildComment(comment))

}

func CommentRemove(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	user := common.GetCurrentUser(ctx)
	if err := services.CommentService.DeleteByUser(user, id); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}
