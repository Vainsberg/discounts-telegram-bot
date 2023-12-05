package main

import (
	"fmt"
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/internal/handler"
	getDiscounts "github.com/Vainsberg/discounts-telegram-bot/internal/handler"
	"github.com/Vainsberg/discounts-telegram-bot/internal/viper"
	"github.com/Vainsberg/discounts-telegram-bot/repository"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error
	cfg, err := viper.NewConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	db := getDiscounts.CreateDB(cfg)
	defer db.Close()
	repository := repository.NewRepository(db)
	handler := handler.NewHandler(repository)
	http.HandleFunc("/discount", handler.GetDiscounts)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
}
