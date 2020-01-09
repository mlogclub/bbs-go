package simple

import (
	"github.com/kataras/iris/v12"

	"github.com/mlogclub/simple/strcase"
)

type QueryParams struct {
	Ctx iris.Context
	SqlCnd
}

func NewQueryParams(ctx iris.Context) *QueryParams {
	return &QueryParams{
		Ctx: ctx,
	}
}

func (q *QueryParams) getValueByColumn(column string) string {
	if q.Ctx == nil {
		return ""
	}
	fieldName := strcase.ToLowerCamel(column)
	return q.Ctx.FormValue(fieldName)
}

func (q *QueryParams) EqByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Eq(column, value)
	}
	return q
}

func (q *QueryParams) NotEqByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.NotEq(column, value)
	}
	return q
}

func (q *QueryParams) GtByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Gt(column, value)
	}
	return q
}

func (q *QueryParams) GteByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Gte(column, value)
	}
	return q
}

func (q *QueryParams) LtByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Lt(column, value)
	}
	return q
}

func (q *QueryParams) LteByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Lte(column, value)
	}
	return q
}

func (q *QueryParams) LikeByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Like(column, value)
	}
	return q
}

func (q *QueryParams) PageByReq() *QueryParams {
	if q.Ctx == nil {
		return q
	}
	paging := GetPaging(q.Ctx)
	q.Page(paging.Page, paging.Limit)
	return q
}

func (q *QueryParams) Asc(column string) *QueryParams {
	q.Orders = append(q.Orders, OrderByCol{Column: column, Asc: true})
	return q
}

func (q *QueryParams) Desc(column string) *QueryParams {
	q.Orders = append(q.Orders, OrderByCol{Column: column, Asc: false})
	return q
}

func (q *QueryParams) Limit(limit int) *QueryParams {
	q.Page(1, limit)
	return q
}

func (q *QueryParams) Page(page, limit int) *QueryParams {
	if q.Paging == nil {
		q.Paging = &Paging{Page: page, Limit: limit}
	} else {
		q.Paging.Page = page
		q.Paging.Limit = limit
	}
	return q
}
