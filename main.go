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
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	fmt.Println("Starting application...")
	
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}
	
	var baseLogger log.Logger
	if environment == "production" {
		baseLogger = log.Production()
	} else {
		baseLogger = log.Development()
	}
	
	baseLogger.Info("Starting application with custom logger...")
	
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
	userRepo := repository.NewUserRepository(queries, repoLogger.With("repository", "User"))
	authRepo := repository.NewAuthRepository(queries, repoLogger.With("repository", "Auth"))
	contentRepo := repository.NewContentRepository(queries, repoLogger.With("repository", "Content"))
	analyticsRepo := repository.NewAnalyticsRepository(queries, repoLogger.With("repository", "Analytics"))
	
	// Initialize services with service-specific logger
	appLogger.Info("Initializing services...")
	userService := service.NewUserService(userRepo, serviceLogger.With("service", "User"))
	authService := service.NewAuthService(userRepo, authRepo, cfg.JWTSecret, cfg.GetTokenDuration(), 
		serviceLogger.With("service", "Auth"))
	contentService := service.NewContentService(contentRepo, userRepo, 
		serviceLogger.With("service", "Content"))
	analyticsService := service.NewAnalyticsService(analyticsRepo, contentRepo, userRepo, 
		serviceLogger.With("service", "Analytics"))
	
	// Initialize handlers with handler-specific logger
	appLogger.Info("Initializing handlers...")
	userHandler := user.NewHandler(userService, handlerLogger.With("handler", "User"))
	authHandler := auth.NewHandler(authService, handlerLogger.With("handler", "Auth"))
	contentHandler := content.NewHandler(contentService, handlerLogger.With("handler", "Content"))
	analyticsHandler := analytics.NewHandler(analyticsService, handlerLogger.With("handler", "Analytics"))
	
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