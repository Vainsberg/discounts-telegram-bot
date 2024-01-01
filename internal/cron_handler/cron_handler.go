package cronhandler

import (
	"net/http"

	"go.uber.org/zap"
)

func HandleCron() {
	var logger zap.Logger
	resp, err := http.Post("http://localhost:8080/discount/update", "application/json", nil)
	if err != nil {
		logger.Info("Ошибка при выполнении запроса:", zap.Error(err))
		return
	}
	defer resp.Body.Close()
}
