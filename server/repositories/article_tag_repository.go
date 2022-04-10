package repositories

import (
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"bbs-go/model"
)

var ArticleTagRepository = newArticleTagRepository()

func newArticleTagRepository() *articleTagRepository {
	return &articleTagRepository{}
}

type articleTagRepository struct {
}

func (r *articleTagRepository) Get(db *gorm.DB, id int64) *model.ArticleTag {
	ret := &model.ArticleTag{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *articleTagRepository) Take(db *gorm.DB, where ...interface{}) *model.ArticleTag {
	ret := &model.ArticleTag{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *articleTagRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.ArticleTag) {
	cnd.Find(db, &list)
	return
}

func (r *articleTagRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.ArticleTag {
	ret := &model.ArticleTag{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *articleTagRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []model.ArticleTag, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *articleTagRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []model.ArticleTag, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.ArticleTag{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *articleTagRepository) Create(db *gorm.DB, t *model.ArticleTag) (err error) {
	err = db.Create(t).Error
	return
}

func (r *articleTagRepository) Update(db *gorm.DB, t *model.ArticleTag) (err error) {
	err = db.Save(t).Error
	return
}

func (r *articleTagRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.ArticleTag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *articleTagRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.ArticleTag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *articleTagRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.ArticleTag{}, "id = ?", id)
}

func (r *articleTagRepository) AddArticleTags(db *gorm.DB, articleId int64, tagIds []int64) {
	if articleId <= 0 || len(tagIds) == 0 {
		return
	}

	for _, tagId := range tagIds {
		_ = r.Create(db, &model.ArticleTag{
			ArticleId:  articleId,
			TagId:      tagId,
			CreateTime: dates.NowTimestamp(),
		})
	}
}

func (r *articleTagRepository) DeleteArticleTags(db *gorm.DB, articleId int64) {
	if articleId <= 0 {
		return
	}
	db.Where("article_id = ?", articleId).Delete(model.ArticleTag{})
}

func (r *articleTagRepository) DeleteArticleTag(db *gorm.DB, articleId, tagId int64) {
	if articleId <= 0 {
		return
	}
	db.Where("article_id = ? and tag_id = ?", articleId, tagId).Delete(model.ArticleTag{})
}

func (r *articleTagRepository) FindByArticleId(db *gorm.DB, articleId int64) []model.ArticleTag {
	return r.Find(db, sqls.NewCnd().Where("article_id = ?", articleId))
}
