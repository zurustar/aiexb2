// backend/internal/handler/router.go
package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/your-org/esms/internal/repository"
	"github.com/your-org/esms/internal/service"
)

// Router はアプリケーションのルーターを設定します
type Router struct {
	router *mux.Router
	mw     *Middleware
}

// NewRouter は新しいRouterを作成します
func NewRouter(
	authService *service.AuthService,
	reservationService *service.ReservationService,
	approvalService *service.ApprovalService,
	userRepo repository.UserRepository,
	resourceRepo repository.ResourceRepository,
) *Router {
	r := mux.NewRouter()
	mw := NewMiddleware(authService)

	router := &Router{
		router: r,
		mw:     mw,
	}

	// グローバルミドルウェア（全てのルートに適用）
	r.Use(mw.CORS)
	r.Use(mw.Logging)
	r.Use(mw.RateLimit)

	// 公開エンドポイント（認証不要）
	// ヘルスチェック
	r.HandleFunc("/health", HealthCheck).Methods("GET")

	// 認証エンドポイント（認証不要）
	authHandler := NewAuthHandler(authService)
	authHandler.RegisterRoutes(r)

	// 認証が必要なルート
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(mw.Authentication)
	protected.Use(mw.CSRF)

	reservationHandler := NewReservationHandler(reservationService, approvalService)
	reservationHandler.RegisterRoutes(protected)

	resourceHandler := NewResourceHandler(resourceRepo)
	resourceHandler.RegisterRoutes(protected)

	userHandler := NewUserHandler(userRepo)
	userHandler.RegisterRoutes(protected)

	// カスタム404/405ハンドラー
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	r.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowedHandler)

	return router
}

// GetRouter はmux.Routerを返します
func (r *Router) GetRouter() *mux.Router {
	return r.router
}

// HealthCheck はヘルスチェックエンドポイント
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// NotFoundHandler は404エラーハンドラー
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	WriteError(w, http.StatusNotFound, "NOT_FOUND", "The requested resource was not found")
}

// MethodNotAllowedHandler は405エラーハンドラー
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The request method is not allowed for this resource")
}
