package api

import (
	"net/http"
	"time"

	"github.com/0xsj/gin-sqlc/api/auth"
	"github.com/0xsj/gin-sqlc/api/content"
	"github.com/0xsj/gin-sqlc/api/user"
	"github.com/0xsj/gin-sqlc/config"
	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/middleware"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config config.Config
	router *gin.Engine
	store  db.Querier
	log    log.Logger
	// server *http.Server
}
func NewServer(config config.Config, store db.Querier, log log.Logger) *Server {
	router := gin.New()
	
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	server := &Server{
		config: config,
		router: router,
		store:  store,
		log:    log,
	}
	
	return server
}

func (s *Server) RegisterHandlers(
	userHandler *user.Handler,
	authHandler *auth.Handler,
	contentHandler *content.Handler,
	authService service.AuthService,
) {
	// Create middleware instances
	authMiddleware := middleware.AuthMiddleware(authService)
	adminMiddleware := middleware.AdminMiddleware()
	
	// Public routes
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
			authGroup.POST("/verify-email", authHandler.VerifyEmail)
		}
		
		// Public user routes
		publicUserGroup := publicRoutes.Group("/users")
		{
			publicUserGroup.GET("/username/:username", userHandler.GetUserByUsername)
			publicUserGroup.GET("/handle/:handle", userHandler.GetUserbyHandle)
		}
		
		// Public content routes (for viewing only)
		publicContentGroup := publicRoutes.Group("/content")
		{
			publicContentGroup.GET("/user/:user_id", contentHandler.GetUserContentItems)
		}
	}
	
	// Protected routes (authentication required)
	protectedRoutes := s.router.Group("/api")
	protectedRoutes.Use(authMiddleware)
	{
		// Auth routes that require authentication
		authGroup := protectedRoutes.Group("/auth")
		{
			authGroup.POST("/logout", authHandler.Logout)
		}
		
		// User routes that require authentication
		userGroup := protectedRoutes.Group("/users")
		{
			userGroup.GET("/:id", userHandler.GetUser)
			userGroup.PUT("/:id", userHandler.UpdateUser)
			userGroup.PATCH("/:id/handle", userHandler.UpdateHandle)
			userGroup.PATCH("/:id/onboarded", userHandler.UpdateOnboardedStatus)
			userGroup.DELETE("/:id", userHandler.DeleteUser)
		}
		
		// Content routes that require authentication
		contentGroup := protectedRoutes.Group("/content")
		{
			contentGroup.POST("", contentHandler.CreateContentItem)
			contentGroup.GET("/:id", contentHandler.GetContentItem)
			contentGroup.PUT("/:id", contentHandler.UpdateContentItem)
			contentGroup.PATCH("/:id/position", contentHandler.UpdateContentItemPosition)
			contentGroup.DELETE("/:id", contentHandler.DeleteContentItem)
		}
	}
	
	// Admin routes
	adminRoutes := protectedRoutes.Group("/admin")
	adminRoutes.Use(adminMiddleware)
	{
		adminRoutes.PATCH("/users/:id/premium", userHandler.UpdatePremiumStatus)
		adminRoutes.PATCH("/users/:id/admin", userHandler.UpdateAdminstatus)
		
	}
	
	// Health check route
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
}

func (s *Server) MountHandlers() {
	api := s.router.Group("/api")
	api.POST("/users")
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) Router() *gin.Engine {
	return s.router
}
