interface Workout {
    id: number;
    name: string;
    started_at?: string;
    duration?: string;
    status?: string;
    completed: boolean;
    cardio_distance?: number;
    cardio_time?: number;
    has_valid_cardio_data?: boolean;
}

interface ShowMyWorkoutsResult {
    items: Workout[];
    pagination: Pagination;
}

interface Pagination {
    limit: number;
    offset: number;
    total: number;
}

interface User {
    id: number;
    first_name: string;
    username?: string;
    photo_url?: string;
}

interface ExerciseType {
    id: number;
    name: string;
    url: string;
    exercise_group_type_code: string;
    rest_in_seconds: number;
    accent: string;
    secondary_accent: string;
    units: string;
    description: string;
}

interface ExerciseStat {
    id: number;
    date: string;
    sets?: FormattedSet[];
}

interface ExerciseStatsResponse {
    items: ExerciseStat[];
    // если у тебя есть total / has_more — добавь сюда
    total?: number;
}

interface Set {
    ID: number;
    Reps: number;
    FactReps: number;
    Weight: number;
    FactWeight: number;
    Minutes: number;
    FactMinutes: number;
    Meters: number;
    FactMeters: number;
    Completed: boolean;
    CompletedAt: string;
    Index: number;
}

interface Exercise {
    ID: number;
    ExerciseType: ExerciseType;
    Sets: Set[];
    Index: number;
}

interface ReadWorkoutDTO {
    progress: WorkoutProgress;
    stats: WorkoutStatistic;
    user_first_name?: string;
}

interface ShareDTO {
    token: string;
    share_url: string;
    created_at: string;
}

interface WorkoutProgress {
    workout: FormattedWorkout;
    TotalExercises: number;
    CompletedExercises: number;
    TotalSets: number;
    CompletedSets: number;
    ProgressPercent: number;
    RemainingMin: number | null;
    SessionStarted: boolean;

    EstimatedCalories?: number;
    EstimatedDurationMin?: number;
    UserWeightKg?: number;
}

interface WorkoutStatistic {
    day_type: any;
    workout_day: any;
    total_weight: number;
    completed_exercises: number;
    cardio_time: number; // в минутах
    exercise_map: Record<number, FormattedExercise>;
    exercise_weight_map: Record<number, number>;
    exercise_time_map: Record<number, number>;
}

interface FormattedWorkout {
    id: number
    user_id: number
    status: string
    started_at: string
    duration: string
    ended_at: string
    day_type_name: string
    completed: boolean
    exercises: FormattedExercise[]
}

interface FormattedExercise {
    id: number
    name: string
    url: string
    group_name: string
    rest_in_seconds: number
    accent: string
    units: string
    description: string
    index: number
    sets: FormattedSet[]
    sum_weight: number
}

interface FormattedSet {
    id: number;
    reps: number;
    fact_reps: number;
    weight: number;
    fact_weight: number;
    minutes: number;
    fact_minutes: number;
    meters: number;
    fact_meters: number;
    formatted_string: string;
    completed: boolean;
    completed_at: string;
    index: number;
}

interface Measurement {
    id: number;
    user_id: number;
    created_at: string; // ISO string
    shoulders: number;
    chest: number;
    hand_left: number;
    hand_right: number;
    waist: number;
    buttocks: number;
    hip_left: number;
    hip_right: number;
    calf_left: number;
    calf_right: number;
    weight: number;
}

interface MeasurementTypeDTO {
    code: string;
    name: string;
}

interface FindWithOffsetLimitMeasurement {
    items: Measurement[];
    count: number;
}


interface ToCreateMeasurement {
    user_id: number;
    shoulders: number;
    chest: number;
    hand_left: number;
    hand_right: number;
    waist: number;
    buttocks: number;
    hip_left: number;
    hip_right: number;
    calf_left: number;
    calf_right: number;
    weight: number;
}

type Group = {
    code: string;
    name: string;
};

interface WorkoutDayTypeDTO {
    id: number;
    program_id: number;   // оставил snake_case, чтобы совпадало с API
    name: string;
    preset: string;
    created_at: string;   // ISO string
}

interface StartWorkoutDTO {
    workout_id: number;
}

interface CurrentExerciseSession {
    exercise: FormattedExercise;
    exercise_type: ExerciseType;
    day_type: WorkoutDayTypeDTO;
    workout: FormattedWorkout;
    exercise_index: number;
}

interface ProgramDTO {
    id: number;
    user_id: number;
    name: string;
    summary: string;
    warnings: string[];
    notes: string[];
    created_at: string;
    day_types: WorkoutDayTypeDTO[];
}

type SetDTO = {
    reps: number;
    weight: number;
    minutes: number;
    meters: number;
};

type ExerciseDTO = {
    id: number;
    name: string;
    units: string;
    sets: SetDTO[];
};

type PresetListDTO = {
    exercises: ExerciseDTO[];
};

type TimerDTO = {
    id: number;
}

type IconDTO = {
    name: string;
}

interface UserProfile {
    id: number;
    first_name: string;
    last_name?: string;
    username?: string;
    email?: string;
    icon: string;
    birth_date?: string;  // "YYYY-MM-DD"
    gender?: "male" | "female";
    weight_kg?: number;
    height_cm?: number;
}

interface UpdateProfileRequest {
    birth_date?: string | null;
    gender?: "male" | "female" | null;
    weight_kg?: number | null;
    height_cm?: number | null;
}

interface PreviewCaloriesResponse {
    calories: number | null;
    duration_min: number | null;
    weight_note: "actual" | "current" | "missing";
}

interface AIProgramPromptRequest {
    mode: "create_program" | "improve_existing_program";
    program_id?: number | null;
    goal: string;
    level: string;
    days_per_week: number;
    session_duration_min: number;
    location: string;
    limitations: string[];
    focus: string[];
    notes: string;
}

interface AIProgramPromptResponse {
    context: AIProgramContext;
    system_prompt: string;
    user_prompt: string;
    output_schema: unknown;
}

interface AIProgramContext {
    generated_for: string;
    request: AIProgramPromptRequest;
    user_profile: {
        user_id: number;
        first_name?: string;
        age?: number;
        gender?: string;
        height_cm?: number;
        weight_kg?: number;
    };
    current_program?: {
        id: number;
        name: string;
        is_active: boolean;
        days: Array<{
            id: number;
            name: string;
            preset?: string;
            exercises: Array<{
                exercise_type_id: number;
                name: string;
                group_code?: string;
                group_name?: string;
                units: string;
                rest_in_seconds: number;
                sets: Array<{reps?: number; weight?: number; minutes?: number; meters?: number}>;
            }>;
        }>;
    };
    training_summary: {
        period_days: number;
        loaded_workouts: number;
        completed_workouts: number;
        avg_workouts_per_week: number;
        first_workout_date?: string;
        last_workout_date?: string;
        has_history: boolean;
        consistency: string;
        exercise_progress: Array<{
            exercise_type_id: number;
            name: string;
            sessions_count: number;
            recent_completion_rate: number;
            trend: string;
            recommendation_signal: string;
        }>;
    };
    measurement_summary: {
        has_measurements: boolean;
        last_date?: string;
        last_weight_kg?: number;
        weight_change_kg?: number;
        loaded_measurements: number;
    };
    provided_exercise_catalog: Array<{id: number; name: string; group_code: string; group_name: string; units: string}>;
    available_group_codes: string[];
    compatibility_notes: string[];
}


interface AIApplyProgramRequest {
    program: AIGeneratedProgram;
    warnings?: string[];
    validation_notes?: string[];
    activate: boolean;
}

interface AIApplyProgramResponse {
    program_id: number;
    name: string;
    days_count: number;
    rules_count: number;
}

interface AIGeneratedProgram {
    name: string;
    days: AIGeneratedProgramDay[];
}

interface AIGeneratedProgramDay {
    name: string;
    focus?: string[];
    exercises: AIGeneratedProgramExercise[];
}

interface AIGeneratedProgramExercise {
    exercise_type_id: number;
    sets: AIGeneratedProgramSet[];
    rest_in_seconds?: number;
    reason: string;
    progression_rule: string;
}

interface AIGeneratedProgramSet {
    reps?: number;
    weight?: number;
    minutes?: number;
    meters?: number;
}
