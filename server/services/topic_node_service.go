package services

import (
	"bbs-go/model/constants"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/model"
	"bbs-go/repositories"
)

var TopicNodeService = newTopicNodeService()

func newTopicNodeService() *topicNodeService {
	return &topicNodeService{}
}

type topicNodeService struct {
}

func (s *topicNodeService) Get(id int64) *model.TopicNode {
	return repositories.TopicNodeRepository.Get(sqls.DB(), id)
}

func (s *topicNodeService) Take(where ...interface{}) *model.TopicNode {
	return repositories.TopicNodeRepository.Take(sqls.DB(), where...)
}

func (s *topicNodeService) Find(cnd *sqls.Cnd) []model.TopicNode {
	return repositories.TopicNodeRepository.Find(sqls.DB(), cnd)
}

func (s *topicNodeService) FindOne(cnd *sqls.Cnd) *model.TopicNode {
	return repositories.TopicNodeRepository.FindOne(sqls.DB(), cnd)
}

func (s *topicNodeService) FindPageByParams(params *params.QueryParams) (list []model.TopicNode, paging *sqls.Paging) {
	return repositories.TopicNodeRepository.FindPageByParams(sqls.DB(), params)
}

func (s *topicNodeService) FindPageByCnd(cnd *sqls.Cnd) (list []model.TopicNode, paging *sqls.Paging) {
	return repositories.TopicNodeRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *topicNodeService) Create(t *model.TopicNode) error {
	return repositories.TopicNodeRepository.Create(sqls.DB(), t)
}

func (s *topicNodeService) Update(t *model.TopicNode) error {
	return repositories.TopicNodeRepository.Update(sqls.DB(), t)
}

func (s *topicNodeService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.TopicNodeRepository.Updates(sqls.DB(), id, columns)
}

func (s *topicNodeService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.TopicNodeRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *topicNodeService) Delete(id int64) {
	repositories.TopicNodeRepository.Delete(sqls.DB(), id)
}

func (s *topicNodeService) GetNodes() []model.TopicNode {
	return repositories.TopicNodeRepository.Find(sqls.DB(), sqls.NewCnd().Eq("status", constants.StatusOk).Asc("sort_no").Desc("id"))
}
