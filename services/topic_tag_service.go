package services

import (
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/simple"
)

var TopicTagService = newTopicTagService()

func newTopicTagService() *topicTagService {
	return &topicTagService{}
}

type topicTagService struct {
}

func (this *topicTagService) Get(id int64) *model.TopicTag {
	return repositories.TopicTagRepository.Get(simple.GetDB(), id)
}

func (this *topicTagService) Take(where ...interface{}) *model.TopicTag {
	return repositories.TopicTagRepository.Take(simple.GetDB(), where...)
}

func (this *topicTagService) QueryCnd(cnd *simple.QueryCnd) (list []model.TopicTag, err error) {
	return repositories.TopicTagRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *topicTagService) Query(queries *simple.ParamQueries) (list []model.TopicTag, paging *simple.Paging) {
	return repositories.TopicTagRepository.Query(simple.GetDB(), queries)
}

func (this *topicTagService) Create(t *model.TopicTag) error {
	return repositories.TopicTagRepository.Create(simple.GetDB(), t)
}

func (this *topicTagService) Update(t *model.TopicTag) error {
	return repositories.TopicTagRepository.Update(simple.GetDB(), t)
}

func (this *topicTagService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.TopicTagRepository.Updates(simple.GetDB(), id, columns)
}

func (this *topicTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.TopicTagRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *topicTagService) Delete(id int64) {
	repositories.TopicTagRepository.Delete(simple.GetDB(), id)
}
