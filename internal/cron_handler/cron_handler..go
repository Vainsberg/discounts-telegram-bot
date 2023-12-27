package cronhandler

import (
	"fmt"
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/internal/response"
)

func HandleCron() {
	ApiURL := "http://localhost:8080/discount/update"
	resp, err := http.Post(ApiURL, "application/json", nil)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
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
