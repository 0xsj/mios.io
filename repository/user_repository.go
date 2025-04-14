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
}

type CreateUserParams struct {}

type UpdateUserParams struct {}

type SQLCUserRepository struct {
	db *db.Queries
}

