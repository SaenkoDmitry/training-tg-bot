package models

import "time"

type WorkoutSession struct {
	ID                   int64 `gorm:"primaryKey"`
	WorkoutDayID         int64
	CurrentExerciseIndex int
	StartedAt            time.Time
	IsActive             bool
}

func (*WorkoutSession) TableName() string {
	return "training.workout_sessions"
}
