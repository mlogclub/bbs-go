package services

import (
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/simple"
)

var ArticleTagService = newArticleTagService()

func newArticleTagService() *articleTagService {
	return &articleTagService{
		ArticleTagRepository: repositories.NewArticleTagRepository(),
	}
}

type articleTagService struct {
	ArticleTagRepository *repositories.ArticleTagRepository
}

func (this *articleTagService) Get(id int64) *model.ArticleTag {
	return this.ArticleTagRepository.Get(simple.GetDB(), id)
}

func (this *articleTagService) Take(where ...interface{}) *model.ArticleTag {
	return this.ArticleTagRepository.Take(simple.GetDB(), where...)
}

func (this *articleTagService) QueryCnd(cnd *simple.QueryCnd) (list []model.ArticleTag, err error) {
	return this.ArticleTagRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *articleTagService) Query(queries *simple.ParamQueries) (list []model.ArticleTag, paging *simple.Paging) {
	return this.ArticleTagRepository.Query(simple.GetDB(), queries)
}

func (this *articleTagService) Create(t *model.ArticleTag) error {
	return this.ArticleTagRepository.Create(simple.GetDB(), t)
}

func (this *articleTagService) Update(t *model.ArticleTag) error {
	return this.ArticleTagRepository.Update(simple.GetDB(), t)
}

func (this *articleTagService) Updates(id int64, columns map[string]interface{}) error {
	return this.ArticleTagRepository.Updates(simple.GetDB(), id, columns)
}

func (this *articleTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.ArticleTagRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *articleTagService) Delete(id int64) {
	this.ArticleTagRepository.Delete(simple.GetDB(), id)
}
