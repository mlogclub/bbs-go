package repositories

import (
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"bbs-go/internal/models"
)

var UserRepository = newUserRepository()

func newUserRepository() *userRepository {
	return &userRepository{}
}

type userRepository struct {
}

func (r *userRepository) Get(db *gorm.DB, id int64) *models.User {
	ret := &models.User{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userRepository) Take(db *gorm.DB, where ...interface{}) *models.User {
	ret := &models.User{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.User) {
	cnd.Find(db, &list)
	return
}

func (r *userRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.User {
	ret := &models.User{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *userRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.User, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *userRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.User, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.User{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *userRepository) Create(db *gorm.DB, t *models.User) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userRepository) Update(db *gorm.DB, t *models.User) (err error) {
	err = db.Save(t).Error
	return
}

func (r *userRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.User{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *userRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.User{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *userRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.User{}, "id = ?", id)
}

func (r *userRepository) GetByEmail(db *gorm.DB, email string) *models.User {
	return r.Take(db, "email = ?", email)
}

func (r *userRepository) GetByUsername(db *gorm.DB, username string) *models.User {
	return r.Take(db, "username = ?", username)
}
