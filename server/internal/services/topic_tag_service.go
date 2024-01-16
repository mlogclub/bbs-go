package services

import (
	"bbs-go/internal/models/constants"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/models"
	"bbs-go/internal/repositories"
)

var TopicTagService = newTopicTagService()

func newTopicTagService() *topicTagService {
	return &topicTagService{}
}

type topicTagService struct {
}

func (s *topicTagService) Get(id int64) *models.TopicTag {
	return repositories.TopicTagRepository.Get(sqls.DB(), id)
}

func (s *topicTagService) Take(where ...interface{}) *models.TopicTag {
	return repositories.TopicTagRepository.Take(sqls.DB(), where...)
}

func (s *topicTagService) Find(cnd *sqls.Cnd) []models.TopicTag {
	return repositories.TopicTagRepository.Find(sqls.DB(), cnd)
}

func (s *topicTagService) FindOne(cnd *sqls.Cnd) *models.TopicTag {
	return repositories.TopicTagRepository.FindOne(sqls.DB(), cnd)
}

func (s *topicTagService) FindPageByParams(params *params.QueryParams) (list []models.TopicTag, paging *sqls.Paging) {
	return repositories.TopicTagRepository.FindPageByParams(sqls.DB(), params)
}

func (s *topicTagService) FindPageByCnd(cnd *sqls.Cnd) (list []models.TopicTag, paging *sqls.Paging) {
	return repositories.TopicTagRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *topicTagService) Create(t *models.TopicTag) error {
	return repositories.TopicTagRepository.Create(sqls.DB(), t)
}

func (s *topicTagService) Update(t *models.TopicTag) error {
	return repositories.TopicTagRepository.Update(sqls.DB(), t)
}

func (s *topicTagService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.TopicTagRepository.Updates(sqls.DB(), id, columns)
}

func (s *topicTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.TopicTagRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *topicTagService) DeleteByTopicId(topicId int64) {
	sqls.DB().Model(models.TopicTag{}).Where("topic_id = ?", topicId).UpdateColumn("status", constants.StatusDeleted)
}

func (s *topicTagService) UndeleteByTopicId(topicId int64) {
	sqls.DB().Model(models.TopicTag{}).Where("topic_id = ?", topicId).UpdateColumn("status", constants.StatusOk)
}
