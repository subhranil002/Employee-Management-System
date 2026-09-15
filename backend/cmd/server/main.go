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
	"github.com/subhranil002/GO-Cognito/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.New()
	slog.SetDefault(log)

	// Load configuration
	cfg, err := appconfig.Load()
	if err != nil {
		log.Error("configuration error", "error", err)
		os.Exit(1)
	}

	// Initialize verifier and employee handlers
	tokenVerifier := middleware.NewTokenVerifier(
		cfg.CognitoIssuer(),
		cfg.CognitoClientID,
	)
	employeeClient := employee.NewClient(cfg.EmployeeAPIURL)
	employeeHandler := employee.NewHandler(employeeClient)

	// Setup routes and middleware
	mux := router.Setup(tokenVerifier, employeeHandler)
	var handler http.Handler = mux
	handler = middleware.CORS(cfg.AllowedOrigin)(handler)
	handler = middleware.Logger(log)(handler)

	// Start HTTP server
	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Info("server started", "addr", cfg.HTTPAddr, "environment", cfg.AppEnv)
	if err = server.ListenAndServe(); err != nil {
		log.Error("server stopped unexpectedly", "error", err)
		os.Exit(1)
	}
}
