package models

import "time"

type WorkoutProgramProgressionRule struct {
	ID               int64 `gorm:"primaryKey;autoIncrement"`
	WorkoutProgramID int64
	WorkoutDayTypeID *int64
	ExerciseTypeID   *int64
	Rule             string
	Reason           string
	Source           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (*WorkoutProgramProgressionRule) TableName() string {
	return "workout_program_progression_rules"
}
