package api

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/uploader"
	"io"
	"log/slog"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"

	"bbs-go/internal/services"
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
		return web.JsonError(err)
	}
	defer file.Close()

	if header.Size > constants.UploadMaxBytes {
		return web.JsonErrorMsg("图片不能超过" + strconv.Itoa(constants.UploadMaxM) + "M")
	}

	contentType := header.Header.Get("Content-Type")
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return web.JsonError(err)
	}

	slog.Info("上传文件：", slog.Any("filename", header.Filename), slog.Any("size", header.Size))

	url, err := uploader.PutImage(fileBytes, contentType)
	if err != nil {
		return web.JsonError(err)
	}
	return web.NewEmptyRspBuilder().Put("url", url).JsonResult()
}
