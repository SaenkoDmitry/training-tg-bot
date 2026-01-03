package exercises

import (
	"sort"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	Get(exerciseID int64) (models.Exercise, error)
	FindAllByWorkoutID(workoutDayID int64) ([]models.Exercise, error)
	Delete(workoutID int64) error
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (u *repoImpl) Get(exerciseID int64) (models.Exercise, error) {
	var exercise models.Exercise
	u.db.Preload("Sets").First(&exercise, exerciseID)
	sort.Slice(exercise.Sets, func(i, j int) bool {
		return exercise.Sets[i].Index < exercise.Sets[j].Index
	})
	return exercise, nil
}

func (u *repoImpl) Delete(workoutID int64) error {
	u.db.Where("workout_day_id = ?", workoutID).Delete(&models.Exercise{})
	return nil
}

func (u *repoImpl) FindAllByWorkoutID(workoutDayID int64) ([]models.Exercise, error) {
	var exercises []models.Exercise
	u.db.Where("workout_day_id = ?", workoutDayID).Find(&exercises)
	for _, exercise := range exercises {
		sort.Slice(exercise.Sets, func(i, j int) bool {
			return exercise.Sets[i].Index < exercise.Sets[j].Index
		})
	}
	return exercises, nil
}
