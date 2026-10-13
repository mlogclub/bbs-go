package google

import "testing"

func TestNewGoogleOAuth_TrimSpace(t *testing.T) {
	client := NewGoogleOAuth("  my_google_client_id  ", "  my_google_secret \n", " http://example.com/google_callback ")
	if client.config.ClientID != "my_google_client_id" {
		t.Fatalf("expected client ID to be trimmed, got %q", client.config.ClientID)
	}
	if client.config.ClientSecret != "my_google_secret" {
		t.Fatalf("expected client secret to be trimmed, got %q", client.config.ClientSecret)
	}
	if client.config.RedirectURL != "http://example.com/google_callback" {
		t.Fatalf("expected redirect URL to be trimmed, got %q", client.config.RedirectURL)
	}
}
