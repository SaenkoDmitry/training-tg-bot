interface Workout {
    id: number;
    name: string;
    started_at?: string;
    duration?: string;
    status?: string;
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
    ID: number;
    Name: string;
    Url: string;
    ExerciseGroupTypeCode: string;
    RestInSeconds: number;
    Accent: string;
    Units: string;
    Description: string;
}

interface Set {
    ID: number;
    Reps: number;
    FactReps: number;
    Weight: number;
    FactWeight: number;
    Completed: boolean;
    Index: number;
}

interface Exercise {
    ID: number;
    ExerciseType: ExerciseType;
    Sets: Set[];
    Index: number;
}

interface Workout {
    ID: number;
    Exercises: Exercise[];
    StartedAt: string;
    EndedAt: string | null;
    Completed: boolean;
    WorkoutDayType: {
        Name: string;
    };
}

interface WorkoutProgress {
    Workout: Workout;
    TotalExercises: number;
    CompletedExercises: number;
    TotalSets: number;
    CompletedSets: number;
    ProgressPercent: number;
    RemainingMin: number | null;
    SessionStarted: boolean;
}

interface WorkoutStatistic {
    DayType: any;
    WorkoutDay: any;
    TotalWeight: number;
    CompletedExercises: number;
    CardioTime: number; // в минутах
    ExerciseTypesMap: Record<number, any>;
    ExerciseWeightMap: Record<number, number>;
    ExerciseTimeMap: Record<number, number>;
}

interface ReadWorkoutDTO {
    Progress: WorkoutProgress;
    Stats: WorkoutStatistic;
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
