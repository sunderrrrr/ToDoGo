package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type ConnConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

const (
	userTable      = "users"
	todoListsTable = "todo_lists"
	userListsTable = "users_lists"
	todoItemsTable = "todo_items"
	listItemsTable = "lists_items"
)

func NewPostgresDB(cfg ConnConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("postgres.go: error connecting to database: %s", err.Error())
	}
	return db, nil
}
