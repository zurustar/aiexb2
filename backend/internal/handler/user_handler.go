// backend/internal/handler/user_handler.go
package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/repository"
	"github.com/your-org/esms/internal/service"
)

// UserHandler はユーザー関連のHTTPハンドラー
type UserHandler struct {
	userRepo repository.UserRepository
}

// NewUserHandler は新しいUserHandlerを作成します
func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// RegisterRoutes はルートを登録します
func (h *UserHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/users", h.ListUsers).Methods("GET")
	r.HandleFunc("/api/v1/users/{id}", h.GetUser).Methods("GET")
	r.HandleFunc("/api/v1/users/me", h.GetCurrentUser).Methods("GET")
}

// ListUsers はユーザー一覧を取得します（管理者のみ）
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(ContextKeySession).(*service.Session)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	// 管理者のみアクセス可能
	if session.Role != domain.RoleAdmin {
		WriteError(w, http.StatusForbidden, "FORBIDDEN", "Admin access required")
		return
	}

	// TODO: リポジトリにListメソッドを追加するか、別の方法で一覧取得
	WriteJSON(w, http.StatusOK, []interface{}{})
}

// GetUser はユーザーを取得します（管理者のみ）
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(ContextKeySession).(*service.Session)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	// 管理者のみアクセス可能
	if session.Role != domain.RoleAdmin {
		WriteError(w, http.StatusForbidden, "FORBIDDEN", "Admin access required")
		return
	}

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	WriteJSON(w, http.StatusOK, user)
}

// GetCurrentUser は現在のユーザー情報を取得します（全ての認証済みユーザー）
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(ContextKeySession).(*service.Session)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), session.UserID)
	if err != nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	WriteJSON(w, http.StatusOK, user)
}
