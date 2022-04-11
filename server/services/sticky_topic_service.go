package services

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/repositories"
	"errors"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var StickyTopicService = newStickyTopicService()

func newStickyTopicService() *stickyTopicService {
	return &stickyTopicService{}
}

type stickyTopicService struct {
}

func (s *stickyTopicService) Get(id int64) *model.StickyTopic {
	return repositories.StickyTopicRepository.Get(sqls.DB(), id)
}

func (s *stickyTopicService) Take(where ...interface{}) *model.StickyTopic {
	return repositories.StickyTopicRepository.Take(sqls.DB(), where...)
}

func (s *stickyTopicService) Find(cnd *sqls.Cnd) []model.StickyTopic {
	return repositories.StickyTopicRepository.Find(sqls.DB(), cnd)
}

func (s *stickyTopicService) FindOne(cnd *sqls.Cnd) *model.StickyTopic {
	return repositories.StickyTopicRepository.FindOne(sqls.DB(), cnd)
}

func (s *stickyTopicService) FindPageByParams(params *params.QueryParams) (list []model.StickyTopic, paging *sqls.Paging) {
	return repositories.StickyTopicRepository.FindPageByParams(sqls.DB(), params)
}

func (s *stickyTopicService) FindPageByCnd(cnd *sqls.Cnd) (list []model.StickyTopic, paging *sqls.Paging) {
	return repositories.StickyTopicRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *stickyTopicService) Count(cnd *sqls.Cnd) int64 {
	return repositories.StickyTopicRepository.Count(sqls.DB(), cnd)
}

func (s *stickyTopicService) Create(t *model.StickyTopic) error {
	return repositories.StickyTopicRepository.Create(sqls.DB(), t)
}

func (s *stickyTopicService) Update(t *model.StickyTopic) error {
	return repositories.StickyTopicRepository.Update(sqls.DB(), t)
}

func (s *stickyTopicService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.StickyTopicRepository.Updates(sqls.DB(), id, columns)
}

func (s *stickyTopicService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.StickyTopicRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *stickyTopicService) Delete(id int64) {
	repositories.StickyTopicRepository.Delete(sqls.DB(), id)
}

func (s *stickyTopicService) FindByTopicId(topicId int64) *model.StickyTopic {
	return repositories.StickyTopicRepository.FindOne(sqls.DB(), sqls.NewCnd().Where("topic_id = ?", topicId))
}

// GetStickyTopics 获取设备
func (s *stickyTopicService) GetStickyTopics(nodeId int64, limit int) []model.Topic {
	if nodeId < 0 {
		nodeId = 0
	}
	return repositories.TopicRepository.FindBySql(sqls.DB(), "select t.* from t_topic t left join t_sticky_topic st on st.topic_id = t.id "+
		" where st.node_id in (0, ?) order by st.id desc limit ?", nodeId, limit)
}

func (s *stickyTopicService) AddSticky(nodeId, topicId int64) error {
	if nodeId < 0 {
		nodeId = 0
	}
	topic := TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return errors.New("话题不存在")
	}
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.StickyTopic{}, "topic_id = ?", topicId).Error; err != nil {
			return err
		}
		return repositories.StickyTopicRepository.Create(tx, &model.StickyTopic{
			TopicId:    topicId,
			NodeId:     nodeId,
			CreateTime: dates.NowTimestamp(),
		})
	})
}

func (s *stickyTopicService) RemoveSticky(topicId int64) error {
	return sqls.DB().Delete(&model.StickyTopic{}, "topic_id = ?", topicId).Error
}
