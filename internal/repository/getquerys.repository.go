package repository

import (
	"database/sql"
	"log"

	"github.com/Vainsberg/discounts-telegram-bot/internal/dto"
)

type RepositoryQuerys struct {
	db *sql.DB
}

func NewRepositoryQuerys(db *sql.DB) *RepositoryQuerys {
	return &RepositoryQuerys{db: db}
}

func (r *RepositoryQuerys) GetQuerys() (dto.GetQuerys, error) {

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
