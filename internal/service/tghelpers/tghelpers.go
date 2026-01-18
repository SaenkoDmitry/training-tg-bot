package tghelpers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PinMessage(bot *tgbotapi.BotAPI, chatID int64, message tgbotapi.Message) {
	pinChatMessageConfig := tgbotapi.PinChatMessageConfig{
		ChatID:              chatID,
		MessageID:           message.MessageID,
		DisableNotification: false,
	}
	if _, err := bot.Request(pinChatMessageConfig); err != nil {
		fmt.Println("cannot pin message:", message.MessageID)
	}
}

func UnpinMessage(bot *tgbotapi.BotAPI, chatID int64, message tgbotapi.Message) {
	unpinChatMessageConfig := tgbotapi.UnpinChatMessageConfig{
		ChatID:    chatID,
		MessageID: message.MessageID,
	}
	if _, err := bot.Request(unpinChatMessageConfig); err != nil {
		fmt.Println("cannot pin message:", message.MessageID)
	}
}

func SendMessage(bot *tgbotapi.BotAPI, msg tgbotapi.Chattable, method string) (tgbotapi.Message, error) {
	message, err := bot.Send(msg)
	if err != nil {
		fmt.Printf("%s: SendMessage: error: %v", method, err)
	}
	return message, err
}
