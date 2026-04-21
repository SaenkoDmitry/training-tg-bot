import { useEffect } from "react";
import { useRestTimer } from "../context/RestTimerContext";
import Button from "./Button";
import "../styles/RestTimer.css";
import { Pause, Play, RotateCcw } from "lucide-react";
import { toast } from "react-hot-toast";
import {startTimer} from "../api/timers.ts";

type Props = {
    seconds: number;
    autoStartTrigger?: number;
    workoutID?: number;
};

export default function RestTimer({ seconds, autoStartTrigger, workoutID }: Props) {
    const {
        remaining,
        running,
        start: localStart,
        pause,
        reset,
        seconds: totalSeconds
    } = useRestTimer();

    // 🔥 функция старта с API
    const start = async (secs: number) => {
        if (secs <= 0) return;

        // серверный таймер и push
        if (!workoutID) return;
        try {
            await startTimer(workoutID, secs).then((resp: TimerDTO) => {
                console.log("started timer:", resp.id)
                if (resp?.id != null) {
                    localStorage.setItem("currentTimerID", resp.id.toString());
                }
                localStart(secs, resp.id); // локальный таймер
            });
        } catch (err) {
            console.error("Failed to start server timer", err);
            toast.error("Не удалось зарегистрировать таймер на сервере");
        }
    };

    // 🔥 автостарт после завершения подхода
    useEffect(() => {
        if (!autoStartTrigger || !workoutID) return;
        localStorage.setItem("floatingTimerWorkoutID", workoutID.toString());
        start(seconds);
    }, [autoStartTrigger, workoutID]);

    // 🔹 уведомление и вибрация при завершении таймера (только локально)
    useEffect(() => {
        if (remaining === 0 && running) {
            // Вибрация
            navigator.vibrate?.([300, 150, 300]);

            // Toast уведомление
            toast.success("Таймер завершён!");

            // очищаем ID
            localStorage.removeItem("floatingTimerWorkoutID");
        }
    }, [remaining, running]);

    const format = (t: number) => {
        const m = Math.floor(t / 60);
        const s = t % 60;
        return `${m}:${s.toString().padStart(2, "0")}`;
    };

    const progress = totalSeconds > 0 ? 1 - remaining / totalSeconds : 0;
    const radius = 28;
    const circumference = 2 * Math.PI * radius;

    return (
        <div className={`rest-timer ${running ? "active" : ""}`}>
            <div className="timer-inner">
                <div className="circle">
                    <svg width="70" height="70">
                        <circle
                            className="bg"
                            strokeWidth="8"
                            r={radius}
                            cx="35"
                            cy="35"
                        />
                        <circle
                            className="progress"
                            strokeWidth="8"
                            r={radius}
                            cx="35"
                            cy="35"
                            strokeDasharray={circumference}
                            strokeDashoffset={circumference * (1 - progress)}
                        />
                    </svg>
                    <div className="time">{format(remaining)}</div>
                </div>

                <div className="actions">
                    <Button
                        variant={running ? "primary" : "active"}
                        onClick={() => {
                            if (running) {
                                pause();
                            } else if (remaining > 0) {
                                start(remaining);
                            } else {
                                localStorage.setItem("floatingTimerWorkoutID", workoutID?.toString() ?? "");
                                start(seconds);
                            }
                        }}
                    >
                        {running ? <Pause size={14} /> : <Play size={14} />}{" "}
                        {running ? "Пауза" : "Старт"}
                    </Button>

                    <Button variant="ghost" onClick={reset}>
                        <RotateCcw size={14} /> Сброс
                    </Button>
                </div>
            </div>
        </div>
    );
}
