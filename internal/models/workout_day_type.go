package models

import "time"

type WorkoutDayType struct {
	ID               int64 `gorm:"primaryKey"`
	WorkoutProgramID int64
	Name             string
	Preset           string // id_1:reps_0*weight_0,reps_1*weight_1;id_2:reps_2_0*weight_2_0;... // 3:17*100,15*160,12*200,12*240,12*260;2:14*40,14*40,14*45,14*50
	CreatedAt        time.Time
}

func (*WorkoutDayType) TableName() string {
	return "training.workout_day_types"
}
