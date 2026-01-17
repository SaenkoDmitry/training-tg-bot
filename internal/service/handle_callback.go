package service

import (
	"errors"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"strings"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *serviceImpl) HandleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	fmt.Println("HandleCallback:", data)

	switch {
	case data == "back_to_menu":
		s.sendMainMenu(chatID, callback.From)

	case strings.HasPrefix(data, "program_"):
		s.programCases(data, chatID)

	case strings.HasPrefix(data, "workout_"):
		s.workoutCases(data, chatID)

	case strings.HasPrefix(data, "timer_"):
		s.timerCases(data, chatID)

	case strings.HasPrefix(data, "set_"):
		s.setCases(data, chatID)

	case strings.HasPrefix(data, "exercise_"):
		s.exerciseCases(data, chatID)

	case strings.HasPrefix(data, "change_"):
		s.changeCases(data, chatID)

	case strings.HasPrefix(data, "stats_"):
		s.statsCases(data, chatID)
	}
}

func (s *serviceImpl) selectExerciseForCurrentWorkout(chatID int64, workoutID int64, exerciseGroupCode string) {
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
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) showStatistics(chatID int64, period string) {
	method := "showStatistics"
	user, err := s.GetUserByChatID(chatID)
	if err != nil {
		return
	}

	text := s.statisticsService.ShowPeriodStatistics(user.ID, period)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Html"
	_, err = s.bot.Send(msg)
	handleErr(method, err)
}

func (s *serviceImpl) GetUserByChatID(chatID int64) (*models.User, error) {
	user, err := s.usersRepo.GetByChatID(chatID)
	if err != nil {
		if errors.Is(err, users.NotFoundUserErr) {
			msg := tgbotapi.NewMessage(chatID, "Сначала создайте пользователя в боте, через команду /start")
			_, err = s.bot.Send(msg)
			if err != nil {
				fmt.Printf("Error is: %v\n", err)
			}
		}
	}
	return user, nil
}
