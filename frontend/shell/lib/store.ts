'use client';

import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface User {
    id: string;
    name: string;
    email: string;
    avatar?: string;
}

interface AuthState {
    user: User | null;
    isAuthenticated: boolean;
    login: (user: User) => void;
    logout: () => void;
}

export const useAuthStore = create<AuthState>()(
    persist(
        (set) => ({
            user: null,
            isAuthenticated: false,
            login: (user) => set({ user, isAuthenticated: true }),
            logout: () => set({ user: null, isAuthenticated: false }),
        }),
        { name: 'auth-storage' }
    )
);

interface GamificationState {
    balance: number;
    streak: number;
    lastCheckIn: string | null;
    setBalance: (balance: number) => void;
    addCoins: (amount: number) => void;
    setStreak: (streak: number) => void;
    setLastCheckIn: (date: string) => void;
}

export const useGamificationStore = create<GamificationState>()(
    persist(
        (set) => ({
            balance: 0,
            streak: 0,
            lastCheckIn: null,
            setBalance: (balance) => set({ balance }),
            addCoins: (amount) => set((state) => ({ balance: state.balance + amount })),
            setStreak: (streak) => set({ streak }),
            setLastCheckIn: (date) => set({ lastCheckIn: date }),
        }),
        { name: 'gamification-storage' }
    )
);

interface CartItem {
    id: string;
    productId: string;
    name: string;
    price: number;
    quantity: number;
    image?: string;
}

interface CartState {
    items: CartItem[];
    couponCode: string | null;
    discount: number;
    addItem: (item: Omit<CartItem, 'id'>) => void;
    removeItem: (id: string) => void;
    updateQuantity: (id: string, quantity: number) => void;
    applyCoupon: (code: string, discount: number) => void;
    clearCart: () => void;
    getTotal: () => number;
}

export const useCartStore = create<CartState>()(
    persist(
        (set, get) => ({
            items: [],
            couponCode: null,
            discount: 0,
            addItem: (item) =>
                set((state) => ({
                    items: [...state.items, { ...item, id: crypto.randomUUID() }],
                })),
            removeItem: (id) =>
                set((state) => ({
                    items: state.items.filter((i) => i.id !== id),
                })),
            updateQuantity: (id, quantity) =>
                set((state) => ({
                    items: state.items.map((i) =>
                        i.id === id ? { ...i, quantity } : i
                    ),
                })),
            applyCoupon: (code, discount) => set({ couponCode: code, discount }),
            clearCart: () => set({ items: [], couponCode: null, discount: 0 }),
            getTotal: () => {
                const state = get();
                const subtotal = state.items.reduce(
                    (sum, i) => sum + i.price * i.quantity,
                    0
                );
                return subtotal - state.discount;
            },
        }),
        { name: 'cart-storage' }
    )
);

interface NotificationState {
    notifications: Notification[];
    addNotification: (notification: Omit<Notification, 'id' | 'timestamp'>) => void;
    removeNotification: (id: string) => void;
    clearAll: () => void;
}

interface Notification {
    id: string;
    type: 'success' | 'error' | 'info' | 'warning';
    title: string;
    message: string;
    timestamp: Date;
}

export const useNotificationStore = create<NotificationState>((set) => ({
    notifications: [],
    addNotification: (notification) =>
        set((state) => ({
            notifications: [
                ...state.notifications,
                { ...notification, id: crypto.randomUUID(), timestamp: new Date() },
            ],
        })),
    removeNotification: (id) =>
        set((state) => ({
            notifications: state.notifications.filter((n) => n.id !== id),
        })),
    clearAll: () => set({ notifications: [] }),
}));
