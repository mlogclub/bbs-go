package repositories

import (
	"bbs-go/model"
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"
)

var PollAnswerRepository = newPollAnswerRepository()

func newPollAnswerRepository() *pollAnswerRepository {
	return &pollAnswerRepository{}
}

type pollAnswerRepository struct {
}

func (r *pollAnswerRepository) Get(db *gorm.DB, id int64) *model.PollAnswer {
	ret := &model.PollAnswer{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *pollAnswerRepository) Take(db *gorm.DB, where ...interface{}) *model.PollAnswer {
	ret := &model.PollAnswer{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *pollAnswerRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.PollAnswer) {
	cnd.Find(db, &list)
	return
}

func (r *pollAnswerRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) *model.PollAnswer {
	ret := &model.PollAnswer{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *pollAnswerRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.PollAnswer, paging *simple.Paging) {
	return r.FindPageByCnd(db, &params.SqlCnd)
}

func (r *pollAnswerRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.PollAnswer, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.PollAnswer{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *pollAnswerRepository) Count(db *gorm.DB, cnd *simple.SqlCnd) int {
	return cnd.Count(db, &model.PollAnswer{})
}

func (r *pollAnswerRepository) Create(db *gorm.DB, t *model.PollAnswer) (err error) {
	err = db.Create(t).Error
	return
}

func (r *pollAnswerRepository) Update(db *gorm.DB, t *model.PollAnswer) (err error) {
	err = db.Save(t).Error
	return
}

func (r *pollAnswerRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.PollAnswer{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *pollAnswerRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.PollAnswer{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *pollAnswerRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.PollAnswer{}, "id = ?", id)
}
