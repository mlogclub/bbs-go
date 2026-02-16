package render

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/resp"
)

func BuildNode(node *models.TopicNode) *resp.NodeResponse {
	if node == nil {
		return nil
	}
	return &resp.NodeResponse{
		Id:          node.Id,
		Name:        node.Name,
		Logo:        node.Logo,
		Description: node.Description,
	}
}

func BuildNodes(nodes []models.TopicNode) []resp.NodeResponse {
	if len(nodes) == 0 {
		return nil
	}
	var ret []resp.NodeResponse
	for _, node := range nodes {
		ret = append(ret, *BuildNode(&node))
	}
	return ret
}
