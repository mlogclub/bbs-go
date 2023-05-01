package services

import (
	"bbs-go/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/model"
)

var AdvertService = newAdvertService()

func newAdvertService() *advertService {
	return &advertService{}
}

type advertService struct {
}

func (s *advertService) Get(id int64) *model.Advert {
	return repositories.AdvertRepository.Get(sqls.DB(), id)
}

func (s *advertService) Take(where ...interface{}) *model.Advert {
	return repositories.AdvertRepository.Take(sqls.DB(), where...)
}

func (s *advertService) Find(cnd *sqls.Cnd) []model.Advert {
	return repositories.AdvertRepository.Find(sqls.DB(), cnd)
}

func (s *advertService) FindOne(cnd *sqls.Cnd) *model.Advert {
	return repositories.AdvertRepository.FindOne(sqls.DB(), cnd)
}

func (s *advertService) FindPageByParams(params *params.QueryParams) (list []model.Advert, paging *sqls.Paging) {
	return repositories.AdvertRepository.FindPageByParams(sqls.DB(), params)
}

func (s *advertService) FindPageByCnd(cnd *sqls.Cnd) (list []model.Advert, paging *sqls.Paging) {
	return repositories.AdvertRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *advertService) Create(t *model.Advert) error {
	return repositories.AdvertRepository.Create(sqls.DB(), t)
}

func (s *advertService) Update(t *model.Advert) error {
	return repositories.AdvertRepository.Update(sqls.DB(), t)
}

func (s *advertService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.AdvertRepository.Updates(sqls.DB(), id, columns)
}

func (s *advertService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.AdvertRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *advertService) Delete(id int64) {
	repositories.AdvertRepository.Delete(sqls.DB(), id)
}
