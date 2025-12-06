'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';

interface Product {
    id: string;
    name: string;
    price: number;
    originalPrice: number;
    discount: number;
    sold: number;
    rating: number;
    reviews: number;
    image: string;
    category: string;
    isNew?: boolean;
    isBestSeller?: boolean;
}

export default function ProductsPage() {
    const [products, setProducts] = useState<Product[]>([]);
    const [filteredProducts, setFilteredProducts] = useState<Product[]>([]);
    const [selectedCategory, setSelectedCategory] = useState('Tất cả');
    const [sortBy, setSortBy] = useState('popular');
    const [searchQuery, setSearchQuery] = useState('');
    const [cart, setCart] = useState<string[]>([]);
    const [notification, setNotification] = useState<string | null>(null);
    const [priceRange, setPriceRange] = useState<[number, number]>([0, 50000000]);

    const categories = ['Tất cả', 'Điện thoại', 'Laptop', 'Thời trang', 'Làm đẹp', 'Nhà cửa', 'Giày dép'];

    useEffect(() => {
        // Load products
        setProducts([
            { id: 'p1', name: 'iPhone 15 Pro Max 256GB', price: 29990000, originalPrice: 34990000, discount: 14, sold: 1234, rating: 4.9, reviews: 856, image: '📱', category: 'Điện thoại', isBestSeller: true },
            { id: 'p2', name: 'Samsung Galaxy S24 Ultra', price: 25990000, originalPrice: 29990000, discount: 13, sold: 987, rating: 4.8, reviews: 543, image: '📲', category: 'Điện thoại' },
            { id: 'p3', name: 'MacBook Air M3 13"', price: 27990000, originalPrice: 31990000, discount: 12, sold: 654, rating: 4.9, reviews: 432, image: '💻', category: 'Laptop', isBestSeller: true },
            { id: 'p4', name: 'Áo Hoodie Unisex', price: 299000, originalPrice: 450000, discount: 34, sold: 5432, rating: 4.7, reviews: 1234, image: '👕', category: 'Thời trang', isNew: true },
            { id: 'p5', name: 'Giày Nike Air Max', price: 3990000, originalPrice: 4590000, discount: 13, sold: 2345, rating: 4.8, reviews: 876, image: '👟', category: 'Giày dép' },
            { id: 'p6', name: 'Son Dưỡng Môi', price: 150000, originalPrice: 250000, discount: 40, sold: 8765, rating: 4.6, reviews: 3456, image: '💄', category: 'Làm đẹp', isNew: true },
            { id: 'p7', name: 'Nồi Chiên Không Dầu', price: 1990000, originalPrice: 2990000, discount: 33, sold: 3456, rating: 4.8, reviews: 2345, image: '🍳', category: 'Nhà cửa' },
            { id: 'p8', name: 'Laptop Dell XPS 13', price: 32990000, originalPrice: 38990000, discount: 15, sold: 432, rating: 4.7, reviews: 234, image: '💻', category: 'Laptop' },
            { id: 'p9', name: 'Quần Jean Slim Fit', price: 399000, originalPrice: 599000, discount: 33, sold: 4567, rating: 4.5, reviews: 1876, image: '👖', category: 'Thời trang' },
            { id: 'p10', name: 'Serum Vitamin C', price: 350000, originalPrice: 500000, discount: 30, sold: 6543, rating: 4.7, reviews: 2987, image: '🧴', category: 'Làm đẹp', isBestSeller: true },
            { id: 'p11', name: 'Đèn Bàn LED', price: 299000, originalPrice: 450000, discount: 34, sold: 2345, rating: 4.6, reviews: 876, image: '💡', category: 'Nhà cửa' },
            { id: 'p12', name: 'Giày Adidas Ultraboost', price: 4290000, originalPrice: 4990000, discount: 14, sold: 1876, rating: 4.9, reviews: 654, image: '👟', category: 'Giày dép', isNew: true },
        ]);
    }, []);

    useEffect(() => {
        let result = [...products];

        // Category filter
        if (selectedCategory !== 'Tất cả') {
            result = result.filter(p => p.category === selectedCategory);
        }

        // Search filter
        if (searchQuery) {
            result = result.filter(p =>
                p.name.toLowerCase().includes(searchQuery.toLowerCase())
            );
        }

        // Price filter
        result = result.filter(p => p.price >= priceRange[0] && p.price <= priceRange[1]);

        // Sorting
        switch (sortBy) {
            case 'price-asc':
                result.sort((a, b) => a.price - b.price);
                break;
            case 'price-desc':
                result.sort((a, b) => b.price - a.price);
                break;
            case 'newest':
                result.sort((a, b) => (b.isNew ? 1 : 0) - (a.isNew ? 1 : 0));
                break;
            case 'best-seller':
                result.sort((a, b) => b.sold - a.sold);
                break;
            default: // popular
                result.sort((a, b) => b.rating - a.rating);
        }

        setFilteredProducts(result);
    }, [products, selectedCategory, sortBy, searchQuery, priceRange]);

    const addToCart = (productId: string, productName: string) => {
        setCart(prev => [...prev, productId]);
        setNotification(`✓ Đã thêm "${productName}" vào giỏ hàng`);
        setTimeout(() => setNotification(null), 3000);
    };

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    return (
        <div className="min-h-screen bg-[#F5F5F5]">
            {/* Notification Toast */}
            {notification && (
                <div className="fixed top-20 right-4 z-50 bg-green-500 text-white px-6 py-3 rounded-lg shadow-lg animate-pulse">
                    {notification}
                </div>
            )}

            {/* Cart Badge */}
            <Link
                href="/cart"
                className="fixed bottom-6 right-6 z-40 bg-[#EE4D2D] text-white p-4 rounded-full shadow-lg hover:bg-[#D73211] transition-colors"
            >
                <span className="text-2xl">🛒</span>
                {cart.length > 0 && (
                    <span className="absolute -top-1 -right-1 bg-yellow-400 text-black text-xs font-bold w-6 h-6 rounded-full flex items-center justify-center">
                        {cart.length}
                    </span>
                )}
            </Link>

            {/* Header */}
            <div className="bg-white border-b sticky top-14 z-30">
                <div className="container mx-auto px-4 py-3">
                    <div className="flex items-center gap-4">
                        {/* Search */}
                        <div className="flex-1 relative">
                            <input
                                type="text"
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                placeholder="Tìm kiếm sản phẩm..."
                                className="w-full border border-gray-300 rounded-sm px-4 py-2 pr-10 focus:outline-none focus:border-[#EE4D2D]"
                            />
                            <span className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400">🔍</span>
                        </div>

                        {/* Sort */}
                        <select
                            value={sortBy}
                            onChange={(e) => setSortBy(e.target.value)}
                            className="border border-gray-300 rounded-sm px-4 py-2 focus:outline-none focus:border-[#EE4D2D]"
                        >
                            <option value="popular">Phổ biến</option>
                            <option value="newest">Mới nhất</option>
                            <option value="best-seller">Bán chạy</option>
                            <option value="price-asc">Giá: Thấp → Cao</option>
                            <option value="price-desc">Giá: Cao → Thấp</option>
                        </select>
                    </div>

                    {/* Categories */}
                    <div className="flex gap-2 mt-3 overflow-x-auto pb-2">
                        {categories.map(cat => (
                            <button
                                key={cat}
                                onClick={() => setSelectedCategory(cat)}
                                className={`px-4 py-1.5 rounded-full text-sm whitespace-nowrap transition-colors ${selectedCategory === cat
                                        ? 'bg-[#EE4D2D] text-white'
                                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                                    }`}
                            >
                                {cat}
                            </button>
                        ))}
                    </div>
                </div>
            </div>

            {/* Products Grid */}
            <div className="container mx-auto px-4 py-6">
                <div className="flex items-center justify-between mb-4">
                    <h1 className="text-lg font-bold">
                        {selectedCategory === 'Tất cả' ? 'Tất cả sản phẩm' : selectedCategory}
                        <span className="text-gray-500 font-normal ml-2">({filteredProducts.length} sản phẩm)</span>
                    </h1>
                </div>

                {filteredProducts.length === 0 ? (
                    <div className="bg-white rounded p-12 text-center">
                        <div className="text-6xl mb-4">🔍</div>
                        <p className="text-gray-500">Không tìm thấy sản phẩm phù hợp</p>
                    </div>
                ) : (
                    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-3">
                        {filteredProducts.map(product => (
                            <div
                                key={product.id}
                                className="bg-white rounded-sm overflow-hidden hover:shadow-lg transition-shadow group cursor-pointer border border-transparent hover:border-[#EE4D2D]"
                            >
                                {/* Image */}
                                <div className="relative aspect-square bg-gray-100 flex items-center justify-center text-6xl overflow-hidden">
                                    <span className="group-hover:scale-110 transition-transform">{product.image}</span>

                                    {/* Badges */}
                                    {product.discount > 0 && (
                                        <span className="absolute top-0 right-0 bg-[#EE4D2D] text-white text-xs px-2 py-1">
                                            -{product.discount}%
                                        </span>
                                    )}
                                    {product.isNew && (
                                        <span className="absolute top-0 left-0 bg-green-500 text-white text-xs px-2 py-1">
                                            Mới
                                        </span>
                                    )}
                                    {product.isBestSeller && !product.isNew && (
                                        <span className="absolute top-0 left-0 bg-yellow-500 text-white text-xs px-2 py-1">
                                            Bán chạy
                                        </span>
                                    )}

                                    {/* Quick Add */}
                                    <button
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            addToCart(product.id, product.name);
                                        }}
                                        className="absolute bottom-2 right-2 bg-[#EE4D2D] text-white p-2 rounded-full opacity-0 group-hover:opacity-100 transition-opacity hover:bg-[#D73211]"
                                    >
                                        🛒
                                    </button>
                                </div>

                                {/* Info */}
                                <div className="p-3">
                                    <h3 className="text-sm line-clamp-2 h-10 mb-2">{product.name}</h3>

                                    <div className="flex items-baseline gap-2 mb-1">
                                        <span className="text-[#EE4D2D] font-bold">₫{formatPrice(product.price)}</span>
                                        {product.originalPrice > product.price && (
                                            <span className="text-gray-400 text-xs line-through">₫{formatPrice(product.originalPrice)}</span>
                                        )}
                                    </div>

                                    <div className="flex items-center justify-between text-xs text-gray-500">
                                        <span className="flex items-center gap-1">
                                            <span className="text-yellow-400">★</span>
                                            {product.rating}
                                        </span>
                                        <span>Đã bán {product.sold > 1000 ? `${(product.sold / 1000).toFixed(1)}k` : product.sold}</span>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
