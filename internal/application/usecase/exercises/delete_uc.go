package exercises

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sessions"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
)

type DeleteUseCase struct {
	exercisesRepo exercises.Repo
	workoutsRepo  workouts.Repo
	sessionsRepo  sessions.Repo
}

func NewDeleteUseCase(exercisesRepo exercises.Repo, workoutsRepo workouts.Repo, sessionsRepo sessions.Repo) *DeleteUseCase {
	return &DeleteUseCase{
		exercisesRepo: exercisesRepo,
		workoutsRepo:  workoutsRepo,
		sessionsRepo:  sessionsRepo,
	}
}

func (uc *DeleteUseCase) Name() string {
	return "Удалить упражнение"
}

func (uc *DeleteUseCase) Execute(exerciseID int64) (int64, error) {
	exercise, err := uc.exercisesRepo.Get(exerciseID)
	if err != nil {
		return 0, err
	}
	workoutID := exercise.WorkoutDayID
	deleteExIndex := exercise.Index

	err = uc.exercisesRepo.Delete(exerciseID)
	if err != nil {
		return 0, err
	}

	workout, err := uc.workoutsRepo.Get(workoutID)
	if err != nil {
		return 0, err
	}

	// обновляем следом идущие индексы
	for index, ex := range workout.Exercises {
		if ex.Index > deleteExIndex {
			workout.Exercises[index].Index--
			_ = uc.exercisesRepo.Save(&workout.Exercises[index])
		}
	}

	session, err := uc.sessionsRepo.GetByWorkoutID(workoutID)
	if err != nil {
		return 0, err
	}

	if session.CurrentExerciseIndex >= len(workout.Exercises) {
		session.CurrentExerciseIndex = len(workout.Exercises) - 1
		_ = uc.sessionsRepo.Save(&session)
	}

	return workoutID, nil
}
