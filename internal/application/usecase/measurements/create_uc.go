package measurements

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/measurements"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
)

type CreateUseCase struct {
	measurementsRepo measurements.Repo
}

func NewCreateUseCase(measurementsRepo measurements.Repo) *CreateUseCase {
	return &CreateUseCase{
		measurementsRepo: measurementsRepo,
	}
}

func (uc *CreateUseCase) Name() string {
	return "Добавить новый измерение тела"
}

func (uc *CreateUseCase) Execute(measurement *models.Measurement) (*dto.Measurement, error) {
	err := uc.measurementsRepo.Save(measurement)
	if err != nil {
		return nil, err
	}

	return &dto.Measurement{
		ID:        measurement.ID,
		UserID:    measurement.UserID,
		CreatedAt: utils.FormatDate(measurement.CreatedAt),
		Shoulders: utils.FormatCentimeters(measurement.Shoulders),
		Chest:     utils.FormatCentimeters(measurement.Chest),
		HandLeft:  utils.FormatCentimeters(measurement.HandLeft),
		HandRight: utils.FormatCentimeters(measurement.HandRight),
		Waist:     utils.FormatCentimeters(measurement.Waist),
		Buttocks:  utils.FormatCentimeters(measurement.Buttocks),
		HipLeft:   utils.FormatCentimeters(measurement.HipLeft),
		HipRight:  utils.FormatCentimeters(measurement.HipRight),
		CalfLeft:  utils.FormatCentimeters(measurement.CalfLeft),
		CalfRight: utils.FormatCentimeters(measurement.CalfRight),
		Weight:    utils.FormatKilograms(measurement.Weight),
	}, nil
}
