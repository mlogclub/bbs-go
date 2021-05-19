package render

import "bbs-go/model"

func BuildNode(node *model.TopicNode) *model.NodeResponse {
	if node == nil {
		return nil
	}
	return &model.NodeResponse{
		NodeId:      node.Id,
		Name:        node.Name,
		Logo:        node.Logo,
		Description: node.Description,
	}
}

func BuildNodes(nodes []model.TopicNode) []model.NodeResponse {
	if len(nodes) == 0 {
		return nil
	}
	var ret []model.NodeResponse
	for _, node := range nodes {
		ret = append(ret, *BuildNode(&node))
	}
	return ret
}
