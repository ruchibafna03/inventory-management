package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/val/inventory/internal/repository"
)

func NewRouter(db *sqlx.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
	}))

	// Repos
	itemRepo     := repository.NewItemRepo(db)
	issueRepo    := repository.NewIssueRepo(db)
	receiptRepo  := repository.NewReceiptRepo(db)
	saleRepo     := repository.NewSaleRepo(db)
	purchRepo    := repository.NewPurchaseRepo(db)
	accountRepo  := repository.NewAccountRepo(db)
	rateRepo     := repository.NewRateRepo(db)
	lotRepo      := repository.NewLotRepo(db)
	utilRepo     := repository.NewUtilityRepo(db)

	// Handlers
	ih  := &ItemHandler{repo: itemRepo}
	ish := &IssueHandler{repo: issueRepo}
	rh  := &ReceiptHandler{repo: receiptRepo}
	sh  := &SaleHandler{repo: saleRepo}
	ph  := &PurchaseHandler{repo: purchRepo}
	ah  := &AccountHandler{repo: accountRepo}
	rah := &RateHandler{repo: rateRepo}
	lh  := &LotHandler{repo: lotRepo}
	uh  := &UtilityHandler{repo: utilRepo}

	r.Route("/api/v1", func(r chi.Router) {
		// Items
		r.Route("/items", func(r chi.Router) {
			r.Get("/", ih.List)
			r.Post("/", ih.Create)
			r.Get("/tags", ih.Tags)
			r.Route("/{itcd}", func(r chi.Router) {
				r.Get("/", ih.Get)
				r.Put("/", ih.Update)
				r.Delete("/", ih.Delete)
			})
		})

		// Issues
		r.Route("/issues", func(r chi.Router) {
			r.Get("/", ish.List)
			r.Post("/", ish.Create)
			r.Route("/{tNo}", func(r chi.Router) {
				r.Get("/", ish.Get)
				r.Put("/", ish.Update)
				r.Delete("/", ish.Delete)
			})
		})

		// Receipts
		r.Route("/receipts", func(r chi.Router) {
			r.Get("/", rh.List)
			r.Post("/", rh.Create)
			r.Route("/{tNo}", func(r chi.Router) {
				r.Get("/", rh.Get)
				r.Put("/", rh.Update)
				r.Delete("/", rh.Delete)
			})
		})

		// Sales
		r.Route("/sales", func(r chi.Router) {
			r.Get("/", sh.List)
			r.Post("/", sh.Create)
			r.Route("/{vouchNo}", func(r chi.Router) {
				r.Get("/", sh.Get)
				r.Put("/", sh.Update)
				r.Delete("/", sh.Delete)
			})
		})

		// Purchases
		r.Route("/purchases", func(r chi.Router) {
			r.Get("/", ph.List)
			r.Post("/", ph.Create)
			r.Route("/{vouchNo}", func(r chi.Router) {
				r.Get("/", ph.Get)
				r.Put("/", ph.Update)
				r.Delete("/", ph.Delete)
			})
		})

		// Accounts
		r.Route("/accounts", func(r chi.Router) {
			r.Get("/", ah.List)
			r.Post("/", ah.Create)
			r.Route("/{acCode}", func(r chi.Router) {
				r.Get("/", ah.Get)
				r.Put("/", ah.Update)
				r.Delete("/", ah.Delete)
				r.Get("/address", ah.GetAddress)
				r.Put("/address", ah.UpsertAddress)
			})
		})
		r.Route("/groups", func(r chi.Router) {
			r.Get("/", ah.ListGroups)
			r.Post("/", ah.CreateGroup)
			r.Put("/{gCode}", ah.UpdateGroup)
			r.Delete("/{gCode}", ah.DeleteGroup)
		})

		// Rates
		r.Route("/rates", func(r chi.Router) {
			r.Get("/", rah.List)
			r.Get("/latest", rah.Latest)
			r.Post("/", rah.Create)
			r.Delete("/{id}", rah.Delete)
		})

		// Lots
		r.Route("/lots", func(r chi.Router) {
			r.Get("/", lh.List)
			r.Post("/", lh.Create)
			r.Route("/{tNo}", func(r chi.Router) {
				r.Get("/", lh.Get)
				r.Put("/", lh.Update)
				r.Delete("/", lh.Delete)
			})
		})

		// Utilities
		r.Route("/utilities", func(r chi.Router) {
			r.Get("/summary", uh.DataSummary)
			r.Get("/orphan-issues", uh.OrphanIssueDetails)
			r.Get("/orphan-receipts", uh.OrphanReceiptDetails)
			r.Get("/stock-position", uh.StockPosition)
			r.Get("/item-history/{itcd}", uh.ItemHistory)
			r.Post("/change-item-code", uh.ChangeItemCode)
			r.Get("/code-change-history", uh.CodeChangeHistory)
			r.Get("/blocked-items", uh.BlockedItems)
			r.Post("/blocked-items/{itcd}", uh.BlockItem)
			r.Delete("/blocked-items/{itcd}", uh.UnblockItem)
			r.Get("/users", uh.ListUsers)
			r.Put("/password", uh.ChangePassword)
		})
	})

	return r
}
