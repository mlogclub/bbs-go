package api

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/models/resp"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/services"
)

type AttachmentController struct {
	Ctx iris.Context
}

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
func (c *AttachmentController) PostUpload() *web.JsonResult {
	user, err := common.CheckLogin(c.Ctx)
	if err != nil {
		return web.JsonError(err)
	}

	if err := services.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	cfg := services.SysConfigService.GetAttachmentConfig()
	if !cfg.Enabled {
		return web.JsonErrorMsg(locales.Get("attachment.disabled"))
	}

	file, header, err := c.Ctx.FormFile("file")
	if err != nil {
		return web.JsonError(err)
	}
	defer file.Close()

	maxBytes := int64(cfg.MaxSizeMB) * 1024 * 1024
	if header.Size > maxBytes {
		return web.JsonErrorMsg(locales.Getf("attachment.too_large", cfg.MaxSizeMB))
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	downloadScore, _ := params.GetInt(c.Ctx, "downloadScore")

	var (
		body io.Reader = file
		size int64     = header.Size
	)
	// 客户端未提供 Size 时回退为读入内存（如 chunked 上传）
	if size <= 0 {
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return web.JsonError(err)
		}
		body = bytes.NewReader(fileBytes)
		size = int64(len(fileBytes))
	}

	att, err := services.AttachmentService.Upload(user.Id, header.Filename, body, size, contentType, downloadScore)
	if err != nil {
		return web.JsonError(err)
	}

	return web.JsonData(resp.AttachmentResponse{
		Id:            att.Id,
		FileName:      att.FileName,
		FileSize:      att.FileSize,
		DownloadScore: att.DownloadScore,
		Downloaded:    false,
	})
}

// GetDownloadBy 下载附件：鉴权后 302 到实际地址（id 为附件 UUID）。
func (c *AttachmentController) GetDownloadBy(id string) {
	user, err := common.CheckLogin(c.Ctx)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusUnauthorized)
		msg := locales.Get("errors.not_login")
		c.Ctx.ContentType("text/html; charset=utf-8")
		c.Ctx.WriteString(attachmentErrorHTML(msg, msg))
		return
	}

	if strs.IsBlank(id) {
		c.Ctx.StatusCode(iris.StatusNotFound)
		msg := locales.Get("attachment.not_found")
		c.Ctx.ContentType("text/html; charset=utf-8")
		c.Ctx.WriteString(attachmentErrorHTML(msg, msg))
		return
	}

	redirectURL, err := services.AttachmentService.Download(id, user.Id)
	if err != nil {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		msg := locales.Get(err.Error())
		c.Ctx.ContentType("text/html; charset=utf-8")
		c.Ctx.WriteString(attachmentErrorHTML(msg, msg))
		return
	}
	if strs.IsBlank(redirectURL) {
		c.Ctx.StatusCode(iris.StatusNotFound)
		msg := locales.Get("attachment.file_missing")
		c.Ctx.ContentType("text/html; charset=utf-8")
		c.Ctx.WriteString(attachmentErrorHTML(msg, msg))
		return
	}

	c.Ctx.Redirect(redirectURL, 302)
}

// PostUpdateDownloadScore 更新附件下载积分
func (c *AttachmentController) PostUpdate_download_score() *web.JsonResult {
	user, err := common.CheckLogin(c.Ctx)
	if err != nil {
		return web.JsonError(err)
	}

	type PatchDownloadScoreReq struct {
		Id            string `json:"id"`
		DownloadScore int    `json:"downloadScore"`
	}

	var body PatchDownloadScoreReq
	if err := c.Ctx.ReadJSON(&body); err != nil {
		return web.JsonErrorMsg("invalid body")
	}
	att, err := services.AttachmentService.UpdateDownloadScore(body.Id, user.Id, body.DownloadScore)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(resp.AttachmentResponse{
		Id:            att.Id,
		FileName:      att.FileName,
		FileSize:      att.FileSize,
		DownloadScore: att.DownloadScore,
	})
}
