package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GulzhanKarakul/payment-service/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryMiddleware_CatchesPanic_Returns500(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("PANIC: something went wrong")
	})

	h := middleware.Recovery(testLogger())(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		h.ServeHTTP(rec, req)
	})

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}