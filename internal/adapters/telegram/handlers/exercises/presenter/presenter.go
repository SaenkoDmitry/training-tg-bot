package presenter

import (
	"fmt"
	"strings"

	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Presenter struct {
	bot *tgbotapi.BotAPI
}

func NewPresenter(bot *tgbotapi.BotAPI) *Presenter {
	return &Presenter{bot: bot}
}

func (p *Presenter) ShowCurrentSession(chatID int64, res *dto.CurrentExerciseSession) {
	var text strings.Builder

	dayType := res.DayType
	exerciseIndex := res.ExerciseIndex
	workoutDay := res.WorkoutDay
	exercise := res.Exercise
	exerciseObj := res.ExerciseObj
	workoutID := workoutDay.ID

	text.WriteString(fmt.Sprintf("<b>%s</b>\n\n", dayType.Name))
	text.WriteString(fmt.Sprintf("<b>Упражнение %d/%d:</b> %s\n\n", exerciseIndex+1, len(workoutDay.Exercises), exerciseObj.Name))
	if exerciseObj.Accent != "" {
		text.WriteString(fmt.Sprintf("<b>Основные мышцы:</b> %s\n\n", exerciseObj.Accent))
	}
	if exerciseObj.SecondaryAccent != "" {
		text.WriteString(fmt.Sprintf("<b>Второстепенные мышцы:</b> %s\n\n", exerciseObj.SecondaryAccent))
	}

	text.WriteString("<b>Подходы:</b>\n")
	for _, set := range exercise.Sets {
		text.WriteString(set.FormattedString)
	}

	var changeSettingsButtons []tgbotapi.InlineKeyboardButton
	if strings.Contains(exercise.Units, constants.MinutesUnit) {
		changeSettingsButtons = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.Minutes, fmt.Sprintf("change_minutes_ex_%d", exercise.ID)),
		)
	}

	if strings.Contains(exercise.Units, constants.MetersUnit) {
		changeSettingsButtons = append(changeSettingsButtons,
			tgbotapi.NewInlineKeyboardButtonData(messages.Meters, fmt.Sprintf("change_meters_ex_%d", exercise.ID)),
		)
	}

	if strings.Contains(exercise.Units, constants.RepsUnit) {
		changeSettingsButtons = append(changeSettingsButtons,
			tgbotapi.NewInlineKeyboardButtonData(messages.Reps, fmt.Sprintf("change_reps_ex_%d", exercise.ID)),
		)
	}

	if strings.Contains(exercise.Units, constants.WeightUnit) {
		changeSettingsButtons = append(changeSettingsButtons,
			tgbotapi.NewInlineKeyboardButtonData(messages.Weight, fmt.Sprintf("change_weight_ex_%d", exercise.ID)),
		)
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.DoneSet, fmt.Sprintf("set_complete_%d", exercise.ID)),
		tgbotapi.NewInlineKeyboardButtonData(messages.AddSet, fmt.Sprintf("set_add_one_%d", exercise.ID)),
		tgbotapi.NewInlineKeyboardButtonData(messages.RemoveSet, fmt.Sprintf("set_remove_last_%d", exercise.ID)),
		tgbotapi.NewInlineKeyboardButtonData(messages.Timer, fmt.Sprintf("timer_start_%d_ex_%d", exerciseObj.RestInSeconds, exercise.ID)),
	))

	if len(changeSettingsButtons) > 0 {
		rows = append(rows, changeSettingsButtons)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.Technique, fmt.Sprintf("exercise_show_hint_%d_%d", workoutDay.ID, exerciseObj.ID)),
		tgbotapi.NewInlineKeyboardButtonData(messages.EndWorkout, fmt.Sprintf("workout_confirm_finish_%d", workoutID)),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.Prev, fmt.Sprintf("exercise_move_to_prev_%d", workoutID)),
		tgbotapi.NewInlineKeyboardButtonData(messages.Progress, fmt.Sprintf("workout_show_progress_%d", workoutID)),
		tgbotapi.NewInlineKeyboardButtonData(messages.DropExercise, fmt.Sprintf("exercise_confirm_delete_%d", exercise.ID)),
		tgbotapi.NewInlineKeyboardButtonData(messages.Next, fmt.Sprintf("exercise_move_to_next_%d", workoutID)),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard

	p.bot.Send(msg)
}

func (p *Presenter) ShowNoExercises(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "❌ В этой тренировке нет упражнений.")
	p.bot.Send(msg)
}

func (p *Presenter) ShowNotFoundExercise(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "❌ Упражнение не найдено.")
	p.bot.Send(msg)
}

func (p *Presenter) ShowSelectExerciseForProgramDayDialog(chatID, dayTypeID int64, group *dto.Group, exerciseTypes []*dto.ExerciseTypeDTO) {
	text := fmt.Sprintf("<b>Тип:</b> %s \n\n %s", group.Name, messages.SelectExercise)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	for _, exerciseType := range exerciseTypes {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				exerciseType.Name,
				fmt.Sprintf("change_program_day_add_exercise_%d_%d", dayTypeID, exerciseType.ID),
			),
		))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) ShowConfirmDeleteDialog(chatID int64, res *dto.ConfirmDeleteExercise) {
	exercise := res.Exercise
	exerciseObj := res.ExerciseObj
	text := fmt.Sprintf("🗑️ <b>Удаление упражнения из тренировочного дня</b>\n\n"+
		"Вы уверены, что хотите удалить упражнение:\n"+
		"<b>%s</b>?\n\n"+
		"⚠️ Это действие нельзя отменить!", exerciseObj.Name)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.YesDelete,
				fmt.Sprintf("exercise_delete_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.NoCancel,
				fmt.Sprintf("workout_start_%d", exercise.WorkoutDayID)),
		),
	)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) CompleteAllExercises(chatID, workoutID int64) {
	msg := tgbotapi.NewMessage(chatID,
		"🎉 Вы завершили все упражнения в этой тренировке!\n\n"+
			"Хотите завершить тренировку или добавить еще упражнения?")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏁 Завершить",
				fmt.Sprintf("workout_finish_%d", workoutID)),
			tgbotapi.NewInlineKeyboardButtonData("➕ Еще упражнение",
				fmt.Sprintf("exercise_add_for_current_workout_%d", workoutID)),
		),
	)
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) ShowHint(chatID int64, res *dto.GetExerciseType, workoutID int64) {
	exerciseType := res.ExerciseType
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	if workoutID == 0 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.BackTo, fmt.Sprintf("exercise_show_list_%s", exerciseType.ExerciseGroupTypeCode)),
		))
	} else {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.BackTo, fmt.Sprintf("exercise_show_current_session_%d", workoutID)),
		))
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	text := ""
	if exerciseType.Description != "" {
		text += exerciseType.Description + "\n\n"
	}
	if exerciseType.Url != "" {
		text += utils.WrapYandexLink(exerciseType.Url)
	}
	if text == "" {
		text += "¯\\_(ツ)_/¯"
	}
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) AddExerciseDialog(chatID, workoutID int64, groups []dto.Group) {
	text := messages.SelectGroupOfMuscle
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
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
	p.bot.Send(msg)
}

func (p *Presenter) ShowSelectExerciseForCurrentWorkoutDialog(chatID, workoutID int64, group *dto.Group, exerciseTypes []*dto.ExerciseTypeDTO) {
	text := fmt.Sprintf("<b>Тип:</b> %s \n\n %s", group.Name, messages.SelectExercise)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	for _, exercise := range exerciseTypes {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				exercise.Name,
				fmt.Sprintf("exercise_add_specific_for_current_workout_%d_%d", workoutID, exercise.ID),
			),
		))
	}
	fmt.Println("rows", len(rows), rows)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	_, err := p.bot.Send(msg)
	fmt.Println("err:", err)
}

func (p *Presenter) ShowAllGroups(chatID int64, groups []dto.Group) {
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	for i, group := range groups {
		if i%3 == 0 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow())
		}
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1],
			tgbotapi.NewInlineKeyboardButtonData(group.Name, fmt.Sprintf("exercise_show_list_%s", group.Code)),
		)
	}
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.BackTo, "/menu"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg := tgbotapi.NewMessage(chatID, messages.SelectGroupOfMuscle)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) ShowAllExercises(chatID int64, exerciseTypes []*dto.ExerciseTypeDTO, groupName string) {
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, ex := range exerciseTypes {
		buttons = append(buttons,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(ex.Name, fmt.Sprintf("exercise_show_info_%d", ex.ID)),
			),
		)
	}
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.BackTo, "exercise_show_all_groups"),
	))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	text := fmt.Sprintf("<b>Тип:</b> %s \n\n %s", groupName, messages.SelectExercise)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}
