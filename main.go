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
	"github.com/0xsj/gin-sqlc/api/openapi"
	api "github.com/0xsj/gin-sqlc/api/server"
	"github.com/0xsj/gin-sqlc/api/user"
	"github.com/0xsj/gin-sqlc/config"
	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/middleware"
	"github.com/0xsj/gin-sqlc/pkg/redis"
	"github.com/0xsj/gin-sqlc/repository"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	openapiMiddleware "github.com/oapi-codegen/gin-middleware"
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
	appLogger := baseLogger.WithLayer("App")
	repoLogger := baseLogger.WithLayer("Repository")
	serviceLogger := baseLogger.WithLayer("Service")
	handlerLogger := baseLogger.WithLayer("Handler")
	middlewareLogger := baseLogger.WithLayer("Middleware")
	serverLogger := baseLogger.WithLayer("Server")
	openapiLogger := baseLogger.WithLayer("OpenAPI")
	redisLogger := baseLogger.WithLayer("redis")
	
	appLogger.Info("Loading configuration...")
	cfg := config.LoadConfig("dev", ".")
	appLogger.Debugf("Loaded configuration: %+v", cfg)
	redisClient, err := redis.NewClient(cfg, redisLogger)
	if err != nil {
		appLogger.Fatalf("Failed to connect to Redis: %v", err)
	}

	defer redisClient.Close()

	
	if cfg.DBUsername == "" || cfg.DBPassword == "" || cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" {
		appLogger.Fatal("ERROR: Database configuration values are missing")
		return
	}
	
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	appLogger.Debugf("Database URL: %s", dbURL)
	
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
	
	appLogger.Info("Initializing database queries...")
	queries := db.New(dbpool)
	
	appLogger.Info("Initializing repositories...")
	userRepo := repository.NewUserRepository(queries, repoLogger.With("repository", "User"))
	authRepo := repository.NewAuthRepository(queries, repoLogger.With("repository", "Auth"))
	contentRepo := repository.NewContentRepository(queries, repoLogger.With("repository", "Content"))
	analyticsRepo := repository.NewAnalyticsRepository(queries, repoLogger.With("repository", "Analytics"))
	
	appLogger.Info("Initializing services...")
	userService := service.NewUserService(userRepo, serviceLogger.With("service", "User"))
	authService := service.NewAuthService(userRepo, authRepo, cfg.JWTSecret, cfg.GetTokenDuration(), 
		serviceLogger.With("service", "Auth"))
	contentService := service.NewContentService(contentRepo, userRepo, 
		serviceLogger.With("service", "Content"))
	analyticsService := service.NewAnalyticsService(analyticsRepo, contentRepo, userRepo, 
		serviceLogger.With("service", "Analytics"))
	
	appLogger.Info("Initializing handlers...")
	userHandler := user.NewHandler(userService, handlerLogger.With("handler", "User"))
	authHandler := auth.NewHandler(authService, handlerLogger.With("handler", "Auth"))
	contentHandler := content.NewHandler(contentService, handlerLogger.With("handler", "Content"))
	analyticsHandler := analytics.NewHandler(analyticsService, handlerLogger.With("handler", "Analytics"))
	
	appLogger.Info("Initializing OpenAPI handler...")
	openapiHandler := openapi.NewHandler(
		authService,  
		userService,  
		contentService, 
		analyticsService,
		openapiLogger,
	)
	
	appLogger.Info("Setting up server...")
	server := api.NewServer(cfg, queries, serverLogger, redisClient)
	
	server.Router().Use(middleware.LoggingMiddleware(middlewareLogger))
	
	swagger, err := openapi.GetSwagger()
	if err == nil {
		swagger.Servers = nil
		
		appLogger.Info("Adding OpenAPI validation middleware...")
		server.Router().Use(openapiMiddleware.OapiRequestValidator(swagger))
	} else {
		appLogger.Errorf("Error loading OpenAPI spec: %v", err)
	}
	
	server.RegisterHandlers(userHandler, authHandler, contentHandler, authService, analyticsHandler)
	
	appLogger.Info("Registering OpenAPI handlers...")
	openapi.RegisterOpenAPIHandlers(server.Router(), openapiHandler) 
	
	appLogger.Infof("Starting HTTP server on %s:%s...", cfg.Host, cfg.Port)
	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
		if err := server.Start(addr); err != nil && err != http.ErrServerClosed {
			appLogger.Fatalf("Server error: %v", err)
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	appLogger.Info("Shutdown signal received...")
	appLogger.Info("Server successfully shut down")
}