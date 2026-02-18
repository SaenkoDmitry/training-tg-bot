import React, {createContext, useContext, useEffect, useRef, useState} from 'react';

type AuthContextType = {
    user: User | null;
    loading: boolean;
    logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType>(null as any);

export const useAuth = () => useContext(AuthContext);

export const AuthProvider: React.FC<{children: React.ReactNode}> = ({children}) => {
    
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);

    const loadMe = async () => {
        const res = await fetch('/api/me', {credentials: 'include'});
        if (res.status !== 200) return setUser(null);
        setUser(await res.json());
    };

    const logout = async () => {
        await fetch('/api/logout', { method: 'POST', credentials: 'include' });
        setUser(null);
    };

    useEffect(() => {
        loadMe().finally(() => setLoading(false));
    }, []);

    return (
        <AuthContext.Provider value={{user, loading, logout}}>
            {children}
        </AuthContext.Provider>
    );
};
