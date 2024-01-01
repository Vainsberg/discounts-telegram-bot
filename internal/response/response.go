package response

import "github.com/Vainsberg/discounts-telegram-bot/internal/dto"

type RequestDiscounts struct {
	Items []dto.Item `json:"items"`
}

type ProductDiscount struct {
	Items []dto.Item `json:"items"`
}

type SubscriptionRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type TextMessage struct {
	Items []dto.Item `json:"items"`
}
