// API client for connecting frontend to backend services

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost';

const SERVICES = {
    recommendation: `${API_BASE}:8081`,
    auction: `${API_BASE}:8082`,
    referral: `${API_BASE}:8083`,
    tracking: `${API_BASE}:8084`,
};

class ApiClient {
    private baseUrl: string;

    constructor(service: keyof typeof SERVICES) {
        this.baseUrl = SERVICES[service];
    }

    async get<T>(path: string): Promise<T> {
        const response = await fetch(`${this.baseUrl}${path}`, {
            headers: { 'Content-Type': 'application/json' },
        });
        if (!response.ok) {
            throw new Error(`API Error: ${response.statusText}`);
        }
        return response.json();
    }

    async post<T>(path: string, data: unknown): Promise<T> {
        const response = await fetch(`${this.baseUrl}${path}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        });
        if (!response.ok) {
            throw new Error(`API Error: ${response.statusText}`);
        }
        return response.json();
    }
}

// Recommendation API
export const recommendationApi = {
    client: new ApiClient('recommendation'),

    getAlsoBought: (productId: string) =>
        recommendationApi.client.get(`/api/v1/recommendations/also-bought?product_id=${productId}`),

    getPersonalized: (userId: string) =>
        recommendationApi.client.get(`/api/v1/recommendations/personalized?user_id=${userId}`),

    getTrending: (location: string) =>
        recommendationApi.client.get(`/api/v1/recommendations/trending?location=${encodeURIComponent(location)}`),

    getBecauseViewed: (productId: string) =>
        recommendationApi.client.get(`/api/v1/recommendations/because-viewed?product_id=${productId}`),
};

// Auction API
export const auctionApi = {
    client: new ApiClient('auction'),

    getAuctions: () =>
        auctionApi.client.get<{ auctions: Auction[]; total: number }>('/api/v1/auctions'),

    getAuction: (id: string) =>
        auctionApi.client.get<Auction>(`/api/v1/auctions/${id}`),

    placeBid: (auctionId: string, userId: string, userName: string, amount: number) =>
        auctionApi.client.post(`/api/v1/auctions/${auctionId}/bid`, { user_id: userId, user_name: userName, amount }),

    getBidHistory: (auctionId: string) =>
        auctionApi.client.get<{ bids: Bid[] }>(`/api/v1/auctions/${auctionId}/bids`),
};

// Referral API
export const referralApi = {
    client: new ApiClient('referral'),

    generateCode: (userId: string, userName: string) =>
        referralApi.client.post<ReferralCode>('/api/v1/referrals/generate-code', { user_id: userId, user_name: userName }),

    getStats: (userId: string) =>
        referralApi.client.get<ReferralStats>(`/api/v1/referrals/stats?user_id=${userId}`),

    getReferrals: (userId: string) =>
        referralApi.client.get<{ referrals: Referral[]; total: number }>(`/api/v1/referrals?user_id=${userId}`),

    redeemCode: (code: string, userId: string, userName: string) =>
        referralApi.client.post('/api/v1/referrals/redeem', { code, user_id: userId, user_name: userName }),
};

// Tracking API
export const trackingApi = {
    client: new ApiClient('tracking'),

    getTracking: (orderId: string) =>
        trackingApi.client.get<OrderTracking>(`/api/v1/tracking/${orderId}`),

    subscribeLive: (orderId: string, onUpdate: (data: LocationUpdate) => void) => {
        const eventSource = new EventSource(`${SERVICES.tracking}/api/v1/tracking/${orderId}/live`);
        eventSource.onmessage = (event) => {
            onUpdate(JSON.parse(event.data));
        };
        return () => eventSource.close();
    },
};

// Types
export interface Auction {
    id: string;
    product_id: string;
    product_name: string;
    product_image: string;
    original_price: number;
    starting_bid: number;
    current_bid: number;
    min_increment: number;
    bid_count: number;
    start_time: string;
    end_time: string;
    status: 'pending' | 'active' | 'ending' | 'ended';
    highest_bidder?: { user_id: string; name: string; avatar_url: string };
}

export interface Bid {
    id: string;
    auction_id: string;
    user_id: string;
    amount: number;
    created_at: string;
    is_winning: boolean;
}

export interface ReferralCode {
    code: string;
    user_id: string;
    link: string;
    uses: number;
    max_uses: number;
    reward_per_use: number;
}

export interface ReferralStats {
    user_id: string;
    total_referrals: number;
    completed_count: number;
    pending_count: number;
    total_earned: number;
    pending_earnings: number;
}

export interface Referral {
    id: string;
    referrer_id: string;
    referred_id: string;
    referred_name: string;
    referred_avatar: string;
    code: string;
    status: 'pending' | 'completed' | 'expired';
    reward: number;
    joined_at: string;
    completed_at?: string;
}

export interface OrderTracking {
    order_id: string;
    status: string;
    estimated_delivery: string;
    carrier: string;
    tracking_number: string;
    current_location?: { lat: number; lng: number; address: string };
    driver?: { id: string; name: string; phone: string; vehicle: string; plate_no: string; avatar_url: string };
    steps: TrackingStep[];
}

export interface TrackingStep {
    status: string;
    description: string;
    time: string;
    location?: string;
    completed: boolean;
}

export interface LocationUpdate {
    order_id: string;
    lat: number;
    lng: number;
    eta_minutes: number;
    distance_km: number;
}
