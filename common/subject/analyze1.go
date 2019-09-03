package subject

import "strings"

// Go语言
type subjectAnalyzer1 struct {
}

func (this *subjectAnalyzer1) GetSubjectId() int64 {
	return 1
}
func (this *subjectAnalyzer1) IsMatch(userId int64, title, content string) bool {
	title = strings.ToLower(title)
	if strings.Contains(title, "go") || strings.Contains(title, "go语言") || strings.Contains(title, "golang") {
		return true
	}
	return false
}
