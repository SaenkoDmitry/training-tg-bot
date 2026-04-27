interface WorkoutCardProps {
    w: Workout;
    idx: number;
}

export default function WorkoutCard({w, idx}: WorkoutCardProps) {
    return (
        <div style={{padding: `var(--card-padding)`, borderRadius: `var(--radius)`, boxShadow: 'var(--shadow-sm)'}}>
            <h2 style={{margin: 0}}>{idx}.{w.name}</h2>
            <div style={{paddingTop: 4, fontWeight: 600, opacity: 0.6}}>{w.started_at}</div>
            <div style={{paddingTop: 4, fontWeight: 600, opacity: 0.6}}>{w.status}</div>
        </div>
    );
}
