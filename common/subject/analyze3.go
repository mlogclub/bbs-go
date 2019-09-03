package subject

import "strings"

// Python
type subjectAnalyzer3 struct {
}

func (this *subjectAnalyzer3) GetSubjectId() int64 {
	return 3
}
func (this *subjectAnalyzer3) IsMatch(userId int64, title, content string) bool {
	title = strings.ToLower(title)
	if strings.Contains(title, "python") {
		return true
	}
	return false
}
