package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- RequestID ---

func TestRequestID_GeneratesID(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		id := GetRequestID(r)
		assert.NotEmpty(t, id) // ID must be in context
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.NotEmpty(t, w.Header().Get("X-Request-ID")) // ID in response header
}

func TestRequestID_HonoursExisting(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "my-trace-id")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, "my-trace-id", w.Header().Get("X-Request-ID"))
}

// --- Recovery ---

func TestRecovery_CatchesPanic(t *testing.T) {
	handler := Recovery(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		panic("something broke")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req) // should NOT panic the test

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- Chain ---

func TestChain_OrderIsCorrect(t *testing.T) {
	order := []string{}

	m1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m1")
			next.ServeHTTP(w, r)
		})
	}
	m2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m2")
			next.ServeHTTP(w, r)
		})
	}
	core := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	})

	chained := Chain(core, m1, m2)
	chained.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	// m1 listed first → runs first
	assert.Equal(t, []string{"m1", "m2", "handler"}, order)
}

// --- Authenticate ---

func TestAuthenticate_MissingToken(t *testing.T) {
	handler := Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	handler := Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
