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
	return repositories.ArticleTagRepository.Get(simple.DB(), id)
}

func (this *articleTagService) Take(where ...interface{}) *model.ArticleTag {
	return repositories.ArticleTagRepository.Take(simple.DB(), where...)
}

func (this *articleTagService) Find(cnd *simple.SqlCnd) []model.ArticleTag {
	return repositories.ArticleTagRepository.Find(simple.DB(), cnd)
}

func (this *articleTagService) FindPageByParams(params *simple.QueryParams) (list []model.ArticleTag, paging *simple.Paging) {
	return repositories.ArticleTagRepository.FindPageByParams(simple.DB(), params)
}

func (this *articleTagService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.ArticleTag, paging *simple.Paging) {
	return repositories.ArticleTagRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *articleTagService) Create(t *model.ArticleTag) error {
	return repositories.ArticleTagRepository.Create(simple.DB(), t)
}

func (this *articleTagService) Update(t *model.ArticleTag) error {
	return repositories.ArticleTagRepository.Update(simple.DB(), t)
}

func (this *articleTagService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ArticleTagRepository.Updates(simple.DB(), id, columns)
}

func (this *articleTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ArticleTagRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (this *articleTagService) DeleteByArticleId(topicId int64) {
	simple.DB().Model(model.ArticleTag{}).Where("topic_id = ?", topicId).UpdateColumn("status", model.ArticleTagStatusDeleted)
}
