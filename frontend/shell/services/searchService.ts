// Search Service - Autocomplete and search history

import { Product, productService } from './productService';

export interface SearchSuggestion {
    type: 'history' | 'trending' | 'product' | 'category';
    text: string;
    productId?: string;
    image?: string;
    price?: number;
}

// Get from localStorage
const getSearchHistory = (): string[] => {
    if (typeof window === 'undefined') return [];
    const saved = localStorage.getItem('searchHistory');
    if (saved) {
        try {
            return JSON.parse(saved);
        } catch (e) {
            return [];
        }
    }
    return [];
};

let searchHistory = getSearchHistory();

const saveHistory = () => {
    if (typeof window !== 'undefined') {
        localStorage.setItem('searchHistory', JSON.stringify(searchHistory));
    }
};

// Trending searches (mock)
const TRENDING_SEARCHES = [
    'iPhone 15 Pro Max',
    'Áo hoodie form rộng',
    'Son Dior',
    'Giày Nike Air Force 1',
    'Tai nghe bluetooth',
    'Kem chống nắng',
    'Laptop gaming',
    'Váy đầm nữ',
];

export const searchService = {
    // Get autocomplete suggestions
    getSuggestions: async (query: string): Promise<SearchSuggestion[]> => {
        await new Promise(resolve => setTimeout(resolve, 150));

        if (!query.trim()) {
            // Return history + trending when no query
            const historySuggestions: SearchSuggestion[] = searchHistory.slice(0, 5).map(text => ({
                type: 'history',
                text,
            }));

            const trendingSuggestions: SearchSuggestion[] = TRENDING_SEARCHES.slice(0, 5).map(text => ({
                type: 'trending',
                text,
            }));

            return [...historySuggestions, ...trendingSuggestions];
        }

        const lowerQuery = query.toLowerCase();
        const suggestions: SearchSuggestion[] = [];

        // History matches
        searchHistory
            .filter(h => h.toLowerCase().includes(lowerQuery))
            .slice(0, 3)
            .forEach(text => {
                suggestions.push({ type: 'history', text });
            });

        // Trending matches
        TRENDING_SEARCHES
            .filter(t => t.toLowerCase().includes(lowerQuery))
            .slice(0, 3)
            .forEach(text => {
                suggestions.push({ type: 'trending', text });
            });

        // Product matches
        const response = await productService.getProducts({ search: query });
        const productList = Array.isArray(response) ? response : (response as { products: Product[] }).products || [];
        productList.slice(0, 4).forEach((product: Product) => {
            suggestions.push({
                type: 'product',
                text: product.name,
                productId: product.id,
                image: product.thumbnail,
                price: product.price,
            });
        });

        // Category suggestions
        const categories = ['Điện thoại', 'Máy tính', 'Thời trang', 'Làm đẹp', 'Giày dép'];
        categories
            .filter(c => c.toLowerCase().includes(lowerQuery))
            .slice(0, 2)
            .forEach(text => {
                suggestions.push({ type: 'category', text });
            });

        return suggestions.slice(0, 10);
    },

    // Add to history
    addToHistory: (query: string): void => {
        if (!query.trim()) return;

        // Remove if exists and add to front
        searchHistory = searchHistory.filter(h => h.toLowerCase() !== query.toLowerCase());
        searchHistory.unshift(query);

        // Keep only last 20
        searchHistory = searchHistory.slice(0, 20);
        saveHistory();
    },

    // Clear history
    clearHistory: (): void => {
        searchHistory = [];
        saveHistory();
    },

    // Get trending searches
    getTrending: (): string[] => TRENDING_SEARCHES,

    // Remove from history
    removeFromHistory: (query: string): void => {
        searchHistory = searchHistory.filter(h => h !== query);
        saveHistory();
    },
};

export default searchService;
