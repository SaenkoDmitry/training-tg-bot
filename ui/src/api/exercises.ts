import { api } from "./client";

export const getExerciseGroups = () =>
    api<Group[]>("/api/exercise-groups");

export const getExerciseTypesByGroup = (group: string) =>
    api<ExerciseType[]>(`/api/exercise-groups/${group}`);
