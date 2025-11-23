// backend/pkg/oidc/client_test.go
package oidc_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/esms/pkg/oidc"
)

func TestNewClient(t *testing.T) {
	// 注意: このテストは実際のOIDCプロバイダーが必要
	// 統合テスト環境でのみ実行されるべき
	t.Skip("Requires actual OIDC provider")

	cfg := &oidc.Config{
		IssuerURL:    "https://accounts.google.com",
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"openid", "email", "profile"},
	}

	ctx := context.Background()
	client, err := oidc.NewClient(ctx, cfg)

	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestClient_GetAuthURL(t *testing.T) {
	t.Skip("Requires actual OIDC provider")

	cfg := &oidc.Config{
		IssuerURL:    "https://accounts.google.com",
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"openid", "email", "profile"},
	}

	ctx := context.Background()
	client, err := oidc.NewClient(ctx, cfg)
	assert.NoError(t, err)

	state := "random-state-string"
	authURL := client.GetAuthURL(state)

	assert.NotEmpty(t, authURL)
	assert.Contains(t, authURL, "client_id=test-client-id")
	assert.Contains(t, authURL, "state=random-state-string")
	assert.Contains(t, authURL, "redirect_uri=")
}

func TestClient_ValidateToken(t *testing.T) {
	t.Skip("Requires actual OIDC provider")

	cfg := &oidc.Config{
		IssuerURL:    "https://accounts.google.com",
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"openid", "email", "profile"},
	}

	ctx := context.Background()
	client, err := oidc.NewClient(ctx, cfg)
	assert.NoError(t, err)

	// 空のトークンは無効
	err = client.ValidateToken(ctx, "")
	assert.Error(t, err)
	assert.Equal(t, oidc.ErrInvalidToken, err)

	// 有効なトークン（簡易チェック）
	err = client.ValidateToken(ctx, "valid-token")
	assert.NoError(t, err)
}

func TestTokenClaims_Validation(t *testing.T) {
	// ユニットテスト: クレーム構造のテスト
	now := time.Now()
	claims := &oidc.TokenClaims{
		Issuer:    "https://accounts.google.com",
		Subject:   "user-123",
		Audience:  []string{"test-client-id"},
		ExpiresAt: now.Add(1 * time.Hour).Unix(),
		IssuedAt:  now.Unix(),
		Email:     "test@example.com",
		Name:      "Test User",
	}

	assert.Equal(t, "https://accounts.google.com", claims.Issuer)
	assert.Equal(t, "user-123", claims.Subject)
	assert.Equal(t, "test@example.com", claims.Email)

	// 有効期限検証
	expTime := time.Unix(claims.ExpiresAt, 0)
	assert.True(t, now.Before(expTime))
}

func TestUserInfo_Structure(t *testing.T) {
	// ユニットテスト: UserInfo構造のテスト
	userInfo := &oidc.UserInfo{
		Subject:       "user-123",
		Email:         "test@example.com",
		EmailVerified: true,
		Name:          "Test User",
		Picture:       "https://example.com/avatar.jpg",
	}

	assert.Equal(t, "user-123", userInfo.Subject)
	assert.Equal(t, "test@example.com", userInfo.Email)
	assert.True(t, userInfo.EmailVerified)
	assert.Equal(t, "Test User", userInfo.Name)
	assert.NotEmpty(t, userInfo.Picture)
}

// 統合テスト用のヘルパー関数
func setupTestOIDCClient(t *testing.T) *oidc.Client {
	t.Helper()

	cfg := &oidc.Config{
		IssuerURL:    "https://accounts.google.com",
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"openid", "email", "profile"},
	}

	ctx := context.Background()
	client, err := oidc.NewClient(ctx, cfg)
	if err != nil {
		t.Skip("OIDC provider not available")
	}

	return client
}
