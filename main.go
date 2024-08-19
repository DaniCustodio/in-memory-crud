package main

import (
	"log/slog"
	"main/api"
	"net/http"
	"time"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to run the code", "error", err)
	}
	slog.Info("code ran successfully")
}

func run() error {
	handler := api.NewHandler()

	server := http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
