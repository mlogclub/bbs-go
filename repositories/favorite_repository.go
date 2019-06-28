package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/simple"
)

type FavoriteRepository struct {
}

func NewFavoriteRepository() *FavoriteRepository {
	return &FavoriteRepository{}
}

func (this *FavoriteRepository) Get(db *gorm.DB, id int64) *model.Favorite {
	ret := &model.Favorite{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *FavoriteRepository) Take(db *gorm.DB, where ...interface{}) *model.Favorite {
	ret := &model.Favorite{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *FavoriteRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.Favorite, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *FavoriteRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.Favorite, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.Favorite{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *FavoriteRepository) Create(db *gorm.DB, t *model.Favorite) (err error) {
	err = db.Create(t).Error
	return
}

func (this *FavoriteRepository) Update(db *gorm.DB, t *model.Favorite) (err error) {
	err = db.Save(t).Error
	return
}

func (this *FavoriteRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Favorite{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *FavoriteRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Favorite{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *FavoriteRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.Favorite{}).Delete("id", id)
}
