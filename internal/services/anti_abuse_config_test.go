package services

import (
	"bbs-go/internal/models/dto"
	"testing"

	"github.com/mlogclub/simple/common/jsons"
)

func TestValidateAntiAbuseConfig(t *testing.T) {
	config := dto.DefaultAntiAbuseConfig()
	config.User.Enabled = true
	config.User.Action = dto.AntiAbuseActionReview
	config.IP.Enabled = true

	if err := validateAntiAbuseConfig(jsons.ToJsonStr(config)); err != nil {
		t.Fatalf("expected valid configuration, got %v", err)
	}

	config.IP.Comment.MaxCount = 0
	if err := validateAntiAbuseConfig(jsons.ToJsonStr(config)); err != nil {
		t.Fatalf("expected zero count to disable a rule, got %v", err)
	}

}
