package models

import "time"

type RestTimer struct {
	ID        int64 `gorm:"primaryKey"`
	UserID    int64
	WorkoutID int64
	EndsAt    time.Time
	Canceled  bool
	Sent      bool
	CreatedAt time.Time
}

func (*RestTimer) TableName() string {
	return "rest_timers"
}
