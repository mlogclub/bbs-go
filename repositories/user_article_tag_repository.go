package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/simple"
)

var UserArticleTagRepository = newUserArticleTagRepository()

func newUserArticleTagRepository() *userArticleTagRepository {
	return &userArticleTagRepository{}
}

type userArticleTagRepository struct {
}

func (this *userArticleTagRepository) Get(db *gorm.DB, id int64) *model.UserArticleTag {
	ret := &model.UserArticleTag{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *userArticleTagRepository) Take(db *gorm.DB, where ...interface{}) *model.UserArticleTag {
	ret := &model.UserArticleTag{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *userArticleTagRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.UserArticleTag, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *userArticleTagRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.UserArticleTag, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.UserArticleTag{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *userArticleTagRepository) Create(db *gorm.DB, t *model.UserArticleTag) (err error) {
	err = db.Create(t).Error
	return
}

func (this *userArticleTagRepository) Update(db *gorm.DB, t *model.UserArticleTag) (err error) {
	err = db.Save(t).Error
	return
}

func (this *userArticleTagRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.UserArticleTag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *userArticleTagRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.UserArticleTag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *userArticleTagRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.UserArticleTag{}).Delete("id", id)
}
