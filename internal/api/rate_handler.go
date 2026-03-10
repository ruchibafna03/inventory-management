package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/val/inventory/internal/models"
	"github.com/val/inventory/internal/repository"
)

type RateHandler struct {
	repo *repository.RateRepo
}

func (h *RateHandler) List(w http.ResponseWriter, r *http.Request) {
	rates, err := h.repo.List()
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, rates)
}

func (h *RateHandler) Latest(w http.ResponseWriter, r *http.Request) {
	rate, err := h.repo.Latest()
	if err != nil {
		respondErr(w, http.StatusNotFound, "no rates found")
		return
	}
	respond(w, http.StatusOK, rate)
}

func (h *RateHandler) Create(w http.ResponseWriter, r *http.Request) {
	var rate models.Rate
	if err := decode(r, &rate); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.repo.Create(&rate); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusCreated, rate)
}

func (h *RateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondErr(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.repo.Delete(id); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// ─── Lot Handler ─────────────────────────────────────────────────────────────

type LotHandler struct {
	repo *repository.LotRepo
}

func (h *LotHandler) List(w http.ResponseWriter, r *http.Request) {
	acCode := r.URL.Query().Get("ac_code")
	lots, err := h.repo.List(acCode)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, lots)
}

func (h *LotHandler) Get(w http.ResponseWriter, r *http.Request) {
	tNo := chi.URLParam(r, "tNo")
	lot, err := h.repo.Get(tNo)
	if err != nil {
		respondErr(w, http.StatusNotFound, "lot not found")
		return
	}
	respond(w, http.StatusOK, lot)
}

func (h *LotHandler) Create(w http.ResponseWriter, r *http.Request) {
	var lot models.Lot
	if err := decode(r, &lot); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.repo.Create(&lot); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusCreated, lot)
}

func (h *LotHandler) Update(w http.ResponseWriter, r *http.Request) {
	var lot models.Lot
	if err := decode(r, &lot); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	lot.TNo = chi.URLParam(r, "tNo")
	if err := h.repo.Update(&lot); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, lot)
}

func (h *LotHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tNo := chi.URLParam(r, "tNo")
	if err := h.repo.Delete(tNo); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]string{"status": "deleted"})
}
