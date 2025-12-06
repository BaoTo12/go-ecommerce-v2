// Product Service - API layer for frontend
// This connects to the backend product-service or uses mock data

export interface Product {
    id: string;
    name: string;
    description: string;
    price: number;
    originalPrice: number;
    discount: number;
    currency: string;
    category: string;
    categoryId: string;
    images: string[];
    thumbnail: string;
    rating: number;
    reviews: number;
    sold: number;
    soldDisplay: string;
    stock: number;
    location: string;
    shop: Shop;
    variants?: Variant[];
    tags?: string[];
    isOfficial?: boolean;
    isFavorite?: boolean;
    freeShip?: boolean;
    createdAt: string;
}

export interface Shop {
    id: string;
    name: string;
    avatar: string;
    rating: number;
    products: number;
    responseRate: string;
    location: string;
    isOfficial: boolean;
}

export interface Variant {
    id: string;
    name: string;
    options: string[];
    selected?: number;
}

export interface Category {
    id: string;
    name: string;
    icon: string;
    image: string;
    productCount: number;
}

// Real product data with actual images from Unsplash/Placeholder services
const PRODUCTS: Product[] = [
    {
        id: 'p1',
        name: 'iPhone 15 Pro Max 256GB Titan Xanh Chính Hãng VN/A Bảo Hành 12 Tháng',
        description: 'iPhone 15 Pro Max với chip A17 Pro mạnh mẽ nhất, camera 48MP, màn hình Super Retina XDR 6.7 inch, thời lượng pin cả ngày. Thiết kế titan cao cấp, nhẹ hơn và bền hơn.',
        price: 29990000,
        originalPrice: 34990000,
        discount: 14,
        currency: 'VND',
        category: 'Điện thoại',
        categoryId: 'phones',
        images: [
            'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=600',
            'https://images.unsplash.com/photo-1510557880182-3d4d3cba35a5?w=600',
            'https://images.unsplash.com/photo-1592750475338-74b7b21085ab?w=600',
        ],
        thumbnail: 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=300',
        rating: 4.9,
        reviews: 8560,
        sold: 12300,
        soldDisplay: '12.3k',
        stock: 156,
        location: 'TP. Hồ Chí Minh',
        shop: {
            id: 'shop1',
            name: 'Apple Store Official',
            avatar: 'https://ui-avatars.com/api/?name=Apple&background=000&color=fff',
            rating: 4.9,
            products: 156,
            responseRate: '95%',
            location: 'TP. Hồ Chí Minh',
            isOfficial: true,
        },
        variants: [
            { id: 'v1', name: 'Màu sắc', options: ['Titan Xanh', 'Titan Đen', 'Titan Trắng', 'Titan Vàng'] },
            { id: 'v2', name: 'Dung lượng', options: ['256GB', '512GB', '1TB'] },
        ],
        tags: ['Chính hãng', 'Trả góp 0%'],
        isOfficial: true,
        freeShip: true,
        createdAt: '2024-01-15',
    },
    {
        id: 'p2',
        name: 'Samsung Galaxy S24 Ultra 512GB Xám Titan Chính Hãng',
        description: 'Galaxy S24 Ultra với S Pen tích hợp, camera 200MP, màn hình Dynamic AMOLED 2X 6.8 inch, chip Snapdragon 8 Gen 3.',
        price: 25990000,
        originalPrice: 29990000,
        discount: 13,
        currency: 'VND',
        category: 'Điện thoại',
        categoryId: 'phones',
        images: [
            'https://images.unsplash.com/photo-1610945415295-d9bbf067e59c?w=600',
            'https://images.unsplash.com/photo-1585060544812-6b45742d762f?w=600',
        ],
        thumbnail: 'https://images.unsplash.com/photo-1610945415295-d9bbf067e59c?w=300',
        rating: 4.8,
        reviews: 5430,
        sold: 8700,
        soldDisplay: '8.7k',
        stock: 89,
        location: 'Hà Nội',
        shop: {
            id: 'shop2',
            name: 'Samsung Official Store',
            avatar: 'https://ui-avatars.com/api/?name=Samsung&background=1428a0&color=fff',
            rating: 4.9,
            products: 234,
            responseRate: '97%',
            location: 'Hà Nội',
            isOfficial: true,
        },
        isOfficial: true,
        freeShip: true,
        createdAt: '2024-02-10',
    },
    {
        id: 'p3',
        name: 'MacBook Air M3 13 inch 256GB Space Gray 2024',
        description: 'MacBook Air với chip M3 thế hệ mới, 8GB RAM, 256GB SSD, màn hình Liquid Retina 13.6 inch sắc nét.',
        price: 27990000,
        originalPrice: 31990000,
        discount: 12,
        currency: 'VND',
        category: 'Laptop',
        categoryId: 'laptops',
        images: [
            'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=600',
            'https://images.unsplash.com/photo-1541807084-5c52b6b3adef?w=600',
        ],
        thumbnail: 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=300',
        rating: 4.9,
        reviews: 2340,
        sold: 3200,
        soldDisplay: '3.2k',
        stock: 45,
        location: 'TP. Hồ Chí Minh',
        shop: {
            id: 'shop1',
            name: 'Apple Store Official',
            avatar: 'https://ui-avatars.com/api/?name=Apple&background=000&color=fff',
            rating: 4.9,
            products: 156,
            responseRate: '95%',
            location: 'TP. Hồ Chí Minh',
            isOfficial: true,
        },
        isOfficial: true,
        freeShip: true,
        createdAt: '2024-03-05',
    },
    {
        id: 'p4',
        name: 'Áo Hoodie Unisex Form Rộng Nỉ Cotton Dày Dặn Premium',
        description: 'Áo hoodie unisex chất liệu cotton dày dặn, form rộng thoải mái, phù hợp mọi dáng người. Nhiều màu sắc để lựa chọn.',
        price: 199000,
        originalPrice: 350000,
        discount: 43,
        currency: 'VND',
        category: 'Thời trang',
        categoryId: 'fashion',
        images: [
            'https://images.unsplash.com/photo-1556821840-3a63f95609a7?w=600',
            'https://images.unsplash.com/photo-1578681994506-b8f463449011?w=600',
        ],
        thumbnail: 'https://images.unsplash.com/photo-1556821840-3a63f95609a7?w=300',
        rating: 4.7,
        reviews: 12340,
        sold: 45200,
        soldDisplay: '45.2k',
        stock: 999,
        location: 'Hà Nội',
        shop: {
            id: 'shop3',
            name: 'Fashion Store VN',
            avatar: 'https://ui-avatars.com/api/?name=Fashion&background=e91e63&color=fff',
            rating: 4.7,
            products: 567,
            responseRate: '92%',
            location: 'Hà Nội',
            isOfficial: false,
        },
        variants: [
            { id: 'v1', name: 'Màu sắc', options: ['Đen', 'Trắng', 'Xám', 'Nâu', 'Be'] },
            { id: 'v2', name: 'Size', options: ['M', 'L', 'XL', 'XXL'] },
        ],
        isFavorite: true,
        freeShip: true,
        createdAt: '2024-01-20',
    },
    {
        id: 'p5',
        name: 'Giày Nike Air Force 1 07 Low White Chính Hãng',
        description: 'Nike Air Force 1 chính hãng, đệm Air êm ái, thiết kế iconic từ năm 1982, phù hợp mọi phong cách.',
        price: 2590000,
        originalPrice: 3200000,
        discount: 19,
        currency: 'VND',
        category: 'Giày dép',
        categoryId: 'shoes',
        images: [
            'https://images.unsplash.com/photo-1600269452121-4f2416e55c28?w=600',
            'https://images.unsplash.com/photo-1595950653106-6c9ebd614d3a?w=600',
        ],
        thumbnail: 'https://images.unsplash.com/photo-1600269452121-4f2416e55c28?w=300',
        rating: 4.8,
        reviews: 3456,
        sold: 5200,
        soldDisplay: '5.2k',
        stock: 78,
        location: 'TP. Hồ Chí Minh',
        shop: {
            id: 'shop4',
            name: 'Nike Official Store',
            avatar: 'https://ui-avatars.com/api/?name=Nike&background=000&color=fff',
            rating: 4.8,
            products: 234,
            responseRate: '96%',
            location: 'TP. Hồ Chí Minh',
            isOfficial: true,
        },
        variants: [
            { id: 'v1', name: 'Size', options: ['38', '39', '40', '41', '42', '43', '44'] },
        ],
        isOfficial: true,
        createdAt: '2024-02-28',
    },
    {
        id: 'p6',
        name: 'Son Dưỡng Môi Dior Addict Lip Glow Fullsize',
        description: 'Son dưỡng môi Dior Addict Lip Glow, dưỡng ẩm và tạo màu tự nhiên, công nghệ Color Reviver phản ứng với độ pH của môi.',
        price: 950000,
        originalPrice: 1200000,
        discount: 21,
        currency: 'VND',
        category: 'Làm đẹp',
        categoryId: 'beauty',
        images: [
            'https://images.unsplash.com/photo-1586495777744-4413f21062fa?w=600',
            'https://images.unsplash.com/photo-1596462502278-27bfdc403348?w=600',
        ],
        thumbnail: 'https://images.unsplash.com/photo-1586495777744-4413f21062fa?w=300',
        rating: 4.9,
        reviews: 8765,
        sold: 18700,
        soldDisplay: '18.7k',
        stock: 234,
        location: 'TP. Hồ Chí Minh',
        shop: {
            id: 'shop5',
            name: 'Dior Beauty Official',
            avatar: 'https://ui-avatars.com/api/?name=Dior&background=9c27b0&color=fff',
            rating: 4.9,
            products: 89,
            responseRate: '98%',
            location: 'TP. Hồ Chí Minh',
            isOfficial: true,
        },
        variants: [
            { id: 'v1', name: 'Màu', options: ['001 Pink', '004 Coral', '006 Berry', '008 Ultra Pink'] },
        ],
        isFavorite: true,
        createdAt: '2024-01-10',
    },
    {
        id: 'p7',
        name: 'Nồi Chiên Không Dầu Lock&Lock 5.2L Digital',
        description: 'Nồi chiên không dầu Lock&Lock 5.2L, bảng điều khiển cảm ứng, 8 chế độ nấu, giảm đến 80% dầu mỡ.',
        price: 1290000,
        originalPrice: 2490000,
        discount: 48,
        currency: 'VND',
        category: 'Nhà cửa',
        categoryId: 'home',
        images: [
            'https://images.unsplash.com/photo-1585515320310-259814833e62?w=600',
            'https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?w=600',
        ],
        thumbnail: 'https://images.unsplash.com/photo-1585515320310-259814833e62?w=300',
        rating: 4.8,
        reviews: 4567,
        sold: 23400,
        soldDisplay: '23.4k',
        stock: 156,
        location: 'Hà Nội',
        shop: {
            id: 'shop6',
            name: 'Lock&Lock Official',
            avatar: 'https://ui-avatars.com/api/?name=LL&background=4caf50&color=fff',
            rating: 4.8,
            products: 345,
            responseRate: '94%',
            location: 'Hà Nội',
            isOfficial: true,
        },
        isOfficial: true,
        freeShip: true,
        createdAt: '2024-02-15',
    },
    {
        id: 'p8',
        name: 'Laptop Dell XPS 13 Plus Intel Core i7 Gen 13',
        description: 'Dell XPS 13 Plus với thiết kế siêu mỏng, màn hình OLED 13.4 inch, chip Intel Core i7 thế hệ 13, 16GB RAM.',
        price: 32990000,
        originalPrice: 38990000,
        discount: 15,
        currency: 'VND',
        category: 'Laptop',
        categoryId: 'laptops',
        images: [
            'https://images.unsplash.com/photo-1593642702821-c8da6771f0c6?w=600',
            'https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=600',
        ],
        thumbnail: 'https://images.unsplash.com/photo-1593642702821-c8da6771f0c6?w=300',
        rating: 4.7,
        reviews: 1234,
        sold: 1200,
        soldDisplay: '1.2k',
        stock: 34,
        location: 'TP. Hồ Chí Minh',
        shop: {
            id: 'shop7',
            name: 'Dell Official Store',
            avatar: 'https://ui-avatars.com/api/?name=Dell&background=0076ce&color=fff',
            rating: 4.7,
            products: 178,
            responseRate: '93%',
            location: 'TP. Hồ Chí Minh',
            isOfficial: true,
        },
        isOfficial: true,
        freeShip: true,
        createdAt: '2024-03-01',
    },
    {
        id: 'p9',
        name: 'Quần Jean Nam Slim Fit Cao Cấp Dáng Ôm Vừa',
        description: 'Quần jean nam slim fit, chất liệu denim cao cấp, co giãn nhẹ, form dáng ôm vừa thoải mái.',
        price: 299000,
        originalPrice: 450000,
        discount: 34,
        currency: 'VND',
        category: 'Thời trang',
        categoryId: 'fashion',
        images: [
            'https://images.unsplash.com/photo-1542272604-787c3835535d?w=600',
            'https://images.unsplash.com/photo-1541099649105-f69ad21f3246?w=600',
        ],
        thumbnail: 'https://images.unsplash.com/photo-1542272604-787c3835535d?w=300',
        rating: 4.6,
        reviews: 6780,
        sold: 67800,
        soldDisplay: '67.8k',
        stock: 456,
        location: 'Hà Nội',
        shop: {
            id: 'shop3',
            name: 'Fashion Store VN',
            avatar: 'https://ui-avatars.com/api/?name=Fashion&background=e91e63&color=fff',
            rating: 4.7,
            products: 567,
            responseRate: '92%',
            location: 'Hà Nội',
            isOfficial: false,
        },
        variants: [
            { id: 'v1', name: 'Màu', options: ['Xanh đậm', 'Xanh nhạt', 'Đen'] },
            { id: 'v2', name: 'Size', options: ['29', '30', '31', '32', '33', '34'] },
        ],
        isFavorite: true,
        createdAt: '2024-01-25',
    },
    {
        id: 'p10',
        name: 'Serum Vitamin C The Ordinary 30ml Chính Hãng',
        description: 'Serum Vitamin C The Ordinary giúp làm sáng da, mờ thâm, chống oxy hóa. Công thức 23% + HA Spheres 2%.',
        price: 350000,
        originalPrice: 500000,
        discount: 30,
        currency: 'VND',
        category: 'Làm đẹp',
        categoryId: 'beauty',
        images: [
            'https://images.unsplash.com/photo-1620916566398-39f1143ab7be?w=600',
            'https://images.unsplash.com/photo-1608248543803-ba4f8c70ae0b?w=600',
        ],
        thumbnail: 'https://images.unsplash.com/photo-1620916566398-39f1143ab7be?w=300',
        rating: 4.8,
        reviews: 9876,
        sold: 34500,
        soldDisplay: '34.5k',
        stock: 567,
        location: 'TP. Hồ Chí Minh',
        shop: {
            id: 'shop8',
            name: 'Beauty Zone Official',
            avatar: 'https://ui-avatars.com/api/?name=BZ&background=ff4081&color=fff',
            rating: 4.8,
            products: 456,
            responseRate: '96%',
            location: 'TP. Hồ Chí Minh',
            isOfficial: true,
        },
        freeShip: true,
        createdAt: '2024-02-20',
    },
    {
        id: 'p11',
        name: 'Đèn Bàn LED Chống Cận 3 Chế Độ Sáng USB',
        description: 'Đèn bàn LED chống cận, 3 chế độ ánh sáng, điều chỉnh độ sáng, cổng sạc USB tiện lợi.',
        price: 189000,
        originalPrice: 320000,
        discount: 41,
        currency: 'VND',
        category: 'Nhà cửa',
        categoryId: 'home',
        images: [
            'https://images.unsplash.com/photo-1507473885765-e6ed057f782c?w=600',
            'https://images.unsplash.com/photo-1534105615256-13940a56ff44?w=600',
        ],
        thumbnail: 'https://images.unsplash.com/photo-1507473885765-e6ed057f782c?w=300',
        rating: 4.5,
        reviews: 3456,
        sold: 12100,
        soldDisplay: '12.1k',
        stock: 234,
        location: 'Hà Nội',
        shop: {
            id: 'shop9',
            name: 'Home Living Store',
            avatar: 'https://ui-avatars.com/api/?name=HL&background=795548&color=fff',
            rating: 4.5,
            products: 678,
            responseRate: '90%',
            location: 'Hà Nội',
            isOfficial: false,
        },
        createdAt: '2024-03-10',
    },
    {
        id: 'p12',
        name: 'Giày Adidas Ultraboost 23 Chính Hãng',
        description: 'Adidas Ultraboost 23, công nghệ Boost êm ái, upper Primeknit thoáng khí, đế Continental chống trượt.',
        price: 4290000,
        originalPrice: 4990000,
        discount: 14,
        currency: 'VND',
        category: 'Giày dép',
        categoryId: 'shoes',
        images: [
            'https://images.unsplash.com/photo-1587563871167-1ee9c731aefb?w=600',
            'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=600',
        ],
        thumbnail: 'https://images.unsplash.com/photo-1587563871167-1ee9c731aefb?w=300',
        rating: 4.9,
        reviews: 2345,
        sold: 2800,
        soldDisplay: '2.8k',
        stock: 67,
        location: 'TP. Hồ Chí Minh',
        shop: {
            id: 'shop10',
            name: 'Adidas Official Store',
            avatar: 'https://ui-avatars.com/api/?name=Adidas&background=000&color=fff',
            rating: 4.9,
            products: 289,
            responseRate: '97%',
            location: 'TP. Hồ Chí Minh',
            isOfficial: true,
        },
        variants: [
            { id: 'v1', name: 'Size', options: ['38', '39', '40', '41', '42', '43'] },
        ],
        isOfficial: true,
        freeShip: true,
        createdAt: '2024-02-25',
    },
    {
        id: 'p13',
        name: 'Túi Xách Nữ Charles & Keith Authentic',
        description: 'Túi xách nữ Charles & Keith chính hãng, thiết kế thanh lịch, chất liệu da tổng hợp cao cấp.',
        price: 890000,
        originalPrice: 1290000,
        discount: 31,
        currency: 'VND',
        category: 'Túi ví',
        categoryId: 'bags',
        images: [
            'https://images.unsplash.com/photo-1584917865442-de89df76afd3?w=600',
            'https://images.unsplash.com/photo-1548036328-c9fa89d128fa?w=600',
        ],
        thumbnail: 'https://images.unsplash.com/photo-1584917865442-de89df76afd3?w=300',
        rating: 4.7,
        reviews: 4567,
        sold: 9100,
        soldDisplay: '9.1k',
        stock: 89,
        location: 'Hà Nội',
        shop: {
            id: 'shop11',
            name: 'Charles & Keith VN',
            avatar: 'https://ui-avatars.com/api/?name=CK&background=212121&color=fff',
            rating: 4.7,
            products: 234,
            responseRate: '95%',
            location: 'Hà Nội',
            isOfficial: true,
        },
        isFavorite: true,
        createdAt: '2024-01-30',
    },
    {
        id: 'p14',
        name: 'Đồng Hồ Casio G-Shock GA-2100 Chính Hãng',
        description: 'Casio G-Shock GA-2100 chính hãng, thiết kế mỏng, chống nước 200m, đèn LED, dây nhựa bền bỉ.',
        price: 2890000,
        originalPrice: 3500000,
        discount: 17,
        currency: 'VND',
        category: 'Đồng hồ',
        categoryId: 'watches',
        images: [
            'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=600',
            'https://images.unsplash.com/photo-1546868871-7041f2a55e12?w=600',
        ],
        thumbnail: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=300',
        rating: 4.8,
        reviews: 3456,
        sold: 4500,
        soldDisplay: '4.5k',
        stock: 56,
        location: 'TP. Hồ Chí Minh',
        shop: {
            id: 'shop12',
            name: 'Casio Official Store',
            avatar: 'https://ui-avatars.com/api/?name=Casio&background=1565c0&color=fff',
            rating: 4.8,
            products: 189,
            responseRate: '96%',
            location: 'TP. Hồ Chí Minh',
            isOfficial: true,
        },
        isOfficial: true,
        createdAt: '2024-02-18',
    },
    {
        id: 'p15',
        name: 'Tai Nghe Bluetooth Apple AirPods Pro 2 USB-C',
        description: 'AirPods Pro 2 với chip H2, chống ồn chủ động, âm thanh không gian, USB-C tiện lợi, pin 30h với case.',
        price: 4990000,
        originalPrice: 6990000,
        discount: 29,
        currency: 'VND',
        category: 'Điện thoại',
        categoryId: 'phones',
        images: [
            'https://images.unsplash.com/photo-1600294037681-c80b4cb5b434?w=600',
            'https://images.unsplash.com/photo-1606220945770-b5b6c2c55bf1?w=600',
        ],
        thumbnail: 'https://images.unsplash.com/photo-1600294037681-c80b4cb5b434?w=300',
        rating: 4.9,
        reviews: 7654,
        sold: 15100,
        soldDisplay: '15.1k',
        stock: 123,
        location: 'TP. Hồ Chí Minh',
        shop: {
            id: 'shop1',
            name: 'Apple Store Official',
            avatar: 'https://ui-avatars.com/api/?name=Apple&background=000&color=fff',
            rating: 4.9,
            products: 156,
            responseRate: '95%',
            location: 'TP. Hồ Chí Minh',
            isOfficial: true,
        },
        isOfficial: true,
        freeShip: true,
        createdAt: '2024-03-08',
    },
];

const CATEGORIES: Category[] = [
    { id: 'phones', name: 'Điện thoại', icon: '📱', image: 'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=100', productCount: 156 },
    { id: 'laptops', name: 'Laptop', icon: '💻', image: 'https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=100', productCount: 89 },
    { id: 'fashion', name: 'Thời trang', icon: '👕', image: 'https://images.unsplash.com/photo-1445205170230-053b83016050?w=100', productCount: 345 },
    { id: 'beauty', name: 'Làm đẹp', icon: '💄', image: 'https://images.unsplash.com/photo-1596462502278-27bfdc403348?w=100', productCount: 234 },
    { id: 'home', name: 'Nhà cửa', icon: '🏠', image: 'https://images.unsplash.com/photo-1484101403633-562f891dc89a?w=100', productCount: 178 },
    { id: 'shoes', name: 'Giày dép', icon: '👟', image: 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=100', productCount: 123 },
    { id: 'bags', name: 'Túi ví', icon: '👜', image: 'https://images.unsplash.com/photo-1548036328-c9fa89d128fa?w=100', productCount: 98 },
    { id: 'watches', name: 'Đồng hồ', icon: '⌚', image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=100', productCount: 67 },
    { id: 'electronics', name: 'Điện tử', icon: '📺', image: 'https://images.unsplash.com/photo-1468495244123-6c6c332eeece?w=100', productCount: 145 },
    { id: 'sports', name: 'Thể thao', icon: '🏃', image: 'https://images.unsplash.com/photo-1461896836934- voices-0-of-faith?w=100', productCount: 89 },
];

// Product Service API
export const productService = {
    // Get all products
    getProducts: async (params?: {
        category?: string;
        search?: string;
        sort?: string;
        page?: number;
        limit?: number;
    }): Promise<{ products: Product[]; total: number }> => {
        // Simulate API delay
        await new Promise(resolve => setTimeout(resolve, 100));

        let result = [...PRODUCTS];

        // Filter by category
        if (params?.category && params.category !== 'Tất cả') {
            result = result.filter(p => p.category === params.category);
        }

        // Search
        if (params?.search) {
            const query = params.search.toLowerCase();
            result = result.filter(p =>
                p.name.toLowerCase().includes(query) ||
                p.description.toLowerCase().includes(query)
            );
        }

        // Sort
        switch (params?.sort) {
            case 'price-asc':
                result.sort((a, b) => a.price - b.price);
                break;
            case 'price-desc':
                result.sort((a, b) => b.price - a.price);
                break;
            case 'newest':
                result.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
                break;
            case 'best-seller':
                result.sort((a, b) => b.sold - a.sold);
                break;
            default:
                result.sort((a, b) => b.rating - a.rating);
        }

        return { products: result, total: result.length };
    },

    // Get product by ID
    getProduct: async (id: string): Promise<Product | null> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        return PRODUCTS.find(p => p.id === id) || null;
    },

    // Get categories
    getCategories: async (): Promise<Category[]> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        return CATEGORIES;
    },

    // Get flash sale products
    getFlashSaleProducts: async (): Promise<Product[]> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        return PRODUCTS.filter(p => p.discount >= 20).slice(0, 6);
    },

    // Get recommended products
    getRecommendedProducts: async (): Promise<Product[]> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        return PRODUCTS.sort((a, b) => b.rating - a.rating).slice(0, 12);
    },
};

export default productService;
