package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RequestDiscounts struct {
	Items struct {
		Name      string `json:"name"`
		Price_rur int    `json:"price_rur"`
		Url       string `json:"url"`
		Image     string `json:"image"`
	} `json:"items"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	var RequestDisc RequestDiscounts

	query := r.URL.Query()
	queryText := query.Get("query")
	responseText := query.Get("response")

	if queryText == "" {

		fmt.Println(http.StatusBadRequest, w)
	}

	ApiPlatiRu := "https://plati.io/api/search.ashx?query=" + queryText + "&response=" + responseText
	ApiPlatiRuResp, err := http.Get(ApiPlatiRu)
	if err != nil {

		fmt.Println(http.StatusBadRequest, w)
		fmt.Println(err)

	}

	resp, err := io.ReadAll(ApiPlatiRuResp.Body)
	if err != nil {

		fmt.Println(http.StatusBadRequest, w)
		fmt.Println(err)

	}
	defer r.Body.Close()

	err = json.Unmarshal(resp, &RequestDisc)
	if err != nil {

		fmt.Println(http.StatusBadRequest, w)
		fmt.Println(err)

	}
	respText := fmt.Sprintf("Игра: %s, цена: %d, ссылка: %s, %s ", RequestDisc.Items.Name, RequestDisc.Items.Price_rur, RequestDisc.Items.Url, RequestDisc.Items.Image)
	w.Write([]byte(respText))
}

func main() {
	http.HandleFunc("/discount", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
