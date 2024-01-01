package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/internal/client"
	pkg "github.com/Vainsberg/discounts-telegram-bot/internal/pkg"
	"github.com/Vainsberg/discounts-telegram-bot/internal/response"
	"github.com/Vainsberg/discounts-telegram-bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Handler struct {
	Logger  *zap.Logger
	Bot     *tgbotapi.BotAPI
	Service service.Service
}

func NewHandler(logger *zap.Logger, bot *tgbotapi.BotAPI, service service.Service) *Handler {
	return &Handler{
		Logger:  logger,
		Bot:     bot,
		Service: service,
	}
}

func (h *Handler) GetDiscounts(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Launch GetDiscounts")

	query, err := pkg.GetQuery(r.URL.Query().Get("query"))
	if err != nil {
		h.Logger.Info("GetQuery error:", zap.Error(err))
		return
	}
	ReplaceText := pkg.ReplaceSpaceUrl(query)
	discountsByGoods := h.Service.DiscountsRepository.GetDiscountsByGoods(ReplaceText)

	if len(discountsByGoods.Items) != 0 {
		w.Write(client.ConvertRequestDiscountsToJSON(discountsByGoods))
		return
	}

	goods := h.Service.FetchAndSaveGoods(ReplaceText)

	bytes, err := json.Marshal(goods)
	if err != nil {
		h.Logger.Info("Error encoding JSON response", zap.Error(err))
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

func (h *Handler) AddSubscription(w http.ResponseWriter, r *http.Request) {
	var subscriptionRequest response.SubscriptionRequest
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

	err = json.Unmarshal(body, &subscriptionRequest)
	if err != nil {
		h.Logger.Info("Ошибка при разборе JSON:", zap.Error(err))
		return
	}
	h.Service.AddLinked(subscriptionRequest.ChatID, subscriptionRequest.Text)
}

func (h *Handler) GetQuerysCron(w http.ResponseWriter, r *http.Request) {
	err := h.Service.ProcessQueryAndFetchGoods()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
