package exercisegrouptypes

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	Get(code string) (models.ExerciseGroupType, error)
	GetAll() ([]models.ExerciseGroupType, error)
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (u *repoImpl) Get(code string) (group models.ExerciseGroupType, err error) {
	u.db.First(&group, code)
	return group, nil
}

func (u *repoImpl) GetAll() (groups []models.ExerciseGroupType, err error) {
	u.db.Order("code ASC").Find(&groups)
	return groups, nil
}
