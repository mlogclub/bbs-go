package services

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var UserFeedService = newUserFeedService()

func newUserFeedService() *userFeedService {
	return &userFeedService{}
}

type userFeedService struct {
}

func (s *userFeedService) Get(id int64) *model.UserFeed {
	return repositories.UserFeedRepository.Get(sqls.DB(), id)
}

func (s *userFeedService) Take(where ...interface{}) *model.UserFeed {
	return repositories.UserFeedRepository.Take(sqls.DB(), where...)
}

func (s *userFeedService) Find(cnd *sqls.Cnd) []model.UserFeed {
	return repositories.UserFeedRepository.Find(sqls.DB(), cnd)
}

func (s *userFeedService) FindOne(cnd *sqls.Cnd) *model.UserFeed {
	return repositories.UserFeedRepository.FindOne(sqls.DB(), cnd)
}

func (s *userFeedService) FindPageByParams(params *params.QueryParams) (list []model.UserFeed, paging *sqls.Paging) {
	return repositories.UserFeedRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userFeedService) FindPageByCnd(cnd *sqls.Cnd) (list []model.UserFeed, paging *sqls.Paging) {
	return repositories.UserFeedRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userFeedService) Count(cnd *sqls.Cnd) int64 {
	return repositories.UserFeedRepository.Count(sqls.DB(), cnd)
}

func (s *userFeedService) Create(t *model.UserFeed) error {
	return repositories.UserFeedRepository.Create(sqls.DB(), t)
}

func (s *userFeedService) Update(t *model.UserFeed) error {
	return repositories.UserFeedRepository.Update(sqls.DB(), t)
}

func (s *userFeedService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserFeedRepository.Updates(sqls.DB(), id, columns)
}

func (s *userFeedService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserFeedRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *userFeedService) Delete(id int64) {
	repositories.UserFeedRepository.Delete(sqls.DB(), id)
}

func (s *userFeedService) DeleteByUser(userId, authorId int64) {
	sqls.DB().Where("user_id = ? and author_id = ?", userId, authorId).Delete(model.UserFeed{})
}

func (s *userFeedService) DeleteByDataId(dataId int64, dataType string) {
	sqls.DB().Where("data_id = ? and data_type = ?", dataId, dataType).Delete(model.UserFeed{})
}

func (s *userFeedService) GetTopics(userId int64, cursor int64) (topics []model.Topic, nextCursor int64, hasMore bool) {
	var limit = 20
	cnd := sqls.NewCnd().Eq("user_id", userId)
	cnd.Eq("data_type", constants.EntityTopic)
	if cursor > 0 {
		cnd.Lt("create_time", cursor)
	}
	cnd.Desc("create_time").Limit(limit)

	userFeeds := repositories.UserFeedRepository.Find(sqls.DB(), cnd)
	if len(userFeeds) > 0 {
		nextCursor = userFeeds[len(userFeeds)-1].CreateTime
		hasMore = len(userFeeds) >= limit
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
