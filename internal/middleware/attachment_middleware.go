package middleware

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"bbs-go/internal/models/constants"
	"bbs-go/internal/repositories"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/sqls"
)

// AttachmentMiddleware 对 /res/uploads/attachments/* 的请求设置 Content-Disposition: attachment，触发浏览器下载而非预览
// 从 path 中解析附件 ID（存储 key 格式为 prefix/2006/01/02/uuid.ext，uuid 即附件 ID），查库后用附件的 FileName 作为下载文件名
func AttachmentMiddleware(ctx iris.Context) {
	path := ctx.Path()
	if ctx.Method() != http.MethodGet && ctx.Method() != http.MethodHead {
		ctx.Next()
		return
	}
	if !strings.HasPrefix(path, "/res/uploads/") {
		ctx.Next()
		return
	}
	if !strings.Contains(path, "/attachments/") {
		ctx.Next()
		return
	}

	base := filepath.Base(path)
	if base == "" || base == "." {
		ctx.Next()
		return
	}
	// 存储 key 格式：attachments/2006/01/02/uuid.ext 或 test/attachments/2006/01/02/uuid.ext，最后一段为 uuid.ext
	attachmentId := strings.TrimSuffix(base, filepath.Ext(base))
	downloadName := base
	if attachmentId != "" {
		att := repositories.AttachmentRepository.Get(sqls.DB(), attachmentId)
		if att != nil && att.Status == constants.StatusOk && att.FileName != "" {
			downloadName = filepath.Base(att.FileName)
		}
	}
	ctx.ResponseWriter().Header().Set("Content-Disposition", `attachment; filename="`+url.QueryEscape(downloadName)+`"`)

	ctx.Next()
}
