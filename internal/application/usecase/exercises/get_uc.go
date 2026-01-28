package exercises

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
)

type GetUseCase struct {
	exercisesRepo     exercises.Repo
	exerciseTypesRepo exercisetypes.Repo
}

func NewGetUseCase(
	exercisesRepo exercises.Repo,
	exerciseTypesRepo exercisetypes.Repo,
) *GetUseCase {
	return &GetUseCase{
		exercisesRepo:     exercisesRepo,
		exerciseTypesRepo: exerciseTypesRepo,
	}
}

func (uc *GetUseCase) Name() string {
	return "Показать данные об упражнении"
}

func (uc *GetUseCase) Execute(exerciseTypeID int64) (*dto.GetExercise, error) {
	exType, err := uc.exerciseTypesRepo.Get(exerciseTypeID)
	if err != nil {
		return nil, err
	}

	return &dto.GetExercise{ExerciseType: exType}, nil
}
