package cronhandler

import (
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/internal/response"
	"go.uber.org/zap"
)

func HandleCron() {
	var logger zap.Logger
	resp, err := http.Post("http://localhost:8080/discount/update", "application/json", nil)
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
