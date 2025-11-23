// backend/pkg/oidc/client.go
package oidc

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenExpired    = errors.New("token expired")
	ErrInvalidIssuer   = errors.New("invalid issuer")
	ErrInvalidAudience = errors.New("invalid audience")
	ErrInvalidConfig   = errors.New("invalid configuration")
	ErrInvalidNonce    = errors.New("invalid nonce")
	ErrInvalidAtHash   = errors.New("invalid at_hash")
	ErrFutureIssuedAt  = errors.New("token issued in the future")
)

// ClockSkew はトークン検証時の許容クロックスキュー
const ClockSkew = 1 * time.Minute

// TimeFunc は現在時刻を返す関数（テスト用に注入可能）
var TimeFunc = time.Now

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
	// Config検証
	if cfg == nil {
		return nil, ErrInvalidConfig
	}
	if cfg.IssuerURL == "" || cfg.ClientID == "" {
		return nil, fmt.Errorf("%w: issuer URL and client ID are required", ErrInvalidConfig)
	}

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

// ValidateToken はアクセストークンを検証します
// JWT形式のアクセストークンの場合、署名検証と有効期限チェックを行います
func (c *Client) ValidateToken(ctx context.Context, accessToken string) error {
	if accessToken == "" {
		return ErrInvalidToken
	}

	// JWT形式かチェック（3つのパートがドットで区切られている）
	parts := strings.Split(accessToken, ".")
	if len(parts) != 3 {
		// Opaque tokenの場合、イントロスペクションエンドポイントが必要
		// 現状は存在チェックのみ
		return nil
	}

	// JWTの場合、oidc.Verifierを使用して検証
	// アクセストークン用のVerifierを作成（audience検証なし）
	verifier := c.provider.Verifier(&oidc.Config{
		SkipClientIDCheck: true, // アクセストークンにはaudience要求なし
	})

	_, err := verifier.Verify(ctx, accessToken)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidToken, err)
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
// TODO: 本番環境では以下のクレームの検証も必要:
// - nonce: リプレイアタック防止（Authorization Code Flowで必須）
// - at_hash: アクセストークンのハッシュ検証（Implicit Flowで必須）
// - azp: Authorized party（複数audienceの場合に必須）
type TokenClaims struct {
	Issuer    string          `json:"iss"`
	Subject   string          `json:"sub"`
	Audience  audienceWrapper `json:"aud"` // 文字列または配列に対応
	ExpiresAt int64           `json:"exp"` // Unix timestamp
	IssuedAt  int64           `json:"iat"` // Unix timestamp
	Email     string          `json:"email"`
	Name      string          `json:"name"`
	// オプショナルなクレーム（将来の拡張用）
	Nonce  string `json:"nonce,omitempty"`   // リプレイアタック防止
	AtHash string `json:"at_hash,omitempty"` // アクセストークンハッシュ
	Azp    string `json:"azp,omitempty"`     // Authorized party
}

// audienceWrapper はaudienceクレームの柔軟な型対応
type audienceWrapper []string

func (a *audienceWrapper) UnmarshalJSON(data []byte) error {
	// 文字列の場合
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*a = []string{single}
		return nil
	}

	// 配列の場合
	var multiple []string
	if err := json.Unmarshal(data, &multiple); err == nil {
		*a = multiple
		return nil
	}

	return fmt.Errorf("audience must be string or array")
}

func (a audienceWrapper) Contains(audience string) bool {
	for _, aud := range a {
		if aud == audience {
			return true
		}
	}
	return false
}

// ParseIDTokenClaims はIDトークンのクレームをパースします
// nonceとaccessTokenを指定することで追加検証を行います
func (c *Client) ParseIDTokenClaims(ctx context.Context, rawIDToken string) (*TokenClaims, error) {
	return c.ParseIDTokenClaimsWithValidation(ctx, rawIDToken, "", "")
}

// ParseIDTokenClaimsWithValidation はIDトークンのクレームをパースし、nonce/at_hashを検証します
func (c *Client) ParseIDTokenClaimsWithValidation(ctx context.Context, rawIDToken, expectedNonce, accessToken string) (*TokenClaims, error) {
	if c.config == nil {
		return nil, ErrInvalidConfig
	}

	idToken, err := c.VerifyIDToken(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	var claims TokenClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	// 追加の検証（クロックスキュー許容）
	if err := c.validateClaims(&claims, expectedNonce, accessToken); err != nil {
		return nil, err
	}

	return &claims, nil
}

// validateClaims はクレームの追加検証を行います
func (c *Client) validateClaims(claims *TokenClaims, expectedNonce, accessToken string) error {
	if c.config == nil {
		return ErrInvalidConfig
	}

	// Issuer検証
	if claims.Issuer != c.config.IssuerURL {
		return fmt.Errorf("%w: expected %s, got %s", ErrInvalidIssuer, c.config.IssuerURL, claims.Issuer)
	}

	// Audience検証
	if !claims.Audience.Contains(c.config.ClientID) {
		return fmt.Errorf("%w: client ID %s not found in audience", ErrInvalidAudience, c.config.ClientID)
	}

	// 有効期限検証（クロックスキュー許容）
	now := TimeFunc()
	expTime := time.Unix(claims.ExpiresAt, 0)
	if now.After(expTime.Add(ClockSkew)) {
		return fmt.Errorf("%w: token expired at %v (now: %v)", ErrTokenExpired, expTime, now)
	}

	// iat（発行時刻）の未来値チェック（クロックスキュー許容）
	iatTime := time.Unix(claims.IssuedAt, 0)
	if now.Add(ClockSkew).Before(iatTime) {
		return fmt.Errorf("%w: token issued at %v (now: %v)", ErrFutureIssuedAt, iatTime, now)
	}

	// nonce検証
	if expectedNonce != "" {
		if claims.Nonce == "" {
			return fmt.Errorf("%w: nonce claim is missing", ErrInvalidNonce)
		}
		if claims.Nonce != expectedNonce {
			return fmt.Errorf("%w: expected %s, got %s", ErrInvalidNonce, expectedNonce, claims.Nonce)
		}
	}

	// at_hash検証（アクセストークンが提供された場合）
	if accessToken != "" && claims.AtHash != "" {
		if !verifyAtHash(claims.AtHash, accessToken) {
			return fmt.Errorf("%w: at_hash verification failed", ErrInvalidAtHash)
		}
	}

	return nil
}

// verifyAtHash はat_hashクレームを検証します
// OIDC仕様: at_hash = base64url(left_half(sha256(access_token)))
func verifyAtHash(atHash, accessToken string) bool {
	// アクセストークンのSHA256ハッシュを計算
	hash := sha256.Sum256([]byte(accessToken))
	// 左半分を取得
	leftHalf := hash[:len(hash)/2]
	// Base64 URL エンコード（パディングなし）
	expected := base64.RawURLEncoding.EncodeToString(leftHalf)

	return atHash == expected
}
