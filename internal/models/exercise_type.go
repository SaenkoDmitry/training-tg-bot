package models

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"strings"
)

type ExerciseType struct {
	ID                    int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                  string `json:"name"`
	Url                   string `json:"url"`
	ExerciseGroupTypeCode string `json:"exercise_group_type_code"`
	RestInSeconds         int    `json:"rest_in_seconds"`
	Accent                string `json:"accent"`
	Units                 string `json:"units"`
	Description           string `json:"description"`
}

func (*ExerciseType) TableName() string {
	return "exercise_types"
}

func (t *ExerciseType) ContainsReps() bool {
	if t == nil {
		return false
	}
	return strings.Contains(t.Units, constants.RepsUnit)
}

func (t *ExerciseType) ContainsWeight() bool {
	if t == nil {
		return false
	}
	return strings.Contains(t.Units, constants.WeightUnit)
}

func (t *ExerciseType) ContainsMinutes() bool {
	if t == nil {
		return false
	}
	return strings.Contains(t.Units, constants.MinutesUnit)
}

func (t *ExerciseType) ContainsMeters() bool {
	if t == nil {
		return false
	}
	return strings.Contains(t.Units, constants.MetersUnit)
}
