package handler

import (
	"net/http"
	"backstock/internal/model"
	"backstock/internal/store"
)

type AuditHandler struct {
	store *store.Store
}

func NewAuditHandler(s *store.Store) *AuditHandler {
	return &AuditHandler{store: s}
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	entries, err := h.store.ListAuditLog()
	if err != nil {
		writeError(w, "failed to list audit log", http.StatusInternalServerError, err)
		return
	}
	if entries == nil {
		entries = []model.AuditEntry{}
	}
	writeJSON(w, http.StatusOK, entries)
}
