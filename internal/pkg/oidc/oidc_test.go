package oidc

import (
	"bbs-go/internal/models/dto"
	"testing"
)

func TestValidateProvider(t *testing.T) {
	valid := dto.OIDCProviderConfig{Key: "company_sso", Name: "Company SSO", Issuer: "https://sso.example.com/realms/main", ClientId: "bbs-go"}
	if err := ValidateProvider(valid); err != nil {
		t.Fatalf("expected valid provider: %v", err)
	}
	valid.Issuer = "http://sso.example.com"
	if err := ValidateProvider(valid); err == nil {
		t.Fatal("HTTP issuer must be rejected")
	}
	valid.Issuer = "https://sso.example.com"
	valid.Key = "Company SSO"
	if err := ValidateProvider(valid); err == nil {
		t.Fatal("invalid provider key must be rejected")
	}
}

func TestChallengeIsDeterministic(t *testing.T) {
	if got, want := Challenge("dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"), "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"; got != want {
		t.Fatalf("challenge=%q want=%q", got, want)
	}
}
