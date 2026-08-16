package api

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/handlers/render"
	"bbs-go/internal/models/req"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/oidc"
	"bbs-go/internal/services"
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

const oidcCallbackPath = "/user/signin/callback/oidc"

func LoginOIDCLoginConfig(ctx *gin.Context) {
	providerKey := strings.TrimSpace(ctx.Query("provider"))
	provider, ok := services.OIDCProvider(providerKey)
	if !ok || !provider.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("OIDC provider is unavailable"))
		return
	}
	state, err := oidc.RandomURLValue()
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	nonce, err := oidc.RandomURLValue()
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	verifier, err := oidc.RandomURLValue()
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	redirect := strings.TrimSpace(ctx.Query("redirect"))
	if redirect == "" || !strings.HasPrefix(redirect, "/") || strings.HasPrefix(redirect, "//") {
		redirect = "/"
	}
	cache.OIDCLoginStateCache.Put(state, &cache.OIDCLoginStateData{ProviderKey: providerKey, Redirect: redirect, Nonce: nonce, Verifier: verifier})
	authURL, err := oidc.AuthorizationURL(context.Background(), provider, bbsurls.AbsUrl(oidcCallbackPath), state, nonce, verifier)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, map[string]string{"authUrl": authURL})
}

func LoginOIDCLoginSubmit(ctx *gin.Context) {
	var request req.OAuthCodeStateReq
	if err := ginx.Bind(ctx, &request); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	if strings.TrimSpace(request.Code) == "" || strings.TrimSpace(request.State) == "" {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("OIDC code and state are required"))
		return
	}
	state := cache.OIDCLoginStateCache.Take(request.State)
	if state == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("OIDC login request expired or was already used"))
		return
	}
	provider, ok := services.OIDCProvider(state.ProviderKey)
	if !ok || !provider.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("OIDC provider is unavailable"))
		return
	}
	claims, err := oidc.ExchangeAndVerify(context.Background(), provider, bbsurls.AbsUrl(oidcCallbackPath), request.Code, state.Verifier, state.Nonce)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	user, err := services.OIDCLogin(provider, claims)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, render.BuildLoginSuccess(ctx, user, state.Redirect))
}
