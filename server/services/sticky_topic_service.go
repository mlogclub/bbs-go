package services

import (
	"bbs-go/model"
	"bbs-go/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var StickyTopicService = newStickyTopicService()

func newStickyTopicService() *stickyTopicService {
	return &stickyTopicService {}
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

