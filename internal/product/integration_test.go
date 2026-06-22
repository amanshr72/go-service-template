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

func setupRouter() http.Handler {
	return NewRouter(NewService(NewInMemoryRepository()))
}

func TestProduct_CreateAndGetByID(t *testing.T) {
	router := setupRouter()
	payload, _ := json.Marshal(CreateProductRequest{
		Name: "MacBook", Price: 120000, Description: "laptop", Category: "electronics", Stock: 5,
	})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var created Product
	_ = json.NewDecoder(rec.Body).Decode(&created)

	req2 := httptest.NewRequest(http.MethodGet, "/"+strconv.FormatInt(created.Id, 10), nil)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	assert.Equal(t, http.StatusOK, rec2.Code)

	var fetched Product
	_ = json.NewDecoder(rec2.Body).Decode(&fetched)
	assert.Equal(t, created.Id, fetched.Id)
	assert.Equal(t, "MacBook", fetched.Name)
}

func TestProduct_GetAll(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestProduct_GetByID_NotFound(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
