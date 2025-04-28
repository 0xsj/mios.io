package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xsj/gin-sqlc/api/analytics"
	authapi "github.com/0xsj/gin-sqlc/api/auth"
	"github.com/0xsj/gin-sqlc/api/content"
	api "github.com/0xsj/gin-sqlc/api/server"
	userapi "github.com/0xsj/gin-sqlc/api/user"
	"github.com/0xsj/gin-sqlc/config"
	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/0xsj/gin-sqlc/log"
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

	logger := &log.EmptyLogger{}

	fmt.Println("Loading configuration...")
	cfg := config.LoadConfig("dev", ".")
	fmt.Printf("Loaded configuration: %+v\n", cfg)

	if cfg.DBUsername == "" || cfg.DBPassword == "" || cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" {
		fmt.Println("ERROR: Database configuration values are missing")
		os.Exit(1)
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	fmt.Printf("Database URL: %s\n", dbURL)

	fmt.Println("Connecting to database...")
	dbpool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		logger.Fatalf("Unable to connect to database: %v", err)
	}

	fmt.Println("Testing database connection with ping...")
	err = dbpool.Ping(context.Background())
	if err != nil {
		fmt.Printf("Database ping failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Database connection successful!")

	defer dbpool.Close()

	fmt.Println("Initializing database queries...")
	queries := db.New(dbpool)

	fmt.Println("Initializing repositories...")
	userRepo := repository.NewUserRepository(queries)
	authRepo := repository.NewAuthRepository(queries)
	contentRepo := repository.NewContentRepository(queries)
	analyticsRepo := repository.NewAnalyticsRepository(queries)

	fmt.Println("Initializing services...")
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, authRepo, cfg.JWTSecret, cfg.GetTokenDuration())
	contentService := service.NewContentService(contentRepo, userRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo, contentRepo, userRepo)
	

	fmt.Println("Initializing handlers...")
	userHandler := userapi.NewHandler(userService)
	authHandler := authapi.NewHandler(authService)
	contentHandler := content.NewHandler(contentService)
	analyticsHandler := analytics.NewHandler(analyticsService)

	fmt.Println("Setting up server...")
	server := api.NewServer(cfg, queries, logger)
	
	server.RegisterHandlers(userHandler, authHandler, contentHandler, authService, analyticsHandler)

	fmt.Printf("Starting HTTP server on %s:%s...\n", cfg.Host, cfg.Port)
	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
		if err := server.Start(addr); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
			logger.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutdown signal received...")
	logger.Info("Shutting down server")


	fmt.Println("Server successfully shut down")
	logger.Info("Server exited properly")
}