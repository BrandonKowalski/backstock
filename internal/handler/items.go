package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"backstock/internal/model"
	"backstock/internal/store"

	"github.com/go-chi/chi/v5"
)

type ItemHandler struct {
	store *store.Store
}

func NewItemHandler(s *store.Store) *ItemHandler {
	return &ItemHandler{store: s}
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	f := model.ItemFilter{
		Location: r.URL.Query().Get("location"),
		Category: r.URL.Query().Get("category"),
		Sort:     r.URL.Query().Get("sort"),
		Search:   r.URL.Query().Get("search"),
	}
	items, err := h.store.ListItems(f)
	if err != nil {
		writeError(w, "failed to list items", http.StatusInternalServerError, err)
		return
	}
	if items == nil {
		items = []model.Item{}
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *ItemHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	item, err := h.store.GetItem(id)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		writeError(w, "failed to get item", http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

type createItemRequest struct {
	Name                 string   `json:"name"`
	IsFood               bool     `json:"is_food"`
	UnitID               *int     `json:"unit_id"`
	PackageSize          *float64 `json:"package_size"`
	ExpirationDate       *string  `json:"expiration_date"`
	BestByDate           *string  `json:"best_by_date"`
	LowQuantityThreshold *float64 `json:"low_quantity_threshold"`
	CategoryIDs          []int    `json:"category_ids"`
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	item := model.Item{
		Name:                 req.Name,
		IsFood:               req.IsFood,
		UnitID:               req.UnitID,
		PackageSize:          req.PackageSize,
		ExpirationDate:       req.ExpirationDate,
		BestByDate:           req.BestByDate,
		LowQuantityThreshold: req.LowQuantityThreshold,
	}
	if err := h.store.CreateItem(&item, req.CategoryIDs); err != nil {
		writeError(w, "failed to create item", http.StatusInternalServerError, err)
		return
	}

	created, err := h.store.GetItem(item.ID)
	if err != nil {
		writeError(w, "failed to get created item", http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req createItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	item := model.Item{
		ID:                   id,
		Name:                 req.Name,
		IsFood:               req.IsFood,
		UnitID:               req.UnitID,
		PackageSize:          req.PackageSize,
		ExpirationDate:       req.ExpirationDate,
		BestByDate:           req.BestByDate,
		LowQuantityThreshold: req.LowQuantityThreshold,
	}
	if err := h.store.UpdateItem(&item, req.CategoryIDs); err != nil {
		writeError(w, "failed to update item", http.StatusInternalServerError, err)
		return
	}

	updated, err := h.store.GetItem(id)
	if err != nil {
		writeError(w, "failed to get updated item", http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.store.DeleteItem(id); err != nil {
		writeError(w, "failed to delete item", http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
