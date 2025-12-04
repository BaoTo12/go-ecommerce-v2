// API Configuration and Base Client
const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

async function apiRequest<T>(
    endpoint: string,
    options: RequestInit = {}
): Promise<T> {
    const response = await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            ...options.headers,
        },
    });

    if (!response.ok) {
        const error = await response.json().catch(() => ({ message: 'Request failed' }));
        throw new Error(error.message || `HTTP ${response.status}`);
    }

    return response.json();
}

// Flash Sale API
export const flashSaleApi = {
    getChallenge: (saleId: string, userId: string) =>
        apiRequest<{ challenge: string; difficulty: number }>(
            `/api/v1/flash-sale/challenge?sale_id=${saleId}&user_id=${userId}`
        ),

    attemptPurchase: (data: {
        flash_sale_id: string;
        user_id: string;
        quantity: number;
        challenge: string;
        nonce: string;
    }) =>
        apiRequest<{ reservation_id: string; expires_at: string }>(
            '/api/v1/flash-sale/purchase',
            { method: 'POST', body: JSON.stringify(data) }
        ),

    getActiveSales: () =>
        apiRequest<FlashSale[]>('/api/v1/flash-sales/active'),
};

// Gamification API
export const gamificationApi = {
    getBalance: (userId: string) =>
        apiRequest<CoinWallet>(`/api/v1/coins/balance?user_id=${userId}`),

    earnCoins: (data: { user_id: string; amount: number; source: string; description: string }) =>
        apiRequest<CoinWallet>('/api/v1/coins/earn', {
            method: 'POST',
            body: JSON.stringify(data),
        }),

    dailyCheckIn: (userId: string) =>
        apiRequest<{ reward: number; streak: number; streak_bonus: boolean }>(
            '/api/v1/check-in',
            { method: 'POST', body: JSON.stringify({ user_id: userId }) }
        ),

    spinLuckyDraw: (userId: string, spinCost: number) =>
        apiRequest<LuckyDrawResult>('/api/v1/lucky-draw/spin', {
            method: 'POST',
            body: JSON.stringify({ user_id: userId, spin_cost: spinCost }),
        }),

    getMissions: (userId: string) =>
        apiRequest<{ missions: Mission[]; user_progress: UserMission[] }>(
            `/api/v1/missions?user_id=${userId}`
        ),
};

// Fraud Detection API
export const fraudApi = {
    checkTransaction: (data: {
        transaction_id: string;
        user_id: string;
        amount: number;
        currency: string;
        ip: string;
        device_id: string;
    }) =>
        apiRequest<FraudCheck>('/api/v1/fraud/check', {
            method: 'POST',
            body: JSON.stringify(data),
        }),

    getAlerts: () => apiRequest<FraudAlert[]>('/api/v1/fraud/alerts'),

    getFraudHistory: (userId: string) =>
        apiRequest<FraudCheck[]>(`/api/v1/fraud/history?user_id=${userId}`),
};

// Analytics API
export const analyticsApi = {
    trackEvent: (data: {
        user_id: string;
        session_id: string;
        event_type: string;
        properties: Record<string, unknown>;
    }) =>
        apiRequest<{ event_id: string }>('/api/v1/track', {
            method: 'POST',
            body: JSON.stringify(data),
        }),

    getDashboard: () => apiRequest<DashboardMetrics>('/api/v1/analytics/dashboard'),

    getSalesReport: (period: string) =>
        apiRequest<SalesReport>(`/api/v1/analytics/sales?period=${period}`),

    getConversionFunnel: (name: string) =>
        apiRequest<ConversionFunnel>(`/api/v1/analytics/funnel?name=${name}`),
};

// Pricing API
export const pricingApi = {
    getPrice: (productId: string) =>
        apiRequest<ProductPrice>(`/api/v1/prices?product_id=${productId}`),

    optimizePrice: (productId: string) =>
        apiRequest<ProductPrice>('/api/v1/prices/optimize', {
            method: 'POST',
            body: JSON.stringify({ product_id: productId }),
        }),
};

// Coupon API  
export const couponApi = {
    validate: (code: string, userId: string, cartTotal: number) =>
        apiRequest<CouponValidation>('/api/v1/coupons/validate', {
            method: 'POST',
            body: JSON.stringify({ code, user_id: userId, cart_total: cartTotal }),
        }),

    apply: (code: string, userId: string, orderId: string) =>
        apiRequest<{ success: boolean; discount: number }>('/api/v1/coupons/apply', {
            method: 'POST',
            body: JSON.stringify({ code, user_id: userId, order_id: orderId }),
        }),

    getActiveCoupons: () => apiRequest<Coupon[]>('/api/v1/coupons/active'),
};

// Livestream API
export const livestreamApi = {
    getStreams: () => apiRequest<Livestream[]>('/api/v1/streams'),
    getStream: (streamId: string) => apiRequest<Livestream>(`/api/v1/streams/${streamId}`),
};

// Types
export interface FlashSale {
    id: string;
    product_id: string;
    original_price: number;
    sale_price: number;
    discount_percent: number;
    total_quantity: number;
    sold_quantity: number;
    max_per_user: number;
    status: string;
    start_time: string;
    end_time: string;
}

export interface CoinWallet {
    user_id: string;
    balance: number;
    lifetime: number;
    updated_at: string;
}

export interface LuckyDrawResult {
    id: string;
    user_id: string;
    prize_id: string;
    prize: {
        id: string;
        name: string;
        type: string;
        value: number;
    };
}

export interface Mission {
    id: string;
    name: string;
    description: string;
    type: string;
    target: number;
    reward: number;
}

export interface UserMission {
    user_id: string;
    mission_id: string;
    progress: number;
    completed: boolean;
    claimed_at?: string;
}

export interface FraudCheck {
    check_id: string;
    transaction_id: string;
    score: number;
    risk_level: string;
    decision: string;
    reasons: string[];
    processing_time: number;
}

export interface FraudAlert {
    id: string;
    fraud_check_id: string;
    alert_type: string;
    severity: string;
    message: string;
    acknowledged: boolean;
}

export interface DashboardMetrics {
    active_users: number;
    orders_in_progress: number;
    current_revenue: number;
    today_orders: number;
    today_revenue: number;
    today_page_views: number;
    orders_change: number;
    revenue_change: number;
}

export interface SalesReport {
    period: string;
    total_orders: number;
    total_revenue: number;
    avg_order_value: number;
}

export interface ConversionFunnel {
    name: string;
    steps: { name: string; users: number; conversion: number }[];
    overall_rate: number;
}

export interface ProductPrice {
    product_id: string;
    base_price: number;
    optimized_price: number;
    competitor_prices: number[];
}

export interface Coupon {
    id: string;
    code: string;
    discount_type: string;
    discount_value: number;
    min_purchase: number;
    max_discount: number;
}

export interface CouponValidation {
    valid: boolean;
    discount: number;
    message: string;
}

export interface Livestream {
    id: string;
    host_id: string;
    title: string;
    status: string;
    viewer_count: number;
    product_ids: string[];
}
