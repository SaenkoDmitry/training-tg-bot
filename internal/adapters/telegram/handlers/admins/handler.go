package admins

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/common"
	userusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/users"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type Handler struct {
	presenter       *Presenter
	commonPresenter *common.Presenter

	findUsersUC *userusecases.FindUseCase
}

func NewHandler(bot *tgbotapi.BotAPI, findUsersUC *userusecases.FindUseCase) *Handler {
	return &Handler{
		presenter:       NewPresenter(bot),
		commonPresenter: common.NewPresenter(bot),
		findUsersUC:     findUsersUC,
	}
}

func (h *Handler) RouteCallback(chatID int64, data string) {
	switch {
	case strings.HasPrefix(data, "/admin/users"):
		h.users(chatID)
	}
}

func (h *Handler) users(chatID int64) {
	users, err := h.findUsersUC.Execute(chatID, 0)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.findUsersUC.Name())
		return
	}
	h.presenter.ShowTopUsers(chatID, users)
}
