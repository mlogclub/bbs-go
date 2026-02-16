package common

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/html"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/pkg/text"
)

func GetSummary(contentType constants.ContentType, content string) (summary string) {
	if contentType == constants.ContentTypeMarkdown {
		summary = markdown.GetSummary(content, constants.SummaryLen)
	} else if contentType == constants.ContentTypeHtml {
		summary = html.GetSummary(content, constants.SummaryLen)
	} else {
		summary = text.GetSummary(content, constants.SummaryLen)
	}
	return
}

// GetMarkdownSummary 截取markdown摘要
func GetMarkdownSummary(markdownStr string) string {
	return markdown.GetSummary(markdownStr, constants.SummaryLen)
}

func Distinct[T any](input []T, getKey func(T) any) (output []T) {
	tempMap := map[any]byte{}
	for _, item := range input {
		l := len(tempMap)
		tempMap[getKey(item)] = 0
		if len(tempMap) != l { // 数量发生变化，说明不存在
			output = append(output, item)
		}
	}
	return
}

func StrRight(str string, size int) string {
	if str == "" || size < 0 {
		return ""
	}
	if len(str) <= size {
		return str
	}
	return str[len(str)-size:]
}
