package services

import (
	"bbs-go/model"
	"bbs-go/repositories"
	"github.com/mlogclub/simple"
)

var PollAnswerService = newPollAnswerService()

func newPollAnswerService() *pollAnswerService {
	return &pollAnswerService{}
}

type pollAnswerService struct {
}

func (s *pollAnswerService) Get(id int64) *model.PollAnswer {
	return repositories.PollAnswerRepository.Get(simple.DB(), id)
}

func (s *pollAnswerService) Take(where ...interface{}) *model.PollAnswer {
	return repositories.PollAnswerRepository.Take(simple.DB(), where...)
}

func (s *pollAnswerService) Find(cnd *simple.SqlCnd) []model.PollAnswer {
	return repositories.PollAnswerRepository.Find(simple.DB(), cnd)
}

func (s *pollAnswerService) FindOne(cnd *simple.SqlCnd) *model.PollAnswer {
	return repositories.PollAnswerRepository.FindOne(simple.DB(), cnd)
}

func (s *pollAnswerService) FindPageByParams(params *simple.QueryParams) (list []model.PollAnswer, paging *simple.Paging) {
	return repositories.PollAnswerRepository.FindPageByParams(simple.DB(), params)
}

func (s *pollAnswerService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.PollAnswer, paging *simple.Paging) {
	return repositories.PollAnswerRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *pollAnswerService) Count(cnd *simple.SqlCnd) int {
	return repositories.PollAnswerRepository.Count(simple.DB(), cnd)
}

func (s *pollAnswerService) Create(t *model.PollAnswer) error {
	return repositories.PollAnswerRepository.Create(simple.DB(), t)
}

func (s *pollAnswerService) Update(t *model.PollAnswer) error {
	return repositories.PollAnswerRepository.Update(simple.DB(), t)
}

func (s *pollAnswerService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.PollAnswerRepository.Updates(simple.DB(), id, columns)
}

func (s *pollAnswerService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.PollAnswerRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *pollAnswerService) Delete(id int64) {
	repositories.PollAnswerRepository.Delete(simple.DB(), id)
}
