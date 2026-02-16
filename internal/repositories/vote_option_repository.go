package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var VoteOptionRepository = newVoteOptionRepository()

func newVoteOptionRepository() *voteOptionRepository {
	return &voteOptionRepository{}
}

type voteOptionRepository struct {
}

func (r *voteOptionRepository) Get(db *gorm.DB, id int64) *models.VoteOption {
	ret := &models.VoteOption{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *voteOptionRepository) Take(db *gorm.DB, where ...interface{}) *models.VoteOption {
	ret := &models.VoteOption{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *voteOptionRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.VoteOption) {
	cnd.Find(db, &list)
	return
}

func (r *voteOptionRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.VoteOption {
	ret := &models.VoteOption{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *voteOptionRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.VoteOption, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *voteOptionRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.VoteOption, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.VoteOption{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *voteOptionRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (list []models.VoteOption) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *voteOptionRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *voteOptionRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.VoteOption{})
}

func (r *voteOptionRepository) Create(db *gorm.DB, t *models.VoteOption) (err error) {
	err = db.Create(t).Error
	return
}

func (r *voteOptionRepository) Update(db *gorm.DB, t *models.VoteOption) (err error) {
	err = db.Save(t).Error
	return
}

func (r *voteOptionRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.VoteOption{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *voteOptionRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.VoteOption{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *voteOptionRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.VoteOption{}, "id = ?", id)
}

