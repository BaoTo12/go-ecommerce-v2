// Cart Service - Handles shopping cart with localStorage persistence

import { Product } from './productService';

export interface CartItem {
    id: string;
    productId: string;
    product: Product;
    quantity: number;
    variant?: string;
    selected: boolean;
    addedAt: string;
}

export interface Cart {
    items: CartItem[];
    updatedAt: string;
}

// Initialize cart from localStorage
const getStoredCart = (): Cart => {
    if (typeof window === 'undefined') {
        return { items: [], updatedAt: new Date().toISOString() };
    }
    const saved = localStorage.getItem('cart');
    if (saved) {
        try {
            return JSON.parse(saved);
        } catch (e) {
            console.error('Failed to parse cart:', e);
        }
    }
    return { items: [], updatedAt: new Date().toISOString() };
};

let cart: Cart = getStoredCart();

const saveCart = () => {
    if (typeof window !== 'undefined') {
        cart.updatedAt = new Date().toISOString();
        localStorage.setItem('cart', JSON.stringify(cart));
    }
};

export const cartService = {
    // Get cart
    getCart: async (): Promise<Cart> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        cart = getStoredCart(); // Refresh from storage
        return cart;
    },

    // Get cart items count
    getItemCount: (): number => {
        return cart.items.reduce((sum, item) => sum + item.quantity, 0);
    },

    // Add item to cart
    addItem: async (product: Product, quantity: number = 1, variant?: string): Promise<CartItem> => {
        await new Promise(resolve => setTimeout(resolve, 100));

        // Check if item already exists
        const existingIndex = cart.items.findIndex(
            item => item.productId === product.id && item.variant === variant
        );

        if (existingIndex >= 0) {
            cart.items[existingIndex].quantity += quantity;
            saveCart();
            return cart.items[existingIndex];
        }

        const newItem: CartItem = {
            id: `cart_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
            productId: product.id,
            product,
            quantity,
            variant,
            selected: true,
            addedAt: new Date().toISOString(),
        };

        cart.items.push(newItem);
        saveCart();
        return newItem;
    },

    // Update item quantity
    updateQuantity: async (itemId: string, quantity: number): Promise<boolean> => {
        await new Promise(resolve => setTimeout(resolve, 50));

        const item = cart.items.find(i => i.id === itemId);
        if (item) {
            if (quantity <= 0) {
                cart.items = cart.items.filter(i => i.id !== itemId);
            } else {
                item.quantity = quantity;
            }
            saveCart();
            return true;
        }
        return false;
    },

    // Remove item
    removeItem: async (itemId: string): Promise<boolean> => {
        await new Promise(resolve => setTimeout(resolve, 50));

        const initialLength = cart.items.length;
        cart.items = cart.items.filter(i => i.id !== itemId);

        if (cart.items.length !== initialLength) {
            saveCart();
            return true;
        }
        return false;
    },

    // Toggle item selection
    toggleSelection: async (itemId: string): Promise<boolean> => {
        const item = cart.items.find(i => i.id === itemId);
        if (item) {
            item.selected = !item.selected;
            saveCart();
            return true;
        }
        return false;
    },

    // Select all items
    selectAll: async (selected: boolean): Promise<void> => {
        cart.items.forEach(item => {
            item.selected = selected;
        });
        saveCart();
    },

    // Get selected items
    getSelectedItems: (): CartItem[] => {
        return cart.items.filter(item => item.selected);
    },

    // Calculate totals
    calculateTotals: (): { subtotal: number; itemCount: number; selectedCount: number } => {
        const selectedItems = cart.items.filter(item => item.selected);
        return {
            subtotal: selectedItems.reduce((sum, item) => sum + item.product.price * item.quantity, 0),
            itemCount: cart.items.reduce((sum, item) => sum + item.quantity, 0),
            selectedCount: selectedItems.reduce((sum, item) => sum + item.quantity, 0),
        };
    },

    // Clear cart
    clearCart: async (): Promise<void> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        cart.items = [];
        saveCart();
    },

    // Clear selected items (after checkout)
    clearSelected: async (): Promise<void> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        cart.items = cart.items.filter(item => !item.selected);
        saveCart();
    },
};

export default cartService;
