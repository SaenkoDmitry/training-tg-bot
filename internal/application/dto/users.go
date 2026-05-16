package dto

import "github.com/SaenkoDmitry/training-tg-bot/internal/models"

type TelegramUser struct {
	ID           int64  `json:"id"` // chat_id
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	PhotoURL     string `json:"photo_url"`
	AuthDate     int64  `json:"auth_date"`
	Hash         string `json:"hash"`
}

type YandexProfile struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	DefaultEmail string `json:"default_email"`
	Birthday     string `json:"birthday"` // формат "YYYY-MM-DD"
	Sex          string `json:"sex"`      // "male" | "female" | null
}

type UserProfile struct {
	ID        int64    `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Icon      string   `json:"icon"`
	BirthDate *string  `json:"birth_date,omitempty"` // "YYYY-MM-DD"
	Gender    *string  `json:"gender,omitempty"`     // "male" | "female"
	WeightKg  *float64 `json:"weight_kg,omitempty"`
	HeightCm  *int     `json:"height_cm,omitempty"`
}

func MapToProfile(user *models.User) *UserProfile {
	profile := &UserProfile{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Icon:      user.Icon,
		Gender:    user.Gender,
		WeightKg:  user.WeightKg,
		HeightCm:  user.HeightCm,
	}
	if user.BirthDate != nil {
		profile.BirthDate = new(user.BirthDate.Format("2006-01-02"))
	}
	return profile
}

type UpdateProfileRequest struct {
	BirthDate *string  `json:"birth_date"` // "YYYY-MM-DD" или null
	Gender    *string  `json:"gender"`     // "male" | "female" или null
	WeightKg  *float64 `json:"weight_kg"`  // или null
	HeightCm  *int     `json:"height_cm"`  // или null
}
