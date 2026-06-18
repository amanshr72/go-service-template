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
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	return NewRouter(svc)
}

func TestProduct_CreateAndGetByID(t *testing.T) {
	router := setupRouter()

	payload := map[string]any{
		"name":  "MacBook",
		"price": 120000.0,
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewBuffer(body),
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var created Product

	err := json.NewDecoder(rec.Body).Decode(&created)
	assert.NoError(t, err)

	req2 := httptest.NewRequest(
		http.MethodGet,
		"/"+strconv.Itoa(created.ID),
		nil,
	)

	rec2 := httptest.NewRecorder()

	router.ServeHTTP(rec2, req2)

	assert.Equal(t, http.StatusOK, rec2.Code)

	var fetched Product

	err = json.NewDecoder(rec2.Body).Decode(&fetched)
	assert.NoError(t, err)

	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, "MacBook", fetched.Name)
}

func TestProduct_GetAll(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestProduct_GetByID_NotFound(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/999",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
