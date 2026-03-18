import {Container, Loader} from "lucide-react";
import React, {useEffect, useState} from "react";
import Toast from "../components/Toast.tsx";
import {getExerciseGroups, getExerciseTypesByGroup} from "../api/exercises.ts";
import {useNavigate, useParams} from "react-router-dom";
import Button from "../components/Button.tsx";

const StatsPageGroupExercise: React.FC = () => {
    const [toast, setToast] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();
    const {groupCode, exerciseID} = useParams();
    const [groupsMap, setGroupsMap] = useState<Record<string, Group>>({});
    const [exercisesMap, setExercisesMap] = useState<Record<number, ExerciseType>>({});

    useEffect(() => {
        if (!groupCode) return;

        const fetchGroupData = async () => {
            try {
                setLoading(true);

                const [exerciseTypes, groups]: [ExerciseType[], Group[]] = await Promise.all([
                    getExerciseTypesByGroup(groupCode),
                    getExerciseGroups()
                ]);

                const groupsMap = groups.reduce<Record<string, Group>>((acc, group) => {
                    acc[group.code] = group;
                    return acc;
                }, {});

                const exercisesMap = exerciseTypes.reduce<Record<number, ExerciseType>>((acc, ex) => {
                    acc[ex.id] = ex;
                    return acc;
                }, {});

                setExercisesMap(exercisesMap);
                setGroupsMap(groupsMap);

            } catch (err: any) {
                setError(err.message || 'Не удалось загрузить данные');
            } finally {
                setLoading(false);
            }
        }

        fetchGroupData();
    }, [groupCode]);

    if (loading) return <Loader/>;
    if (error) return <p style={{color: 'red'}}>{error}</p>;

    return (
        <div>

            <div className={"page stack"}>

                {groupCode && <h1>Динамика: {groupsMap[groupCode].name}</h1>}
                {groupCode && <div style={{color: '--color-text-muted'}}><b>{exercisesMap[Number(exerciseID)].name}</b></div>}
                {<Button variant="ghost"><Container/> В разработке</Button>}
            </div>

            {toast && <Toast message={toast} onClose={() => setToast(null)}/>}
        </div>
    );
}

export default StatsPageGroupExercise;
