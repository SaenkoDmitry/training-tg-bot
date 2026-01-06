package utils

import (
	"fmt"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
)

func GetWorkoutNameByID(ID string) string {
	switch ID {
	case constants.LegsAndShouldersWorkoutID:
		return constants.LegsAndShouldersWorkoutName
	case constants.BackAndBicepsWorkoutID:
		return constants.BackAndBicepsWorkoutName
	case constants.ChestAndTricepsID:
		return constants.ChestAndTricepsName
	case constants.CardioID:
		return constants.CardioName
	}
	return ""
}

func FormatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d ч %d мин", hours, minutes)
	} else if minutes > 0 {
		return fmt.Sprintf("%d мин", minutes)
	}
	return fmt.Sprintf("%d сек", seconds)
}

func BetweenTimes(startedAt time.Time, endedAt *time.Time) string {
	duration := endedAt.Sub(startedAt)
	return FormatDuration(duration)
}
