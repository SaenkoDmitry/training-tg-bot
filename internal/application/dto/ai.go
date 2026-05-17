package dto

type AIProgramPromptRequest struct {
	Mode               string   `json:"mode"`
	ProgramID          *int64   `json:"program_id"`
	Goal               string   `json:"goal"`
	Level              string   `json:"level"`
	DaysPerWeek        int      `json:"days_per_week"`
	SessionDurationMin int      `json:"session_duration_min"`
	Location           string   `json:"location"`
	Limitations        []string `json:"limitations"`
	Focus              []string `json:"focus"`
	Notes              string   `json:"notes"`
}

type AIProgramPromptResponse struct {
	Context      *AIProgramContext `json:"context"`
	SystemPrompt string            `json:"system_prompt"`
	UserPrompt   string            `json:"user_prompt"`
	OutputSchema any               `json:"output_schema"`
}

type AIProgramContext struct {
	GeneratedFor        string                      `json:"generated_for"`
	Request             AIProgramPromptRequest      `json:"request"`
	UserProfile         AIUserProfileContext        `json:"user_profile"`
	CurrentProgram      *AICurrentProgramContext    `json:"current_program,omitempty"`
	TrainingSummary     AITrainingSummaryContext    `json:"training_summary"`
	MeasurementSummary  AIMeasurementSummaryContext `json:"measurement_summary"`
	ExerciseCatalog     []AIExerciseCatalogItem     `json:"provided_exercise_catalog"`
	AvailableGroupCodes []string                    `json:"available_group_codes"`
	CompatibilityNotes  []string                    `json:"compatibility_notes"`
}

type AIUserProfileContext struct {
	UserID    int64    `json:"user_id"`
	FirstName string   `json:"first_name,omitempty"`
	Age       *int     `json:"age,omitempty"`
	Gender    *string  `json:"gender,omitempty"`
	HeightCm  *int     `json:"height_cm,omitempty"`
	WeightKg  *float64 `json:"weight_kg,omitempty"`
}

type AICurrentProgramContext struct {
	ID       int64                 `json:"id"`
	Name     string                `json:"name"`
	IsActive bool                  `json:"is_active"`
	Days     []AIProgramDayContext `json:"days"`
}

type AIProgramDayContext struct {
	ID        int64                      `json:"id"`
	Name      string                     `json:"name"`
	Preset    string                     `json:"preset,omitempty"`
	Exercises []AIProgramExerciseContext `json:"exercises"`
}

type AIProgramExerciseContext struct {
	ExerciseTypeID int64                 `json:"exercise_type_id"`
	Name           string                `json:"name"`
	GroupCode      string                `json:"group_code,omitempty"`
	GroupName      string                `json:"group_name,omitempty"`
	Units          string                `json:"units"`
	RestInSeconds  int                   `json:"rest_in_seconds"`
	Sets           []AIProgramSetContext `json:"sets"`
}

type AIProgramSetContext struct {
	Reps    int     `json:"reps,omitempty"`
	Weight  float32 `json:"weight,omitempty"`
	Minutes int     `json:"minutes,omitempty"`
	Meters  int     `json:"meters,omitempty"`
}

type AITrainingSummaryContext struct {
	PeriodDays         int                         `json:"period_days"`
	LoadedWorkouts     int                         `json:"loaded_workouts"`
	CompletedWorkouts  int                         `json:"completed_workouts"`
	AvgWorkoutsPerWeek float64                     `json:"avg_workouts_per_week"`
	FirstWorkoutDate   string                      `json:"first_workout_date,omitempty"`
	LastWorkoutDate    string                      `json:"last_workout_date,omitempty"`
	HasHistory         bool                        `json:"has_history"`
	Consistency        string                      `json:"consistency"`
	ExerciseProgress   []AIExerciseProgressContext `json:"exercise_progress"`
}

type AIExerciseProgressContext struct {
	ExerciseTypeID       int64                      `json:"exercise_type_id"`
	Name                 string                     `json:"name"`
	GroupCode            string                     `json:"group_code,omitempty"`
	GroupName            string                     `json:"group_name,omitempty"`
	Units                string                     `json:"units"`
	SessionsCount        int                        `json:"sessions_count"`
	RecentCompletionRate float64                    `json:"recent_completion_rate"`
	Trend                string                     `json:"trend"`
	BestWeight           float32                    `json:"best_weight,omitempty"`
	BestVolume           float32                    `json:"best_volume,omitempty"`
	BestMinutes          int                        `json:"best_minutes,omitempty"`
	BestMeters           int                        `json:"best_meters,omitempty"`
	LastSessions         []AIExerciseSessionContext `json:"last_sessions"`
	RecommendationSignal string                     `json:"recommendation_signal"`
}

type AIExerciseSessionContext struct {
	Date      string                  `json:"date"`
	DayName   string                  `json:"day_name,omitempty"`
	Sets      []AICompletedSetContext `json:"sets"`
	Volume    float32                 `json:"volume,omitempty"`
	Completed bool                    `json:"completed"`
}

type AICompletedSetContext struct {
	Reps      int     `json:"reps,omitempty"`
	Weight    float32 `json:"weight,omitempty"`
	Minutes   int     `json:"minutes,omitempty"`
	Meters    int     `json:"meters,omitempty"`
	Completed bool    `json:"completed"`
}

type AIMeasurementSummaryContext struct {
	HasMeasurements    bool     `json:"has_measurements"`
	LastDate           string   `json:"last_date,omitempty"`
	LastWeightKg       *float64 `json:"last_weight_kg,omitempty"`
	WeightChangeKg     *float64 `json:"weight_change_kg,omitempty"`
	LoadedMeasurements int      `json:"loaded_measurements"`
}

type AIExerciseCatalogItem struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	GroupCode       string `json:"group_code"`
	GroupName       string `json:"group_name"`
	Accent          string `json:"accent,omitempty"`
	SecondaryAccent string `json:"secondary_accent,omitempty"`
	Units           string `json:"units"`
	RestInSeconds   int    `json:"rest_in_seconds"`
}

type AIApplyProgramRequest struct {
	Program  AIGeneratedProgram `json:"program"`
	Warnings []string           `json:"warnings"`
	Notes    []string           `json:"validation_notes"`
	Activate bool               `json:"activate"`
}

type AIApplyProgramResponse struct {
	ProgramID  int64  `json:"program_id"`
	Name       string `json:"name"`
	DaysCount  int    `json:"days_count"`
	RulesCount int    `json:"rules_count"`
}

type AIGeneratedProgram struct {
	Name string                  `json:"name"`
	Days []AIGeneratedProgramDay `json:"days"`
}

type AIGeneratedProgramDay struct {
	Name      string                       `json:"name"`
	Focus     []string                     `json:"focus"`
	Exercises []AIGeneratedProgramExercise `json:"exercises"`
}

type AIGeneratedProgramExercise struct {
	ExerciseTypeID  int64                   `json:"exercise_type_id"`
	Sets            []AIGeneratedProgramSet `json:"sets"`
	RestInSeconds   int                     `json:"rest_in_seconds"`
	Reason          string                  `json:"reason"`
	ProgressionRule string                  `json:"progression_rule"`
}

type AIGeneratedProgramSet struct {
	Reps    int     `json:"reps"`
	Weight  float32 `json:"weight"`
	Minutes int     `json:"minutes"`
	Meters  int     `json:"meters"`
}
