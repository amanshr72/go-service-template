package product

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON"})
		return
	}
	p, err := h.svc.Create(req)
	if err != nil {
		var ve *ValidationError
		if errors.As(err, &ve) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "validation failed", Fields: &ve.Fields})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.GetAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, products)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}
	p, err := h.svc.GetByID(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "product not found"})
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
