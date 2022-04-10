package repositories

import (
	"bbs-go/model"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var EmailCodeRepository = newEmailCodeRepository()

func newEmailCodeRepository() *emailCodeRepository {
	return &emailCodeRepository{}
}

type emailCodeRepository struct {
}

func (r *emailCodeRepository) Get(db *gorm.DB, id int64) *model.EmailCode {
	ret := &model.EmailCode{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *emailCodeRepository) Take(db *gorm.DB, where ...interface{}) *model.EmailCode {
	ret := &model.EmailCode{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *emailCodeRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.EmailCode) {
	cnd.Find(db, &list)
	return
}

func (r *emailCodeRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.EmailCode {
	ret := &model.EmailCode{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *emailCodeRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []model.EmailCode, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *emailCodeRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []model.EmailCode, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.EmailCode{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *emailCodeRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &model.EmailCode{})
}

func (r *emailCodeRepository) Create(db *gorm.DB, t *model.EmailCode) (err error) {
	err = db.Create(t).Error
	return
}

func (r *emailCodeRepository) Update(db *gorm.DB, t *model.EmailCode) (err error) {
	err = db.Save(t).Error
	return
}

func (r *emailCodeRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.EmailCode{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *emailCodeRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.EmailCode{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *emailCodeRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.EmailCode{}, "id = ?", id)
}
