import {useState} from "react";
import {api} from "../api/client.ts";

export function useVideo(originalUrl: string | null) {
    const [videoUrl, setVideoUrl] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetchVideo = async () => {
        if (!originalUrl) return;

        const CACHE_KEY = `video-${originalUrl}`;
        const cached = localStorage.getItem(CACHE_KEY);

        if (cached) {
            const {url, expires} = JSON.parse(cached);
            if (Date.now() < expires) {
                setVideoUrl(url);
                return;
            }
        }

        try {
            setLoading(true);
            const data = await api<{ url: string }>(
                `/api/video/link?url=${encodeURIComponent(originalUrl)}`
            );
            setVideoUrl(data.url);

            localStorage.setItem(
                CACHE_KEY,
                JSON.stringify({url: data.url, expires: Date.now() + 4 * 60 * 1000})
            );
        } catch (e: any) {
            setError("Ошибка загрузки видео 😢");
        } finally {
            setLoading(false);
        }
    };

    return {videoUrl, loading, error, fetchVideo};
}
