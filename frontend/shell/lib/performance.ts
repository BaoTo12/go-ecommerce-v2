// Frontend performance utilities

/**
 * Debounces a function to limit how often it can be called
 */
export function debounce<T extends (...args: unknown[]) => unknown>(
    fn: T,
    delay: number
): (...args: Parameters<T>) => void {
    let timeoutId: NodeJS.Timeout;
    return (...args: Parameters<T>) => {
        clearTimeout(timeoutId);
        timeoutId = setTimeout(() => fn(...args), delay);
    };
}

/**
 * Throttles a function to run at most once per interval
 */
export function throttle<T extends (...args: unknown[]) => unknown>(
    fn: T,
    interval: number
): (...args: Parameters<T>) => void {
    let lastTime = 0;
    return (...args: Parameters<T>) => {
        const now = Date.now();
        if (now - lastTime >= interval) {
            lastTime = now;
            fn(...args);
        }
    };
}

/**
 * Memoizes an async function with TTL
 */
export function memoize<T extends (...args: unknown[]) => Promise<unknown>>(
    fn: T,
    ttl: number = 60000 // 1 minute default
): T {
    const cache = new Map<string, { value: unknown; expires: number }>();

    return (async (...args: Parameters<T>) => {
        const key = JSON.stringify(args);
        const cached = cache.get(key);

        if (cached && cached.expires > Date.now()) {
            return cached.value;
        }

        const result = await fn(...args);
        cache.set(key, { value: result, expires: Date.now() + ttl });

        // Cleanup expired entries periodically
        if (cache.size > 100) {
            const now = Date.now();
            cache.forEach((v, k) => {
                if (v.expires < now) cache.delete(k);
            });
        }

        return result;
    }) as T;
}

/**
 * Lazy loads an image and returns a promise
 */
export function lazyLoadImage(src: string): Promise<HTMLImageElement> {
    return new Promise((resolve, reject) => {
        const img = new Image();
        img.onload = () => resolve(img);
        img.onerror = reject;
        img.src = src;
    });
}

/**
 * Preloads critical images
 */
export function preloadImages(urls: string[]): Promise<HTMLImageElement[]> {
    return Promise.all(urls.map(lazyLoadImage));
}

/**
 * Creates an intersection observer for lazy loading
 */
export function createLazyLoader(
    callback: (entry: IntersectionObserverEntry) => void,
    options?: IntersectionObserverInit
): IntersectionObserver {
    return new IntersectionObserver((entries) => {
        entries.forEach((entry) => {
            if (entry.isIntersecting) {
                callback(entry);
            }
        });
    }, {
        rootMargin: '50px',
        threshold: 0.1,
        ...options,
    });
}

/**
 * Batches multiple updates into a single requestAnimationFrame
 */
export function batchUpdates(updates: (() => void)[]): void {
    requestAnimationFrame(() => {
        updates.forEach((update) => update());
    });
}

/**
 * Measures component render time
 */
export function measureRender(name: string): { start: () => void; end: () => void } {
    let startTime: number;
    return {
        start: () => {
            startTime = performance.now();
        },
        end: () => {
            const duration = performance.now() - startTime;
            if (duration > 16) { // More than 1 frame at 60fps
                console.warn(`[Perf] ${name} took ${duration.toFixed(2)}ms`);
            }
        },
    };
}

/**
 * Prefetches a route for faster navigation
 */
export function prefetchRoute(href: string): void {
    if (typeof window !== 'undefined') {
        const link = document.createElement('link');
        link.rel = 'prefetch';
        link.href = href;
        document.head.appendChild(link);
    }
}

/**
 * Chunks an array for pagination
 */
export function chunk<T>(array: T[], size: number): T[][] {
    const chunks: T[][] = [];
    for (let i = 0; i < array.length; i += size) {
        chunks.push(array.slice(i, i + size));
    }
    return chunks;
}

/**
 * Virtual list helper - calculates visible items
 */
export function getVisibleRange(
    scrollTop: number,
    viewportHeight: number,
    itemHeight: number,
    totalItems: number,
    overscan: number = 3
): { start: number; end: number } {
    const start = Math.max(0, Math.floor(scrollTop / itemHeight) - overscan);
    const end = Math.min(totalItems, Math.ceil((scrollTop + viewportHeight) / itemHeight) + overscan);
    return { start, end };
}
