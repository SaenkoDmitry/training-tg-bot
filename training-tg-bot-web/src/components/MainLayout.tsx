import React, {useState} from 'react';
import {Link, useLocation} from 'react-router-dom';
import {useAuth} from '../context/AuthContext';
import TelegramLogin from "../pages/TelegramLogin.tsx";

const tabs = [
    {name: 'Мои тренировки', path: '/'},
    {name: 'Статистика', path: '/stats'},
    {name: 'Мои программы', path: '/programs'},
    {name: 'Замеры', path: '/measurements'},
    {name: 'Библиотека упражнений', path: '/library'},
];

const MainLayout: React.FC<{children: React.ReactNode}> = ({children}) => {
    const location = useLocation();
    const [menuOpen, setMenuOpen] = useState(false);
    const {user, logout, widgetRef, loading} = useAuth();

    const isMobile = window.innerWidth <= 768;

    return (
        <div style={{minHeight: '100vh', fontFamily: 'Arial, sans-serif'}}>

            {/* ---------------- NAVBAR ---------------- */}
            <nav
                style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    padding: '0.5rem 1rem',
                    borderBottom: '1px solid #ddd',
                    background: '#fafafa',
                    position: 'sticky',
                    top: 0,
                    zIndex: 10,
                    gap: 12
                }}
            >
                {/* ---------- Tabs (desktop) ---------- */}
                <div
                    style={{
                        display: 'flex',
                        gap: '0.5rem',
                        flex: 1,
                        overflowX: 'auto',
                    }}
                >
                    {tabs.map((tab) => (
                        <Link
                            key={tab.path}
                            to={tab.path}
                            style={{
                                padding: '0.5rem 1rem',
                                borderBottom:
                                    location.pathname === tab.path ? '3px solid #4caf50' : 'none',
                                color: location.pathname === tab.path ? '#4caf50' : '#333',
                                fontWeight: location.pathname === tab.path ? 'bold' : 'normal',
                                textDecoration: 'none',
                                whiteSpace: 'nowrap',
                                flexShrink: 0,
                                display: !isMobile ? 'block' : 'none',
                            }}
                        >
                            {tab.name}
                        </Link>
                    ))}
                </div>

                {/* ---------- Right side ---------- */}
                <div style={{display: 'flex', alignItems: 'center', gap: 12}}>

                    {/* Telegram login */}
                    {!loading && <TelegramLogin />}

                    {/* User info */}
                    {user && (
                        <>
                            <span>Привет, {user.first_name} 👋</span>
                            <button
                                onClick={logout}
                                style={{
                                    background: '#eee',
                                    border: '1px solid #ccc',
                                    borderRadius: 6,
                                    padding: '0.4rem 0.8rem',
                                    cursor: 'pointer',
                                }}
                            >
                                Выйти
                            </button>
                        </>
                    )}

                    {/* Start workout button (только на странице тренировок) */}
                    {user && location.pathname === '/' && (
                        <button
                            style={{
                                background: '#4caf50',
                                color: '#fff',
                                border: 'none',
                                borderRadius: 8,
                                padding: '0.5rem 1rem',
                                cursor: 'pointer',
                            }}
                            onClick={() => alert('Начало новой тренировки!')}
                        >
                            Начать тренировку
                        </button>
                    )}

                    {/* Burger */}
                    {isMobile && (
                        <div onClick={() => setMenuOpen(!menuOpen)} style={{cursor: 'pointer'}}>
                            <div style={{width: 25, height: 3, background: '#333', margin: '4px 0'}}/>
                            <div style={{width: 25, height: 3, background: '#333', margin: '4px 0'}}/>
                            <div style={{width: 25, height: 3, background: '#333', margin: '4px 0'}}/>
                        </div>
                    )}
                </div>
            </nav>

            {/* ---------------- MOBILE MENU ---------------- */}
            {menuOpen && (
                <div
                    style={{
                        position: 'fixed',
                        inset: 0,
                        background: 'rgba(0,0,0,0.7)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        zIndex: 20,
                    }}
                    onClick={() => setMenuOpen(false)}
                >
                    <div
                        style={{
                            background: '#fff',
                            padding: '2rem',
                            borderRadius: 8,
                            display: 'flex',
                            flexDirection: 'column',
                            gap: '1rem',
                            minWidth: 200,
                        }}
                        onClick={(e) => e.stopPropagation()}
                    >
                        {tabs.map((tab) => (
                            <Link
                                key={tab.path}
                                to={tab.path}
                                onClick={() => setMenuOpen(false)}
                                style={{
                                    textAlign: 'center',
                                    textDecoration: 'none',
                                    color: '#333',
                                    fontWeight:
                                        location.pathname === tab.path ? 'bold' : 'normal',
                                }}
                            >
                                {tab.name}
                            </Link>
                        ))}
                    </div>
                </div>
            )}

            {/* ---------------- CONTENT ---------------- */}
            <div style={{padding: '1rem'}}>
                {children}
            </div>
        </div>
    );
};

export default MainLayout;
