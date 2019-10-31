package services

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var TopicLikeService = newTopicLikeService()

func newTopicLikeService() *topicLikeService {
	return &topicLikeService{}
}

type topicLikeService struct {
}

func (this *topicLikeService) Get(id int64) *model.TopicLike {
	return repositories.TopicLikeRepository.Get(simple.GetDB(), id)
}

func (this *topicLikeService) Take(where ...interface{}) *model.TopicLike {
	return repositories.TopicLikeRepository.Take(simple.GetDB(), where...)
}

func (this *topicLikeService) QueryCnd(cnd *simple.QueryCnd) (list []model.TopicLike, err error) {
	return repositories.TopicLikeRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *topicLikeService) Query(params *simple.ParamQueries) (list []model.TopicLike, paging *simple.Paging) {
	return repositories.TopicLikeRepository.Query(simple.GetDB(), queries)
}

func (this *topicLikeService) Create(t *model.TopicLike) error {
	return repositories.TopicLikeRepository.Create(simple.GetDB(), t)
}

func (this *topicLikeService) Update(t *model.TopicLike) error {
	return repositories.TopicLikeRepository.Update(simple.GetDB(), t)
}

func (this *topicLikeService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.TopicLikeRepository.Updates(simple.GetDB(), id, columns)
}

func (this *topicLikeService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.TopicLikeRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *topicLikeService) Delete(id int64) {
	repositories.TopicLikeRepository.Delete(simple.GetDB(), id)
}

// 统计数量
func (this *topicLikeService) Count(topicId int64) int64 {
	var count int64 = 0
	simple.GetDB().Model(&model.TopicLike{}).Where("topic_id = ?", topicId).Count(&count)
	return count
}

func (this *topicLikeService) Like(userId int64, topicId int64) error {
	topic := repositories.TopicRepository.Get(simple.GetDB(), topicId)
	if topic == nil || topic.Status != model.TopicStatusOk {
		return errors.New("话题不存在")
	}

	// 判断是否已经点赞了
	topicLike := repositories.TopicLikeRepository.Take(simple.GetDB(), "user_id = ? and topic_id = ?", userId, topicId)
	if topicLike != nil {
		return errors.New("已点赞")
	}

	return simple.Tx(simple.GetDB(), func(tx *gorm.DB) error {
		// 点赞
		err := repositories.TopicLikeRepository.Create(tx, &model.TopicLike{
			UserId:     userId,
			TopicId:    topicId,
			CreateTime: simple.NowTimestamp(),
		})
		if err != nil {
			return err
		}

		// 更新帖子点赞数
		return simple.GetDB().Exec("update t_topic set like_count = like_count + 1 where id = ?", topicId).Error
	})
}
