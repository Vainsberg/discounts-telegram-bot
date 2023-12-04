package main

import (
	"fmt"
	"net/http"

	getDiscounts "github.com/Vainsberg/discounts-telegram-bot/internal/GetDiscounts"
	"github.com/Vainsberg/discounts-telegram-bot/repository"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error
	db := getDiscounts.CreateDB()
	defer db.Close()
	repository := repository.NewRepository(db)
	handler := getDiscounts.NewHandler(repository)
	http.HandleFunc("/discount", handler.GetDiscounts)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
}
