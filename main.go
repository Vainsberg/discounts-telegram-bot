package main

import (
	"bytes"
	"encoding/json"
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
			if update.Message != nil {
				bottg.HandleRequest(bot, update.Message, &update)
				userText := update.Message.Text
				replyKeyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Подписаться на скидки", userText),
					),
				)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нажми на кнопку для подписки на товар :)")
				msg.ReplyMarkup = replyKeyboard
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
			} else if update.CallbackQuery != nil {
				callback := update.CallbackQuery

				if callback.Data != "" {
					userText := callback.Data
					chatID := callback.Message.Chat.ID
					ApiURL := "http://localhost:8080/subscribe"
					payload := fmt.Sprintf("chat_id=%d&text=%s", chatID, userText)
					payloadmash, err := json.Marshal(payload)
					if err != nil {
						log.Println("Ошибка JSON:", err)
						return
					}
					resp, err := http.Post(ApiURL, "application/json", bytes.NewBuffer(payloadmash))
					if err != nil {
						fmt.Println("Ошибка при выполнении запроса:", err)
						return
					}
					defer resp.Body.Close()

					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Вы подписались на товар")
					if _, err := bot.Request(callback); err != nil {
						panic(err)
					}
				}
			}
		}
	}()

	db := db.CreateDB(cfg)
	defer db.Close()

	repositoryGoods := repository.NewRepository(db)
	RepositorySubs := repository.NewRepositorySubs(db)
	api := client.NewPlatiClient("https://plati.io")
	handler := handler.NewHandler(repositoryGoods, api, RepositorySubs)
	http.HandleFunc("/discount", handler.GetDiscounts)
	http.HandleFunc("/subscribe", handler.AddSubscription)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
}
