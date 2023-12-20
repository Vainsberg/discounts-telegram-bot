package dto

type Item struct {
	Name      string `json:"name"`
	Price_rur int    `json:"price_rur"`
	Url       string `json:"url"`
	Image     string `json:"image"`
}

type LinkedAndAccounts struct {
	Name      string
	Price_rur int
	Url       string
	Image     string
}

type GetQuerysItem struct {
	Query string `json:"query"`
}

type GetQuerys struct {
	Items []struct {
		Query string `json:"query"`
	} `json:"items"`
}
