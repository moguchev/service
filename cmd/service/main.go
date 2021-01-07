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

	emp_uc "github.com/moguchev/service/internal/employees/usecase"

	delivery "github.com/moguchev/service/internal/employees/delivery/http"

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
	// !errors.Is(err, migrate.ErrNoChange) do not work(
	if err != nil && strings.Compare(err.Error(), migrate.ErrNoChange.Error()) != 0 {
		log.WithError(err).Fatal("migrate db")
	}

	// Create Repository level
	// Create Usecase level
	employeeUC := emp_uc.NewEmployeesUsecase()

	// Create Router
	router := mux.NewRouter()

	// Set Middlewares
	mw := middleware.InitMiddleware(log)
	router.Use(mw.RecoverMiddleware)
	router.Use(mw.CORSMiddleware)

	// Set Handlers
	base := router.PathPrefix(cfg.Server.APIBasePath).Subrouter()

	delivery.SetEmployeesHandler(base, employeeUC)

	// Make Server
	bctx := logger.WithLogger(context.Background(), log)
	srv := http.Server{
		Handler: router,
		Addr:    cfg.Server.Address,
		BaseContext: func(net.Listener) context.Context {
			return bctx
		},
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		err = srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	group.Go(func() error {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop,
			syscall.SIGINT,
			syscall.SIGTERM,
		)
		<-stop
		log.Info("stop service")
		cancel()
		err = srv.Shutdown(context.Background())
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
