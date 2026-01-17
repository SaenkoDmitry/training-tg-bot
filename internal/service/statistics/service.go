package statistics

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/daytypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisegrouptypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
	"strings"
	"time"
)

type Service interface {
	ShowWorkoutStatistics(workoutID int64) string
	ShowPeriodStatistics(userID int64, period string) string
}

type serviceImpl struct {
	usersRepo              users.Repo
	dayTypesRepo           daytypes.Repo
	workoutsRepo           workouts.Repo
	exerciseTypesRepo      exercisetypes.Repo
	exerciseGroupTypesRepo exercisegrouptypes.Repo
}

func NewService(
	usersRepo users.Repo,
	dayTypesRepo daytypes.Repo,
	workoutsRepo workouts.Repo,
	exerciseTypesRepo exercisetypes.Repo,
	exerciseGroupTypesRepo exercisegrouptypes.Repo,
) Service {
	return &serviceImpl{
		usersRepo:              usersRepo,
		workoutsRepo:           workoutsRepo,
		dayTypesRepo:           dayTypesRepo,
		exerciseTypesRepo:      exerciseTypesRepo,
		exerciseGroupTypesRepo: exerciseGroupTypesRepo,
	}
}

func (s *serviceImpl) ShowWorkoutStatistics(workoutID int64) string {
	workoutDay, err := s.workoutsRepo.Get(workoutID)
	if err != nil {
		return ""
	}

	dayType, err := s.dayTypesRepo.Get(workoutDay.WorkoutDayTypeID)
	if err != nil {
		return ""
	}

	totalWeight := 0.0
	completedExercises := 0
	totalTime := 0

	var text strings.Builder
	text.WriteString(messages.Statistics + fmt.Sprintf(": %s\n\n", dayType.Name))

	if workoutDay.EndedAt != nil {
		text.WriteString(messages.WorkoutTime + fmt.Sprintf(": %s\n", utils.BetweenTimes(workoutDay.StartedAt, workoutDay.EndedAt)))
	}

	text.WriteString(fmt.Sprintf(messages.WorkoutDate+": %s\n\n", workoutDay.StartedAt.Add(3*time.Hour).Format("02.01.2006 15:04")))

	for _, exercise := range workoutDay.Exercises {
		if exercise.CompletedSets() == 0 {
			continue
		}

		exerciseObj, getErr := s.exerciseTypesRepo.Get(exercise.ExerciseTypeID)
		if getErr != nil {
			continue
		}

		completedExercises++
		exerciseTime := 0
		exerciseWeight := 0.0
		maxWeight := 0.0

		for _, set := range exercise.Sets {
			if !set.Completed {
				continue
			}
			exerciseWeight += float64(set.GetRealWeight()) * float64(set.GetRealReps())
			exerciseTime += set.GetRealMinutes()
			maxWeight = max(maxWeight, float64(set.GetRealWeight()))
		}
		totalWeight += exerciseWeight
		totalTime += exerciseTime

		lastSet := exercise.Sets[len(exercise.Sets)-1]
		text.WriteString(fmt.Sprintf("• <b>%s:</b> \n", exerciseObj.Name))
		if lastSet.GetRealReps() > 0 {
			text.WriteString(fmt.Sprintf("  • Выполнено: %d из %d подходов\n", exercise.CompletedSets(), len(exercise.Sets)))
			text.WriteString(fmt.Sprintf("  • Рабочий вес: %d * %.0f кг \n", lastSet.GetRealReps(), lastSet.GetRealWeight()))
			text.WriteString(fmt.Sprintf("  • Общий вес: %.0f кг \n\n", exerciseWeight))
		} else if lastSet.GetRealMinutes() > 0 {
			text.WriteString(fmt.Sprintf("  • Общее время: %d минут \n\n", exerciseTime))
		}
	}

	text.WriteString(messages.Summary + "\n")
	text.WriteString(fmt.Sprintf("• Упражнений: %d/%d\n", completedExercises, len(workoutDay.Exercises)))
	if totalWeight > 0 {
		text.WriteString(fmt.Sprintf("• Общий тоннаж: %.0f кг\n", totalWeight))
	}
	if totalTime > 0 {
		text.WriteString(fmt.Sprintf("• Общее время: %d минут\n", totalTime))
	}
	return text.String()
}

const (
	week  = "week"
	month = "month"
)

func (s *serviceImpl) ShowPeriodStatistics(userID int64, period string) string {
	workoutObjs, _ := s.workoutsRepo.FindAll(userID)

	completedWorkouts := 0
	sumTime := time.Duration(0)
	cardioTime := 0
	for _, w := range workoutObjs {
		if !w.Completed {
			continue
		}
		switch period {
		case week:
			if time.Since(w.StartedAt).Hours() > 7*24 {
				continue
			}
		case month:
			if time.Since(w.StartedAt).Hours() > 30*24 {
				continue
			}
		default:
		}

		completedWorkouts++
		sumTime += w.EndedAt.Sub(*&w.StartedAt)
		for _, e := range w.Exercises {
			if len(e.Sets) == 0 {
				continue
			}
			for _, s := range e.Sets {
				if !s.Completed {
					continue
				}
				if s.GetRealMinutes() > 0 {
					cardioTime += s.GetRealMinutes()
				}
			}
		}
	}

	avgTime := time.Duration(0)
	if completedWorkouts != 0 {
		avgTime = sumTime / time.Duration(completedWorkouts)
	}

	var statsText strings.Builder
	switch period {
	case week:
		statsText.WriteString(messages.StatisticsWeek)
	case month:
		statsText.WriteString(messages.StatisticsMonth)
	default:
		statsText.WriteString(messages.StatisticsAll)
	}
	statsText.WriteString("\n\n")

	statsText.WriteString(messages.EndsWorkouts + fmt.Sprintf(": %d\n", completedWorkouts))
	statsText.WriteString(messages.AvgWorkoutTime + fmt.Sprintf(": %s\n", utils.FormatDuration(avgTime)))
	statsText.WriteString(messages.OverallWorkoutTime + fmt.Sprintf(": %d мин\n", cardioTime))

	return statsText.String()
}
