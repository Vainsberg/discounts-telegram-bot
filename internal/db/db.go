package db

import (
	"database/sql"
	"log"

	"github.com/Vainsberg/discounts-telegram-bot/internal/viper"
)

func CreateDB(cfg *viper.Config) *sql.DB {
	var err error
	db, err := sql.Open("mysql", cfg.DbUser+":"+cfg.DbPass+"@tcp(127.0.0.1:3306)/discounts")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS goods (
			id INTEGER PRIMARY KEY AUTO_INCREMENT,
		    name TEXT,
			price_ru  FLOAT,
			url TEXT,
			image TEXT,
			dt DATETIME DEFAULT CURRENT_TIMESTAMP,
			query TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS subscriptions (
			id INTEGER PRIMARY KEY AUTO_INCREMENT,
			chat_id BIGINT,
			query TEXT
		)
		  `)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
