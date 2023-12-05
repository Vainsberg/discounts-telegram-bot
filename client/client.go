package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/internal/response"
)

func Client(queryText string, w http.ResponseWriter, r *http.Request, response response.RequestDiscounts) (int, error) {
	var respText []byte
	var err error
	if len(response.Items) == 0 {
		ApiPlatiRu := "https://plati.io/api/search.ashx?query=" + queryText + "&response=json"
		ApiPlatiRuResp, err := http.Get(ApiPlatiRu)
		if err != nil {
			fmt.Println(http.StatusBadRequest)
			fmt.Errorf("plati ru api error: %s", err)
			return 0, err
		}

		resp, err := io.ReadAll(ApiPlatiRuResp.Body)
		if err != nil {
			fmt.Println(http.StatusBadRequest, w)
			fmt.Errorf("plati ru api readAll error: %s", err)
			return 0, err
		}
		defer r.Body.Close()

		err = json.Unmarshal(resp, &response)
		if err != nil {
			fmt.Println(http.StatusBadRequest)
			fmt.Errorf("json.Unmarshal error: %s", err)
			return 0, err
		}

		respText, err = json.Marshal(response)
		if err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			fmt.Errorf("json.Marshal error: %s", err)
			return 0, err
		}

	} else if len(response.Items) != 0 {
		respText, err = json.Marshal(response)
		if err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			fmt.Println(err)
			return 0, err
		}
	}
	return w.Write(respText)
}
