// Voucher Service - Coupon codes and discounts

export interface Voucher {
    id: string;
    code: string;
    type: 'percentage' | 'fixed' | 'freeship';
    value: number;
    maxDiscount?: number;
    minOrder: number;
    description: string;
    expiresAt: string;
    usageLimit: number;
    usedCount: number;
    isCollected: boolean;
    category?: string;
    shopId?: string;
    shopName?: string;
}

// Mock vouchers
const ALL_VOUCHERS: Voucher[] = [
    {
        id: 'v1',
        code: 'GIAM50K',
        type: 'fixed',
        value: 50000,
        minOrder: 200000,
        description: 'Giảm ₫50.000 cho đơn từ ₫200.000',
        expiresAt: '2024-12-31',
        usageLimit: 1000,
        usedCount: 456,
        isCollected: true,
    },
    {
        id: 'v2',
        code: 'FREESHIP',
        type: 'freeship',
        value: 30000,
        minOrder: 0,
        description: 'Miễn phí vận chuyển đến ₫30.000',
        expiresAt: '2024-12-31',
        usageLimit: 5000,
        usedCount: 2341,
        isCollected: true,
    },
    {
        id: 'v3',
        code: 'NEWUSER',
        type: 'percentage',
        value: 15,
        maxDiscount: 100000,
        minOrder: 100000,
        description: 'Giảm 15% tối đa ₫100.000 cho khách mới',
        expiresAt: '2024-12-15',
        usageLimit: 500,
        usedCount: 123,
        isCollected: false,
    },
    {
        id: 'v4',
        code: 'APPLE20',
        type: 'percentage',
        value: 5,
        maxDiscount: 500000,
        minOrder: 5000000,
        description: 'Giảm 5% đơn từ ₫5.000.000',
        expiresAt: '2024-12-31',
        usageLimit: 100,
        usedCount: 45,
        isCollected: false,
        shopId: 'shop1',
        shopName: 'Apple Store Official',
    },
    {
        id: 'v5',
        code: 'BEAUTY30',
        type: 'fixed',
        value: 30000,
        minOrder: 150000,
        description: 'Giảm ₫30.000 cho đơn Làm đẹp từ ₫150.000',
        expiresAt: '2024-12-25',
        usageLimit: 2000,
        usedCount: 890,
        isCollected: false,
        category: 'beauty',
    },
    {
        id: 'v6',
        code: 'FLASH1212',
        type: 'percentage',
        value: 20,
        maxDiscount: 200000,
        minOrder: 300000,
        description: 'Flash Sale 12.12 - Giảm 20%',
        expiresAt: '2024-12-12',
        usageLimit: 10000,
        usedCount: 4567,
        isCollected: true,
    },
];

let userVouchers = ALL_VOUCHERS.filter(v => v.isCollected);

export const voucherService = {
    // Get all available vouchers
    getAllVouchers: async (): Promise<Voucher[]> => {
        await new Promise(resolve => setTimeout(resolve, 100));
        return ALL_VOUCHERS;
    },

    // Get user's collected vouchers
    getUserVouchers: async (): Promise<Voucher[]> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        return userVouchers;
    },

    // Collect a voucher
    collectVoucher: async (voucherId: string): Promise<boolean> => {
        await new Promise(resolve => setTimeout(resolve, 100));

        const voucher = ALL_VOUCHERS.find(v => v.id === voucherId);
        if (!voucher) return false;
        if (voucher.isCollected) return false;
        if (voucher.usedCount >= voucher.usageLimit) return false;

        voucher.isCollected = true;
        userVouchers = ALL_VOUCHERS.filter(v => v.isCollected);
        return true;
    },

    // Apply voucher to order
    applyVoucher: async (code: string, orderTotal: number): Promise<{
        success: boolean;
        discount: number;
        error?: string;
        voucher?: Voucher;
    }> => {
        await new Promise(resolve => setTimeout(resolve, 100));

        const voucher = userVouchers.find(v => v.code === code);
        if (!voucher) {
            return { success: false, discount: 0, error: 'Mã voucher không tồn tại hoặc chưa được lưu' };
        }

        if (new Date(voucher.expiresAt) < new Date()) {
            return { success: false, discount: 0, error: 'Voucher đã hết hạn' };
        }

        if (orderTotal < voucher.minOrder) {
            return {
                success: false,
                discount: 0,
                error: `Đơn hàng tối thiểu ₫${voucher.minOrder.toLocaleString('vi-VN')}`
            };
        }

        let discount = 0;
        if (voucher.type === 'fixed') {
            discount = voucher.value;
        } else if (voucher.type === 'percentage') {
            discount = Math.floor(orderTotal * voucher.value / 100);
            if (voucher.maxDiscount) {
                discount = Math.min(discount, voucher.maxDiscount);
            }
        } else if (voucher.type === 'freeship') {
            discount = voucher.value;
        }

        return { success: true, discount, voucher };
    },

    // Get vouchers by category
    getVouchersByCategory: async (category: string): Promise<Voucher[]> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        return ALL_VOUCHERS.filter(v => v.category === category);
    },

    // Get shop vouchers
    getShopVouchers: async (shopId: string): Promise<Voucher[]> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        return ALL_VOUCHERS.filter(v => v.shopId === shopId);
    },
};

export default voucherService;
