// Wishlist Service - Handles favorites/wishlist with localStorage persistence

import { Product } from './productService';

export interface WishlistItem {
    id: string;
    productId: string;
    product: Product;
    addedAt: string;
}

export interface Wishlist {
    items: WishlistItem[];
    updatedAt: string;
}

// Initialize wishlist from localStorage
const getStoredWishlist = (): Wishlist => {
    if (typeof window === 'undefined') {
        return { items: [], updatedAt: new Date().toISOString() };
    }
    const saved = localStorage.getItem('wishlist');
    if (saved) {
        try {
            return JSON.parse(saved);
        } catch (e) {
            console.error('Failed to parse wishlist:', e);
        }
    }
    return { items: [], updatedAt: new Date().toISOString() };
};

let wishlist: Wishlist = getStoredWishlist();

const saveWishlist = () => {
    if (typeof window !== 'undefined') {
        wishlist.updatedAt = new Date().toISOString();
        localStorage.setItem('wishlist', JSON.stringify(wishlist));
    }
};

export const wishlistService = {
    // Get wishlist
    getWishlist: async (): Promise<Wishlist> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        wishlist = getStoredWishlist();
        return wishlist;
    },

    // Get wishlist count
    getItemCount: (): number => {
        return wishlist.items.length;
    },

    // Check if product is in wishlist
    isInWishlist: (productId: string): boolean => {
        return wishlist.items.some(item => item.productId === productId);
    },

    // Add to wishlist
    addItem: async (product: Product): Promise<WishlistItem | null> => {
        await new Promise(resolve => setTimeout(resolve, 100));

        // Check if already in wishlist
        if (wishlist.items.some(item => item.productId === product.id)) {
            return null;
        }

        const newItem: WishlistItem = {
            id: `wish_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
            productId: product.id,
            product,
            addedAt: new Date().toISOString(),
        };

        wishlist.items.push(newItem);
        saveWishlist();
        return newItem;
    },

    // Remove from wishlist
    removeItem: async (productId: string): Promise<boolean> => {
        await new Promise(resolve => setTimeout(resolve, 50));

        const initialLength = wishlist.items.length;
        wishlist.items = wishlist.items.filter(i => i.productId !== productId);

        if (wishlist.items.length !== initialLength) {
            saveWishlist();
            return true;
        }
        return false;
    },

    // Toggle wishlist (add if not in, remove if in)
    toggleItem: async (product: Product): Promise<{ added: boolean }> => {
        const isIn = wishlist.items.some(item => item.productId === product.id);

        if (isIn) {
            await wishlistService.removeItem(product.id);
            return { added: false };
        } else {
            await wishlistService.addItem(product);
            return { added: true };
        }
    },

    // Clear wishlist
    clearWishlist: async (): Promise<void> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        wishlist.items = [];
        saveWishlist();
    },
};

export default wishlistService;
