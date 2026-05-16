package workouts

import (
	"errors"

	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/calculator"
)

type CalculateWorkoutCaloriesUC struct {
	workoutRepo  workouts.Repo
	exerciseRepo exercises.Repo
	userRepo     users.Repo
}

func NewCalculateWorkoutCaloriesUC(
	workoutsRepo workouts.Repo,
	exerciseRepo exercises.Repo,
	userRepo users.Repo,
) *CalculateWorkoutCaloriesUC {
	return &CalculateWorkoutCaloriesUC{
		workoutRepo:  workoutsRepo,
		exerciseRepo: exerciseRepo,
		userRepo:     userRepo,
	}
}

func (uc *CalculateWorkoutCaloriesUC) Name() string {
	return "Расчет калорий"
}

var (
	ErrWeightRequired = errors.New("weight required")
	ErrGenderRequired = errors.New("gender required")
)

func (uc *CalculateWorkoutCaloriesUC) Execute(workoutID int64) (*dto.CaloriesCalc, error) {
	workout, err := uc.workoutRepo.Get(workoutID)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.GetByID(workout.UserID)
	if err != nil {
		return nil, err
	}

	if user.WeightKg == nil || *user.WeightKg == 0 {
		return nil, ErrWeightRequired
	}

	if user.Gender == nil || *user.Gender == "" {
		return nil, ErrGenderRequired
	}

	exerciseObjs, err := uc.exerciseRepo.FindAllByWorkoutID(workoutID)
	if err != nil {
		return nil, err
	}

	calc := calculator.NewCalculator(*user.WeightKg, *user.Gender, user.BirthDate)
	calories, durationMin := calc.CalculateWorkout(exerciseObjs)

	workout.EstimatedCalories = &calories
	workout.EstimatedDurationMinutes = &durationMin
	workout.UserWeightKg = user.WeightKg

	return &dto.CaloriesCalc{
		Calories:    calories,
		DurationMin: durationMin,
		UserWeight:  user.WeightKg,
	}, nil
}
