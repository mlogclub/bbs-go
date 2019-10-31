
package services

import (
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
	"github.com/mlogclub/simple"
)

var LinkService = newLinkService()

func newLinkService() *linkService {
	return &linkService {}
}

type linkService struct {
}

func (this *linkService) Get(id int64) *model.Link {
	return repositories.LinkRepository.Get(simple.GetDB(), id)
}

func (this *linkService) Take(where ...interface{}) *model.Link {
	return repositories.LinkRepository.Take(simple.GetDB(), where...)
}

func (this *linkService) QueryCnd(cnd *simple.QueryCnd) (list []model.Link, err error) {
	return repositories.LinkRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *linkService) Query(params *simple.ParamQueries) (list []model.Link, paging *simple.Paging) {
	return repositories.LinkRepository.Query(simple.GetDB(), queries)
}

func (this *linkService) Create(t *model.Link) error {
	return repositories.LinkRepository.Create(simple.GetDB(), t)
}

func (this *linkService) Update(t *model.Link) error {
	return repositories.LinkRepository.Update(simple.GetDB(), t)
}

func (this *linkService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.LinkRepository.Updates(simple.GetDB(), id, columns)
}

func (this *linkService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.LinkRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *linkService) Delete(id int64) {
	repositories.LinkRepository.Delete(simple.GetDB(), id)
}

