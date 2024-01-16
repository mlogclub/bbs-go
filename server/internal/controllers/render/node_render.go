package render

import "bbs-go/internal/models"

func BuildNode(node *models.TopicNode) *models.NodeResponse {
	if node == nil {
		return nil
	}
	return &models.NodeResponse{
		Id:          node.Id,
		Name:        node.Name,
		Logo:        node.Logo,
		Description: node.Description,
	}
}

func BuildNodes(nodes []models.TopicNode) []models.NodeResponse {
	if len(nodes) == 0 {
		return nil
	}
	var ret []models.NodeResponse
	for _, node := range nodes {
		ret = append(ret, *BuildNode(&node))
	}
	return ret
}
