import {useEffect} from "react";
import {useNavigate} from "react-router-dom";
import {useAuth} from "../context/AuthContext.tsx";

const AuthTelegram = () => {
    const navigate = useNavigate();
    const {refreshUser} = useAuth();

    useEffect(() => {
        const hash = window.location.hash;

        if (!hash.startsWith("#tgAuthResult=")) {
            navigate("/profile");
            return;
        }

        const encoded = hash.replace("#tgAuthResult=", "");
        const decoded = JSON.parse(atob(encoded));

        fetch("/api/telegram/login", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(decoded),
        })
            .then(res => res.json())
            .then(data => {
                localStorage.setItem("token", data.token);
                refreshUser();       // ← обновляем user сразу
                navigate("/");             // ← редирект на главную
            });
    }, []);

    return null;
};

export default AuthTelegram;
