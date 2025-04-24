// user/request.go
package user

import "github.com/google/uuid"

type CreateUserRequest struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Handle          string `json:"handle"`
	Bio             string `json:"bio"`
	ProfileImageURL string `json:"profile_image_url"`
	LayoutVersion   string `json:"layout_version"`
	CustomDomain    string `json:"custom_domain"`
	IsPremium       bool   `json:"is_premium"`
	IsAdmin         bool   `json:"is_admin"`
	Onboarded       bool   `json:"onboarded"`
}

type UserResponse struct {
	ID              uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	Handle          string    `json:"handle,omitempty"`
	Email           string    `json:"email"`
	FirstName       string    `json:"first_name,omitempty"`
	LastName        string    `json:"last_name,omitempty"`
	Bio             string    `json:"bio,omitempty"`
	ProfileImageURL string    `json:"profile_image_url,omitempty"`
	LayoutVersion   string    `json:"layout_version,omitempty"`
	CustomDomain    string    `json:"custom_domain,omitempty"`
	IsPremium       bool      `json:"is_premium"`
	IsAdmin         bool      `json:"is_admin"`
	Onboarded       bool      `json:"onboarded"`
}

type UpdateUserRequest struct {
	Username        *string `json:"username"`
	Email           *string `json:"email"`
	FirstName       *string `json:"first_name"`
	LastName        *string `json:"last_name"`
	Bio             *string `json:"bio"`
	ProfileImageURL *string `json:"profile_image_url"`
	LayoutVersion   *string `json:"layout_version"`
	CustomDomain    *string `json:"custom_domain"`
}

type UpdateHandleRequest struct {
	Handle string `json:"handle" binding:"required"`
}

type UpdatePremiumStatusRequest struct {
	IsPremium bool `json:"is_premium"`
}

type UpdateAdminStatusRequest struct {
	IsAdmin bool `json:"is_admin"`
}

type UpdateOnboardedStatusRequest struct {
	Onboarded bool `json:"onboarded"`
}
