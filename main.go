package main

import (
	"fmt"
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/internal/client"
	"github.com/Vainsberg/discounts-telegram-bot/internal/db"
	"github.com/Vainsberg/discounts-telegram-bot/internal/handler"
	"github.com/Vainsberg/discounts-telegram-bot/internal/repository"
	"github.com/Vainsberg/discounts-telegram-bot/internal/viper"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error
	cfg, err := viper.NewConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	db := db.CreateDB(cfg)
	defer db.Close()
	repository := repository.NewRepository(db)
	api := client.NewPlatiClient("https://plati.io")
	handler := handler.NewHandler(repository, api)
	http.HandleFunc("/discount", handler.GetDiscounts)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
}
