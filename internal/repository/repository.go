package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Vainsberg/discounts-telegram-bot/internal/dto"
	"github.com/Vainsberg/discounts-telegram-bot/internal/response"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {

	return &Repository{db: db}

}

func (r *Repository) GetDiscountsByGoods(queryText string) response.RequestDiscounts {

	rows, err := r.db.Query("SELECT name, price_ru, url, image FROM goods WHERE query = ? AND dt >= CURRENT_TIMESTAMP() - INTERVAL 24 HOUR;", queryText)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	responseN := response.RequestDiscounts{}
	for rows.Next() {
		var item dto.Item
		err := rows.Scan(&item.Name, &item.Price_rur, &item.Url, &item.Image)
		if err != nil {
			log.Fatal(err)

		}
		responseN.Items = append(responseN.Items, item)
	}

	return responseN

}

func (r *Repository) SaveGood(name string, price_rur float64, url string, image string, queryText string) {
	responseN := response.RequestDiscounts{}
	for _, v := range responseN.Items {
		_, err := r.db.Exec("INSERT INTO goods (name, price_ru, url, image, dt, query) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP(), ?)", v.Name, v.Price_rur, v.Url, v.Image, queryText)
		if err != nil {
			log.Fatal(err)
			fmt.Println(err)
			return
		}
	}
}
