package exercisetypes

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	Get(exerciseTypeID int64) (models.ExerciseType, error)
	GetAll() ([]models.ExerciseType, error)
	GetAllByGroup(code string) ([]models.ExerciseType, error)
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (u *repoImpl) Get(exerciseTypeID int64) (exerciseType models.ExerciseType, err error) {
	u.db.First(&exerciseType, exerciseTypeID)
	return exerciseType, nil
}

func (u *repoImpl) GetAll() (exerciseTypes []models.ExerciseType, err error) {
	u.db.Order("id ASC").Find(&exerciseTypes)
	return exerciseTypes, nil
}

func (u *repoImpl) GetAllByGroup(code string) (exerciseTypes []models.ExerciseType, err error) {
	u.db.Where("exercise_group_type_code = ?", code).Order("id ASC").Find(&exerciseTypes)
	return exerciseTypes, nil
}
