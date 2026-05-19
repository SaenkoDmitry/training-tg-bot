package models

import (
	"time"

	"github.com/lib/pq"
)

type WorkoutProgram struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64
	Name      string
	CreatedAt time.Time
	DayTypes  []WorkoutDayType `gorm:"foreignKey:WorkoutProgramID;constraint:OnDelete:CASCADE"`

	Summary         *string        `db:"summary"`
	Warnings        pq.StringArray `gorm:"type:text[];column:warnings"`
	ValidationNotes pq.StringArray `gorm:"type:text[];column:validation_notes"`
}

func (*WorkoutProgram) TableName() string {
	return "workout_programs"
}
