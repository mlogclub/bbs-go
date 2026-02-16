package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var VoteRecordRepository = newVoteRecordRepository()

func newVoteRecordRepository() *voteRecordRepository {
	return &voteRecordRepository{}
}

type voteRecordRepository struct {
}

func (r *voteRecordRepository) Get(db *gorm.DB, id int64) *models.VoteRecord {
	ret := &models.VoteRecord{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *voteRecordRepository) Take(db *gorm.DB, where ...interface{}) *models.VoteRecord {
	ret := &models.VoteRecord{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *voteRecordRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.VoteRecord) {
	cnd.Find(db, &list)
	return
}

func (r *voteRecordRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.VoteRecord {
	ret := &models.VoteRecord{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *voteRecordRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.VoteRecord, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *voteRecordRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.VoteRecord, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.VoteRecord{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *voteRecordRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (list []models.VoteRecord) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *voteRecordRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *voteRecordRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.VoteRecord{})
}

func (r *voteRecordRepository) Create(db *gorm.DB, t *models.VoteRecord) (err error) {
	err = db.Create(t).Error
	return
}

func (r *voteRecordRepository) Update(db *gorm.DB, t *models.VoteRecord) (err error) {
	err = db.Save(t).Error
	return
}

func (r *voteRecordRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.VoteRecord{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *voteRecordRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.VoteRecord{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *voteRecordRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.VoteRecord{}, "id = ?", id)
}

