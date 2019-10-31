
package repositories

import (
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/simple"
	"github.com/jinzhu/gorm"
)

var ProjectRepository = newProjectRepository()

func newProjectRepository() *projectRepository {
	return &projectRepository{}
}

type projectRepository struct {
}

func (this *projectRepository) Get(db *gorm.DB, id int64) *model.Project {
	ret := &model.Project{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *projectRepository) Take(db *gorm.DB, where ...interface{}) *model.Project {
	ret := &model.Project{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *projectRepository) QueryCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Project, err error) {
	err = cnd.Exec(db).Find(&list).Error
	return
}

func (this *projectRepository) Query(db *gorm.DB, params *simple.QueryParams) (list []model.Project, paging *simple.Paging) {
	params.StartQuery(db).Find(&list)
    params.StartCount(db).Model(&model.Project{}).Count(&params.Paging.Total)
	paging = params.Paging
	return
}

func (this *projectRepository) Create(db *gorm.DB, t *model.Project) (err error) {
	err = db.Create(t).Error
	return
}

func (this *projectRepository) Update(db *gorm.DB, t *model.Project) (err error) {
	err = db.Save(t).Error
	return
}

func (this *projectRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Project{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *projectRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Project{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *projectRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.Project{}).Delete("id", id)
}

