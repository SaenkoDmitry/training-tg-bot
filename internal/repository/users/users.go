package users

import (
	"errors"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	GetTop10() ([]models.User, error)
	FindTopN(offset, limit int) ([]UserWithCount, error)
	Save(user *models.User) error
	Create(chatID int64, from *tgbotapi.User) (*models.User, error)
	GetByID(ID int64) (*models.User, error)
	GetByChatID(chatID int64) (*models.User, error)
	ChangeIcon(userID int64, name string) error
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (u *repoImpl) Save(user *models.User) error {
	u.db.Save(user)
	return nil
}

func (u *repoImpl) Create(chatID int64, from *tgbotapi.User) (*models.User, error) {
	user := models.User{
		ChatID:       chatID,
		Username:     from.UserName,
		FirstName:    from.FirstName,
		LastName:     from.LastName,
		LanguageCode: from.LanguageCode,
		CreatedAt:    time.Now(),
	}
	err := u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&user).Error
	})
	return &user, err
}

var (
	NotFoundUserErr = errors.New("not found user")
)

func (u *repoImpl) GetByChatID(chatID int64) (*models.User, error) {
	var user models.User

	result := u.db.Where("chat_id = ?", chatID).First(&user)

	if result.Error != nil {
		return nil, NotFoundUserErr
	}
	return &user, nil
}

func (u *repoImpl) GetByID(ID int64) (*models.User, error) {
	var user models.User

	result := u.db.Where("id = ?", ID).First(&user)

	if result.Error != nil {
		return nil, NotFoundUserErr
	}
	return &user, nil
}

func (u *repoImpl) GetTop10() ([]models.User, error) {
	var users []models.User
	tx := u.db.Order("created_at DESC").Limit(10).Find(&users)
	if tx.Error != nil {
		return []models.User{}, tx.Error
	}
	return users, nil
}

type UserWithCount struct {
	models.User
	WorkoutCount int64 `gorm:"column:workout_count"`
}

func (u *repoImpl) FindTopN(offset, limit int) ([]UserWithCount, error) {
	var users []UserWithCount
	tx := u.db.
		Table("users").
		Select("users.*, COUNT(w.id) as workout_count").
		Joins(`
        LEFT JOIN workout_days w 
        ON w.user_id = users.id 
        AND w.completed = true
    `).
		Group("users.id").
		Order("workout_count DESC").
		Limit(limit).
		Offset(offset).
		Scan(&users)
	if tx.Error != nil {
		return []UserWithCount{}, tx.Error
	}
	return users, nil
}

func (u *repoImpl) ChangeIcon(userID int64, name string) error {
	return u.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("icon", name).Error
}
