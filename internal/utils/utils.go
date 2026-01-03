package utils

import "github.com/SaenkoDmitry/training-tg-bot/internal/constants"

func GetWorkoutNameByID(ID string) string {
	switch ID {
	case constants.LegsWorkoutID:
		return constants.LegsWorkoutName
	case constants.BackAndArmsWorkoutID:
		return constants.BackAndArmsWorkoutName
	case constants.ChestAndShouldersID:
		return constants.ChestAndShouldersName
	}
	return ""
}
