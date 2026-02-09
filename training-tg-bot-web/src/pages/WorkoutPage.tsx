import {useParams} from 'react-router-dom';
import {useEffect, useState} from 'react';

const WorkoutPage = () => {
    const {id} = useParams<{ id: string }>();
    const [data, setData] = useState<ReadWorkoutDTO | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchWorkout = async () => {
            try {
                const res = await fetch(`/api/workouts/${id}`);
                if (!res.ok) throw new Error(`Ошибка: ${res.status}`);
                const json: ReadWorkoutDTO = await res.json();
                setData(json);
            } catch (err: any) {
                setError(err.message || 'Не удалось загрузить данные');
            } finally {
                setLoading(false);
            }
        };

        fetchWorkout();
    }, [id]);

    if (loading) return <p>Загрузка...</p>;
    if (error) return <p style={{color: 'red'}}>{error}</p>;
    if (!data) return <p>Данные тренировки не найдены</p>;

    const {Progress, Stats} = data;
    const {Workout, ProgressPercent, RemainingMin, SessionStarted, CompletedExercises, TotalExercises} = Progress;

    return (
        <div style={{maxWidth: '700px', margin: '0 auto', padding: '1rem'}}>
            <h2>{Workout.WorkoutDayType?.Name || `Тренировка ${Workout.ID}`}</h2>
            <p>
                Статус: {Workout.Completed ? 'Завершена' : SessionStarted ? 'В процессе' : 'Не начата'}
            </p>
            {RemainingMin !== undefined && RemainingMin > 0 && <p>Оставшееся время: {RemainingMin} мин</p>}

            {/* Прогресс тренировки */}
            <div style={{margin: '1rem 0'}}>
                <div style={{background: '#eee', borderRadius: '8px', overflow: 'hidden', height: '20px'}}>
                    <div
                        style={{
                            width: `${ProgressPercent}%`,
                            background: '#4caf50',
                            height: '100%',
                            transition: 'width 0.3s',
                        }}
                    />
                </div>
                <p>{ProgressPercent}% выполнено</p>
            </div>

            {/* Упражнения */}
            <h3>Упражнения ({CompletedExercises}/{TotalExercises})</h3>
            <ul style={{listStyle: "none", padding: 0}}>
                {Workout.Exercises.map((ex: any) => (
                    <li
                        key={ex.ID}
                        style={{
                            border: "1px solid #ddd",
                            borderRadius: "8px",
                            padding: "0.5rem",
                            marginBottom: "0.5rem",
                        }}
                    >
                        <strong>{ex.ExerciseType.Name}</strong>
                        <ul style={{paddingLeft: "1rem"}}>
                            {ex.Sets.map((set: any) => (
                                <li key={set.ID} style={{marginBottom: "0.5rem"}}>
                                    Подход {set.Index + 1}: {set.Reps} повторений{" "}
                                    {set.Weight > 0 ? `с весом ${set.Weight} кг` : ""} —{" "}
                                    {set.Completed ? "✅" : "❌"}
                                </li>
                            ))}
                            <div
                                style={{background: "#eee", height: "8px", borderRadius: "4px", overflow: "hidden", marginTop: "2px"}}>
                                <div style={{
                                    width: `${ex.Sets.filter((set: any) => set.Completed).length / ex.Sets.length * 100}%`,
                                    height: "100%",
                                    background: "#4caf50",
                                    transition: "width 0.3s",
                                }}/>
                            </div>
                        </ul>
                    </li>
                ))}
            </ul>

            <h3>Статистика</h3>

            <div
                style={{
                    display: 'grid',
                    background: '#f9f9f9',
                    padding: '1rem',
                    borderRadius: '8px',
                }}
            >
                <div>
                    {Stats.CardioTime > 0 && <p><strong>🫀 Время кардио:</strong> {Stats.CardioTime} мин</p>}
                    {Stats.TotalWeight > 0 && <p><strong>🏋 Общий вес:</strong> {Stats.TotalWeight} кг</p>}
                </div>
            </div>
        </div>
    );
};

export default WorkoutPage;
