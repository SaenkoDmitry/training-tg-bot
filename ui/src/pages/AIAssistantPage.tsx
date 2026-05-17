import React, {useEffect, useMemo, useState} from "react";
import Button from "../components/Button";
import {buildAIProgramPrompt, createProgramFromAI, getAIProgramContext} from "../api/ai";
import {getPrograms} from "../api/programs";
import {Bot, Clipboard, Loader, Save, Sparkles} from "lucide-react";
import "../styles/AIAssistantPage.css";

const defaultRequest: AIProgramPromptRequest = {
    mode: "create_program",
    program_id: null,
    goal: "muscle_gain",
    level: "beginner",
    days_per_week: 3,
    session_duration_min: 60,
    location: "gym",
    limitations: [],
    focus: [],
    notes: "",
};

const goals = [
    {value: "muscle_gain", label: "Набор мышц"},
    {value: "fat_loss", label: "Похудение"},
    {value: "strength", label: "Сила"},
    {value: "endurance", label: "Выносливость"},
    {value: "health", label: "Здоровье/форма"},
    {value: "return_after_break", label: "Возврат после перерыва"},
];

const levels = [
    {value: "beginner", label: "Новичок"},
    {value: "intermediate", label: "Средний"},
    {value: "advanced", label: "Опытный"},
];

const locations = [
    {value: "gym", label: "Зал"},
    {value: "home", label: "Дом"},
    {value: "outdoor", label: "Улица"},
];

export default function AIAssistantPage() {
    const [request, setRequest] = useState<AIProgramPromptRequest>(defaultRequest);
    const [programs, setPrograms] = useState<ProgramDTO[]>([]);
    const [previewContext, setPreviewContext] = useState<AIProgramContext | null>(null);
    const [result, setResult] = useState<AIProgramPromptResponse | null>(null);
    const [loading, setLoading] = useState(false);
    const [applying, setApplying] = useState(false);
    const [aiJson, setAIJson] = useState("");
    const [activateCreated, setActivateCreated] = useState(true);
    const [toast, setToast] = useState<string | null>(null);

    const [focusRaw, setFocusRaw] = useState("");
    const [limitationsRaw, setLimitationsRaw] = useState("");

    const selectedProgramID = request.program_id || undefined;

    useEffect(() => {
        getPrograms().then(setPrograms).catch(() => setPrograms([]));
    }, []);

    useEffect(() => {
        getAIProgramContext(selectedProgramID)
            .then(setPreviewContext)
            .catch(() => setPreviewContext(null));
    }, [selectedProgramID]);

    const selectedProgram = useMemo(
        () => programs.find((program) => program.id === request.program_id),
        [programs, request.program_id]
    );

    const showToast = (text: string) => {
        setToast(text);
        setTimeout(() => setToast(null), 2500);
    };

    const update = <K extends keyof AIProgramPromptRequest>(key: K, value: AIProgramPromptRequest[K]) => {
        setRequest(prev => ({...prev, [key]: value}));
    };

    const splitText = (value: string) => value
        .split(/[\n,]/)
        .map(item => item.trim())
        .filter(Boolean);

    const buildPrompt = async () => {
        try {
            setLoading(true);
            setResult(null);
            const payload = {
                ...request,
                program_id: request.mode === "improve_existing_program" ? request.program_id : null,
            };
            const response = await buildAIProgramPrompt(payload);
            setResult(response);
            setPreviewContext(response.context);
            showToast("✅ Prompt собран");
        } catch (e) {
            showToast(e instanceof Error ? `❌ ${e.message}` : "❌ Ошибка сборки prompt");
        } finally {
            setLoading(false);
        }
    };

    const copyPrompt = async () => {
        if (!result) return;
        await navigator.clipboard.writeText(
            `system prompt:\n${result.system_prompt}\n\n` +
            `user prompt:\n${result.user_prompt}\n\n` +
            `output schema:\n${JSON.stringify(result.output_schema, null, 2)}`
        );
        showToast("📋 Prompt скопирован");
    };

    const applyProgram = async () => {
        try {
            setApplying(true);
            const parsed = JSON.parse(aiJson);
            const program = parsed.program || parsed;
            const response = await createProgramFromAI({
                program,
                warnings: parsed.warnings || [],
                validation_notes: parsed.validation_notes || [],
                activate: activateCreated,
            });
            showToast(`✅ Создана программа «${response.name}» (${response.days_count} дн.)`);
        } catch (e) {
            showToast(e instanceof Error ? `❌ ${e.message}` : "❌ Не удалось создать программу");
        } finally {
            setApplying(false);
        }
    };

    return (
        <div className="page stack ai-page">
            <div className="ai-hero">
                <div>
                    <h1><Bot/> AI-помощник</h1>
                    <p>
                        Собирает профиль, историю тренировок, замеры, текущую программу и пожелания в компактный JSON-контекст и prompt для нейросети.
                    </p>
                </div>
            </div>

            <section className="ai-card">
                <h2>1. Что нужно сделать?</h2>
                <div className="ai-grid two">
                    <label>
                        Режим
                        <select value={request.mode} onChange={(e) => update("mode", e.target.value as AIProgramPromptRequest["mode"])}>
                            <option value="create_program">Создать новую программу</option>
                            <option value="improve_existing_program">Скорректировать текущую</option>
                        </select>
                    </label>

                    <label>
                        Программа для коррекции
                        <select
                            value={request.program_id || ""}
                            disabled={request.mode !== "improve_existing_program"}
                            onChange={(e) => update("program_id", e.target.value ? Number(e.target.value) : null)}
                        >
                            <option value="">Активная или отсутствует</option>
                            {programs.map((program) => (
                                <option key={program.id} value={program.id}>{program.name}{program.is_active ? " · активная" : ""}</option>
                            ))}
                        </select>
                    </label>
                </div>
            </section>

            <section className="ai-card">
                <h2>2. Цель и ограничения</h2>
                <div className="ai-grid two">
                    <label>
                        Цель
                        <select value={request.goal} onChange={(e) => update("goal", e.target.value)}>
                            {goals.map(goal => <option key={goal.value} value={goal.value}>{goal.label}</option>)}
                        </select>
                    </label>
                    <label>
                        Уровень
                        <select value={request.level} onChange={(e) => update("level", e.target.value)}>
                            {levels.map(level => <option key={level.value} value={level.value}>{level.label}</option>)}
                        </select>
                    </label>
                    <label>
                        Дней в неделю
                        <input type="number" min={1} max={7} value={request.days_per_week} onChange={(e) => update("days_per_week", Number(e.target.value))}/>
                    </label>
                    <label>
                        Длительность тренировки, минут
                        <input type="number" min={20} max={180} step={5} value={request.session_duration_min} onChange={(e) => update("session_duration_min", Number(e.target.value))}/>
                    </label>
                    <label>
                        Где тренируешься
                        <select value={request.location} onChange={(e) => update("location", e.target.value)}>
                            {locations.map(location => <option key={location.value} value={location.value}>{location.label}</option>)}
                        </select>
                    </label>
                    <label>
                        Фокусные группы через запятую
                        <input value={focusRaw}
                               onChange={(e) => {
                                   setFocusRaw(e.target.value);
                                   update("focus", splitText(e.target.value));
                               }} placeholder="chest, back, legs, biceps, triceps, deltas, press, cardio, buttocks"
                        />
                    </label>
                </div>
                <label>
                    Ограничения / травмы / запреты через запятую или с новой строки
                    <textarea rows={3} placeholder="не нагружать поясницу, не делать бег"
                              value={limitationsRaw}
                              onChange={(e) => {
                                  setLimitationsRaw(e.target.value);
                                  update("focus", splitText(e.target.value));
                              }}/>
                </label>
                <label>
                    Пожелания свободным текстом
                    <textarea rows={4} placeholder="Хочу уложиться в час, подтянуть грудь и спину, без тяжелой становой" value={request.notes} onChange={(e) => update("notes", e.target.value)}/>
                </label>
                <Button variant="active" onClick={buildPrompt} disabled={loading}>
                    {loading ? <Loader size={14}/> : <Sparkles size={14}/>} Собрать prompt
                </Button>
            </section>

            <section className="ai-card ai-summary">
                <h2>3. Контекст, который уйдет в prompt</h2>
                <div className="ai-metrics">
                    <div><b>{previewContext?.training_summary.completed_workouts ?? 0}</b><span>тренировок</span></div>
                    <div><b>{previewContext?.training_summary.exercise_progress.length ?? 0}</b><span>упражнений в истории</span></div>
                    <div><b>{previewContext?.provided_exercise_catalog.length ?? 0}</b><span>упражнений в каталоге</span></div>
                    <div><b>{previewContext?.measurement_summary.has_measurements ? "есть" : "нет"}</b><span>замеры</span></div>
                </div>
                <p>
                    Текущая программа: <b>{selectedProgram?.name || previewContext?.current_program?.name || "не выбрана / отсутствует"}</b>
                </p>
            </section>

            {result && (
                <section className="ai-card">
                    <div className="ai-result-title">
                        <h2>4. Итоговый prompt</h2>
                        <Button onClick={copyPrompt}><Clipboard size={14}/> Копировать</Button>
                    </div>
                    <h3>System prompt</h3>
                    <pre>{result.system_prompt}</pre>
                    <h3>User prompt</h3>
                    <pre>{result.user_prompt}</pre>
                    <h3>JSON schema для ответа</h3>
                    <pre>{JSON.stringify(result.output_schema, null, 2)}</pre>
                </section>
            )}

            {result && (
                <section className="ai-card">
                    <h2>5. Создать программу из JSON ответа AI</h2>
                    <p className="ai-help">
                        Вставьте сюда полный JSON-ответ нейросети. Можно вставить как весь объект с полем <code>program</code>, так и сам объект программы.
                    </p>
                    <label className="ai-checkbox">
                        <input
                            type="checkbox"
                            checked={activateCreated}
                            onChange={(e) => setActivateCreated(e.target.checked)}
                        />
                        Сделать созданную программу активной
                    </label>
                    <textarea
                        rows={10}
                        placeholder={`{
  "summary": "...",
  "program": {
    "name": "AI программа",
    "days": []
  },
  "warnings": [],
  "validation_notes": []
}`}
                        value={aiJson}
                        onChange={(e) => setAIJson(e.target.value)}
                    />
                    <Button variant="active" onClick={applyProgram} disabled={applying || !aiJson.trim()}>
                        {applying ? <Loader size={14}/> : <Save size={14}/>} Создать программу
                    </Button>
                </section>
            )}

            {toast && <div className="toast">{toast}</div>}
        </div>
    );
}
