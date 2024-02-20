package main

import (
	"fmt"
	"github.com/twibber/core/app/routes"
	"github.com/twibber/core/cfg"
	"log/slog"
)

// main is the entry point for the application
func main() {
	// Log the server start
	slog.With("port", cfg.Config.Port, "debug", cfg.Config.Debug).Info("starting server")

	// Configure the routes and start the server
	if err := routes.Configure().Listen(fmt.Sprintf("%s:%s", "0.0.0.0", cfg.Config.Port)); err != nil {
		// if the server fails to start, panic with the error
		panic(err)
	}
}
