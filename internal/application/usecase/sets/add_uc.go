package sets

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisegrouptypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sets"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/docgenerator"
	summarysvc "github.com/SaenkoDmitry/training-tg-bot/internal/service/summary"
)

type AddOneMoreUseCase struct {
	usersRepo              users.Repo
	exerciseGroupTypesRepo exercisegrouptypes.Repo
	workoutsRepo           workouts.Repo
	exercisesRepo          exercises.Repo
	summaryService         summarysvc.Service
	docGeneratorService    docgenerator.Service
	setsRepo               sets.Repo
}

func NewAddOneMoreUseCase(
	setsRepo sets.Repo,
	exercisesRepo exercises.Repo,
) *AddOneMoreUseCase {
	return &AddOneMoreUseCase{
		setsRepo:      setsRepo,
		exercisesRepo: exercisesRepo,
	}
}

func (uc *AddOneMoreUseCase) Name() string {
	return "Добавить подход"
}

func (uc *AddOneMoreUseCase) Execute(exerciseID int64) (*dto.AddOneMoreSet, error) {
	ex, err := uc.exercisesRepo.Get(exerciseID)
	if err != nil {
		return nil, err
	}
	nextSet := &models.Set{
		ExerciseID: ex.ID,
	}

	if len(ex.Sets) > 0 {
		lastSet := ex.Sets[len(ex.Sets)-1]
		if ex.GetExerciseType().ContainsReps() {
			nextSet.Reps = lastSet.Reps
			if lastSet.FactReps > 0 {
				nextSet.Reps = lastSet.FactReps
			}
		}
		if ex.GetExerciseType().ContainsWeight() {
			nextSet.Weight = lastSet.Weight
			if lastSet.FactWeight > 0 {
				nextSet.Weight = lastSet.FactWeight
			}
		}
		if ex.GetExerciseType().ContainsMinutes() {
			nextSet.Minutes = lastSet.Minutes
			if lastSet.FactMinutes > 0 {
				nextSet.Minutes = lastSet.FactMinutes
			}
		}
		if ex.GetExerciseType().ContainsMeters() {
			nextSet.Meters = lastSet.Meters
			if lastSet.FactMeters > 0 {
				nextSet.Meters = lastSet.FactMeters
			}
		}
		nextSet.Index = lastSet.Index + 1
	} else {
		switch {
		case ex.GetExerciseType().ContainsReps():
			nextSet.Reps = 10
		case ex.GetExerciseType().ContainsWeight():
			nextSet.Weight = 10
		case ex.GetExerciseType().ContainsMinutes():
			nextSet.Minutes = 10
		case ex.GetExerciseType().ContainsMeters():
			nextSet.Meters = 10
		}
	}

	err = uc.setsRepo.Save(nextSet)
	if err != nil {
		fmt.Println("cannot create set:", err.Error())
		return nil, err
	}

	return &dto.AddOneMoreSet{
		WorkoutID: ex.WorkoutDayID,
	}, nil
}
