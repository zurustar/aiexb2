// backend/internal/util/logger.go
package util

import (
	"context"
	"io"
	"log/slog"
	"os"
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
		Level: slog.LevelInfo,
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
