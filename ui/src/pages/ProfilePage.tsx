import React, {useEffect, useState} from "react";
import {useAuth} from "../context/AuthContext";
import Button from "../components/Button";
import Toast from "../components/Toast";
import {
    Bell,
    BellOff,
    Cake,
    ChartNoAxesCombined,
    LogOut,
    Moon,
    Pencil,
    Ruler,
    Sun,
    User,
    UserLock,
    Weight
} from "lucide-react";
import type {IconName} from "../components/IconPicker";
import IconPicker, {ICONS} from "../components/IconPicker";
import {subscribePush, unsubscribePush} from "../api/subscribePush";
import {useUserIcon} from "../hooks/useUserIcons.ts";
import {useNavigate} from "react-router-dom";
import {getVapidKey} from "../api/vapid.ts";
import {getProfile, updateProfile} from "../api/profile";

const ProfilePage: React.FC = () => {
    const {user, logout, loading} = useAuth();
    const [toast, setToast] = useState<string | null>(null);
    const [notificationsEnabled, setNotificationsEnabled] = useState(false);
    const [checking, setChecking] = useState(true);
    const navigate = useNavigate();

    const [iconModalOpen, setIconModalOpen] = useState(false);
    const {icon, updateIcon} = useUserIcon();
    const CurrentIcon = ICONS[icon];

    const [profile, setProfile] = useState<UserProfile | null>(null);
    const [profileLoading, setProfileLoading] = useState(false);
    const [editMode, setEditMode] = useState(false);
    const [saving, setSaving] = useState(false);

    const [form, setForm] = useState({
        birth_date: "",
        gender: "" as "male" | "female" | "",
        weight_kg: "",
        height_cm: "",
    });

    const isMobile = window.innerWidth <= 768;

    const [darkMode, setDarkMode] = useState<boolean>(() => {
        const saved = localStorage.getItem("darkMode");
        if (saved !== null) return saved === "true";
        return window.matchMedia?.("(prefers-color-scheme: dark)").matches ?? false;
    });

    useEffect(() => {
        const root = document.documentElement;
        if (darkMode) root.classList.add("dark-theme");
        else root.classList.remove("dark-theme");
        localStorage.setItem("darkMode", darkMode.toString());
    }, [darkMode]);

    // Загрузка профиля
    useEffect(() => {
        if (!user) return;

        setProfileLoading(true);
        getProfile()
            .then((p: UserProfile) => {
                setProfile(p);
                setForm({
                    birth_date: p.birth_date ?? "",
                    gender: p.gender ?? "",
                    weight_kg: p.weight_kg?.toString() ?? "",
                    height_cm: p.height_cm?.toString() ?? "",
                });
            })
            .catch(() => setToast("Ошибка загрузки профиля ❌"))
            .finally(() => setProfileLoading(false));
    }, [user]);

    // Проверка подписки
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

    const saveIcon = async (name: IconName) => {
        updateIcon(name)
            .then(() => {
                setToast("Иконка обновлена ✅");
                setIconModalOpen(false);
            })
            .catch(() => setToast("Ошибка при сохранении ❌"));
    };

    const toggleNotifications = async () => {
        if (!("serviceWorker" in navigator)) return;
        const registration = await navigator.serviceWorker.ready;
        const subscription = await registration.pushManager.getSubscription();

        if (subscription) {
            await subscription.unsubscribe();
            await unsubscribePush();
            setNotificationsEnabled(false);
            setToast("Уведомления выключены ❌");
            return;
        }

        const permission = await Notification.requestPermission();
        if (permission === "granted") {
            const vapidPublicKey = await getVapidKey();
            await subscribePush(vapidPublicKey);
            setNotificationsEnabled(true);
            setToast("Уведомления включены ✅");
        } else {
            setToast("Разрешение на уведомления не выдано");
        }
    };

    const handleSaveProfile = async () => {
        setSaving(true);
        try {
            const payload: UpdateProfileRequest = {
                birth_date: form.birth_date || null,
                gender: form.gender || null,
                weight_kg: form.weight_kg ? parseFloat(form.weight_kg) : null,
                height_cm: form.height_cm ? parseInt(form.height_cm) : null,
            };

            await updateProfile(payload);
            setProfile(prev => prev ? {...prev, ...payload} : null);
            setEditMode(false);
            setToast("Профиль обновлён ✅");
        } catch {
            setToast("Ошибка сохранения ❌");
        } finally {
            setSaving(false);
        }
    };

    const cancelEdit = () => {
        setEditMode(false);
        if (profile) {
            setForm({
                birth_date: profile.birth_date ?? "",
                gender: profile.gender ?? "",
                weight_kg: profile.weight_kg?.toString() ?? "",
                height_cm: profile.height_cm?.toString() ?? "",
            });
        }
    };

    const inputStyle: React.CSSProperties = {
        width: "100%",
        padding: "10px 12px",
        borderRadius: "var(--radius-md)",
        border: "1px solid var(--color-border)",
        background: "var(--color-input)",
        color: "var(--color-text)",
        fontSize: 14,
    };

    const labelStyle: React.CSSProperties = {
        fontSize: 14,
        fontWeight: 500,
        color: "var(--color-text-secondary)",
        marginBottom: 4,
        display: "block",
    };

    return (
        <div className={"page stack"}>
            {isMobile && (
                <div style={{display: "flex", justifyContent: "flex-end"}}>
                    <Button variant="attention" onClick={() => setDarkMode(!darkMode)}>
                        {darkMode ? <Sun/> : <Moon/>}
                    </Button>
                </div>
            )}

            {/* NOT LOGGED IN */}
            {!loading && !user && (
                <div className="card" style={{textAlign: "center"}}>
                    <UserLock size={42}/>
                    <div style={{fontSize: 16, fontWeight: 600, marginBottom: 16}}>
                        Войдите в аккаунт
                    </div>
                    <div className="stack">
                        <Button variant="primary" onClick={() => {
                            const state = crypto.randomUUID();
                            localStorage.setItem("oauth_state", state);
                            window.location.href = `/api/telegram/login?origin=${encodeURIComponent(window.location.origin)}&state=${state}`;
                        }}>
                            Войти через Telegram
                        </Button>
                        <Button variant="danger" onClick={() => {
                            const state = crypto.randomUUID();
                            localStorage.setItem("oauth_state", state);
                            window.location.href = `/api/yandex/login?origin=${encodeURIComponent(window.location.origin)}&state=${state}`;
                        }}>
                            Войти через Yandex
                        </Button>
                    </div>
                </div>
            )}

            {/* LOGGED IN */}
            {user && (
                <>
                    {/* Карточка профиля */}
                    <div className="card" style={{position: "relative", textAlign: "center"}}>
                        <div onClick={() => setIconModalOpen(true)} style={{
                            position: "absolute", top: 'var(--card-gap)', right: 'var(--card-gap)', cursor: "pointer", opacity: 0.5
                        }}>
                            <Button variant="ghost" size="sm">
                                <Pencil size={14}/>
                            </Button>
                        </div>

                        <CurrentIcon size={40}/>
                        <div style={{fontSize: 18, fontWeight: 600}}>
                            {user.first_name}
                        </div>
                        {user.username && (
                            <div style={{opacity: 0.6, fontSize: 14}}>@{user.username}</div>
                        )}
                    </div>

                    {/* Данные профиля */}
                    <div className="card" style={{display: "flex", flexDirection: "column", gap: 'var(--card-gap)'}}>
                        <div style={{display: "flex", justifyContent: "space-between", alignItems: "center"}}>
                            <h3 style={{margin: 0, fontSize: 16}}>Личные данные</h3>
                            {!editMode && (
                                <Button style={{opacity: 0.5}} variant="ghost" onClick={() => setEditMode(true)}
                                        size="sm">
                                    <Pencil size={14}/>
                                </Button>
                            )}
                        </div>

                        {profileLoading ? (
                            <div style={{textAlign: "center", padding: "20px 0", opacity: 0.5}}>
                                Загрузка...
                            </div>
                        ) : editMode ? (
                            <>
                                <div>
                                    <label style={labelStyle}><Cake size={12}/> Дата рождения</label>
                                    <input
                                        type="date"
                                        value={form.birth_date}
                                        onChange={e => setForm(f => ({...f, birth_date: e.target.value}))}
                                        style={inputStyle}
                                    />
                                </div>

                                <div>
                                    <label style={labelStyle}><User size={12}/> Пол</label>
                                    <select
                                        value={form.gender}
                                        onChange={e => setForm(f => ({...f, gender: e.target.value as any}))}
                                        style={inputStyle}
                                    >
                                        <option value="">Не указан</option>
                                        <option value="male">Мужской</option>
                                        <option value="female">Женский</option>
                                    </select>
                                </div>

                                <div style={{display: "grid", gridTemplateColumns: "1fr 1fr", gap: 'var(--card-gap)'}}>
                                    <div>
                                        <label style={labelStyle}><Weight size={12}/> Вес, кг</label>
                                        <input
                                            type="number"
                                            step="0.1"
                                            min="20"
                                            max="300"
                                            placeholder="75"
                                            value={form.weight_kg}
                                            onChange={e => setForm(f => ({...f, weight_kg: e.target.value}))}
                                            style={inputStyle}
                                        />
                                    </div>
                                    <div>
                                        <label style={labelStyle}><Ruler size={12}/> Рост, см</label>
                                        <input
                                            type="number"
                                            min="100"
                                            max="250"
                                            placeholder="180"
                                            value={form.height_cm}
                                            onChange={e => setForm(f => ({...f, height_cm: e.target.value}))}
                                            style={inputStyle}
                                        />
                                    </div>
                                </div>

                                <div style={{display: "flex", gap: 8, marginTop: 8}}>
                                    <Button variant="active" onClick={handleSaveProfile} disabled={saving}
                                            style={{flex: 1}}>
                                        {saving ? "Сохранение..." : "Сохранить"}
                                    </Button>
                                    <Button variant="ghost" onClick={cancelEdit}>
                                        Отмена
                                    </Button>
                                </div>
                            </>
                        ) : (
                            <div style={{display: "flex", flexDirection: "column", gap: 'var(--card-gap)'}}>
                                <ProfileRow icon={<Cake size={18}/>} label="Дата рождения"
                                            value={profile?.birth_date ? new Date(profile.birth_date).toLocaleDateString("ru-RU") : "—"}/>
                                <ProfileRow icon={<User size={18}/>} label="Пол"
                                            value={profile?.gender === "male" ? "Мужской" : profile?.gender === "female" ? "Женский" : "—"}/>
                                <ProfileRow icon={<Weight size={18}/>} label="Вес"
                                            value={profile?.weight_kg ? `${profile.weight_kg} кг` : "—"}/>
                                <ProfileRow icon={<Ruler size={18}/>} label="Рост"
                                            value={profile?.height_cm ? `${profile.height_cm} см` : "—"}/>
                            </div>
                        )}
                    </div>

                    {/* Остальные кнопки */}
                    {!checking && (
                        <Button variant={notificationsEnabled ? "ghost" : "active"} onClick={toggleNotifications}>
                            {notificationsEnabled ? <><BellOff size={16}/> Выключить уведомления</>
                                : <><Bell size={16}/> Включить уведомления</>}
                        </Button>
                    )}

                    <Button variant="primary" onClick={() => navigate(`/statistics`)}>
                        <ChartNoAxesCombined/> Динамика
                    </Button>

                    <Button variant="danger" onClick={logout} style={{width: "100%", height: 48, fontSize: 16}}>
                        <LogOut/> Выйти
                    </Button>

                    {toast && <Toast message={toast} onClose={() => setToast(null)}/>}
                    {iconModalOpen && (
                        <IconPicker selected={icon} onSelect={saveIcon} onClose={() => setIconModalOpen(false)}/>
                    )}
                </>
            )}
        </div>
    );
};

const ProfileRow: React.FC<{ icon: React.ReactNode; label: string; value: string }> = ({icon, label, value}) => (
    <div style={{display: "flex", alignItems: "center", gap: 'var(--card-gap)', padding: "6px 0"}}>
        <div style={{opacity: 0.5, display: "flex", alignItems: "center"}}>{icon}</div>
        <div style={{flex: 1, fontSize: 14}}>
            <span style={{opacity: 0.6}}>{label}</span>
        </div>
        <div style={{fontWeight: 500, fontSize: 14}}>{value}</div>
    </div>
);

export default ProfilePage;
