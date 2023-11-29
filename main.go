package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RequestDiscounts struct {
	Items []struct {
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

	if queryText == "" {

		fmt.Println(http.StatusBadRequest, w)
	}

	ApiPlatiRu := "https://plati.io/api/search.ashx?query=" + queryText + "&response=json"
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
	respText, err := json.Marshal(RequestDisc)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.Write(respText)
}

func main() {
	http.HandleFunc("/discount", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
