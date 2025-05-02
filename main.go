package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xsj/gin-sqlc/api/analytics"
	"github.com/0xsj/gin-sqlc/api/auth"
	"github.com/0xsj/gin-sqlc/api/content"
	api "github.com/0xsj/gin-sqlc/api/server"
	"github.com/0xsj/gin-sqlc/api/user"
	"github.com/0xsj/gin-sqlc/config"
	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/middleware"
	"github.com/0xsj/gin-sqlc/pkg/password"
	"github.com/0xsj/gin-sqlc/repository"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/jackc/pgx/v4/pgxpool"
)

func testPasswordVerification() {
	storedHash := "$2a$12$zS4uugKZD/axLQwjvSkGx.bIau3FX5UPox/digU9Quv9ujw9gpVDO"
	storedSalt := "lNj+J85F8862k7icgRKChQ=="
	plainPassword := "Password123!"
	
	err := password.VerifyPassword(plainPassword, storedHash, storedSalt)
	fmt.Printf("Password verification result: %v\n", err)
}

func main() {
	testPasswordVerification()
	fmt.Println("Starting application...")

	// Get environment
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	// Initialize logger
	logger := log.NewZapLogger(environment)
	logger.Info("Starting application with ZapLogger...")

	// Load configuration
	fmt.Println("Loading configuration...")
	cfg := config.LoadConfig("dev", ".")
	logger.Debugf("Loaded configuration: %+v", cfg)

	if cfg.DBUsername == "" || cfg.DBPassword == "" || cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" {
		logger.Fatal("ERROR: Database configuration values are missing")
		return
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	logger.Debugf("Database URL: %s", dbURL)

	logger.Info("Connecting to database...")
	dbpool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		logger.Fatalf("Database connection error: %v", err)
	}

	logger.Info("Testing database connection with ping...")
	err = dbpool.Ping(context.Background())
	if err != nil {
		logger.Fatalf("Database ping failed: %v", err)
	}
	logger.Info("Database connection successful!")

	defer dbpool.Close()

	logger.Info("Initializing database queries...")
	queries := db.New(dbpool)

	logger.Info("Initializing repositories...")
	userRepo := repository.NewUserRepository(queries)
	authRepo := repository.NewAuthRepository(queries)
	contentRepo := repository.NewContentRepository(queries)
	analyticsRepo := repository.NewAnalyticsRepository(queries)

	logger.Info("Initializing services...")
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, authRepo, cfg.JWTSecret, cfg.GetTokenDuration())
	contentService := service.NewContentService(contentRepo, userRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo, contentRepo, userRepo)
	
	logger.Info("Initializing handlers...")
	userHandler := user.NewHandler(userService)
	authHandler := auth.NewHandler(authService)
	contentHandler := content.NewHandler(contentService)
	analyticsHandler := analytics.NewHandler(analyticsService)

	logger.Info("Setting up server...")
	server := api.NewServer(cfg, queries, logger)
	
	// Add middleware to the router
	server.Router().Use(middleware.LoggingMiddleware(logger))
	
	// Register handlers
	server.RegisterHandlers(userHandler, authHandler, contentHandler, authService, analyticsHandler)

	logger.Infof("Starting HTTP server on %s:%s...", cfg.Host, cfg.Port)
	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
		if err := server.Start(addr); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutdown signal received...")
	logger.Info("Server successfully shut down")
}