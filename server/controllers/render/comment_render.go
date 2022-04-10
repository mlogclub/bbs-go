package render

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/pkg/markdown"
	"bbs-go/services"
	"html"
	"strconv"

	"github.com/mlogclub/simple/common/arrays"
	"github.com/mlogclub/simple/web"
)

func BuildComment(comment *model.Comment) *model.CommentResponse {
	return doBuildComment(comment, nil, true, true)
}

func BuildComments(comments []model.Comment, currentUser *model.User, isBuildReplies, isBuildQuote bool) []model.CommentResponse {
	if len(comments) == 0 {
		return nil
	}

	likedCommentIds := getLikedCommentIds(comments, currentUser)

	var ret []model.CommentResponse
	for _, comment := range comments {
		item := doBuildComment(&comment, currentUser, isBuildReplies, isBuildQuote)
		item.Liked = arrays.Contains(comment.Id, likedCommentIds)
		ret = append(ret, *item)
	}
	return ret
}

func getLikedCommentIds(comments []model.Comment, currentUser *model.User) (likedCommentIds []int64) {
	if currentUser == nil || len(comments) == 0 {
		return
	}
	var commentIds []int64
	for _, comment := range comments {
		commentIds = append(commentIds, comment.Id)
	}
	likedCommentIds = services.UserLikeService.IsLiked(currentUser.Id, constants.EntityComment, commentIds)
	return
}

// doBuildComment 渲染评论
// isBuildReplies 是否渲染评论的二级回复，一级评论时要设置为true，其他时候都为false
// isBuildQuote 是否渲染评论的引用，二级回复时要设置为true，其他时候都为false
func doBuildComment(comment *model.Comment, currentUser *model.User, isBuildReplies, isBuildQuote bool) *model.CommentResponse {
	if comment == nil {
		return nil
	}

	ret := &model.CommentResponse{
		CommentId:    comment.Id,
		User:         BuildUserInfoDefaultIfNull(comment.UserId),
		EntityType:   comment.EntityType,
		EntityId:     comment.EntityId,
		QuoteId:      comment.QuoteId,
		LikeCount:    comment.LikeCount,
		CommentCount: comment.CommentCount,
		Status:       comment.Status,
		CreateTime:   comment.CreateTime,
	}

	if comment.Status == constants.StatusOk {
		if comment.ContentType == constants.ContentTypeMarkdown {
			content := markdown.ToHTML(comment.Content)
			ret.Content = handleHtmlContent(content)
		} else if comment.ContentType == constants.ContentTypeHtml {
			ret.Content = handleHtmlContent(comment.Content)
		} else {
			ret.Content = html.EscapeString(comment.Content)
		}
		ret.ImageList = buildImageList(comment.ImageList)
	} else {
		ret.Content = "内容已删除"
	}

	if isBuildReplies && comment.CommentCount > 0 {
		var repliesLimit int64 = 3
		replies, nextCursor, _ := services.CommentService.GetReplies(comment.Id, 0, int(repliesLimit))
		//var replyResults []model.CommentResponse
		//for _, reply := range replies {
		//	replyResults = append(replyResults, *doBuildComment(&reply, false, true))
		//}
		replyResults := BuildComments(replies, currentUser, false, true)
		ret.Replies = &web.CursorResult{
			Results: replyResults,
			Cursor:  strconv.FormatInt(nextCursor, 10),
			HasMore: comment.CommentCount > repliesLimit,
		}
	}

	if isBuildQuote && comment.QuoteId > 0 {
		quote := doBuildComment(services.CommentService.Get(comment.QuoteId), currentUser, false, false)
		if quote != nil {
			ret.Quote = quote
		}
	}

	return ret
}
