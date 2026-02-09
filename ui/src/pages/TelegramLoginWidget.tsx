import React, {useEffect, useRef} from 'react';
import {useAuth} from '../context/AuthContext.tsx';

const TelegramLoginWidget: React.FC = () => {
    const {user} = useAuth();
    const widgetRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (user || !widgetRef.current) return;

        widgetRef.current.innerHTML = '';

        const isDev = process.env.NODE_ENV === 'development';

        const botUsername = isDev
            ? 'fitness_gym_buddy_dev_bot'
            : 'form_journey_bot';

        console.log("botUsername", botUsername);

        const script = document.createElement('script');
        script.src = 'https://telegram.org/js/telegram-widget.js?15';
        script.async = true;
        script.setAttribute('data-telegram-login', botUsername);
        script.setAttribute('data-size', 'large');
        script.setAttribute('data-userpic', 'true');
        script.setAttribute('data-request-access', 'write');

        widgetRef.current.appendChild(script);
    }, [user]);

    if (user) return null;

    return <div ref={widgetRef}/>;
};

export default TelegramLoginWidget;
