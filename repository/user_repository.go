package repository

import (
	"context"

	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/google/uuid"
)


type UserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (db.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (db.User, error)
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
	UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error
	UpdateEmail(ctx context.Context, userID uuid.UUID, email string) error
	UpdatePremiumStatus(ctx context.Context, userID uuid.UUID, isPremium bool) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type CreateUserParams struct {}

type UpdateUserParams struct {}

type SQLCUserRepository struct {
	db *db.Queries
}

