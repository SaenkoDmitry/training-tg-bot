package service

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *serviceImpl) HandleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := message.Text

	fmt.Println("HandleMessage:", text)

	user, _ := s.usersRepo.GetUser(chatID, message.From.UserName)

	switch {
	case text == "/start" || text == "/menu" || text == "🔙 В меню":
		s.sendMainMenu(chatID)

	case text == "/start_workout" || text == "▶️ Начать тренировку":
		s.showWorkoutTypeMenu(chatID)

	case text == "/stats" || text == "📊 Статистика":
		s.showStatsMenu(chatID, user.ID)

	case text == "📋 Мои тренировки" || text == "/workouts":
		s.showMyWorkouts(chatID)

		// default:
		// 	handleState(chatID, user.ID, text)
	}
}

func (s *serviceImpl) sendMainMenu(chatID int64) {
	text := "🏋️‍♂️ *Добро пожаловать в Бот для тренировок!* \n\n Выберите действие:"

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("▶️ Начать тренировку"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📋 Мои тренировки"),
			tgbotapi.NewKeyboardButton("📊 Статистика"),
		),
	)
	keyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) showWorkoutTypeMenu(chatID int64) {
	text := "Выберите тип тренировки:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🦵 Ноги", "create_workout_legs"),
			tgbotapi.NewInlineKeyboardButtonData("🏋️‍♂️ Спина", "create_workout_back"),
		),
		// tgbotapi.NewInlineKeyboardRow(
		// tgbotapi.NewInlineKeyboardButtonData("💪 Руки", "create_workout_arms"),
		// 	tgbotapi.NewInlineKeyboardButtonData("🫀 Грудь", "create_workout_chest"),
		// ),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("🌀 Плечи", "create_workout_shoulders"),
		// 	tgbotapi.NewInlineKeyboardButtonData("⚡️ Кардио", "create_workout_cardio"),
		// ),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) showMyWorkouts(chatID int64) {
	user := s.usersRepo.GetUserByChatID(chatID)

	workouts, _ := s.workoutsRepo.Find(user.ID)

	if len(workouts) == 0 {
		msg := tgbotapi.NewMessage(chatID, "📭 У вас пока нет созданных тренировок.\n\nСоздайте первую тренировку!")
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙 В меню", "back_to_menu"),
			),
		)
		msg.ReplyMarkup = keyboard
		s.bot.Send(msg)
		return
	}

	text := "📋 *Ваши тренировки:*\n\n"
	for i, workout := range workouts {
		status := "🟡 Активна"
		if workout.Completed {
			status = "✅ Завершена"
		}
		date := workout.StartedAt.Format("02.01.2006")
		text += fmt.Sprintf("%d. *%s* - %s\n   📅 %s\n\n",
			i+1, workout.Name, status, date)
	}

	text += "Выберите тренировку для просмотра:"

	var rows [][]tgbotapi.InlineKeyboardButton
	for i, workout := range workouts {
		if i%2 == 0 {
			rows = append(rows, []tgbotapi.InlineKeyboardButton{})
		}
		rowIndex := len(rows) - 1
		buttonText := fmt.Sprintf("%s %d", workout.Name, i+1)
		rows[rowIndex] = append(rows[rowIndex],
			tgbotapi.NewInlineKeyboardButtonData(buttonText,
				fmt.Sprintf("view_workout_%d", workout.ID)))
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔙 В меню", "back_to_menu"),
	})

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) showStatsMenu(chatID int64, userID int64) {
	text := "📊 *Статистика тренировок*\n\n Выберите период:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📅 За неделю", "stats_week"),
			tgbotapi.NewInlineKeyboardButtonData("🗓️ За месяц", "stats_month"),
			tgbotapi.NewInlineKeyboardButtonData("📈 Общая", "stats_all"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}
