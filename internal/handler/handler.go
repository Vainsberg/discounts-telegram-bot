package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/internal/client"
	pkg "github.com/Vainsberg/discounts-telegram-bot/internal/pkg"
	"github.com/Vainsberg/discounts-telegram-bot/internal/repository"
	"github.com/Vainsberg/discounts-telegram-bot/internal/response"
)

type Handler struct {
	DiscountsRepository  repository.Repository
	DiscountsPlatiClient client.PlatiClient
}

func NewHandler(repos *repository.Repository, plati *client.PlatiClient) *Handler {
	return &Handler{DiscountsRepository: *repos,
		DiscountsPlatiClient: *plati}
}

func (h *Handler) GetDiscounts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	queryText := query.Get("query")
	if queryText == "" {
		fmt.Println(http.StatusBadRequest, w)
		return
	}
	pkg.Check(queryText)
	responseN := h.DiscountsRepository.GetDiscountsByGoods(queryText)

	if len(responseN.Items) == 0 {
		goods, err := h.DiscountsPlatiClient.GetGoodsClient(queryText)
		if err != nil {
			fmt.Errorf("DiscountsPlatiClient error: %s", err)
			return
		}
		for _, v := range goods.Items {
			h.DiscountsRepository.SaveGood(v.Name, float64(v.Price_rur), v.Url, v.Image, queryText)
		}

	} else if len(responseN.Items) != 0 {
		var response response.RequestDiscounts
		respText, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		w.Write(respText)
	}
}
