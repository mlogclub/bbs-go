package repositories

import (
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"

	"server/model"
)

var UserRefereeRepository = newUserRefereeRepository()

func newUserRefereeRepository() *userRefereeRepository {
	return &userRefereeRepository{}
}

type userRefereeRepository struct {
}

func (r *userRefereeRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.UserReferee {
	ret := &model.UserReferee{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *userRefereeRepository) Create(db *gorm.DB, t *model.UserReferee) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userRefereeRepository) Update(db *gorm.DB, t *model.UserReferee) (err error) {
	err = db.Save(t).Error
	return
}
