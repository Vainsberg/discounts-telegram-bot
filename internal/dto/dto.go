package dto

type Item struct {
	Name      string `json:"name"`
	Price_rur int    `json:"price_rur"`
	Url       string `json:"url"`
	Image     string `json:"image"`
}
