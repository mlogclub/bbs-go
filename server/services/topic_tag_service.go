package services

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var TopicTagService = newTopicTagService()

func newTopicTagService() *topicTagService {
	return &topicTagService{}
}

type topicTagService struct {
}

func (this *topicTagService) Get(id int64) *model.TopicTag {
	return repositories.TopicTagRepository.Get(simple.DB(), id)
}

func (this *topicTagService) Take(where ...interface{}) *model.TopicTag {
	return repositories.TopicTagRepository.Take(simple.DB(), where...)
}

func (this *topicTagService) Find(cnd *simple.SqlCnd) []model.TopicTag {
	return repositories.TopicTagRepository.Find(simple.DB(), cnd)
}

func (this *topicTagService) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.TopicTag) {
	cnd.FindOne(db, &ret)
	return
}

func (this *topicTagService) FindPageByParams(params *simple.QueryParams) (list []model.TopicTag, paging *simple.Paging) {
	return repositories.TopicTagRepository.FindPageByParams(simple.DB(), params)
}

func (this *topicTagService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.TopicTag, paging *simple.Paging) {
	return repositories.TopicTagRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *topicTagService) Create(t *model.TopicTag) error {
	return repositories.TopicTagRepository.Create(simple.DB(), t)
}

func (this *topicTagService) Update(t *model.TopicTag) error {
	return repositories.TopicTagRepository.Update(simple.DB(), t)
}

func (this *topicTagService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.TopicTagRepository.Updates(simple.DB(), id, columns)
}

func (this *topicTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.TopicTagRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (this *topicTagService) DeleteByTopicId(topicId int64) {
	simple.DB().Model(model.TopicTag{}).Where("topic_id = ?", topicId).UpdateColumn("status", model.TopicTagStatusDeleted)
}
