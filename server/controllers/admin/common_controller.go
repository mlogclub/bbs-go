package admin

import (
	"runtime"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
)

type CommonController struct {
	Ctx iris.Context
}

func (c *CommonController) GetSysteminfo() *simple.JsonResult {
	return simple.NewEmptyRspBuilder().
		Put("os", runtime.GOOS).
		Put("arch", runtime.GOARCH).
		Put("numCpu", runtime.NumCPU()).
		Put("goroot", runtime.GOROOT()).
		Put("goversion", runtime.Version()).
		JsonResult()
}
