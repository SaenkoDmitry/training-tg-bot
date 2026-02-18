import React, {useEffect, useState} from "react";
import {useAuth} from "../context/AuthContext";
import Button from "../components/Button";
import {subscribePush, unsubscribePush} from "../api/subscribePush.ts";
import Toast from "../components/Toast.tsx";

const VAPID_PUBLIC_KEY = 'BK0VOgS6oooJu5aKXkg0Amn6zVTWqEjjHjlxFJE4lMygZ_Wyp_D1LCVR3LkCEiOF4hHsCRDCNEa-TMlkR22LEms';

const ProfilePage: React.FC = () => {
    const {user, logout, loading} = useAuth();
    const [toast, setToast] = useState<string | null>(null);
    const [notificationsEnabled, setNotificationsEnabled] = useState(false);
    const [checking, setChecking] = useState(true);

    // Проверяем есть ли активная подписка
    useEffect(() => {
        const checkSubscription = async () => {
            if (!("serviceWorker" in navigator)) {
                setChecking(false);
                return;
            }

            const registration = await navigator.serviceWorker.ready;
            const subscription = await registration.pushManager.getSubscription();
            setNotificationsEnabled(!!subscription);
            setChecking(false);
        };

        checkSubscription();
    }, []);

    const toggleNotifications = async () => {
        if (!("serviceWorker" in navigator)) return;

        const registration = await navigator.serviceWorker.ready;
        const subscription = await registration.pushManager.getSubscription();

        // Если подписка есть → отключаем
        if (subscription) {
            await subscription.unsubscribe();
            await unsubscribePush();
            setNotificationsEnabled(false);
            setToast("Уведомления выключены ❌");
            return;
        }

        // Если нет → включаем
        const permission = await Notification.requestPermission();
        if (permission === "granted") {
            await subscribePush(VAPID_PUBLIC_KEY);
            setNotificationsEnabled(true);
            setToast("Уведомления включены ✅");
        } else {
            setToast("Разрешение на уведомления не выдано");
        }
    };

    return (
        <div
            style={{
                maxWidth: 420,
                margin: '0 auto',
                padding: '1rem',
                display: 'flex',
                flexDirection: 'column',
                gap: 20,
            }}
        >
            {/* ---------------- NOT LOGGED IN ---------------- */}
            {!loading && !user && (
                <div
                    style={{
                        background: '#fff',
                        borderRadius: 20,
                        padding: '2rem 1.5rem',
                        boxShadow: '0 6px 20px rgba(0,0,0,0.06)',
                        textAlign: 'center',
                    }}
                >
                    <div style={{fontSize: 42, marginBottom: 12}}>🔐</div>

                    <div
                        style={{
                            fontSize: 16,
                            fontWeight: 600,
                            marginBottom: 16,
                        }}
                    >
                        Войдите в аккаунт
                    </div>

                    <Button
                        variant={"primary"}
                        onClick={() => {
                            const origin = window.location.origin;

                            window.location.href =
                                `/api/telegram/login?origin=${encodeURIComponent(origin)}`;
                        }}
                    >
                        Войти через Telegram
                    </Button>
                </div>
            )}

            {/* ---------------- LOGGED IN ---------------- */}
            {user && (
                <>
                    <div
                        style={{
                            background: '#fff',
                            borderRadius: 20,
                            padding: '1.5rem',
                            boxShadow: '0 6px 20px rgba(0,0,0,0.06)',
                            textAlign: 'center',
                        }}
                    >
                        <div style={{fontSize: 42, marginBottom: 8}}>👤</div>

                        <div style={{fontSize: 18, fontWeight: 600}}>
                            {user.first_name}
                        </div>

                        {user.username && (
                            <div style={{opacity: 0.6, fontSize: 14}}>
                                @{user.username}
                            </div>
                        )}
                    </div>

                    <Button
                        variant="danger"
                        onClick={logout}
                        style={{
                            width: '100%',
                            height: 48,
                            fontSize: 16,
                            borderRadius: 14,
                        }}
                    >
                        Выйти из аккаунта
                    </Button>

                    {!checking && (
                        <Button
                            variant={notificationsEnabled ? "ghost" : "active"}
                            onClick={toggleNotifications}
                        >
                            {notificationsEnabled
                                ? "Выключить уведомления"
                                : "Включить уведомления"}
                        </Button>
                    )}

                    {toast && <Toast message={toast} onClose={() => setToast(null)}/>}
                </>
            )}
        </div>
    );
};

export default ProfilePage;
