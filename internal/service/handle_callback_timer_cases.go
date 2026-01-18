package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/tghelpers"
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
	method := "startRestTimerWithExercise"
	if seconds == 0 {
		msg := tgbotapi.NewMessage(chatID, messages.RestNotSupported)
		msg.ParseMode = constants.HtmlParseMode
		_, _ = tghelpers.SendMessage(s.bot, msg, method)
		return
	}

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(messages.RestTimer, seconds))
	newTimerID := s.timerStore.NewTimer(chatID)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.CancelTimer, fmt.Sprintf("timer_unpin_and_cancel_%s", newTimerID)),
		),
	)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard

	var message tgbotapi.Message
	message, _ = tghelpers.SendMessage(s.bot, msg, method)
	tghelpers.PinMessage(s.bot, chatID, message)

	go func() {
		remaining := seconds

		for remaining > 0 {
			time.Sleep(1 * time.Second)
			remaining--
			if !s.timerStore.HasTimer(chatID, newTimerID) {
				tghelpers.UnpinMessage(s.bot, chatID, message)
				editMsg := tgbotapi.NewEditMessageText(chatID, message.MessageID, messages.TimerCanceled)
				editMsg.ParseMode = constants.HtmlParseMode
				_, _ = tghelpers.SendMessage(s.bot, editMsg, method)
				return
			}

			if remaining%10 == 0 || remaining <= 20 {
				editMsg := tgbotapi.NewEditMessageTextAndMarkup(chatID, message.MessageID,
					fmt.Sprintf(messages.RestTimer, remaining), keyboard)
				editMsg.ParseMode = constants.HtmlParseMode
				message, _ = tghelpers.SendMessage(s.bot, editMsg, method)
			}
		}

		editMsg := tgbotapi.NewEditMessageText(
			chatID,
			message.MessageID,
			messages.RestIsEnded,
		)
		editMsg.ParseMode = constants.HtmlParseMode

		editMessage, err := tghelpers.SendMessage(s.bot, editMsg, method)
		if err != nil {
			return
		}

		tghelpers.UnpinMessage(s.bot, chatID, editMessage)

		exercise, _ := s.exercisesRepo.Get(exerciseID)
		s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
	}()
}
