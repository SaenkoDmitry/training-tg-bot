package ai

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisegrouptypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/measurements"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
)

const (
	defaultHistoryLimit      = 15
	defaultMeasurementLimit  = 10
	defaultHistoryPeriodDays = 90
)

type BuildProgramPromptUseCase struct {
	usersRepo              users.Repo
	programsRepo           programs.Repo
	workoutsRepo           workouts.Repo
	measurementsRepo       measurements.Repo
	exerciseTypesRepo      exercisetypes.Repo
	exerciseGroupTypesRepo exercisegrouptypes.Repo
}

func NewBuildProgramPromptUseCase(
	usersRepo users.Repo,
	programsRepo programs.Repo,
	workoutsRepo workouts.Repo,
	measurementsRepo measurements.Repo,
	exerciseTypesRepo exercisetypes.Repo,
	exerciseGroupTypesRepo exercisegrouptypes.Repo,
) *BuildProgramPromptUseCase {
	return &BuildProgramPromptUseCase{
		usersRepo:              usersRepo,
		programsRepo:           programsRepo,
		workoutsRepo:           workoutsRepo,
		measurementsRepo:       measurementsRepo,
		exerciseTypesRepo:      exerciseTypesRepo,
		exerciseGroupTypesRepo: exerciseGroupTypesRepo,
	}
}

func (uc *BuildProgramPromptUseCase) Execute(userID int64, req dto.AIProgramPromptRequest) (*dto.AIProgramPromptResponse, error) {
	normalizeRequest(&req)

	ctx, err := uc.BuildContext(userID, req)
	if err != nil {
		return nil, err
	}

	payload, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return nil, err
	}

	return &dto.AIProgramPromptResponse{
		Context:      ctx,
		SystemPrompt: buildSystemPrompt(),
		UserPrompt:   buildUserPrompt(req, string(payload)),
		OutputSchema: buildOutputSchema(),
	}, nil
}

func (uc *BuildProgramPromptUseCase) BuildContext(userID int64, req dto.AIProgramPromptRequest) (*dto.AIProgramContext, error) {
	user, err := uc.usersRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	groups, err := uc.exerciseGroupTypesRepo.GetAll()
	if err != nil {
		return nil, err
	}
	groupsMap := make(map[string]string, len(groups))
	groupCodes := make([]string, 0, len(groups))
	for _, g := range groups {
		groupsMap[g.Code] = g.Name
		groupCodes = append(groupCodes, g.Code)
	}
	sort.Strings(groupCodes)

	exerciseTypes, err := uc.exerciseTypesRepo.GetAll()
	if err != nil {
		return nil, err
	}
	exerciseTypeMap := make(map[int64]models.ExerciseType, len(exerciseTypes))
	catalog := make([]dto.AIExerciseCatalogItem, 0, len(exerciseTypes))
	for _, ex := range exerciseTypes {
		exerciseTypeMap[ex.ID] = ex
		catalog = append(catalog, dto.AIExerciseCatalogItem{
			ID:              ex.ID,
			Name:            ex.Name,
			GroupCode:       ex.ExerciseGroupTypeCode,
			GroupName:       groupsMap[ex.ExerciseGroupTypeCode],
			Accent:          ex.Accent,
			SecondaryAccent: ex.SecondaryAccent,
			Units:           ex.Units,
			RestInSeconds:   ex.RestInSeconds,
		})
	}

	currentProgram, err := uc.buildCurrentProgram(user, req, exerciseTypeMap, groupsMap)
	if err != nil {
		return nil, err
	}

	workoutObjs, err := uc.workoutsRepo.Find(userID, 0, defaultHistoryLimit)
	if err != nil {
		return nil, err
	}

	measurementObjs, err := uc.measurementsRepo.FindAllLimitOffset(userID, defaultMeasurementLimit, 0)
	if err != nil {
		return nil, err
	}

	return &dto.AIProgramContext{
		GeneratedFor:        time.Now().UTC().Format(time.RFC3339),
		Request:             req,
		UserProfile:         buildUserProfile(user),
		CurrentProgram:      currentProgram,
		TrainingSummary:     buildTrainingSummary(workoutObjs, groupsMap),
		MeasurementSummary:  buildMeasurementSummary(measurementObjs),
		ExerciseCatalog:     catalog,
		AvailableGroupCodes: groupCodes,
		CompatibilityNotes: []string{
			"AI must return exercises by exercise_type_id from provided_exercise_catalog only.",
			"Program days can be converted to the current preset string format after backend validation.",
			"Progression rules should be stored separately from static preset sets because preset currently contains only planned sets.",
		},
	}, nil
}

func (uc *BuildProgramPromptUseCase) buildCurrentProgram(
	user *models.User,
	req dto.AIProgramPromptRequest,
	exerciseTypeMap map[int64]models.ExerciseType,
	groupsMap map[string]string,
) (*dto.AICurrentProgramContext, error) {
	var programID *int64
	if req.ProgramID != nil {
		programID = req.ProgramID
	} else if user.ActiveProgramID != nil {
		programID = user.ActiveProgramID
	}
	if programID == nil {
		return nil, nil
	}

	programObj, err := uc.programsRepo.Get(*programID)
	if err != nil {
		return nil, err
	}
	if programObj.UserID != user.ID {
		return nil, fmt.Errorf("program does not belong to user")
	}

	program := &dto.AICurrentProgramContext{
		ID:       programObj.ID,
		Name:     programObj.Name,
		IsActive: user.ActiveProgramID != nil && *user.ActiveProgramID == programObj.ID,
		Days:     make([]dto.AIProgramDayContext, 0, len(programObj.DayTypes)),
	}

	for _, day := range programObj.DayTypes {
		dayCtx := dto.AIProgramDayContext{
			ID:        day.ID,
			Name:      day.Name,
			Preset:    day.Preset,
			Exercises: make([]dto.AIProgramExerciseContext, 0),
		}
		for _, ex := range utils.SplitPreset(day.Preset) {
			exerciseType, ok := exerciseTypeMap[ex.ID]
			if !ok {
				continue
			}
			sets := make([]dto.AIProgramSetContext, 0, len(ex.Sets))
			for _, set := range ex.Sets {
				sets = append(sets, dto.AIProgramSetContext{
					Reps:    set.Reps,
					Weight:  set.Weight,
					Minutes: set.Minutes,
				})
			}
			dayCtx.Exercises = append(dayCtx.Exercises, dto.AIProgramExerciseContext{
				ExerciseTypeID: exerciseType.ID,
				Name:           exerciseType.Name,
				GroupCode:      exerciseType.ExerciseGroupTypeCode,
				GroupName:      groupsMap[exerciseType.ExerciseGroupTypeCode],
				Units:          exerciseType.Units,
				RestInSeconds:  exerciseType.RestInSeconds,
				Sets:           sets,
			})
		}
		program.Days = append(program.Days, dayCtx)
	}

	return program, nil
}

func buildUserProfile(user *models.User) dto.AIUserProfileContext {
	profile := dto.AIUserProfileContext{
		UserID:    user.ID,
		FirstName: user.FirstName,
		Gender:    user.Gender,
		HeightCm:  user.HeightCm,
		WeightKg:  user.WeightKg,
	}
	if user.BirthDate != nil {
		age := time.Now().Year() - user.BirthDate.Year()
		birthdayThisYear := time.Date(time.Now().Year(), user.BirthDate.Month(), user.BirthDate.Day(), 0, 0, 0, 0, time.Local)
		if time.Now().Before(birthdayThisYear) {
			age--
		}
		profile.Age = &age
	}
	return profile
}

type exerciseAccumulator struct {
	ExerciseTypeID int64
	Name           string
	GroupCode      string
	GroupName      string
	Units          string
	Sessions       []dto.AIExerciseSessionContext
	TotalSets      int
	CompletedSets  int
	BestWeight     float32
	BestVolume     float32
	BestMinutes    int
	BestMeters     int
}

func buildTrainingSummary(workoutObjs []models.WorkoutDay, groupsMap map[string]string) dto.AITrainingSummaryContext {
	summary := dto.AITrainingSummaryContext{
		PeriodDays:       defaultHistoryPeriodDays,
		LoadedWorkouts:   len(workoutObjs),
		ExerciseProgress: make([]dto.AIExerciseProgressContext, 0),
	}
	if len(workoutObjs) == 0 {
		summary.Consistency = "no_history"
		return summary
	}

	accs := make(map[int64]*exerciseAccumulator)
	var firstDate, lastDate time.Time

	for i, workout := range workoutObjs {
		if i == 0 || workout.StartedAt.After(lastDate) {
			lastDate = workout.StartedAt
		}
		if i == 0 || workout.StartedAt.Before(firstDate) {
			firstDate = workout.StartedAt
		}
		if workout.Completed {
			summary.CompletedWorkouts++
		}

		for _, exercise := range workout.Exercises {
			if exercise.ExerciseType == nil {
				continue
			}
			exType := exercise.ExerciseType
			acc := accs[exType.ID]
			if acc == nil {
				acc = &exerciseAccumulator{
					ExerciseTypeID: exType.ID,
					Name:           exType.Name,
					GroupCode:      exType.ExerciseGroupTypeCode,
					GroupName:      groupsMap[exType.ExerciseGroupTypeCode],
					Units:          exType.Units,
				}
				accs[exType.ID] = acc
			}

			session := dto.AIExerciseSessionContext{
				Date:      workout.StartedAt.Format("2006-01-02"),
				Completed: workout.Completed,
				Sets:      make([]dto.AICompletedSetContext, 0, len(exercise.Sets)),
			}
			if workout.WorkoutDayType != nil {
				session.DayName = workout.WorkoutDayType.Name
			}

			for _, set := range exercise.Sets {
				reps := set.GetRealReps()
				weight := set.GetRealWeight()
				minutes := set.GetRealMinutes()
				meters := set.GetRealMeters()
				volume := weight * float32(reps)
				session.Volume += volume
				if weight > acc.BestWeight {
					acc.BestWeight = weight
				}
				if session.Volume > acc.BestVolume {
					acc.BestVolume = session.Volume
				}
				if minutes > acc.BestMinutes {
					acc.BestMinutes = minutes
				}
				if meters > acc.BestMeters {
					acc.BestMeters = meters
				}

				acc.TotalSets++
				if set.Completed {
					acc.CompletedSets++
				}

				session.Sets = append(session.Sets, dto.AICompletedSetContext{
					Reps:      reps,
					Weight:    weight,
					Minutes:   minutes,
					Meters:    meters,
					Completed: set.Completed,
				})
			}
			acc.Sessions = append(acc.Sessions, session)
		}
	}

	summary.HasHistory = summary.CompletedWorkouts > 0
	summary.FirstWorkoutDate = firstDate.Format("2006-01-02")
	summary.LastWorkoutDate = lastDate.Format("2006-01-02")
	periodDays := int(math.Ceil(lastDate.Sub(firstDate).Hours()/24)) + 1
	if periodDays < 1 {
		periodDays = 1
	}
	if periodDays > defaultHistoryPeriodDays {
		periodDays = defaultHistoryPeriodDays
	}
	summary.PeriodDays = periodDays
	summary.AvgWorkoutsPerWeek = roundFloat(float64(summary.CompletedWorkouts)/float64(periodDays)*7, 2)
	summary.Consistency = consistencyLabel(summary.AvgWorkoutsPerWeek)

	progress := make([]dto.AIExerciseProgressContext, 0, len(accs))
	for _, acc := range accs {
		sort.Slice(acc.Sessions, func(i, j int) bool {
			return acc.Sessions[i].Date > acc.Sessions[j].Date
		})
		lastSessions := acc.Sessions
		if len(lastSessions) > 3 {
			lastSessions = lastSessions[:3]
		}
		completionRate := 0.0
		if acc.TotalSets > 0 {
			completionRate = roundFloat(float64(acc.CompletedSets)/float64(acc.TotalSets), 2)
		}
		trend := calculateTrend(acc.Sessions)
		progress = append(progress, dto.AIExerciseProgressContext{
			ExerciseTypeID:       acc.ExerciseTypeID,
			Name:                 acc.Name,
			GroupCode:            acc.GroupCode,
			GroupName:            acc.GroupName,
			Units:                acc.Units,
			SessionsCount:        len(acc.Sessions),
			RecentCompletionRate: completionRate,
			Trend:                trend,
			BestWeight:           acc.BestWeight,
			BestVolume:           acc.BestVolume,
			BestMinutes:          acc.BestMinutes,
			BestMeters:           acc.BestMeters,
			LastSessions:         lastSessions,
			RecommendationSignal: recommendationSignal(completionRate, trend, len(acc.Sessions)),
		})
	}
	sort.Slice(progress, func(i, j int) bool {
		if progress[i].SessionsCount == progress[j].SessionsCount {
			return progress[i].ExerciseTypeID < progress[j].ExerciseTypeID
		}
		return progress[i].SessionsCount > progress[j].SessionsCount
	})
	summary.ExerciseProgress = progress

	return summary
}

func buildMeasurementSummary(measurements []models.Measurement) dto.AIMeasurementSummaryContext {
	summary := dto.AIMeasurementSummaryContext{LoadedMeasurements: len(measurements)}
	if len(measurements) == 0 {
		return summary
	}
	summary.HasMeasurements = true
	newest := measurements[0]
	oldest := measurements[len(measurements)-1]
	summary.LastDate = newest.CreatedAt.Format("2006-01-02")
	lastWeight := roundFloat(float64(newest.Weight)/1000, 2)
	summary.LastWeightKg = &lastWeight
	if oldest.Weight > 0 && newest.Weight > 0 && newest.ID != oldest.ID {
		change := roundFloat(float64(newest.Weight-oldest.Weight)/1000, 2)
		summary.WeightChangeKg = &change
	}
	return summary
}

func calculateTrend(sessions []dto.AIExerciseSessionContext) string {
	if len(sessions) < 3 {
		return "insufficient_data"
	}
	sort.Slice(sessions, func(i, j int) bool { return sessions[i].Date > sessions[j].Date })
	recentCount := min(3, len(sessions))
	olderStart := recentCount
	olderEnd := min(olderStart+3, len(sessions))
	if olderEnd <= olderStart {
		return "insufficient_data"
	}
	recentAvg := avgVolume(sessions[:recentCount])
	olderAvg := avgVolume(sessions[olderStart:olderEnd])
	if olderAvg == 0 && recentAvg == 0 {
		return "flat"
	}
	if recentAvg > olderAvg*1.05 {
		return "up"
	}
	if recentAvg < olderAvg*0.95 {
		return "down"
	}
	return "flat"
}

func avgVolume(sessions []dto.AIExerciseSessionContext) float32 {
	if len(sessions) == 0 {
		return 0
	}
	var sum float32
	for _, s := range sessions {
		sum += s.Volume
	}
	return sum / float32(len(sessions))
}

func consistencyLabel(avg float64) string {
	switch {
	case avg == 0:
		return "no_completed_workouts"
	case avg < 1.5:
		return "irregular"
	case avg < 3:
		return "moderate"
	default:
		return "consistent"
	}
}

func recommendationSignal(completionRate float64, trend string, sessionsCount int) string {
	if sessionsCount < 2 {
		return "insufficient_data"
	}
	if completionRate >= 0.9 && (trend == "up" || trend == "flat") {
		return "can_progress_slowly"
	}
	if completionRate < 0.7 || trend == "down" {
		return "reduce_or_stabilize_load"
	}
	return "keep_and_observe"
}

func normalizeRequest(req *dto.AIProgramPromptRequest) {
	req.Mode = strings.TrimSpace(req.Mode)
	if req.Mode == "" {
		req.Mode = "create_program"
	}
	if req.DaysPerWeek <= 0 {
		req.DaysPerWeek = 3
	}
	if req.SessionDurationMin <= 0 {
		req.SessionDurationMin = 60
	}
	if req.Level == "" {
		req.Level = "beginner"
	}
	if req.Location == "" {
		req.Location = "gym"
	}
}

func buildSystemPrompt() string {
	return strings.TrimSpace(`Ты — помощник по планированию тренировок внутри фитнес-трекера Form Journey.

Задача: составить новую программу или предложить корректировку существующей программы на основе JSON-контекста пользователя.

Жесткие правила:
1. Используй только упражнения из provided_exercise_catalog.
2. Никогда не придумывай exercise_type_id.
3. Не давай медицинских диагнозов и не лечи травмы.
4. Если пользователь указал боль, травму или заболевание, добавь предупреждение обратиться к специалисту.
5. Не повышай объем или интенсивность резко; прогрессия должна быть плавной.
6. Учитывай уровень, регулярность тренировок, цель, ограничения, доступное время и историю выполнения подходов.
7. Для каждого упражнения верни reason и progression_rule.
8. Ответ должен быть валидным JSON по output_schema. Не добавляй Markdown и свободный текст вне JSON.`)
}

func buildUserPrompt(req dto.AIProgramPromptRequest, contextJSON string) string {
	modeDescription := "составь новую тренировочную программу"
	if req.Mode == "improve_existing_program" {
		modeDescription = "скорректируй текущую тренировочную программу"
	}
	return fmt.Sprintf(strings.TrimSpace(`Режим: %s.

Пожелание пользователя: %s

JSON-контекст приложения:
%s

Верни JSON с программой/изменениями. Если истории тренировок или текущей программы нет, используй безопасные предположения и явно укажи это в validation_notes.`), modeDescription, emptyFallback(req.Notes, "пользователь не оставил отдельный комментарий"), contextJSON)
}

func buildOutputSchema() any {
	return map[string]any{
		"type":     "object",
		"required": []string{"summary", "program", "warnings", "validation_notes"},
		"properties": map[string]any{
			"summary": map[string]any{"type": "string"},
			"program": map[string]any{
				"type":     "object",
				"required": []string{"name", "days"},
				"properties": map[string]any{
					"name": map[string]any{"type": "string"},
					"days": map[string]any{"type": "array", "items": map[string]any{
						"type":     "object",
						"required": []string{"name", "focus", "exercises"},
						"properties": map[string]any{
							"name":  map[string]any{"type": "string"},
							"focus": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
							"exercises": map[string]any{"type": "array", "items": map[string]any{
								"type":     "object",
								"required": []string{"exercise_type_id", "sets", "reason", "progression_rule"},
								"properties": map[string]any{
									"exercise_type_id": map[string]any{"type": "integer"},
									"sets": map[string]any{"type": "array", "items": map[string]any{"type": "object", "properties": map[string]any{
										"reps": map[string]any{"type": "integer"}, "weight": map[string]any{"type": "number"}, "minutes": map[string]any{"type": "integer"}, "meters": map[string]any{"type": "integer"},
									}}},
									"rest_in_seconds":  map[string]any{"type": "integer"},
									"reason":           map[string]any{"type": "string"},
									"progression_rule": map[string]any{"type": "string"},
								},
							}},
						},
					}},
				},
			},
			"changes":          map[string]any{"type": "array", "items": map[string]any{"type": "object"}},
			"warnings":         map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"validation_notes": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
		},
	}
}

func roundFloat(value float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Round(value*factor) / factor
}

func emptyFallback(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
