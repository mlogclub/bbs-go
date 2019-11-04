package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
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

func (this *projectRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Project) {
	cnd.Find(db, &list)
	return
}

func (this *projectRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.Project) {
	cnd.FindOne(db, &ret)
	return
}

func (this *projectRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.Project, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *projectRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Project, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Project{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
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
