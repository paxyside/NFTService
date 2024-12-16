package application

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"nft_service/infrastructure/config"
	"nft_service/infrastructure/database"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func StartApp() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	l = l.With("app_info", slog.GroupValue(
		slog.String("os", runtime.GOOS),
		slog.String("go_version", runtime.Version()),
		slog.Int("num_cpu", runtime.NumCPU()),
		slog.Int("num_goroutine", runtime.NumGoroutine()),
	))
	slog.SetDefault(l)

	cfg, err := config.LoadConfig()
	if err != nil {
		l.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	db, err := database.Init(cfg.DBURI)
	if err != nil {
		l.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	server, err := setupServer(ctx, db, cfg)
	if err != nil {
		l.Error("failed to setup server", slog.Any("error", err))
		os.Exit(1)
	}

	serverAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      server,
	}

	l.Info("starting server", slog.String("address", serverAddr))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Error("listen", "address", serverAddr, "error", err)
		}
	}()

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			l.Error("listen", "address", ":6060", "error", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	l.Info("shutting down server...")

	cancel()
	ctxShutdown, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		l.Error("server shutdown:", slog.Any("error", err))
	}

	<-ctxShutdown.Done()
	l.Info("server exiting")
}
