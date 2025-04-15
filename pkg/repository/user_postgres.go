package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (ur *UserPostgres) ResetPassword(username string, oldPassword string, newPassword string) error {
	tx, err := ur.db.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("UPADTE %s SET password_hash=$1 WHERE username=$2 and password_hash=$3", userTable)
	_, err = tx.Exec(query, newPassword, username, oldPassword)
	if err != nil {
		tx.Rollback()
	}
	return tx.Commit()
}
