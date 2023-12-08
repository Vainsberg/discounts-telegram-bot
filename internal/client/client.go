package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/internal/response"
)

type PlatiClient struct {
	BaseURL string
}

func NewPlatiClient(baseURL string) *PlatiClient {
	return &PlatiClient{BaseURL: baseURL}
}

func (c *PlatiClient) GetGoodsClient(queryText string) (*response.RequestDiscounts, error) {
	url := fmt.Sprintf("%s/api/search.ashx?query=%s&response=json", c.BaseURL, queryText)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("plati ru api get request error: %s, %s", url, err)
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("plati ru api read body error: %s, %s", url, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var discounts response.RequestDiscounts
	err = json.Unmarshal(r, &discounts)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %s", err)
	}

	return &discounts, nil
}

func DateFromDatebase(response response.RequestDiscounts) []byte {
	respText, err := json.Marshal(response)
	if err != nil {
		fmt.Errorf("Error encoding JSON response %s", err)
		panic(fmt.Sprintf("Error encoding JSON response: %s", err))
	}
	return respText
}
