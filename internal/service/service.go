package service

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Vainsberg/discounts-telegram-bot/internal/client"
	cronhandler "github.com/Vainsberg/discounts-telegram-bot/internal/cron_handler"
	"github.com/Vainsberg/discounts-telegram-bot/internal/dto"
	"github.com/Vainsberg/discounts-telegram-bot/internal/pkg"
	"github.com/Vainsberg/discounts-telegram-bot/internal/repository"
	"github.com/Vainsberg/discounts-telegram-bot/internal/response"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Service struct {
	SubsRepository       repository.RepositorySubs
	Logger               *zap.Logger
	DiscountsPlatiClient client.PlatiClient
	DiscountsRepository  repository.Repository
	Bot                  *tgbotapi.BotAPI
}

func NewService(
	logger *zap.Logger,
	discountsPlatiClient *client.PlatiClient,
	repositorySubs *repository.RepositorySubs,
	discountsRepository *repository.Repository,
	bot *tgbotapi.BotAPI,
) *Service {
	return &Service{
		Logger:               logger,
		DiscountsPlatiClient: *discountsPlatiClient,
		SubsRepository:       *repositorySubs,
		DiscountsRepository:  *discountsRepository,
		Bot:                  bot,
	}
}

func (s *Service) ProcessQueryAndFetchGoods() error {
	querys, err := s.SubsRepository.GetQuerys()
	if err != nil {
		s.Logger.Info("DiscountsPlatiClient error:", zap.Error(err))
		return err
	}

	for _, v := range querys.Items {
		query := pkg.ReplaceSpaceUrl(v.Query)
		s.Logger.Info("Sleep 5 second...", zap.String("query", query))
		time.Sleep(5 * time.Second)
		goods, err := s.DiscountsPlatiClient.GetGoodsClient(query)
		if err != nil {
			s.Logger.Info("DiscountsPlatiClient error:", zap.Error(err))
			return err
		}

		for _, el := range goods.Items {
			RequestDiscounts := dto.Item{
				Name:      el.Name,
				Price_rur: el.Price_rur,
				Url:       el.Url,
				Image:     el.Image,
			}
			s.processGoodsItem(RequestDiscounts, query)
		}
	}
	return nil
}

func (s *Service) processGoodsItem(item dto.Item, query string) {
	s.DiscountsRepository.SaveGood(item.Name, float64(item.Price_rur), item.Url, item.Image, query)
	pastPrice := s.DiscountsRepository.SearchAveragePrice(float64(item.Price_rur), item.Url)

	if pastPrice < float64(item.Price_rur) {
		rebate := pkg.CalculatePercentageDifference(pastPrice, float64(item.Price_rur))

		if rebate >= 20 {
			s.Logger.Info("Finding the updated price", zap.Int("Price_rur", item.Price_rur))
			chatsID := s.SubsRepository.SearchChatID(query)
			goodDiscounts := cronhandler.ProductDiscounts(item.Name, float64(item.Price_rur), item.Url, item.Image)

			for _, v := range goodDiscounts.Items {
				productDiscount := dto.Item{
					Name:      v.Name,
					Price_rur: v.Price_rur,
					Url:       v.Url,
					Image:     v.Image,
				}
				s.sendDiscount(productDiscount, chatsID)
			}
		}
	}
}

func (s *Service) sendDiscount(item dto.Item, chatsID int64) {
	text := fmt.Sprintf(
		"*%s*\n"+
			"*Rub* _%v_\n"+
			"*Ссылка* _%s_\n",
		item.Name, item.Price_rur, item.Url)
	respy, err := http.Get("https:" + item.Image)
	if err != nil {
		log.Println("Ошибка при получении изображения:", err)
		return
	}
	defer respy.Body.Close()

	reader := tgbotapi.FileReader{Name: "photo.jpg", Reader: respy.Body}
	photoMsg := tgbotapi.NewPhoto(chatsID, reader)
	_, err = s.Bot.Send(photoMsg)
	if err != nil {
		log.Println("Ошибка при отправке сообщения боту:", err)
	}

	msg := tgbotapi.NewMessage(chatsID, text)
	msg.ParseMode = "markdown"
	s.Logger.Info("Product withdrawal with a discount")
	_, err = s.Bot.Send(msg)
}

func (s *Service) FetchAndSaveGoods(CheckQueryText string) *response.RequestDiscounts {
	goods, err := s.DiscountsPlatiClient.GetGoodsClient(CheckQueryText)
	if err != nil {
		s.Logger.Info("DiscountsPlatiClient error:", zap.Error(err))
		return nil
	}

	for _, v := range goods.Items {
		err := s.DiscountsRepository.SaveGood(v.Name, float64(v.Price_rur), v.Url, v.Image, CheckQueryText)
		if err != nil {
			s.Logger.Info("Error SaveGood", zap.Error(err))
		}
	}
	return goods
}

func (s *Service) AddLinked(chatID int64, requestText string) error {
	return s.SubsRepository.AddLincked(chatID, requestText)
}
