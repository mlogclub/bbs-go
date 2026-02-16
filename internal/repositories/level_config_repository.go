package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var LevelConfigRepository = newLevelConfigRepository()

func newLevelConfigRepository() *levelConfigRepository {
	return &levelConfigRepository{}
}

type levelConfigRepository struct {
}

func (r *levelConfigRepository) Get(db *gorm.DB, id int64) *models.LevelConfig {
	ret := &models.LevelConfig{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *levelConfigRepository) Take(db *gorm.DB, where ...interface{}) *models.LevelConfig {
	ret := &models.LevelConfig{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *levelConfigRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.LevelConfig) {
	cnd.Find(db, &list)
	return
}

func (r *levelConfigRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.LevelConfig {
	ret := &models.LevelConfig{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *levelConfigRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.LevelConfig, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *levelConfigRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.LevelConfig, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.LevelConfig{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *levelConfigRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (list []models.LevelConfig) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *levelConfigRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *levelConfigRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.LevelConfig{})
}

func (r *levelConfigRepository) Create(db *gorm.DB, t *models.LevelConfig) (err error) {
	err = db.Create(t).Error
	return
}

func (r *levelConfigRepository) Update(db *gorm.DB, t *models.LevelConfig) (err error) {
	err = db.Save(t).Error
	return
}

func (r *levelConfigRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.LevelConfig{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *levelConfigRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.LevelConfig{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *levelConfigRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.LevelConfig{}, "id = ?", id)
}

