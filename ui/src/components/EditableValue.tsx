import { useEffect, useState } from "react";
import "../styles/EditableValue.css";

type Props = {
    fact: number;       // фактическое значение
    planned?: number;   // запланированное
    suffix: string;
    completed?: boolean;
    onSave?: (v: number) => void;
};

export default function EditableValue({
                                          fact,
                                          planned,
                                          suffix,
                                          completed,
                                          onSave,
                                      }: Props) {
    const [localStr, setLocalStr] = useState((fact || 0) as string);

    // синхронизация локального значения при изменении fact извне
    useEffect(() => {
        setLocalStr((fact || 0) as string);
    }, [fact]);

    const handleSave = () => {
        if (onSave) {
            // если ничего не введено, подставляем план
            const valueToSend = localStr
                ? parseFloat(localStr)           // конвертируем в число
                : planned || 0;                  // если пустое, берем план или 0
            onSave(valueToSend);
        }
    };

    if (completed) {
        return (
            <div className="completed-row">
                <div className="value-text"
                     style={{
                         color: fact >= planned ? "var(--color-active)" : "var(--color-danger)"
                     }}
                >
                    {!planned || planned === fact
                        ? `${fact} ${suffix}`
                        : <span>{planned} → {fact} {suffix}</span>}
                </div>
            </div>
        );
    }

    return (
        <div className="editable-wrapper">
            <input
                type="text"
                inputMode="decimal"       // цифровая клавиатура на мобильных
                value={localStr || ""}    // хранить строку, не число
                placeholder={planned?.toString()}
                className="edit-input"
                onChange={e => {
                    let val = e.target.value;

                    // Заменяем запятую на точку
                    val = val.replace(",", ".");

                    // Разрешаем цифры, пустую строку и точку с максимум 1 цифрой после
                    if (/^\d*\.?\d?$/.test(val)) {
                        setLocalStr(val);  // сохраняем как строку
                    }
                }}
                onKeyDown={e => {
                    if (e.key === "Enter") {
                        setLocalStr(localStr);  // конвертируем в число при сохранении
                        handleSave();
                    }
                }}
                onBlur={() => {
                    setLocalStr(localStr);    // конвертируем в число при потере фокуса
                    handleSave();
                }}
            />

            <span className="input-suffix">{suffix}</span>
        </div>
    );
}
