package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/val/inventory/internal/models"
	"github.com/val/inventory/internal/repository"
)

type IssueHandler struct {
	repo *repository.IssueRepo
}

type issueRequest struct {
	models.Issue
	Details []models.IssueDetail `json:"details"`
}

func (h *IssueHandler) List(w http.ResponseWriter, r *http.Request) {
	page, perPage := pageParams(r)
	acCode := r.URL.Query().Get("ac_code")
	tag := r.URL.Query().Get("tag")

	issues, total, err := h.repo.List(page, perPage, acCode, tag)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, models.PagedResult[models.Issue]{
		Data: issues, Total: total, Page: page, PerPage: perPage,
	})
}

func (h *IssueHandler) Get(w http.ResponseWriter, r *http.Request) {
	tNo := chi.URLParam(r, "tNo")
	issue, err := h.repo.Get(tNo)
	if err != nil {
		respondErr(w, http.StatusNotFound, "issue not found")
		return
	}
	respond(w, http.StatusOK, issue)
}

func (h *IssueHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req issueRequest
	if err := decode(r, &req); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.repo.Create(&req.Issue, req.Details); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusCreated, req.Issue)
}

func (h *IssueHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req issueRequest
	if err := decode(r, &req); err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	req.Issue.TNo = chi.URLParam(r, "tNo")
	if err := h.repo.Update(&req.Issue, req.Details); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, req.Issue)
}

func (h *IssueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tNo := chi.URLParam(r, "tNo")
	if err := h.repo.Delete(tNo); err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]string{"status": "deleted"})
}
