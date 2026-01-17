package service

import "strings"

func (s *serviceImpl) statsCases(data string, chatID int64) {
	switch {
	case strings.HasPrefix(data, "stats_"):
		period := strings.TrimPrefix(data, "stats_")
		s.showStatistics(chatID, period)
	}
}
