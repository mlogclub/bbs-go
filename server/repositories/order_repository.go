package repositories

import (
	"bbs-go/model"
	"gorm.io/gorm"
)

var OrderRepository = newOrderRepository()

func newOrderRepository() *orderRepository {
	return &orderRepository{}
}

type orderRepository struct {
}

func (r *orderRepository) Get(db *gorm.DB, orderId string) *model.Order {
	ret := &model.Order{}
	if err := db.First(ret, "out_trade_no = ?", orderId).Error; err != nil {
		return nil
	}
	return ret
}

func (r *orderRepository) Create(db *gorm.DB, t *model.Order) (err error) {
	err = db.Create(t).Error
	return
}

func (r *orderRepository) UpdateColumn(db *gorm.DB, orderId string, name string, value interface{}) (err error) {
	err = db.Model(&model.Order{}).Where("out_trade_no = ?", orderId).UpdateColumn(name, value).Error
	return
}
