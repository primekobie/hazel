package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/primekobie/hazel/handlers"
	"github.com/primekobie/hazel/mail"
	"github.com/primekobie/hazel/postgres"
	"github.com/primekobie/hazel/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

//	@title			Hazel Project Management API
//	@version		1.0
//	@description	This is the backend API for the Hazel project management application.
//	@contact.name	API Support
//	@contact.url	https://github.com/primekobie/hazel
//	@contact.email	support@hazel.local
func main() {

	_ = godotenv.Load()

	setupLogging()

	cfg := loadConfig()

	db, err := pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	mailer := mail.NewMailer(cfg.MailConfig)
	userService := services.NewUserService(postgres.NewUserStore(db), mailer)
	workspaceService := services.NewWorkspaceService(postgres.NewWorkspaceStore(db))

	handler := handlers.NewHandler(userService, workspaceService)

	app := newApplication(handler, cfg.ServerAddress)

	// Graceful shutdown setup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("Starting server")
		serverErr <- app.start()
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			panic(err)
		}
	case sig := <-stop:
		slog.Info("Shutting down server", "signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.shutdown(ctx); err != nil {
			slog.Error("Graceful shutdown failed", "error", err)
		}
	}
}
