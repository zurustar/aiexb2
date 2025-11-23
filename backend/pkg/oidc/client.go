// backend/pkg/oidc/client.go
package oidc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenExpired    = errors.New("token expired")
	ErrInvalidIssuer   = errors.New("invalid issuer")
	ErrInvalidAudience = errors.New("invalid audience")
)

// Config はOIDCクライアントの設定
type Config struct {
	IssuerURL    string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// Client はOIDCクライアント
type Client struct {
	config       *Config
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	oauth2Config *oauth2.Config
}

// NewClient は新しいOIDCクライアントを作成します
func NewClient(ctx context.Context, cfg *Config) (*Client, error) {
	// OIDC Discovery
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	// IDトークン検証器を作成
	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.ClientID,
	})

	// OAuth2設定
	oauth2Config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       cfg.Scopes,
	}

	return &Client{
		config:       cfg,
		provider:     provider,
		verifier:     verifier,
		oauth2Config: oauth2Config,
	}, nil
}

// GetAuthURL は認証URLを生成します
func (c *Client) GetAuthURL(state string) string {
	return c.oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCode は認証コードをトークンに交換します
func (c *Client) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := c.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return token, nil
}

// VerifyIDToken はIDトークンを検証します
func (c *Client) VerifyIDToken(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	idToken, err := c.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	return idToken, nil
}

// UserInfo はユーザー情報を表します
type UserInfo struct {
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// GetUserInfo はIDトークンからユーザー情報を取得します
func (c *Client) GetUserInfo(ctx context.Context, idToken *oidc.IDToken) (*UserInfo, error) {
	var userInfo UserInfo
	if err := idToken.Claims(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &userInfo, nil
}

// ValidateToken はアクセストークンを検証します（簡易実装）
func (c *Client) ValidateToken(ctx context.Context, accessToken string) error {
	// 実際の実装では、トークンイントロスペクションエンドポイントを使用するか、
	// JWTの場合は署名検証を行う
	// ここでは簡易的にトークンの存在チェックのみ
	if accessToken == "" {
		return ErrInvalidToken
	}

	return nil
}

// RefreshToken はリフレッシュトークンを使用して新しいアクセストークンを取得します
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	tokenSource := c.oauth2Config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: refreshToken,
	})

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return newToken, nil
}

// TokenClaims はトークンのクレームを表します
type TokenClaims struct {
	Issuer    string    `json:"iss"`
	Subject   string    `json:"sub"`
	Audience  []string  `json:"aud"`
	ExpiresAt time.Time `json:"exp"`
	IssuedAt  time.Time `json:"iat"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
}

// ParseIDTokenClaims はIDトークンのクレームをパースします
func (c *Client) ParseIDTokenClaims(ctx context.Context, rawIDToken string) (*TokenClaims, error) {
	idToken, err := c.VerifyIDToken(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	var claims TokenClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	// 追加の検証
	if claims.Issuer != c.config.IssuerURL {
		return nil, ErrInvalidIssuer
	}

	if time.Now().After(claims.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	return &claims, nil
}
