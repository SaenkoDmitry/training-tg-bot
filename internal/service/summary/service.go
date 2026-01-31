package summary

import (
	"math"
	"sort"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
)

type Service interface {
	BuildTotal(workouts []models.WorkoutDay, groupCodesMap map[string]string) map[string]*ExerciseSummary
	BuildByDate(workouts []models.WorkoutDay) map[string]*DateSummary
	BuildExerciseProgressByDates(workouts []models.WorkoutDay) []*ExerciseProgressByDates
	BuildByWeekAndExType(workouts []models.WorkoutDay, groupCodesMap map[string]string) map[utils.DateRange]map[string]*WeekSummary
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
		date := w.StartedAt.Add(3 * time.Hour).Format("2006-01-02")

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
		date := w.StartedAt.Add(3 * time.Hour).Format("2006-01-02")

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

func (s *serviceImpl) BuildExerciseProgressByDates(workouts []models.WorkoutDay) []*ExerciseProgressByDates {
	exerciseWithProgressesMap := make(map[string]map[string]*Progress) // exName -> date -> progress

	for _, w := range workouts {
		if !w.Completed {
			continue
		}

		currDate := w.StartedAt.Add(3 * time.Hour).Format("02.01.2006")
		thisWeekRange := utils.GetThisWeekRange(w.StartedAt)

		var key string
		for _, e := range w.Exercises {
			if e.CompletedSets() == 0 {
				continue
			}

			if e.ExerciseType.ExerciseGroupTypeCode == "cardio" {
				key = thisWeekRange.Format()
			} else {
				key = currDate
			}

			sumWeight := float32(0)
			countOfReps := 0

			if _, ok := exerciseWithProgressesMap[e.ExerciseType.Name]; !ok {
				exerciseWithProgressesMap[e.ExerciseType.Name] = make(map[string]*Progress)
			}
			if _, ok := exerciseWithProgressesMap[e.ExerciseType.Name][key]; !ok {
				exerciseWithProgressesMap[e.ExerciseType.Name][key] = &Progress{
					Units:     e.ExerciseType.Units,
					GroupCode: e.ExerciseType.ExerciseGroupTypeCode,
				}
			}

			tempProgress := exerciseWithProgressesMap[e.ExerciseType.Name][key]

			for _, set := range e.Sets {
				if !set.Completed {
					continue
				}

				countOfReps += set.GetRealReps()
				sumWeight += set.GetRealWeight() * float32(set.GetRealReps())
				if tempProgress.MaxWeight < set.GetRealWeight() ||
					tempProgress.MaxWeight == set.GetRealWeight() && tempProgress.MaxReps < set.GetRealReps() {
					tempProgress.MaxWeight = set.GetRealWeight()
					tempProgress.MaxReps = set.GetRealReps()
				}

				tempProgress.SumMinutes += set.GetRealMinutes()
				if tempProgress.MinMinutes == 0 {
					tempProgress.MinMinutes = set.GetRealMinutes()
					tempProgress.MaxMinutes = set.GetRealMinutes()
				} else {
					tempProgress.MinMinutes = min(tempProgress.MinMinutes, set.GetRealMinutes())
					tempProgress.MaxMinutes = max(tempProgress.MaxMinutes, set.GetRealMinutes())
				}

				tempProgress.SumMeters += set.GetRealMeters()
				if tempProgress.MinMeters == 0 {
					tempProgress.MinMeters = set.GetRealMeters()
					tempProgress.MaxMeters = set.GetRealMeters()
				} else {
					tempProgress.MinMeters = min(tempProgress.MinMeters, set.GetRealMeters())
					tempProgress.MaxMeters = max(tempProgress.MaxMeters, set.GetRealMeters())
				}
			}
			tempProgress.AvgWeight = sumWeight / float32(countOfReps)
		}
	}

	result := make([]*ExerciseProgressByDates, 0, len(exerciseWithProgressesMap))
	for exName, progressByDate := range exerciseWithProgressesMap {
		units := ""
		groupCode := ""
		progressByDates := make([]*DateWithProgress, 0)
		for date, progress := range progressByDate {
			units = progress.Units
			groupCode = progress.GroupCode
			progressByDates = append(progressByDates, &DateWithProgress{
				Date:     date,
				Progress: progress,
			})
		}
		sort.Slice(progressByDates, func(i, j int) bool {
			return progressByDates[i].Date < progressByDates[j].Date
		})
		result = append(result, &ExerciseProgressByDates{
			ExerciseName:          exName,
			DateWithProgress:      progressByDates,
			ExerciseUnitType:      units,
			ExerciseGroupTypeCode: groupCode,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].ExerciseGroupTypeCode == "cardio" {
			return true
		}
		if result[j].ExerciseGroupTypeCode == "cardio" {
			return true
		}
		return result[i].ExerciseGroupTypeCode < result[j].ExerciseGroupTypeCode
	})

	return result
}

func (s *serviceImpl) BuildByWeekAndExType(workouts []models.WorkoutDay, groupCodesMap map[string]string) map[utils.DateRange]map[string]*WeekSummary {
	result := make(map[utils.DateRange]map[string]*WeekSummary)
	for _, w := range workouts {
		thisWeek := utils.GetThisWeekRange(w.StartedAt)
		if _, ok := result[thisWeek]; !ok {
			result[thisWeek] = map[string]*WeekSummary{}
		}
		for _, e := range w.Exercises {
			groupName := groupCodesMap[e.ExerciseType.ExerciseGroupTypeCode]
			if _, ok := result[thisWeek][groupName]; !ok {
				result[thisWeek][groupName] = &WeekSummary{}
			}
			for _, set := range e.Sets {
				result[thisWeek][groupName].SumWeight += float32(set.GetRealReps()) * set.GetRealWeight()
				result[thisWeek][groupName].SumMinutes += set.GetRealMinutes()
				result[thisWeek][groupName].SumMeters += set.GetRealMeters()
			}
		}
	}
	return result
}
