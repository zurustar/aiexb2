// backend/internal/handler/auth_handler.go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/your-org/esms/internal/service"
)

// AuthHandler は認証関連のHTTPハンドラー
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler は新しいAuthHandlerを作成します
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
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

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"session_id": "session-token", // 実際のセッションIDを返す
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

	// セッションIDを取得（簡易実装）
	sessionID := r.Header.Get("Authorization")

	err := h.authService.Logout(r.Context(), sessionID, session.UserID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "LOGOUT_FAILED", err.Error())
		return
	}

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
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	session, err := h.authService.RefreshSession(r.Context(), req.SessionID)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "REFRESH_FAILED", err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"session_id": req.SessionID,
		"expires_at": session.ExpiresAt,
	})
}
