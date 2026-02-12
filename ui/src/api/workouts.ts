import {api} from "./client.ts";

export const deleteWorkout = (id: number) =>
    api(`/api/workouts/${id}`, {method: "DELETE"});

export const getWorkouts = (offset, limit: number) =>
    api<ShowMyWorkoutsResult>(`/api/workouts?offset=${offset}&limit=${limit}`);
