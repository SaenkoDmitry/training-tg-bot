package models

type ExerciseGroupType struct {
	Code string `gorm:"primaryKey"`
	Name string
}

func (*ExerciseGroupType) TableName() string {
	return "training.exercise_group_types"
}
