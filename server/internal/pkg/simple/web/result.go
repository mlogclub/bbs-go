package web

import "bbs-go/internal/pkg/simple/sqls"

// PageResult 分页返回数据
type PageResult struct {
	Page    *sqls.Paging `json:"page"`    // 分页信息
	Results interface{}  `json:"results"` // 数据
}

// CursorResult Cursor分页返回数据
type CursorResult struct {
	Results interface{} `json:"results"` // 数据
	Cursor  string      `json:"cursor"`  // 下一页
	HasMore bool        `json:"hasMore"` // 是否还有数据
}
