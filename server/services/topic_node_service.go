package services

import (
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/repositories"
)

var TopicNodeService = newTopicNodeService()

func newTopicNodeService() *topicNodeService {
	return &topicNodeService{}
}

type topicNodeService struct {
}

func (this *topicNodeService) Get(id int64) *model.TopicNode {
	return repositories.TopicNodeRepository.Get(simple.DB(), id)
}

func (this *topicNodeService) Take(where ...interface{}) *model.TopicNode {
	return repositories.TopicNodeRepository.Take(simple.DB(), where...)
}

func (this *topicNodeService) Find(cnd *simple.SqlCnd) []model.TopicNode {
	return repositories.TopicNodeRepository.Find(simple.DB(), cnd)
}

func (this *topicNodeService) FindOne(cnd *simple.SqlCnd) *model.TopicNode {
	return repositories.TopicNodeRepository.FindOne(simple.DB(), cnd)
}

func (this *topicNodeService) FindPageByParams(params *simple.QueryParams) (list []model.TopicNode, paging *simple.Paging) {
	return repositories.TopicNodeRepository.FindPageByParams(simple.DB(), params)
}

func (this *topicNodeService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.TopicNode, paging *simple.Paging) {
	return repositories.TopicNodeRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *topicNodeService) Create(t *model.TopicNode) error {
	return repositories.TopicNodeRepository.Create(simple.DB(), t)
}

func (this *topicNodeService) Update(t *model.TopicNode) error {
	return repositories.TopicNodeRepository.Update(simple.DB(), t)
}

func (this *topicNodeService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.TopicNodeRepository.Updates(simple.DB(), id, columns)
}

func (this *topicNodeService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.TopicNodeRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (this *topicNodeService) Delete(id int64) {
	repositories.TopicNodeRepository.Delete(simple.DB(), id)
}

func (this *topicNodeService) GetNodes() []model.TopicNode {
	return repositories.TopicNodeRepository.Find(simple.DB(), simple.NewSqlCnd().Eq("status", model.StatusOk).Asc("sort_no").Desc("id"))
}
