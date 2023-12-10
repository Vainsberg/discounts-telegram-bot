package pkg

import (
	"errors"
	"fmt"
	"net/http"
)

func GetQuery(text string) (string, error) {
	if text == "" {
		fmt.Println(http.StatusBadRequest)
		return "", errors.New("empty query")
	}
	return text, nil
}
