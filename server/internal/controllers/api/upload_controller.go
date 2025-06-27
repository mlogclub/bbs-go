package api

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/locales"
	"io"
	"log/slog"

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
		return web.JsonErrorMsg(locales.Getf("upload.image_too_large", constants.UploadMaxM))
	}

	contentType := header.Header.Get("Content-Type")
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return web.JsonError(err)
	}

	slog.Info("上传文件：", slog.Any("filename", header.Filename), slog.Any("size", header.Size))

	url, err := services.UploadService.PutImage(fileBytes, contentType)
	if err != nil {
		return web.JsonError(err)
	}
	return web.NewEmptyRspBuilder().Put("url", url).JsonResult()
}
