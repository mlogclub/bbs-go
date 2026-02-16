package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var UserTaskEventService = newUserTaskEventService()

func newUserTaskEventService() *userTaskEventService {
	return &userTaskEventService{}
}

type userTaskEventService struct {
}

func (s *userTaskEventService) Get(id int64) *models.UserTaskEvent {
	return repositories.UserTaskEventRepository.Get(sqls.DB(), id)
}

func (s *userTaskEventService) Take(where ...interface{}) *models.UserTaskEvent {
	return repositories.UserTaskEventRepository.Take(sqls.DB(), where...)
}

func (s *userTaskEventService) Find(cnd *sqls.Cnd) []models.UserTaskEvent {
	return repositories.UserTaskEventRepository.Find(sqls.DB(), cnd)
}

func (s *userTaskEventService) FindOne(cnd *sqls.Cnd) *models.UserTaskEvent {
	return repositories.UserTaskEventRepository.FindOne(sqls.DB(), cnd)
}

func (s *userTaskEventService) FindPageByParams(params *params.QueryParams) (list []models.UserTaskEvent, paging *sqls.Paging) {
	return repositories.UserTaskEventRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userTaskEventService) FindPageByCnd(cnd *sqls.Cnd) (list []models.UserTaskEvent, paging *sqls.Paging) {
	return repositories.UserTaskEventRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userTaskEventService) Count(cnd *sqls.Cnd) int64 {
	return repositories.UserTaskEventRepository.Count(sqls.DB(), cnd)
}

func (s *userTaskEventService) Create(t *models.UserTaskEvent) error {
	return repositories.UserTaskEventRepository.Create(sqls.DB(), t)
}

func (s *userTaskEventService) Update(t *models.UserTaskEvent) error {
	return repositories.UserTaskEventRepository.Update(sqls.DB(), t)
}

func (s *userTaskEventService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserTaskEventRepository.Updates(sqls.DB(), id, columns)
}

func (s *userTaskEventService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserTaskEventRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *userTaskEventService) Delete(id int64) {
	repositories.UserTaskEventRepository.Delete(sqls.DB(), id)
}

