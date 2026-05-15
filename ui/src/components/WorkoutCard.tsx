interface WorkoutCardProps {
    w: Workout;
    idx: number;
}

export default function WorkoutCard({w, idx}: WorkoutCardProps) {
    return (
        <div style={{padding: `var(--card-padding)`, borderRadius: `var(--radius)`, boxShadow: 'var(--shadow-sm)'}}>
            <h2 style={{margin: 0}}>{idx}.{w.name}</h2>
            <div style={{paddingTop: 4, fontWeight: 600, opacity: 0.6}}>{w.started_at}</div>
            <div style={{paddingTop: 4, fontWeight: 600, opacity: 0.6}}>
                {w.status}
                {!w.has_valid_cardio_data && w.duration && <> ⏱️{w.duration}</>}
                {w.has_valid_cardio_data && w.cardio_distance > 0 && <> 🏃🏼{w.cardio_distance / 1000}км за ⏱️{w.cardio_time}мин</>}
            </div>
        </div>
    );
}
