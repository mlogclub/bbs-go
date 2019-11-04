package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
)

var ArticleTagRepository = newArticleTagRepository()

func newArticleTagRepository() *articleTagRepository {
	return &articleTagRepository{}
}

type articleTagRepository struct {
}

func (this *articleTagRepository) Get(db *gorm.DB, id int64) *model.ArticleTag {
	ret := &model.ArticleTag{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *articleTagRepository) Take(db *gorm.DB, where ...interface{}) *model.ArticleTag {
	ret := &model.ArticleTag{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *articleTagRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.ArticleTag) {
	cnd.Find(db, &list)
	return
}

func (this *articleTagRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.ArticleTag) {
	cnd.FindOne(db, &ret)
	return
}

func (this *articleTagRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.ArticleTag, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *articleTagRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.ArticleTag, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.ArticleTag{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *articleTagRepository) Create(db *gorm.DB, t *model.ArticleTag) (err error) {
	err = db.Create(t).Error
	return
}

func (this *articleTagRepository) Update(db *gorm.DB, t *model.ArticleTag) (err error) {
	err = db.Save(t).Error
	return
}

func (this *articleTagRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.ArticleTag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *articleTagRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.ArticleTag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *articleTagRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.ArticleTag{}, "id = ?", id)
}

func (this *articleTagRepository) AddArticleTags(db *gorm.DB, articleId int64, tagIds []int64) {
	if articleId <= 0 || len(tagIds) == 0 {
		return
	}

	for _, tagId := range tagIds {
		_ = this.Create(db, &model.ArticleTag{
			ArticleId:  articleId,
			TagId:      tagId,
			CreateTime: simple.NowTimestamp(),
		})
	}
}

func (this *articleTagRepository) DeleteArticleTags(db *gorm.DB, articleId int64) {
	if articleId <= 0 {
		return
	}
	db.Where("article_id = ?", articleId).Delete(model.ArticleTag{})
}

func (this *articleTagRepository) FindByArticleId(db *gorm.DB, articleId int64) []model.ArticleTag {
	return this.Find(db, simple.NewSqlCnd().Where("article_id = ?", articleId))
}
