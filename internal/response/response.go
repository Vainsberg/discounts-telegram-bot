package response

type RequestDiscounts struct {
	Items []struct {
		Name      string `json:"name"`
		Price_rur int    `json:"price_rur"`
		Url       string `json:"url"`
		Image     string `json:"image"`
	} `json:"items"`
}

type ResponseQuery struct {
	Query string `json:"query"`
}

type ProductDiscount struct {
	Items []struct {
		Name      string `json:"name"`
		Price_rur int    `json:"price_rur"`
		Url       string `json:"url"`
		Image     string `json:"image"`
	} `json:"items"`
}

type SubscriptionRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type TextMessage struct {
	Items []struct {
		Name      string `json:"name"`
		Price_rur int    `json:"price_rur"`
		Url       string `json:"url"`
		Image     string `json:"image"`
	} `json:"items"`
}
