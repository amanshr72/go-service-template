package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin_Success(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	body, _ := json.Marshal(LoginInput{Email: "admin@test.com", Password: "password"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp TokenResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp.Token) // got a token back
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	body, _ := json.Marshal(LoginInput{Email: "wrong@test.com", Password: "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestValidateToken_RoundTrip(t *testing.T) {
	// Sign a token then immediately validate it
	resp, err := Login(LoginInput{Email: "admin@test.com", Password: "password"})
	assert.NoError(t, err)

	claims, err := ValidateToken(resp.Token)
	assert.NoError(t, err)
	assert.Equal(t, "admin@test.com", claims.Email)
	assert.Equal(t, 1, claims.UserID)
}
