import React, {createContext, useCallback, useContext, useEffect, useState,} from "react";
import {changeUserIcon, getUserIcon} from "../api/users";
import type {IconName} from "../components/IconPicker";
import {ICONS} from "../components/IconPicker";

interface ContextType {
    icon: IconName;
    loading: boolean;
    updateIcon: (name: IconName) => Promise<void>;
}

const UserIconContext = createContext<ContextType | null>(null);

export const UserIconProvider: React.FC<{ children: React.ReactNode }> = ({children}) => {
    const [icon, setIcon] = useState<IconName>("Smile");
    const [loading, setLoading] = useState(true);
    const [loaded, setLoaded] = useState(false); // 🔥 чтобы не грузить повторно

    const loadIcon = useCallback(async () => {
        if (loaded) return; // уже загружали

        try {
            const name = await getUserIcon();

            if (ICONS[name as IconName]) {
                setIcon(name as IconName);
            }
        } catch {
            setIcon("Smile");
        } finally {
            setLoading(false);
            setLoaded(true);
        }
    }, [loaded]);

    useEffect(() => {
        loadIcon();
    }, [loadIcon]);

    const updateIcon = async (name: IconName) => {
        const prev = icon;
        setIcon(name); // optimistic

        try {
            await changeUserIcon(name);
        } catch {
            setIcon(prev);
            throw new Error("Failed to update icon");
        }
    };

    return (
        <UserIconContext.Provider value={{icon, loading, updateIcon}}>
            {children}
        </UserIconContext.Provider>
    );
};

export const useUserIcon = () => {
    const ctx = useContext(UserIconContext);
    if (!ctx) {
        throw new Error("useUserIcon must be used inside UserIconProvider");
    }
    return ctx;
};
