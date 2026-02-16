package repositories

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var ThirdUserRepository = newThirdUserRepository()

func newThirdUserRepository() *thirdUserRepository {
	return &thirdUserRepository{}
}

type thirdUserRepository struct {
}

func (r *thirdUserRepository) Get(db *gorm.DB, id int64) *models.ThirdUser {
	ret := &models.ThirdUser{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *thirdUserRepository) Take(db *gorm.DB, where ...interface{}) *models.ThirdUser {
	ret := &models.ThirdUser{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *thirdUserRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.ThirdUser) {
	cnd.Find(db, &list)
	return
}

func (r *thirdUserRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.ThirdUser {
	ret := &models.ThirdUser{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *thirdUserRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.ThirdUser, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *thirdUserRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.ThirdUser, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.ThirdUser{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *thirdUserRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (list []models.ThirdUser) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *thirdUserRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *thirdUserRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.ThirdUser{})
}

func (r *thirdUserRepository) Create(db *gorm.DB, t *models.ThirdUser) (err error) {
	err = db.Create(t).Error
	return
}

func (r *thirdUserRepository) Update(db *gorm.DB, t *models.ThirdUser) (err error) {
	err = db.Save(t).Error
	return
}

func (r *thirdUserRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.ThirdUser{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *thirdUserRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.ThirdUser{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *thirdUserRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.ThirdUser{}, "id = ?", id)
}

func (r *thirdUserRepository) GetByOpenId(db *gorm.DB, openId string, thirdType constants.ThirdType) *models.ThirdUser {
	return r.FindOne(db, sqls.NewCnd().Where("open_id = ? and third_type = ?", openId, thirdType))
}

func (r *thirdUserRepository) GetByUserId(db *gorm.DB, userId int64, thirdType constants.ThirdType) *models.ThirdUser {
	return r.FindOne(db, sqls.NewCnd().Where("user_id = ? and third_type = ?", userId, thirdType))
}
