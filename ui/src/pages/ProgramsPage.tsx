import { useEffect, useState } from "react";
import {
    getPrograms,
    createProgram,
    deleteProgram,
    chooseProgram,
    renameProgram,
} from "../api/programs";
import ProgramCard from "../components/ProgramCard";
import Button from "../components/Button";
import { useNavigate } from "react-router-dom";
import "./ProgramBase.css";

export default function ProgramsPage() {
    const [programs, setPrograms] = useState<any[]>([]);
    const navigate = useNavigate();

    const load = async () => {
        setPrograms(await getPrograms());
    };

    useEffect(() => {
        load();
    }, []);

    const handleCreate = async () => {
        const name = prompt("Название программы");
        if (!name) return;

        await createProgram(name);
        load();
    };

    const handleRename = async (id: number, oldName: string) => {
        const name = prompt("Новое название", oldName);
        if (!name) return;

        await renameProgram(id, name);
        load();
    };

    const handleActivate = async (id: number) => {
        await chooseProgram(id);
        // После успешного вызова обновляем локальный стейт
        setPrograms((prev) =>
            prev.map((p) => ({
                ...p,
                is_active: p.id === id, // только выбранная программа активна
            }))
        );
    };

    return (
        <div className="page stack">
            <h1>Программы</h1>

            <Button variant="primary" onClick={handleCreate}>
                + Новая программа
            </Button>

            {programs.map((p) => (
                <ProgramCard
                    key={p.id}
                    name={p.name}
                    active={p.is_active}
                    onOpen={() => navigate(`/programs/${p.id}`)}
                    onActivate={() => handleActivate(p.id)}
                    onRename={() => handleRename(p.id, p.name)}
                    onDelete={async () => {
                        await deleteProgram(p.id);
                        load();
                    }}
                />
            ))}

        </div>
    );
}
