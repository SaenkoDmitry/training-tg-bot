package router

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/admins"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/changes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/daytypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/exports"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/sets"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/stats"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/timers"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/workouts"
	userusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/users"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	bot *tgbotapi.BotAPI

	workoutsHandler  *workouts.Handler
	timersHandler    *timers.Handler
	statsHandler     *stats.Handler
	setsHandler      *sets.Handler
	programsHandler  *programs.Handler
	exportsHandler   *exports.Handler
	exercisesHandler *exercises.Handler
	changesHandler   *changes.Handler
	dayTypesHandler  *daytypes.Handler
	adminsHandler    *admins.Handler
	createUserUC     *userusecases.CreateUseCase
	getUserUC        *userusecases.GetUseCase
}

func New(
	bot *tgbotapi.BotAPI,
	createUserUC *userusecases.CreateUseCase,
	getUserUC *userusecases.GetUseCase,
	adminsHandler *admins.Handler,
	workoutsHandler *workouts.Handler,
	timersHandler *timers.Handler,
	statsHandler *stats.Handler,
	setsHandler *sets.Handler,
	programsHandler *programs.Handler,
	exportsHandler *exports.Handler,
	exercisesHandler *exercises.Handler,
	changesHandler *changes.Handler,
	dayTypesHandler *daytypes.Handler,
) *Router {
	return &Router{
		bot:              bot,
		createUserUC:     createUserUC,
		getUserUC:        getUserUC,
		adminsHandler:    adminsHandler,
		workoutsHandler:  workoutsHandler,
		timersHandler:    timersHandler,
		statsHandler:     statsHandler,
		setsHandler:      setsHandler,
		programsHandler:  programsHandler,
		exportsHandler:   exportsHandler,
		exercisesHandler: exercisesHandler,
		changesHandler:   changesHandler,
		dayTypesHandler:  dayTypesHandler,
	}
}

func (r *Router) HandleUpdate(update tgbotapi.Update) {
	defer func() {
		if recovery := recover(); recovery != nil {
			fmt.Printf("Recovered: %v\n", recovery)
		}
	}()
	switch {
	case update.Message != nil:
		r.routeMessage(update.Message)

	case update.CallbackQuery != nil:
		r.routeCallback(update.CallbackQuery)
	}
}
