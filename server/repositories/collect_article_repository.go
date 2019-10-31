
package repositories

import (
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/simple"
	"github.com/jinzhu/gorm"
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

func (this *collectArticleRepository) QueryCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.CollectArticle, err error) {
	err = cnd.Exec(db).Find(&list).Error
	return
}

func (this *collectArticleRepository) Query(db *gorm.DB, params *simple.QueryParams) (list []model.CollectArticle, paging *simple.Paging) {
	params.StartQuery(db).Find(&list)
    params.StartCount(db).Model(&model.CollectArticle{}).Count(&params.Paging.Total)
	paging = params.Paging
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

