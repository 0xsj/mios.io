package user

import (
	"net/http"

	"github.com/0xsj/gin-sqlc/api"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	userService service.UserService
}

func NewHandler(userService service.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}


func (h *Handler) RegisterRoutes(r *gin.Engine) {
	userGroup := r.Group("/api/users")
	{
		userGroup.POST("", h.CreateUser)
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	input := service.CreateUserInput{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	user, err := h.userService.CreateUser(c, input)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	id, _ := uuid.Parse(user.ID)
	response := UserResponse{
		ID:       id,
		Username: user.Username,
		Email:    user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsPremium: user.IsPremium,
	}

	api.RespondWithSuccess(c, response, "User created successfully", http.StatusCreated)
}

