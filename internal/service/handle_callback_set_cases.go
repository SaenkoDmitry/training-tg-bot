package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
)

func (s *serviceImpl) setCases(data string, chatID int64) {
	switch {
	case strings.HasPrefix(data, "set_complete_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "set_complete_"), 10, 64)
		s.completeExerciseSet(chatID, exerciseID)

	case strings.HasPrefix(data, "set_add_one_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "set_add_one_"), 10, 64)
		s.addOneMoreSet(chatID, exerciseID)

	case strings.HasPrefix(data, "set_remove_last_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "set_remove_last_"), 10, 64)
		s.removeLastSet(chatID, exerciseID)
	}
}

func (s *serviceImpl) removeLastSet(chatID int64, exerciseID int64) {
	exercise, err := s.exercisesRepo.Get(exerciseID)
	if err != nil || len(exercise.Sets) == 0 {
		return
	}
	if len(exercise.Sets) == 1 {
		msg := tgbotapi.NewMessage(chatID, "Нельзя удалить единственный подход, удалите упражнение целиком кликом на 🗑")
		msg.ParseMode = "Html"
		s.bot.Send(msg)
		return
	}

	lastSet := exercise.Sets[len(exercise.Sets)-1]
	err = s.setsRepo.Delete(lastSet.ID)
	if err != nil {
		fmt.Println("cannot remove set:", err.Error())
		return
	}

	msg := tgbotapi.NewMessage(chatID, "✅ <b>Подход удален!</b>")
	msg.ParseMode = "Html"
	s.bot.Send(msg)

	s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
}

func (s *serviceImpl) addOneMoreSet(chatID int64, exerciseID int64) {
	exercise, err := s.exercisesRepo.Get(exerciseID)
	if err != nil || len(exercise.Sets) == 0 {
		return
	}
	lastSet := exercise.Sets[len(exercise.Sets)-1]
	err = s.setsRepo.Save(&models.Set{
		ExerciseID: exercise.ID,
		Reps:       lastSet.Reps,
		Weight:     lastSet.Weight,
		Minutes:    lastSet.Minutes,
		Meters:     lastSet.Meters,
		Index:      lastSet.Index + 1,
	})
	if err != nil {
		fmt.Println("cannot create set:", err.Error())
		return
	}

	msg := tgbotapi.NewMessage(chatID, "✅ <b>Еще один подход добавлен!</b>")
	msg.ParseMode = "Html"
	s.bot.Send(msg)

	s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
}

func (s *serviceImpl) completeExerciseSet(chatID int64, exerciseID int64) {
	exercise, _ := s.exercisesRepo.Get(exerciseID)

	nextSet := exercise.NextSet()

	if nextSet.ID != 0 {
		nextSet.Completed = true
		now := time.Now()
		nextSet.CompletedAt = &now
		s.setsRepo.Save(&nextSet)
	} else {
		s.moveToNextExercise(chatID, exercise.WorkoutDayID)
		return
	}

	text := fmt.Sprintf("✅ *Подход завершен!*\n\n")
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	s.bot.Send(msg)

	exerciseType, err := s.exerciseTypesRepo.Get(exercise.ExerciseTypeID)
	if err != nil {
		s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
		return
	}

	if nextSet.ID == exercise.LastSet().ID {
		s.moveToNextExercise(chatID, exercise.WorkoutDayID)
	} else {
		s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
	}

	if exerciseType.RestInSeconds > 0 {
		s.startRestTimerWithExercise(chatID, exerciseType.RestInSeconds, exerciseID)
	}
}
