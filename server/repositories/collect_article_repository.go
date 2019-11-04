package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
)

var CollectArticleRepository = newCollectArticleRepository()

func newCollectArticleRepository() *collectArticleRepository {
	return &collectArticleRepository{}
}

type collectArticleRepository struct {
}

func (this *collectArticleRepository) Get(db *gorm.DB, id int64) *model.CollectArticle {
	ret := &model.CollectArticle{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *collectArticleRepository) Take(db *gorm.DB, where ...interface{}) *model.CollectArticle {
	ret := &model.CollectArticle{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *collectArticleRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.CollectArticle) {
	cnd.Find(db, &list)
	return
}

func (this *collectArticleRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.CollectArticle) {
	cnd.FindOne(db, &ret)
	return
}

func (this *collectArticleRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.CollectArticle, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *collectArticleRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.CollectArticle, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.CollectArticle{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *collectArticleRepository) Create(db *gorm.DB, t *model.CollectArticle) (err error) {
	err = db.Create(t).Error
	return
}

func (this *collectArticleRepository) Update(db *gorm.DB, t *model.CollectArticle) (err error) {
	err = db.Save(t).Error
	return
}

func (this *collectArticleRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.CollectArticle{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *collectArticleRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.CollectArticle{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *collectArticleRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.CollectArticle{}, "id = ?", id)
}
