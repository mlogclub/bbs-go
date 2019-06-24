package services

import (
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
)

type ScanArticleShareCallback func(shares []model.ArticleShare)

type ArticleShareService struct {
	ArticleShareRepository *repositories.ArticleShareRepository
}

func NewArticleShareService() *ArticleShareService {
	return &ArticleShareService{
		ArticleShareRepository: repositories.NewArticleShareRepository(),
	}
}

func (this *ArticleShareService) Get(id int64) *model.ArticleShare {
	return this.ArticleShareRepository.Get(simple.GetDB(), id)
}

func (this *ArticleShareService) Take(where ...interface{}) *model.ArticleShare {
	return this.ArticleShareRepository.Take(simple.GetDB(), where...)
}

func (this *ArticleShareService) QueryCnd(cnd *simple.QueryCnd) (list []model.ArticleShare, err error) {
	return this.ArticleShareRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *ArticleShareService) Query(queries *simple.ParamQueries) (list []model.ArticleShare, paging *simple.Paging) {
	return this.ArticleShareRepository.Query(simple.GetDB(), queries)
}

func (this *ArticleShareService) Create(t *model.ArticleShare) error {
	return this.ArticleShareRepository.Create(simple.GetDB(), t)
}

func (this *ArticleShareService) Update(t *model.ArticleShare) error {
	return this.ArticleShareRepository.Update(simple.GetDB(), t)
}

func (this *ArticleShareService) Updates(id int64, columns map[string]interface{}) error {
	return this.ArticleShareRepository.Updates(simple.GetDB(), id, columns)
}

func (this *ArticleShareService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.ArticleShareRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *ArticleShareService) Delete(id int64) {
	this.ArticleShareRepository.Delete(simple.GetDB(), id)
}

// 扫描
func (this *ArticleShareService) Scan(cb ScanArticleShareCallback) {
	var cursor int64
	for {
		list, err := this.ArticleShareRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("id > ? and status = ? ",
			cursor, model.ArticleShareStatusOk).Order("id asc").Size(300))
		if err != nil {
			logrus.Error("查询失败...")
			break
		}
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		cb(list)
	}
}

// 扫描
func (this *ArticleShareService) ScanWithDate(dateFrom, dateTo int64, cb ScanArticleShareCallback) {
	var cursor int64
	for {
		list, err := this.ArticleShareRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("id > ? and status = ? and create_time >= ? and create_time < ?",
			cursor, model.ArticleShareStatusOk, dateFrom, dateTo).Order("id asc").Size(300))
		if err != nil {
			break
		}
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		cb(list)
	}
}
