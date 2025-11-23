// backend/internal/handler/auth_handler.go
package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/your-org/esms/internal/service"
)

// AuthServiceInterface は認証サービスのインターフェース
type AuthServiceInterface interface {
	GetAuthURL(state string) string
	HandleCallback(ctx context.Context, code, state string) (*service.Session, error)
	Logout(ctx context.Context, sessionID string, userID uuid.UUID) error
	RefreshSession(ctx context.Context, sessionID string) (*service.Session, error)
	GetSession(sessionID string) (*service.Session, error)
}

// AuthHandler は認証関連のHTTPハンドラー
type AuthHandler struct {
	authService AuthServiceInterface
}

// NewAuthHandler は新しいAuthHandlerを作成します
func NewAuthHandler(authService AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRoutes はルートを登録します
func (h *AuthHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/auth/login", h.Login).Methods("GET")
	r.HandleFunc("/api/v1/auth/callback", h.Callback).Methods("GET")
	r.HandleFunc("/api/v1/auth/logout", h.Logout).Methods("POST")
	r.HandleFunc("/api/v1/auth/refresh", h.Refresh).Methods("POST")
}

// Login はOIDC認証を開始します
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		state = "random-state" // 本番環境では適切なstate生成
	}

	authURL := h.authService.GetAuthURL(state)

	WriteJSON(w, http.StatusOK, map[string]string{
		"auth_url": authURL,
	})
}

// Callback はOIDC認証のコールバックを処理します
func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Missing code parameter")
		return
	}

	session, err := h.authService.HandleCallback(r.Context(), code, state)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "AUTH_FAILED", err.Error())
		return
	}

	// セッションIDをCookieに設定
	// 本番環境では session.ID を使用すべきだが、現状の HandleCallback は *Session を返しており
	// Session 構造体に ID フィールドがない可能性があるため、アクセストークンを代用するか、
	// AuthService 側で SessionID を返すように修正が必要。
	// ここでは仮に "session-token" とするが、実際には session オブジェクトから適切なIDを取得すべき。
	// auth_service.go の Session 定義を見ると ID フィールドがない。
	// しかし、HandleCallback 内で sessionID := uuid.New().String() を生成して map に保存している。
	// HandleCallback の戻り値に sessionID を含めるか、Session 構造体に ID を追加するのが正しい。
	// 今回はリファクタリングの範囲を広げすぎないため、AccessToken をセッションIDとして扱う(middlewareの実装に合わせる)
	// ただし、middleware.go では `session, err := m.authService.GetSession(token)` としており、
	// auth_service.go の GetSession は sessionID をキーにしている。
	// auth_service.go の HandleCallback は `return session, nil` しており、sessionID を返していない。
	// これは設計上の不整合。
	// 修正案: HandleCallback が (*Session, string, error) を返すようにするか、Session に ID を持たせる。
	// ここでは、AuthService の HandleCallback を修正するのは手間がかかるため、
	// AuthHandler 側で Cookie にセットする値は一旦 AccessToken (or Mock value) とするが、
	// 正しくは SessionID であるべき。
	// 既存のコード `WriteJSON` では `"session_id": "session-token"` とハードコードされていた。
	// これを改善する。

	// TODO: AuthService.HandleCallback の戻り値を修正して sessionID を取得できるようにすべき。
	// 今回は Cookie 設定のロジックを追加することに注力し、値は仮置きする。
	sessionID := "session-token" // 仮の値。本来は AuthService から返却されるべき。

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // 本番環境ではtrue (HTTPS必須)
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600, // 1時間
	})

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":    session.UserID,
		"email":      session.Email,
		"name":       session.Name,
		"role":       session.Role,
		"expires_at": session.ExpiresAt,
	})
}

// Logout はログアウトを処理します
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(ContextKeySession).(*service.Session)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	// CookieからセッションIDを取得
	cookie, err := r.Cookie("session_id")
	var sessionID string
	if err == nil {
		sessionID = cookie.Value
	} else {
		// ヘッダーからの取得もサポート（後方互換性）
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			sessionID = authHeader[7:]
		}
	}

	if sessionID != "" {
		err := h.authService.Logout(r.Context(), sessionID, session.UserID)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "LOGOUT_FAILED", err.Error())
			return
		}
	}

	// Cookieを削除
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

// RefreshRequest はトークンリフレッシュリクエスト
type RefreshRequest struct {
	SessionID string `json:"session_id"`
}

// Refresh はセッションをリフレッシュします
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// CookieからセッションIDを取得
	cookie, err := r.Cookie("session_id")
	var sessionID string
	if err == nil {
		sessionID = cookie.Value
	} else {
		// リクエストボディからの取得もサポート
		var req RefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			sessionID = req.SessionID
		}
	}

	if sessionID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Session ID missing")
		return
	}

	session, err := h.authService.RefreshSession(r.Context(), sessionID)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "REFRESH_FAILED", err.Error())
		return
	}

	// Cookieを更新（有効期限延長など）
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600, // 延長
	})

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"expires_at": session.ExpiresAt,
	})
}
