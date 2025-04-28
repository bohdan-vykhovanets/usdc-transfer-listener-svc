package service

import (
	"context"
	er "errors"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/data/postgres"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/service/background"
	"github.com/ethereum/go-ethereum/ethclient"
	"net"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/config"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type service struct {
	log      *logan.Entry
	copus    types.Copus
	listener net.Listener
}

func (s *service) run(cfg config.Config, appCtx context.Context) error {
	r := s.router(cfg)
	server := &http.Server{Handler: r}

	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	go func() {
		<-appCtx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*5)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.log.Warn("Failed to shutdown http server gracefully")
		} else {
			s.log.Warn("Http server shutdown gracefully")
		}
	}()

	s.log.Info("Service started")
	err := server.Serve(s.listener)
	if err != nil && !er.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "Server failed")
	}
	s.log.Info("Serve function finished")
	return nil
}

func newService(cfg config.Config) *service {
	return &service{
		log:      cfg.Log(),
		copus:    cfg.Copus(),
		listener: cfg.Listener(),
	}
}

func Run(cfg config.Config) {
	appCtx, stopSignalListener := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignalListener()

	var wg sync.WaitGroup
	logger := cfg.Log()

	var bckCtx context.Context = appCtx
	bckCtx = background.CtxLog(logger)(bckCtx)
	bckCtx = background.CtxDb(postgres.NewMainQ(cfg.DB()))(bckCtx)

	svc := newService(cfg)

	node := cfg.Node()
	mainnetSocket := node.GetNodeUrl()
	client, err := ethclient.Dial(mainnetSocket)
	if err != nil {
		logger.WithError(err).Fatal("Failed to dial Ethereum client")
	}

	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.WithField("panic", r).Error("panic recovered")
			}
		}()
		background.Transfers(bckCtx, client, &wg)
	}()

	runErr := svc.run(cfg, appCtx)

	logger.Info("Waiting for background job to finish")
	waitCtx, waitCancel := context.WithTimeout(context.Background(), time.Second*15)
	defer waitCancel()

	waitChan := make(chan struct{})
	go func() {
		defer close(waitChan)
		wg.Wait()
	}()

	select {
	case <-waitChan:
		logger.Info("Background job finished")
	case <-waitCtx.Done():
		logger.Info("Timed out waiting for background job to finish")
	}

	if runErr != nil {
		logger.WithError(runErr).Fatal("Service run failed")
	}
}
