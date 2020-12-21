package render

import (
	"bbs-go/common/markdown"
	"bbs-go/model"
	"bbs-go/services"
)

func BuildTopic(topic *model.Topic) *model.TopicResponse {
	if topic == nil {
		return nil
	}

	rsp := &model.TopicResponse{}

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

func BuildSimpleTopics(topics []model.Topic) []model.TopicSimpleResponse {
	if topics == nil || len(topics) == 0 {
		return nil
	}
	var responses []model.TopicSimpleResponse
	for _, topic := range topics {
		responses = append(responses, *BuildSimpleTopic(&topic))
	}
	return responses
}
