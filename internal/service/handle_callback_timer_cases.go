package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
)

func (s *serviceImpl) timerCases(data string, chatID int64) {
	switch {
	case strings.HasPrefix(data, "timer_unpin_and_cancel_"):
		timerID := strings.TrimPrefix(data, "timer_unpin_and_cancel_")
		s.unpinAndCancelTimer(chatID, timerID)

	case strings.HasPrefix(data, "timer_start_"):
		parts := strings.Split(data, "_")
		if len(parts) >= 5 && parts[3] == "ex" {
			seconds, _ := strconv.Atoi(parts[2])
			exerciseID, _ := strconv.ParseInt(parts[4], 10, 64)
			s.startRestTimerWithExercise(chatID, seconds, exerciseID)
		}
	}
}

func (s *serviceImpl) unpinAndCancelTimer(chatID int64, timerID string) {
	s.timerStore.StopTimer(chatID, timerID)
}

func (s *serviceImpl) startRestTimerWithExercise(chatID int64, seconds int, exerciseID int64) {
	if seconds == 0 {
		msg := tgbotapi.NewMessage(chatID, "У этого упражнения не предусмотрен отдых! 😐")
		msg.ParseMode = "Html"
		s.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(messages.RestTimer, seconds))
	newTimerID := s.timerStore.NewTimer(chatID)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отменить", fmt.Sprintf("timer_unpin_and_cancel_%s", newTimerID)),
		),
	)
	msg.ReplyMarkup = keyboard

	var message tgbotapi.Message
	message, _ = s.bot.Send(msg)
	s.pinMessage(chatID, message)

	go func() {
		remaining := seconds

		for remaining > 0 {
			time.Sleep(1 * time.Second)
			remaining--
			if !s.timerStore.HasTimer(chatID, newTimerID) {
				s.unpinMessage(chatID, message)
				editMsg := tgbotapi.NewEditMessageText(chatID, message.MessageID, "Таймер отменен")
				s.bot.Send(editMsg)
				return
			}

			var err error
			if remaining%10 == 0 || remaining <= 20 {
				editMsg := tgbotapi.NewEditMessageTextAndMarkup(chatID, message.MessageID,
					fmt.Sprintf(messages.RestTimer, remaining), keyboard)
				if message, err = s.bot.Send(editMsg); err != nil {
					fmt.Println("cannot edit msg")
				}
			}
		}

		editMsg := tgbotapi.NewEditMessageText(
			chatID,
			message.MessageID,
			"🔔 *Время отдыха закончилось!*\n\n Приступайте к следующему подходу! 💪",
		)
		editMsg.ParseMode = "Markdown"
		editMessage, _ := s.bot.Send(editMsg)
		s.unpinMessage(chatID, editMessage)

		exercise, _ := s.exercisesRepo.Get(exerciseID)

		s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
	}()
}

func (s *serviceImpl) pinMessage(chatID int64, message tgbotapi.Message) {
	pinChatMessageConfig := tgbotapi.PinChatMessageConfig{
		ChatID:              chatID,
		MessageID:           message.MessageID,
		DisableNotification: false,
	}
	if _, err := s.bot.Request(pinChatMessageConfig); err != nil {
		fmt.Println("cannot pin message:", message.MessageID)
	}
}

func (s *serviceImpl) unpinMessage(chatID int64, message tgbotapi.Message) {
	unpinChatMessageConfig := tgbotapi.UnpinChatMessageConfig{
		ChatID:    chatID,
		MessageID: message.MessageID,
	}
	if _, err := s.bot.Request(unpinChatMessageConfig); err != nil {
		fmt.Println("cannot pin message:", message.MessageID)
	}
}
