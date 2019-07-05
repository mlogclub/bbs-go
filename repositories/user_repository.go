package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
)

var UserRepository = newUserRepository()

func newUserRepository() *userRepository {
	return &userRepository{}
}

type userRepository struct {
}

func (this *userRepository) Get(db *gorm.DB, id int64) *model.User {
	ret := &model.User{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *userRepository) Take(db *gorm.DB, where ...interface{}) *model.User {
	ret := &model.User{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *userRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.User, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *userRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.User, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.User{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *userRepository) Create(db *gorm.DB, t *model.User) (err error) {
	err = db.Create(t).Error
	return
}

func (this *userRepository) Update(db *gorm.DB, t *model.User) (err error) {
	err = db.Save(t).Error
	return
}

func (this *userRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.User{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *userRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.User{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *userRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.User{}).Delete("id", id)
}

func (this *userRepository) GetByEmail(db *gorm.DB, email string) *model.User {
	return this.Take(db, "email = ?", email)
}

func (this *userRepository) GetByUsername(db *gorm.DB, username string) *model.User {
	return this.Take(db, "username = ?", username)
}
