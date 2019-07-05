package services

import (
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/simple"
)

var TopicTagService = newTopicTagService()

func newTopicTagService() *topicTagService {
	return &topicTagService{
		TopicTagRepository: repositories.NewTopicTagRepository(),
	}
}

type topicTagService struct {
	TopicTagRepository *repositories.TopicTagRepository
}

func (this *topicTagService) Get(id int64) *model.TopicTag {
	return this.TopicTagRepository.Get(simple.GetDB(), id)
}

func (this *topicTagService) Take(where ...interface{}) *model.TopicTag {
	return this.TopicTagRepository.Take(simple.GetDB(), where...)
}

func (this *topicTagService) QueryCnd(cnd *simple.QueryCnd) (list []model.TopicTag, err error) {
	return this.TopicTagRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *topicTagService) Query(queries *simple.ParamQueries) (list []model.TopicTag, paging *simple.Paging) {
	return this.TopicTagRepository.Query(simple.GetDB(), queries)
}

func (this *topicTagService) Create(t *model.TopicTag) error {
	return this.TopicTagRepository.Create(simple.GetDB(), t)
}

func (this *topicTagService) Update(t *model.TopicTag) error {
	return this.TopicTagRepository.Update(simple.GetDB(), t)
}

func (this *topicTagService) Updates(id int64, columns map[string]interface{}) error {
	return this.TopicTagRepository.Updates(simple.GetDB(), id, columns)
}

func (this *topicTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.TopicTagRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *topicTagService) Delete(id int64) {
	this.TopicTagRepository.Delete(simple.GetDB(), id)
}
