package main

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/api"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func initServer(container *usecase.Container) {
	r := chi.NewRouter()

	// Средства middleware chi
	r.Use(middleware.Logger)    // лог запросов
	r.Use(middleware.Recoverer) // recovery от паник
	r.Use(middleware.RequestID) // уникальный ID запроса

	s := api.New(container)

	r.Route("/api/login", func(r chi.Router) {
		r.Post("/", api.LoginHandler)
	})
	r.Route("/api/logout", func(r chi.Router) {
		r.Post("/", api.LogoutHandler)
	})

	r.Route("/api/me", func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Get("/", api.MeHandler)
	})

	r.Route("/api/workouts", func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Get("/", s.GetAllWorkouts)          // GET /api/workouts
		r.Get("/{workout_id}", s.ReadWorkout) // GET /api/workouts/123
	})

	r.Route("/api/measurements", func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Get("/", s.GetMeasurements)
		r.Post("/", s.CreateMeasurement)
	})

	r.Route("/api/exercise-groups", func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Get("/", s.GetExerciseGroups)
		r.Get("/{group}", s.GetExerciseTypesByGroup)
	})

	// UI (React build)
	r.Handle("/*", serveSPA("./ui"))

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}

func serveSPA(staticDir string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(staticDir))

	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(staticDir, r.URL.Path)

		// если файл существует — отдать его (js/css/png/...)
		if _, err := os.Stat(path); err == nil {
			fs.ServeHTTP(w, r)
			return
		}

		// иначе всегда index.html (SPA fallback)
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	}
}
