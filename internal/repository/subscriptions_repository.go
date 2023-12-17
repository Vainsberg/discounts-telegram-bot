package repository

import (
	"database/sql"
	"log"
)

type RepositorySubs struct {
	db *sql.DB
}

func NewRepositorySubs(db *sql.DB) *RepositorySubs {
	return &RepositorySubs{db: db}
}

func (r *RepositorySubs) AddLincked(chatID int, text string) error {

	_, err := r.db.Exec("INSERT INTO linked_accounts (name, query) VALUES (?, ?);", chatID, text)
	if err != nil {
		log.Fatal(err)
	}
	return nil

}
