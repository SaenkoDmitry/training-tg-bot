package models

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"strings"
)

type ExerciseType struct {
	ID                    int64 `gorm:"primaryKey;autoIncrement"`
	Name                  string
	Url                   string
	ExerciseGroupTypeCode string
	RestInSeconds         int
	Accent                string
	Units                 string
	Description           string
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
