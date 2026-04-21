interface WorkoutCardProps {
    w: Workout;
    idx: number;
}

export default function WorkoutCard({ w, idx }: WorkoutCardProps) {
    return (
        <div
            style={{
                padding: 16,
                borderRadius: 16,
                boxShadow: '0 4px 12px rgba(0,0,0,0.08)',
                transition: '0.2s',
            }}
        >

            <h2 style={{ margin: 0 }}>{idx}.{w.name}</h2>

            <p style={{ margin: '6px 0', opacity: 0.7 }}>
                {w.started_at}
            </p>

            <p
                style={{
                    margin: 0,
                    fontWeight: 600,
                    color:
                        w.status === 'finished'
                            ? 'var(--color-primary-hover)'
                            : w.status === 'in_progress'
                                ? '#f9a825'
                                : '#999',
                }}
            >
                {w.status}
            </p>
        </div>
    );
}
