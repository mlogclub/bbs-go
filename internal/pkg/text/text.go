package text

import (
	"strings"

	"github.com/mlogclub/simple/common/strs"
)

// GetSummary 获取summary
func GetSummary(s string, length int) string {
	s = strings.TrimSpace(s)
	summary := strs.Substr(s, 0, length)
	if strs.RuneLen(s) > length {
		summary += "..."
	}
	return summary
}
