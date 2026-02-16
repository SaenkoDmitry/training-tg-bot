package pushsubscriptions

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
)

type Repo interface {
	Create(userID int64, sub dto.PushSubscription) error
	Delete(subID int64) error
	FindAll(userID int64) ([]*models.PushSubscription, error)
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (r *repoImpl) Create(userID int64, sub dto.PushSubscription) error {
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "endpoint"}},
		DoNothing: true,
	}).Create(&models.PushSubscription{
		UserID:   userID,
		Endpoint: sub.Endpoint,
		P256dh:   sub.Keys.P256dh,
		Auth:     sub.Keys.Auth,
	}).Error
	return err
}
func (r *repoImpl) FindAll(userID int64) ([]*models.PushSubscription, error) {
	var subs []*models.PushSubscription
	err := r.db.Where("user_id = ?", userID).Find(&subs).Error
	return subs, err
}

func (r *repoImpl) Delete(subID int64) error {
	return r.db.Where("id = ?", subID).Delete(&models.PushSubscription{}).Error
}
