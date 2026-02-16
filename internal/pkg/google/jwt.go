package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// VerifyJWTWithGoogleAPI 使用 Google API 验证 JWT（推荐方式）
// 这个方法会调用 Google 的 tokeninfo 端点进行验证
func VerifyJWTWithGoogleAPI(ctx context.Context, idToken string) (*GoogleUserInfo, error) {
	// 调用 Google 的 tokeninfo 端点验证 JWT
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to verify JWT: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("JWT verification failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tokenInfo struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified string `json:"email_verified"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Aud           string `json:"aud"`
		Iss           string `json:"iss"`
		Exp           string `json:"exp"`
	}

	if err := json.Unmarshal(body, &tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token info: %w", err)
	}

	// 验证签发者
	if tokenInfo.Iss != "https://accounts.google.com" && tokenInfo.Iss != "accounts.google.com" {
		return nil, fmt.Errorf("invalid JWT issuer: %s", tokenInfo.Iss)
	}

	// 转换为 GoogleUserInfo 格式
	emailVerified := false
	if tokenInfo.EmailVerified == "true" {
		emailVerified = true
	}

	userInfo := &GoogleUserInfo{
		ID:            tokenInfo.Sub, // sub 字段等同于 UserInfo API 的 id
		Email:         tokenInfo.Email,
		VerifiedEmail: emailVerified,
		Name:          tokenInfo.Name,
		GivenName:     tokenInfo.GivenName,
		FamilyName:    tokenInfo.FamilyName,
		Picture:       tokenInfo.Picture,
	}

	return userInfo, nil
}
