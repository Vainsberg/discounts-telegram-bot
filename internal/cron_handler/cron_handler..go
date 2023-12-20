package cronhandler

import (
	"fmt"
	"net/http"
)

func HandleCron() {
	ApiURL := "http://localhost:8080/discounts/update"
	resp, err := http.Post(ApiURL, "application/json", nil)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer resp.Body.Close()

}
