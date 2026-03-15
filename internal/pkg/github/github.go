package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const (
	AuthorizationCallbackURL = "/user/signin/callback/github" // AuthorizationCallbackURL GitHub 授权回调地址
)

var githubEndpoint = oauth2.Endpoint{
	AuthURL:  "https://github.com/login/oauth/authorize",
	TokenURL: "https://github.com/login/oauth/access_token",
}

type GithubOAuth struct {
	config *oauth2.Config
}

type GithubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

type githubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

func NewGithubOAuth(clientId, clientSecret, redirectURI string) *GithubOAuth {
	config := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes: []string{
			"read:user",
			"user:email",
		},
		Endpoint: githubEndpoint,
	}
	return &GithubOAuth{config: config}
}

func (g *GithubOAuth) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (g *GithubOAuth) GetUserInfo(ctx context.Context, code string) (*GithubUserInfo, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := g.config.Client(ctx, token)
	client.Timeout = 10 * time.Second

	user, err := g.getUser(client)
	if err != nil {
		return nil, err
	}
	// GitHub /user 可能拿不到 email（用户未公开），补查 /user/emails
	if user.Email == "" {
		if email, err := g.getPrimaryVerifiedEmail(client); err == nil && email != "" {
			user.Email = email
		}
	}
	return user, nil
}

func (g *GithubOAuth) getUser(client *http.Client) (*GithubUserInfo, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d, body: %s", resp.StatusCode, string(body))
	}

	var out GithubUserInfo
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}
	return &out, nil
}

func (g *GithubOAuth) getPrimaryVerifiedEmail(client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get user emails: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get user emails: status %d, body: %s", resp.StatusCode, string(body))
	}

	var emails []githubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", fmt.Errorf("failed to unmarshal emails: %w", err)
	}
	for _, e := range emails {
		if e.Primary && e.Verified && e.Email != "" {
			return e.Email, nil
		}
	}
	for _, e := range emails {
		if e.Verified && e.Email != "" {
			return e.Email, nil
		}
	}
	return "", nil
}
