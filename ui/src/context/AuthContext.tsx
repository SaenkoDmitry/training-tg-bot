import React, {createContext, useContext, useEffect, useRef, useState} from 'react';

type AuthContextType = {
    user: User | null;
    loading: boolean;
    logout: () => Promise<void>;
    refreshUser: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType>(null as any);

export const useAuth = () => useContext(AuthContext);

export const AuthProvider: React.FC<{children: React.ReactNode}> = ({children}) => {

    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);

    const refreshUser = async () => {
        setLoading(true);
        await loadMe().finally(() => setLoading(false));
    };

    const loadMe = async () => {
        const token = localStorage.getItem("token");

        if (!token) {
            setUser(null);
            return;
        }

        const res = await fetch('/api/me', {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });

        if (res.status !== 200) {
            setUser(null);
            return;
        }

        setUser(await res.json());
    };

    const logout = async () => {
        localStorage.removeItem("token");
        setUser(null);
    };

    useEffect(() => {
        loadMe().finally(() => setLoading(false));
    }, []);

    return (
        <AuthContext.Provider value={{ user, loading, logout, refreshUser }}>
            {children}
        </AuthContext.Provider>
    );
};
