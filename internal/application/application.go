package application

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"nft_service/infrastructure/config"
	"nft_service/infrastructure/database"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func StartApp() {
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
	l.Info("connected to database")
	defer db.Close()

	server, err := setupServer(db, cfg)
	if err != nil {
		l.Error("failed to setup server", slog.Any("error", err))
		os.Exit(1)
	}

	serverAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	srv := &http.Server{
		Addr:         serverAddr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      server,
	}

	l.Info("starting server", slog.String("address", serverAddr))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Error("listen", "address", serverAddr, "error", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case d := <-quit:
		l.Info("servers stopped", "reason", d.String())
	}

	l.Info("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		l.Error("server shutdown:", slog.Any("error", err))
	}

	select {
	case <-ctx.Done():
		l.Warn("server shutdown: timeout of 5 seconds.")
	}

	l.Info("server exiting")

	os.Exit(0)
}
