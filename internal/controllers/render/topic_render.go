package render

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/resp"
	"bbs-go/internal/pkg/common"
	html2 "bbs-go/internal/pkg/html"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/pkg/text"
	"bbs-go/internal/services"
	"html"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/arrays"
	"github.com/mlogclub/simple/common/strs"
)

func BuildTopic(ctx iris.Context, topic *models.Topic) *resp.TopicResponse {
	rsp := _buildTopic(topic, true)
	if rsp == nil {
		return nil
	}

	if currentUser := common.GetCurrentUser(ctx); currentUser != nil {
		rsp.Liked = services.UserLikeService.Exists(currentUser.Id, constants.EntityTopic, topic.Id)
		rsp.Favorited = services.FavoriteService.IsFavorited(currentUser.Id, constants.EntityTopic, topic.Id)
	}

	if vote := services.VoteService.Get(topic.VoteId); vote != nil {
		rsp.Vote = BuildVote(ctx, vote)
	}

	// 附件仅在帖子详情接口返回。
	list := services.AttachmentService.ListByTopicId(topic.Id)
	if len(list) > 0 {
		var currentUser *models.User
		if u := common.GetCurrentUser(ctx); u != nil {
			currentUser = u
		}
		rsp.Attachments = BuildAttachmentResponses(list, currentUser)
	}

	return rsp
}

// BuildAttachmentResponses 将附件列表转为 AttachmentResponse 列表；currentUser 为 nil 时 downloaded 均为 false（如编辑表单）
func BuildAttachmentResponses(list []models.Attachment, currentUser *models.User) []resp.AttachmentResponse {
	if len(list) == 0 {
		return nil
	}
	atts := make([]resp.AttachmentResponse, 0, len(list))
	downloadedMap := make(map[string]bool)
	if currentUser != nil && len(list) > 0 {
		attachmentIds := make([]string, 0, len(list))
		for _, att := range list {
			attachmentIds = append(attachmentIds, att.Id)
		}
		for _, attachmentId := range services.AttachmentService.FindDownloadedAttachmentIds(currentUser.Id, attachmentIds) {
			downloadedMap[attachmentId] = true
		}
	}

	for _, att := range list {
		atts = append(atts, resp.AttachmentResponse{
			Id:            att.Id,
			FileName:      att.FileName,
			FileSize:      att.FileSize,
			DownloadScore: att.DownloadScore,
			DownloadCount: att.DownloadCount,
			Downloaded:    downloadedMap[att.Id],
		})
	}
	return atts
}

func BuildSimpleTopic(topic *models.Topic) *resp.TopicResponse {
	buildContent := constants.IsTweetTopicType(topic.Type) // 动态时渲染内容
	return _buildTopic(topic, buildContent)
}

func BuildSimpleTopics(ctx iris.Context, topics []models.Topic) []resp.TopicResponse {
	if len(topics) == 0 {
		return nil
	}

	var likedTopicIds []int64
	if currentUser := common.GetCurrentUser(ctx); currentUser != nil {
		var topicIds []int64
		for _, topic := range topics {
			topicIds = append(topicIds, topic.Id)
		}
		likedTopicIds = services.UserLikeService.IsLiked(currentUser.Id, constants.EntityTopic, topicIds)
	}

	var responses []resp.TopicResponse
	for _, topic := range topics {
		item := BuildSimpleTopic(&topic)
		item.Liked = arrays.Contains(topic.Id, likedTopicIds)
		if vote := services.VoteService.Get(topic.VoteId); vote != nil {
			item.Vote = BuildVote(ctx, vote)
		}
		responses = append(responses, *item)
	}
	return responses
}

func _buildTopic(topic *models.Topic, buildContent bool) *resp.TopicResponse {
	if topic == nil {
		return nil
	}

	rsp := &resp.TopicResponse{}

	rsp.Id = idcodec.Encode(topic.Id)
	rsp.Type = topic.Type
	rsp.QaStatus = topic.QaStatus
	rsp.AcceptedCommentId = topic.AcceptedCommentId
	rsp.SolvedAt = topic.SolvedAt
	rsp.BountyScore = topic.BountyScore
	rsp.Title = topic.Title
	rsp.User = BuildUserInfoDefaultIfNull(topic.UserId)
	rsp.LastCommentTime = topic.LastCommentTime
	rsp.CreateTime = topic.CreateTime
	rsp.ViewCount = topic.ViewCount
	rsp.CommentCount = topic.CommentCount
	rsp.LikeCount = topic.LikeCount
	rsp.Recommend = topic.Recommend
	rsp.RecommendTime = topic.RecommendTime
	rsp.Sticky = topic.Sticky
	rsp.StickyTime = topic.StickyTime
	rsp.Status = topic.Status
	rsp.IpLocation = topic.IpLocation

	// 构建内容
	if buildContent {
		if !constants.IsTweetTopicType(topic.Type) {
			contentHtml := topic.Content
			if topic.ContentType == constants.ContentTypeMarkdown {
				contentHtml = markdown.ToHTML(topic.Content)
			}
			rsp.Content = handleHtmlContent(contentHtml)
		} else {
			rsp.Content = html.EscapeString(topic.Content)
		}
	} else {
		if !constants.IsTweetTopicType(topic.Type) {
			contentHtml := topic.Content
			if topic.ContentType == constants.ContentTypeMarkdown {
				contentHtml = markdown.ToHTML(topic.Content)
			}
			rsp.Summary = html2.GetSummary(contentHtml, 128)
		} else {
			rsp.Summary = text.GetSummary(topic.Content, 128)
		}
	}

	if constants.IsTweetTopicType(topic.Type) {
		if strs.IsBlank(topic.Content) {
			rsp.Content = "分享图片"
		} else {
			rsp.Content = html.EscapeString(topic.Content)
		}
		rsp.ImageList = BuildImageList(topic.ImageList)
	}

	if topic.NodeId > 0 {
		node := services.TopicNodeService.Get(topic.NodeId)
		rsp.Node = BuildNode(node)
	}

	tags := services.TopicService.GetTopicTags(topic.Id)
	rsp.Tags = BuildTags(tags)

	return rsp
}
