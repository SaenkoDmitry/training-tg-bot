import EditableValue from "./EditableValue";
import Button from "./Button.tsx";

type Props = {
    set: FormattedSet;
    index: number;
    onDelete: () => void;
    onComplete: () => void;
    onChange: (id: number, reps: number, weight: number, minutes: number, meters: number) => void;
};

export default function SetRow({ set, index, onDelete, onComplete, onChange }: Props) {
    return (
        <div className={`set-row ${set.completed ? "done" : ""}`}>
            <div className="set-index">{index + 1}</div>

            {set.reps > 0 && (
                <EditableValue
                    fact={set.fact_reps}
                    planned={set.reps}
                    suffix="повт."
                    completed={set.completed}
                    onSave={(v) => onChange(set.id, v, set.fact_weight, set.fact_minutes, set.fact_meters)}
                />
            )}

            {set.weight > 0 && (
                <EditableValue
                    fact={set.fact_weight}
                    planned={set.weight}
                    suffix="кг"
                    completed={set.completed}
                    onSave={(v) => onChange(set.id, set.fact_reps, v, set.fact_minutes, set.fact_meters)}
                />
            )}

            {set.minutes > 0 && (
                <EditableValue
                    fact={set.fact_minutes}
                    planned={set.minutes}
                    suffix="мин"
                    completed={set.completed}
                    onSave={(v) => onChange(set.id, set.fact_reps, set.fact_weight, v, set.fact_meters)}
                />
            )}

            {set.meters > 0 && (
                <EditableValue
                    fact={set.fact_meters}
                    planned={set.meters}
                    suffix="м"
                    completed={set.completed}
                    onSave={(v) => onChange(set.id, set.fact_reps, set.fact_weight, set.fact_minutes, v)}
                />
            )}

            <div className="set-actions">
                <Button variant={"active"} onClick={onComplete}>✓</Button>
                <Button variant={"danger"} onClick={onDelete}>✕</Button>
            </div>
        </div>
    );
}
