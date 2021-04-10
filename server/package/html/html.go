package html

import (
	"github.com/mlogclub/simple"
)

func GetSummary(htmlStr string, summaryLen int) string {
	if summaryLen <= 0 || simple.IsEmpty(htmlStr) {
		return ""
	}
	text := simple.GetHtmlText(htmlStr)
	return simple.GetSummary(text, summaryLen)
}
