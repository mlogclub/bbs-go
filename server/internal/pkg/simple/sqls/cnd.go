package sqls

import (
	"log/slog"

	"gorm.io/gorm"
)

type Cnd struct {
	SelectCols []string     // 要查询的字段，如果为空，表示查询所有字段
	Params     []ParamPair  // 参数
	Orders     []OrderByCol // 排序
	Paging     *Paging      // 分页
}

type ParamPair struct {
	Query string        // 查询
	Args  []interface{} // 参数
}

// OrderByCol 排序信息
type OrderByCol struct {
	Column string // 排序字段
	Asc    bool   // 是否正序
}

func NewCnd() *Cnd {
	return &Cnd{}
}

func (s *Cnd) Cols(selectCols ...string) *Cnd {
	if len(selectCols) > 0 {
		s.SelectCols = append(s.SelectCols, selectCols...)
	}
	return s
}

func (s *Cnd) Eq(column string, args ...interface{}) *Cnd {
	s.Where(KeywordWrap(column)+" = ?", args)
	return s
}

func (s *Cnd) NotEq(column string, args ...interface{}) *Cnd {
	s.Where(KeywordWrap(column)+" <> ?", args)
	return s
}

func (s *Cnd) Gt(column string, args ...interface{}) *Cnd {
	s.Where(KeywordWrap(column)+" > ?", args)
	return s
}

func (s *Cnd) Gte(column string, args ...interface{}) *Cnd {
	s.Where(KeywordWrap(column)+" >= ?", args)
	return s
}

func (s *Cnd) Lt(column string, args ...interface{}) *Cnd {
	s.Where(KeywordWrap(column)+" < ?", args)
	return s
}

func (s *Cnd) Lte(column string, args ...interface{}) *Cnd {
	s.Where(KeywordWrap(column)+" <= ?", args)
	return s
}

func (s *Cnd) Like(column string, str string) *Cnd {
	s.Where(KeywordWrap(column)+" LIKE ?", "%"+str+"%")
	return s
}

func (s *Cnd) Starting(column string, str string) *Cnd {
	s.Where(KeywordWrap(column)+" LIKE ?", str+"%")
	return s
}

func (s *Cnd) Ending(column string, str string) *Cnd {
	s.Where(KeywordWrap(column)+" LIKE ?", "%"+str)
	return s
}

func (s *Cnd) In(column string, params interface{}) *Cnd {
	s.Where(KeywordWrap(column)+" in (?) ", params)
	return s
}

func (s *Cnd) NotIn(column string, params interface{}) *Cnd {
	s.Where(KeywordWrap(column)+" not in (?) ", params)
	return s
}

func (s *Cnd) Where(query string, args ...interface{}) *Cnd {
	s.Params = append(s.Params, ParamPair{query, args})
	return s
}

func (s *Cnd) Asc(column string) *Cnd {
	s.Orders = append(s.Orders, OrderByCol{Column: KeywordWrap(column), Asc: true})
	return s
}

func (s *Cnd) Desc(column string) *Cnd {
	s.Orders = append(s.Orders, OrderByCol{Column: KeywordWrap(column), Asc: false})
	return s
}

func (s *Cnd) Limit(limit int) *Cnd {
	s.Page(1, limit)
	return s
}

func (s *Cnd) Page(page, limit int) *Cnd {
	if s.Paging == nil {
		s.Paging = &Paging{Page: page, Limit: limit}
	} else {
		s.Paging.Page = page
		s.Paging.Limit = limit
	}
	return s
}

func (s *Cnd) Build(db *gorm.DB) *gorm.DB {
	ret := db

	if len(s.SelectCols) > 0 {
		cols := make([]string, len(s.SelectCols))
		for i, col := range s.SelectCols {
			cols[i] = KeywordWrap(col)
		}
		ret = ret.Select(cols)
	}

	// where
	if len(s.Params) > 0 {
		for _, param := range s.Params {
			ret = ret.Where(param.Query, param.Args...)
		}
	}

	// order
	if len(s.Orders) > 0 {
		for _, order := range s.Orders {
			if order.Asc {
				ret = ret.Order(order.Column + " ASC")
			} else {
				ret = ret.Order(order.Column + " DESC")
			}
		}
	}

	// limit
	if s.Paging != nil && s.Paging.Limit > 0 {
		ret = ret.Limit(s.Paging.Limit)
	}

	// offset
	if s.Paging != nil && s.Paging.Offset() > 0 {
		ret = ret.Offset(s.Paging.Offset())
	}
	return ret
}

func (s *Cnd) Find(db *gorm.DB, out interface{}) {
	if err := s.Build(db).Find(out).Error; err != nil {
		slog.Error(err.Error(), slog.Any("error", err))
	}
}

func (s *Cnd) FindOne(db *gorm.DB, out interface{}) error {
	if err := s.Limit(1).Build(db).First(out).Error; err != nil {
		return err
	}
	return nil
}

func (s *Cnd) Count(db *gorm.DB, model interface{}) int64 {
	ret := db.Model(model)

	// where
	if len(s.Params) > 0 {
		for _, query := range s.Params {
			ret = ret.Where(query.Query, query.Args...)
		}
	}

	var count int64
	if err := ret.Count(&count).Error; err != nil {
		slog.Error(err.Error(), slog.Any("error", err))
	}
	return count
}
