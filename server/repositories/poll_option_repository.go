package repositories

import (
	"bbs-go/model"
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"
)

var PollOptionRepository = newPollOptionRepository()

func newPollOptionRepository() *pollOptionRepository {
	return &pollOptionRepository{}
}

type pollOptionRepository struct {
}

func (r *pollOptionRepository) Get(db *gorm.DB, id int64) *model.PollOption {
	ret := &model.PollOption{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *pollOptionRepository) Take(db *gorm.DB, where ...interface{}) *model.PollOption {
	ret := &model.PollOption{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *pollOptionRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.PollOption) {
	cnd.Find(db, &list)
	return
}

func (r *pollOptionRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) *model.PollOption {
	ret := &model.PollOption{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *pollOptionRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.PollOption, paging *simple.Paging) {
	return r.FindPageByCnd(db, &params.SqlCnd)
}

func (r *pollOptionRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.PollOption, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.PollOption{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *pollOptionRepository) Count(db *gorm.DB, cnd *simple.SqlCnd) int {
	return cnd.Count(db, &model.PollOption{})
}

func (r *pollOptionRepository) Create(db *gorm.DB, t *model.PollOption) (err error) {
	err = db.Create(t).Error
	return
}

func (r *pollOptionRepository) Update(db *gorm.DB, t *model.PollOption) (err error) {
	err = db.Save(t).Error
	return
}

func (r *pollOptionRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.PollOption{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *pollOptionRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.PollOption{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *pollOptionRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.PollOption{}, "id = ?", id)
}
