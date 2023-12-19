package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/internal/client"
	pkg "github.com/Vainsberg/discounts-telegram-bot/internal/pkg"
	"github.com/Vainsberg/discounts-telegram-bot/internal/repository"
)

type Handler struct {
	DiscountsRepository  repository.Repository
	DiscountsPlatiClient client.PlatiClient
	SubsRepository       repository.RepositorySubs
}

func NewHandler(repos *repository.Repository, plati *client.PlatiClient, subs *repository.RepositorySubs) *Handler {
	return &Handler{
		DiscountsRepository:  *repos,
		DiscountsPlatiClient: *plati,
		SubsRepository:       *subs,
	}
}

func (h *Handler) GetDiscounts(w http.ResponseWriter, r *http.Request) {
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
