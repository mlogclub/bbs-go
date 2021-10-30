package services

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/repositories"
	"github.com/mlogclub/simple"
)

var UserFeedService = newUserFeedService()

func newUserFeedService() *userFeedService {
	return &userFeedService{}
}

type userFeedService struct {
}

func (s *userFeedService) Get(id int64) *model.UserFeed {
	return repositories.UserFeedRepository.Get(simple.DB(), id)
}

func (s *userFeedService) Take(where ...interface{}) *model.UserFeed {
	return repositories.UserFeedRepository.Take(simple.DB(), where...)
}

func (s *userFeedService) Find(cnd *simple.SqlCnd) []model.UserFeed {
	return repositories.UserFeedRepository.Find(simple.DB(), cnd)
}

func (s *userFeedService) FindOne(cnd *simple.SqlCnd) *model.UserFeed {
	return repositories.UserFeedRepository.FindOne(simple.DB(), cnd)
}

func (s *userFeedService) FindPageByParams(params *simple.QueryParams) (list []model.UserFeed, paging *simple.Paging) {
	return repositories.UserFeedRepository.FindPageByParams(simple.DB(), params)
}

func (s *userFeedService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.UserFeed, paging *simple.Paging) {
	return repositories.UserFeedRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *userFeedService) Count(cnd *simple.SqlCnd) int64 {
	return repositories.UserFeedRepository.Count(simple.DB(), cnd)
}

func (s *userFeedService) Create(t *model.UserFeed) error {
	return repositories.UserFeedRepository.Create(simple.DB(), t)
}

func (s *userFeedService) Update(t *model.UserFeed) error {
	return repositories.UserFeedRepository.Update(simple.DB(), t)
}

func (s *userFeedService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserFeedRepository.Updates(simple.DB(), id, columns)
}

func (s *userFeedService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserFeedRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *userFeedService) Delete(id int64) {
	repositories.UserFeedRepository.Delete(simple.DB(), id)
}

func (s *userFeedService) DeleteByUser(userId, authorId int64) {
	simple.DB().Where("user_id = ? and author_id = ?", userId, authorId).Delete(model.UserFeed{})
}

func (s *userFeedService) DeleteByDataId(dataId int64, dataType string) {
	simple.DB().Where("data_id = ? and data_type = ?", dataId, dataType).Delete(model.UserFeed{})
}

func (s *userFeedService) GetTopics(userId int64, cursor int64) (topics []model.Topic, nextCursor int64) {
	cnd := simple.NewSqlCnd()
	if userId > 0 {
		cnd.Eq("user_id", userId)
	}
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	cnd.Eq("data_type", constants.EntityTopic).Desc("id").Limit(20)

	userFeeds := repositories.UserFeedRepository.Find(simple.DB(), cnd)
	if len(userFeeds) > 0 {
		nextCursor = userFeeds[len(userFeeds)-1].Id
	} else {
		nextCursor = cursor
	}

	var topicIds []int64
	for _, item := range userFeeds {
		topicIds = append(topicIds, item.DataId)
	}
	topics = TopicService.GetTopicByIds(topicIds)

	return
}
