package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/val/inventory/internal/models"
	"github.com/val/inventory/internal/repository"
)

type ReceiptHandler struct {
	repo *repository.ReceiptRepo
}

type receiptRequest struct {
	models.Receipt
	Details []models.ReceiptDetail `json:"details"`
}

func (h *ReceiptHandler) List(w http.ResponseWriter, r *http.Request) {
	page, perPage := pageParams(r)
	acCode := r.URL.Query().Get("ac_code")

	receipts, total, err := h.repo.List(page, perPage, acCode)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, models.PagedResult[models.Receipt]{
		Data: receipts, Total: total, Page: page, PerPage: perPage,
	})
}

func (h *ReceiptHandler) Get(w http.ResponseWriter, r *http.Request) {
	tNo := chi.URLParam(r, "tNo")
	recpt, err := h.repo.Get(tNo)
	if err != nil {
		respondErr(w, http.StatusNotFound, "receipt not found")
		return
	}
	respond(w, http.StatusOK, recpt)
}

func (h *ReceiptHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req receiptRequest
	if err := decode(r, &req); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.repo.Create(&req.Receipt, req.Details); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusCreated, req.Receipt)
}

func (h *ReceiptHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req receiptRequest
	if err := decode(r, &req); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	req.Receipt.TNo = chi.URLParam(r, "tNo")
	if err := h.repo.Update(&req.Receipt, req.Details); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, req.Receipt)
}

func (h *ReceiptHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tNo := chi.URLParam(r, "tNo")
	if err := h.repo.Delete(tNo); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]string{"status": "deleted"})
}
