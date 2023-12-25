package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Vainsberg/discounts-telegram-bot/internal/client"
	pkg "github.com/Vainsberg/discounts-telegram-bot/internal/pkg"
	"github.com/Vainsberg/discounts-telegram-bot/internal/repository"
	"go.uber.org/zap"
)

type Handler struct {
	Logger               *zap.Logger
	DiscountsRepository  repository.Repository
	DiscountsPlatiClient client.PlatiClient
	SubsRepository       repository.RepositorySubs
	RepositoryQuerys     repository.RepositorySubs
}

func NewHandler(logger *zap.Logger, repos *repository.Repository, plati *client.PlatiClient, subs *repository.RepositorySubs, querys *repository.RepositorySubs) *Handler {
	return &Handler{
		Logger:               logger,
		DiscountsRepository:  *repos,
		DiscountsPlatiClient: *plati,
		SubsRepository:       *subs,
		RepositoryQuerys:     *querys,
	}
}

func (h *Handler) GetDiscounts(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Launch GetDiscounts")

	query, err := pkg.GetQuery(r.URL.Query().Get("query"))
	if err != nil {
		fmt.Errorf("GetQuery error: %s", err)
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
		fmt.Errorf("DiscountsPlatiClient error: %s", err)
		return
	}

	for _, v := range goods.Items {
		err := h.DiscountsRepository.SaveGood(v.Name, float64(v.Price_rur), v.Url, v.Image, CheckQueryText)
		if err != nil {
			fmt.Println(err)
		}
	}

	respText, err := json.Marshal(goods)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.Write(respText)
}

type SubscriptionRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func (h *Handler) AddSubscription(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("AddSubscription")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("Read body error: %s", err)
		return

	}

	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	var result SubscriptionRequest

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Ошибка при разборе JSON:", err)
		return
	}
	h.SubsRepository.AddLincked(result.ChatID, result.Text)
}

func (h *Handler) GetQuerysCron(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("GetQuerysCron")

	query, err := h.RepositoryQuerys.GetQuerys()
	if err != nil {
		fmt.Errorf("DiscountsPlatiClient error: %s", err)
		return
	}

	for _, v := range query.Items {
		time.Sleep(5 * time.Second)
		goods, err := h.DiscountsPlatiClient.GetGoodsClient(v.Query)
		if err != nil {
			fmt.Errorf("DiscountsPlatiClient error: %s", err)
			return
		}
		for _, el := range goods.Items {
			h.DiscountsRepository.SaveGood(el.Name, float64(el.Price_rur), el.Url, el.Image, v.Query)
			pastPrice := h.DiscountsRepository.SearchDiscount(float64(el.Price_rur), el.Url)

			if pastPrice < float64(el.Price_rur) {
				rebate := pkg.CalculatePercentageDifference(pastPrice, float64(el.Price_rur))

				if rebate >= 15 && rebate <= 20 {

				}
			}

		}
	}

}
