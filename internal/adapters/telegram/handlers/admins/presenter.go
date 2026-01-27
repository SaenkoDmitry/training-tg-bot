package admins

import (
	"bytes"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Presenter struct {
	bot *tgbotapi.BotAPI
}

func (p Presenter) ShowTopUsers(chatID int64, users []users.UserWithCount) {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	var text bytes.Buffer
	text.WriteString(fmt.Sprintf("<b>%s:</b>\n\n", messages.Users))
	for i, u := range users {
		if i%2 == 0 {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow())
		}
		text.WriteString(fmt.Sprintf("%d. %s (%d тренировок)\n\n", i+1, u.FullName(), u.WorkoutCount))
		rows[len(rows)-1] = append(rows[len(rows)-1],
			tgbotapi.NewInlineKeyboardButtonData(u.Username, fmt.Sprintf("workout_show_by_user_id_%d", u.ID)),
		)
	}
	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	p.bot.Send(msg)
}

func NewPresenter(bot *tgbotapi.BotAPI) *Presenter {
	return &Presenter{bot: bot}
}
