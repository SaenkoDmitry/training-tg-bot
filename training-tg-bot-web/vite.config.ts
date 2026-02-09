import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
    plugins: [react()],
    server: {
        // разрешаем локальный сервер для туннельных хостов
        host: true,
        strictPort: false,
        port: 5173,
        allowedHosts: [
            'localhost',
            '127.0.0.1',
            '2237bc05745959.lhr.life',
        ],
        proxy: {
            '/api': 'http://localhost:8080',
        }
    }
})
