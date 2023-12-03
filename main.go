package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/internal/pkg"
	"github.com/Vainsberg/discounts-telegram-bot/internal/viper"
	"github.com/Vainsberg/discounts-telegram-bot/repository"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Handler struct {
	discountsRepository repository.Repository
}

func NewHandler(repos *repository.Repository) *Handler {
	return &Handler{discountsRepository: *repos}
}
func (h *Handler) GetDiscounts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	queryText := query.Get("query")
	if queryText == "" {
		fmt.Println(http.StatusBadRequest, w)
		return
	}
	pkg.Check(queryText)

	responseN := h.discountsRepository.GetDiscountsByGoods(queryText)

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
			h.discountsRepository.GetDiscountsAddendumByGoods(v.Name, float64(v.Price_rur), v.Url, v.Image, queryText)
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
	var cfg viper.Config
	var err error
	db, err = sql.Open("mysql", cfg.DbUser+":"+cfg.DbPass+"@tcp(127.0.0.1:3306)/discounts")
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

	http.HandleFunc("/discount", GetDiscounts)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
}
