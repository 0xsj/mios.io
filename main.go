// main.go - Updated to include file service
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xsj/mios.io/api/analytics"
	"github.com/0xsj/mios.io/api/auth"
	"github.com/0xsj/mios.io/api/content"
	"github.com/0xsj/mios.io/api/file"
	"github.com/0xsj/mios.io/api/link_metadata"
	api "github.com/0xsj/mios.io/api/server"
	"github.com/0xsj/mios.io/api/user"
	"github.com/0xsj/mios.io/config"
	db "github.com/0xsj/mios.io/db/sqlc"
	"github.com/0xsj/mios.io/log"
	"github.com/0xsj/mios.io/middleware"
	"github.com/0xsj/mios.io/pkg/email"
	"github.com/0xsj/mios.io/pkg/redis"
	"github.com/0xsj/mios.io/pkg/storage"
	"github.com/0xsj/mios.io/repository"
	"github.com/0xsj/mios.io/service"
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
	appLogger := baseLogger.WithLayer("App")
	repoLogger := baseLogger.WithLayer("Repository")
	serviceLogger := baseLogger.WithLayer("Service")
	handlerLogger := baseLogger.WithLayer("Handler")
	middlewareLogger := baseLogger.WithLayer("Middleware")
	serverLogger := baseLogger.WithLayer("Server")
	redisLogger := baseLogger.WithLayer("redis")
	storageLogger := baseLogger.WithLayer("Storage")

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

	templateManager, err := email.NewTemplateManager("./pkg/email/templates")
	if err != nil {
		appLogger.Fatalf("Failed to initialize email template manager: %v", err)
	}

	baseURL := fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port)
	if cfg.Environment == "production" {
		baseURL = "https://appreciate.it"
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

	// Initialize storage
	appLogger.Info("Initializing storage...")
	var storageService storage.Storage
	
	switch cfg.StorageProvider {
	case "s3":
		appLogger.Info("Using S3 storage")
		s3Config := storage.S3Config{
			Region:      cfg.S3Region,
			Bucket:      cfg.S3Bucket,
			CDNDomain:   cfg.StorageCDNDomain,
			AccessKeyID: cfg.S3AccessKeyID,
			SecretKey:   cfg.S3SecretAccessKey,
		}
		storageService, err = storage.NewS3Storage(s3Config, storageLogger)
		if err != nil {
			appLogger.Fatalf("Failed to initialize S3 storage: %v", err)
		}
	case "local":
		appLogger.Info("Using local storage")
		storageService = storage.NewLocalStorage(cfg.StorageBasePath, cfg.StorageBaseURL, storageLogger)
		
		// Create uploads directory if it doesn't exist
		if err := os.MkdirAll(cfg.StorageBasePath, 0755); err != nil {
			appLogger.Fatalf("Failed to create uploads directory: %v", err)
		}
	default:
		appLogger.Fatalf("Unknown storage provider: %s", cfg.StorageProvider)
	}

	appLogger.Info("Initializing repositories...")
	userRepo := repository.NewUserRepository(queries, repoLogger.With("repository", "User"))
	authRepo := repository.NewAuthRepository(queries, repoLogger.With("repository", "Auth"))
	contentRepo := repository.NewContentRepository(queries, repoLogger.With("repository", "Content"))
	analyticsRepo := repository.NewAnalyticsRepository(queries, repoLogger.With("repository", "Analytics"))
	linkMetadataRepo := repository.NewLinkMetadataRepository(queries, repoLogger.With("repository", "LinkMetadata"))
	emailClient := email.NewEmailClient(baseLogger.WithLayer("Email"), templateManager)

	appLogger.Info("Initializing services...")
	userService := service.NewUserService(userRepo, serviceLogger.With("service", "User"))
	authService := service.NewAuthService(
		userRepo,
		authRepo,
		emailClient,
		cfg.JWTSecret,
		cfg.GetTokenDuration(),
		serviceLogger.With("service", "Auth"),
		baseURL,
	)
	contentService := service.NewContentService(contentRepo, userRepo,
		serviceLogger.With("service", "Content"))
	analyticsService := service.NewAnalyticsService(analyticsRepo, contentRepo, userRepo,
		serviceLogger.With("service", "Analytics"))
	linkMetadataService := service.NewLinkMetadataService(linkMetadataRepo,
		serviceLogger.With("service", "LinkMetadata"))
	
	// Initialize file service
	fileServiceConfig := service.FileServiceConfig{
		MaxFileSize:   cfg.MaxFileSize,
		MaxAvatarSize: cfg.MaxAvatarSize,
		AllowedImageTypes: []string{
			"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp",
		},
		AllowedVideoTypes: []string{
			"video/mp4", "video/webm", "video/ogg", "video/avi", "video/mov",
		},
		AllowedFileTypes: []string{
			"application/pdf", "text/plain", "application/json",
		},
		CDNDomain: cfg.StorageCDNDomain,
	}
	fileService := service.NewFileService(storageService, fileServiceConfig, serviceLogger.With("service", "File"))

	appLogger.Info("Initializing handlers...")
	userHandler := user.NewHandler(userService, handlerLogger.With("handler", "User"))
	authHandler := auth.NewHandler(authService, handlerLogger.With("handler", "Auth"))
	contentHandler := content.NewHandler(contentService, handlerLogger.With("handler", "Content"))
	analyticsHandler := analytics.NewHandler(analyticsService, handlerLogger.With("handler", "Analytics"))
	linkMetadataHandler := link_metadata.NewHandler(linkMetadataService, handlerLogger.With("handler", "LinkMetadata"))
	fileHandler := file.NewHandler(fileService, handlerLogger.With("handler", "File"))

	appLogger.Info("Initializing OpenAPI handler...")

	appLogger.Info("Setting up server...")
	server, err := api.NewServer(cfg, queries, serverLogger, redisClient)
	if err != nil {
		appLogger.Fatalf("Failed to initialize server: %v", err)
	}

	server.Router().Use(middleware.LoggingMiddleware(middlewareLogger))

	// Serve static files for local storage
	if cfg.StorageProvider == "local" {
		server.Router().Static("/uploads", cfg.StorageBasePath)
	}

	server.RegisterHandlers(userHandler, authHandler, contentHandler, authService, analyticsHandler, linkMetadataHandler, fileHandler)

	appLogger.Info("Registering OpenAPI handlers...")

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