package ginx

import (
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

func Bind(ctx *gin.Context, obj any) error {
	return ctx.ShouldBind(obj)
}

func BindJSON(ctx *gin.Context, obj any) error {
	return ctx.ShouldBindJSON(obj)
}

func BindQuery(ctx *gin.Context, obj any) error {
	return ctx.ShouldBindQuery(obj)
}

func BindForm(ctx *gin.Context, obj any) error {
	return ctx.ShouldBind(obj)
}

func FormValues(ctx *gin.Context) url.Values {
	_ = ctx.Request.ParseMultipartForm(32 << 20)
	_ = ctx.Request.ParseForm()
	values := url.Values{}
	for k, v := range ctx.Request.URL.Query() {
		values[k] = append(values[k], v...)
	}
	for k, v := range ctx.Request.PostForm {
		values[k] = append(values[k], v...)
	}
	return values
}

func FormValue(ctx *gin.Context, name string) string {
	if value := ctx.Query(name); value != "" {
		return value
	}
	return ctx.PostForm(name)
}

func FormValueDefault(ctx *gin.Context, name, def string) string {
	if value := FormValue(ctx, name); value != "" {
		return value
	}
	return def
}

func GetCookie(ctx *gin.Context, name string) string {
	value, _ := ctx.Cookie(name)
	return value
}

func SetCookieKV(ctx *gin.Context, name, value string, opts ...CookieOption) {
	options := cookieOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	maxAge := 0
	if options.expires > 0 {
		maxAge = int(options.expires.Seconds())
	}
	ctx.SetCookie(name, value, maxAge, "/", "", false, options.httpOnly)
}

func RemoveCookie(ctx *gin.Context, name string) {
	ctx.SetCookie(name, "", -1, "/", "", false, true)
}

type CookieOption func(*cookieOptions)

type cookieOptions struct {
	httpOnly bool
	expires  time.Duration
}

func CookieHTTPOnly(enabled bool) CookieOption {
	return func(opts *cookieOptions) {
		opts.httpOnly = enabled
	}
}

func CookieExpires(duration time.Duration) CookieOption {
	return func(opts *cookieOptions) {
		opts.expires = duration
	}
}
