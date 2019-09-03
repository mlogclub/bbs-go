package subject

var analyzers = []subjectAnalyzer{
	&subjectAnalyzer1{}, &subjectAnalyzer2{}, &subjectAnalyzer3{}, &subjectAnalyzer4{}, &subjectAnalyzer5{},
}

// 分析内容属于哪个专栏
func AnalyzeSubjects(userId int64, title, content string) (subjectIds []int64) {
	for _, analyzer := range analyzers {
		if analyzer.IsMatch(userId, title, content) {
			subjectIds = append(subjectIds, analyzer.GetSubjectId())
		}
	}
	return
}

type subjectAnalyzer interface {
	GetSubjectId() int64
	IsMatch(userId int64, title, content string) bool
}
