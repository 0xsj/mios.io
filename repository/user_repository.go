package repository

import (
	"context"
	"errors"

	db "github.com/0xsj/gin-sqlc/db/sqlc"
	"github.com/google/uuid"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateKey   = errors.New("duplicate key violation")
	ErrDatabase       = errors.New("database error")

	ErrInvalidInput = errors.New("invalid input parameters")
    ErrPermissionDenied = errors.New("permission denied for operation")
    ErrForeignKeyViolation = errors.New("foreign key constraint violation")
    ErrTransactionFailed = errors.New("database transaction failed")
    ErrConnectionFailed = errors.New("database connection failed")
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (*db.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*db.User, error)
	GetUserByUsername(ctx context.Context, username string) (*db.User, error)
	GetUserByHandle(ctx context.Context, handle string) (*db.User, error)
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
	UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error
	UpdateHandle(ctx context.Context, userID uuid.UUID, handle string) error
	UpdateEmail(ctx context.Context, userID uuid.UUID, email string) error
	UpdatePremiumStatus(ctx context.Context, userID uuid.UUID, isPremium bool) error
	UpdateAdminStatus(ctx context.Context, userID uuid.UUID, isAdmin bool) error
	UpdateOnboardedStatus(ctx context.Context, userID uuid.UUID, onboarded bool) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type CreateUserParams struct {
	Username        string
	Handle          string
	Email           string
	FirstName       string
	LastName        string
	ProfileImageURL string
	Bio             string
	LayoutVersion   string
	CustomDomain    string
	IsPremium       bool
	IsAdmin         bool
	Onboarded       bool
}

type UpdateUserParams struct {
	UserID          uuid.UUID
	FirstName       string
	LastName        string
	ProfileImageURL string
	Bio             string
	LayoutVersion   string
	CustomDomain    string
}

type SQLCUserRepository struct {
	db *db.Queries
}

func NewUserRepository(db *db.Queries) UserRepository {
	return &SQLCUserRepository{
		db: db,
	}
}

func (r *SQLCUserRepository) CreateUser(ctx context.Context, arg CreateUserParams) (*db.User, error){}

func (r *SQLCUserRepository) GetUser(ctx context.Context, userID uuid.UUID) (*db.User, error ){}

func (r *SQLCUserRepository) GetUserByUsername(ctx context.Context, username string) (*db.User, error) {}

func (r *SQLCUserRepository) GetUserByHandle(ctx context.Context, handle string) (*db.User, error) {}

func (r *SQLCUserRepository) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {}

func (r *SQLCUserRepository) UpdateUser(ctx context.Context, arg UpdateUserParams) error {}

func (r *SQLCUserRepository) UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error {}

func (r *SQLCUserRepository) UpdateHandle(ctx context.Context, userID uuid.UUID, handle string) error {}

func (r *SQLCUserRepository) UpdateEmail(ctx context.Context, userID uuid.UUID, email string) error {}

func (r *SQLCUserRepository) UpdatePremiumStatus(ctx context.Context, userID uuid.UUID, isPremium bool) {}

func (r *SQLCUserRepository) UpdateAdminStatus(ctx context.Context, userID uuid.UUID, isAdmin bool) error {}

func (r *SQLCUserRepository) UpdateOnboardedStatus(ctx context.Context, userID uuid.UUID, onboarded bool) {}

func (r *SQLCUserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {}