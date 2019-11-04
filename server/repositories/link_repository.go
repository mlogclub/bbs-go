package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
)

var LinkRepository = newLinkRepository()

func newLinkRepository() *linkRepository {
	return &linkRepository{}
}

type linkRepository struct {
}

func (this *linkRepository) Get(db *gorm.DB, id int64) *model.Link {
	ret := &model.Link{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *linkRepository) Take(db *gorm.DB, where ...interface{}) *model.Link {
	ret := &model.Link{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *linkRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Link) {
	cnd.Find(db, &list)
	return
}

func (this *linkRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.Link) {
	cnd.FindOne(db, &ret)
	return
}

func (this *linkRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.Link, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *linkRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Link, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Link{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *linkRepository) Create(db *gorm.DB, t *model.Link) (err error) {
	err = db.Create(t).Error
	return
}

func (this *linkRepository) Update(db *gorm.DB, t *model.Link) (err error) {
	err = db.Save(t).Error
	return
}

func (this *linkRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Link{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *linkRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Link{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *linkRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Link{}, "id = ?", id)
}
