package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/val/inventory/internal/repository"
)

type UtilityHandler struct {
	repo *repository.UtilityRepo
}

func (h *UtilityHandler) DataSummary(w http.ResponseWriter, r *http.Request) {
	counts, err := h.repo.DataSummary()
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, counts)
}

func (h *UtilityHandler) OrphanIssueDetails(w http.ResponseWriter, r *http.Request) {
	rows, err := h.repo.OrphanIssueDetails()
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, rows)
}

func (h *UtilityHandler) OrphanReceiptDetails(w http.ResponseWriter, r *http.Request) {
	rows, err := h.repo.OrphanReceiptDetails()
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, rows)
}

func (h *UtilityHandler) StockPosition(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.StockPosition()
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, items)
}

func (h *UtilityHandler) ItemHistory(w http.ResponseWriter, r *http.Request) {
	itcd := chi.URLParam(r, "itcd")
	events, err := h.repo.ItemHistory(itcd)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, events)
}

func (h *UtilityHandler) ChangeItemCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	if err := decode(r, &req); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.From == "" || req.To == "" {
		respondErr(w, http.StatusBadRequest, "from and to are required")
		return
	}
	if err := h.repo.ChangeItemCode(req.From, req.To); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *UtilityHandler) CodeChangeHistory(w http.ResponseWriter, r *http.Request) {
	rows, err := h.repo.CodeChangeHistory()
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, rows)
}

func (h *UtilityHandler) BlockedItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.BlockedItems()
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, items)
}

func (h *UtilityHandler) BlockItem(w http.ResponseWriter, r *http.Request) {
	itcd := chi.URLParam(r, "itcd")
	if err := h.repo.BlockItem(itcd); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]string{"status": "blocked"})
}

func (h *UtilityHandler) UnblockItem(w http.ResponseWriter, r *http.Request) {
	itcd := chi.URLParam(r, "itcd")
	if err := h.repo.UnblockItem(itcd); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]string{"status": "unblocked"})
}

func (h *UtilityHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.ListUsers()
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, users)
}

func (h *UtilityHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username    string `json:"username"`
		NewPassword string `json:"new_password"`
	}
	if err := decode(r, &req); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Username == "" || req.NewPassword == "" {
		respondErr(w, http.StatusBadRequest, "username and new_password are required")
		return
	}
	if len(req.NewPassword) < 4 {
		respondErr(w, http.StatusBadRequest, "password must be at least 4 characters")
		return
	}
	if err := h.repo.ChangePassword(req.Username, req.NewPassword); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]string{"status": "password changed"})
}
