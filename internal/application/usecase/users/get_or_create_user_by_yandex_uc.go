package users

import (
	"errors"

	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type GetOrCreateUserByYandexUseCase struct {
	usersRepo users.Repo
}

func NewGetOrCreateUserByYandexUseCase(usersRepo users.Repo) *GetOrCreateUserByYandexUseCase {
	return &GetOrCreateUserByYandexUseCase{
		usersRepo: usersRepo,
	}
}

func (uc *GetOrCreateUserByYandexUseCase) Name() string {
	return "Создать яндекс аккаунт в системе или найти существующий"
}

func (uc *GetOrCreateUserByYandexUseCase) Execute(profile *dto.YandexProfile) (*models.User, error) {
	user, err := uc.usersRepo.GetByYandexID(profile.ID)
	if err != nil && !errors.Is(err, users.NotFoundUserErr) {
		return nil, err
	}

	if err != nil && errors.Is(err, users.NotFoundUserErr) {
		user, err = uc.usersRepo.CreateYandex(profile)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
