package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
)

var GithubUserRepository = newGithubUserRepository()

func newGithubUserRepository() *githubUserRepository {
	return &githubUserRepository{}
}

type githubUserRepository struct {
}

func (this *githubUserRepository) Get(db *gorm.DB, id int64) *model.GithubUser {
	ret := &model.GithubUser{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *githubUserRepository) Take(db *gorm.DB, where ...interface{}) *model.GithubUser {
	ret := &model.GithubUser{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *githubUserRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.GithubUser, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *githubUserRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.GithubUser, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.GithubUser{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *githubUserRepository) Create(db *gorm.DB, t *model.GithubUser) (err error) {
	err = db.Create(t).Error
	return
}

func (this *githubUserRepository) Update(db *gorm.DB, t *model.GithubUser) (err error) {
	err = db.Save(t).Error
	return
}

func (this *githubUserRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.GithubUser{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *githubUserRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.GithubUser{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *githubUserRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.GithubUser{}, "id = ?", id)
}

func (this *githubUserRepository) GetByGithubId(db *gorm.DB, githubId int64) *model.GithubUser {
	return this.Take(db, "github_id = ?", githubId)
}
