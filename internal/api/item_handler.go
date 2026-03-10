package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/val/inventory/internal/models"
	"github.com/val/inventory/internal/repository"
)

type ItemHandler struct {
	repo *repository.ItemRepo
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	page, perPage := pageParams(r)
	search := r.URL.Query().Get("search")
	cat := r.URL.Query().Get("cat")

	items, total, err := h.repo.List(page, perPage, search, cat)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, models.PagedResult[models.Item]{
		Data: items, Total: total, Page: page, PerPage: perPage,
	})
}

func (h *ItemHandler) Get(w http.ResponseWriter, r *http.Request) {
	itcd := chi.URLParam(r, "itcd")
	item, err := h.repo.Get(itcd)
	if err != nil {
		respondErr(w, http.StatusNotFound, "item not found")
		return
	}
	respond(w, http.StatusOK, item)
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := decode(r, &item); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.repo.Create(&item); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusCreated, item)
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := decode(r, &item); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	item.ItCd = chi.URLParam(r, "itcd")
	if err := h.repo.Update(&item); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, item)
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	itcd := chi.URLParam(r, "itcd")
	if err := h.repo.Delete(itcd); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *ItemHandler) Tags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.repo.Tags()
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, tags)
}
