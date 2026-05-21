package api

import (
	"bytes"
	"io"
	"log/slog"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"

	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/services"
)

func UploadHandle(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	defer file.Close()

	if header.Size > constants.UploadMaxBytes {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Getf("upload.image_too_large", constants.UploadMaxM)))
		return
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
			ginx.WriteJSON(ctx, err)
			return
		}
		body = bytes.NewReader(fileBytes)
		size = int64(len(fileBytes))
	}

	url, err := services.UploadService.PutImageStream(body, size, contentType)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, map[string]any{"url": url})

}
