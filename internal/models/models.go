package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
	"gorm.io/gorm"
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
	Exercises []Exercise `gorm:"foreignKey:WorkoutDayID;constraint:OnDelete:CASCADE"`
	StartedAt time.Time
	EndedAt   *time.Time
	Completed bool
}

func (wd *WorkoutDay) BeforeDelete(tx *gorm.DB) (err error) {
	// Удаляем все упражнения дня (что вызовет BeforeDelete для каждого Exercise)
	return tx.Where("workout_day_id = ?", wd.ID).Delete(&Exercise{}).Error
}

func (w *WorkoutDay) Status() string {
	if !w.Completed {
		return fmt.Sprintf("🟡 Активна")
	}
	if w.EndedAt != nil {
		return fmt.Sprintf("✅ Завершена в %s", w.EndedAt.Add(3*time.Hour).Format("15:04"))
	}

	return fmt.Sprintf("✅ Завершена")
}

func (w *WorkoutDay) String() string {
	var text strings.Builder

	text.WriteString(fmt.Sprintf("<b>Тип:</b> %s \n", utils.GetWorkoutNameByID(w.Name)))
	text.WriteString(fmt.Sprintf("<b>Статус:</b> %s\n", w.Status()))
	text.WriteString(fmt.Sprintf("<b>Дата:</b> 📅 %s\n\n", w.StartedAt.Add(3*time.Hour).Format("02.01.2006")))
	text.WriteString("<b>Упражнения:</b>\n")

	for i, exercise := range w.Exercises {
		exerciseObj, ok := constants.AllExercises[exercise.Name]
		if !ok {
			continue
		}
		text.WriteString(fmt.Sprintf("<b>%s %d. %s</b>\n", exercise.Status(), i+1, exerciseObj.GetName()))

		for _, set := range exercise.Sets {
			text.WriteString(set.String(w.Completed))
		}
		text.WriteString("\n")
	}

	return text.String()
}

type Exercise struct {
	ID            int64 `gorm:"primaryKey"`
	WorkoutDayID  int64
	Name          string
	Sets          []Set `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE"`
	RestInSeconds int
	Index         int
}

func (e *Exercise) BeforeDelete(tx *gorm.DB) (err error) {
	// Удаляем все сеты упражнения
	return tx.Where("exercise_id = ?", e.ID).Delete(&Set{}).Error
}

func (e *Exercise) Status() string {
	completedExerciseSets := e.CompletedSets()
	allSets := len(e.Sets)

	status := "🔴"
	if completedExerciseSets >= allSets {
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
	FactReps    int
	Weight      float32
	FactWeight  float32
	Minutes     int
	FactMinutes int
	Completed   bool
	CompletedAt *time.Time
	Index       int
}

func (s *Set) String(done bool) string {
	var text strings.Builder

	if s.Minutes > 0 {
		text.WriteString(fmt.Sprintf("• %s минут: ", s.FormatMinutes()))
	} else {
		text.WriteString(fmt.Sprintf("• %s повторов по %s кг: ", s.FormatReps(), s.FormatWeight()))
	}

	if s.Completed {
		text.WriteString(fmt.Sprintf("✅, %s", s.CompletedAt.Add(3*time.Hour).Format("15:04:05")))
	} else {
		if done {
			text.WriteString("💔")
		} else {
			text.WriteString("🚀")
		}
	}
	text.WriteString("\n")
	return text.String()
}

func (s *Set) FormatReps() string {
	if s.FactReps != 0 {
		return fmt.Sprintf("<strike>%d</strike> <b>%d</b>", s.Reps, s.FactReps)
	}
	return fmt.Sprintf("%d", s.Reps)
}

func (s *Set) FormatWeight() string {
	if s.FactWeight != float32(0) {
		return fmt.Sprintf("<strike>%.0f</strike> <b>%.0f</b>", s.Weight, s.FactWeight)
	}
	return fmt.Sprintf("%.0f", s.Weight)
}

func (s *Set) FormatMinutes() string {
	if s.FactMinutes != 0 {
		return fmt.Sprintf("<strike>%d</strike> <b>%d</b>", s.Minutes, s.FactMinutes)
	}
	return fmt.Sprintf("%d", s.Minutes)
}

func (s *Set) GetRealReps() int {
	if s == nil {
		return 0
	}
	if s.FactReps > 0 {
		return s.FactReps
	}
	return s.Reps
}

func (s *Set) GetRealWeight() float32 {
	if s == nil {
		return 0
	}
	if s.FactWeight > 0 {
		return s.FactWeight
	}
	return s.Weight
}

func (s *Set) GetRealMinutes() int {
	if s == nil {
		return 0
	}
	if s.FactMinutes > 0 {
		return s.FactMinutes
	}
	return s.Minutes
}

type WorkoutSession struct {
	ID                   int64 `gorm:"primaryKey"`
	WorkoutDayID         int64
	CurrentExerciseIndex int
	StartedAt            time.Time
	IsActive             bool
}
