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

type StockHandler struct {
	store *store.Store
}

func NewStockHandler(s *store.Store) *StockHandler {
	return &StockHandler{store: s}
}

func (h *StockHandler) ListForItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	stocks, err := h.store.ListStock(itemID)
	if err != nil {
		writeError(w, "failed to list stock", http.StatusInternalServerError, err)
		return
	}
	if stocks == nil {
		stocks = []model.Stock{}
	}
	writeJSON(w, http.StatusOK, stocks)
}

func (h *StockHandler) Add(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var st model.Stock
	if err := json.NewDecoder(r.Body).Decode(&st); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	st.ItemID = itemID

	if err := h.store.AddStock(&st); err != nil {
		writeError(w, "failed to add stock", http.StatusInternalServerError, err)
		return
	}

	// Audit log
	if item, err := h.store.GetItem(itemID); err == nil {
		_ = h.store.AddAuditEntry(item.Name, st.Quantity)
	}

	created, err := h.store.GetStock(st.ID)
	if err != nil {
		writeError(w, "failed to get created stock", http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *StockHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "stockID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var body struct {
		Quantity float64 `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.store.UpdateStock(id, body.Quantity); err != nil {
		writeError(w, "failed to update stock", http.StatusInternalServerError, err)
		return
	}

	if body.Quantity <= 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	updated, err := h.store.GetStock(id)
	if err != nil {
		writeError(w, "failed to get updated stock", http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *StockHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "stockID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.store.DeleteStock(id); err != nil {
		writeError(w, "failed to delete stock", http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *StockHandler) Move(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "stockID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req model.StockMoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.store.MoveStock(id, req); err != nil {
		http.Error(w, "failed to move stock", http.StatusBadRequest)
		return
	}

	st, err := h.store.GetStock(id)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		writeError(w, "failed to get moved stock", http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, st)
}
