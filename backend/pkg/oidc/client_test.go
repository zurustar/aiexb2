// backend/pkg/oidc/client_test.go
package oidc_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/your-org/esms/pkg/oidc"
)

// mockOIDCServer はテスト用のモックOIDCプロバイダ
type mockOIDCServer struct {
	server     *httptest.Server
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
}

// newMockOIDCServer はモックOIDCサーバーを作成します
func newMockOIDCServer(t *testing.T) *mockOIDCServer {
	t.Helper()

	// RSAキーペアを生成
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	mock := &mockOIDCServer{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}

	// HTTPサーバーを作成
	mux := http.NewServeMux()

	// Discovery endpoint
	mux.HandleFunc("/.well-known/openid-configuration", mock.handleDiscovery)

	// JWKS endpoint
	mux.HandleFunc("/jwks", mock.handleJWKS)

	mock.server = httptest.NewServer(mux)
	mock.issuer = mock.server.URL

	return mock
}

// Close はモックサーバーを閉じます
func (m *mockOIDCServer) Close() {
	m.server.Close()
}

// handleDiscovery はDiscoveryエンドポイントを処理します
func (m *mockOIDCServer) handleDiscovery(w http.ResponseWriter, r *http.Request) {
	discovery := map[string]interface{}{
		"issuer":                                m.issuer,
		"authorization_endpoint":                m.issuer + "/authorize",
		"token_endpoint":                        m.issuer + "/token",
		"jwks_uri":                              m.issuer + "/jwks",
		"response_types_supported":              []string{"code"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discovery)
}

// handleJWKS はJWKSエンドポイントを処理します
func (m *mockOIDCServer) handleJWKS(w http.ResponseWriter, r *http.Request) {
	// 公開鍵をJWK形式で返す
	n := base64.RawURLEncoding.EncodeToString(m.publicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString([]byte{1, 0, 1}) // 65537

	jwks := map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"kty": "RSA",
				"use": "sig",
				"kid": "test-key-id",
				"n":   n,
				"e":   e,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)
}

// generateIDToken はテスト用のIDトークンを生成します
func (m *mockOIDCServer) generateIDToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-key-id"

	return token.SignedString(m.privateKey)
}

// TestGeneratePKCEChallenge はPKCE生成をテストします
func TestGeneratePKCEChallenge(t *testing.T) {
	pkce, err := oidc.GeneratePKCEChallenge()

	require.NoError(t, err)
	require.NotNil(t, pkce)
	assert.NotEmpty(t, pkce.CodeVerifier)
	assert.NotEmpty(t, pkce.CodeChallenge)

	// code_challengeがcode_verifierのSHA256ハッシュであることを確認
	hash := sha256.Sum256([]byte(pkce.CodeVerifier))
	expectedChallenge := base64.RawURLEncoding.EncodeToString(hash[:])
	assert.Equal(t, expectedChallenge, pkce.CodeChallenge)
}

// TestGenerateNonce はnonce生成をテストします
func TestGenerateNonce(t *testing.T) {
	nonce, err := oidc.GenerateNonce()

	require.NoError(t, err)
	assert.NotEmpty(t, nonce)

	// 2つのnonceは異なるべき
	nonce2, err := oidc.GenerateNonce()
	require.NoError(t, err)
	assert.NotEqual(t, nonce, nonce2)
}

// TestNewClient_WithMockServer はモックサーバーでのクライアント作成をテストします
func TestNewClient_WithMockServer(t *testing.T) {
	mock := newMockOIDCServer(t)
	defer mock.Close()

	cfg := &oidc.Config{
		IssuerURL:    mock.issuer,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"openid", "email", "profile"},
	}

	ctx := context.Background()
	client, err := oidc.NewClient(ctx, cfg)

	require.NoError(t, err)
	require.NotNil(t, client)
}

// TestNewClient_InvalidConfig は無効な設定でのエラーをテストします
func TestNewClient_InvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *oidc.Config
	}{
		{
			name:   "Nil config",
			config: nil,
		},
		{
			name: "Empty issuer",
			config: &oidc.Config{
				ClientID: "test-client-id",
			},
		},
		{
			name: "Empty client ID",
			config: &oidc.Config{
				IssuerURL: "https://example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			client, err := oidc.NewClient(ctx, tt.config)

			assert.Error(t, err)
			assert.Nil(t, client)
			assert.ErrorIs(t, err, oidc.ErrInvalidConfig)
		})
	}
}

// TestGetAuthURL_WithPKCEAndNonce はPKCEとnonceを含む認証URLをテストします
func TestGetAuthURL_WithPKCEAndNonce(t *testing.T) {
	mock := newMockOIDCServer(t)
	defer mock.Close()

	cfg := &oidc.Config{
		IssuerURL:    mock.issuer,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"openid", "email"},
	}

	ctx := context.Background()
	client, err := oidc.NewClient(ctx, cfg)
	require.NoError(t, err)

	pkce, err := oidc.GeneratePKCEChallenge()
	require.NoError(t, err)

	nonce, err := oidc.GenerateNonce()
	require.NoError(t, err)

	authURL := client.GetAuthURL(oidc.AuthURLParams{
		State:         "test-state",
		Nonce:         nonce,
		CodeChallenge: pkce.CodeChallenge,
	})

	assert.NotEmpty(t, authURL)
	assert.Contains(t, authURL, "client_id=test-client-id")
	assert.Contains(t, authURL, "state=test-state")
	assert.Contains(t, authURL, "nonce="+nonce)
	assert.Contains(t, authURL, "code_challenge="+pkce.CodeChallenge)
	assert.Contains(t, authURL, "code_challenge_method=S256")
}

// TestVerifyIDToken_ValidToken は有効なIDトークンの検証をテストします
func TestVerifyIDToken_ValidToken(t *testing.T) {
	mock := newMockOIDCServer(t)
	defer mock.Close()

	cfg := &oidc.Config{
		IssuerURL:    mock.issuer,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
	}

	ctx := context.Background()
	client, err := oidc.NewClient(ctx, cfg)
	require.NoError(t, err)

	// 有効なIDトークンを生成
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   mock.issuer,
		"sub":   "user-123",
		"aud":   "test-client-id",
		"exp":   now.Add(1 * time.Hour).Unix(),
		"iat":   now.Unix(),
		"email": "test@example.com",
		"name":  "Test User",
	}

	idToken, err := mock.generateIDToken(claims)
	require.NoError(t, err)

	// IDトークンを検証
	verified, err := client.VerifyIDToken(ctx, idToken)
	require.NoError(t, err)
	require.NotNil(t, verified)
}

// TestVerifyIDToken_WithNonce はnonceを含むIDトークンの検証をテストします
func TestVerifyIDToken_WithNonce(t *testing.T) {
	mock := newMockOIDCServer(t)
	defer mock.Close()

	cfg := &oidc.Config{
		IssuerURL:    mock.issuer,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	ctx := context.Background()
	client, err := oidc.NewClient(ctx, cfg)
	require.NoError(t, err)

	nonce := "test-nonce-12345"
	now := time.Now()

	claims := jwt.MapClaims{
		"iss":   mock.issuer,
		"sub":   "user-123",
		"aud":   "test-client-id",
		"exp":   now.Add(1 * time.Hour).Unix(),
		"iat":   now.Unix(),
		"nonce": nonce,
	}

	idToken, err := mock.generateIDToken(claims)
	require.NoError(t, err)

	// nonceを指定して検証
	parsedClaims, err := client.ParseIDTokenClaimsWithValidation(ctx, idToken, nonce, "")
	require.NoError(t, err)
	assert.Equal(t, nonce, parsedClaims.Nonce)
}

// TestVerifyIDToken_InvalidNonce は無効なnonceでのエラーをテストします
func TestVerifyIDToken_InvalidNonce(t *testing.T) {
	mock := newMockOIDCServer(t)
	defer mock.Close()

	cfg := &oidc.Config{
		IssuerURL:    mock.issuer,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	ctx := context.Background()
	client, err := oidc.NewClient(ctx, cfg)
	require.NoError(t, err)

	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   mock.issuer,
		"sub":   "user-123",
		"aud":   "test-client-id",
		"exp":   now.Add(1 * time.Hour).Unix(),
		"iat":   now.Unix(),
		"nonce": "wrong-nonce",
	}

	idToken, err := mock.generateIDToken(claims)
	require.NoError(t, err)

	// 異なるnonceで検証
	_, err = client.ParseIDTokenClaimsWithValidation(ctx, idToken, "expected-nonce", "")
	assert.Error(t, err)
	assert.ErrorIs(t, err, oidc.ErrInvalidNonce)
}

// TestVerifyIDToken_ExpiredToken は期限切れトークンのエラーをテストします
func TestVerifyIDToken_ExpiredToken(t *testing.T) {
	mock := newMockOIDCServer(t)
	defer mock.Close()

	cfg := &oidc.Config{
		IssuerURL:    mock.issuer,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	ctx := context.Background()
	client, err := oidc.NewClient(ctx, cfg)
	require.NoError(t, err)

	// 期限切れのIDトークンを生成
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": mock.issuer,
		"sub": "user-123",
		"aud": "test-client-id",
		"exp": now.Add(-1 * time.Hour).Unix(), // 1時間前に期限切れ
		"iat": now.Add(-2 * time.Hour).Unix(),
	}

	idToken, err := mock.generateIDToken(claims)
	require.NoError(t, err)

	// 検証はOIDCライブラリ内部で失敗する
	_, err = client.VerifyIDToken(ctx, idToken)
	assert.Error(t, err)
}

// TestVerifyAtHash はat_hash検証をテストします
func TestVerifyAtHash(t *testing.T) {
	mock := newMockOIDCServer(t)
	defer mock.Close()

	cfg := &oidc.Config{
		IssuerURL:    mock.issuer,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	ctx := context.Background()
	client, err := oidc.NewClient(ctx, cfg)
	require.NoError(t, err)

	accessToken := "test-access-token-12345"

	// at_hashを計算
	hash := sha256.Sum256([]byte(accessToken))
	leftHalf := hash[:len(hash)/2]
	atHash := base64.RawURLEncoding.EncodeToString(leftHalf)

	now := time.Now()
	claims := jwt.MapClaims{
		"iss":     mock.issuer,
		"sub":     "user-123",
		"aud":     "test-client-id",
		"exp":     now.Add(1 * time.Hour).Unix(),
		"iat":     now.Unix(),
		"at_hash": atHash,
	}

	idToken, err := mock.generateIDToken(claims)
	require.NoError(t, err)

	// at_hashを指定して検証
	parsedClaims, err := client.ParseIDTokenClaimsWithValidation(ctx, idToken, "", accessToken)
	require.NoError(t, err)
	assert.Equal(t, atHash, parsedClaims.AtHash)
}

// TestVerifyAtHash_Invalid は無効なat_hashでのエラーをテストします
func TestVerifyAtHash_Invalid(t *testing.T) {
	mock := newMockOIDCServer(t)
	defer mock.Close()

	cfg := &oidc.Config{
		IssuerURL:    mock.issuer,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	ctx := context.Background()
	client, err := oidc.NewClient(ctx, cfg)
	require.NoError(t, err)

	now := time.Now()
	claims := jwt.MapClaims{
		"iss":     mock.issuer,
		"sub":     "user-123",
		"aud":     "test-client-id",
		"exp":     now.Add(1 * time.Hour).Unix(),
		"iat":     now.Unix(),
		"at_hash": "invalid-at-hash",
	}

	idToken, err := mock.generateIDToken(claims)
	require.NoError(t, err)

	// 異なるアクセストークンで検証
	_, err = client.ParseIDTokenClaimsWithValidation(ctx, idToken, "", "different-access-token")
	assert.Error(t, err)
	assert.ErrorIs(t, err, oidc.ErrInvalidAtHash)
}
