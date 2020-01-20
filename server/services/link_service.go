package services

import (
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/repositories"
)

var LinkService = newLinkService()

func newLinkService() *linkService {
	return &linkService{}
}

type linkService struct {
}

func (s *linkService) Get(id int64) *model.Link {
	return repositories.LinkRepository.Get(simple.DB(), id)
}

func (s *linkService) Take(where ...interface{}) *model.Link {
	return repositories.LinkRepository.Take(simple.DB(), where...)
}

func (s *linkService) Find(cnd *simple.SqlCnd) []model.Link {
	return repositories.LinkRepository.Find(simple.DB(), cnd)
}

func (s *linkService) FindOne(cnd *simple.SqlCnd) *model.Link {
	return repositories.LinkRepository.FindOne(simple.DB(), cnd)
}

func (s *linkService) FindPageByParams(params *simple.QueryParams) (list []model.Link, paging *simple.Paging) {
	return repositories.LinkRepository.FindPageByParams(simple.DB(), params)
}

func (s *linkService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Link, paging *simple.Paging) {
	return repositories.LinkRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *linkService) Create(t *model.Link) error {
	return repositories.LinkRepository.Create(simple.DB(), t)
}

func (s *linkService) Update(t *model.Link) error {
	return repositories.LinkRepository.Update(simple.DB(), t)
}

func (s *linkService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.LinkRepository.Updates(simple.DB(), id, columns)
}

func (s *linkService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.LinkRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *linkService) Delete(id int64) {
	repositories.LinkRepository.Delete(simple.DB(), id)
}
