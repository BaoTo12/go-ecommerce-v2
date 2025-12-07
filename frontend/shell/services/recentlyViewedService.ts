// Recently Viewed Service - Track product browsing history

import { Product } from './productService';

const MAX_HISTORY = 20;

export interface RecentlyViewedItem {
    product: Product;
    viewedAt: string;
}

// Get from localStorage
const getStoredHistory = (): RecentlyViewedItem[] => {
    if (typeof window === 'undefined') return [];
    const saved = localStorage.getItem('recentlyViewed');
    if (saved) {
        try {
            return JSON.parse(saved);
        } catch (e) {
            console.error('Failed to parse recently viewed:', e);
        }
    }
    return [];
};

let history = getStoredHistory();

const saveHistory = () => {
    if (typeof window !== 'undefined') {
        localStorage.setItem('recentlyViewed', JSON.stringify(history));
    }
};

export const recentlyViewedService = {
    // Get recently viewed products
    getRecentlyViewed: async (): Promise<RecentlyViewedItem[]> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        history = getStoredHistory();
        return history;
    },

    // Add product to history
    addToHistory: (product: Product): void => {
        // Remove if already exists
        history = history.filter(item => item.product.id !== product.id);

        // Add to beginning
        history.unshift({
            product,
            viewedAt: new Date().toISOString(),
        });

        // Trim to max size
        if (history.length > MAX_HISTORY) {
            history = history.slice(0, MAX_HISTORY);
        }

        saveHistory();
    },

    // Clear history
    clearHistory: async (): Promise<void> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        history = [];
        saveHistory();
    },

    // Get count
    getCount: (): number => {
        return history.length;
    },
};

export default recentlyViewedService;
