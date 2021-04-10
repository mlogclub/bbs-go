package app

import (
	"bbs-go/package/common"
)

func StartOn() {
	if !common.IsProd() {
		return
	}

	// 开启定时任务
	startSchedule()
}
