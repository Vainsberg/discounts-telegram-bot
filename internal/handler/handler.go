package handler

import (
	"fmt"
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/client"
	"github.com/Vainsberg/discounts-telegram-bot/internal/pkg"
	"github.com/Vainsberg/discounts-telegram-bot/internal/response"
	"github.com/Vainsberg/discounts-telegram-bot/repository"
)

type Handler struct {
	DiscountsRepository repository.Repository
}

func NewHandler(repos *repository.Repository) *Handler {
	return &Handler{DiscountsRepository: *repos}
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
		client.Client(queryText, w, r, responseN)
		for _, v := range responseN.Items {
			h.DiscountsRepository.GetDiscountsAddendumByGoods(v.Name, float64(v.Price_rur), v.Url, v.Image, queryText)
		}

	} else if len(responseN.Items) != 0 {
		var response response.RequestDiscounts
		client.Client(queryText, w, r, response)

	}
}
