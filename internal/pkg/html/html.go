package html

import (
	"log/slog"
	"strings"

	"bbs-go/internal/pkg/text"

	"github.com/PuerkitoBio/goquery"
	"github.com/mlogclub/simple/common/strs"
)

func GetSummary(htmlStr string, summaryLen int) string {
	if summaryLen <= 0 || strs.IsEmpty(htmlStr) {
		return ""
	}
	return text.GetSummary(GetHtmlText(htmlStr), summaryLen)
}

// GetHtmlText 获取html文本
func GetHtmlText(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return ""
	}
	return doc.Text()
}
