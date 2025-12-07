// AI Recommendation Service - "Customers also bought" and personalized suggestions

import { Product, productService } from './productService';
import { recentlyViewedService } from './recentlyViewedService';

export interface RecommendationResult {
    products: Product[];
    reason: string;
    confidence: number;
}

// Simulated user preferences based on browsing history
interface UserPreferences {
    categories: Record<string, number>;
    priceRange: { min: number; max: number };
    brands: string[];
}

// Build preferences from browsing history
const buildUserPreferences = async (): Promise<UserPreferences> => {
    const history = await recentlyViewedService.getRecentlyViewed();
    const categories: Record<string, number> = {};
    const prices: number[] = [];
    const brands: string[] = [];

    history.forEach(item => {
        const cat = item.product.category;
        categories[cat] = (categories[cat] || 0) + 1;
        prices.push(item.product.price);
        if (item.product.shop?.name) {
            brands.push(item.product.shop.name);
        }
    });

    return {
        categories,
        priceRange: {
            min: prices.length ? Math.min(...prices) * 0.5 : 0,
            max: prices.length ? Math.max(...prices) * 1.5 : 50000000,
        },
        brands: [...new Set(brands)],
    };
};

// Cosine similarity for product matching
const calculateSimilarity = (product1: Product, product2: Product): number => {
    let score = 0;

    // Same category
    if (product1.category === product2.category) score += 0.4;

    // Similar price range (within 30%)
    const priceDiff = Math.abs(product1.price - product2.price) / Math.max(product1.price, product2.price);
    if (priceDiff < 0.3) score += 0.3 * (1 - priceDiff);

    // Same shop
    if (product1.shop?.id === product2.shop?.id) score += 0.2;

    // Similar rating
    const ratingDiff = Math.abs(product1.rating - product2.rating);
    if (ratingDiff < 1) score += 0.1 * (1 - ratingDiff);

    return score;
};

export const recommendationService = {
    // Get "Customers also bought" recommendations
    getAlsoBought: async (productId: string): Promise<RecommendationResult> => {
        await new Promise(resolve => setTimeout(resolve, 200));

        const currentProduct = await productService.getProduct(productId);
        if (!currentProduct) {
            return { products: [], reason: '', confidence: 0 };
        }

        const allProducts = await productService.getProducts({ limit: 50 });
        const productList = Array.isArray(allProducts) ? allProducts : allProducts.products || [];

        // Calculate similarity scores
        const scored = productList
            .filter(p => p.id !== productId)
            .map(p => ({
                product: p,
                score: calculateSimilarity(currentProduct, p),
            }))
            .sort((a, b) => b.score - a.score)
            .slice(0, 6);

        return {
            products: scored.map(s => s.product),
            reason: 'Khách hàng cũng mua',
            confidence: scored[0]?.score || 0,
        };
    },

    // Get "You may also like" based on browsing history
    getPersonalized: async (): Promise<RecommendationResult> => {
        await new Promise(resolve => setTimeout(resolve, 250));

        const preferences = await buildUserPreferences();
        const allProducts = await productService.getProducts({ limit: 50 });
        const productList = Array.isArray(allProducts) ? allProducts : allProducts.products || [];

        // Score products based on preferences
        const scored = productList
            .map(p => {
                let score = 0;

                // Category preference
                const catScore = preferences.categories[p.category] || 0;
                score += catScore * 0.4;

                // Price range fit
                if (p.price >= preferences.priceRange.min && p.price <= preferences.priceRange.max) {
                    score += 0.3;
                }

                // Brand preference
                if (preferences.brands.includes(p.shop?.name || '')) {
                    score += 0.2;
                }

                // High rating bonus
                if (p.rating >= 4.5) score += 0.1;

                return { product: p, score };
            })
            .sort((a, b) => b.score - a.score)
            .slice(0, 12);

        return {
            products: scored.map(s => s.product),
            reason: 'Dành riêng cho bạn',
            confidence: scored[0]?.score || 0,
        };
    },

    // Get "Frequently bought together"
    getFrequentlyBoughtTogether: async (productId: string): Promise<RecommendationResult> => {
        await new Promise(resolve => setTimeout(resolve, 150));

        const currentProduct = await productService.getProduct(productId);
        if (!currentProduct) {
            return { products: [], reason: '', confidence: 0 };
        }

        // Simulate complementary products
        const allProducts = await productService.getProducts({ limit: 30 });
        const productList = Array.isArray(allProducts) ? allProducts : allProducts.products || [];

        const complementary = productList
            .filter(p => p.id !== productId && p.category !== currentProduct.category)
            .slice(0, 3);

        return {
            products: complementary,
            reason: 'Thường được mua cùng',
            confidence: 0.7,
        };
    },

    // Get "Because you viewed X"
    getBecauseYouViewed: async (productId: string): Promise<RecommendationResult> => {
        await new Promise(resolve => setTimeout(resolve, 200));

        const viewedProduct = await productService.getProduct(productId);
        if (!viewedProduct) {
            return { products: [], reason: '', confidence: 0 };
        }

        const similar = await productService.getProducts({
            category: viewedProduct.category,
            limit: 8
        });
        const productList = Array.isArray(similar) ? similar : similar.products || [];

        return {
            products: productList.filter(p => p.id !== productId).slice(0, 6),
            reason: `Vì bạn đã xem "${viewedProduct.name.substring(0, 30)}..."`,
            confidence: 0.8,
        };
    },

    // Get trending in your area
    getTrendingNearYou: async (location: string = 'TP. Hồ Chí Minh'): Promise<RecommendationResult> => {
        await new Promise(resolve => setTimeout(resolve, 150));

        const allProducts = await productService.getProducts({ limit: 30 });
        const productList = Array.isArray(allProducts) ? allProducts : allProducts.products || [];

        const localProducts = productList
            .filter(p => p.location === location)
            .sort((a, b) => b.sold - a.sold)
            .slice(0, 6);

        return {
            products: localProducts,
            reason: `Đang hot tại ${location}`,
            confidence: 0.75,
        };
    },
};

export default recommendationService;
