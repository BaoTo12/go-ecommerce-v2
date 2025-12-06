'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';

interface Product {
    id: string;
    name: string;
    price: number;
    originalPrice: number;
    sold: string;
    rating: number;
    image: string;
    category: string;
    location: string;
    isOfficial?: boolean;
    isFavorite?: boolean;
    freeShip?: boolean;
}

export default function ProductsPage() {
    const [products, setProducts] = useState<Product[]>([]);
    const [filteredProducts, setFilteredProducts] = useState<Product[]>([]);
    const [selectedCategory, setSelectedCategory] = useState('Tất cả');
    const [sortBy, setSortBy] = useState('popular');
    const [searchQuery, setSearchQuery] = useState('');
    const [priceRange, setPriceRange] = useState({ min: '', max: '' });
    const [notification, setNotification] = useState<string | null>(null);

    const categories = ['Tất cả', 'Điện thoại', 'Laptop', 'Thời trang', 'Làm đẹp', 'Nhà cửa', 'Giày dép', 'Túi ví', 'Đồng hồ'];

    useEffect(() => {
        setProducts([
            { id: 'p1', name: 'iPhone 15 Pro Max 256GB Titan Xanh Chính Hãng VN/A Bảo Hành 12 Tháng', price: 29990000, originalPrice: 34990000, sold: '12.3k', rating: 4.9, image: '📱', category: 'Điện thoại', location: 'TP. Hồ Chí Minh', isOfficial: true, freeShip: true },
            { id: 'p2', name: 'Samsung Galaxy S24 Ultra 512GB Xám Titan Chính Hãng', price: 25990000, originalPrice: 29990000, sold: '8.7k', rating: 4.8, image: '📲', category: 'Điện thoại', location: 'Hà Nội', isOfficial: true, freeShip: true },
            { id: 'p3', name: 'MacBook Air M3 13 inch 256GB Space Gray 2024', price: 27990000, originalPrice: 31990000, sold: '3.2k', rating: 4.9, image: '💻', category: 'Laptop', location: 'TP. Hồ Chí Minh', isOfficial: true, freeShip: true },
            { id: 'p4', name: 'Áo Hoodie Unisex Form Rộng Nỉ Cotton Dày Dặn Premium', price: 199000, originalPrice: 350000, sold: '45.2k', rating: 4.7, image: '👕', category: 'Thời trang', location: 'Hà Nội', isFavorite: true, freeShip: true },
            { id: 'p5', name: 'Giày Nike Air Force 1 07 Low White Chính Hãng', price: 2590000, originalPrice: 3200000, sold: '5.2k', rating: 4.8, image: '👟', category: 'Giày dép', location: 'TP. Hồ Chí Minh', isOfficial: true },
            { id: 'p6', name: 'Son Dưỡng Môi Dior Addict Lip Glow Fullsize', price: 950000, originalPrice: 1200000, sold: '18.7k', rating: 4.9, image: '💄', category: 'Làm đẹp', location: 'TP. Hồ Chí Minh', isFavorite: true },
            { id: 'p7', name: 'Nồi Chiên Không Dầu Lock&Lock 5.2L Digital', price: 1290000, originalPrice: 2490000, sold: '23.4k', rating: 4.8, image: '🍳', category: 'Nhà cửa', location: 'Hà Nội', freeShip: true },
            { id: 'p8', name: 'Laptop Dell XPS 13 Plus Intel Core i7 Gen 13', price: 32990000, originalPrice: 38990000, sold: '1.2k', rating: 4.7, image: '💻', category: 'Laptop', location: 'TP. Hồ Chí Minh', isOfficial: true },
            { id: 'p9', name: 'Quần Jean Nam Slim Fit Cao Cấp Dáng Ôm Vừa', price: 299000, originalPrice: 450000, sold: '67.8k', rating: 4.6, image: '👖', category: 'Thời trang', location: 'Hà Nội', isFavorite: true },
            { id: 'p10', name: 'Serum Vitamin C The Ordinary 30ml Chính Hãng', price: 350000, originalPrice: 500000, sold: '34.5k', rating: 4.8, image: '🧴', category: 'Làm đẹp', location: 'TP. Hồ Chí Minh', freeShip: true },
            { id: 'p11', name: 'Đèn Bàn LED Chống Cận 3 Chế Độ Sáng USB', price: 189000, originalPrice: 320000, sold: '12.1k', rating: 4.5, image: '💡', category: 'Nhà cửa', location: 'Hà Nội' },
            { id: 'p12', name: 'Giày Adidas Ultraboost 23 Chính Hãng', price: 4290000, originalPrice: 4990000, sold: '2.8k', rating: 4.9, image: '👟', category: 'Giày dép', location: 'TP. Hồ Chí Minh', isOfficial: true, freeShip: true },
            { id: 'p13', name: 'Túi Xách Nữ Charles & Keith Authentic', price: 890000, originalPrice: 1290000, sold: '9.1k', rating: 4.7, image: '👜', category: 'Túi ví', location: 'Hà Nội', isFavorite: true },
            { id: 'p14', name: 'Đồng Hồ Casio G-Shock GA-2100 Chính Hãng', price: 2890000, originalPrice: 3500000, sold: '4.5k', rating: 4.8, image: '⌚', category: 'Đồng hồ', location: 'TP. Hồ Chí Minh', isOfficial: true },
            { id: 'p15', name: 'Tai Nghe Bluetooth Apple AirPods Pro 2 USB-C', price: 4990000, originalPrice: 6990000, sold: '15.1k', rating: 4.9, image: '🎧', category: 'Điện thoại', location: 'TP. Hồ Chí Minh', isOfficial: true, freeShip: true },
        ]);
    }, []);

    useEffect(() => {
        let result = [...products];

        if (selectedCategory !== 'Tất cả') {
            result = result.filter(p => p.category === selectedCategory);
        }

        if (searchQuery) {
            result = result.filter(p =>
                p.name.toLowerCase().includes(searchQuery.toLowerCase())
            );
        }

        switch (sortBy) {
            case 'price-asc':
                result.sort((a, b) => a.price - b.price);
                break;
            case 'price-desc':
                result.sort((a, b) => b.price - a.price);
                break;
            case 'newest':
                result.sort((a, b) => parseFloat(b.sold) - parseFloat(a.sold));
                break;
            case 'best-seller':
                result.sort((a, b) => parseFloat(b.sold.replace('k', '000')) - parseFloat(a.sold.replace('k', '000')));
                break;
            default:
                result.sort((a, b) => b.rating - a.rating);
        }

        setFilteredProducts(result);
    }, [products, selectedCategory, sortBy, searchQuery]);

    const addToCart = (productName: string, e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setNotification(`Đã thêm "${productName.substring(0, 30)}..." vào giỏ hàng`);
        setTimeout(() => setNotification(null), 2000);
    };

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);
    const getDiscount = (price: number, original: number) => Math.round((1 - price / original) * 100);

    return (
        <div className="min-h-screen bg-[#f5f5f5]">
            {/* Toast */}
            {notification && (
                <div className="toast toast-success">{notification}</div>
            )}

            <div className="container mx-auto px-4 py-4">
                <div className="flex gap-4">
                    {/* Sidebar Filters */}
                    <aside className="w-[190px] flex-shrink-0 hidden lg:block">
                        <div className="bg-white rounded-sm shadow-sm p-4">
                            <h3 className="font-bold text-sm mb-3 flex items-center gap-2">
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h7" />
                                </svg>
                                Bộ Lọc Tìm Kiếm
                            </h3>

                            {/* Categories */}
                            <div className="border-b pb-4 mb-4">
                                <h4 className="text-sm font-medium mb-2">Theo Danh Mục</h4>
                                <div className="space-y-2">
                                    {categories.map(cat => (
                                        <button
                                            key={cat}
                                            onClick={() => setSelectedCategory(cat)}
                                            className={`block w-full text-left text-sm py-1 px-2 rounded transition-colors ${selectedCategory === cat
                                                    ? 'text-[#ee4d2d] bg-[#fef6f5]'
                                                    : 'text-gray-600 hover:text-[#ee4d2d]'
                                                }`}
                                        >
                                            {cat}
                                        </button>
                                    ))}
                                </div>
                            </div>

                            {/* Price Range */}
                            <div className="border-b pb-4 mb-4">
                                <h4 className="text-sm font-medium mb-2">Khoảng Giá</h4>
                                <div className="flex items-center gap-2">
                                    <input
                                        type="text"
                                        placeholder="₫ TỪ"
                                        value={priceRange.min}
                                        onChange={(e) => setPriceRange({ ...priceRange, min: e.target.value })}
                                        className="w-full border px-2 py-1 text-sm outline-none focus:border-[#ee4d2d]"
                                    />
                                    <span className="text-gray-400">-</span>
                                    <input
                                        type="text"
                                        placeholder="₫ ĐẾN"
                                        value={priceRange.max}
                                        onChange={(e) => setPriceRange({ ...priceRange, max: e.target.value })}
                                        className="w-full border px-2 py-1 text-sm outline-none focus:border-[#ee4d2d]"
                                    />
                                </div>
                                <button className="w-full mt-2 py-1 bg-[#ee4d2d] text-white text-sm hover:opacity-90">
                                    ÁP DỤNG
                                </button>
                            </div>

                            {/* Rating Filter */}
                            <div>
                                <h4 className="text-sm font-medium mb-2">Đánh Giá</h4>
                                {[5, 4, 3].map(stars => (
                                    <button key={stars} className="flex items-center gap-1 py-1 text-sm text-gray-600 hover:text-[#ee4d2d]">
                                        {[...Array(5)].map((_, i) => (
                                            <span key={i} className={i < stars ? 'star-rating' : 'text-gray-300'}>★</span>
                                        ))}
                                        <span>trở lên</span>
                                    </button>
                                ))}
                            </div>
                        </div>
                    </aside>

                    {/* Main Content */}
                    <main className="flex-1">
                        {/* Search Bar */}
                        <div className="bg-white rounded-sm shadow-sm p-3 mb-3">
                            <input
                                type="text"
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                placeholder="Tìm sản phẩm trong danh mục này..."
                                className="w-full border px-3 py-2 text-sm outline-none focus:border-[#ee4d2d] rounded-sm"
                            />
                        </div>

                        {/* Sort Bar */}
                        <div className="bg-[#ededed] rounded-sm p-3 mb-3 flex items-center gap-2">
                            <span className="text-gray-500 text-sm">Sắp xếp theo</span>
                            {[
                                { value: 'popular', label: 'Phổ Biến' },
                                { value: 'newest', label: 'Mới Nhất' },
                                { value: 'best-seller', label: 'Bán Chạy' },
                            ].map(opt => (
                                <button
                                    key={opt.value}
                                    onClick={() => setSortBy(opt.value)}
                                    className={`px-3 py-1.5 text-sm rounded-sm transition-colors ${sortBy === opt.value
                                            ? 'bg-[#ee4d2d] text-white'
                                            : 'bg-white text-gray-700 hover:bg-gray-100'
                                        }`}
                                >
                                    {opt.label}
                                </button>
                            ))}
                            <select
                                value={sortBy.startsWith('price') ? sortBy : ''}
                                onChange={(e) => e.target.value && setSortBy(e.target.value)}
                                className="ml-auto px-3 py-1.5 text-sm border bg-white outline-none"
                            >
                                <option value="">Giá</option>
                                <option value="price-asc">Giá: Thấp đến Cao</option>
                                <option value="price-desc">Giá: Cao đến Thấp</option>
                            </select>
                        </div>

                        {/* Products Grid */}
                        {filteredProducts.length === 0 ? (
                            <div className="bg-white rounded-sm p-12 text-center">
                                <div className="text-5xl mb-4">🔍</div>
                                <p className="text-gray-500">Không tìm thấy sản phẩm phù hợp</p>
                            </div>
                        ) : (
                            <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-[10px]">
                                {filteredProducts.map(product => (
                                    <Link key={product.id} href={`/products/${product.id}`} className="product-card group">
                                        <div className="relative aspect-square bg-gray-50 flex items-center justify-center overflow-hidden">
                                            <span className="text-6xl product-image">{product.image}</span>

                                            {/* Discount badge */}
                                            {product.originalPrice > product.price && (
                                                <div className="discount-badge">
                                                    -{getDiscount(product.price, product.originalPrice)}%
                                                </div>
                                            )}

                                            {/* Top left badges */}
                                            <div className="absolute top-0 left-0 flex flex-col gap-0.5">
                                                {product.isOfficial && (
                                                    <div className="bg-[#ee4d2d] text-white text-[10px] px-1">Mall</div>
                                                )}
                                                {product.isFavorite && (
                                                    <div className="bg-[#ee4d2d] text-white text-[10px] px-1">Yêu thích</div>
                                                )}
                                            </div>

                                            {/* Add to cart */}
                                            <button
                                                onClick={(e) => addToCart(product.name, e)}
                                                className="absolute bottom-2 right-2 w-8 h-8 bg-[#ee4d2d] text-white rounded-sm flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity text-sm hover:bg-[#d73211]"
                                            >
                                                +
                                            </button>
                                        </div>

                                        <div className="p-2">
                                            <h3 className="text-xs line-clamp-2 h-8 mb-1">{product.name}</h3>

                                            {/* Free ship badge */}
                                            {product.freeShip && (
                                                <div className="inline-flex items-center text-[10px] text-[#00bfa5] border border-[#00bfa5] px-1 mb-1">
                                                    <svg className="w-3 h-3 mr-0.5" fill="currentColor" viewBox="0 0 24 24">
                                                        <path d="M20 8h-3V4H3c-1.1 0-2 .9-2 2v11h2c0 1.66 1.34 3 3 3s3-1.34 3-3h6c0 1.66 1.34 3 3 3s3-1.34 3-3h2v-5l-3-4zM6 18.5c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5zm13.5-9l1.96 2.5H17V9.5h2.5zm-1.5 9c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5z" />
                                                    </svg>
                                                    Miễn phí vận chuyển
                                                </div>
                                            )}

                                            <div className="flex items-end justify-between">
                                                <div>
                                                    <span className="price-current text-sm font-medium">₫{formatPrice(product.price)}</span>
                                                    {product.originalPrice > product.price && (
                                                        <span className="price-original block">₫{formatPrice(product.originalPrice)}</span>
                                                    )}
                                                </div>
                                            </div>

                                            <div className="flex items-center justify-between mt-1 text-[11px] text-gray-500">
                                                <span className="flex items-center gap-0.5">
                                                    <span className="star-rating">★</span> {product.rating}
                                                </span>
                                                <span>Đã bán {product.sold}</span>
                                            </div>
                                            <div className="text-[11px] text-gray-400 mt-0.5">{product.location}</div>
                                        </div>
                                    </Link>
                                ))}
                            </div>
                        )}
                    </main>
                </div>
            </div>
        </div>
    );
}
