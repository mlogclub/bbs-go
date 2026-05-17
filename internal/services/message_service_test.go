package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/pkg/msg"
	"testing"
)

func TestBuildEmailNoticeSubjectAvoidsBlankSitePrefix(t *testing.T) {
	got := MessageService.buildEmailNoticeSubject("", "你的话题被设为推荐")
	if got != "你的话题被设为推荐" {
		t.Fatalf("expected subject without blank site prefix, got %q", got)
	}
}

func TestBuildEmailNoticeSubjectIncludesSiteTitle(t *testing.T) {
	got := MessageService.buildEmailNoticeSubject("BBS-GO", "你的话题被设为推荐")
	if got != "BBS-GO - 你的话题被设为推荐" {
		t.Fatalf("expected subject with site title, got %q", got)
	}
}

func TestBuildEmailNoticeContentFallsBackToNoticeTitle(t *testing.T) {
	got := MessageService.buildEmailNoticeContent("", "你的话题被设为推荐")
	if got != "你的话题被设为推荐" {
		t.Fatalf("expected notice title fallback, got %q", got)
	}
}

func TestBuildEmailNoticeDetailURLUsesTopicForRecommend(t *testing.T) {
	idcodec.Init(1)

	got := MessageService.buildEmailNoticeDetailURL(&models.Message{
		Type:      int(msg.TypeTopicRecommend),
		ExtraData: `{"topicId":123}`,
	})

	if got != "/topic/"+idcodec.Encode(123) {
		t.Fatalf("expected topic detail url, got %q", got)
	}
}
