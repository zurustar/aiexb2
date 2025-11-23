// backend/internal/handler/middleware.go
package handler

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/your-org/esms/internal/service"
	"golang.org/x/time/rate"
)

// ContextKey はコンテキストキーの型
type ContextKey string

const (
	// ContextKeySession はセッション情報のコンテキストキー
	ContextKeySession ContextKey = "session"
	// ContextKeyRequestID はリクエストIDのコンテキストキー
	ContextKeyRequestID ContextKey = "request_id"
)

// AuthServiceInterface は認証サービスのインターフェース
type AuthServiceInterface interface {
	GetSession(sessionID string) (*service.Session, error)
}

// Middleware はミドルウェアの集合
type Middleware struct {
	authService AuthServiceInterface
	rateLimiter *RateLimiter
}

// NewMiddleware は新しいMiddlewareを作成します
func NewMiddleware(authService AuthServiceInterface) *Middleware {
	return &Middleware{
		authService: authService,
		rateLimiter: NewRateLimiter(100, 10), // 100 req/sec, burst 10
	}
}

// CORS はCORSヘッダーを設定するミドルウェア
func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 許可するオリジンのリスト（環境変数等から取得するのが望ましい）
		allowedOrigins := map[string]bool{
			"http://localhost:3000": true,
			"http://localhost:8080": true,
		}

		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// プリフライトリクエストの処理
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logging はリクエストをログに記録するミドルウェア
func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// レスポンスライター
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		log.Printf(
			"%s %s %d %s %s",
			r.Method,
			r.RequestURI,
			lrw.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}

// loggingResponseWriter はステータスコードを記録するResponseWriter
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Authentication は認証を行うミドルウェア
func (m *Middleware) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authorizationヘッダーからトークンを取得
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Bearer トークンを抽出
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// セッションIDからセッションを取得
		session, err := m.authService.GetSession(token)
		if err != nil {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		// セッション情報をコンテキストに追加
		ctx := context.WithValue(r.Context(), ContextKeySession, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole は指定されたロールを要求するミドルウェア
func (m *Middleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, ok := r.Context().Value(ContextKeySession).(*service.Session)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// ロールチェック（簡易実装）
			// 本番環境ではdomain.Roleを使用
			if string(session.Role) != role && session.Role != "admin" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimit はレート制限を適用するミドルウェア
func (m *Middleware) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// IPアドレスベースのレート制限
		ip := getIP(r)
		if !m.rateLimiter.Allow(ip) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CSRF はCSRF対策を行うミドルウェア（簡易実装）
func (m *Middleware) CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// GET, HEAD, OPTIONS は CSRF チェック不要
		if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// X-CSRF-Token ヘッダーをチェック
		csrfToken := r.Header.Get("X-CSRF-Token")
		if csrfToken == "" {
			http.Error(w, "CSRF token missing", http.StatusForbidden)
			return
		}

		// TODO: トークンの検証
		// 本番環境では適切なCSRFトークン検証を実装

		next.ServeHTTP(w, r)
	})
}

// getIP はリクエストからIPアドレスを取得します
func getIP(r *http.Request) string {
	// X-Forwarded-For ヘッダーをチェック
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// X-Real-IP ヘッダーをチェック
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// RemoteAddr を使用
	ip := r.RemoteAddr
	// ポート番号を削除
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

// RateLimiter はIPアドレスベースのレート制限
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter は新しいRateLimiterを作成します
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

// Allow はIPアドレスに対してリクエストを許可するか判定します
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
	}
	rl.mu.Unlock()

	return limiter.Allow()
}
