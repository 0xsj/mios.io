package repository

import db "github.com/0xsj/gin-sqlc/db/sqlc"

type AuthRepository interface {}

type CreateAuthParams struct {}

type SQLCAuthRepository struct {
	db *db.Queries
}

func NewAuthRepository(db *db.Queries) AuthRepository {
	return &SQLCAuthRepository{
		db: db,
	}
}

func CreateAuth(){}

func GetAuthByUserID(){}

func UpdatePassword(){}

func UpdateResetToken(){}

func VerifyEmail(){}

func UpdateLastLogin(){}