package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"backstock/internal/model"
	"backstock/internal/store"

	"github.com/go-chi/chi/v5"
)

type UnitHandler struct {
	store *store.Store
}

func NewUnitHandler(s *store.Store) *UnitHandler {
	return &UnitHandler{store: s}
}

func (h *UnitHandler) List(w http.ResponseWriter, r *http.Request) {
	units, err := h.store.ListUnits()
	if err != nil {
		writeError(w, "failed to list units", http.StatusInternalServerError, err)
		return
	}
	if units == nil {
		units = []model.Unit{}
	}
	writeJSON(w, http.StatusOK, units)
}

func (h *UnitHandler) Create(w http.ResponseWriter, r *http.Request) {
	var u model.Unit
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if u.Name == "" || u.Abbreviation == "" {
		http.Error(w, "name and abbreviation are required", http.StatusBadRequest)
		return
	}
	if err := h.store.CreateUnit(&u); err != nil {
		writeError(w, "failed to create unit", http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, u)
}

func (h *UnitHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var u model.Unit
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	u.ID = id
	if err := h.store.UpdateUnit(&u); err != nil {
		writeError(w, "failed to update unit", http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (h *UnitHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.store.DeleteUnit(id); err != nil {
		writeError(w, "failed to delete unit", http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
