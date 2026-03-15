package api

import (
	"bytes"
	"io"
	"log/slog"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"

	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/services"
)

type UploadController struct {
	Ctx iris.Context
}

func (c *UploadController) Post() *web.JsonResult {
	user := common.GetCurrentUser(c.Ctx)
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
	slog.Info("上传文件：", slog.Any("filename", header.Filename), slog.Any("size", header.Size))

	var body io.Reader
	var size int64
	if header.Size > 0 {
		body, size = file, header.Size
	} else {
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return web.JsonError(err)
		}
		body = bytes.NewReader(fileBytes)
		size = int64(len(fileBytes))
	}

	url, err := services.UploadService.PutImageStream(body, size, contentType)
	if err != nil {
		return web.JsonError(err)
	}
	return web.NewEmptyRspBuilder().Put("url", url).JsonResult()
}
