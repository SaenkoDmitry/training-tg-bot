import React, { useEffect, useRef } from "react";
import { useAuth } from "../context/AuthContext";

declare global {
    interface Window {
        onTelegramAuth: (user: any) => void;
    }
}

const TelegramLoginWidget: React.FC = () => {
    const { refreshUser } = useAuth();
    const widgetRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (!widgetRef.current) return;

        widgetRef.current.innerHTML = "";

        const botUsername =
            process.env.NODE_ENV === "development"
                ? "fitness_gym_buddy_dev_bot"
                : "form_journey_bot";

        window.onTelegramAuth = async (user: any) => {
            const res = await fetch("/api/telegram/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(user),
            });

            if (!res.ok) {
                alert("Login failed");
                return;
            }

            const data = await res.json();
            localStorage.setItem("token", data.token);

            await refreshUser();
        };

        const script = document.createElement("script");
        script.src = "https://telegram.org/js/telegram-widget.js?22";
        script.async = true;
        script.setAttribute("data-telegram-login", botUsername);
        script.setAttribute("data-size", "large");
        script.setAttribute("data-userpic", "true");
        script.setAttribute("data-onauth", "onTelegramAuth(user)");

        widgetRef.current.appendChild(script);
    }, []);

    return <div ref={widgetRef} />;
};

export default TelegramLoginWidget;
