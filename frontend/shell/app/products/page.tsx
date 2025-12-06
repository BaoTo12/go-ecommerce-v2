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
    const [notification, setNotification] = useState<string | null>(null);

    const categories = ['Tất cả', 'Điện thoại', 'Laptop', 'Thời trang', 'Làm đẹp', 'Nhà cửa', 'Giày dép'];

    useEffect(() => {
        setProducts([
            { id: 'p1', name: 'iPhone 15 Pro Max 256GB', price: 29990000, originalPrice: 34990000, discount: 14, sold: 1234, rating: 4.9, reviews: 856, image: '📱', category: 'Điện thoại', isBestSeller: true },
            { id: 'p2', name: 'Samsung Galaxy S24 Ultra', price: 25990000, originalPrice: 29990000, discount: 13, sold: 987, rating: 4.8, reviews: 543, image: '📲', category: 'Điện thoại' },
            { id: 'p3', name: 'MacBook Air M3 13"', price: 27990000, originalPrice: 31990000, discount: 12, sold: 654, rating: 4.9, reviews: 432, image: '💻', category: 'Laptop', isBestSeller: true },
            { id: 'p4', name: 'Áo Hoodie Unisex', price: 299000, originalPrice: 450000, discount: 34, sold: 5432, rating: 4.7, reviews: 1234, image: '👕', category: 'Thời trang', isNew: true },
            { id: 'p5', name: 'Giày Nike Air Max', price: 3990000, originalPrice: 4590000, discount: 13, sold: 2345, rating: 4.8, reviews: 876, image: '👟', category: 'Giày dép' },
            { id: 'p6', name: 'Son Dưỡng Môi Dior', price: 950000, originalPrice: 1200000, discount: 21, sold: 8765, rating: 4.6, reviews: 3456, image: '💄', category: 'Làm đẹp', isNew: true },
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
                result.sort((a, b) => (b.isNew ? 1 : 0) - (a.isNew ? 1 : 0));
                break;
            case 'best-seller':
                result.sort((a, b) => b.sold - a.sold);
                break;
            default:
                result.sort((a, b) => b.rating - a.rating);
        }

        setFilteredProducts(result);
    }, [products, selectedCategory, sortBy, searchQuery]);

    const addToCart = (productName: string) => {
        setNotification(`✓ Đã thêm "${productName}" vào giỏ hàng`);
        setTimeout(() => setNotification(null), 3000);
    };

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    return (
        <div className="min-h-screen bg-[#F5F5F5] animate-fade-in">
            {/* Notification Toast */}
            {notification && (
                <div className="fixed top-24 right-4 z-50 bg-green-500 text-white px-6 py-3 rounded-2xl shadow-xl animate-slide-in-right">
                    {notification}
                </div>
            )}

            {/* Floating Cart Button */}
            <Link
                href="/cart"
                className="fixed bottom-6 right-6 z-40 w-14 h-14 bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] text-white rounded-full shadow-xl flex items-center justify-center text-2xl hover:scale-110 transition-transform"
            >
                🛒
            </Link>

            {/* Header */}
            <div className="bg-white border-b sticky top-[104px] z-30">
                <div className="container mx-auto px-4 py-4">
                    <div className="flex flex-col md:flex-row items-center gap-4">
                        {/* Search */}
                        <div className="flex-1 relative w-full">
                            <input
                                type="text"
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                placeholder="Tìm kiếm sản phẩm..."
                                className="w-full border-2 border-gray-200 rounded-full px-6 py-3 pr-12 focus:outline-none focus:border-[#EE4D2D] transition-colors"
                            />
                            <span className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 text-xl">🔍</span>
                        </div>

                        {/* Sort */}
                        <select
                            value={sortBy}
                            onChange={(e) => setSortBy(e.target.value)}
                            className="border-2 border-gray-200 rounded-full px-6 py-3 focus:outline-none focus:border-[#EE4D2D] bg-white"
                        >
                            <option value="popular">Phổ biến</option>
                            <option value="newest">Mới nhất</option>
                            <option value="best-seller">Bán chạy</option>
                            <option value="price-asc">Giá: Thấp → Cao</option>
                            <option value="price-desc">Giá: Cao → Thấp</option>
                        </select>
                    </div>

                    {/* Categories */}
                    <div className="flex gap-2 mt-4 overflow-x-auto pb-2">
                        {categories.map(cat => (
                            <button
                                key={cat}
                                onClick={() => setSelectedCategory(cat)}
                                className={`px-5 py-2 rounded-full text-sm whitespace-nowrap transition-all font-medium ${selectedCategory === cat
                                        ? 'bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] text-white shadow-lg'
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
                    <h1 className="text-xl font-bold">
                        {selectedCategory === 'Tất cả' ? 'Tất cả sản phẩm' : selectedCategory}
                        <span className="text-gray-500 font-normal ml-2">({filteredProducts.length} sản phẩm)</span>
                    </h1>
                </div>

                {filteredProducts.length === 0 ? (
                    <div className="bg-white rounded-2xl p-12 text-center">
                        <div className="text-6xl mb-4">🔍</div>
                        <p className="text-gray-500">Không tìm thấy sản phẩm phù hợp</p>
                    </div>
                ) : (
                    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
                        {filteredProducts.map((product, index) => (
                            <div
                                key={product.id}
                                className="bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-xl transition-all group cursor-pointer border-2 border-transparent hover:border-[#EE4D2D] animate-slide-up"
                                style={{ animationDelay: `${index * 50}ms` }}
                            >
                                {/* Image - Clickable to detail */}
                                <Link href={`/products/${product.id}`}>
                                    <div className="relative aspect-square bg-gradient-to-br from-gray-50 to-gray-100 flex items-center justify-center overflow-hidden">
                                        <span className="text-6xl group-hover:scale-125 transition-transform duration-500">{product.image}</span>

                                        {/* Badges */}
                                        {product.discount > 0 && (
                                            <span className="absolute top-2 right-2 bg-[#EE4D2D] text-white text-xs px-2 py-1 rounded-full font-bold">
                                                -{product.discount}%
                                            </span>
                                        )}
                                        {product.isNew && (
                                            <span className="absolute top-2 left-2 bg-green-500 text-white text-xs px-2 py-1 rounded-full font-bold">
                                                Mới
                                            </span>
                                        )}
                                        {product.isBestSeller && !product.isNew && (
                                            <span className="absolute top-2 left-2 bg-yellow-500 text-white text-xs px-2 py-1 rounded-full font-bold">
                                                🔥 Hot
                                            </span>
                                        )}
                                    </div>
                                </Link>

                                {/* Info */}
                                <div className="p-4">
                                    <Link href={`/products/${product.id}`}>
                                        <h3 className="text-sm line-clamp-2 h-10 mb-2 hover:text-[#EE4D2D] transition-colors">{product.name}</h3>
                                    </Link>

                                    <div className="flex items-baseline gap-2 mb-2">
                                        <span className="text-[#EE4D2D] font-bold text-lg">₫{formatPrice(product.price)}</span>
                                    </div>
                                    {product.originalPrice > product.price && (
                                        <span className="text-gray-400 text-xs line-through">₫{formatPrice(product.originalPrice)}</span>
                                    )}

                                    <div className="flex items-center justify-between text-xs text-gray-500 mt-2">
                                        <span className="flex items-center gap-1">
                                            <span className="text-yellow-400">★</span>
                                            {product.rating}
                                        </span>
                                        <span>Đã bán {product.sold > 1000 ? `${(product.sold / 1000).toFixed(1)}k` : product.sold}</span>
                                    </div>

                                    {/* Add to cart button */}
                                    <button
                                        onClick={(e) => {
                                            e.preventDefault();
                                            addToCart(product.name);
                                        }}
                                        className="w-full mt-3 py-2 bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] text-white rounded-full font-medium opacity-0 group-hover:opacity-100 transition-all hover:shadow-lg flex items-center justify-center gap-2"
                                    >
                                        <span>🛒</span> Thêm vào giỏ
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
