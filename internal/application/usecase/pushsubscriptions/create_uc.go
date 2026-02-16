package pushsubscriptions

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/pushsubscriptions"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type CreateUseCase struct {
	pushSubscriptionsRepo pushsubscriptions.Repo
	usersRepo             users.Repo
}

func NewCreateUseCase(
	pushSubscriptionsRepo pushsubscriptions.Repo,
	usersRepo users.Repo,
) *CreateUseCase {
	return &CreateUseCase{
		pushSubscriptionsRepo: pushSubscriptionsRepo,
		usersRepo:             usersRepo,
	}
}

func (uc *CreateUseCase) Name() string {
	return "Создать подписку"
}

func (uc *CreateUseCase) Execute(chatID int64, sub dto.PushSubscription) error {
	user, err := uc.usersRepo.GetByChatID(chatID)
	if err != nil {
		return err
	}

	err = uc.pushSubscriptionsRepo.Create(user.ID, sub)
	if err != nil {
		return err
	}

	return nil
}
