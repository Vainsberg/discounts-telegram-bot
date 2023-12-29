package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Vainsberg/discounts-telegram-bot/internal/client"
	cronhandler "github.com/Vainsberg/discounts-telegram-bot/internal/cron_handler"
	pkg "github.com/Vainsberg/discounts-telegram-bot/internal/pkg"
	"github.com/Vainsberg/discounts-telegram-bot/internal/repository"
	"github.com/Vainsberg/discounts-telegram-bot/internal/response"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Handler struct {
	Logger               *zap.Logger
	DiscountsRepository  repository.Repository
	DiscountsPlatiClient client.PlatiClient
	SubsRepository       repository.RepositorySubs
	RepositoryQuerys     repository.RepositorySubs
	Bot                  *tgbotapi.BotAPI
}

func NewHandler(logger *zap.Logger, repos *repository.Repository, plati *client.PlatiClient, subs *repository.RepositorySubs, querys *repository.RepositorySubs, bot *tgbotapi.BotAPI) *Handler {
	return &Handler{
		Logger:               logger,
		DiscountsRepository:  *repos,
		DiscountsPlatiClient: *plati,
		SubsRepository:       *subs,
		RepositoryQuerys:     *querys,
		Bot:                  bot,
	}
}

func (h *Handler) GetDiscounts(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Launch GetDiscounts")

	query, err := pkg.GetQuery(r.URL.Query().Get("query"))
	if err != nil {
		h.Logger.Info("GetQuery error:", zap.Error(err))
		return
	}
	CheckQueryText := pkg.Check(query)
	responseN := h.DiscountsRepository.GetDiscountsByGoods(CheckQueryText)

	if len(responseN.Items) != 0 {
		w.Write(client.DateFromDatebase(responseN))
		return
	}

	goods, err := h.DiscountsPlatiClient.GetGoodsClient(CheckQueryText)
	if err != nil {
		h.Logger.Info("DiscountsPlatiClient error:", zap.Error(err))
		return
	}

	for _, v := range goods.Items {
		err := h.DiscountsRepository.SaveGood(v.Name, float64(v.Price_rur), v.Url, v.Image, CheckQueryText)
		if err != nil {
			h.Logger.Info("Error SaveGood", zap.Error(err))
		}
	}

	respText, err := json.Marshal(goods)
	if err != nil {
		h.Logger.Info("Error encoding JSON response", zap.Error(err))
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(respText)
}

func (h *Handler) AddSubscription(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("AddSubscription")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Info("Read body error:", zap.Error(err))
		return
	}

	if err != nil {
		h.Logger.Info("Ошибка при чтении тела запроса", zap.Error(err))
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	var result response.SubscriptionRequest

	err = json.Unmarshal(body, &result)
	if err != nil {
		h.Logger.Info("Ошибка при разборе JSON:", zap.Error(err))
		return
	}
	h.SubsRepository.AddLincked(result.ChatID, result.Text)
}

func (h *Handler) GetQuerysCron(w http.ResponseWriter, r *http.Request) {
	query, err := h.RepositoryQuerys.GetQuerys()
	if err != nil {
		h.Logger.Info("DiscountsPlatiClient error:", zap.Error(err))
		return
	}

	for _, v := range query.Items {
		h.Logger.Info("Sleep 5 second...", zap.String("query", v.Query))
		time.Sleep(5 * time.Second)
		goods, err := h.DiscountsPlatiClient.GetGoodsClient(v.Query)
		if err != nil {
			h.Logger.Info("DiscountsPlatiClient error:", zap.Error(err))
			return
		}

		for _, el := range goods.Items {
			h.DiscountsRepository.SaveGood(el.Name, float64(el.Price_rur), el.Url, el.Image, v.Query)
			pastPrice := h.DiscountsRepository.SearchAveragePrice(float64(el.Price_rur), el.Url)

			if pastPrice < float64(el.Price_rur) {
				rebate := pkg.CalculatePercentageDifference(pastPrice, float64(el.Price_rur))

				if rebate >= 20 {
					h.Logger.Info("Finding the updated price", zap.Int("Price_rur", el.Price_rur))
					chat := h.SubsRepository.SearchChatID(v.Query)
					goodDiscounts := cronhandler.ProductDiscounts(el.Name, float64(el.Price_rur), el.Url, el.Image)

					for _, v := range goodDiscounts.Items {
						text := fmt.Sprintf(
							"*%s*\n"+
								"*Rub* _%v_\n"+
								"*Ссылка* _%s_\n",
							v.Name, v.Price_rur, v.Url)
						chatID := chat
						respy, err := http.Get("https:" + v.Image)
						if err != nil {
							log.Println("Ошибка при получении изображения:", err)
							continue
						}
						defer respy.Body.Close()

						reader := tgbotapi.FileReader{Name: "photo.jpg", Reader: respy.Body}
						photoMsg := tgbotapi.NewPhoto(chatID, reader)
						_, err = h.Bot.Send(photoMsg)
						if err != nil {
							log.Println("Ошибка при отправке сообщения боту:", err)
						}

						msg := tgbotapi.NewMessage(chatID, text)
						msg.ParseMode = "markdown"
						h.Logger.Info("Product withdrawal with a discount")
						_, err = h.Bot.Send(msg)
					}
				}
			}
		}
	}
}
