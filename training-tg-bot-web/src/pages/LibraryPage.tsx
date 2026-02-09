import React, {useEffect, useState} from 'react';
import {useAuth} from '../context/AuthContext';
import SafeTextRenderer from "../components/SafeTextRenderer.tsx";

const LibraryPage: React.FC = () => {
    const {user} = useAuth();

    const [groups, setGroups] = useState<Group[]>([]);
    const [selectedGroup, setSelectedGroup] = useState<string | null>(null);
    const [exercises, setExercises] = useState<ExerciseType[]>([]);
    const [loading, setLoading] = useState(false);
    const [openedId, setOpenedId] = useState<number | null>(null);

    // -------- load groups --------
    useEffect(() => {
        if (!user) return;

        fetch('/api/exercise-groups', {credentials: 'include'})
            .then(res => res.json())
            .then(data => {
                setGroups(data.groups);
                if (data.groups.length) setSelectedGroup(data.groups[0].code);
            });
    }, [user]);

    // -------- load exercises --------
    useEffect(() => {
        if (!selectedGroup) return;

        setLoading(true);

        fetch(`/api/exercise-groups/${selectedGroup}`, {credentials: 'include'})
            .then(res => res.json())
            .then(data => setExercises(data.exercise_types))
            .finally(() => setLoading(false));
    }, [selectedGroup]);

    if (!user) return <h3>Войдите через Telegram 👆</h3>;

    return (
        <div>
            <h1>Библиотека упражнений</h1>

            {/* ---------- GROUP TABS ---------- */}
            <div
                style={{
                    display: 'grid',
                    gridTemplateColumns: 'repeat(auto-fill, minmax(120px, 1fr))',
                    gap: 12,          // ← больше расстояние
                    padding: '4px 2px',
                    marginBottom: 20
                }}
            >
                {groups.map(g => (
                    <button
                        key={g.code}
                        onClick={() => setSelectedGroup(g.code)}
                        style={{
                            padding: '8px 10px',
                            borderRadius: 12,
                            border: '1px solid #ddd',
                            background: selectedGroup === g.code ? '#4caf50' : '#fff',
                            color: selectedGroup === g.code ? '#fff' : '#333',
                            cursor: 'pointer',
                            fontSize: 14,
                            fontWeight: 500
                        }}
                    >
                        {g.name}
                    </button>
                ))}
            </div>

            {/* ---------- EXERCISES ---------- */}
            {loading && <p>Загрузка...</p>}

            <div style={{display: 'flex', flexDirection: 'column', gap: 12}}>
                {exercises.map((ex, index) => {
                    const isOpen = openedId === ex.ID;
                    const softBg = index % 2 === 0 ? '#fff' : '#f9f9f9'; // мягкое чередование

                    return (
                        <div
                            key={ex.ID}
                            style={{
                                border: '1px solid #eee',
                                borderRadius: 12,
                                padding: 12,
                                transition: 'all 0.2s ease',
                                boxShadow: isOpen ? '0 6px 12px rgba(0,0,0,0.15)' : '0 2px 4px rgba(0,0,0,0.05)',
                                backgroundColor: isOpen ? '#e0f8e1' : softBg,
                                cursor: 'pointer',
                            }}
                            onClick={() => setOpenedId(prev => (prev === ex.ID ? null : ex.ID))}
                            onMouseEnter={e => {
                                if (!isOpen) e.currentTarget.style.boxShadow = '0 4px 8px rgba(0,0,0,0.1)';
                                if (!isOpen) e.currentTarget.style.backgroundColor = '#f6f6f6';
                            }}
                            onMouseLeave={e => {
                                if (!isOpen) e.currentTarget.style.boxShadow = '0 2px 4px rgba(0,0,0,0.05)';
                                if (!isOpen) e.currentTarget.style.backgroundColor = softBg;
                            }}
                        >
                            <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center'}}>
                                <strong>{ex.Name}</strong>
                                {/* стрелка */}
                                <span
                                    style={{
                                        display: 'inline-block',
                                        transition: 'transform 0.3s ease',
                                        transform: isOpen ? 'rotate(90deg)' : 'rotate(0deg)',
                                    }}
                                >
            ▶
          </span>
                            </div>

                            {/* раскрытие с плавным эффектом */}
                            <div
                                style={{
                                    maxHeight: isOpen ? 1000 : 0,
                                    overflow: 'hidden',
                                    transition: 'max-height 0.3s ease',
                                    marginTop: isOpen ? 10 : 0,
                                }}
                            >
                                {ex.Description && (
                                    <p>
                                        <SafeTextRenderer html={ex.Description}/>
                                    </p>
                                )}

                                <p>
                                    <b>Отдых: </b>{ex.RestInSeconds}s · <b>Ед.:</b> {ex.Units}
                                </p>

                                {ex.Url && (
                                    <a
                                        href={ex.Url}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                        style={{
                                            display: 'inline-block',
                                            padding: '10px 16px',
                                            backgroundColor: '#ffdb4d',
                                            color: 'black',
                                            borderRadius: 8,
                                            boxShadow: '0 4px 8px rgba(0,0,0,0.15)',
                                            transition: 'all 0.2s ease',
                                            marginTop: 6,
                                        }}
                                        onMouseEnter={e => {
                                            e.currentTarget.style.backgroundColor = '#ffcf27';
                                            e.currentTarget.style.transform = 'translateY(-2px)';
                                            e.currentTarget.style.boxShadow = '0 6px 12px rgba(0,0,0,0.2)';
                                        }}
                                        onMouseLeave={e => {
                                            e.currentTarget.style.backgroundColor = '#ffdb4d';
                                            e.currentTarget.style.transform = 'translateY(0)';
                                            e.currentTarget.style.boxShadow = '0 4px 8px rgba(0,0,0,0.15)';
                                        }}
                                    >
                                        Смотреть технику (🛸 я.диск)
                                    </a>
                                )}
                            </div>
                        </div>
                    );
                })}
            </div>
        </div>
    );
};

export default LibraryPage;
