package service

import (
	"bytes"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/tghelpers"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
)

func (s *serviceImpl) programCases(data string, chatID, userID int64) {
	switch {
	case strings.HasPrefix(data, "program_create"):
		s.createProgram(chatID, userID)

	case strings.HasPrefix(data, "program_management"):
		s.programManagement(chatID)

	case strings.HasPrefix(data, "program_edit_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "program_edit_"), 10, 64)
		s.editProgram(chatID, programID)

	case strings.HasPrefix(data, "program_day_edit_"):
		dayTypeID, _ := strconv.ParseInt(strings.TrimPrefix(data, "program_day_edit_"), 10, 64)
		s.addNewDayTypeExercise(chatID, dayTypeID)

	case strings.HasPrefix(data, "program_change_name_of_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "program_change_name_of_"), 10, 64)
		s.askForNewProgramName(chatID, programID)

	case strings.HasPrefix(data, "program_change_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "program_change_"), 10, 64)
		s.changeProgram(chatID, programID)

	case strings.HasPrefix(data, "program_confirm_delete_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "program_confirm_delete_"), 10, 64)
		s.confirmDeleteProgram(chatID, programID)

	case strings.HasPrefix(data, "program_delete_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "program_delete_"), 10, 64)
		s.deleteProgram(chatID, programID)
	}
}

func (s *serviceImpl) createProgram(chatID, userID int64) {
	programs, err := s.programsRepo.FindAll(userID)
	if err != nil {
		return
	}

	_, err = s.programsRepo.Create(userID, fmt.Sprintf("#%d", len(programs)+1))
	if err != nil {
		return
	}

	s.programManagement(chatID)
}

func (s *serviceImpl) editProgram(chatID int64, programID int64) {
	method := "editProgram"
	program, err := s.programsRepo.Get(programID)
	if err != nil {
		return
	}

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	text := &bytes.Buffer{}

	text.WriteString(fmt.Sprintf("<b>Программа: %s</b>\n\n", program.Name))
	text.WriteString("<b>Список дней:</b>\n\n")
	for i, dayType := range program.DayTypes {
		if i%2 == 0 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow())
		}
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1],
			tgbotapi.NewInlineKeyboardButtonData(dayType.Name, fmt.Sprintf("program_day_edit_%d", dayType.ID)),
		)

		text.WriteString(fmt.Sprintf("<b>%d. %s</b>\n", i+1, dayType.Name))
		text.WriteString(fmt.Sprintf("%s \n\n", s.formatPreset(dayType.Preset)))
	}
	text.WriteString("<b>Выберите день, в который хотите добавить упражнения:</b>")

	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("➕ Добавить день", fmt.Sprintf("change_day_name_%d", programID)),
		tgbotapi.NewInlineKeyboardButtonData("🎟️ Переименовать", fmt.Sprintf("program_change_name_of_%d", programID)),
	))
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👑 Выбрать текущей", fmt.Sprintf("program_change_%d", programID)),
		tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить", fmt.Sprintf("program_confirm_delete_%d", programID)),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
}

func (s *serviceImpl) confirmDeleteProgram(chatID, programID int64) {
	method := "confirmDeleteProgram"

	program, err := s.programsRepo.Get(programID)
	if err != nil {
		return
	}

	text := fmt.Sprintf("🗑️ *Удаление программы*\n\n"+
		"Вы уверены, что хотите удалить программу:\n"+
		"*%s*?\n\n"+
		"⚠️ Это действие нельзя отменить!", program.Name)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, удалить",
				fmt.Sprintf("program_delete_%d", programID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, отмена",
				fmt.Sprintf("program_edit_%d", programID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.MarkdownParseMode
	msg.ReplyMarkup = keyboard
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
}

func (s *serviceImpl) deleteProgram(chatID, programID int64) {
	method := "deleteProgram"
	fmt.Println("deleteProgram")

	user, err := s.usersRepo.GetByChatID(chatID)
	if err != nil {
		s.handleGetUserErr(chatID, method, err)
		return
	}

	if *user.ActiveProgramID == programID {
		msg := tgbotapi.NewMessage(chatID, "Нельзя удалить текущую программу 😓")
		_, _ = tghelpers.SendMessage(s.bot, msg, method)
		return
	}

	program, err := s.programsRepo.Get(programID)
	if err != nil {
		return
	}

	err = s.programsRepo.Delete(&program)
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, "✅ Успешно удалено!")
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
	s.programManagement(chatID)
}

func (s *serviceImpl) changeProgram(chatID, programID int64) {
	method := "changeProgram"
	fmt.Sprintf("%s: programID: %d", method, programID)

	user, err := s.usersRepo.GetByChatID(chatID)
	if err != nil {
		s.handleGetUserErr(chatID, method, err)
		return
	}

	*user.ActiveProgramID = programID
	err = s.usersRepo.Save(user)
	if err != nil {
		fmt.Printf("%s: %s\n", method, err.Error())
	}

	msg := tgbotapi.NewMessage(chatID, "✅ Успешно обновлено!")
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
	s.programManagement(chatID)
}

func (s *serviceImpl) formatPreset(preset string) string {
	exercises := utils.SplitPreset(preset)

	buffer := &bytes.Buffer{}
	for _, ex := range exercises {

		exerciseType, err := s.exerciseTypesRepo.Get(ex.ID)
		if err != nil {
			continue
		}
		buffer.WriteString(fmt.Sprintf("• *%s*\n", exerciseType.Name))
		buffer.WriteString(fmt.Sprintf("    • "))
		for i, set := range ex.Sets {
			if i > 0 {
				buffer.WriteString(", ")
			}
			if set.Minutes > 0 {
				buffer.WriteString(fmt.Sprintf("%d мин", set.Minutes))
			} else {
				buffer.WriteString(fmt.Sprintf("%d \\* %.0f кг", set.Reps, set.Weight))
			}
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (s *serviceImpl) addNewDayTypeExercise(chatID, dayTypeID int64) {
	method := "addNewDayTypeExercise"
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
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1],
			tgbotapi.NewInlineKeyboardButtonData(group.Name, fmt.Sprintf("exercise_select_for_program_day_%d_%s", dayTypeID, group.Code)),
		)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
}

func (s *serviceImpl) programManagement(chatID int64) {
	method := "programManagement"

	user, err := s.usersRepo.GetByChatID(chatID)
	if err != nil {
		s.handleGetUserErr(chatID, method, err)
		return
	}

	programs, err := s.programsRepo.FindAll(user.ID)
	if err != nil {
		return
	}

	addNewProgram := tgbotapi.NewInlineKeyboardButtonData("➕ Добавить новую", "program_create")

	if len(programs) == 0 {
		msg := tgbotapi.NewMessage(chatID, "🥲 У вас нет тренировочных программ, создайте первую!")
		msg.ParseMode = constants.MarkdownParseMode
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(addNewProgram))
		_, err = s.bot.Send(msg)
		handleErr(method, err)
		return
	}

	text := &bytes.Buffer{}
	text.WriteString("*Ваши программы:*\n\n")

	var rows [][]tgbotapi.InlineKeyboardButton
	for i, program := range programs {
		if i%2 == 0 {
			rows = append(rows, []tgbotapi.InlineKeyboardButton{})
		}

		if program.ID == *user.ActiveProgramID {
			text.WriteString(fmt.Sprintf("• 🟢 *%s* \n  📅 %s\n\n", program.Name, program.CreatedAt.Format("02.01.2006 15:04")))
		} else {
			text.WriteString(fmt.Sprintf("• *%s* \n 📅 %s\n\n", program.Name, program.CreatedAt.Format("02.01.2006 15:04")))
		}

		rows[len(rows)-1] = append(rows[len(rows)-1],
			tgbotapi.NewInlineKeyboardButtonData(program.Name, fmt.Sprintf("program_edit_%d", program.ID)))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(addNewProgram))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = constants.MarkdownParseMode
	msg.ReplyMarkup = keyboard
	_, err = s.bot.Send(msg)
	handleErr(method, err)
}
