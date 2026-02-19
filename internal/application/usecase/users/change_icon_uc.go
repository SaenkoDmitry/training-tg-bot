package users

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type ChangeIconUseCase struct {
	usersRepo users.Repo
}

func NewChangeIconUseCase(usersRepo users.Repo) *ChangeIconUseCase {
	return &ChangeIconUseCase{
		usersRepo: usersRepo,
	}
}

func (uc *ChangeIconUseCase) Name() string {
	return "Сменить иконку"
}

func (uc *ChangeIconUseCase) Execute(userID int64, name string) error {
	return uc.usersRepo.ChangeIcon(userID, name)
}
