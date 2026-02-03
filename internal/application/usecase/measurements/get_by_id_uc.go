package measurements

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/measurements"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
)

type GetByIDUseCase struct {
	measurementsRepo measurements.Repo
}

func NewGetByIDUseCase(
	measurementsRepo measurements.Repo,
) *GetByIDUseCase {
	return &GetByIDUseCase{
		measurementsRepo: measurementsRepo,
	}
}

func (uc *GetByIDUseCase) Name() string {
	return "Показать измерение"
}

func (uc *GetByIDUseCase) Execute(measurementID int64) (*dto.Measurement, error) {
	measurementObj, err := uc.measurementsRepo.Get(measurementID)
	if err != nil {
		return nil, err
	}

	return &dto.Measurement{
		ID:        measurementObj.ID,
		CreatedAt: utils.FormatDate(measurementObj.CreatedAt),
		Shoulders: utils.FormatCentimeters(measurementObj.Shoulders),
		Chest:     utils.FormatCentimeters(measurementObj.Chest),
		HandLeft:  utils.FormatCentimeters(measurementObj.HandLeft),
		HandRight: utils.FormatCentimeters(measurementObj.HandRight),
		Waist:     utils.FormatCentimeters(measurementObj.Waist),
		Buttocks:  utils.FormatCentimeters(measurementObj.Buttocks),
		HipLeft:   utils.FormatCentimeters(measurementObj.HipLeft),
		HipRight:  utils.FormatCentimeters(measurementObj.HipRight),
		CalfLeft:  utils.FormatCentimeters(measurementObj.CalfLeft),
		CalfRight: utils.FormatCentimeters(measurementObj.CalfRight),
		Weight:    utils.FormatKilograms(measurementObj.Weight),
	}, nil
}
