package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/common/strs"

	"bbs-go/internal/models/req"
	"bbs-go/internal/models/resp"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/services"
)

// attachmentErrorHTML 返回附件下载异常时的友好 HTML 页面
func attachmentErrorHTML(title, message string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>%s</title>
<style>
  *{box-sizing:border-box}
  body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;max-width:420px;margin:0 auto;padding:48px 24px;min-height:100vh;display:flex;align-items:center;justify-content:center;background:#f8f9fa;color:#1a1a1a}
  .card{background:#fff;border-radius:12px;padding:32px 28px;text-align:center;box-shadow:0 1px 3px rgba(0,0,0,.06)}
  .card .title{font-size:1.125rem;font-weight:600;margin:0 0 12px;color:#1a1a1a}
  .card .msg{font-size:0.9375rem;line-height:1.6;margin:0;color:#555}
</style>
</head>
<body><div class="card"><p class="title">%s</p><p class="msg">%s</p></div></body>
</html>`, title, title, message)
}

// PostUpload 上传附件（发帖前或发帖时）
func AttachmentUpload(ctx *gin.Context) {
	user, err := common.CheckLogin(ctx)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	if err := services.UserService.CheckPostStatus(user); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	cfg := services.SysConfigService.GetAttachmentConfig()
	if !cfg.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("attachment.disabled")))
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	defer file.Close()

	maxBytes := int64(cfg.MaxSizeMB) * 1024 * 1024
	if header.Size > maxBytes {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Getf("attachment.too_large", cfg.MaxSizeMB)))
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	downloadScore, _ := params.GetInt(ctx, "downloadScore")

	var (
		body io.Reader = file
		size int64     = header.Size
	)
	// 客户端未提供 Size 时回退为读入内存（如 chunked 上传）
	if size <= 0 {
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			ginx.WriteJSON(ctx, err)
			return
		}
		body = bytes.NewReader(fileBytes)
		size = int64(len(fileBytes))
	}

	att, err := services.AttachmentService.Upload(user.Id, header.Filename, body, size, contentType, downloadScore)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	ginx.WriteJSON(ctx, resp.AttachmentResponse{
		Id:            att.Id,
		FileName:      att.FileName,
		FileSize:      att.FileSize,
		DownloadScore: att.DownloadScore,
		Downloaded:    false,
	})

}

// GetDownloadBy 下载附件：鉴权后 302 到实际地址（id 为附件 UUID）。
func AttachmentDownload(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := common.CheckLogin(ctx)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		msg := locales.Get("errors.not_login")
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.Writer.WriteString(attachmentErrorHTML(msg, msg))
		return
	}

	if strs.IsBlank(id) {
		ctx.Status(http.StatusNotFound)
		msg := locales.Get("attachment.not_found")
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.Writer.WriteString(attachmentErrorHTML(msg, msg))
		return
	}

	redirectURL, err := services.AttachmentService.Download(id, user.Id)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		msg := locales.Get(err.Error())
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.Writer.WriteString(attachmentErrorHTML(msg, msg))
		return
	}
	if strs.IsBlank(redirectURL) {
		ctx.Status(http.StatusNotFound)
		msg := locales.Get("attachment.file_missing")
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.Writer.WriteString(attachmentErrorHTML(msg, msg))
		return
	}

	ctx.Redirect(302, redirectURL)
}

// PostUpdateDownloadScore 更新附件下载积分
func AttachmentUpdateDownloadScore(ctx *gin.Context) {
	user, err := common.CheckLogin(ctx)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	var body req.PatchDownloadScoreReq
	if err := ginx.BindJSON(ctx, &body); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("invalid body"))
		return
	}
	att, err := services.AttachmentService.UpdateDownloadScore(body.Id, user.Id, body.DownloadScore)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, resp.AttachmentResponse{
		Id:            att.Id,
		FileName:      att.FileName,
		FileSize:      att.FileSize,
		DownloadScore: att.DownloadScore,
	})

}
