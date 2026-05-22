package services

import "testing"

func TestParseModulesConfig_BackfillsQaFromTopicForLegacyConfig(t *testing.T) {
	cfg := parseModulesConfig(`{"tweet":true,"topic":true,"article":false}`)

	if !cfg.QA {
		t.Fatalf("expected legacy config without qa to keep QA enabled when topic is enabled")
	}
}

func TestParseModulesConfig_RespectsExplicitQaSwitch(t *testing.T) {
	cfg := parseModulesConfig(`{"tweet":true,"topic":true,"qa":false,"article":true}`)

	if cfg.QA {
		t.Fatalf("expected explicit qa=false to disable QA independently from topic")
	}
}
