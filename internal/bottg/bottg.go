package bottg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleRequest(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	type TextMessage struct {
		Items []struct {
			Name      string `json:"name"`
			Price_rur int    `json:"price_rur"`
			Url       string `json:"url"`
			Image     string `json:"image"`
		} `json:"items"`
	}

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
	var result TextMessage

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Ошибка при разборе JSON:", err)
		return
	}
	var text string
	for _, v := range result.Items {
		text = fmt.Sprintf(
			"*%s*\n"+
				"*Rub* _%v_\n"+
				"*Ссылка* _%s_\n"+
				"*Фото* _%s_\n",
			v.Name, v.Price_rur, v.Url, v.Image)
		reply := tgbotapi.NewMessage(message.Chat.ID, text)
		reply.ParseMode = "Markdown"
		_, err = bot.Send(reply)
		if err != nil {
			log.Println("Ошибка при отправке сообщения боту:", err)
		}
	}
}
