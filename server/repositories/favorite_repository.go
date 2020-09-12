package repositories

import (
	"github.com/mlogclub/simple"
	"gorm.io/gorm"

	"bbs-go/model"
)

var FavoriteRepository = newFavoriteRepository()

func newFavoriteRepository() *favoriteRepository {
	return &favoriteRepository{}
}

type favoriteRepository struct {
}

func (r *favoriteRepository) Get(db *gorm.DB, id int64) *model.Favorite {
	ret := &model.Favorite{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *favoriteRepository) Take(db *gorm.DB, where ...interface{}) *model.Favorite {
	ret := &model.Favorite{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *favoriteRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Favorite) {
	cnd.Find(db, &list)
	return
}

func (r *favoriteRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) *model.Favorite {
	ret := &model.Favorite{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *favoriteRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.Favorite, paging *simple.Paging) {
	return r.FindPageByCnd(db, &params.SqlCnd)
}

func (r *favoriteRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Favorite, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Favorite{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *favoriteRepository) Create(db *gorm.DB, t *model.Favorite) (err error) {
	err = db.Create(t).Error
	return
}

func (r *favoriteRepository) Update(db *gorm.DB, t *model.Favorite) (err error) {
	err = db.Save(t).Error
	return
}

func (r *favoriteRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Favorite{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *favoriteRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Favorite{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *favoriteRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Favorite{}, "id = ?", id)
}
