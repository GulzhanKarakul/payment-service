package middleware_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GulzhanKarakul/payment-service/internal/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestLoggerMiddleware_SetsRequestIDHeader(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := middleware.Logger(testLogger())(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	
	requestID := rec.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID)

	_, err := uuid.Parse(requestID)
	require.NoError(t, err, "X-Request-ID is valid UUID")
}