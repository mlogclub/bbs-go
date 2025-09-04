package render

import (
	"html"

	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	html2 "bbs-go/internal/pkg/html"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/pkg/text"
	"bbs-go/internal/services"

	"bbs-go/internal/pkg/simple/common/arrays"
	"bbs-go/internal/pkg/simple/common/strs"
)

func BuildTopic(topic *models.Topic, currentUser *models.User) *models.TopicResponse {
	resp := _buildTopic(topic, true)
	if currentUser != nil {
		resp.Liked = services.UserLikeService.Exists(currentUser.Id, constants.EntityTopic, topic.Id)
		resp.Favorited = services.FavoriteService.IsFavorited(currentUser.Id, constants.EntityTopic, topic.Id)
	}
	return resp
}

func BuildSimpleTopic(topic *models.Topic) *models.TopicResponse {
	buildContent := topic.Type == constants.TopicTypeTweet // 动态时渲染内容
	return _buildTopic(topic, buildContent)
}

func BuildSimpleTopics(topics []models.Topic, currentUser *models.User) []models.TopicResponse {
	if len(topics) == 0 {
		return nil
	}

	var likedTopicIds []int64
	if currentUser != nil {
		var topicIds []int64
		for _, topic := range topics {
			topicIds = append(topicIds, topic.Id)
		}
		likedTopicIds = services.UserLikeService.IsLiked(currentUser.Id, constants.EntityTopic, topicIds)
	}

	var responses []models.TopicResponse
	for _, topic := range topics {
		item := BuildSimpleTopic(&topic)
		item.Liked = arrays.Contains(topic.Id, likedTopicIds)
		responses = append(responses, *item)
	}
	return responses
}

func _buildTopic(topic *models.Topic, buildContent bool) *models.TopicResponse {
	if topic == nil {
		return nil
	}

	rsp := &models.TopicResponse{}

	rsp.Id = topic.Id
	rsp.Type = topic.Type
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
		if topic.Type == constants.TopicTypeTopic {
			contentHtml := topic.Content
			if topic.ContentType == constants.ContentTypeMarkdown {
				contentHtml = markdown.ToHTML(topic.Content)
			}
			rsp.Content = handleHtmlContent(contentHtml)
		} else {
			rsp.Content = html.EscapeString(topic.Content)
		}
	} else {
		if topic.Type == constants.TopicTypeTopic {
			contentHtml := topic.Content
			if topic.ContentType == constants.ContentTypeMarkdown {
				contentHtml = markdown.ToHTML(topic.Content)
			}
			rsp.Summary = html2.GetSummary(contentHtml, 128)
		} else {
			rsp.Summary = text.GetSummary(topic.Content, 128)
		}
	}

	if topic.Type == constants.TopicTypeTweet {
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
