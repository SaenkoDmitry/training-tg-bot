package docgenerator

import (
	"sort"
	"strconv"

	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/summary"
	"github.com/xuri/excelize/v2"
)

func (s *serviceImpl) writeByDateSummarySheet(f *excelize.File, summary map[string]*summary.DateSummary) {
	sheet := ByDateSummarySheet
	_, _ = f.NewSheet(sheet)

	headers := []string{
		messages.WorkoutDate,
		"Тренировок",
		"Упражнений",
		"Сетов",
		"Общий объём (кг)",
		"Макс вес (кг)",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheet, cell, h)
	}

	row := 2
	dates := make([]string, 0, len(summary))
	for d := range summary {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	for _, date := range dates {
		sum := summary[date]

		_ = f.SetCellValue(sheet, "A"+strconv.Itoa(row), date)
		_ = f.SetCellValue(sheet, "B"+strconv.Itoa(row), sum.Workouts)
		_ = f.SetCellValue(sheet, "C"+strconv.Itoa(row), len(sum.Exercises))
		_ = f.SetCellValue(sheet, "D"+strconv.Itoa(row), sum.Sets)
		_ = f.SetCellValue(sheet, "E"+strconv.Itoa(row), sum.TotalVolume)
		_ = f.SetCellValue(sheet, "F"+strconv.Itoa(row), sum.MaxWeight)
		row++
	}
}
