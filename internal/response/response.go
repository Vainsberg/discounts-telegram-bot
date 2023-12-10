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
