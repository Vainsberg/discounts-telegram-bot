package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/internal/bottg"
	"github.com/Vainsberg/discounts-telegram-bot/internal/client"
	"github.com/Vainsberg/discounts-telegram-bot/internal/db"
	"github.com/Vainsberg/discounts-telegram-bot/internal/handler"
	"github.com/Vainsberg/discounts-telegram-bot/internal/repository"
	"github.com/Vainsberg/discounts-telegram-bot/internal/viper"
	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	var err error
	cfg, err := viper.NewConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Apikey)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}
			if update.Message.Text == "/start" {
				replyKeyboard := tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("Подписаться на скидки"),
					),
				)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Нажми кнопку:")
				msg.ReplyMarkup = replyKeyboard
				_, err := bot.Send(msg)
				if err != nil {
					log.Println("Ошибка при отправке сообщения боту:", err)
				}
			} else {
				bottg.HandleRequest(bot, update.Message)
			}
		}
	}()

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
