package cronhandler

import (
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/internal/response"
	"go.uber.org/zap"
)

func HandleCron() {
	var logger zap.Logger
	ApiURL := "http://localhost:8080/discount/update"
	resp, err := http.Post(ApiURL, "application/json", nil)
	if err != nil {
		logger.Info("Ошибка при выполнении запроса:", zap.Error(err))
		return
	}
	defer resp.Body.Close()
}

func ProductDiscounts(name string, price_rur float64, url string, image string) response.ProductDiscount {
	var product response.ProductDiscount

	product.Items = append(product.Items, struct {
		Name      string `json:"name"`
		Price_rur int    `json:"price_rur"`
		Url       string `json:"url"`
		Image     string `json:"image"`
	}{
		Name:      name,
		Price_rur: int(price_rur),
		Url:       url,
		Image:     image,
	})
	return product
}
