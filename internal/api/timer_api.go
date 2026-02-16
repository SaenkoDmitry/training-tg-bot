package api

import (
	"encoding/json"
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (s *serviceImpl) StartTimer(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req struct {
		WorkoutID int64 `json:"workout_id"`
		Seconds   int   `json:"seconds"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	timer, err := s.timerManager.Start(claims.ChatID, req.WorkoutID, req.Seconds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(timer)
}

func (s *serviceImpl) CancelTimer(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	timerIDStr := chi.URLParam(r, "id")
	timerID, err := strconv.ParseInt(timerIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.timerManager.Cancel(timerID, claims.ChatID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			w.WriteHeader(http.StatusNotFound)
		case err.Error() == "forbidden":
			w.WriteHeader(http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
