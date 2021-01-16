package services

import (
	"bbs-go/model/constants"
	"errors"
	"github.com/mlogclub/simple/date"

	"github.com/mlogclub/simple"
	"gorm.io/gorm"

	"bbs-go/model"
	"bbs-go/repositories"
)

var UserLikeService = newUserLikeService()

func newUserLikeService() *userLikeService {
	return &userLikeService{}
}

type userLikeService struct {
}

func (s *userLikeService) Get(id int64) *model.UserLike {
	return repositories.UserLikeRepository.Get(simple.DB(), id)
}

func (s *userLikeService) Take(where ...interface{}) *model.UserLike {
	return repositories.UserLikeRepository.Take(simple.DB(), where...)
}

func (s *userLikeService) Find(cnd *simple.SqlCnd) []model.UserLike {
	return repositories.UserLikeRepository.Find(simple.DB(), cnd)
}

func (s *userLikeService) FindOne(cnd *simple.SqlCnd) *model.UserLike {
	return repositories.UserLikeRepository.FindOne(simple.DB(), cnd)
}

func (s *userLikeService) FindPageByParams(params *simple.QueryParams) (list []model.UserLike, paging *simple.Paging) {
	return repositories.UserLikeRepository.FindPageByParams(simple.DB(), params)
}

func (s *userLikeService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.UserLike, paging *simple.Paging) {
	return repositories.UserLikeRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *userLikeService) Create(t *model.UserLike) error {
	return repositories.UserLikeRepository.Create(simple.DB(), t)
}

func (s *userLikeService) Update(t *model.UserLike) error {
	return repositories.UserLikeRepository.Update(simple.DB(), t)
}

func (s *userLikeService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserLikeRepository.Updates(simple.DB(), id, columns)
}

func (s *userLikeService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserLikeRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *userLikeService) Delete(id int64) {
	repositories.UserLikeRepository.Delete(simple.DB(), id)
}

// 统计数量
func (s *userLikeService) Count(entityType string, entityId int64) int64 {
	var count int64 = 0
	simple.DB().Model(&model.UserLike{}).Where("entity_type = ?", entityType).Where("entity_id = ?", entityId).Count(&count)
	return count
}

// 最近点赞
func (s *userLikeService) Recent(entityType string, entityId int64, count int) []model.UserLike {
	return s.Find(simple.NewSqlCnd().Eq("entity_type", entityType).Eq("entity_id", entityId).Desc("id").Limit(count))
}

// Exists 是否点赞
func (s *userLikeService) Exists(userId int64, entityType string, entityId int64) bool {
	return repositories.UserLikeRepository.FindOne(simple.DB(), simple.NewSqlCnd().Eq("user_id", userId).
		Eq("entity_type", entityType).Eq("entity_id", entityId)) != nil
}

// 是否点赞，返回已点赞实体编号
func (s *userLikeService) IsLiked(userId int64, entityType string, entityIds []int64) (likedEntityIds []int64) {
	list := repositories.UserLikeRepository.Find(simple.DB(), simple.NewSqlCnd().Eq("user_id", userId).
		Eq("entity_type", entityType).In("entity_id", entityIds))
	for _, like := range list {
		likedEntityIds = append(likedEntityIds, like.EntityId)
	}
	return
}

// 话题点赞
func (s *userLikeService) TopicLike(userId int64, topicId int64) error {
	topic := repositories.TopicRepository.Get(simple.DB(), topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return errors.New("话题不存在")
	}

	if err := simple.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.like(tx, userId, constants.EntityTopic, topicId); err != nil {
			return err
		}
		// 更新点赞数
		return tx.Exec("update t_topic set like_count = like_count + 1 where id = ?", topicId).Error
	}); err != nil {
		return err
	}

	// 发送消息
	MessageService.SendTopicLikeMsg(topicId, userId)

	return nil
}

func (s *userLikeService) like(db *gorm.DB, userId int64, entityType string, entityId int64) error {
	// 判断是否已经点赞了
	if s.Exists(userId, entityType, entityId) {
		return errors.New("已点赞")
	}
	// 点赞
	return repositories.UserLikeRepository.Create(db, &model.UserLike{
		UserId:     userId,
		EntityType: entityType,
		EntityId:   entityId,
		CreateTime: date.NowTimestamp(),
	})
}
