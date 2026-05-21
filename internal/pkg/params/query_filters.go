package params

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/common/strs/strcase"
	"github.com/mlogclub/simple/sqls"
)

type QueryOp string

const (
	Eq       QueryOp = "eq"
	Gt       QueryOp = "gt"
	Lt       QueryOp = "lt"
	Gte      QueryOp = "gte"
	Lte      QueryOp = "lte"
	Like     QueryOp = "like"
	In       QueryOp = "in"
	Starting QueryOp = "starting"
	Ending   QueryOp = "ending"
)

type QueryFilter struct {
	ParamName    string
	Op           QueryOp
	ColumnName   string
	ValueWrapper func(origin string) string
}

func NewPagedSqlCnd(ctx *gin.Context, filters ...QueryFilter) *sqls.Cnd {
	cnd := NewSqlCnd(ctx, filters...)
	p := GetPaging(ctx)
	cnd.Page(p.Page, p.Limit)
	return cnd
}

func NewSqlCnd(ctx *gin.Context, filters ...QueryFilter) *sqls.Cnd {
	cnd := sqls.NewCnd()
	for _, filter := range filters {
		columnName := filter.ColumnName
		paramValue := QueryValue(ctx, filter.ParamName)
		if strs.IsBlank(string(filter.Op)) {
			filter.Op = Eq
		}
		if filter.ValueWrapper != nil {
			paramValue = filter.ValueWrapper(paramValue)
		}
		if strs.IsBlank(paramValue) {
			continue
		}
		if strs.IsBlank(columnName) {
			columnName = strcase.ToSnake(filter.ParamName)
		}
		switch filter.Op {
		case Eq:
			cnd.Eq(columnName, paramValue)
		case Gt:
			cnd.Gt(columnName, paramValue)
		case Lt:
			cnd.Lt(columnName, paramValue)
		case Gte:
			cnd.Gte(columnName, paramValue)
		case Lte:
			cnd.Lte(columnName, paramValue)
		case Like:
			cnd.Like(columnName, paramValue)
		case Starting:
			cnd.Starting(columnName, paramValue)
		case Ending:
			cnd.Ending(columnName, paramValue)
		case In:
			cnd.In(columnName, strings.Split(paramValue, ","))
		}
	}
	return cnd
}
