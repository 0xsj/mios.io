package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authapi "github.com/0xsj/gin-sqlc/api/auth"
	userapi "github.com/0xsj/gin-sqlc/api/user"
	"github.com/0xsj/gin-sqlc/config"
	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/repository"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	// Print startup message
	fmt.Println("Starting application...")

	logger := &log.EmptyLogger{}

	// Load configuration and print debug info
	fmt.Println("Loading configuration...")
	cfg := config.LoadConfig("dev", ".")
	fmt.Printf("Loaded configuration: %+v\n", cfg)

	// Check for empty configuration values
	if cfg.DBUsername == "" || cfg.DBPassword == "" || cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" {
		fmt.Println("ERROR: Database configuration values are missing")
		os.Exit(1)
	}

	// Build database URL and print it for debugging
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	fmt.Printf("Database URL: %s\n", dbURL)

	// Connect to database
	fmt.Println("Connecting to database...")
	dbpool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		logger.Fatalf("Unable to connect to database: %v", err)
	}

	// Test database connection with ping
	fmt.Println("Testing database connection with ping...")
	err = dbpool.Ping(context.Background())
	if err != nil {
		fmt.Printf("Database ping failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Database connection successful!")

	defer dbpool.Close()

	// Initialize queries
	fmt.Println("Initializing database queries...")
	queries := db.New(dbpool)

	// Initialize repository
	fmt.Println("Initializing user repository...")
	userRepo := repository.NewUserRepository(queries)
	authRepo := repository.NewAuthRepository(queries)

	// Initialize service
	fmt.Println("Initializing user service...")
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, authRepo, cfg.JWTSecret, cfg.GetTokenDuration())

	// Initialize handler
	fmt.Println("Initializing user handler...")
	userHandler := userapi.NewHandler(userService)
	authHandler := authapi.NewHandler(authService)

	// Setup router
	fmt.Println("Setting up Gin router...")
	router := gin.Default()

	// Register routes
	fmt.Println("Registering API routes...")
	userHandler.RegisterRoutes(router)
	authHandler.RegisterRoutes(router)

	// Print registered routes for debugging
	for _, route := range router.Routes() {
		fmt.Printf("Registered route: %s %s\n", route.Method, route.Path)
	}

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: router,
	}

	// Start server
	fmt.Printf("Starting HTTP server on %s:%s...\n", cfg.Host, cfg.Port)
	go func() {
		logger.Infof("Starting server on %s:%s", cfg.Host, cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
			logger.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutdown signal received...")
	logger.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown error: %v\n", err)
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server successfully shut down")
	logger.Info("Server exited properly")
}
