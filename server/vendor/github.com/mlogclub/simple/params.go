package simple

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/iris-contrib/schema"
	"github.com/kataras/iris/v12"
)

var (
	decoder = schema.NewDecoder() // form, url, schema.
)

func init() {
	decoder.AddAliasTag("form", "json")
	decoder.ZeroEmpty(true)
}

// param error
func paramError(name string) error {
	return errors.New(fmt.Sprintf("unable to find param value '%s'", name))
}

// ReadForm read object from FormData
func ReadForm(ctx iris.Context, obj interface{}) error {
	values := ctx.FormValues()
	if len(values) == 0 {
		return nil
	}
	decoder := schema.NewDecoder()
	decoder.ZeroEmpty(true)
	return decoder.Decode(obj, values)
}

func FormValue(ctx iris.Context, name string) string {
	return ctx.FormValue(name)
}

func FormValueRequired(ctx iris.Context, name string) (string, error) {
	str := FormValue(ctx, name)
	if len(str) == 0 {
		return "", errors.New("参数：" + name + "不能为空")
	}
	return str, nil
}

func FormValueDefault(ctx iris.Context, name, def string) string {
	return ctx.FormValueDefault(name, def)
}

func FormValueInt(ctx iris.Context, name string) (int, error) {
	str := ctx.FormValue(name)
	if str == "" {
		return 0, paramError(name)
	}
	return strconv.Atoi(str)
}

func FormValueIntDefault(ctx iris.Context, name string, def int) int {
	if v, err := FormValueInt(ctx, name); err == nil {
		return v
	}
	return def
}

func FormValueInt64(ctx iris.Context, name string) (int64, error) {
	str := ctx.FormValue(name)
	if str == "" {
		return 0, paramError(name)
	}
	return strconv.ParseInt(str, 10, 64)
}

func FormValueInt64Default(ctx iris.Context, name string, def int64) int64 {
	if v, err := FormValueInt64(ctx, name); err == nil {
		return v
	}
	return def
}

func FormValueInt64Array(ctx iris.Context, name string) []int64 {
	str := ctx.FormValue(name)
	if str == "" {
		return nil
	}
	ss := strings.Split(str, ",")
	if len(ss) == 0 {
		return nil
	}
	var ret []int64
	for _, v := range ss {
		item, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		ret = append(ret, item)
	}
	return ret
}

func FormValueStringArray(ctx iris.Context, name string) []string {
	str := ctx.FormValue(name)
	if len(str) == 0 {
		return nil
	}
	ss := strings.Split(str, ",")
	if len(ss) == 0 {
		return nil
	}
	var ret []string
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		ret = append(ret, s)
	}
	return ret
}

func FormValueBool(ctx iris.Context, name string) (bool, error) {
	str := ctx.FormValue(name)
	if str == "" {
		return false, paramError(name)
	}
	return strconv.ParseBool(str)
}

func GetPaging(ctx iris.Context) *Paging {
	page := FormValueIntDefault(ctx, "page", 1)
	limit := FormValueIntDefault(ctx, "limit", 20)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return &Paging{Page: page, Limit: limit}
}
