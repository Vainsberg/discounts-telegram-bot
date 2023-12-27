package bottg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Vainsberg/discounts-telegram-bot/internal/pkg"
	"github.com/Vainsberg/discounts-telegram-bot/internal/response"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleRequest(bot *tgbotapi.BotAPI, message *tgbotapi.Message, update *tgbotapi.Update) {
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
	checkUserText := pkg.Check(userText)

	resp, err := http.Get("http://localhost:8080/discount?query=" + checkUserText + "&response=json")
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
	var result response.TextMessage

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Ошибка при разборе JSON:", err)
		return
	}
	var text string
	slice := result.Items
	for _, v := range slice {
		text = fmt.Sprintf(
			"*%s*\n"+
				"*Rub* _%v_\n"+
				"*Ссылка* _%s_\n",
			v.Name, v.Price_rur, v.Url)
		chatID := update.Message.Chat.ID
		respy, err := http.Get("https:" + v.Image)
		if err != nil {
			log.Println("Ошибка при получении изображения:", err)
			continue
		}
		defer respy.Body.Close()

		reader := tgbotapi.FileReader{Name: "photo.jpg", Reader: respy.Body}
		photoMsg := tgbotapi.NewPhoto(chatID, reader)
		_, err = bot.Send(photoMsg)
		if err != nil {
			log.Println("Ошибка при отправке сообщения боту:", err)
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		msg.ParseMode = "markdown"
		_, err = bot.Send(msg)
	}
}

func HandleCallback(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	callback := update.CallbackQuery
	if callback.Data != "" {
		userText := callback.Data
		chatID := callback.Message.Chat.ID
		ApiURL := "http://localhost:8080/subscribe"

		payload := response.SubscriptionRequest{ChatID: chatID, Text: userText}
		payloadmash, err := json.Marshal(payload)
		if err != nil {
			fmt.Errorf("Ошибка Marshal JSON: %s", err)
			return
		}

		resp, err := http.Post(ApiURL, "application/json", bytes.NewBuffer(payloadmash))
		if err != nil {
			fmt.Println("Ошибка при выполнении запроса:", err)
			return
		}
		defer resp.Body.Close()

		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Вы подписались на товар")
		if _, err := bot.Request(callback); err != nil {
			fmt.Errorf("Ошибка отправки сообщения callback: %s", err)
			return
		}
	}
}
