package render

import (
	"bbs-go/internal/models/resp"
	"bbs-go/internal/pkg/search"
	"bbs-go/internal/services"
)

func BuildSearchTopics(docs []search.TopicDocument) []resp.SearchTopicResponse {
	var items []resp.SearchTopicResponse
	for _, doc := range docs {
		items = append(items, BuildSearchTopic(doc))
	}
	return items
}

func BuildSearchTopic(doc search.TopicDocument) resp.SearchTopicResponse {
	rsp := resp.SearchTopicResponse{
		Id:         doc.Id,
		Title:      doc.Title,
		Summary:    doc.Content,
		CreateTime: doc.CreateTime,
		User:       BuildUserInfoDefaultIfNull(doc.UserId),
	}

	if doc.NodeId > 0 {
		node := services.TopicNodeService.Get(doc.NodeId)
		rsp.Node = BuildNode(node)
	}

	tags := services.TopicService.GetTopicTags(doc.Id)
	rsp.Tags = BuildTags(tags)
	return rsp
}
