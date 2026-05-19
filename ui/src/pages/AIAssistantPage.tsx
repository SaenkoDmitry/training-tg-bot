import React, {useEffect, useMemo, useState} from "react";
import Button from "../components/Button";
import {buildAIProgramPrompt, createProgramFromAI, getAIProgramContext} from "../api/ai";
import {getPrograms} from "../api/programs";
import {getExerciseGroups} from "../api/exercises";
import {Bot, Clipboard, Loader, Save, Sparkles, X} from "lucide-react";
import "../styles/AIAssistantPage.css";

// Тип для группы мышц
type ExerciseGroup = {
    code: string;
    name: string;
};

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

    const [limitationsRaw, setLimitationsRaw] = useState("");

    // --- Состояния для автокомплита фокусных групп ---
    // Храним массив объектов Group
    const [exerciseGroups, setExerciseGroups] = useState<ExerciseGroup[]>([]);
    const [focusInputValue, setFocusInputValue] = useState("");
    const [showFocusSuggestions, setShowFocusSuggestions] = useState(false);
    // -------------------------------------------------

    const selectedProgramID = request.program_id || undefined;

    useEffect(() => {
        getPrograms().then(setPrograms).catch(() => setPrograms([]));
    }, []);

    // Загрузка групп мышц
    useEffect(() => {
        getExerciseGroups()
            .then(data => {
                // API возвращает массив объектов [{code: "chest", name: "Грудь"}, ...]
                setExerciseGroups(data || []);
            })
            .catch(() => setExerciseGroups([]));
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

    // Логика фильтрации подсказок
    const filteredSuggestions = useMemo(() => {
        const search = focusInputValue.toLowerCase().trim();

        return exerciseGroups
            .filter(group => {
                // Фильтруем по совпадению названия или кода
                const matchesSearch = !search ||
                    group.name.toLowerCase().includes(search) ||
                    group.code.toLowerCase().includes(search);
                // Исключаем уже выбранные (сравниваем по name, так как в focus храним name)
                const notSelected = !request.focus.includes(group.name);
                return matchesSearch && notSelected;
            });
    }, [focusInputValue, exerciseGroups, request.focus]);

    const addFocusItem = (group: ExerciseGroup) => {
        if (!request.focus.includes(group.name)) {
            update("focus", [...request.focus, group.name]);
        }
        setFocusInputValue("");
    };

    const removeFocusItem = (item: string) => {
        setRequest(prev => ({
            ...prev,
            focus: prev.focus.filter(i => i !== item)
        }));
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
        <div className="page stack">
            <h1 style={{textAlign: "center"}}><Bot size={24}/> AI-помощник</h1>

            <div className="ai-hero">
                <div>
                    <p style={{
                        backgroundColor: "var(--color-card-alt)",
                        fontWeight: 600,
                        padding: "var(--card-padding)",
                        borderRadius: "var(--radius)"
                    }}>
                        Помогает собрать промпт для создания индивидуальной программы тренировок.
                    </p>
                    <p>
                        Использует профиль, историю тренировок, замеры, текущую программу и пожелания в компактный JSON-контекст и prompt для нейросети.
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

                    {/* --- Поле Фокусные группы с автокомплитом --- */}
                    <label className="ai-autocomplete-label">
                        Фокусные группы
                        <div className="ai-tags-input-container">
                            <div className="ai-tags-wrapper">
                                {request.focus.map((item, index) => (
                                    <span key={`${item}-${index}`} className="ai-tag">
                                        {item}
                                        <button
                                            type="button"
                                            onPointerDown={(e) => {
                                                e.preventDefault();
                                                e.stopPropagation();
                                                removeFocusItem(item);
                                            }}
                                            className="ai-tag-remove"
                                        >
                                            <X size={14} />
                                        </button>
                                    </span>
                                ))}
                                <input
                                    type="text"
                                    value={focusInputValue}
                                    onChange={(e) => {
                                        setFocusInputValue(e.target.value);
                                        setShowFocusSuggestions(true);
                                    }}
                                    onFocus={() => setShowFocusSuggestions(true)}
                                    onBlur={() => setTimeout(() => setShowFocusSuggestions(false), 200)}
                                    placeholder={request.focus.length ? "Добавить..." : "Выберите группы мышц"}
                                    className="ai-tags-input"
                                />
                            </div>
                            {/* Показываем список, если есть фокус и есть варианты (даже если строка поиска пуста) */}
                            {showFocusSuggestions && filteredSuggestions.length > 0 && (
                                <ul className="ai-suggestions-list">
                                    {filteredSuggestions.map((group) => (
                                        <li
                                            key={group.code}
                                            onMouseDown={() => addFocusItem(group)}
                                            className="ai-suggestion-item"
                                        >
                                            {group.name}
                                        </li>
                                    ))}
                                </ul>
                            )}
                        </div>
                    </label>
                    {/* ------------------------------------------- */}

                </div>
                <label>
                    Ограничения / травмы / запреты через запятую или с новой строки
                    <textarea rows={3} placeholder="не нагружать поясницу, не делать бег"
                              value={limitationsRaw}
                              onChange={(e) => {
                                  setLimitationsRaw(e.target.value);
                                  update("limitations", splitText(e.target.value));
                              }}/>
                </label>
                <label>
                    Пожелания свободным текстом
                    <textarea rows={4} placeholder="Хочу уложиться в час, подтянуть грудь и спину, без тяжелой становой" value={request.notes} onChange={(e) => update("notes", e.target.value)}/>
                </label>
            </section>

            {/* Остальной код без изменений */}
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

            <Button variant="active" onClick={buildPrompt} disabled={loading} style={{width: "100%"}}>
                {loading ? <Loader size={14}/> : <Sparkles size={14}/>} Собрать prompt
            </Button>

            {result && (
                <section className="ai-card">
                    <div className="ai-result-title">
                        <h2>4. Итоговый prompt</h2>
                    </div>
                    <Button onClick={copyPrompt} style={{width: "100%"}}><Clipboard size={14}/> Копировать все</Button>

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
                    <Button
                        variant="active"
                        onClick={applyProgram}
                        disabled={applying || !aiJson.trim()}
                        style={{width: "100%", marginTop: "var(--card-gap)"}}
                    >
                        {applying ? <Loader size={14}/> : <Save size={14}/>} Создать программу
                    </Button>
                </section>
            )}

            {toast && <div className="toast">{toast}</div>}
        </div>
    );
}
