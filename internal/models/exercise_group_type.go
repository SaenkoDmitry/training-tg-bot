package models

type ExerciseGroupType struct {
	Code string `gorm:"primaryKey" json:"code"`
	Name string `json:"name"`
}

func (*ExerciseGroupType) TableName() string {
	return "exercise_group_types"
}
