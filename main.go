package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/internal/pkg"
	"github.com/Vainsberg/discounts-telegram-bot/internal/response"
	"github.com/Vainsberg/discounts-telegram-bot/internal/viper"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	queryText := query.Get("query")
	if queryText == "" {
		fmt.Println(http.StatusBadRequest, w)
		return
	}
	pkg.Check(queryText)

	rows, err := db.Query("SELECT name, price_ru, url, image FROM goods WHERE query = ? AND dt >= CURRENT_TIMESTAMP() - INTERVAL 24 HOUR;", queryText)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var responseN response.RequestDiscounts
	responseN.Items = []struct {
		Name      string `json:"name"`
		Price_rur int    `json:"price_rur"`
		Url       string `json:"url"`
		Image     string `json:"image"`
	}{}

	for rows.Next() {
		var item struct {
			Name      string `json:"name"`
			Price_rur int    `json:"price_rur"`
			Url       string `json:"url"`
			Image     string `json:"image"`
		}

		err := rows.Scan(&item.Name, &item.Price_rur, &item.Url, &item.Image)
		if err != nil {
			log.Fatal(err)
			return
		}
		responseN.Items = append(responseN.Items, item)
	}

	if len(responseN.Items) == 0 {
		ApiPlatiRu := "https://plati.io/api/search.ashx?query=" + queryText + "&response=json"
		ApiPlatiRuResp, err := http.Get(ApiPlatiRu)
		if err != nil {
			fmt.Println(http.StatusBadRequest, w)
			fmt.Errorf("plati ru api error: %s", err)
			return
		}

		resp, err := io.ReadAll(ApiPlatiRuResp.Body)
		if err != nil {
			fmt.Println(http.StatusBadRequest, w)
			fmt.Errorf("plati ru api readAll error: %s", err)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(resp, &responseN)
		if err != nil {
			fmt.Println(http.StatusBadRequest, w)
			fmt.Errorf("json.Unmarshal error: %s", err)
			return
		}

		respText, err := json.Marshal(responseN)
		if err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			fmt.Errorf("json.Marshal error: %s", err)
			return
		}
		w.Write(respText)

		for _, v := range responseN.Items {
			_, err = db.Exec("INSERT INTO goods (name, price_ru, url, image, dt, query) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP(), ?)", v.Name, v.Price_rur, v.Url, v.Image, queryText)
			if err != nil {
				log.Fatal(err)
				fmt.Println(err)
				return
			}
		}

	} else if len(responseN.Items) != 0 {
		respText, err := json.Marshal(responseN)
		if err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		w.Write(respText)
	}
}

func main() {
	var err error

	db, err = sql.Open("mysql", viper.User+":"+viper.Pass+"@tcp(127.0.0.1:3306)/discounts")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS goods (
			id INTEGER PRIMARY KEY AUTO_INCREMENT,
		    name TEXT,
			price_ru  FLOAT,
			url TEXT,
			image TEXT,
			dt DATETIME DEFAULT CURRENT_TIMESTAMP,
			query TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/discount", handler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
}
