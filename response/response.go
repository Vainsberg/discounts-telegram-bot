package response

import "time"

type RequestDiscounts struct {
	Items []struct {
		Name      string `json:"name"`
		Price_rur int    `json:"price_rur"`
		Url       string `json:"url"`
		Image     string `json:"image"`
	} `json:"items"`
}

type ResponseQueryandTime struct {
	Query string    `json:"query"`
	Dt    time.Time `json:"dt"`
}
