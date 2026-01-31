package docgenerator

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"

	"github.com/xuri/excelize/v2"

	"github.com/SaenkoDmitry/training-tg-bot/internal/service/summary"
)

func (s *serviceImpl) writeAllProgressCharts(
	f *excelize.File,
	exerciseProgressByDates []*summary.ExerciseProgressByDates,
	redHeaderStyle int,
	greenHeaderStyle int,
	blueHeaderStyle int,
) {

	sheet := ByExerciseSummarySheet
	_, _ = f.NewSheet(sheet)

	row := 1
	for i, e := range exerciseProgressByDates {
		style := blueHeaderStyle
		switch {
		case i%3 == 1:
			style = redHeaderStyle
		case i%3 == 2:
			style = greenHeaderStyle
		}
		row = s.writeProgressChart(f, sheet, e.ExerciseName, e.DateWithProgress, e.ExerciseUnitType, row, style)
	}
}

func (s *serviceImpl) writeProgressChart(
	f *excelize.File,
	sheet string,
	exercise string,
	dateWithProgresses []*summary.DateWithProgress,
	units string,
	firstRow int,
	headerStyle int,
) int {
	switch {
	case strings.Contains(units, constants.RepsUnit):
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", firstRow), messages.WorkoutDate)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", firstRow), "Макс вес (кг)")
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", firstRow), "Макс кол-во повторов")
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", firstRow), "Средний вес (кг)")
	case strings.Contains(units, constants.MinutesUnit):
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", firstRow), messages.WorkoutDate)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", firstRow), "Макс время (минут)")
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", firstRow), "Мин время (минут)")
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", firstRow), "Всего (минут)")
	case strings.Contains(units, constants.MetersUnit):
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", firstRow), messages.WorkoutDate)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", firstRow), "Макс дистанция (метры)")
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", firstRow), "Мин дистанция (метры)")
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", firstRow), "Всего (метров)")
	}
	_ = f.SetRowStyle(ByExerciseSummarySheet, firstRow, firstRow, headerStyle)

	lastRow := firstRow + len(dateWithProgresses)

	row := firstRow + 1
	for _, d := range dateWithProgresses {
		switch {
		case strings.Contains(units, constants.RepsUnit):
			_ = f.SetCellValue(sheet, "A"+strconv.Itoa(row), d.Date)
			_ = f.SetCellValue(sheet, "B"+strconv.Itoa(row), float64(d.Progress.MaxWeight))
			_ = f.SetCellValue(sheet, "C"+strconv.Itoa(row), float64(d.Progress.MaxReps))
			_ = f.SetCellValue(sheet, "D"+strconv.Itoa(row), math.Round(float64(d.Progress.AvgWeight)))
		case strings.Contains(units, constants.MinutesUnit):
			_ = f.SetCellValue(sheet, "A"+strconv.Itoa(row), d.Date)
			_ = f.SetCellValue(sheet, "B"+strconv.Itoa(row), float64(d.Progress.MaxMinutes))
			_ = f.SetCellValue(sheet, "C"+strconv.Itoa(row), float64(d.Progress.MinMinutes))
			_ = f.SetCellValue(sheet, "D"+strconv.Itoa(row), math.Round(float64(d.Progress.SumMinutes)))
		case strings.Contains(units, constants.MetersUnit):
			_ = f.SetCellValue(sheet, "A"+strconv.Itoa(row), d.Date)
			_ = f.SetCellValue(sheet, "B"+strconv.Itoa(row), float64(d.Progress.MaxMeters))
			_ = f.SetCellValue(sheet, "C"+strconv.Itoa(row), float64(d.Progress.MinMeters))
			_ = f.SetCellValue(sheet, "D"+strconv.Itoa(row), math.Round(float64(d.Progress.SumMeters)))
		}
		row++
	}

	chart := &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{
			{
				Name:       fmt.Sprintf("%s!$B$%d", sheet, firstRow),
				Categories: fmt.Sprintf("%s!$A$%d:$A$%d", sheet, firstRow+1, lastRow),
				Values:     fmt.Sprintf("%s!$B$%d:$B$%d", sheet, firstRow+1, lastRow),
				Marker: excelize.ChartMarker{
					Symbol: "circle",
					Size:   6,
				},
			},
			{
				Name:       fmt.Sprintf("%s!$C$%d", sheet, firstRow),
				Categories: fmt.Sprintf("%s!$A$%d:$A$%d", sheet, firstRow+1, lastRow),
				Values:     fmt.Sprintf("%s!$C$%d:$C$%d", sheet, firstRow+1, lastRow),
				Marker: excelize.ChartMarker{
					Symbol: "circle",
					Size:   6,
				},
			},
			{
				Name:       fmt.Sprintf("%s!$D$%d", sheet, firstRow),
				Categories: fmt.Sprintf("%s!$A$%d:$A$%d", sheet, firstRow+1, lastRow),
				Values:     fmt.Sprintf("%s!$D$%d:$D$%d", sheet, firstRow+1, lastRow),
				Marker: excelize.ChartMarker{
					Symbol: "circle",
					Size:   6,
				},
			},
		},
		Title: []excelize.RichTextRun{
			{Text: exercise},
		},
		Legend: excelize.ChartLegend{
			Position: "bottom",
		},
	}

	err := f.AddChart(sheet, fmt.Sprintf("F%d", firstRow+1), chart)
	if err != nil {
		fmt.Println("error while build chart:", err.Error())
	}

	return firstRow + betweenExerciseRows
}

const (
	betweenExerciseRows = 17
)
