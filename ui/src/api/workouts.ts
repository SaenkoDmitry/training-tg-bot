import {api} from "./client.ts";

export const deleteWorkout = (id: number) =>
    api(`/api/workouts/${id}`, {method: "DELETE"});

export const getWorkouts = (offset, limit: number) =>
    api<ShowMyWorkoutsResult>(`/api/workouts?offset=${offset}&limit=${limit}`);

export const startWorkout = (dayTypeID: number) =>
    api<StartWorkoutDTO>("/api/workouts/start", {
        method: "POST",
        body: JSON.stringify({
            day_type_id: dayTypeID,
        }),
    });

export const finishWorkout = (workoutID: number) =>
    api<StartWorkoutDTO>(`/api/workouts/${workoutID}/finish`, {
        method: "POST",
    });

export const getWorkout = (id: number) =>
    api<ReadWorkoutDTO>(`/api/workouts/${id}`);

export const createShare = async (workoutId: number): Promise<{ token: string; share_url: string }> => {
    return api<ShareDTO>(`/api/workouts/${workoutId}/share`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
    });
};

export const getPublicWorkout = async (token: string): Promise<ReadWorkoutDTO> => {
    return api<ReadWorkoutDTO>(`/api/public/workouts/${token}`, {
        method: 'GET',
    });
};
