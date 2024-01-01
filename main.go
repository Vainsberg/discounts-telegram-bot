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
	"github.com/Vainsberg/discounts-telegram-bot/internal/service"
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
	bottg.RunBot(bot)

	logger, err := zap.NewProduction()
	if err != nil {
		panic("Error create logger")
	}
	defer logger.Sync()

	db := db.CreateDB(cfg)
	defer db.Close()

	repositoryGoods := repository.NewRepository(db)
	RepositorySubs := repository.NewRepositorySubs(db)
	api := client.NewPlatiClient("https://plati.io")
	service := service.NewService(logger, *api, RepositorySubs, *repositoryGoods, *bot)
	handler := handler.NewHandler(logger, bot, *service)
	http.HandleFunc("/discount", handler.GetDiscounts)
	http.HandleFunc("/subscribe", handler.AddSubscription)
	http.HandleFunc("/discount/update", handler.GetQuerysCron)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
}
