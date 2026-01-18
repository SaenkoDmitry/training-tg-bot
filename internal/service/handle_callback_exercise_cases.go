package service

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/tghelpers"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
)

func (s *serviceImpl) exerciseCases(data string, chatID int64) {
	switch {
	case strings.HasPrefix(data, "exercise_move_to_prev_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "exercise_move_to_prev_"), 10, 64)
		s.moveToPrevExercise(chatID, workoutDayID)

	case strings.HasPrefix(data, "exercise_show_current_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "exercise_show_current_"), 10, 64)
		s.showCurrentExerciseSession(chatID, workoutDayID)

	case strings.HasPrefix(data, "exercise_move_to_next_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "exercise_move_to_next_"), 10, 64)
		s.moveToNextExercise(chatID, workoutDayID)

	case strings.HasPrefix(data, "exercise_show_hint_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "exercise_show_hint_"), 10, 64)
		s.showExerciseHint(chatID, workoutID)

	case strings.HasPrefix(data, "exercise_add_for_current_workout_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "exercise_add_for_current_workout_"), 10, 64)
		s.addExercise(chatID, workoutDayID)

	case strings.HasPrefix(data, "exercise_select_for_current_workout_"):
		text := strings.TrimPrefix(data, "exercise_select_for_current_workout_")
		if arr := strings.Split(text, "_"); len(arr) == 2 {
			workoutDayID, _ := strconv.ParseInt(arr[0], 10, 64)
			code := arr[1]
			s.selectExerciseForCurrentWorkout(chatID, workoutDayID, code)
		}

	case strings.HasPrefix(data, "exercise_add_specific_for_current_workout_"):
		text := strings.TrimPrefix(data, "exercise_add_specific_for_current_workout_")
		if arr := strings.Split(text, "_"); len(arr) == 2 {
			workoutID, _ := strconv.ParseInt(arr[0], 10, 64)
			internalExerciseID, _ := strconv.ParseInt(arr[1], 10, 64)
			s.addSpecificExerciseForCurrentWorkout(chatID, workoutID, internalExerciseID)
		}

	case strings.HasPrefix(data, "exercise_confirm_delete_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "exercise_confirm_delete_"), 10, 64)
		s.confirmDeleteExercise(chatID, exerciseID)

	case strings.HasPrefix(data, "exercise_delete_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "exercise_delete_"), 10, 64)
		s.deleteExercise(chatID, exerciseID)

	case strings.HasPrefix(data, "exercise_select_for_program_day_"):
		text := strings.Split(strings.TrimPrefix(data, "exercise_select_for_program_day_"), "_")
		if len(text) < 2 {
			return
		}
		dayTypeID, _ := strconv.ParseInt(text[0], 10, 64)
		exerciseGroupCode := text[1]
		s.selectExerciseForProgramDay(chatID, dayTypeID, exerciseGroupCode)

	case strings.HasPrefix(data, "exercise_add_specific_for_program_day_"):
		text := strings.Split(strings.TrimPrefix(data, "exercise_add_specific_for_program_day_"), "_")
		if len(text) < 2 {
			return
		}
		dayTypeID, _ := strconv.ParseInt(text[0], 10, 64)
		exerciseTypeID, _ := strconv.ParseInt(text[1], 10, 64)
		s.askForPreset(chatID, dayTypeID, exerciseTypeID)
	}
}

func (s *serviceImpl) confirmDeleteExercise(chatID int64, exerciseID int64) {
	method := "confirmDeleteExercise"
	exercise, _ := s.exercisesRepo.Get(exerciseID)

	exerciseObj, err := s.exerciseTypesRepo.Get(exercise.ExerciseTypeID)
	if err != nil {
		return
	}

	text := fmt.Sprintf("🗑️ <b>Удаление упражнения из тренировочного дня</b>\n\n"+
		"Вы уверены, что хотите удалить упражнение:\n"+
		"<b>%s</b>?\n\n"+
		"⚠️ Это действие нельзя отменить!", exerciseObj.Name)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, удалить",
				fmt.Sprintf("exercise_delete_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, отмена",
				fmt.Sprintf("workout_start_%d", exercise.WorkoutDayID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
}

func (s *serviceImpl) deleteExercise(chatID int64, exerciseID int64) {
	exercise, err := s.exercisesRepo.Get(exerciseID)
	if err != nil {
		return
	}

	err = s.exercisesRepo.Delete(exerciseID)
	if err != nil {
		return
	}

	s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
}

func (s *serviceImpl) showCurrentExerciseSession(chatID int64, workoutID int64) {
	method := "showCurrentExerciseSession"
	if workoutID == 0 {
		return
	}

	workoutDay, _ := s.workoutsRepo.Get(workoutID)
	if len(workoutDay.Exercises) == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ В этой тренировке нет упражнений.")
		_, _ = tghelpers.SendMessage(s.bot, msg, method)
		return
	}

	session, _ := s.sessionsRepo.GetByWorkoutID(workoutID)

	exerciseIndex := session.CurrentExerciseIndex
	if exerciseIndex >= len(workoutDay.Exercises) {
		exerciseIndex = 0
	}

	exercise := workoutDay.Exercises[exerciseIndex]

	var text strings.Builder

	exerciseObj, err := s.exerciseTypesRepo.Get(exercise.ExerciseTypeID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "❌ Упражнение не найдено.")
		_, _ = tghelpers.SendMessage(s.bot, msg, method)
		return
	}

	dayType, err := s.dayTypesRepo.Get(workoutDay.WorkoutDayTypeID)
	if err != nil {
		return
	}

	text.WriteString(fmt.Sprintf("<b>%s</b>\n\n", dayType.Name))
	text.WriteString(fmt.Sprintf("<b>Упражнение %d/%d:</b> %s\n\n", exerciseIndex+1, len(workoutDay.Exercises), exerciseObj.Name))
	if exerciseObj.Accent != "" {
		text.WriteString(fmt.Sprintf("<b>Акцент:</b> %s\n\n", exerciseObj.Accent))
	}

	for _, set := range exercise.Sets {
		text.WriteString(set.String(workoutDay.Completed))
	}

	var changeSettingsButtons []tgbotapi.InlineKeyboardButton
	if len(exercise.Sets) > 0 && exercise.Sets[0].Minutes > 0 {
		changeSettingsButtons = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.Minutes, fmt.Sprintf("change_minutes_ex_%d", exercise.ID)),
		)
	}

	if len(exercise.Sets) > 0 && exercise.Sets[0].Meters > 0 {
		changeSettingsButtons = append(changeSettingsButtons,
			tgbotapi.NewInlineKeyboardButtonData(messages.Meters, fmt.Sprintf("change_meters_ex_%d", exercise.ID)),
		)
	}

	if len(exercise.Sets) > 0 && exercise.Sets[0].Reps > 0 {
		changeSettingsButtons = append(changeSettingsButtons,
			tgbotapi.NewInlineKeyboardButtonData(messages.Reps, fmt.Sprintf("change_reps_ex_%d", exercise.ID)),
		)
	}

	if len(exercise.Sets) > 0 && exercise.Sets[0].Weight > 0 {
		changeSettingsButtons = append(changeSettingsButtons,
			tgbotapi.NewInlineKeyboardButtonData(messages.Weight, fmt.Sprintf("change_weight_ex_%d", exercise.ID)),
		)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.DoneSet, fmt.Sprintf("set_complete_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.AddSet, fmt.Sprintf("set_add_one_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.RemoveSet, fmt.Sprintf("set_remove_last_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.Timer, fmt.Sprintf("timer_start_%d_ex_%d", exercise.ExerciseType.RestInSeconds, exercise.ID)),
		),
		changeSettingsButtons,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.Technique, fmt.Sprintf("exercise_show_hint_%d", workoutID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.EndWorkout, fmt.Sprintf("workout_confirm_finish_%d", workoutID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.Prev, fmt.Sprintf("exercise_move_to_prev_%d", workoutID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.Progress, fmt.Sprintf("workout_show_progress_%d", workoutID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.DropExercise, fmt.Sprintf("exercise_confirm_delete_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.Next, fmt.Sprintf("exercise_move_to_next_%d", workoutID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
}

func (s *serviceImpl) moveToExercise(chatID int64, workoutID int64, next bool) {
	method := "moveToExercise"
	session, _ := s.sessionsRepo.GetByWorkoutID(workoutID)

	if session.ID == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ Активная сессия не найдена")
		_, _ = tghelpers.SendMessage(s.bot, msg, method)
		return
	}

	exercises, _ := s.exercisesRepo.FindAllByWorkoutID(workoutID)

	if len(exercises) == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ В тренировке нет упражнений")
		_, _ = tghelpers.SendMessage(s.bot, msg, method)
		return
	}

	if next {
		session.CurrentExerciseIndex++
	} else {
		session.CurrentExerciseIndex--
	}

	if session.CurrentExerciseIndex < 0 {
		session.CurrentExerciseIndex = 0
		msg := tgbotapi.NewMessage(chatID,
			"Более ранних упражнений в этой тренировке нет")
		_, _ = tghelpers.SendMessage(s.bot, msg, method)
		s.showCurrentExerciseSession(chatID, workoutID)
		return
	}

	if session.CurrentExerciseIndex >= len(exercises) {
		session.CurrentExerciseIndex = 0
		msg := tgbotapi.NewMessage(chatID,
			"🎉 Вы завершили все упражнения в этой тренировке!\n\n"+
				"Хотите завершить тренировку или добавить еще упражнения?")

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🏁 Завершить",
					fmt.Sprintf("workout_confirm_finish_%d", workoutID)),
				tgbotapi.NewInlineKeyboardButtonData("➕ Еще упражнение",
					fmt.Sprintf("exercise_add_for_current_workout_%d", workoutID)),
			),
		)

		msg.ReplyMarkup = keyboard
		_, _ = tghelpers.SendMessage(s.bot, msg, method)
		return
	}

	s.sessionsRepo.Save(&session)
	s.showCurrentExerciseSession(chatID, workoutID)
}

func (s *serviceImpl) moveToPrevExercise(chatID int64, workoutID int64) {
	s.moveToExercise(chatID, workoutID, false)
}

func (s *serviceImpl) moveToNextExercise(chatID int64, workoutID int64) {
	s.moveToExercise(chatID, workoutID, true)
}

func (s *serviceImpl) showExerciseHint(chatID int64, workoutID int64) {
	method := "showExerciseHint"
	workoutDay, _ := s.workoutsRepo.Get(workoutID)

	if len(workoutDay.Exercises) == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ В этой тренировке нет упражнений.")
		_, _ = tghelpers.SendMessage(s.bot, msg, method)
		return
	}

	session, _ := s.sessionsRepo.GetByWorkoutID(workoutID)

	exerciseIndex := session.CurrentExerciseIndex
	if exerciseIndex >= len(workoutDay.Exercises) {
		exerciseIndex = 0
	}

	exercise := workoutDay.Exercises[exerciseIndex]

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", fmt.Sprintf("exercise_show_current_%d", workoutID)),
	))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	msg := tgbotapi.NewMessage(chatID, utils.WrapYandexLink(exercise.ExerciseType.Url))
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
}

func (s *serviceImpl) addExercise(chatID int64, workoutID int64) {
	method := "addExercise"
	text := messages.SelectGroupOfMuscle

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)

	groups, err := s.exerciseGroupTypesRepo.GetAll()
	if err != nil {
		return
	}

	for i, group := range groups {
		if i%3 == 0 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow())
		}
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1], tgbotapi.NewInlineKeyboardButtonData(group.Name,
			fmt.Sprintf("exercise_select_for_current_workout_%d_%s", workoutID, group.Code)))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
}

func (s *serviceImpl) addSpecificExerciseForCurrentWorkout(chatID int64, workoutID int64, exerciseTypeID int64) {
	method := "addSpecificExerciseForCurrentWorkout"
	exerciseObj, err := s.exerciseTypesRepo.Get(exerciseTypeID)
	if err != nil {
		return
	}

	fmt.Println("newExercise:", exerciseObj)

	workout, _ := s.workoutsRepo.Get(workoutID)
	idx := 0
	if len(workout.Exercises) > 0 {
		lastExercise := workout.Exercises[len(workout.Exercises)-1]
		idx = lastExercise.Index + 1
	}
	newExercise := models.Exercise{
		ExerciseTypeID: exerciseObj.ID,
		Index:          idx,
		WorkoutDayID:   workoutID,
		Sets: []models.Set{
			{Index: 1}, // по дефолту один подход
		},
	}
	exerciseTypeObj, err := s.exerciseTypesRepo.Get(newExercise.ExerciseTypeID)
	if err == nil {
		if strings.Contains(exerciseTypeObj.Units, "meters") {
			newExercise.Sets[0].Meters = 100
		}
		if strings.Contains(exerciseTypeObj.Units, "minutes") {
			newExercise.Sets[0].Minutes = 1
		}
		if strings.Contains(exerciseTypeObj.Units, "reps") {
			newExercise.Sets[0].Reps = 10
		}
		if strings.Contains(exerciseTypeObj.Units, "weight") {
			newExercise.Sets[0].Weight = 10
		}
	}
	workout.Exercises = append(workout.Exercises, newExercise)

	s.workoutsRepo.Save(&workout)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Упражнение <b>'%s'</b> добавлено! ✅", exerciseTypeObj.Name))
	msg.ParseMode = constants.HtmlParseMode
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
	s.showWorkoutProgress(chatID, workoutID)
}

func (s *serviceImpl) selectExerciseForProgramDay(chatID int64, dayTypeID int64, exerciseGroupCode string) {
	method := "selectExerciseForProgramDay"
	group, err := s.exerciseGroupTypesRepo.Get(exerciseGroupCode)
	if err != nil {
		return
	}

	text := fmt.Sprintf("*Тип: %s \n\n Выберите упражнение из списка:*", group.Name)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	exerciseTypes, err := s.exerciseTypesRepo.GetAllByGroup(exerciseGroupCode)
	if err != nil {
		return
	}

	for _, exercise := range exerciseTypes {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				exercise.Name,
				fmt.Sprintf("exercise_add_specific_for_program_day_%d_%d", dayTypeID, exercise.ID),
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
