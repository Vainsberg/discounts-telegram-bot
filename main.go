package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/internal/bottg"
	"github.com/Vainsberg/discounts-telegram-bot/internal/client"
	cronhandler "github.com/Vainsberg/discounts-telegram-bot/internal/cron_handler"
	"github.com/Vainsberg/discounts-telegram-bot/internal/db"
	"github.com/Vainsberg/discounts-telegram-bot/internal/handler"
	"github.com/Vainsberg/discounts-telegram-bot/internal/repository"
	"github.com/Vainsberg/discounts-telegram-bot/internal/viper"
	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron"
	"go.uber.org/zap"
)

func main() {

	var err error
	cfg, err := viper.NewConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	c := cron.New()
	c.AddFunc(cfg.CountCron, func() {
		cronhandler.HandleCron()
	})
	c.Start()

	bot, err := tgbotapi.NewBotAPI(cfg.Apikey)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	logger, err := zap.NewProduction()
	if err != nil {
		panic("Error create logger")
	}
	defer logger.Sync()

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
				bottg.HandleCallback(bot, &update)
			}
		}
	}()

	db := db.CreateDB(cfg)
	defer db.Close()

	repositoryGoods := repository.NewRepository(db)
	RepositorySubs := repository.NewRepositorySubs(db)
	RepositoryQueys := repository.NewRepositorySubs(db)
	api := client.NewPlatiClient("https://plati.io")
	handler := handler.NewHandler(logger, repositoryGoods, api, RepositorySubs, RepositoryQueys, bot)
	http.HandleFunc("/discount", handler.GetDiscounts)
	http.HandleFunc("/subscribe", handler.AddSubscription)
	http.HandleFunc("/discount/update", handler.GetQuerysCron)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
	select {}
}
