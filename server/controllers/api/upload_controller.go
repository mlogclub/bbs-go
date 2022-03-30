package api

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/uploader"
	"io/ioutil"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/sirupsen/logrus"

	"bbs-go/services"
)

type UploadController struct {
	Ctx iris.Context
}

func (c *UploadController) Post() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	file, header, err := c.Ctx.FormFile("image")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	defer file.Close()

	if header.Size > constants.UploadMaxBytes {
		return web.JsonErrorMsg("图片不能超过" + strconv.Itoa(constants.UploadMaxM) + "M")
	}

	contentType := header.Header.Get("Content-Type")
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	logrus.Info("上传文件：", header.Filename, " size:", header.Size)

	url, err := uploader.PutImage(fileBytes, contentType)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.NewEmptyRspBuilder().Put("url", url).JsonResult()
}
