package docgenerator

import (
	"fmt"

	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/xuri/excelize/v2"
)

func (s *serviceImpl) writeWorkoutsSheet(f *excelize.File, workouts []models.WorkoutDay, groupCodesMap map[string]string) {
	sheet := WorkoutSheet
	_, _ = f.NewSheet(sheet)

	headers := []string{
		messages.WorkoutDate,
		"Тренировка",
		"Упражнение",
		"Тип",
		"Номер сета",
		messages.Weight,
		messages.Reps,
		messages.Minutes,
		messages.Meters,
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheet, cell, h)
	}

	row := 2

	for _, w := range workouts {
		for _, e := range w.Exercises {
			for i, set := range e.Sets {
				_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", row), w.StartedAt.Format("2006-01-02"))
				_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", row), w.WorkoutDayType.Name)
				_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", row), e.ExerciseType.Name)
				_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", row), groupCodesMap[e.ExerciseType.ExerciseGroupTypeCode])
				_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", row), i+1)
				_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", row), set.Weight)
				_ = f.SetCellValue(sheet, fmt.Sprintf("G%d", row), set.Reps)
				_ = f.SetCellValue(sheet, fmt.Sprintf("H%d", row), set.Minutes)
				_ = f.SetCellValue(sheet, fmt.Sprintf("I%d", row), set.Meters)
				row++
			}
		}
	}
}
