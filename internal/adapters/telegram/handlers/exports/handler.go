package exports

import (
	exportusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/exports"
	userusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/users"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type Handler struct {
	presenter                   *Presenter
	getUserUC                   *userusecases.GetUseCase
	exportWorkoutsToExcelUC     *exportusecases.ExportWorkoutsToExcelUseCase
	exportMeasurementsToExcelUC *exportusecases.ExportMeasurementsToExcelUseCase
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	getUserUC *userusecases.GetUseCase,
	exportWorkoutsToExcelUC *exportusecases.ExportWorkoutsToExcelUseCase,
	exportMeasurementsToExcelUC *exportusecases.ExportMeasurementsToExcelUseCase,
) *Handler {
	return &Handler{
		presenter:                   NewPresenter(bot),
		getUserUC:                   getUserUC,
		exportWorkoutsToExcelUC:     exportWorkoutsToExcelUC,
		exportMeasurementsToExcelUC: exportMeasurementsToExcelUC,
	}
}

func (h *Handler) RouteCallback(chatID int64, data string) {
	switch {
	case strings.HasPrefix(data, "export_workouts_to_excel"):
		h.exportWorkoutsToExcel(chatID)
	case strings.HasPrefix(data, "export_measurements_to_excel"):
		h.exportMeasurementsToExcel(chatID)
	}
}

func (h *Handler) exportWorkoutsToExcel(chatID int64) {
	user, err := h.getUserUC.Execute(chatID)
	if err != nil {
		return
	}
	buffer, err := h.exportWorkoutsToExcelUC.Execute(user.ID)
	if err != nil {
		h.presenter.CannotDoAction(chatID, h.exportWorkoutsToExcelUC.Name())
		return
	}
	h.presenter.WriteDoc(chatID, buffer, "workouts.xlsx")
}

func (h *Handler) exportMeasurementsToExcel(chatID int64) {
	buffer, err := h.exportMeasurementsToExcelUC.Execute(chatID)
	if err != nil {
		h.presenter.CannotDoAction(chatID, h.exportMeasurementsToExcelUC.Name())
		return
	}
	h.presenter.WriteDoc(chatID, buffer, "measurements.xlsx")
}
