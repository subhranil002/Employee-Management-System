package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"

	"github.com/subhranil002/GO-Cognito/internal/auth"
	appconfig "github.com/subhranil002/GO-Cognito/internal/config"
	"github.com/subhranil002/GO-Cognito/internal/employee"
	"github.com/subhranil002/GO-Cognito/internal/middleware"
	"github.com/subhranil002/GO-Cognito/internal/router"
	"github.com/subhranil002/GO-Cognito/pkg/logger"
)

func main() {
	// Initialize structured logger
	log := logger.New()
	slog.SetDefault(log)

	// Load configuration from environment
	cfg, err := appconfig.Load()
	if err != nil {
		log.Error("configuration error", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Load AWS configuration and SDK credentials
	awsCfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		log.Error("failed to load AWS configuration", "error", err)
		os.Exit(1)
	}

	cognitoClient := cognito.NewFromConfig(awsCfg)

	// Initialize auth service
	authService := auth.NewService(
		cognitoClient,
		cfg.CognitoUserPoolID,
		cfg.CognitoClientID,
		cfg.CognitoClientSecret,
	)

	// Wire token refresher into verifier for automatic token refresh
	tokenRefresher := middleware.NewTokenRefresher(
		cognitoClient,
		cfg.CognitoClientID,
		cfg.CognitoClientSecret,
	)

	// Initialize token verifier
	tokenVerifier, err := middleware.NewTokenVerifier(
		ctx,
		cfg.CognitoIssuer(),
		cfg.CognitoClientID,
		tokenRefresher,
	)
	if err != nil {
		log.Error("failed to create token verifier", "error", err)
		os.Exit(1)
	}

	authHandler := auth.NewHandler(authService)

	// Create HTTP client that proxies employee CRUD requests to GO-CRUD backend
	employeeClient := employee.NewClient(cfg.EmployeeAPIURL)
	employeeHandler := employee.NewHandler(employeeClient)

	// Register routes and setup middleware pipeline
	mux := router.Setup(authHandler, tokenVerifier, employeeHandler)

	var handler http.Handler = mux
	handler = middleware.CORS(cfg.AllowedOrigin)(handler)
	handler = middleware.Logger(log)(handler)

	// Configure HTTP server with timeouts
	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Info(
		"server started",
		"addr", cfg.HTTPAddr,
		"environment", cfg.AppEnv,
	)

	err = server.ListenAndServe()
	if err != nil {
		log.Error("server stopped unexpectedly", "error", err)
		os.Exit(1)
	}
}
