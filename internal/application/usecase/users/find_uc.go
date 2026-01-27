package users

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type FindUseCase struct {
	usersRepo    users.Repo
	programsRepo programs.Repo
}

func NewFindUseCase(usersRepo users.Repo, programsRepo programs.Repo) *FindUseCase {
	return &FindUseCase{
		usersRepo:    usersRepo,
		programsRepo: programsRepo,
	}
}

func (uc *FindUseCase) Name() string {
	return "Показать N самых активных (по числу тренировок) пользователей"
}

const (
	pageLimit = 10
)

func (uc *FindUseCase) Execute(chatID int64, offset int) ([]users.UserWithCount, error) {
	user, err := uc.usersRepo.GetByChatID(chatID)
	if err != nil {
		return nil, err
	}
	if !user.IsAdmin() {
		return nil, errors.New("not allowed")
	}
	usersWithCounts, err := uc.usersRepo.FindTopN(offset, pageLimit)
	if err != nil {
		return nil, err
	}
	return usersWithCounts, nil
}
