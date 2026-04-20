package share

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/share"
)

type CreateShareUC struct {
	shareRepo share.Repo
}

func NewCreateShareUC(shareRepo share.Repo) *CreateShareUC {
	return &CreateShareUC{shareRepo: shareRepo}
}

func (uc *CreateShareUC) Execute(workoutID int64) (*models.WorkoutShare, error) {
	// Генерируем случайный токен 32 байта = 64 hex символа
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(tokenBytes)

	shareModel := &models.WorkoutShare{
		WorkoutDayID: workoutID,
		Token:        token,
		CreatedAt:    time.Now(),
		// todo добавить expires_at через 30 дней ?
		// ExpiresAt: &[]time.Time{time.Now().AddDate(0, 0, 30)}[0],
	}

	if err := uc.shareRepo.Create(shareModel); err != nil {
		return nil, err
	}

	return shareModel, nil
}
