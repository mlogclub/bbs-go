package render

import (
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/markdown"
	"strings"
)

func BuildAboutPage(config dto.AboutPageConfig) map[string]any {
	content := getLocalizedText(config.Content)
	if strings.TrimSpace(content) == "" {
		return map[string]any{
			"content": "",
		}
	}
	return map[string]any{
		"content": handleHtmlContent(markdown.ToHTML(content)),
	}
}

func getLocalizedText(text dto.LocalizedText) string {
	if len(text) == 0 {
		return ""
	}
	if value := strings.TrimSpace(text[string(config.Instance.Language)]); value != "" {
		return value
	}
	if value := strings.TrimSpace(text[string(config.LanguageEnUS)]); value != "" {
		return value
	}
	if value := strings.TrimSpace(text[string(config.LanguageZhCN)]); value != "" {
		return value
	}
	for _, value := range text {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
