import { api } from './client.ts';

// Конвертация VAPID из base64 в Uint8Array
function urlBase64ToUint8Array(base64String: string) {
    const padding = '='.repeat((4 - base64String.length % 4) % 4);
    const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/');
    const rawData = atob(base64);
    return new Uint8Array([...rawData].map(c => c.charCodeAt(0)));
}

export async function subscribePush(VAPID_PUBLIC_KEY: string) {
    if (!('serviceWorker' in navigator) || !('PushManager' in window)) return null;

    const registration = await navigator.serviceWorker.ready;

    const existing = await registration.pushManager.getSubscription();
    if (existing) return existing;

    const subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: urlBase64ToUint8Array(VAPID_PUBLIC_KEY),
    });

    await api('/api/push/subscribe', {
        method: 'POST',
        body: JSON.stringify(subscription),
    });

    return subscription;
}
