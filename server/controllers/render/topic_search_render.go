package render

import (
	"bbs-go/model"
	"bbs-go/package/es"
	"bbs-go/services"
)

func BuildSearchTopics(docs []es.TopicDocument) []model.SearchTopicResponse {
	var items []model.SearchTopicResponse
	for _, doc := range docs {
		items = append(items, BuildSearchTopic(doc))
	}
	return items
}

func BuildSearchTopic(doc es.TopicDocument) model.SearchTopicResponse {
	rsp := model.SearchTopicResponse{
		TopicId:    doc.Id,
		Tags:       nil,
		Title:      doc.Title,
		Summary:    doc.Content,
		CreateTime: doc.CreateTime,
	}

	rsp.User = BuildUserDefaultIfNull(doc.UserId)

	if doc.NodeId > 0 {
		node := services.TopicNodeService.Get(doc.NodeId)
		rsp.Node = BuildNode(node)
	}

	tags := services.TopicService.GetTopicTags(doc.Id)
	rsp.Tags = BuildTags(tags)
	return rsp
}
