export const getVapidKey = async (): Promise<string> => {
    const res = await fetch('/api/vapid-key');
    const data = await res.json();
    return data.public_key;
};