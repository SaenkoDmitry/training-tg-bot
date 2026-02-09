package api

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase"
	"net/http"
)

type Service interface {
	GetAllWorkouts(w http.ResponseWriter, r *http.Request)
	ReadWorkout(w http.ResponseWriter, r *http.Request)
	GetMeasurements(w http.ResponseWriter, r *http.Request)
	CreateMeasurement(w http.ResponseWriter, r *http.Request)
	GetExerciseGroups(w http.ResponseWriter, _ *http.Request)
	GetExerciseTypesByGroup(w http.ResponseWriter, r *http.Request)
}

type serviceImpl struct {
	container *usecase.Container
}

func New(container *usecase.Container) Service {
	return &serviceImpl{
		container: container,
	}
}
