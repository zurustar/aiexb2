// backend/internal/util/logger.go
package util

import (
	"context"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

type contextKey string

const (
	TraceIDKey contextKey = "trace_id"
)

// InitLogger はロガーを初期化します
// w が nil の場合は os.Stdout に出力します
func InitLogger(isDev bool, w io.Writer) *slog.Logger {
	if w == nil {
		w = os.Stdout
	}
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		ReplaceAttr: maskPII,
	}
	if isDev {
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(w, opts)
	} else {
		handler = slog.NewJSONHandler(w, opts)
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

// WithTraceID はコンテキストからトレースIDを取得し、ログフィールドに追加したロガーを返します
func WithTraceID(ctx context.Context) *slog.Logger {
	logger := slog.Default()
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return logger.With("trace_id", traceID)
	}
	return logger
}

// LogError はエラーログを出力するヘルパー関数
func LogError(ctx context.Context, msg string, err error) {
	WithTraceID(ctx).Error(msg, "error", err)
}

// LogInfo は情報ログを出力するヘルパー関数
func LogInfo(ctx context.Context, msg string, args ...any) {
	WithTraceID(ctx).Info(msg, args...)
}

// maskPII はログ出力時にPII（個人情報）をマスクする関数
func maskPII(groups []string, a slog.Attr) slog.Attr {
	// キー名によるフィルタリング
	key := strings.ToLower(a.Key)
	if key == "email" || key == "password" || key == "token" || key == "secret" {
		return slog.String(a.Key, "***MASKED***")
	}

	// 文字列値のコンテンツフィルタリング（簡易的なメールアドレス検出）
	if a.Value.Kind() == slog.KindString {
		val := a.Value.String()
		if strings.Contains(val, "@") {
			// 簡易的なメールアドレス正規表現
			re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
			if re.MatchString(val) {
				masked := re.ReplaceAllString(val, "***@***.***")
				return slog.String(a.Key, masked)
			}
		}
	}

	return a
}
