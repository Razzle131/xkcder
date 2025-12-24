package main

import (
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"yadro.com/course/web/adapters/api"
	"yadro.com/course/web/adapters/rest"
	"yadro.com/course/web/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := mustMakeLogger(cfg.LogLevel)

	api := api.New(log, cfg.ApiAddress)

	mux := http.NewServeMux()
	mux.Handle("/images/logo2.png", rest.NewLogoHandler())
	mux.Handle("/", rest.NewIndexHandler(log, "./html/index.html"))
	mux.Handle("/search", rest.NewSearchHandler(log, "./html/search.html", api))
	mux.Handle("/search/random", rest.NewSearchRandomHandler("./html/search.html", api))
	mux.Handle("/admin", rest.NewAdminHandler(log, "./html/admin.html", cfg.TokenCookie, api))
	mux.Handle("/login", rest.NewLoginHandler(cfg.TokenCookie, api))
	mux.Handle("/update", rest.NewUpdateHandler(cfg.TokenCookie, api))
	mux.Handle("/drop", rest.NewDropHandler(cfg.TokenCookie, api))
	mux.Handle("/stats", rest.NewStatsHandler(api))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Info("Running HTTP server", "address", ":8080")
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error("server closed unexpectedly", "error", err)
		}
	}
}

func mustMakeLogger(logLevel string) *slog.Logger {
	var level slog.Level
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		slog.Error("failed parsing log level", "error", err)
		os.Exit(1)
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})

	return slog.New(handler)
}
