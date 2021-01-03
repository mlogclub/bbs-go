package render

import (
	"bbs-go/common/markdown"
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/services"
	"github.com/mlogclub/simple"
)

func BuildTopic(topic *model.Topic) *model.TopicResponse {
	if topic == nil {
		return nil
	}

	rsp := &model.TopicResponse{}

	rsp.TopicId = topic.Id
	rsp.Type = topic.Type
	rsp.Title = topic.Title
	rsp.User = BuildUserDefaultIfNull(topic.UserId)
	rsp.LastCommentTime = topic.LastCommentTime
	rsp.CreateTime = topic.CreateTime
	rsp.ViewCount = topic.ViewCount
	rsp.CommentCount = topic.CommentCount
	rsp.LikeCount = topic.LikeCount

	if topic.NodeId > 0 {
		node := services.TopicNodeService.Get(topic.NodeId)
		rsp.Node = BuildNode(node)
	}

	tags := services.TopicService.GetTopicTags(topic.Id)
	rsp.Tags = BuildTags(tags)

	content := markdown.ToHTML(topic.Content)
	rsp.Content = handleHtmlContent(content)

	return rsp
}

func BuildSimpleTopic(topic *model.Topic) *model.TopicSimpleResponse {
	if topic == nil {
		return nil
	}

	rsp := &model.TopicSimpleResponse{}

	rsp.TopicId = topic.Id
	rsp.Title = topic.Title
	rsp.User = BuildUserDefaultIfNull(topic.UserId)
	rsp.LastCommentTime = topic.LastCommentTime
	rsp.CreateTime = topic.CreateTime
	rsp.ViewCount = topic.ViewCount
	rsp.CommentCount = topic.CommentCount
	rsp.LikeCount = topic.LikeCount

	if topic.NodeId > 0 {
		node := services.TopicNodeService.Get(topic.NodeId)
		rsp.Node = BuildNode(node)
	}

	tags := services.TopicService.GetTopicTags(topic.Id)
	rsp.Tags = BuildTags(tags)
	return rsp
}

func BuildSimpleTopics(topics []model.Topic, currentUser *model.User) []model.TopicSimpleResponse {
	if topics == nil || len(topics) == 0 {
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

	var responses []model.TopicSimpleResponse
	for _, topic := range topics {
		var (
			liked = simple.Contains(topic.Id, likedTopicIds)
			item  = BuildSimpleTopic(&topic)
		)
		item.Liked = liked
		responses = append(responses, *item)
	}
	return responses
}
