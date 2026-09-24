package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	appconfig "github.com/subhranil002/GO-Cognito/internal/config"
	"github.com/subhranil002/GO-Cognito/internal/employee"
	"github.com/subhranil002/GO-Cognito/internal/middleware"
	"github.com/subhranil002/GO-Cognito/internal/router"
	"github.com/subhranil002/GO-Cognito/internal/user"
	"github.com/subhranil002/GO-Cognito/pkg/logger"
)

func main() {
	// Initialize structured logger
	log := logger.New()
	slog.SetDefault(log)

	// Load environment or Secrets Manager configuration
	cfg, err := appconfig.Load()
	if err != nil {
		log.Error("configuration error", "error", err)
		os.Exit(1)
	}

	userClient := user.NewClient(cfg.EMSAPIURL)
	userHandler := user.NewHandler()

	// Initialize JWT verifier using Cognito JWKS
	verifier := middleware.NewVerifier(
		cfg.CognitoIssuer(),
		cfg.CognitoClientID,
		userClient,
	)

	employeeClient := employee.NewClient(cfg.EMSAPIURL)
	employeeHandler := employee.NewHandler(employeeClient)

	// Register routes and global middleware
	mux := router.Setup(verifier, employeeHandler, userHandler)
	var handler http.Handler = mux
	handler = middleware.CORS(cfg.AllowedOrigin)(handler)
	handler = middleware.Logger(log)(handler)

	// Configure HTTP server timeouts
	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Start HTTP server listener
	log.Info("server started", "addr", cfg.HTTPAddr, "environment", cfg.AppEnv)
	if err = server.ListenAndServe(); err != nil {
		log.Error("server stopped unexpectedly", "error", err)
		os.Exit(1)
	}
}
