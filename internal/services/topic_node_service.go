package services

import (
	"errors"

	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/locales"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"gorm.io/gorm"
)

var TopicNodeService = newTopicNodeService()

func newTopicNodeService() *topicNodeService {
	return &topicNodeService{}
}

type topicNodeService struct {
}

func (s *topicNodeService) Get(id int64) *models.TopicNode {
	return repositories.TopicNodeRepository.Get(sqls.DB(), id)
}

func (s *topicNodeService) Take(where ...interface{}) *models.TopicNode {
	return repositories.TopicNodeRepository.Take(sqls.DB(), where...)
}

func (s *topicNodeService) Find(cnd *sqls.Cnd) []models.TopicNode {
	return repositories.TopicNodeRepository.Find(sqls.DB(), cnd)
}

func (s *topicNodeService) FindOne(cnd *sqls.Cnd) *models.TopicNode {
	return repositories.TopicNodeRepository.FindOne(sqls.DB(), cnd)
}

func (s *topicNodeService) FindPageByParams(params *params.QueryParams) (list []models.TopicNode, paging *sqls.Paging) {
	return repositories.TopicNodeRepository.FindPageByParams(sqls.DB(), params)
}

func (s *topicNodeService) FindPageByCnd(cnd *sqls.Cnd) (list []models.TopicNode, paging *sqls.Paging) {
	return repositories.TopicNodeRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *topicNodeService) Create(t *models.TopicNode) error {
	return repositories.TopicNodeRepository.Create(sqls.DB(), t)
}

func (s *topicNodeService) Update(t *models.TopicNode) error {
	return repositories.TopicNodeRepository.Update(sqls.DB(), t)
}

func (s *topicNodeService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.TopicNodeRepository.Updates(sqls.DB(), id, columns)
}

func (s *topicNodeService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.TopicNodeRepository.UpdateColumn(sqls.DB(), id, name, value)
}

// DeleteWithCheck 删除节点，若为一级且有子节点则返回错误
func (s *topicNodeService) DeleteWithCheck(id int64) error {
	node := s.Get(id)
	if node == nil {
		return nil
	}
	if node.ParentId == 0 {
		children := s.GetChildren(id)
		if len(children) > 0 {
			return errors.New(locales.Get("topic.node.has_children"))
		}
	}
	return repositories.TopicNodeRepository.Updates(sqls.DB(), id, map[string]interface{}{
		"status": constants.StatusDeleted,
	})
}

// GetTopLevelNodes 仅一级节点（parent_id=0），用于导航
func (s *topicNodeService) GetTopLevelNodes() []models.TopicNode {
	return repositories.TopicNodeRepository.Find(sqls.DB(), sqls.NewCnd().
		Eq("status", constants.StatusOk).
		Eq("parent_id", 0).
		Asc("sort_no").Desc("id"))
}

// GetChildren 获取某一级下的二级节点
func (s *topicNodeService) GetChildren(parentId int64) []models.TopicNode {
	return repositories.TopicNodeRepository.Find(sqls.DB(), sqls.NewCnd().
		Eq("status", constants.StatusOk).
		Eq("parent_id", parentId).
		Asc("sort_no").Desc("id"))
}

// GetNodeIdsForList 用于帖子列表筛选：一级返回 [自身+子节点id]，二级返回 [自身]
func (s *topicNodeService) GetNodeIdsForList(nodeId int64) []int64 {
	node := s.Get(nodeId)
	if node == nil {
		return nil
	}
	if node.ParentId == 0 {
		ids := []int64{nodeId}
		for _, c := range s.GetChildren(nodeId) {
			ids = append(ids, c.Id)
		}
		return ids
	}
	return []int64{nodeId}
}

func (s *topicNodeService) GetNodes() []models.TopicNode {
	return repositories.TopicNodeRepository.Find(sqls.DB(), sqls.NewCnd().Eq("status", constants.StatusOk).Asc("sort_no").Desc("id"))
}

func (s *topicNodeService) GetNodesByType(nodeType constants.TopicNodeType) []models.TopicNode {
	return repositories.TopicNodeRepository.Find(sqls.DB(), sqls.NewCnd().
		Eq("status", constants.StatusOk).
		Eq("type", nodeType).
		Asc("sort_no").Desc("id"))
}

func (s *topicNodeService) GetNodesByTopicType(topicType constants.TopicType) []models.TopicNode {
	if topicType == constants.TopicTypeQA {
		return s.GetNodesByType(constants.TopicNodeTypeQA)
	}
	return s.GetNodesByType(constants.TopicNodeTypeNormal)
}

func (s *topicNodeService) GetNextSortNo() int {
	if max := s.FindOne(sqls.NewCnd().Eq("status", constants.StatusOk).Desc("sort_no")); max != nil {
		return max.SortNo + 1
	}
	return 0
}

func (s *topicNodeService) UpdateSort(ids []int64) error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := repositories.TopicNodeRepository.UpdateColumn(tx, id, "sort_no", i); err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateChildrenType 将父节点下所有子节点的 type 更新为指定值（父节点编辑类型时联动）
func (s *topicNodeService) UpdateChildrenType(parentId int64, nodeType constants.TopicNodeType) error {
	children := s.GetChildren(parentId)
	if len(children) == 0 {
		return nil
	}
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		for _, c := range children {
			if err := repositories.TopicNodeRepository.UpdateColumn(tx, c.Id, "type", nodeType); err != nil {
				return err
			}
		}
		return nil
	})
}
