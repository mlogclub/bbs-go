package services

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
)

type ScanTopicCallback func(topics []model.Topic)

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

// 扫描
func (this *TopicService) Scan(cb ScanTopicCallback) {
	var cursor int64
	for {
		list, err := this.TopicRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("id > ?",
			cursor).Order("id asc").Size(300))
		if err != nil {
			break
		}
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		cb(list)
	}
}

// 发表
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

// 更新
func (this *TopicService) Edit(topicId int64, tags []string, title, content string) *simple.CodeError {
	if len(title) == 0 {
		return simple.NewErrorMsg("标题不能为空")
	}

	if simple.RuneLen(title) > 128 {
		return simple.NewErrorMsg("标题长度不能超过128")
	}

	err := simple.Tx(simple.GetDB(), func(tx *gorm.DB) error {
		tagIds := this.TagRepository.GetOrCreates(tx, tags)
		err := this.TopicRepository.Updates(simple.GetDB(), topicId, map[string]interface{}{
			"title":   title,
			"content": content,
		})
		if err != nil {
			return err
		}
		this.TopicTagRepository.RemoveTopicTags(tx, topicId)      // 先删掉所有的标签
		this.TopicTagRepository.AddTopicTags(tx, topicId, tagIds) // 然后重新添加标签
		return nil
	})
	return simple.NewError2(err)
}

// 帖子标签
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

// 浏览数+1
func (this *TopicService) IncrViewCount(topicId int64) {
	simple.GetDB().Exec("update t_topic set view_count = view_count + 1 where id = ?", topicId)
}

// 更新最后回复时间
func (this *TopicService) SetLastCommentTime(topicId, lastCommentTime int64) {
	err := this.UpdateColumn(topicId, "last_comment_time", lastCommentTime)
	if err != nil {
		logrus.Error(err)
	}
}
