package repository

import (
	"database/sql"
	"log"

	"github.com/Vainsberg/discounts-telegram-bot/internal/dto"
)

type RepositorySubs struct {
	db *sql.DB
}

func NewRepositorySubs(db *sql.DB) *RepositorySubs {
	return &RepositorySubs{db: db}
}

func (r *RepositorySubs) AddLincked(chatID int64, text string) error {
	_, err := r.db.Exec("INSERT INTO subscriptions (chat_id, query) VALUES (?, ?);", chatID, text)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (r *RepositorySubs) GetQuerys() (dto.GetQuerys, error) {
	rows, err := r.db.Query("SELECT DISTINCT query FROM subscriptions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	responseN := dto.GetQuerys{}
	for rows.Next() {
		var item dto.GetQuerysItem
		err := rows.Scan(&item.Query)
		if err != nil {
			log.Fatal(err)

		}
		responseN.Items = append(responseN.Items, item)
	}
	return responseN, nil
}

func (r *RepositorySubs) SearchChatID(query string) int64 {
	var ChatIDbyDiscounts int64
	rank := r.db.QueryRow("SELECT chat_id FROM subscriptions WHERE query = ?;", query)
	if err := rank.Scan(&ChatIDbyDiscounts); err != nil && err != sql.ErrNoRows {
		return 0
	}
	return ChatIDbyDiscounts
}
