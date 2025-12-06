// User Service - API layer for user accounts and orders

export interface User {
    id: string;
    name: string;
    email: string;
    phone: string;
    avatar: string;
    joinDate: string;
    addresses: Address[];
    defaultAddressId: string;
}

export interface Address {
    id: string;
    name: string;
    phone: string;
    address: string;
    district: string;
    city: string;
    isDefault: boolean;
    type: 'home' | 'office';
}

export interface Order {
    id: string;
    orderNumber: string;
    status: OrderStatus;
    items: OrderItem[];
    subtotal: number;
    shippingFee: number;
    discount: number;
    total: number;
    paymentMethod: string;
    shippingMethod: string;
    shippingAddress: Address;
    trackingNumber?: string;
    createdAt: string;
    updatedAt: string;
    estimatedDelivery?: string;
    deliveredAt?: string;
}

export type OrderStatus =
    | 'pending_payment'
    | 'processing'
    | 'shipped'
    | 'out_for_delivery'
    | 'delivered'
    | 'cancelled'
    | 'refunded';

export interface OrderItem {
    id: string;
    productId: string;
    name: string;
    image: string;
    variant?: string;
    price: number;
    quantity: number;
    shopName: string;
}

export interface ShippingUpdate {
    id: string;
    timestamp: string;
    location: string;
    status: string;
    description: string;
}

// Mock data
const CURRENT_USER: User = {
    id: 'u1',
    name: 'Nguyễn Văn A',
    email: 'nguyenvana@email.com',
    phone: '0901234567',
    avatar: 'https://ui-avatars.com/api/?name=Nguyen+Van+A&background=ee4d2d&color=fff&size=200',
    joinDate: '2023-06-15',
    addresses: [
        {
            id: 'addr1',
            name: 'Nguyễn Văn A',
            phone: '0901234567',
            address: '123 Đường ABC, Phường XYZ',
            district: 'Quận 1',
            city: 'TP. Hồ Chí Minh',
            isDefault: true,
            type: 'home',
        },
        {
            id: 'addr2',
            name: 'Nguyễn Văn A',
            phone: '0901234567',
            address: '456 Đường DEF, Tầng 5',
            district: 'Quận 3',
            city: 'TP. Hồ Chí Minh',
            isDefault: false,
            type: 'office',
        },
    ],
    defaultAddressId: 'addr1',
};

const ORDERS: Order[] = [
    {
        id: 'ord1',
        orderNumber: 'SP241206001234',
        status: 'out_for_delivery',
        items: [
            { id: 'i1', productId: 'p1', name: 'iPhone 15 Pro Max 256GB Titan Xanh', image: 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=100', variant: 'Xanh Titan, 256GB', price: 29990000, quantity: 1, shopName: 'Apple Store Official' },
        ],
        subtotal: 29990000,
        shippingFee: 0,
        discount: 500000,
        total: 29490000,
        paymentMethod: 'Thanh toán khi nhận hàng',
        shippingMethod: 'Giao Hàng Nhanh',
        shippingAddress: CURRENT_USER.addresses[0],
        trackingNumber: 'GHN123456789',
        createdAt: '2024-12-04T10:30:00',
        updatedAt: '2024-12-06T08:15:00',
        estimatedDelivery: '2024-12-06',
    },
    {
        id: 'ord2',
        orderNumber: 'SP241205005678',
        status: 'delivered',
        items: [
            { id: 'i2', productId: 'p6', name: 'Son Dưỡng Môi Dior Addict Lip Glow', image: 'https://images.unsplash.com/photo-1586495777744-4413f21062fa?w=100', variant: '001 Pink', price: 950000, quantity: 2, shopName: 'Dior Beauty Official' },
            { id: 'i3', productId: 'p10', name: 'Serum Vitamin C The Ordinary 30ml', image: 'https://images.unsplash.com/photo-1620916566398-39f1143ab7be?w=100', price: 350000, quantity: 1, shopName: 'Beauty Zone Official' },
        ],
        subtotal: 2250000,
        shippingFee: 0,
        discount: 50000,
        total: 2200000,
        paymentMethod: 'Ví MoMo',
        shippingMethod: 'Giao Hàng Tiết Kiệm',
        shippingAddress: CURRENT_USER.addresses[0],
        trackingNumber: 'GHTK987654321',
        createdAt: '2024-12-01T14:20:00',
        updatedAt: '2024-12-03T16:45:00',
        deliveredAt: '2024-12-03T16:45:00',
    },
    {
        id: 'ord3',
        orderNumber: 'SP241130009012',
        status: 'processing',
        items: [
            { id: 'i4', productId: 'p4', name: 'Áo Hoodie Unisex Form Rộng Nỉ Cotton', image: 'https://images.unsplash.com/photo-1556821840-3a63f95609a7?w=100', variant: 'Đen, L', price: 199000, quantity: 2, shopName: 'Fashion Store VN' },
        ],
        subtotal: 398000,
        shippingFee: 25000,
        discount: 0,
        total: 423000,
        paymentMethod: 'VNPay QR',
        shippingMethod: 'Giao Hàng Nhanh',
        shippingAddress: CURRENT_USER.addresses[1],
        createdAt: '2024-12-06T09:00:00',
        updatedAt: '2024-12-06T09:30:00',
        estimatedDelivery: '2024-12-08',
    },
    {
        id: 'ord4',
        orderNumber: 'SP241128003456',
        status: 'shipped',
        items: [
            { id: 'i5', productId: 'p5', name: 'Giày Nike Air Force 1 07 Low White', image: 'https://images.unsplash.com/photo-1600269452121-4f2416e55c28?w=100', variant: 'Size 42', price: 2590000, quantity: 1, shopName: 'Nike Official Store' },
        ],
        subtotal: 2590000,
        shippingFee: 0,
        discount: 100000,
        total: 2490000,
        paymentMethod: 'Thẻ tín dụng',
        shippingMethod: 'Giao Hàng Nhanh',
        shippingAddress: CURRENT_USER.addresses[0],
        trackingNumber: 'GHN456789123',
        createdAt: '2024-12-05T11:15:00',
        updatedAt: '2024-12-06T07:00:00',
        estimatedDelivery: '2024-12-07',
    },
];

const SHIPPING_UPDATES: Record<string, ShippingUpdate[]> = {
    'ord1': [
        { id: 's1', timestamp: '2024-12-06T08:15:00', location: 'Quận 1, TP.HCM', status: 'Đang giao hàng', description: 'Shipper đang trên đường giao hàng cho bạn' },
        { id: 's2', timestamp: '2024-12-06T07:00:00', location: 'Bưu cục Quận 1', status: 'Đã đến bưu cục', description: 'Đơn hàng đã đến bưu cục gần bạn' },
        { id: 's3', timestamp: '2024-12-05T18:30:00', location: 'Trung tâm phân loại HCM', status: 'Đang vận chuyển', description: 'Đơn hàng đang được vận chuyển đến bưu cục' },
        { id: 's4', timestamp: '2024-12-05T14:00:00', location: 'Kho hàng Bình Dương', status: 'Đã lấy hàng', description: 'Đơn hàng đã được lấy từ người bán' },
        { id: 's5', timestamp: '2024-12-04T10:30:00', location: '', status: 'Đơn hàng đã đặt', description: 'Đơn hàng của bạn đã được đặt thành công' },
    ],
    'ord4': [
        { id: 's1', timestamp: '2024-12-06T07:00:00', location: 'Trung tâm phân loại HCM', status: 'Đang vận chuyển', description: 'Đơn hàng đang trên đường đến địa chỉ của bạn' },
        { id: 's2', timestamp: '2024-12-05T20:00:00', location: 'Kho hàng Hà Nội', status: 'Đã lấy hàng', description: 'Đơn hàng đã được lấy từ kho Nike' },
        { id: 's3', timestamp: '2024-12-05T11:15:00', location: '', status: 'Đơn hàng đã đặt', description: 'Đơn hàng đã được đặt thành công' },
    ],
};

// User Service API
export const userService = {
    // Get current user
    getCurrentUser: async (): Promise<User> => {
        await new Promise(resolve => setTimeout(resolve, 100));
        return CURRENT_USER;
    },

    // Update user profile
    updateProfile: async (data: Partial<User>): Promise<User> => {
        await new Promise(resolve => setTimeout(resolve, 200));
        Object.assign(CURRENT_USER, data);
        return CURRENT_USER;
    },

    // Get user addresses
    getAddresses: async (): Promise<Address[]> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        return CURRENT_USER.addresses;
    },

    // Add address
    addAddress: async (address: Omit<Address, 'id'>): Promise<Address> => {
        await new Promise(resolve => setTimeout(resolve, 100));
        const newAddress = { ...address, id: `addr${Date.now()}` };
        CURRENT_USER.addresses.push(newAddress);
        return newAddress;
    },

    // Get orders
    getOrders: async (status?: OrderStatus): Promise<Order[]> => {
        await new Promise(resolve => setTimeout(resolve, 100));
        if (status) {
            return ORDERS.filter(o => o.status === status);
        }
        return ORDERS;
    },

    // Get order by ID
    getOrder: async (orderId: string): Promise<Order | null> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        return ORDERS.find(o => o.id === orderId) || null;
    },

    // Get shipping updates
    getShippingUpdates: async (orderId: string): Promise<ShippingUpdate[]> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        return SHIPPING_UPDATES[orderId] || [];
    },

    // Cancel order
    cancelOrder: async (orderId: string): Promise<boolean> => {
        await new Promise(resolve => setTimeout(resolve, 200));
        const order = ORDERS.find(o => o.id === orderId);
        if (order && ['pending_payment', 'processing'].includes(order.status)) {
            order.status = 'cancelled';
            return true;
        }
        return false;
    },
};

export default userService;
