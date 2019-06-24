package services

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/simple"
)

type TopicService struct {
	TopicRepository    *repositories.TopicRepository
	TagRepository      *repositories.TagRepository
	TopicTagRepository *repositories.TopicTagRepository
}

func NewTopicService() *TopicService {
	return &TopicService{
		TopicRepository:    repositories.NewTopicRepository(),
		TagRepository:      repositories.NewTagRepository(),
		TopicTagRepository: repositories.NewTopicTagRepository(),
	}
}

func (this *TopicService) Get(id int64) *model.Topic {
	return this.TopicRepository.Get(simple.GetDB(), id)
}

func (this *TopicService) Take(where ...interface{}) *model.Topic {
	return this.TopicRepository.Take(simple.GetDB(), where...)
}

func (this *TopicService) QueryCnd(cnd *simple.QueryCnd) (list []model.Topic, err error) {
	return this.TopicRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *TopicService) Query(queries *simple.ParamQueries) (list []model.Topic, paging *simple.Paging) {
	return this.TopicRepository.Query(simple.GetDB(), queries)
}

func (this *TopicService) Create(t *model.Topic) error {
	return this.TopicRepository.Create(simple.GetDB(), t)
}

func (this *TopicService) Update(t *model.Topic) error {
	return this.TopicRepository.Update(simple.GetDB(), t)
}

func (this *TopicService) Updates(id int64, columns map[string]interface{}) error {
	return this.TopicRepository.Updates(simple.GetDB(), id, columns)
}

func (this *TopicService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.TopicRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *TopicService) Delete(id int64) {
	this.TopicRepository.Delete(simple.GetDB(), id)
}

func (this *TopicService) Publish(userId int64, tags []string, title, content string) (*model.Topic, *simple.CodeError) {
	if len(title) == 0 {
		return nil, simple.NewErrorMsg("标题不能为空")
	}

	if simple.RuneLen(title) > 128 {
		return nil, simple.NewErrorMsg("标题长度不能超过128")
	}

	now := simple.NowTimestamp()
	topic := &model.Topic{
		UserId:          userId,
		Title:           title,
		Content:         content,
		Status:          model.TopicStatusOk,
		LastCommentTime: now,
		CreateTime:      now,
	}

	err := simple.Tx(simple.GetDB(), func(tx *gorm.DB) error {
		tagIds := this.TagRepository.GetOrCreates(tx, tags)
		err := this.TopicRepository.Create(simple.GetDB(), topic)
		if err != nil {
			return err
		}

		this.TopicTagRepository.AddTopicTags(tx, topic.Id, tagIds)
		return nil
	})
	return topic, simple.NewError2(err)
}

func (this *TopicService) GetTopicTags(topicId int64) []model.Tag {
	topicTags, err := this.TopicTagRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("topic_id = ?", topicId))
	if err != nil {
		return nil
	}

	var tagIds []int64
	for _, topicTag := range topicTags {
		tagIds = append(tagIds, topicTag.TagId)
	}

	return this.TagRepository.GetTagInIds(tagIds)
}

func (this *TopicService) IncrViewCount(topicId int64) {
	simple.GetDB().Exec("update t_topic set view_count = view_count + 1 where id = ?", topicId)
}
