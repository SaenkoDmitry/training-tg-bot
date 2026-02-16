/// <reference lib="webworker" />
import { precacheAndRoute } from 'workbox-precaching'

declare let self: ServiceWorkerGlobalScope

precacheAndRoute(self.__WB_MANIFEST)

// PUSH
self.addEventListener('push', (event) => {
    const data = event.data?.json() || {}

    const options: NotificationOptions = {
        body: data.body || 'Таймер завершён!',
        icon: '/web-app-manifest-192x192.png',
        badge: '/web-app-manifest-192x192.png',
        requireInteraction: true,
        tag: data.tag || '',
        data: {
            url: data.url || '/'
        }
    }

    event.waitUntil(
        self.registration.showNotification(
            data.title || 'Form Journey',
            options
        )
    )
})

// CLICK
self.addEventListener('notificationclick', (event) => {
    event.notification.close()

    const url = event.notification.data?.url || '/'

    event.waitUntil(
        self.clients.matchAll({ type: 'window', includeUncontrolled: true })
            .then((clientsArr) => {
                for (const client of clientsArr) {
                    if (client.url.includes(url) && 'focus' in client)
                        return client.focus()
                }
                return self.clients.openWindow(url)
            })
    )
})
