package subject

import "strings"

// Java
type subjectAnalyzer2 struct {
}

func (this *subjectAnalyzer2) GetSubjectId() int64 {
	return 2
}
func (this *subjectAnalyzer2) IsMatch(userId int64, title, content string) bool {
	title = strings.ToLower(title)
	if strings.Contains(title, "java") || strings.Contains(title, "spring") {
		return true
	}
	return false
}
