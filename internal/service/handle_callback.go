package service

import (
	"errors"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/tghelpers"
	"strings"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *serviceImpl) HandleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	user, err := s.GetUserByChatID(chatID)
	if err != nil {
		return
	}

	fmt.Println("HandleCallback:", data)

	switch {
	case data == "back_to_menu":
		s.sendMainMenu(chatID, callback.From)

	case strings.HasPrefix(data, "program_"):
		s.programCases(data, chatID, user.ID)

	case strings.HasPrefix(data, "workout_"):
		s.workoutCases(data, chatID, user.ID)

	case strings.HasPrefix(data, "timer_"):
		s.timerCases(data, chatID)

	case strings.HasPrefix(data, "set_"):
		s.setCases(data, chatID)

	case strings.HasPrefix(data, "exercise_"):
		s.exerciseCases(data, chatID)

	case strings.HasPrefix(data, "change_"):
		s.changeCases(data, chatID)

	case strings.HasPrefix(data, "stats_"):
		s.statsCases(data, chatID, user.ID)
	}
}

func (s *serviceImpl) selectExerciseForCurrentWorkout(chatID int64, workoutID int64, exerciseGroupCode string) {
	method := "selectExerciseForCurrentWorkout"
	group, err := s.exerciseGroupTypesRepo.Get(exerciseGroupCode)
	if err != nil {
		return
	}

	text := fmt.Sprintf("*Тип:* %s \n\n *Выберите упражнение из списка:*", group.Name)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	exerciseTypes, err := s.exerciseTypesRepo.GetAllByGroup(exerciseGroupCode)
	if err != nil {
		return
	}

	for _, exercise := range exerciseTypes {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				exercise.Name,
				fmt.Sprintf("exercise_add_specific_for_current_workout_%d_%d", workoutID, exercise.ID),
			),
		))
	}
	fmt.Println("rows", len(rows), rows)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.MarkdownParseMode
	msg.ReplyMarkup = keyboard
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
}

func (s *serviceImpl) GetUserByChatID(chatID int64) (*models.User, error) {
	method := "GetUserByChatID"
	user, err := s.usersRepo.GetByChatID(chatID)
	if err != nil {
		if errors.Is(err, users.NotFoundUserErr) {
			msg := tgbotapi.NewMessage(chatID, "Сначала создайте пользователя в боте, через команду /start")
			_, _ = tghelpers.SendMessage(s.bot, msg, method)
		}
	}
	return user, nil
}
