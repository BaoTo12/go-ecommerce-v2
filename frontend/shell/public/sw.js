const CACHE_NAME = 'shopee-cache-v1';
const STATIC_ASSETS = [
    '/',
    '/manifest.json',
    '/offline.html',
];

// Install Service Worker
self.addEventListener('install', (event) => {
    event.waitUntil(
        caches.open(CACHE_NAME).then((cache) => {
            console.log('Caching static assets');
            return cache.addAll(STATIC_ASSETS);
        })
    );
    self.skipWaiting();
});

// Activate and clean up old caches
self.addEventListener('activate', (event) => {
    event.waitUntil(
        caches.keys().then((cacheNames) => {
            return Promise.all(
                cacheNames
                    .filter((name) => name !== CACHE_NAME)
                    .map((name) => caches.delete(name))
            );
        })
    );
    self.clients.claim();
});

// Fetch with cache-first strategy for static assets
self.addEventListener('fetch', (event) => {
    const { request } = event;
    const url = new URL(request.url);

    // Skip non-GET requests
    if (request.method !== 'GET') return;

    // Skip external requests
    if (url.origin !== self.location.origin) return;

    // Cache-first for static assets
    if (isStaticAsset(url.pathname)) {
        event.respondWith(
            caches.match(request).then((cachedResponse) => {
                if (cachedResponse) {
                    // Return cached and update in background
                    fetchAndCache(request);
                    return cachedResponse;
                }
                return fetchAndCache(request);
            })
        );
        return;
    }

    // Network-first for API and dynamic content
    if (url.pathname.startsWith('/api/') || url.pathname.startsWith('/_next/')) {
        event.respondWith(
            fetch(request)
                .then((response) => {
                    // Clone and cache successful responses
                    if (response.ok) {
                        const clone = response.clone();
                        caches.open(CACHE_NAME).then((cache) => cache.put(request, clone));
                    }
                    return response;
                })
                .catch(() => caches.match(request))
        );
        return;
    }

    // Stale-while-revalidate for pages
    event.respondWith(
        caches.match(request).then((cachedResponse) => {
            const fetchPromise = fetch(request).then((networkResponse) => {
                if (networkResponse.ok) {
                    const clone = networkResponse.clone();
                    caches.open(CACHE_NAME).then((cache) => cache.put(request, clone));
                }
                return networkResponse;
            });

            return cachedResponse || fetchPromise;
        })
    );
});

// Helper functions
function isStaticAsset(pathname) {
    return (
        pathname.match(/\.(js|css|png|jpg|jpeg|gif|svg|ico|woff|woff2)$/) ||
        pathname.startsWith('/icons/') ||
        pathname.startsWith('/images/')
    );
}

async function fetchAndCache(request) {
    const response = await fetch(request);
    if (response.ok) {
        const cache = await caches.open(CACHE_NAME);
        cache.put(request, response.clone());
    }
    return response;
}

// Background sync for offline actions
self.addEventListener('sync', (event) => {
    if (event.tag === 'sync-cart') {
        event.waitUntil(syncCart());
    }
    if (event.tag === 'sync-orders') {
        event.waitUntil(syncOrders());
    }
});

async function syncCart() {
    // Sync cart items when back online
    const clients = await self.clients.matchAll();
    clients.forEach((client) => {
        client.postMessage({ type: 'CART_SYNCED' });
    });
}

async function syncOrders() {
    // Sync pending orders when back online
    const clients = await self.clients.matchAll();
    clients.forEach((client) => {
        client.postMessage({ type: 'ORDERS_SYNCED' });
    });
}

// Push notifications
self.addEventListener('push', (event) => {
    if (!event.data) return;

    const data = event.data.json();
    const options = {
        body: data.body,
        icon: '/icons/icon-192x192.png',
        badge: '/icons/badge.png',
        vibrate: [100, 50, 100],
        data: { url: data.url || '/' },
        actions: [
            { action: 'view', title: 'Xem ngay' },
            { action: 'close', title: 'Đóng' },
        ],
    };

    event.waitUntil(
        self.registration.showNotification(data.title, options)
    );
});

// Handle notification click
self.addEventListener('notificationclick', (event) => {
    event.notification.close();

    if (event.action === 'view' || !event.action) {
        const url = event.notification.data?.url || '/';
        event.waitUntil(
            clients.openWindow(url)
        );
    }
});
