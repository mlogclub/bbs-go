package services

import (
	"bbs-go/model/constants"
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/repositories"
)

var TopicTagService = newTopicTagService()

func newTopicTagService() *topicTagService {
	return &topicTagService{}
}

type topicTagService struct {
}

func (s *topicTagService) Get(id int64) *model.TopicTag {
	return repositories.TopicTagRepository.Get(simple.DB(), id)
}

func (s *topicTagService) Take(where ...interface{}) *model.TopicTag {
	return repositories.TopicTagRepository.Take(simple.DB(), where...)
}

func (s *topicTagService) Find(cnd *simple.SqlCnd) []model.TopicTag {
	return repositories.TopicTagRepository.Find(simple.DB(), cnd)
}

func (s *topicTagService) FindOne(cnd *simple.SqlCnd) *model.TopicTag {
	return repositories.TopicTagRepository.FindOne(simple.DB(), cnd)
}

func (s *topicTagService) FindPageByParams(params *simple.QueryParams) (list []model.TopicTag, paging *simple.Paging) {
	return repositories.TopicTagRepository.FindPageByParams(simple.DB(), params)
}

func (s *topicTagService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.TopicTag, paging *simple.Paging) {
	return repositories.TopicTagRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *topicTagService) Create(t *model.TopicTag) error {
	return repositories.TopicTagRepository.Create(simple.DB(), t)
}

func (s *topicTagService) Update(t *model.TopicTag) error {
	return repositories.TopicTagRepository.Update(simple.DB(), t)
}

func (s *topicTagService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.TopicTagRepository.Updates(simple.DB(), id, columns)
}

func (s *topicTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.TopicTagRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *topicTagService) DeleteByTopicId(topicId int64) {
	simple.DB().Model(model.TopicTag{}).Where("topic_id = ?", topicId).UpdateColumn("status", constants.StatusDeleted)
}

func (s *topicTagService) UndeleteByTopicId(topicId int64) {
	simple.DB().Model(model.TopicTag{}).Where("topic_id = ?", topicId).UpdateColumn("status", constants.StatusOk)
}
