package product

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductHandler_Create(t *testing.T) {
	svc := NewService(NewInMemoryRepository())
	router := NewRouter(svc)

	body, _ := json.Marshal(Product{Name: "Keyboard", Price: 49.99})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestProductHandler_GetAll(t *testing.T) {
	svc := NewService(NewInMemoryRepository())
	_, _ = svc.Create("Mouse", 20)
	router := NewRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_GetByID_NotFound(t *testing.T) {
	svc := NewService(NewInMemoryRepository())
	router := NewRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProductHandler_GetByID_Found(t *testing.T) {
	svc := NewService(NewInMemoryRepository())
	p, _ := svc.Create("Monitor", 199)
	router := NewRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/"+strconv.Itoa(p.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
