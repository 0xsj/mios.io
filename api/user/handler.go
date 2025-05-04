package user

import (
	"net/http"

	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/pkg/response"
	"github.com/0xsj/gin-sqlc/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	userService service.UserService
	logger      log.Logger
}

func NewHandler(userService service.UserService, logger log.Logger) *Handler {
	return &Handler{
		userService: userService,
		logger:      logger,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	userGroup := r.Group("/api/users")
	{
		userGroup.POST("", h.CreateUser)
		userGroup.GET("/:id", h.GetUser)
		userGroup.GET("/username/:username", h.GetUserByUsername)
		userGroup.GET("/handle/:handle", h.GetUserByHandle)
		userGroup.GET("/email/:email", h.GetUserByEmail)
		userGroup.PUT("/:id", h.UpdateUser)
		userGroup.PATCH("/:id/handle", h.UpdateHandle)
		userGroup.PATCH("/:id/premium", h.UpdatePremiumStatus)
		userGroup.PATCH("/:id/admin", h.UpdateAdminStatus)
		userGroup.PATCH("/:id/onboarded", h.UpdateOnboardedStatus)
		userGroup.DELETE("/:id", h.DeleteUser)
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	h.logger.Info("CreateUser handler called")

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}
	
	h.logger.Debugf("Received create user request: %+v", req)

	// Default handle to username if not provided
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
		h.logger.Errorf("Failed to create user: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	id, _ := uuid.Parse(user.ID)
	responseData := UserResponse{
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

	h.logger.Infof("User created successfully with ID: %s", user.ID)
	response.Success(c, responseData, "User created successfully", http.StatusCreated)
}

func (h *Handler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	h.logger.Debugf("GetUser handler called for user ID: %s", userID)

	user, err := h.userService.GetUser(c, userID)
	if err != nil {
		h.logger.Warnf("Failed to retrieve user: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	id, _ := uuid.Parse(user.ID)
	responseData := UserResponse{
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

	h.logger.Debugf("User retrieved successfully with ID: %s", user.ID)
	response.Success(c, responseData, "User retrieved successfully")
}

func (h *Handler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	h.logger.Debugf("GetUserByUsername handler called for username: %s", username)
	
	if username == "" {
		h.logger.Warn("Invalid username parameter: empty value")
		response.Error(c, response.ErrBadRequestResponse, "Username cannot be empty")
		return
	}

	user, err := h.userService.GetUserByUsername(c, username)
	if err != nil {
		h.logger.Warnf("Failed to retrieve user by username: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	id, _ := uuid.Parse(user.ID)
	responseData := UserResponse{
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

	h.logger.Debugf("User retrieved successfully by username: %s", username)
	response.Success(c, responseData, "User retrieved successfully")
}

func (h *Handler) GetUserByHandle(c *gin.Context) {
	handle := c.Param("handle")
	h.logger.Debugf("GetUserByHandle handler called for handle: %s", handle)
	
	if handle == "" {
		h.logger.Warn("Invalid handle parameter: empty value")
		response.Error(c, response.ErrBadRequestResponse, "Handle cannot be empty")
		return
	}

	user, err := h.userService.GetUserByHandle(c, handle)
	if err != nil {
		h.logger.Warnf("Failed to retrieve user by handle: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	id, _ := uuid.Parse(user.ID)
	responseData := UserResponse{
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

	h.logger.Debugf("User retrieved successfully by handle: %s", handle)
	response.Success(c, responseData, "User retrieved successfully")
}

func (h *Handler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	h.logger.Debugf("GetUserByEmail handler called for email: %s", email)
	
	if email == "" {
		h.logger.Warn("Invalid email parameter: empty value")
		response.Error(c, response.ErrBadRequestResponse, "Email cannot be empty")
		return
	}

	user, err := h.userService.GetUserByEmail(c, email)
	if err != nil {
		h.logger.Warnf("Failed to retrieve user by email: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	id, _ := uuid.Parse(user.ID)
	responseData := UserResponse{
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

	h.logger.Debugf("User retrieved successfully by email: %s", email)
	response.Success(c, responseData, "User retrieved successfully")
}

func (h *Handler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	h.logger.Infof("UpdateUser handler called for user ID: %s", userID)

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
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
		h.logger.Errorf("Failed to update user: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	id, _ := uuid.Parse(updatedUser.ID)
	responseData := UserResponse{
		ID:        id,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		IsPremium: updatedUser.IsPremium,
	}

	h.logger.Infof("User updated successfully with ID: %s", updatedUser.ID)
	response.Success(c, responseData, "User updated successfully")
}

func (h *Handler) UpdateHandle(c *gin.Context) {
	userID := c.Param("id")
	h.logger.Infof("UpdateHandle handler called for user ID: %s", userID)

	var req UpdateHandleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	if req.Handle == "" {
		h.logger.Warn("Invalid handle parameter: empty value")
		response.Error(c, response.ErrBadRequestResponse, "Handle cannot be empty")
		return
	}

	updatedUser, err := h.userService.UpdateHandle(c, userID, req.Handle)
	if err != nil {
		h.logger.Errorf("Failed to update user handle: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	id, _ := uuid.Parse(updatedUser.ID)
	responseData := UserResponse{
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

	h.logger.Infof("User handle updated successfully to '%s' for user ID: %s", req.Handle, updatedUser.ID)
	response.Success(c, responseData, "User handle updated successfully")
}

func (h *Handler) UpdatePremiumStatus(c *gin.Context) {
	userID := c.Param("id")
	h.logger.Infof("UpdatePremiumStatus handler called for user ID: %s", userID)

	var req UpdatePremiumStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	updatedUser, err := h.userService.UpdatePremiumStatus(c, userID, req.IsPremium)
	if err != nil {
		h.logger.Errorf("Failed to update premium status: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	id, _ := uuid.Parse(updatedUser.ID)
	responseData := UserResponse{
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

	h.logger.Infof("User premium status updated to %v for user ID: %s", req.IsPremium, updatedUser.ID)
	response.Success(c, responseData, "User premium status updated successfully")
}

func (h *Handler) UpdateAdminStatus(c *gin.Context) {
	userID := c.Param("id")
	h.logger.Infof("UpdateAdminStatus handler called for user ID: %s", userID)

	var req UpdateAdminStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	updatedUser, err := h.userService.UpdateAdminStatus(c, userID, req.IsAdmin)
	if err != nil {
		h.logger.Errorf("Failed to update admin status: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	id, _ := uuid.Parse(updatedUser.ID)
	responseData := UserResponse{
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

	h.logger.Infof("User admin status updated to %v for user ID: %s", req.IsAdmin, updatedUser.ID)
	response.Success(c, responseData, "User admin status updated successfully")
}

func (h *Handler) UpdateOnboardedStatus(c *gin.Context) {
	userID := c.Param("id")
	h.logger.Infof("UpdateOnboardedStatus handler called for user ID: %s", userID)

	var req UpdateOnboardedStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid request format: %v", err)
		response.Error(c, response.ErrBadRequestResponse, err.Error())
		return
	}

	updatedUser, err := h.userService.UpdateOnboardedStatus(c, userID, req.Onboarded)
	if err != nil {
		h.logger.Errorf("Failed to update onboarded status: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	id, _ := uuid.Parse(updatedUser.ID)
	responseData := UserResponse{
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

	h.logger.Infof("User onboarded status updated to %v for user ID: %s", req.Onboarded, updatedUser.ID)
	response.Success(c, responseData, "User onboarded status updated successfully")
}

func (h *Handler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	h.logger.Infof("DeleteUser handler called for user ID: %s", userID)

	err := h.userService.DeleteUser(c, userID)
	if err != nil {
		h.logger.Errorf("Failed to delete user: %v", err)
		response.HandleError(c, err, h.logger)
		return
	}

	h.logger.Infof("User deleted successfully with ID: %s", userID)
	response.Success(c, nil, "User deleted successfully", http.StatusOK)
}