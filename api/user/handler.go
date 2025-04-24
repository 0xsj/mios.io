package user

import (
	"fmt"
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
		userGroup.GET("/:id", h.GetUser)
		userGroup.GET("/username/:username", h.GetUserByUsername)
		userGroup.GET("/handle/:handle", h.GetUserbyHandle)
		userGroup.GET("/email/:email", h.GetUserByEmail)
		userGroup.PUT("/:id", h.UpdateUser)
		userGroup.PATCH("/:id/handle", h.UpdateHandle)
		userGroup.PATCH("/:id/premium", h.UpdatePremiumStatus)
		userGroup.PATCH("/:id/admin", h.UpdateAdminstatus)
		userGroup.PATCH("/:id/onboarded", h.UpdateOnboardedStatus)
		userGroup.DELETE("/:id", h.DeleteUser)
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	fmt.Println("CreateUser handler called")

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Request binding error: %v\n", err)
		api.HandleError(c, api.ErrInvalidInput)
		return
	}
	fmt.Printf("Received request: %+v\n", req)

	handle := req.Username
	if req.Handle != "" {
		handle = req.Handle
	}

	input := service.CreateUserInput{
		Username:        req.Username,
		Handle:          handle,
		Email:           req.Email,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Bio:             req.Bio,
		ProfileImageURL: req.ProfileImageURL,
		LayoutVersion:   req.LayoutVersion,
		CustomDomain:    req.CustomDomain,
		IsPremium:       req.IsPremium,
		IsAdmin:         req.IsAdmin,
		Onboarded:       req.Onboarded,
	}

	user, err := h.userService.CreateUser(c, input)
	if err != nil {
		fmt.Printf("Service error: %v\n", err)
		api.HandleError(c, err)
		return
	}

	id, _ := uuid.Parse(user.ID)
	response := UserResponse{
		ID:              id,
		Username:        user.Username,
		Handle:          user.Handle,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Bio:             user.Bio,
		ProfileImageURL: user.ProfileImageURL,
		LayoutVersion:   user.LayoutVersion,
		CustomDomain:    user.CustomDomain,
		IsPremium:       user.IsPremium,
		IsAdmin:         user.IsAdmin,
		Onboarded:       user.Onboarded,
	}

	api.RespondWithSuccess(c, response, "User created successfully", http.StatusCreated)
}

func (h *Handler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.userService.GetUser(c, userID)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	id, _ := uuid.Parse(user.ID)
	response := UserResponse{
		ID:              id,
		Username:        user.Username,
		Handle:          user.Handle,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Bio:             user.Bio,
		ProfileImageURL: user.ProfileImageURL,
		LayoutVersion:   user.LayoutVersion,
		CustomDomain:    user.CustomDomain,
		IsPremium:       user.IsPremium,
		IsAdmin:         user.IsAdmin,
		Onboarded:       user.Onboarded,
	}

	api.RespondWithSuccess(c, response, "User retreived successfully!")
}

func (h *Handler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	user, err := h.userService.GetUserByUsername(c, username)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	id, _ := uuid.Parse(user.ID)
	response := UserResponse{
		ID:              id,
		Username:        user.Username,
		Handle:          user.Handle,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Bio:             user.Bio,
		ProfileImageURL: user.ProfileImageURL,
		LayoutVersion:   user.LayoutVersion,
		CustomDomain:    user.CustomDomain,
		IsPremium:       user.IsPremium,
		IsAdmin:         user.IsAdmin,
		Onboarded:       user.Onboarded,
	}

	api.RespondWithSuccess(c, response, "User retreived successfully")
}

func (h *Handler) GetUserbyHandle(c *gin.Context) {
	handle := c.Param("handle")
	if handle == "" {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	user, err := h.userService.GetUserByHandle(c, handle)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	id, _ := uuid.Parse(user.ID)
	response := UserResponse{
		ID:              id,
		Username:        user.Username,
		Handle:          user.Handle,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Bio:             user.Bio,
		ProfileImageURL: user.ProfileImageURL,
		LayoutVersion:   user.LayoutVersion,
		CustomDomain:    user.CustomDomain,
		IsPremium:       user.IsPremium,
		IsAdmin:         user.IsAdmin,
		Onboarded:       user.Onboarded,
	}

	api.RespondWithSuccess(c, response, "User retreived successfully")
}

func (h *Handler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	user, err := h.userService.GetUserByEmail(c, email)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	id, _ := uuid.Parse(user.ID)
	response := UserResponse{
		ID:              id,
		Username:        user.Username,
		Handle:          user.Handle,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Bio:             user.Bio,
		ProfileImageURL: user.ProfileImageURL,
		LayoutVersion:   user.LayoutVersion,
		CustomDomain:    user.CustomDomain,
		IsPremium:       user.IsPremium,
		IsAdmin:         user.IsAdmin,
		Onboarded:       user.Onboarded,
	}

	api.RespondWithSuccess(c, response, "User retreived successfully")
}

func (h *Handler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Request binding error: %v\n", err)
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	input := service.UpdateUserInput{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	updatedUser, err := h.userService.UpdateUser(c, userID, input)
	if err != nil {
		fmt.Printf("Service error: %v\n", err)
		api.HandleError(c, err)
		return
	}

	id, _ := uuid.Parse(updatedUser.ID)
	response := UserResponse{
		ID:        id,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		IsPremium: updatedUser.IsPremium,
	}

	api.RespondWithSuccess(c, response, "User updated successfully")
}

func (h *Handler) UpdateHandle(c *gin.Context) {
	userID := c.Param("id")

	var req UpdateHandleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	if req.Handle == "" {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	updatedUser, err := h.userService.UpdateHandle(c, userID, req.Handle)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	id, _ := uuid.Parse(updatedUser.ID)
	response := UserResponse{
		ID:              id,
		Username:        updatedUser.Username,
		Handle:          updatedUser.Handle,
		Email:           updatedUser.Email,
		FirstName:       updatedUser.FirstName,
		LastName:        updatedUser.LastName,
		Bio:             updatedUser.Bio,
		ProfileImageURL: updatedUser.ProfileImageURL,
		LayoutVersion:   updatedUser.LayoutVersion,
		CustomDomain:    updatedUser.CustomDomain,
		IsPremium:       updatedUser.IsPremium,
		IsAdmin:         updatedUser.IsAdmin,
		Onboarded:       updatedUser.Onboarded,
	}

	api.RespondWithSuccess(c, response, "User handle updated successfully")
}

func (h *Handler) UpdatePremiumStatus(c *gin.Context) {
	userID := c.Param("id")

	var req UpdatePremiumStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	updatedUser, err := h.userService.UpdatePremiumStatus(c, userID, req.IsPremium)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	id, _ := uuid.Parse(updatedUser.ID)
	response := UserResponse{
		ID:              id,
		Username:        updatedUser.Username,
		Handle:          updatedUser.Handle,
		Email:           updatedUser.Email,
		FirstName:       updatedUser.FirstName,
		LastName:        updatedUser.LastName,
		Bio:             updatedUser.Bio,
		ProfileImageURL: updatedUser.ProfileImageURL,
		LayoutVersion:   updatedUser.LayoutVersion,
		CustomDomain:    updatedUser.CustomDomain,
		IsPremium:       updatedUser.IsPremium,
		IsAdmin:         updatedUser.IsAdmin,
		Onboarded:       updatedUser.Onboarded,
	}

	api.RespondWithSuccess(c, response, "User premium status updated successfully")
}

func (h *Handler) UpdateAdminstatus(c *gin.Context) {
	userID := c.Param("id")
	var req UpdateAdminStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	updatedUser, err := h.userService.UpdateAdminStatus(c, userID, req.IsAdmin)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	id, _ := uuid.Parse(updatedUser.ID)
	response := UserResponse{
		ID:              id,
		Username:        updatedUser.Username,
		Handle:          updatedUser.Handle,
		Email:           updatedUser.Email,
		FirstName:       updatedUser.FirstName,
		LastName:        updatedUser.LastName,
		Bio:             updatedUser.Bio,
		ProfileImageURL: updatedUser.ProfileImageURL,
		LayoutVersion:   updatedUser.LayoutVersion,
		CustomDomain:    updatedUser.CustomDomain,
		IsPremium:       updatedUser.IsPremium,
		IsAdmin:         updatedUser.IsAdmin,
		Onboarded:       updatedUser.Onboarded,
	}

	api.RespondWithSuccess(c, response, "User admin status updated successfully")
}

func (h *Handler) UpdateOnboardedStatus(c *gin.Context) {
	userID := c.Param("id")

	var req UpdateOnboardedStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleError(c, api.ErrInvalidInput)
		return
	}

	updatedUser, err := h.userService.UpdateOnboardedStatus(c, userID, req.Onboarded)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	id, _ := uuid.Parse(updatedUser.ID)
	response := UserResponse{
		ID:              id,
		Username:        updatedUser.Username,
		Handle:          updatedUser.Handle,
		Email:           updatedUser.Email,
		FirstName:       updatedUser.FirstName,
		LastName:        updatedUser.LastName,
		Bio:             updatedUser.Bio,
		ProfileImageURL: updatedUser.ProfileImageURL,
		LayoutVersion:   updatedUser.LayoutVersion,
		CustomDomain:    updatedUser.CustomDomain,
		IsPremium:       updatedUser.IsPremium,
		IsAdmin:         updatedUser.IsAdmin,
		Onboarded:       updatedUser.Onboarded,
	}

	api.RespondWithSuccess(c, response, "Update onboarded ")
}

func (h *Handler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	err := h.userService.DeleteUser(c, userID)
	if err != nil {
		api.HandleError(c, err)
		return
	}

	api.RespondWithSuccess(c, nil, "User deleted successfully", http.StatusOK)
}
