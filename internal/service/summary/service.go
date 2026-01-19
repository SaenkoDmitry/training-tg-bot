package summary

import (
	"math"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
)

type Service interface {
	BuildTotal(workouts []models.WorkoutDay, groupCodesMap map[string]string) map[string]*ExerciseSummary
	BuildByDate(workouts []models.WorkoutDay) map[string]*DateSummary
	BuildExerciseProgress(workouts []models.WorkoutDay, exerciseName string) map[string]*Progress
}

type serviceImpl struct {
}

func NewService() Service {
	return &serviceImpl{}
}

func (s *serviceImpl) BuildTotal(workouts []models.WorkoutDay, groupCodesMap map[string]string) map[string]*ExerciseSummary {
	summary := make(map[string]*ExerciseSummary)

	for _, w := range workouts {
		if !w.Completed {
			continue
		}
		date := w.StartedAt.Format("2006-01-02")

		for _, e := range w.Exercises {
			if e.CompletedSets() == 0 {
				continue
			}

			sum, ok := summary[e.ExerciseType.Name]
			if !ok {
				sum = &ExerciseSummary{
					Workouts: make(map[string]struct{}),
				}
				summary[e.ExerciseType.Name] = sum
			}

			sum.ExerciseType = groupCodesMap[e.ExerciseType.ExerciseGroupTypeCode]
			sum.Workouts[date] = struct{}{}

			for _, set := range e.Sets {
				if !set.Completed {
					continue
				}
				sum.Sets++
				sum.TotalWeight += float64(set.Weight) * float64(set.Reps)
				sum.TotalReps += set.Reps
				sum.TotalMinutes += set.Minutes

				if set.Weight > sum.MaxWeight {
					sum.MaxWeight = set.Weight
				}
			}
		}
	}

	for _, sum := range summary {
		if sum.TotalReps > 0 {
			sum.AvgWeight = math.Round(sum.TotalWeight / float64(sum.TotalReps))
		}
	}

	return summary
}

func (s *serviceImpl) BuildByDate(workouts []models.WorkoutDay) map[string]*DateSummary {
	result := make(map[string]*DateSummary)

	for _, w := range workouts {
		date := w.StartedAt.Format("2006-01-02")

		d, ok := result[date]
		if !ok {
			d = &DateSummary{
				Workouts:  1,
				Exercises: make(map[string]struct{}),
			}
			result[date] = d
		}

		for _, e := range w.Exercises {
			d.Exercises[e.ExerciseType.Name] = struct{}{}

			for _, sum := range e.Sets {
				d.Sets++
				d.TotalVolume += sum.Weight * float32(sum.Reps)

				if sum.Weight > d.MaxWeight {
					d.MaxWeight = sum.Weight
				}
			}
		}
	}

	return result
}

func (s *serviceImpl) BuildExerciseProgress(
	workouts []models.WorkoutDay,
	exerciseName string,
) map[string]*Progress {

	progress := make(map[string]*Progress)

	for _, w := range workouts {
		if !w.Completed {
			continue
		}

		date := w.StartedAt.Format("2006-01-02")

		for _, e := range w.Exercises {
			if e.ExerciseType.Name != exerciseName {
				continue
			}
			if e.CompletedSets() == 0 {
				continue
			}

			sumWeight := float32(0)
			countOfReps := 0

			progress[date] = &Progress{}

			for _, set := range e.Sets {
				if !set.Completed {
					continue
				}

				countOfReps += set.Reps
				sumWeight += set.Weight * float32(set.Reps)
				if progress[date].MaxWeight < set.Weight ||
					progress[date].MaxWeight == set.Weight && progress[date].MaxReps < set.Reps {
					progress[date].MaxWeight = set.Weight
					progress[date].MaxReps = set.Reps
				}

				progress[date].SumMinutes += set.Minutes
				if progress[date].MinMinutes == 0 {
					progress[date].MinMinutes = set.Minutes
					progress[date].MaxMinutes = set.Minutes
				} else {
					progress[date].MinMinutes = min(progress[date].MinMinutes, set.Minutes)
					progress[date].MaxMinutes = max(progress[date].MaxMinutes, set.Minutes)
				}

				progress[date].SumMeters += set.Meters
				if progress[date].MinMeters == 0 {
					progress[date].MinMeters = set.Meters
					progress[date].MaxMeters = set.Meters
				} else {
					progress[date].MinMeters = min(progress[date].MinMeters, set.Meters)
					progress[date].MaxMeters = max(progress[date].MaxMeters, set.Meters)
				}
			}
			progress[date].AvgWeight = sumWeight / float32(countOfReps)
		}
	}

	return progress
}
