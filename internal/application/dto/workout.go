package dto

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
)

type WorkoutItem struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	StartedAt string `json:"started_at"`
	Duration  string `json:"duration"`
	Completed bool   `json:"completed"`
	Status    string `json:"status"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

type ShowMyWorkoutsResult struct {
	Items      []WorkoutItem `json:"items"`
	Pagination Pagination    `json:"pagination"`
}

type ConfirmDeleteWorkout struct {
	WorkoutID   int64
	DayTypeName string
}

type DeleteWorkout struct {
}

type ConfirmFinishWorkout struct {
	DayType models.WorkoutDayType
}

type FinishWorkout struct {
	WorkoutID int64
}

type CreateWorkout struct {
	WorkoutID int64
}

type StartWorkout struct {
}

type WorkoutProgress struct {
	Workout *models.WorkoutDay

	TotalExercises     int
	CompletedExercises int

	TotalSets     int
	CompletedSets int

	ProgressPercent int
	RemainingMin    *int

	SessionStarted bool
}

type WorkoutStatistic struct {
	DayType            models.WorkoutDayType
	WorkoutDay         models.WorkoutDay
	TotalWeight        float64
	CompletedExercises int
	CardioTime         int
	ExerciseTypesMap   map[int64]models.ExerciseType
	ExerciseWeightMap  map[int64]float64
	ExerciseTimeMap    map[int64]int
}

type ShowWorkoutByUserID struct {
	Workouts []models.WorkoutDay
	User     *models.User
}
