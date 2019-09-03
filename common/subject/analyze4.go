package subject

import "strings"

// Android
type subjectAnalyzer4 struct {
}

func (this *subjectAnalyzer4) GetSubjectId() int64 {
	return 4
}
func (this *subjectAnalyzer4) IsMatch(userId int64, title, content string) bool {
	title = strings.ToLower(title)
	if strings.Contains(title, "android") || strings.Contains(title, "安卓") {
		return true
	}
	return false
}
