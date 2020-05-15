package services

import (
	"bbs-go/model"
	"bbs-go/repositories"
	"github.com/mlogclub/simple"
)

var PollOptionService = newPollOptionService()

func newPollOptionService() *pollOptionService {
	return &pollOptionService{}
}

type pollOptionService struct {
}

func (s *pollOptionService) Get(id int64) *model.PollOption {
	return repositories.PollOptionRepository.Get(simple.DB(), id)
}

func (s *pollOptionService) Take(where ...interface{}) *model.PollOption {
	return repositories.PollOptionRepository.Take(simple.DB(), where...)
}

func (s *pollOptionService) Find(cnd *simple.SqlCnd) []model.PollOption {
	return repositories.PollOptionRepository.Find(simple.DB(), cnd)
}

func (s *pollOptionService) FindOne(cnd *simple.SqlCnd) *model.PollOption {
	return repositories.PollOptionRepository.FindOne(simple.DB(), cnd)
}

func (s *pollOptionService) FindPageByParams(params *simple.QueryParams) (list []model.PollOption, paging *simple.Paging) {
	return repositories.PollOptionRepository.FindPageByParams(simple.DB(), params)
}

func (s *pollOptionService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.PollOption, paging *simple.Paging) {
	return repositories.PollOptionRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *pollOptionService) Count(cnd *simple.SqlCnd) int {
	return repositories.PollOptionRepository.Count(simple.DB(), cnd)
}

func (s *pollOptionService) Create(t *model.PollOption) error {
	return repositories.PollOptionRepository.Create(simple.DB(), t)
}

func (s *pollOptionService) Update(t *model.PollOption) error {
	return repositories.PollOptionRepository.Update(simple.DB(), t)
}

func (s *pollOptionService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.PollOptionRepository.Updates(simple.DB(), id, columns)
}

func (s *pollOptionService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.PollOptionRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *pollOptionService) Delete(id int64) {
	repositories.PollOptionRepository.Delete(simple.DB(), id)
}
