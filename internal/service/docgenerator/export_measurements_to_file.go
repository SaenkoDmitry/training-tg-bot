package docgenerator

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/docgenerator/helpers"
	"github.com/xuri/excelize/v2"
)

const (
	MeasurementSheet = messages.Measurements
)

func (s *serviceImpl) ExportMeasurementsToFile(measurements []*dto.Measurement) (*excelize.File, error) {
	f := excelize.NewFile()

	headerStyle := helpers.HeaderStyle(f, constants.SkyBlueColor)
	s.writeMeasurementChartSheet(f, measurements, headerStyle)

	_ = f.SetRowStyle(MeasurementSheet, 1, 1, headerStyle)
	helpers.AutoFitColumns(f, MeasurementSheet, 1, 12)
	_ = f.DeleteSheet(DefaultSheet)

	f.SetActiveSheet(0)
	return f, nil
}
