package main

import (
	"context"
	"fmt"
	"github.com/Lyr-a-Brode/tolmach/api-service/web"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	setupLogger()

	cfg, err := parseConfig()
	if err != nil {
		log.WithError(err).Fatal("unable to parse config")
	}

	if cfg.App.EnableDebug {
		log.SetLevel(log.DebugLevel)
	}

	router := web.NewRouter()

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: router,
	}

	runWithGracefulShutdown(&srv, cfg.App.ShutdownTimeout)
}

func setupLogger() {
	logger := log.StandardLogger()

	logger.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999",
	})

	logger.SetLevel(log.InfoLevel)
}

func runWithGracefulShutdown(srv *http.Server, timeout time.Duration) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("unable to start web server")
		}
	}()

	<-ctx.Done()

	stop()
	log.Info("shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatalf("unable to shutdown web server in %s", timeout)
	}

	log.Info("server stopped")
}
