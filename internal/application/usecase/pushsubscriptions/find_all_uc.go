package pushsubscriptions

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/pushsubscriptions"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type FindAllUseCase struct {
	pushSubscriptionsRepo pushsubscriptions.Repo
	usersRepo             users.Repo
}

func NewFindAllUseCase(
	pushSubscriptionsRepo pushsubscriptions.Repo,
	usersRepo users.Repo,
) *FindAllUseCase {
	return &FindAllUseCase{
		pushSubscriptionsRepo: pushSubscriptionsRepo,
		usersRepo:             usersRepo,
	}
}

func (uc *FindAllUseCase) Name() string {
	return "Получить все пуши пользователя"
}

func (uc *FindAllUseCase) Execute(chatID int64) ([]*models.PushSubscription, error) {
	user, err := uc.usersRepo.GetByChatID(chatID)
	if err != nil {
		return nil, err
	}

	subs, err := uc.pushSubscriptionsRepo.FindAll(user.ID)
	if err != nil {
		return nil, err
	}

	return subs, nil
}
