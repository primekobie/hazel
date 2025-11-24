package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/primekobie/hazel/handlers"
)

type application struct {
	handler *handlers.Handler
	server  *http.Server
}

func newApplication(handler *handlers.Handler, address string) *application {
	server := http.Server{
		Addr: fmt.Sprintf(":%s", address),
	}

	return &application{
		handler: handler,
		server:  &server,
	}
}

func (app *application) start() error {

	app.server.Handler = app.routes()

	return app.server.ListenAndServe()
}

func (app *application) shutdown(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}
