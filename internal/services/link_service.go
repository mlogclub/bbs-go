package services

import (
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/models"
	"bbs-go/internal/repositories"
)

var LinkService = newLinkService()

func newLinkService() *linkService {
	return &linkService{}
}

type linkService struct {
}

func (s *linkService) Get(id int64) *models.Link {
	return repositories.LinkRepository.Get(sqls.DB(), id)
}

func (s *linkService) Take(where ...interface{}) *models.Link {
	return repositories.LinkRepository.Take(sqls.DB(), where...)
}

func (s *linkService) Find(cnd *sqls.Cnd) []models.Link {
	return repositories.LinkRepository.Find(sqls.DB(), cnd)
}

func (s *linkService) FindOne(cnd *sqls.Cnd) *models.Link {
	return repositories.LinkRepository.FindOne(sqls.DB(), cnd)
}

func (s *linkService) FindPageByParams(params *params.QueryParams) (list []models.Link, paging *sqls.Paging) {
	return repositories.LinkRepository.FindPageByParams(sqls.DB(), params)
}

func (s *linkService) FindPageByCnd(cnd *sqls.Cnd) (list []models.Link, paging *sqls.Paging) {
	return repositories.LinkRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *linkService) Create(t *models.Link) error {
	return repositories.LinkRepository.Create(sqls.DB(), t)
}

func (s *linkService) Update(t *models.Link) error {
	return repositories.LinkRepository.Update(sqls.DB(), t)
}

func (s *linkService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.LinkRepository.Updates(sqls.DB(), id, columns)
}

func (s *linkService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.LinkRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *linkService) Delete(id int64) {
	repositories.LinkRepository.Delete(sqls.DB(), id)
}
