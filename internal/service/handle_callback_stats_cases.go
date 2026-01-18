package service

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/tghelpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (s *serviceImpl) statsCases(data string, chatID, userID int64) {
	switch {
	case strings.HasPrefix(data, "stats_"):
		period := strings.TrimPrefix(data, "stats_")
		s.showStatistics(chatID, userID, period)
	}
}

func (s *serviceImpl) showStatistics(chatID, userID int64, period string) {
	method := "showStatistics"
	text := s.statisticsService.ShowPeriodStatistics(userID, period)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
}
