package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	delivery "github.com/moguchev/service/internal/delivery/http"
	repo "github.com/moguchev/service/internal/repository"
	uc "github.com/moguchev/service/internal/usecase"

	"github.com/golang-migrate/migrate"
	"github.com/gorilla/mux"
	"github.com/moguchev/service/config"
	"github.com/moguchev/service/migration"
	"github.com/moguchev/service/pkg/logger"
	"github.com/moguchev/service/pkg/middleware"
	"github.com/moguchev/service/pkg/pgsql"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var log = logrus.New()

func main() {
	configPath := flag.String("c", "config.yaml", "set config path")
	flag.Parse()

	if *configPath == "" {
		log.WithError(fmt.Errorf("config path in blank")).Fatal("find config")
	}

	cfg, err := config.GetConfig(*configPath)
	if err != nil {
		log.WithError(err).Fatal("create config")
	}

	// Create logger
	log = cfg.Log.CreateLogger()

	// Create ctx
	ctx, cancel := context.WithCancel(context.Background())
	ctx = logger.WithLogger(ctx, log)

	// Create DB
	db, err := cfg.DB.CreateDB()
	if err != nil {
		log.WithError(err).Fatal("init db")
	}

	// Migrate DB
	err = pgsql.EnsureDB(db, migration.Assets)
	// !errors.Is(err, migrate.ErrNoChange) do not work
	if err != nil && strings.Compare(err.Error(), migrate.ErrNoChange.Error()) != 0 {
		log.WithError(err).Fatal("migrate db")
	}

	// Create Repository level
	empRepo := repo.NewEmployeesRepository(db)
	// Create Usecase level
	empUC := uc.NewEmployeesUsecase(empRepo)

	// Create Router
	router := mux.NewRouter()

	// Set Middlewares
	mw := middleware.InitMiddleware(log)
	router.Use(mw.RecoverMiddleware)
	router.Use(mw.CORSMiddleware)

	// Set Handlers
	base := router.PathPrefix(cfg.Server.APIBasePath).Subrouter()

	delivery.SetEmployeesHandler(base, empUC)

	// Make Server
	bctx := logger.WithLogger(context.Background(), log)
	srv := http.Server{
		Handler: router,
		Addr:    cfg.Server.Address,
		BaseContext: func(net.Listener) context.Context {
			return bctx
		},
	}

	group, gctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		err = srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	group.Go(func() error {
		sgnl := make(chan os.Signal, 1)
		signal.Notify(sgnl,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)
		stop := <-sgnl
		log.WithField("signal", stop).Info("waiting for all processes to stop")

		cancel()
		err = srv.Shutdown(gctx)
		if err != nil {
			return err
		}
		return nil
	})

	log.Infof("service started at %s", cfg.Server.Address)

	if err = group.Wait(); err != nil {
		log.WithError(err).Fatal()
	}
}
