package services

import (
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/simple"
)

type ArticleTagService struct {
	ArticleTagRepository *repositories.ArticleTagRepository
}

func NewArticleTagService() *ArticleTagService {
	return &ArticleTagService{
		ArticleTagRepository: repositories.NewArticleTagRepository(),
	}
}

func (this *ArticleTagService) Get(id int64) *model.ArticleTag {
	return this.ArticleTagRepository.Get(simple.GetDB(), id)
}

func (this *ArticleTagService) Take(where ...interface{}) *model.ArticleTag {
	return this.ArticleTagRepository.Take(simple.GetDB(), where...)
}

func (this *ArticleTagService) QueryCnd(cnd *simple.QueryCnd) (list []model.ArticleTag, err error) {
	return this.ArticleTagRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *ArticleTagService) Query(queries *simple.ParamQueries) (list []model.ArticleTag, paging *simple.Paging) {
	return this.ArticleTagRepository.Query(simple.GetDB(), queries)
}

func (this *ArticleTagService) Create(t *model.ArticleTag) error {
	return this.ArticleTagRepository.Create(simple.GetDB(), t)
}

func (this *ArticleTagService) Update(t *model.ArticleTag) error {
	return this.ArticleTagRepository.Update(simple.GetDB(), t)
}

func (this *ArticleTagService) Updates(id int64, columns map[string]interface{}) error {
	return this.ArticleTagRepository.Updates(simple.GetDB(), id, columns)
}

func (this *ArticleTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.ArticleTagRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *ArticleTagService) Delete(id int64) {
	this.ArticleTagRepository.Delete(simple.GetDB(), id)
}
