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
	"github.com/0xsj/gin-sqlc/repository"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/jackc/pgx/v4/pgxpool"
)

// func testPasswordVerification() {
// 	storedHash := "$2a$12$zS4uugKZD/axLQwjvSkGx.bIau3FX5UPox/digU9Quv9ujw9gpVDO"
// 	storedSalt := "lNj+J85F8862k7icgRKChQ=="
// 	plainPassword := "Password123!"
// 	err := password.VerifyPassword(plainPassword, storedHash, storedSalt)
// 	fmt.Printf("Password verification result: %v\n", err)
// }

func main() {
	// testPasswordVerification()
	fmt.Println("Starting application...")
	
	// Get environment
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}
	
	// Initialize base logger
	baseLogger := log.NewZapLogger(environment)
	baseLogger.Info("Starting application with ZapLogger...")
	
	// Create layer-specific loggers
	appLogger := baseLogger.WithLayer("App")
	repoLogger := baseLogger.WithLayer("Repository")
	serviceLogger := baseLogger.WithLayer("Service")
	handlerLogger := baseLogger.WithLayer("Handler")
	middlewareLogger := baseLogger.WithLayer("Middleware")
	serverLogger := baseLogger.WithLayer("Server")
	
	// Load configuration
	appLogger.Info("Loading configuration...")
	cfg := config.LoadConfig("dev", ".")
	appLogger.Debugf("Loaded configuration: %+v", cfg)
	
	if cfg.DBUsername == "" || cfg.DBPassword == "" || cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" {
		appLogger.Fatal("ERROR: Database configuration values are missing")
		return
	}
	
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	appLogger.Debugf("Database URL: %s", dbURL)
	
	// Connect to database
	appLogger.Info("Connecting to database...")
	dbpool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		appLogger.Fatalf("Database connection error: %v", err)
	}
	
	appLogger.Info("Testing database connection with ping...")
	err = dbpool.Ping(context.Background())
	if err != nil {
		appLogger.Fatalf("Database ping failed: %v", err)
	}
	
	appLogger.Info("Database connection successful!")
	defer dbpool.Close()
	
	// Initialize database and components
	appLogger.Info("Initializing database queries...")
	queries := db.New(dbpool)
	
	// Initialize repositories with repository-specific logger
	appLogger.Info("Initializing repositories...")
	userRepo := repository.NewUserRepository(queries, repoLogger.WithField("repository", "User"))
	authRepo := repository.NewAuthRepository(queries, repoLogger.WithField("repository", "Auth"))
	contentRepo := repository.NewContentRepository(queries, repoLogger.WithField("repository", "Content"))
	analyticsRepo := repository.NewAnalyticsRepository(queries, repoLogger.WithField("repository", "Analytics"))
	
	// Initialize services with service-specific logger
	appLogger.Info("Initializing services...")
	userService := service.NewUserService(userRepo, serviceLogger.WithField("service", "User"))
	authService := service.NewAuthService(userRepo, authRepo, cfg.JWTSecret, cfg.GetTokenDuration(), 
		serviceLogger.WithField("service", "Auth"))
	contentService := service.NewContentService(contentRepo, userRepo, 
		serviceLogger.WithField("service", "Content"))
	analyticsService := service.NewAnalyticsService(analyticsRepo, contentRepo, userRepo, 
		serviceLogger.WithField("service", "Analytics"))
	
	// Initialize handlers with handler-specific logger
	appLogger.Info("Initializing handlers...")
	userHandler := user.NewHandler(userService, handlerLogger.WithField("handler", "User"))
	authHandler := auth.NewHandler(authService, handlerLogger.WithField("handler", "Auth"))
	contentHandler := content.NewHandler(contentService, handlerLogger.WithField("handler", "Content"))
	analyticsHandler := analytics.NewHandler(analyticsService, handlerLogger.WithField("handler", "Analytics"))
	
	// Setup server
	appLogger.Info("Setting up server...")
	server := api.NewServer(cfg, queries, serverLogger)
	
	// Apply middleware
	server.Router().Use(middleware.LoggingMiddleware(middlewareLogger))
	
	// Register routes
	server.RegisterHandlers(userHandler, authHandler, contentHandler, authService, analyticsHandler)
	
	// Start server
	appLogger.Infof("Starting HTTP server on %s:%s...", cfg.Host, cfg.Port)
	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
		if err := server.Start(addr); err != nil && err != http.ErrServerClosed {
			appLogger.Fatalf("Server error: %v", err)
		}
	}()
	
	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	appLogger.Info("Shutdown signal received...")
	appLogger.Info("Server successfully shut down")
}