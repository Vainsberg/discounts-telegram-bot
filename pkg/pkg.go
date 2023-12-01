package pkg

import (
	"time"

	"github.com/Vainsberg/discounts-telegram-bot/response"
)

var responseone response.RequestDiscounts

var responsetwo response.ResponseQueryandTime

// var Name string
// var Price_ru int
// var Url string
// var Image string
var Query string
var Dt time.Time

func Scan() {

	// 	for _, v := range responseone.Items {
	// 		Name = v.Name
	// 		Price_ru = v.Price_rur
	// 		Url = v.Url
	// 		Image = v.Image

	// 	}
	Query = responsetwo.Query
	Dt = responsetwo.Dt

}
