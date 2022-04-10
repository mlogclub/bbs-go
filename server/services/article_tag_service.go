package services

import (
	"bbs-go/model/constants"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/model"
	"bbs-go/repositories"
)

var ArticleTagService = newArticleTagService()

func newArticleTagService() *articleTagService {
	return &articleTagService{}
}

type articleTagService struct {
}

func (s *articleTagService) Get(id int64) *model.ArticleTag {
	return repositories.ArticleTagRepository.Get(sqls.DB(), id)
}

func (s *articleTagService) Take(where ...interface{}) *model.ArticleTag {
	return repositories.ArticleTagRepository.Take(sqls.DB(), where...)
}

func (s *articleTagService) Find(cnd *sqls.Cnd) []model.ArticleTag {
	return repositories.ArticleTagRepository.Find(sqls.DB(), cnd)
}

func (s *articleTagService) FindPageByParams(params *params.QueryParams) (list []model.ArticleTag, paging *sqls.Paging) {
	return repositories.ArticleTagRepository.FindPageByParams(sqls.DB(), params)
}

func (s *articleTagService) FindPageByCnd(cnd *sqls.Cnd) (list []model.ArticleTag, paging *sqls.Paging) {
	return repositories.ArticleTagRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *articleTagService) Create(t *model.ArticleTag) error {
	return repositories.ArticleTagRepository.Create(sqls.DB(), t)
}

func (s *articleTagService) Update(t *model.ArticleTag) error {
	return repositories.ArticleTagRepository.Update(sqls.DB(), t)
}

func (s *articleTagService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ArticleTagRepository.Updates(sqls.DB(), id, columns)
}

func (s *articleTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ArticleTagRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *articleTagService) DeleteByArticleId(topicId int64) {
	sqls.DB().Model(model.ArticleTag{}).Where("topic_id = ?", topicId).UpdateColumn("status", constants.StatusDeleted)
}
