package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var VoteRepository = newVoteRepository()

func newVoteRepository() *voteRepository {
	return &voteRepository{}
}

type voteRepository struct {
}

func (r *voteRepository) Get(db *gorm.DB, id int64) *models.Vote {
	ret := &models.Vote{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *voteRepository) Take(db *gorm.DB, where ...interface{}) *models.Vote {
	ret := &models.Vote{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *voteRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.Vote) {
	cnd.Find(db, &list)
	return
}

func (r *voteRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.Vote {
	ret := &models.Vote{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *voteRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.Vote, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *voteRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.Vote, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Vote{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *voteRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (list []models.Vote) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *voteRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *voteRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.Vote{})
}

func (r *voteRepository) Create(db *gorm.DB, t *models.Vote) (err error) {
	err = db.Create(t).Error
	return
}

func (r *voteRepository) Update(db *gorm.DB, t *models.Vote) (err error) {
	err = db.Save(t).Error
	return
}

func (r *voteRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Vote{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *voteRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Vote{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *voteRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Vote{}, "id = ?", id)
}

