package services

import (
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var ArticleTagService = newArticleTagService()

func newArticleTagService() *articleTagService {
	return &articleTagService{}
}

type articleTagService struct {
}

func (this *articleTagService) Get(id int64) *model.ArticleTag {
	return repositories.ArticleTagRepository.Get(simple.GetDB(), id)
}

func (this *articleTagService) Take(where ...interface{}) *model.ArticleTag {
	return repositories.ArticleTagRepository.Take(simple.GetDB(), where...)
}

func (this *articleTagService) QueryCnd(cnd *simple.QueryCnd) (list []model.ArticleTag, err error) {
	return repositories.ArticleTagRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *articleTagService) Query(params *simple.ParamQueries) (list []model.ArticleTag, paging *simple.Paging) {
	return repositories.ArticleTagRepository.Query(simple.GetDB(), queries)
}

func (this *articleTagService) Create(t *model.ArticleTag) error {
	return repositories.ArticleTagRepository.Create(simple.GetDB(), t)
}

func (this *articleTagService) Update(t *model.ArticleTag) error {
	return repositories.ArticleTagRepository.Update(simple.GetDB(), t)
}

func (this *articleTagService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ArticleTagRepository.Updates(simple.GetDB(), id, columns)
}

func (this *articleTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ArticleTagRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *articleTagService) DeleteByArticleId(topicId int64) {
	simple.GetDB().Model(model.ArticleTag{}).Where("topic_id = ?", topicId).UpdateColumn("status", model.ArticleTagStatusDeleted)
}
