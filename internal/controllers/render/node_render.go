package render

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/resp"
	"bbs-go/internal/services"

	"github.com/mlogclub/simple/common/strs"
)

func BuildNode(node *models.TopicNode) *resp.NodeResponse {
	if node == nil {
		return nil
	}
	if strs.IsBlank(node.Logo) {
		node.Logo = "/res/images/node_default.svg"
	}
	return &resp.NodeResponse{
		Id:          node.Id,
		ParentId:    node.ParentId,
		Name:        node.Name,
		Type:        node.Type,
		Logo:        node.Logo,
		Description: node.Description,
	}
}

// BuildNodeWithChildren 构建节点详情（一级节点时带子节点，用于列表页展示子分类）
func BuildNodeWithChildren(node *models.TopicNode) *resp.NodeResponse {
	r := BuildNode(node)
	if r == nil {
		return nil
	}
	if node.ParentId == 0 {
		children := services.TopicNodeService.GetChildren(node.Id)
		if len(children) > 0 {
			r.Children = BuildNodes(children)
		}
	}
	return r
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

// BuildNodeTree 将扁平节点列表构建为树形（用于前台节点选择，一级节点包含 children）
func BuildNodeTree(parentId int64, list []models.TopicNode) []resp.NodeResponse {
	var ret []resp.NodeResponse
	for _, node := range list {
		if node.ParentId != parentId {
			continue
		}
		item := BuildNode(&node)
		if item == nil {
			continue
		}
		children := BuildNodeTree(node.Id, list)
		if len(children) > 0 {
			item.Children = children
		}
		ret = append(ret, *item)
	}
	return ret
}

// BuildTopicNodeTree 将扁平节点列表构建为树形（用于后台列表，含 sortNo/status/createTime，children 始终为 slice）
func BuildTopicNodeTree(parentId int64, list []models.TopicNode) []resp.TopicNodeTreeItem {
	var ret []resp.TopicNodeTreeItem
	for _, node := range list {
		if node.ParentId == parentId {
			children := BuildTopicNodeTree(node.Id, list)
			logo := node.Logo
			if strs.IsBlank(logo) {
				logo = "/res/images/node_default.svg"
			}
			ret = append(ret, resp.TopicNodeTreeItem{
				Id:          node.Id,
				ParentId:    node.ParentId,
				Name:        node.Name,
				Type:        node.Type,
				Logo:        logo,
				Description: node.Description,
				SortNo:      node.SortNo,
				Status:      node.Status,
				CreateTime:  node.CreateTime,
				Children:    children, // 叶子节点为 []，保证 Arco Table 树形展示
			})
		}
	}
	return ret
}
