package workouts

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	Create(workoutDay *models.WorkoutDay) error
	Delete(workout *models.WorkoutDay) error
	Save(workout *models.WorkoutDay) error
	Get(workoutID int64) (models.WorkoutDay, error)
	Count(userID int64) (count int64, err error)
	FindAll(userID int64) ([]models.WorkoutDay, error)
	Find(userID int64, offset, limit int) ([]models.WorkoutDay, error)
	FindPreviousByType(userID int64, workoutType string) (models.WorkoutDay, error)
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (u *repoImpl) Create(workoutDay *models.WorkoutDay) error {
	u.db.Create(&workoutDay)
	return nil
}

func (u *repoImpl) Delete(workout *models.WorkoutDay) error {
	u.db.Delete(workout)
	return nil
}

func (u *repoImpl) Save(workout *models.WorkoutDay) error {
	u.db.Save(workout)
	return nil
}

func (u *repoImpl) Get(workoutID int64) (workoutDay models.WorkoutDay, err error) {
	u.db.
		Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB { return db.Order("sets.index ASC") }).
		Preload("Exercises", func(db *gorm.DB) *gorm.DB { return db.Order("exercises.index ASC") }).
		First(&workoutDay, workoutID)

	return workoutDay, nil
}

func (u *repoImpl) Count(userID int64) (count int64, err error) {
	result := u.db.Model(&models.WorkoutDay{}).Where("user_id = ?", userID).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (u *repoImpl) FindAll(userID int64) (workouts []models.WorkoutDay, err error) {
	u.db.
		Where("user_id = ?", userID).
		Order("started_at DESC").
		Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB { return db.Order("sets.index ASC") }).
		Preload("Exercises", func(db *gorm.DB) *gorm.DB { return db.Order("exercises.index ASC") }).
		Find(&workouts)

	return workouts, nil
}

func (u *repoImpl) Find(userID int64, offset, limit int) (workouts []models.WorkoutDay, err error) {
	u.db.
		Where("user_id = ?", userID).
		Order("started_at DESC").
		Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB { return db.Order("sets.index ASC") }).
		Preload("Exercises", func(db *gorm.DB) *gorm.DB { return db.Order("exercises.index ASC") }).
		Offset(offset).
		Limit(limit).
		Find(&workouts)

	return workouts, nil
}

func (u *repoImpl) FindPreviousByType(userID int64, workoutType string) (models.WorkoutDay, error) {
	var workout models.WorkoutDay
	u.db.Where("user_id = ? AND name = ? AND completed = ?", userID, workoutType, true).
		Order("started_at DESC").
		Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB { return db.Order("sets.index ASC") }).
		Preload("Exercises", func(db *gorm.DB) *gorm.DB { return db.Order("exercises.index ASC") }).
		First(&workout)
	return workout, nil
}
