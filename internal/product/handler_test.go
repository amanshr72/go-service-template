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

func validPayload() []byte {
	b, _ := json.Marshal(CreateProductRequest{
		Name: "Keyboard", Price: 49.99, Description: "desc", Category: "cat", Stock: 1,
	})
	return b
}

func TestProductHandler_Create(t *testing.T) {
	router := NewRouter(NewService(NewInMemoryRepository()))
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(validPayload()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestProductHandler_GetAll(t *testing.T) {
	svc := NewService(NewInMemoryRepository())
	_, _ = svc.Create(CreateProductRequest{Name: "Mouse", Price: 20, Description: "d", Category: "c", Stock: 1})
	router := NewRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_GetByID_NotFound(t *testing.T) {
	router := NewRouter(NewService(NewInMemoryRepository()))
	req := httptest.NewRequest(http.MethodGet, "/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProductHandler_GetByID_Found(t *testing.T) {
	svc := NewService(NewInMemoryRepository())
	p, _ := svc.Create(CreateProductRequest{Name: "Monitor", Price: 199, Description: "d", Category: "c", Stock: 1})
	router := NewRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/"+strconv.FormatInt(p.Id, 10), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
