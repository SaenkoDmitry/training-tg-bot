import React, { createContext, useContext, useEffect, useState } from "react";

type ThemeContextType = {
    darkMode: boolean;
    toggleTheme: () => void;
};

const ThemeContext = createContext<ThemeContextType>({
    darkMode: false,
    toggleTheme: () => {},
});

export const useTheme = () => useContext(ThemeContext);

export const ThemeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [darkMode, setDarkMode] = useState<boolean>(() => {
        const saved = localStorage.getItem("darkMode");
        if (saved !== null) return saved === "true";
        return window.matchMedia?.("(prefers-color-scheme: dark)").matches ?? false;
    });

    useEffect(() => {
        const root = document.documentElement;
        if (darkMode) {
            root.classList.add("dark-theme");
        } else {
            root.classList.remove("dark-theme");
        }
        localStorage.setItem("darkMode", darkMode.toString());
    }, [darkMode]);

    const toggleTheme = () => setDarkMode((prev) => !prev);

    return (
        <ThemeContext.Provider value={{ darkMode, toggleTheme }}>
            {children}
        </ThemeContext.Provider>
    );
};
