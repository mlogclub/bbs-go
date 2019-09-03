package subject

import "strings"

// 程序员
type subjectAnalyzer5 struct {
}

func (this *subjectAnalyzer5) GetSubjectId() int64 {
	return 5
}
func (this *subjectAnalyzer5) IsMatch(userId int64, title, content string) bool {
	title = strings.ToLower(title)
	if strings.Contains(title, "程序员") || strings.Contains(title, "码农") {
		return true
	}
	return false
}
