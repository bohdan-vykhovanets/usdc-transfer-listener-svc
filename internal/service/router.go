package service

import (
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/config"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/data/postgres"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/service/background"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/service/handlers"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"net/http"
)

func (s *service) router(cfg config.Config) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxDb(postgres.NewMainQ(cfg.DB())),
			handlers.CtxNode(cfg.Node().GetNodeUrl()),
			background.CtxLog(s.log),
			background.CtxDb(postgres.NewMainQ(cfg.DB())),
		),
	)
	r.Route("/integrations/usdc-transfer-listener-svc", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		r.Get("/transfers", handlers.IndexTransfers)
		r.Post("/transfers/send-tx", handlers.SendTransactionTest)
		r.Post("/transfers/send-token", handlers.SendTokenTest)
	})

	return r
}
