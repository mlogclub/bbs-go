package repositories

import (
	"bbs-go/model"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var StickyTopicRepository = newStickyTopicRepository()

func newStickyTopicRepository() *stickyTopicRepository {
	return &stickyTopicRepository{}
}

type stickyTopicRepository struct {
}

func (r *stickyTopicRepository) Get(db *gorm.DB, id int64) *model.StickyTopic {
	ret := &model.StickyTopic{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *stickyTopicRepository) Take(db *gorm.DB, where ...interface{}) *model.StickyTopic {
	ret := &model.StickyTopic{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *stickyTopicRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.StickyTopic) {
	cnd.Find(db, &list)
	return
}

func (r *stickyTopicRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.StickyTopic {
	ret := &model.StickyTopic{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *stickyTopicRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []model.StickyTopic, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *stickyTopicRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []model.StickyTopic, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.StickyTopic{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *stickyTopicRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (list []model.StickyTopic) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *stickyTopicRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *stickyTopicRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &model.StickyTopic{})
}

func (r *stickyTopicRepository) Create(db *gorm.DB, t *model.StickyTopic) (err error) {
	err = db.Create(t).Error
	return
}

func (r *stickyTopicRepository) Update(db *gorm.DB, t *model.StickyTopic) (err error) {
	err = db.Save(t).Error
	return
}

func (r *stickyTopicRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.StickyTopic{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *stickyTopicRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.StickyTopic{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *stickyTopicRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.StickyTopic{}, "id = ?", id)
}
