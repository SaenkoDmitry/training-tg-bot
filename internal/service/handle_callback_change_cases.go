package service

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
)

func (s *serviceImpl) changeCases(data string, chatID int64) {
	switch {
	case strings.HasPrefix(data, "change_reps_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_reps_ex_"), 10, 64)
		s.askForNewReps(chatID, exerciseID)

	case strings.HasPrefix(data, "change_weight_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_weight_ex_"), 10, 64)
		s.askForNewWeight(chatID, exerciseID)

	case strings.HasPrefix(data, "change_minutes_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_minutes_ex_"), 10, 64)
		s.askForNewMinutes(chatID, exerciseID)

	case strings.HasPrefix(data, "change_meters_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_meters_ex_"), 10, 64)
		s.askForNewMeters(chatID, exerciseID)

	case strings.HasPrefix(data, "change_day_name_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_day_name_"), 10, 64)
		s.askForNewDayName(chatID, programID)
	}
}

func (s *serviceImpl) askForNewReps(chatID int64, exerciseID int64) {
	s.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_reps_%d", exerciseID))
	msg := tgbotapi.NewMessage(chatID, messages.EnterNewReps)
	msg.ParseMode = "Html"
	s.bot.Send(msg)
}

func (s *serviceImpl) askForNewWeight(chatID int64, exerciseID int64) {
	s.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_weight_%d", exerciseID))
	msg := tgbotapi.NewMessage(chatID, messages.EnterNewWeight)
	msg.ParseMode = "Html"
	s.bot.Send(msg)
}

func (s *serviceImpl) askForNewMinutes(chatID int64, exerciseID int64) {
	s.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_minutes_%d", exerciseID))
	msg := tgbotapi.NewMessage(chatID, messages.EnterNewTime)
	msg.ParseMode = "Html"
	s.bot.Send(msg)
}

func (s *serviceImpl) askForNewMeters(chatID int64, exerciseID int64) {
	s.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_meters_%d", exerciseID))
	msg := tgbotapi.NewMessage(chatID, messages.EnterNewMeters)
	msg.ParseMode = "Html"
	s.bot.Send(msg)
}

func (s *serviceImpl) askForNewDayName(chatID, programID int64) {
	s.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_day_name_for_program_%d", programID))
	msg := tgbotapi.NewMessage(chatID, messages.EnterWorkoutDayName)
	msg.ParseMode = "Html"
	s.bot.Send(msg)
}

func (s *serviceImpl) askForNewProgramName(chatID, programID int64) {
	s.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_program_name_%d", programID))
	msg := tgbotapi.NewMessage(chatID, messages.EnterNewProgramName)
	msg.ParseMode = "Html"
	s.bot.Send(msg)
}

func (s *serviceImpl) askForPreset(chatID, dayTypeID, exerciseTypeID int64) {
	s.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_day_preset_%d_%d", dayTypeID, exerciseTypeID))

	exerciseType, _ := s.exerciseTypesRepo.Get(exerciseTypeID)
	exerciseTypeUnits := constants.RepsUnit + "," + constants.WeightUnit
	if exerciseType.Units != "" {
		exerciseTypeUnits = exerciseType.Units
	}
	msg := tgbotapi.NewMessage(chatID, messages.EnterPreset+
		fmt.Sprintf("\n\n<b>Подсказка:</b> для вашего упражнения следует выбрать <b>%s</b> !", exerciseTypeUnits))

	msg.ParseMode = "Html"
	s.bot.Send(msg)
}
