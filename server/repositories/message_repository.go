package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/simple"
)

var MessageRepository = newMessageRepository()

func newMessageRepository() *messageRepository {
	return &messageRepository{}
}

type messageRepository struct {
}

func (this *messageRepository) Get(db *gorm.DB, id int64) *model.Message {
	ret := &model.Message{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *messageRepository) Take(db *gorm.DB, where ...interface{}) *model.Message {
	ret := &model.Message{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *messageRepository) QueryCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Message, err error) {
	err = cnd.Exec(db).Find(&list).Error
	return
}

func (this *messageRepository) Query(db *gorm.DB, params *simple.QueryParams) (list []model.Message, paging *simple.Paging) {
	params.StartQuery(db).Find(&list)
	params.StartCount(db).Model(&model.Message{}).Count(&params.Paging.Total)
	paging = params.Paging
	return
}

func (this *messageRepository) Create(db *gorm.DB, t *model.Message) (err error) {
	err = db.Create(t).Error
	return
}

func (this *messageRepository) Update(db *gorm.DB, t *model.Message) (err error) {
	err = db.Save(t).Error
	return
}

func (this *messageRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Message{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *messageRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Message{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *messageRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Message{}, "id = ?", id)
}
