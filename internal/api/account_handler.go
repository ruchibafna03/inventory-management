package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/val/inventory/internal/models"
	"github.com/val/inventory/internal/repository"
)

type AccountHandler struct {
	repo *repository.AccountRepo
}

func (h *AccountHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.repo.ListGroups()
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, groups)
}

func (h *AccountHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var g models.GroupMaster
	if err := decode(r, &g); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.repo.CreateGroup(&g); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusCreated, g)
}

func (h *AccountHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	var g models.GroupMaster
	if err := decode(r, &g); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	g.GCode = chi.URLParam(r, "gCode")
	if err := h.repo.UpdateGroup(&g); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, g)
}

func (h *AccountHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	gCode := chi.URLParam(r, "gCode")
	if err := h.repo.DeleteGroup(gCode); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	page, perPage := pageParams(r)
	search := r.URL.Query().Get("search")
	gCode := r.URL.Query().Get("g_code")

	accounts, total, err := h.repo.List(page, perPage, search, gCode)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, models.PagedResult[models.AccountMaster]{
		Data: accounts, Total: total, Page: page, PerPage: perPage,
	})
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	acCode := chi.URLParam(r, "acCode")
	acc, err := h.repo.Get(acCode)
	if err != nil {
		respondErr(w, http.StatusNotFound, "account not found")
		return
	}
	respond(w, http.StatusOK, acc)
}

func (h *AccountHandler) GetAddress(w http.ResponseWriter, r *http.Request) {
	acCode := chi.URLParam(r, "acCode")
	addr, err := h.repo.GetAddress(acCode)
	if err != nil {
		respondErr(w, http.StatusNotFound, "address not found")
		return
	}
	respond(w, http.StatusOK, addr)
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var acc models.AccountMaster
	if err := decode(r, &acc); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.repo.Create(&acc); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusCreated, acc)
}

func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	var acc models.AccountMaster
	if err := decode(r, &acc); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	acc.ACCode = chi.URLParam(r, "acCode")
	if err := h.repo.Update(&acc); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, acc)
}

func (h *AccountHandler) UpsertAddress(w http.ResponseWriter, r *http.Request) {
	var addr models.AddressMaster
	if err := decode(r, &addr); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	addr.ACCode = chi.URLParam(r, "acCode")
	if err := h.repo.UpsertAddress(&addr); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, addr)
}

func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	acCode := chi.URLParam(r, "acCode")
	if err := h.repo.Delete(acCode); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]string{"status": "deleted"})
}
