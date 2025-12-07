// Shopee Coins Service - Loyalty points system with gamification

export interface CoinTransaction {
    id: string;
    type: 'earn' | 'spend' | 'expire';
    amount: number;
    description: string;
    orderId?: string;
    timestamp: string;
}

export interface CoinState {
    balance: number;
    lifetime: number;
    expiring: { amount: number; date: string } | null;
    transactions: CoinTransaction[];
    dailyCheckIn: {
        lastCheckIn: string | null;
        streak: number;
        todayReward: number;
    };
}

// Get from localStorage
const getStoredCoins = (): CoinState => {
    if (typeof window === 'undefined') {
        return {
            balance: 0,
            lifetime: 0,
            expiring: null,
            transactions: [],
            dailyCheckIn: { lastCheckIn: null, streak: 0, todayReward: 5 }
        };
    }

    const saved = localStorage.getItem('shopeeCoins');
    if (saved) {
        try {
            return JSON.parse(saved);
        } catch (e) {
            console.error('Failed to parse coins:', e);
        }
    }

    // Default starting bonus
    return {
        balance: 500,
        lifetime: 500,
        expiring: { amount: 100, date: '2024-12-31' },
        transactions: [
            {
                id: 't1',
                type: 'earn',
                amount: 500,
                description: 'Thưởng chào mừng thành viên mới',
                timestamp: new Date().toISOString(),
            }
        ],
        dailyCheckIn: {
            lastCheckIn: null,
            streak: 0,
            todayReward: 5,
        }
    };
};

let coinState = getStoredCoins();

const saveCoins = () => {
    if (typeof window !== 'undefined') {
        localStorage.setItem('shopeeCoins', JSON.stringify(coinState));
    }
};

export const coinsService = {
    // Get coin balance
    getBalance: (): number => coinState.balance,

    // Get full state
    getState: async (): Promise<CoinState> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        coinState = getStoredCoins();
        return coinState;
    },

    // Earn coins
    earnCoins: async (amount: number, description: string, orderId?: string): Promise<CoinTransaction> => {
        await new Promise(resolve => setTimeout(resolve, 100));

        const transaction: CoinTransaction = {
            id: `coin_${Date.now()}`,
            type: 'earn',
            amount,
            description,
            orderId,
            timestamp: new Date().toISOString(),
        };

        coinState.balance += amount;
        coinState.lifetime += amount;
        coinState.transactions.unshift(transaction);
        saveCoins();

        return transaction;
    },

    // Spend coins
    spendCoins: async (amount: number, description: string, orderId?: string): Promise<{ success: boolean; error?: string }> => {
        await new Promise(resolve => setTimeout(resolve, 100));

        if (coinState.balance < amount) {
            return { success: false, error: 'Không đủ xu' };
        }

        const transaction: CoinTransaction = {
            id: `coin_${Date.now()}`,
            type: 'spend',
            amount: -amount,
            description,
            orderId,
            timestamp: new Date().toISOString(),
        };

        coinState.balance -= amount;
        coinState.transactions.unshift(transaction);
        saveCoins();

        return { success: true };
    },

    // Daily check-in
    dailyCheckIn: async (): Promise<{ success: boolean; reward: number; streak: number; error?: string }> => {
        await new Promise(resolve => setTimeout(resolve, 200));

        const today = new Date().toDateString();
        const lastCheckIn = coinState.dailyCheckIn.lastCheckIn;

        if (lastCheckIn === today) {
            return { success: false, reward: 0, streak: coinState.dailyCheckIn.streak, error: 'Bạn đã điểm danh hôm nay rồi!' };
        }

        // Check if yesterday (for streak)
        const yesterday = new Date();
        yesterday.setDate(yesterday.getDate() - 1);
        const wasYesterday = lastCheckIn === yesterday.toDateString();

        const newStreak = wasYesterday ? coinState.dailyCheckIn.streak + 1 : 1;
        const reward = Math.min(5 + (newStreak - 1) * 2, 50); // 5, 7, 9, 11... up to 50

        coinState.dailyCheckIn = {
            lastCheckIn: today,
            streak: newStreak,
            todayReward: reward,
        };

        await coinsService.earnCoins(reward, `Điểm danh ngày ${newStreak}`);

        return { success: true, reward, streak: newStreak };
    },

    // Check if already checked in today
    hasCheckedInToday: (): boolean => {
        const today = new Date().toDateString();
        return coinState.dailyCheckIn.lastCheckIn === today;
    },

    // Get streak
    getStreak: (): number => coinState.dailyCheckIn.streak,

    // Calculate coins from order (1% of order value)
    calculateOrderCoins: (orderTotal: number): number => {
        return Math.floor(orderTotal / 10000); // 1 coin per 10,000 VND
    },

    // Spin the wheel
    spinWheel: async (): Promise<{ reward: number; type: 'coins' | 'voucher' | 'nothing' }> => {
        await new Promise(resolve => setTimeout(resolve, 3000)); // Spinning animation

        const random = Math.random();
        if (random < 0.3) {
            // 30% chance: coins
            const coins = [10, 20, 50, 100][Math.floor(Math.random() * 4)];
            await coinsService.earnCoins(coins, 'Vòng quay may mắn');
            return { reward: coins, type: 'coins' };
        } else if (random < 0.4) {
            // 10% chance: voucher
            return { reward: 10000, type: 'voucher' };
        } else {
            // 60% chance: nothing
            return { reward: 0, type: 'nothing' };
        }
    },
};

export default coinsService;
