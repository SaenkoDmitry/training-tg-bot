package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
)

type User struct {
	ID        int64 `gorm:"primaryKey"`
	Username  string
	ChatID    int64
	CreatedAt time.Time
}

type WorkoutDay struct {
	ID        int64 `gorm:"primaryKey"`
	UserID    int64
	Name      string
	Exercises []Exercise `gorm:"foreignKey:WorkoutDayID"`
	StartedAt time.Time
	EndedAt   *time.Time
	Completed bool
}

func (w *WorkoutDay) Status() string {
	if !w.Completed {
		return fmt.Sprintf("⏳ Активна")
	}
	if w.EndedAt != nil {
		return fmt.Sprintf("✅ Завершена в %s", w.EndedAt.Add(3*time.Hour).Format("15:04"))
	}

	return fmt.Sprintf("✅ Завершена")
}

func (w *WorkoutDay) String() string {
	var text strings.Builder

	text.WriteString(fmt.Sprintf("*Тренировка:* %s \n", utils.GetWorkoutNameByID(w.Name)))
	text.WriteString(fmt.Sprintf("*Статус:* %s\n", w.Status()))
	text.WriteString(fmt.Sprintf("*Дата:* %s\n\n", w.StartedAt.Add(3*time.Hour).Format("02.01.2006")))
	text.WriteString("*Упражнения:*\n")

	for i, exercise := range w.Exercises {
		text.WriteString(fmt.Sprintf("%s %d. %s: \n", exercise.Status(), i+1, exercise.Name))
		lastSet := exercise.Sets[len(exercise.Sets)-1]
		text.WriteString(fmt.Sprintf("Рабочий вес: %d \\* %.0f кг \n\n", lastSet.Reps, lastSet.Weight))
	}

	return text.String()
}

type Exercise struct {
	ID            int64 `gorm:"primaryKey"`
	WorkoutDayID  int64
	Name          string
	Sets          []Set `gorm:"foreignKey:ExerciseID"`
	Hint          string
	RestInSeconds int
}

func (e *Exercise) Status() string {
	completedExerciseSets := e.CompletedSets()
	allSets := len(e.Sets)

	status := "🔴"
	if int(completedExerciseSets) >= allSets {
		status = "🟢"
	} else if completedExerciseSets > 0 {
		status = "🟡"
	}
	return status
}

func (e *Exercise) CompletedSets() int {
	completedSets := 0
	for _, set := range e.Sets {
		if set.Completed {
			completedSets++
		}
	}
	return completedSets
}

func (e *Exercise) NextSet() Set {
	for _, set := range e.Sets {
		if !set.Completed {
			return set
		}
	}
	return Set{}
}

type Set struct {
	ID          int64 `gorm:"primaryKey"`
	ExerciseID  int64
	Reps        int
	Weight      float32
	Completed   bool
	CompletedAt *time.Time
	Index       int
}

type WorkoutSession struct {
	ID                   int64 `gorm:"primaryKey"`
	WorkoutDayID         int64
	CurrentExerciseIndex int
	StartedAt            time.Time
	IsActive             bool
}
