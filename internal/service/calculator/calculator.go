package calculator

import (
	"math"
	"strings"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
)

type CalorieCalculator struct {
	UserWeightKg float64
	Gender       string // "male" | "female"
	BirthDate    *time.Time
}

// MET-значения по группам (Compendium of Physical Activities)
var metValues = map[string]float64{
	"legs":     8.0, // приседания, жим ногами
	"chest":    6.0, // жим лёжа
	"back":     6.5, // тяги
	"biceps":   4.0, // изолированно
	"triceps":  4.0,
	"deltas":   5.0, // плечи
	"press":    3.5, // пресс
	"buttocks": 6.0, // ягодицы
	"cardio":   0.0, // считаем отдельно по скорости
}

// среднее время выполнения подхода по группам мышц
var repDurationSec = map[string]int{
	"legs":     5, // приседания медленные
	"back":     4, // тяги
	"chest":    3, // жим лёжа
	"biceps":   2, // сгибания быстрые
	"triceps":  2,
	"deltas":   3,
	"press":    2, // пресс
	"buttocks": 4,
}

// Кардио MET по типу (приблизительно)
var cardioMETs = map[string]float64{
	"run_treadmill": 9.8,
	"run_outdoor":   9.0,
	"walk":          3.5,
	"bike":          7.5,
	"swim":          8.0,
	"elliptical":    5.0,
}

func NewCalculator(weightKg float64, gender string, birthDate *time.Time) *CalorieCalculator {
	return &CalorieCalculator{
		UserWeightKg: weightKg,
		Gender:       gender,
		BirthDate:    birthDate,
	}
}

func (c *CalorieCalculator) CalculateWorkout(exercises []models.Exercise) (float64, int) {
	var totalCalories float64
	var totalDurationSec int

	for _, ex := range exercises {
		if ex.CompletedSets() == 0 {
			continue
		}
		calories, durationSec := c.calculateExercise(ex)
		totalCalories += calories
		totalDurationSec += durationSec
	}

	return math.Round(totalCalories*100) / 100, totalDurationSec / 60
}

func (c *CalorieCalculator) calculateExercise(ex models.Exercise) (float64, int) {
	if ex.ExerciseType.ExerciseGroupTypeCode == "cardio" {
		return c.calculateCardio(ex, ex.ExerciseType)
	}
	return c.calculateStrength(ex, ex.ExerciseType)
}

// Силовое: MET × вес × время
func (c *CalorieCalculator) calculateStrength(
	ex models.Exercise,
	exType *models.ExerciseType,
) (calories float64, durationSec int) {
	met := metValues[exType.ExerciseGroupTypeCode]
	if met == 0 {
		met = 5.0 // fallback
	}

	// Длительность: работа + отдых
	totalReps := 0
	for _, s := range ex.Sets {

		// скипаем расчет для сета, если еще не выполнено
		if !s.Completed {
			continue
		}

		reps := s.FactReps
		if reps == 0 {
			reps = s.Reps // если факт не указан, берём план
		}
		totalReps += reps
	}

	workSec := totalReps * repDurationSec[ex.ExerciseType.ExerciseGroupTypeCode]

	restSec := 0
	for _, s := range ex.Sets {
		if !s.Completed {
			continue
		}
		restSec += exType.RestInSeconds
	}
	if restSec < 0 {
		restSec = 0
	}
	durationSec = workSec + restSec

	// Корректировка на вес штанги: +10% за каждые 50% от веса тела
	avgWeight := c.avgWeight(ex.Sets)
	weightMultiplier := 1.0 + (avgWeight/c.UserWeightKg)*0.2

	hours := float64(durationSec) / 3600.0
	calories = met * c.UserWeightKg * hours * weightMultiplier * c.genderFactor() * c.ageFactor()

	return
}

// Кардио: MET × вес × время (из minutes/meters)
func (c *CalorieCalculator) calculateCardio(ex models.Exercise, exType *models.ExerciseType) (float64, int) {
	var minutes int
	var meters int

	for _, s := range ex.Sets {
		if !s.Completed {
			continue
		}
		if s.FactMinutes > 0 {
			minutes += s.FactMinutes
		} else if s.Minutes > 0 {
			minutes += s.Minutes
		}
		if s.FactMeters > 0 {
			meters += s.FactMeters
		} else if s.Meters > 0 {
			meters += s.Meters
		}
	}

	// Если есть метры, но нет минут — эстимируем скорость
	if minutes == 0 && meters > 0 {
		minutes = c.estimateMinutesFromMeters(exType.Name, meters)
	}

	met := c.cardioMET(exType.Name, meters, minutes)
	hours := float64(minutes) / 60.0

	calories := met * c.UserWeightKg * hours * c.genderFactor() * c.ageFactor()

	return calories, minutes * 60
}

func (c *CalorieCalculator) avgWeight(sets []models.Set) float64 {
	var total float64
	var count int
	for _, s := range sets {
		if !s.Completed {
			continue
		}
		w := s.FactWeight
		if w == 0 {
			w = s.Weight
		}
		if w > 0 {
			total += float64(w)
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

func (c *CalorieCalculator) genderFactor() float64 {
	if c.Gender == "female" {
		return 0.9
	}
	return 1.0
}

func (c *CalorieCalculator) ageFactor() float64 {
	if c.BirthDate == nil {
		return 1.0
	}
	age := int(time.Since(*c.BirthDate).Hours() / 24 / 365)
	if age <= 30 {
		return 1.0
	}
	// -2% за каждые 10 лет после 30
	return 1.0 - float64(age-30)*0.002
}

func (c *CalorieCalculator) cardioMET(name string, meters, minutes int) float64 {
	// Пытаемся определить по названию
	switch {
	case contains(name, "бег", "run"):
		if minutes > 0 && meters > 0 {
			speed := float64(meters) / float64(minutes) // м/мин
			if speed > 200 {
				return 11.0 // быстрый бег
			}
		}
		return 9.8
	case contains(name, "ходьба", "walk"):
		return 3.5
	case contains(name, "велосипед", "bike"):
		return 7.5
	case contains(name, "бассейн", "плавание", "swim"):
		return 8.0
	case contains(name, "эллипс", "ellip"):
		return 5.0
	default:
		return 6.0
	}
}

func (c *CalorieCalculator) estimateMinutesFromMeters(name string, meters int) int {
	switch {
	case contains(name, "бег", "run"):
		return meters / 150 // ~9 км/ч
	case contains(name, "ходьба", "walk"):
		return meters / 80 // ~5 км/ч
	case contains(name, "плавание", "swim"):
		return meters / 50 // ~3 км/ч
	default:
		return meters / 100
	}
}

func contains(s string, subs ...string) bool {
	s = strings.ToLower(s)
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
