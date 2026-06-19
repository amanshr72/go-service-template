package smoke

import (
	"database/sql"
	"errors"
	"go-crud2/internal/auth"
	"go-crud2/internal/health"
	"go-crud2/internal/notification"
	"go-crud2/internal/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockNotifier struct {
	SendCalled bool
	ShouldFail bool
}

func (m *MockNotifier) SendEmail(req notification.SendEmailRequest) (*notification.SendEmailResponse, error) {
	m.SendCalled = true
	if m.ShouldFail {
		return nil, errors.New("mock vendor failure")
	}
	return &notification.SendEmailResponse{MessageID: "mock_123", Status: "queued"}, nil
}

func setupServer() http.Handler {
	repo := user.NewInMemoryRepository()
	notifier := &MockNotifier{}
	svc := user.NewService(repo, notifier)

	mux := http.NewServeMux()

	var db *sql.DB

	auth.RegisterRoutes(mux)
	health.RegisterRoutes(mux, db)

	_ = user.RegisterRoutes(mux, svc)

	return mux
}

func TestSmoke_Health(t *testing.T) {
	server := setupServer()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	assert.NotEqual(t, http.StatusInternalServerError, rec.Code)
}

func TestSmoke_Users(t *testing.T) {
	server := setupServer()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	assert.NotEqual(t, http.StatusInternalServerError, rec.Code)
}

func TestSmoke_Login(t *testing.T) {
	server := setupServer()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	assert.NotEqual(t, http.StatusInternalServerError, rec.Code)
}

func TestSmoke_GraphQL(t *testing.T) {
	server := setupServer()
	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	assert.NotEqual(t, http.StatusInternalServerError, rec.Code)
}
