import { useEffect, useState, useCallback } from "react";
import { getUserIcon, changeUserIcon } from "../api/profile.ts";
import { ICONS } from "../components/IconPicker";
import type { IconName } from "../components/IconPicker";

const STORAGE_KEY = "user_icon";

const isValidIcon = (value: any): value is IconName => {
    return typeof value === "string" && value in ICONS;
};

export const useUserIcon = () => {
    // 🔥 сразу читаем из localStorage
    const [icon, setIcon] = useState<IconName>(() => {
        const cached = localStorage.getItem(STORAGE_KEY);
        if (isValidIcon(cached)) return cached;
        return "Smile";
    });

    // --- Фоновая синхронизация с сервером ---
    useEffect(() => {
        const syncWithServer = async () => {
            try {
                const serverIcon = await getUserIcon();

                if (isValidIcon(serverIcon) && serverIcon !== icon) {
                    setIcon(serverIcon);
                    localStorage.setItem(STORAGE_KEY, serverIcon);
                }
            } catch {
                // тихо игнорируем
            }
        };

        syncWithServer();
    }, []);

    // --- Обновление ---
    const updateIcon = useCallback(async (name: IconName) => {
        const prev = icon;

        // optimistic
        setIcon(name);
        localStorage.setItem(STORAGE_KEY, name);

        try {
            await changeUserIcon(name);
        } catch {
            // rollback
            setIcon(prev);
            localStorage.setItem(STORAGE_KEY, prev);
            throw new Error("Failed to update icon");
        }
    }, [icon]);

    return {
        icon,
        updateIcon,
    };
};
