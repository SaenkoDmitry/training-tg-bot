import Button from "./Button";
import {useNavigate} from "react-router-dom";
import {PencilLine, Trash2} from "lucide-react";

type Props = {
    day: WorkoutDayTypeDTO;
    programId: number;
    onRename: () => Promise<void>;
    onDelete: (dayId: number) => Promise<void>;
};

export default function DayCard({day, programId, onRename, onDelete}: Props) {
    const navigate = useNavigate();

    return (
        <div
            className="card row"
            style={{
                cursor: "pointer",
                padding: "0.5rem",
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between"
            }}
        >
            <div
                onClick={() => navigate(`/programs/${programId}/days/${day.id}`)}
                style={{flex: 1}}
            >
                {day.name}
            </div>

            <Button onClick={async (e) => {
                e.stopPropagation();
                await onRename();
            }}>
                <PencilLine size={14}/>
            </Button>

            <Button
                style={{marginLeft: 'var(--card-gap)'}}
                variant="danger"
                onClick={async (e) => {
                    e.stopPropagation(); // чтобы клик по кнопке не открывал страницу дня
                    await onDelete(day.id);
                }}
            >
                <Trash2 size={14}/>
            </Button>
        </div>
    );
}
