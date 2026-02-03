package measurements

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/measurements"
)

type DeleteByIDUseCase struct {
	measurementsRepo measurements.Repo
}

func NewDeleteByIDUseCase(
	measurementsRepo measurements.Repo,
) *DeleteByIDUseCase {
	return &DeleteByIDUseCase{
		measurementsRepo: measurementsRepo,
	}
}

func (uc *DeleteByIDUseCase) Name() string {
	return "Удалить измерение"
}

func (uc *DeleteByIDUseCase) Execute(measurementID int64) error {
	return uc.measurementsRepo.DeleteByID(measurementID)
}
