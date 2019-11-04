package services

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var LinkService = newLinkService()

func newLinkService() *linkService {
	return &linkService{}
}

type linkService struct {
}

func (this *linkService) Get(id int64) *model.Link {
	return repositories.LinkRepository.Get(simple.DB(), id)
}

func (this *linkService) Take(where ...interface{}) *model.Link {
	return repositories.LinkRepository.Take(simple.DB(), where...)
}

func (this *linkService) Find(cnd *simple.SqlCnd) []model.Link {
	return repositories.LinkRepository.Find(simple.DB(), cnd)
}

func (this *linkService) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.Link) {
	cnd.FindOne(db, &ret)
	return
}

func (this *linkService) FindPageByParams(params *simple.QueryParams) (list []model.Link, paging *simple.Paging) {
	return repositories.LinkRepository.FindPageByParams(simple.DB(), params)
}

func (this *linkService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Link, paging *simple.Paging) {
	return repositories.LinkRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *linkService) Create(t *model.Link) error {
	return repositories.LinkRepository.Create(simple.DB(), t)
}

func (this *linkService) Update(t *model.Link) error {
	return repositories.LinkRepository.Update(simple.DB(), t)
}

func (this *linkService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.LinkRepository.Updates(simple.DB(), id, columns)
}

func (this *linkService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.LinkRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (this *linkService) Delete(id int64) {
	repositories.LinkRepository.Delete(simple.DB(), id)
}
