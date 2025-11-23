// backend/internal/util/logger_test.go
package util_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/esms/internal/util"
)

func TestInitLogger(t *testing.T) {
	// Dev mode (TextHandler)
	loggerDev := util.InitLogger(true)
	assert.NotNil(t, loggerDev)

	// Prod mode (JSONHandler)
	loggerProd := util.InitLogger(false)
	assert.NotNil(t, loggerProd)
}

func TestWithTraceID(t *testing.T) {
	// バッファに出力するようにロガーを一時的に変更
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	ctx := context.WithValue(context.Background(), util.TraceIDKey, "test-trace-id")

	util.LogInfo(ctx, "test message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)

	assert.Equal(t, "test message", logEntry["msg"])
	assert.Equal(t, "test-trace-id", logEntry["trace_id"])
	assert.Equal(t, "INFO", logEntry["level"])
}
