package github

import "testing"

func TestNewGithubOAuth_TrimSpace(t *testing.T) {
	client := NewGithubOAuth("  my_client_id  ", "  my_secret \n", " http://example.com/callback ")
	if client.config.ClientID != "my_client_id" {
		t.Fatalf("expected client ID to be trimmed, got %q", client.config.ClientID)
	}
	if client.config.ClientSecret != "my_secret" {
		t.Fatalf("expected client secret to be trimmed, got %q", client.config.ClientSecret)
	}
	if client.config.RedirectURL != "http://example.com/callback" {
		t.Fatalf("expected redirect URL to be trimmed, got %q", client.config.RedirectURL)
	}
}
