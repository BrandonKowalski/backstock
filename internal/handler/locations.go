package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"backstock/internal/model"
	"backstock/internal/store"

	"github.com/go-chi/chi/v5"
)

type LocationHandler struct {
	store *store.Store
}

func NewLocationHandler(s *store.Store) *LocationHandler {
	return &LocationHandler{store: s}
}

func (h *LocationHandler) List(w http.ResponseWriter, r *http.Request) {
	locs, err := h.store.ListLocations()
	if err != nil {
		writeError(w, "failed to list locations", http.StatusInternalServerError, err)
		return
	}
	if locs == nil {
		locs = []model.Location{}
	}
	writeJSON(w, http.StatusOK, locs)
}

func (h *LocationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var loc model.Location
	if err := json.NewDecoder(r.Body).Decode(&loc); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if loc.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if err := h.store.CreateLocation(&loc); err != nil {
		writeError(w, "failed to create location", http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, loc)
}

func (h *LocationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var loc model.Location
	if err := json.NewDecoder(r.Body).Decode(&loc); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	loc.ID = id
	if err := h.store.UpdateLocation(&loc); err != nil {
		writeError(w, "failed to update location", http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, loc)
}

func (h *LocationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.store.DeleteLocation(id); err != nil {
		writeError(w, "failed to delete location", http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
