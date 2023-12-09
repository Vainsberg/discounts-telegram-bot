package pkg

import (
	"fmt"
	"net/http"
)

func GetQuery(text string) string {
	if text == "" {
		fmt.Println(http.StatusBadRequest)
		return "error"
	}
	return text
}
