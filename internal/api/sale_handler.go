package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/val/inventory/internal/models"
	"github.com/val/inventory/internal/repository"
)

type SaleHandler struct {
	repo *repository.SaleRepo
}

func (h *SaleHandler) List(w http.ResponseWriter, r *http.Request) {
	page, perPage := pageParams(r)
	acCode := r.URL.Query().Get("ac_code")
	sales, total, err := h.repo.List(page, perPage, acCode)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, models.PagedResult[models.Sale]{
		Data: sales, Total: total, Page: page, PerPage: perPage,
	})
}

func (h *SaleHandler) Get(w http.ResponseWriter, r *http.Request) {
	vouchNo := chi.URLParam(r, "vouchNo")
	sale, err := h.repo.Get(vouchNo)
	if err != nil {
		respondErr(w, http.StatusNotFound, "sale not found")
		return
	}
	respond(w, http.StatusOK, sale)
}

func (h *SaleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var sale models.Sale
	if err := decode(r, &sale); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.repo.Create(&sale); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusCreated, sale)
}

func (h *SaleHandler) Update(w http.ResponseWriter, r *http.Request) {
	var sale models.Sale
	if err := decode(r, &sale); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	sale.VouchNo = chi.URLParam(r, "vouchNo")
	if err := h.repo.Update(&sale); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, sale)
}

func (h *SaleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vouchNo := chi.URLParam(r, "vouchNo")
	if err := h.repo.Delete(vouchNo); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// ─── Purchase Handler ────────────────────────────────────────────────────────

type PurchaseHandler struct {
	repo *repository.PurchaseRepo
}

func (h *PurchaseHandler) List(w http.ResponseWriter, r *http.Request) {
	page, perPage := pageParams(r)
	acCode := r.URL.Query().Get("ac_code")
	purchases, total, err := h.repo.List(page, perPage, acCode)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, models.PagedResult[models.Purchase]{
		Data: purchases, Total: total, Page: page, PerPage: perPage,
	})
}

func (h *PurchaseHandler) Get(w http.ResponseWriter, r *http.Request) {
	vouchNo := chi.URLParam(r, "vouchNo")
	p, err := h.repo.Get(vouchNo)
	if err != nil {
		respondErr(w, http.StatusNotFound, "purchase not found")
		return
	}
	respond(w, http.StatusOK, p)
}

func (h *PurchaseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p models.Purchase
	if err := decode(r, &p); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.repo.Create(&p); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusCreated, p)
}

func (h *PurchaseHandler) Update(w http.ResponseWriter, r *http.Request) {
	var p models.Purchase
	if err := decode(r, &p); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	p.VouchNo = chi.URLParam(r, "vouchNo")
	if err := h.repo.Update(&p); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, p)
}

func (h *PurchaseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vouchNo := chi.URLParam(r, "vouchNo")
	if err := h.repo.Delete(vouchNo); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]string{"status": "deleted"})
}
