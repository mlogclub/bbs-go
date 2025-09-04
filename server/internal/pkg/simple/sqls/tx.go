package sqls

import (
	"gorm.io/gorm"
)

// CallbackFunc 定义回调函数类型
type CallbackFunc func()

// RegisterCallbackFunc 定义注册回调函数类型
type RegisterCallbackFunc func(CallbackFunc)

type TxFunc func(ctx *TxContext) error

type TxContext struct {
	Tx               *gorm.DB
	RegisterCallback RegisterCallbackFunc
}

func WithTransaction(fn TxFunc) error {
	var callbacks []CallbackFunc

	registerCallback := func(fn CallbackFunc) {
		callbacks = append(callbacks, fn)
	}

	err := DB().Transaction(func(tx *gorm.DB) error {
		ctx := &TxContext{
			Tx:               tx,
			RegisterCallback: registerCallback,
		}
		return fn(ctx)
	})

	if err == nil {
		for _, callback := range callbacks {
			if callback != nil {
				callback()
			}
		}
	}

	return err
}
