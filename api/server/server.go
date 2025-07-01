package api

import (
	"fmt"
	"time"

	"github.com/0xsj/mios.io/api/analytics"
	"github.com/0xsj/mios.io/api/auth"
	"github.com/0xsj/mios.io/api/content"
	"github.com/0xsj/mios.io/api/link_metadata" // Added import
	"github.com/0xsj/mios.io/api/user"
	"github.com/0xsj/mios.io/config"
	db "github.com/0xsj/mios.io/db/sqlc"
	"github.com/0xsj/mios.io/log"
	"github.com/0xsj/mios.io/middleware"
	"github.com/0xsj/mios.io/pkg/redis"
	"github.com/0xsj/mios.io/pkg/response"
	"github.com/0xsj/mios.io/service"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config      config.Config
	router      *gin.Engine
	store       db.Querier
	logger      log.Logger
	redisClient *redis.Client
}

func NewServer(config config.Config, store db.Querier, logger log.Logger, redisClient *redis.Client) (*Server, error) {
	router := gin.Default()

	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return nil, fmt.Errorf("failed to set trusted proxies: %w", err)
	}

	router.Use(middleware.RequestLogger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORSMiddleware())

	server := &Server{
		config:      config,
		router:      router,
		store:       store,
		logger:      logger,
		redisClient: redisClient,
	}

	logger.Info("API server initialized successfully")
	return server, nil
}

func (s *Server) RegisterHandlers(
	userHandler *user.Handler,
	authHandler *auth.Handler,
	contentHandler *content.Handler,
	authService service.AuthService,
	analyticsHandler *analytics.Handler,
	linkMetadataHandler *link_metadata.Handler, // Added parameter
) {
	s.logger.Info("Registering API routes")

	authMiddleware := middleware.AuthMiddleware(authService, s.logger)
	adminMiddleware := middleware.AdminMiddleware(s.logger)
	verifiedEmailMiddleware := middleware.RequireVerifiedEmail(authService, s.logger)

	publicRoutes := s.router.Group("/api")
	{
		// Auth routes
		authGroup := publicRoutes.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.RefreshToken)
			authGroup.POST("/forgot-password", authHandler.ForgotPassword)
			authGroup.POST("/reset-password", authHandler.ResetPassword)
			// authGroup.POST("/verify-email", authHandler.VerifyEmail)
		}

		// Public user routes
		publicUserGroup := publicRoutes.Group("/users")
		{
			publicUserGroup.GET("/username/:username", userHandler.GetUserByUsername)
			publicUserGroup.GET("/handle/:handle", userHandler.GetUserByHandle)
		}

		// Public content routes
		publicContentGroup := publicRoutes.Group("/content")
		{
			publicContentGroup.GET("/user/:user_id", contentHandler.GetUserContentItems)
		}

		// Public link metadata routes
		publicMetadataGroup := publicRoutes.Group("/link-metadata")
		{
			publicMetadataGroup.GET("/platforms", linkMetadataHandler.ListPlatforms)
			publicMetadataGroup.GET("/url", linkMetadataHandler.GetLinkMetadata)
		}
	}

	// Protected routes - require authentication
	protectedRoutes := s.router.Group("/api")
	protectedRoutes.Use(authMiddleware)
	{
		// Auth routes that require authentication
		authGroup := protectedRoutes.Group("/auth")
		{
			authGroup.POST("/logout", authHandler.Logout)
		}

		// User routes
		userGroup := protectedRoutes.Group("/users")
		{
			userGroup.GET("/:id", userHandler.GetUser)
			userGroup.PUT("/:id", userHandler.UpdateUser)
			userGroup.PATCH("/:id/handle", userHandler.UpdateHandle)
			userGroup.PATCH("/:id/onboarded", userHandler.UpdateOnboardedStatus)
			userGroup.DELETE("/:id", userHandler.DeleteUser)
		}

		// Content routes
		contentGroup := protectedRoutes.Group("/content")
		{
			// Some operations might need email verification
			verifiedContentGroup := contentGroup.Group("")
			verifiedContentGroup.Use(verifiedEmailMiddleware)
			{
				verifiedContentGroup.POST("", contentHandler.CreateContentItem)
				verifiedContentGroup.PUT("/:id", contentHandler.UpdateContentItem)
				verifiedContentGroup.PATCH("/:id/position", contentHandler.UpdateContentItemPosition)
				verifiedContentGroup.DELETE("/:id", contentHandler.DeleteContentItem)
			}

			// Some operations might not need email verification
			contentGroup.GET("/:id", contentHandler.GetContentItem)
		}

		// Analytics routes
		analyticsGroup := protectedRoutes.Group("/analytics")
		{
			analyticsGroup.POST("/clicks", analyticsHandler.RecordClick)
			analyticsGroup.POST("/page-views", analyticsHandler.RecordPageView)
			analyticsGroup.GET("/items/:id", analyticsHandler.GetContentItemAnalytics)
			analyticsGroup.POST("/items/:id/time-range", analyticsHandler.GetItemAnalyticsByTimeRange)
			analyticsGroup.GET("/users/:id", analyticsHandler.GetUserAnalytics)
			analyticsGroup.POST("/users/:id/time-range", analyticsHandler.GetUserAnalyticsByTimeRange)
			analyticsGroup.POST("/users/:id/page-views", analyticsHandler.GetProfilePageViewsByTimeRange)
			analyticsGroup.GET("/users/:id/dashboard", analyticsHandler.GetProfileDashboard)
			analyticsGroup.POST("/users/:id/referrers", analyticsHandler.GetReferrerAnalytics)
		}

		// Protected link metadata routes
		linkMetadataGroup := protectedRoutes.Group("/link-metadata")
		{
			linkMetadataGroup.POST("/fetch", linkMetadataHandler.FetchLinkMetadata)
		}
	}

	// Admin routes - require authentication and admin role
	adminRoutes := protectedRoutes.Group("/admin")
	adminRoutes.Use(adminMiddleware)
	{
		adminRoutes.PATCH("/users/:id/premium", userHandler.UpdatePremiumStatus)
		adminRoutes.PATCH("/users/:id/admin", userHandler.UpdateAdminStatus)
	}

	// Health check endpoint
	s.router.GET("/health", s.handleHealthCheck)

	s.logger.Info("API routes registered successfully")
}

// GetRedisClient returns the Redis client
func (s *Server) GetRedisClient() *redis.Client {
	return s.redisClient
}

// handleHealthCheck handles the health check endpoint
func (s *Server) handleHealthCheck(c *gin.Context) {
	s.logger.Debug("Health check endpoint called")

	// Gather system health information
	healthInfo := map[string]interface{}{
		"status":      "ok",
		"time":        time.Now().Format(time.RFC3339),
		"environment": s.config.Environment,
		"version":     s.config.Version,
	}

	response.Success(c, healthInfo, "Service is healthy")
}

// Start begins listening for HTTP requests on the specified address
func (s *Server) Start(addr string) error {
	s.logger.Infof("Starting API server on %s", addr)
	return s.router.Run(addr)
}

// Router returns the Gin engine for testing
func (s *Server) Router() *gin.Engine {
	return s.router
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() {
	s.logger.Info("Shutting down API server")
}
