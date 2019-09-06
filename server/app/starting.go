package app

import (
	"github.com/mlogclub/bbs-go/common"
)

func StartOn() {
	if !common.IsProd() {
		return
	}

	// 开启定时任务
	startSchedule()
}
