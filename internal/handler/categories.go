package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"backstock/internal/model"
	"backstock/internal/store"

	"github.com/go-chi/chi/v5"
)

type CategoryHandler struct {
	store *store.Store
}

func NewCategoryHandler(s *store.Store) *CategoryHandler {
	return &CategoryHandler{store: s}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	cats, err := h.store.ListCategories()
	if err != nil {
		writeError(w, "failed to list categories", http.StatusInternalServerError, err)
		return
	}
	if cats == nil {
		cats = []model.Category{}
	}
	writeJSON(w, http.StatusOK, cats)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c model.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if c.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if err := h.store.CreateCategory(&c); err != nil {
		writeError(w, "failed to create category", http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var c model.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	c.ID = id
	if err := h.store.UpdateCategory(&c); err != nil {
		writeError(w, "failed to update category", http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.store.DeleteCategory(id); err != nil {
		writeError(w, "failed to delete category", http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
