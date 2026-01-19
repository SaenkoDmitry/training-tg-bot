package docgenerator

import (
	"strconv"

	"github.com/SaenkoDmitry/training-tg-bot/internal/service/summary"
	"github.com/xuri/excelize/v2"
)

func (s *serviceImpl) writeTotalSummarySheet(f *excelize.File, summary map[string]*summary.ExerciseSummary) {
	sheet := TotalSummarySheet
	_, _ = f.NewSheet(sheet)

	headers := []string{
		"Упражнение",
		"Тип",
		//"Акцент",
		"Тренировок",
		"Сетов",
		"Макс вес (кг)",
		"Средний вес (кг)",
		"Общий объём",
		"Общее время",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheet, cell, h)
	}

	row := 2
	for name, sum := range summary {
		_ = f.SetCellValue(sheet, "A"+strconv.Itoa(row), name)
		_ = f.SetCellValue(sheet, "B"+strconv.Itoa(row), sum.ExerciseType)
		_ = f.SetCellValue(sheet, "C"+strconv.Itoa(row), len(sum.Workouts))
		_ = f.SetCellValue(sheet, "D"+strconv.Itoa(row), sum.Sets)
		_ = f.SetCellValue(sheet, "E"+strconv.Itoa(row), sum.MaxWeight)
		_ = f.SetCellValue(sheet, "F"+strconv.Itoa(row), sum.AvgWeight)
		_ = f.SetCellValue(sheet, "G"+strconv.Itoa(row), sum.TotalWeight)
		_ = f.SetCellValue(sheet, "H"+strconv.Itoa(row), sum.TotalMinutes)

		row++
	}
}
