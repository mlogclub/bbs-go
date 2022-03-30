package api

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/uploader"
	"io/ioutil"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/mvc"
	"github.com/sirupsen/logrus"

	"bbs-go/services"
)

type UploadController struct {
	Ctx iris.Context
}

func (c *UploadController) Post() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return mvc.JsonError(err)
	}

	file, header, err := c.Ctx.FormFile("image")
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	defer file.Close()

	if header.Size > constants.UploadMaxBytes {
		return mvc.JsonErrorMsg("图片不能超过" + strconv.Itoa(constants.UploadMaxM) + "M")
	}

	contentType := header.Header.Get("Content-Type")
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}

	logrus.Info("上传文件：", header.Filename, " size:", header.Size)

	url, err := uploader.PutImage(fileBytes, contentType)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.NewEmptyRspBuilder().Put("url", url).JsonResult()
}
