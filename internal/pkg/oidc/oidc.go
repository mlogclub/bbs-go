package oidc

import (
	"bbs-go/internal/models/dto"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	coreoidc "github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Claims struct {
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	PreferredName string `json:"preferred_username"`
	Picture       string `json:"picture"`
}

func RandomURLValue() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func Challenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func ValidateProvider(cfg dto.OIDCProviderConfig) error {
	if strings.TrimSpace(cfg.Key) == "" || strings.TrimSpace(cfg.Name) == "" || strings.TrimSpace(cfg.ClientId) == "" {
		return fmt.Errorf("OIDC provider key, name and client ID are required")
	}
	if len(cfg.Key) > 64 || !validKey(cfg.Key) {
		return fmt.Errorf("OIDC provider key may contain only lowercase letters, digits, hyphens and underscores")
	}
	u, err := url.Parse(cfg.Issuer)
	if err != nil || u.Scheme != "https" || u.Host == "" {
		return fmt.Errorf("OIDC issuer must be an HTTPS URL")
	}
	return nil
}

func validKey(key string) bool {
	for _, r := range key {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' || r == '_') {
			return false
		}
	}
	return true
}

func provider(ctx context.Context, issuer string) (*coreoidc.Provider, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	return coreoidc.NewProvider(context.WithValue(ctx, oauth2.HTTPClient, client), issuer)
}

func AuthorizationURL(ctx context.Context, cfg dto.OIDCProviderConfig, redirectURI, state, nonce, verifier string) (string, error) {
	if err := ValidateProvider(cfg); err != nil {
		return "", err
	}
	p, err := provider(ctx, cfg.Issuer)
	if err != nil {
		return "", fmt.Errorf("OIDC discovery failed: %w", err)
	}
	scopes := cfg.Scopes
	if len(scopes) == 0 {
		scopes = []string{coreoidc.ScopeOpenID, "profile", "email"}
	}
	hasOpenID := false
	for _, scope := range scopes {
		if scope == coreoidc.ScopeOpenID {
			hasOpenID = true
		}
	}
	if !hasOpenID {
		scopes = append([]string{coreoidc.ScopeOpenID}, scopes...)
	}
	c := oauth2.Config{ClientID: cfg.ClientId, ClientSecret: cfg.ClientSecret, Endpoint: p.Endpoint(), RedirectURL: redirectURI, Scopes: scopes}
	return c.AuthCodeURL(state, oauth2.SetAuthURLParam("nonce", nonce), oauth2.SetAuthURLParam("code_challenge", Challenge(verifier)), oauth2.SetAuthURLParam("code_challenge_method", "S256")), nil
}

func ExchangeAndVerify(ctx context.Context, cfg dto.OIDCProviderConfig, redirectURI, code, verifier, nonce string) (*Claims, error) {
	if err := ValidateProvider(cfg); err != nil {
		return nil, err
	}
	p, err := provider(ctx, cfg.Issuer)
	if err != nil {
		return nil, fmt.Errorf("OIDC discovery failed: %w", err)
	}
	c := oauth2.Config{ClientID: cfg.ClientId, ClientSecret: cfg.ClientSecret, Endpoint: p.Endpoint(), RedirectURL: redirectURI}
	token, err := c.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", verifier))
	if err != nil {
		return nil, fmt.Errorf("OIDC token exchange failed: %w", err)
	}
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		return nil, fmt.Errorf("OIDC token response did not contain an ID token")
	}
	idToken, err := p.Verifier(&coreoidc.Config{ClientID: cfg.ClientId}).Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("OIDC ID token verification failed: %w", err)
	}
	var claims Claims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("OIDC ID token claims failed: %w", err)
	}
	if claims.Subject == "" || idToken.Nonce != nonce {
		return nil, fmt.Errorf("OIDC ID token nonce or subject is invalid")
	}
	if cfg.RequireVerifiedEmail && (!claims.EmailVerified || claims.Email == "") {
		return nil, fmt.Errorf("OIDC provider did not return a verified email address")
	}
	return &claims, nil
}
