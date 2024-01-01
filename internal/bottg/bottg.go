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
	var product response.TextMessage
	var text string
	if message.Text == "" {
		errorMsg := "Неправильная форма заполнения. Пожалуйста, введите нормальное название товара."
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMsg)
		_, err := bot.Send(reply)
		if err != nil {
			log.Println("Ошибка при отправке сообщения боту:", err)
		}
		return
	}
	RepaleseUserText := pkg.ReplaceSpaceUrl(message.Text)

	resp, err := http.Get("http://localhost:8080/discount?query=" + RepaleseUserText + "&response=json")
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

	err = json.Unmarshal(body, &product)
	if err != nil {
		log.Println("Ошибка при разборе JSON:", err)
		return
	}

	for _, v := range product.Items {
		text = fmt.Sprintf(
			"*%s*\n"+
				"*Rub* _%v_\n"+
				"*Ссылка* _%s_\n",
			v.Name, v.Price_rur, v.Url)
		respy, err := http.Get("https:" + v.Image)
		if err != nil {
			log.Println("Ошибка при получении изображения:", err)
			continue
		}
		defer respy.Body.Close()

		reader := tgbotapi.FileReader{Name: "photo.jpg", Reader: respy.Body}
		photoMsg := tgbotapi.NewPhoto(update.Message.Chat.ID, reader)
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
		payload := response.SubscriptionRequest{ChatID: callback.Message.Chat.ID, Text: callback.Data}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			fmt.Errorf("Ошибка Marshal JSON: %s", err)
			return
		}

		resp, err := http.Post("http://localhost:8080/subscribe", "application/json", bytes.NewBuffer(payloadBytes))
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

func RunBot(bot *tgbotapi.BotAPI) {
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message != nil {
				HandleRequest(bot, update.Message, &update)
				userText := update.Message.Text
				replyKeyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Подписаться на скидки", userText),
					),
				)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нажми на кнопку для подписки на товар :)")
				msg.ReplyMarkup = replyKeyboard
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			} else if update.CallbackQuery != nil {
				HandleCallback(bot, &update)
			}
		}
	}()
}
