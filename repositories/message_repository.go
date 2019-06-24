
package repositories

import (
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/simple"
	"github.com/jinzhu/gorm"
)

type MessageRepository struct {
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{}
}

func (this *MessageRepository) Get(db *gorm.DB, id int64) *model.Message {
	ret := &model.Message{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *MessageRepository) Take(db *gorm.DB, where ...interface{}) *model.Message {
	ret := &model.Message{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *MessageRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.Message, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *MessageRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.Message, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
    queries.StartCount(db).Model(&model.Message{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *MessageRepository) Create(db *gorm.DB, t *model.Message) (err error) {
	err = db.Create(t).Error
	return
}

func (this *MessageRepository) Update(db *gorm.DB, t *model.Message) (err error) {
	err = db.Save(t).Error
	return
}

func (this *MessageRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Message{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *MessageRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Message{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *MessageRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.Message{}).Delete("id", id)
}

