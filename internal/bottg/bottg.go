package bottg

import (
	"io"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleRequest(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if message.Text == "" {
		errorMsg := "Неправильная форма заполнения. Пожалуйста, введите нормальное название товара."
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMsg)
		_, err := bot.Send(reply)
		if err != nil {
			log.Println("Ошибка при отправке сообщения боту:", err)
		}
		return
	}
	userText := message.Text

	resp, err := http.Get("http://localhost:8080/discount?query=" + userText + "&response=json")
	if err != nil {
		log.Println("Ошибка при создании HTTP-запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка при чтении тела ответа:", err)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, string(body))
	_, err = bot.Send(reply)
	if err != nil {
		log.Println("Ошибка при отправке сообщения боту:", err)
	}
}
