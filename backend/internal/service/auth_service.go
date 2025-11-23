// backend/internal/service/auth_service.go
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/repository"
	"github.com/your-org/esms/pkg/oidc"
)

var (
	ErrUnauthenticated        = errors.New("unauthenticated")
	ErrInsufficientPermission = errors.New("insufficient permission")
	ErrSessionNotFound        = errors.New("session not found")
	ErrInvalidState           = errors.New("invalid state")
)

// Session はユーザーセッション情報
type Session struct {
	UserID       uuid.UUID
	Email        string
	Name         string
	Role         domain.Role
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// AuthService は認証に関するビジネスロジックを提供します
type AuthService struct {
	oidcClient   *oidc.Client
	userRepo     repository.UserRepository
	auditLogRepo repository.AuditLogRepository
	// セッションストア（簡易実装: メモリ内）
	// 本番環境ではRedisなどの永続化ストアを使用
	sessions map[string]*Session
	// state検証用（簡易実装）
	states map[string]time.Time
}

// NewAuthService は新しいAuthServiceを作成します
func NewAuthService(
	oidcClient *oidc.Client,
	userRepo repository.UserRepository,
	auditLogRepo repository.AuditLogRepository,
) *AuthService {
	return &AuthService{
		oidcClient:   oidcClient,
		userRepo:     userRepo,
		auditLogRepo: auditLogRepo,
		sessions:     make(map[string]*Session),
		states:       make(map[string]time.Time),
	}
}

// GetAuthURL は認証URLを生成します
func (s *AuthService) GetAuthURL(state string) string {
	// state を保存（CSRF対策）
	s.states[state] = time.Now().Add(10 * time.Minute)
	return s.oidcClient.GetAuthURL(state)
}

// HandleCallback は認証コールバックを処理します
func (s *AuthService) HandleCallback(ctx context.Context, code, state string) (*Session, error) {
	// state検証（CSRF対策）
	if !s.validateState(state) {
		return nil, ErrInvalidState
	}

	// 認証コードをトークンに交換
	token, err := s.oidcClient.ExchangeCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// IDトークンを検証
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token in token response")
	}

	idToken, err := s.oidcClient.VerifyIDToken(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	// ユーザー情報を取得
	userInfo, err := s.oidcClient.GetUserInfo(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// ユーザーをDBに同期
	user, err := s.syncUser(ctx, userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to sync user: %w", err)
	}

	// セッションを作成
	sessionID := uuid.New().String()
	session := &Session{
		UserID:       user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Role:         user.Role,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}
	s.sessions[sessionID] = session

	// 監査ログ記録
	auditLog := &domain.AuditLog{
		ID:         uuid.New(),
		UserID:     user.ID,
		Action:     domain.AuditActionLogin,
		TargetType: "user",
		TargetID:   user.ID.String(),
		Details: map[string]interface{}{
			"email": user.Email,
		},
		CreatedAt: time.Now(),
	}
	_ = s.auditLogRepo.Create(ctx, auditLog)

	return session, nil
}

// syncUser はOIDCユーザー情報をDBに同期します
func (s *AuthService) syncUser(ctx context.Context, userInfo *oidc.UserInfo) (*domain.User, error) {
	// メールアドレスで既存ユーザーを検索
	user, err := s.userRepo.GetByEmail(ctx, userInfo.Email)
	if err == nil {
		// 既存ユーザーの情報を更新
		user.Name = userInfo.Name
		user.UpdatedAt = time.Now()
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, err
		}
		return user, nil
	}

	// 新規ユーザーを作成
	user = &domain.User{
		ID:        uuid.New(),
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		Role:      domain.RoleGeneral, // デフォルトロール
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// validateState はstateを検証します
func (s *AuthService) validateState(state string) bool {
	expiresAt, ok := s.states[state]
	if !ok {
		return false
	}

	// 有効期限チェック
	if time.Now().After(expiresAt) {
		delete(s.states, state)
		return false
	}

	// 使用済みstateを削除
	delete(s.states, state)
	return true
}

// GetSession はセッションIDからセッション情報を取得します
func (s *AuthService) GetSession(sessionID string) (*Session, error) {
	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}

	// セッション有効期限チェック
	if time.Now().After(session.ExpiresAt) {
		delete(s.sessions, sessionID)
		return nil, ErrSessionNotFound
	}

	return session, nil
}

// ValidateToken はアクセストークンを検証します
func (s *AuthService) ValidateToken(ctx context.Context, accessToken string) error {
	return s.oidcClient.ValidateToken(ctx, accessToken)
}

// RefreshSession はリフレッシュトークンを使用してセッションを更新します
func (s *AuthService) RefreshSession(ctx context.Context, sessionID string) (*Session, error) {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	// リフレッシュトークンで新しいアクセストークンを取得
	newToken, err := s.oidcClient.RefreshToken(ctx, session.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// セッション情報を更新
	session.AccessToken = newToken.AccessToken
	session.ExpiresAt = newToken.Expiry
	if newToken.RefreshToken != "" {
		session.RefreshToken = newToken.RefreshToken
	}

	return session, nil
}

// Logout はセッションを削除します
func (s *AuthService) Logout(ctx context.Context, sessionID string, userID uuid.UUID) error {
	delete(s.sessions, sessionID)

	// 監査ログ記録
	auditLog := &domain.AuditLog{
		ID:         uuid.New(),
		UserID:     userID,
		Action:     domain.AuditActionLogout,
		TargetType: "user",
		TargetID:   userID.String(),
		CreatedAt:  time.Now(),
	}
	_ = s.auditLogRepo.Create(ctx, auditLog)

	return nil
}

// CheckPermission は権限をチェックします
func (s *AuthService) CheckPermission(session *Session, requiredRole domain.Role) error {
	if session == nil {
		return ErrUnauthenticated
	}

	// ロールの階層: Admin > Manager > General
	roleLevel := map[domain.Role]int{
		domain.RoleGeneral: 1,
		domain.RoleManager: 2,
		domain.RoleAdmin:   3,
	}

	if roleLevel[session.Role] < roleLevel[requiredRole] {
		return ErrInsufficientPermission
	}

	return nil
}
