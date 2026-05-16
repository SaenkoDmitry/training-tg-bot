package users

import (
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type UpdateProfileUseCase struct {
	repo users.Repo
}

func NewUpdateProfileUC(repo users.Repo) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{repo: repo}
}

func (uc *UpdateProfileUseCase) Execute(userID int64, req dto.UpdateProfileRequest) error {
	user, err := uc.repo.GetByID(userID)
	if err != nil {
		return err
	}

	// Валидация и обновление
	if req.BirthDate != nil {
		if *req.BirthDate == "" {
			user.BirthDate = nil
		} else {
			parsed, err := time.Parse("2006-01-02", *req.BirthDate)
			if err != nil {
				return err // или domain.ErrInvalidBirthDate
			}
			user.BirthDate = &parsed
		}
	}

	if req.Gender != nil {
		if *req.Gender == "" {
			user.Gender = nil
		} else if *req.Gender == "male" || *req.Gender == "female" {
			user.Gender = req.Gender
		}
	}

	if req.WeightKg != nil {
		if *req.WeightKg <= 0 {
			user.WeightKg = nil
		} else {
			user.WeightKg = req.WeightKg
		}
	}

	if req.HeightCm != nil {
		if *req.HeightCm <= 0 {
			user.HeightCm = nil
		} else {
			user.HeightCm = req.HeightCm
		}
	}

	return uc.repo.Save(user)
}
