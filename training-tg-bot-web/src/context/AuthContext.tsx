import React, {createContext, useContext, useEffect, useRef, useState} from 'react';

type AuthContextType = {
    user: User | null;
    loading: boolean;
    logout: () => Promise<void>;
    widgetRef: React.RefObject<HTMLDivElement>;
};

const AuthContext = createContext<AuthContextType>(null as any);

export const useAuth = () => useContext(AuthContext);

export const AuthProvider: React.FC<{children: React.ReactNode}> = ({children}) => {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);
    const widgetRef = useRef<HTMLDivElement>(null);

    // ---------- loadMe ----------
    const loadMe = async () => {
        const res = await fetch('/api/me', {credentials: 'include'});
        if (res.status !== 200) return setUser(null);
        setUser(await res.json());
    };

    // ---------- logout ----------
    const logout = async () => {
        await fetch('/api/logout', {
            method: 'POST',
            credentials: 'include',
        });
        setUser(null);
    };

    // ---------- init ----------
    useEffect(() => {
        loadMe().finally(() => setLoading(false));
    }, []);

    // ---------- telegram login ----------
    useEffect(() => {
        const handleMessage = async (event: MessageEvent) => {
            if (event.origin !== 'https://oauth.telegram.org') return;

            let data = event.data;
            if (typeof data === 'string') {
                try { data = JSON.parse(data); } catch { return; }
            }
            if (data?.event !== 'auth_user') return;

            await fetch('/api/login', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                credentials: 'include',
                body: JSON.stringify(data.auth_data),
            });

            await loadMe();

            if (widgetRef.current) widgetRef.current.innerHTML = '';
        };

        window.addEventListener('message', handleMessage);
        return () => window.removeEventListener('message', handleMessage);
    }, []);

    // ---------- telegram widget ----------
    useEffect(() => {
        if (user || !widgetRef.current) return;

        widgetRef.current.innerHTML = '';

        const script = document.createElement('script');
        script.src = 'https://telegram.org/js/telegram-widget.js?15';
        script.async = true;
        script.setAttribute('data-telegram-login', 'fitness_gym_buddy_dev_bot');
        script.setAttribute('data-size', 'large');
        script.setAttribute('data-userpic', 'true');
        script.setAttribute('data-request-access', 'write');

        widgetRef.current.appendChild(script);
    }, [user, widgetRef]);

    return (
        <AuthContext.Provider value={{user, loading, logout, widgetRef}}>
            {children}
        </AuthContext.Provider>
    );
};
