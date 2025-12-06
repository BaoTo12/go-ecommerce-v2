'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';

export default function HomePage() {
    const [countdown, setCountdown] = useState({ hours: 2, minutes: 45, seconds: 30 });
    const [currentBanner, setCurrentBanner] = useState(0);
    const [isLoaded, setIsLoaded] = useState(false);
    const [hoveredCategory, setHoveredCategory] = useState<number | null>(null);

    const categories = [
        { icon: '📱', name: 'Điện Thoại & Phụ Kiện' },
        { icon: '💻', name: 'Máy Tính & Laptop' },
        { icon: '📺', name: 'Thiết Bị Điện Tử' },
        { icon: '📷', name: 'Máy Ảnh & Máy Quay' },
        { icon: '⌚', name: 'Đồng Hồ' },
        { icon: '👟', name: 'Giày Dép Nam' },
        { icon: '👠', name: 'Giày Dép Nữ' },
        { icon: '👜', name: 'Túi Ví Nữ' },
        { icon: '👗', name: 'Thời Trang Nữ' },
        { icon: '👔', name: 'Thời Trang Nam' },
        { icon: '💄', name: 'Sắc Đẹp' },
        { icon: '🏠', name: 'Nhà Cửa & Đời Sống' },
        { icon: '🍼', name: 'Mẹ & Bé' },
        { icon: '🏃', name: 'Thể Thao & Du Lịch' },
        { icon: '🎮', name: 'Gaming & Console' },
    ];

    const flashSaleProducts = [
        { id: 1, name: 'iPhone 15 Pro Max', price: 29990000, originalPrice: 34990000, discount: 14, sold: 87, image: '📱' },
        { id: 2, name: 'AirPods Pro 2', price: 4990000, originalPrice: 6990000, discount: 29, sold: 92, image: '🎧' },
        { id: 3, name: 'MacBook Air M3', price: 24990000, originalPrice: 27990000, discount: 11, sold: 45, image: '💻' },
        { id: 4, name: 'Galaxy Watch 6', price: 6990000, originalPrice: 8990000, discount: 22, sold: 78, image: '⌚' },
        { id: 5, name: 'iPad Pro M2', price: 21990000, originalPrice: 25990000, discount: 15, sold: 34, image: '📟' },
        { id: 6, name: 'Sony WH-1000XM5', price: 6990000, originalPrice: 8990000, discount: 22, sold: 65, image: '🎧' },
    ];

    const recommendedProducts = [
        { id: 'p1', name: 'iPhone 15 Pro Max 256GB Titan Xanh Chính Hãng VN/A', price: 29990000, originalPrice: 34990000, sold: '12.3k', rating: 4.9, image: '📱', location: 'TP. Hồ Chí Minh' },
        { id: 'p2', name: 'Áo Hoodie Unisex Form Rộng Nỉ Cotton Premium Dày Dặn', price: 199000, originalPrice: 350000, sold: '45.2k', rating: 4.8, image: '👕', location: 'Hà Nội' },
        { id: 'p3', name: 'Son Dưỡng Môi Dior Addict Lip Glow chính hãng', price: 950000, originalPrice: 1200000, sold: '8.7k', rating: 4.9, image: '💄', location: 'TP. Hồ Chí Minh' },
        { id: 'p4', name: 'Giày Nike Air Force 1 Low White Chính Hãng', price: 2590000, originalPrice: 3200000, sold: '5.2k', rating: 4.7, image: '👟', location: 'Hà Nội' },
        { id: 'p5', name: 'Tai Nghe Bluetooth Apple AirPods Pro 2 USB-C', price: 4990000, originalPrice: 6990000, sold: '15.1k', rating: 4.9, image: '🎧', location: 'TP. Hồ Chí Minh' },
        { id: 'p6', name: 'Nồi Chiên Không Dầu Lock&Lock 5.2L EJF356BLK', price: 1290000, originalPrice: 2490000, sold: '23.4k', rating: 4.8, image: '🍳', location: 'Hà Nội' },
        { id: 'p7', name: 'Laptop Dell XPS 13 Plus Intel Core i7 Gen 13', price: 32990000, originalPrice: 38990000, sold: '1.2k', rating: 4.7, image: '💻', location: 'TP. Hồ Chí Minh' },
        { id: 'p8', name: 'Bộ Skincare COSRX Advanced Snail 96 Set', price: 450000, originalPrice: 650000, sold: '67.8k', rating: 4.9, image: '🧴', location: 'Hà Nội' },
        { id: 'p9', name: 'Đồng Hồ Casio G-Shock GA-2100 Chính Hãng', price: 2890000, originalPrice: 3500000, sold: '4.5k', rating: 4.8, image: '⌚', location: 'TP. Hồ Chí Minh' },
        { id: 'p10', name: 'Túi Xách Nữ Charles & Keith CNK Authentic', price: 890000, originalPrice: 1290000, sold: '9.1k', rating: 4.6, image: '👜', location: 'Hà Nội' },
    ];

    const banners = [
        { bg: 'linear-gradient(135deg, #f53d2d 0%, #ff6533 100%)', text: '🔥 Flash Sale 12.12', subtitle: 'Giảm đến 50%!' },
        { bg: 'linear-gradient(135deg, #00bfa5 0%, #00897b 100%)', text: '🚚 Miễn Phí Vận Chuyển', subtitle: 'Đơn từ 0Đ' },
        { bg: 'linear-gradient(135deg, #5c6bc0 0%, #3949ab 100%)', text: '💳 Hoàn Tiền 10%', subtitle: 'Qua ShopeePay' },
    ];

    useEffect(() => {
        setIsLoaded(true);
    }, []);

    useEffect(() => {
        const timer = setInterval(() => {
            setCountdown(prev => {
                if (prev.seconds > 0) return { ...prev, seconds: prev.seconds - 1 };
                if (prev.minutes > 0) return { ...prev, minutes: prev.minutes - 1, seconds: 59 };
                if (prev.hours > 0) return { hours: prev.hours - 1, minutes: 59, seconds: 59 };
                return { hours: 2, minutes: 0, seconds: 0 };
            });
        }, 1000);
        return () => clearInterval(timer);
    }, []);

    useEffect(() => {
        const timer = setInterval(() => {
            setCurrentBanner(prev => (prev + 1) % banners.length);
        }, 4000);
        return () => clearInterval(timer);
    }, [banners.length]);

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    return (
        <div className={`min-h-screen bg-[#f5f5f5] ${isLoaded ? 'animate-fade-in' : 'opacity-0'}`}>
            {/* Banner Carousel */}
            <section className="container mx-auto px-4 pt-4">
                <div className="grid grid-cols-3 gap-2">
                    <div className="col-span-2 relative h-[235px] rounded-sm overflow-hidden">
                        {banners.map((banner, i) => (
                            <div
                                key={i}
                                className={`absolute inset-0 flex flex-col items-center justify-center text-white transition-all duration-700 ${currentBanner === i ? 'opacity-100 scale-100' : 'opacity-0 scale-105'
                                    }`}
                                style={{ background: banner.bg }}
                            >
                                <span className="text-4xl font-bold mb-2 animate-fade-in-down">{banner.text}</span>
                                <span className="text-xl opacity-90 animate-fade-in-up">{banner.subtitle}</span>
                            </div>
                        ))}
                        {/* Dots */}
                        <div className="absolute bottom-3 left-1/2 -translate-x-1/2 flex gap-2">
                            {banners.map((_, i) => (
                                <button
                                    key={i}
                                    onClick={() => setCurrentBanner(i)}
                                    className={`transition-all duration-300 ${currentBanner === i ? 'w-6 h-2 bg-white rounded-full' : 'w-2 h-2 bg-white/50 rounded-full hover:bg-white/80'
                                        }`}
                                />
                            ))}
                        </div>
                    </div>
                    <div className="flex flex-col gap-2">
                        <Link href="/live" className="flex-1 bg-gradient-to-r from-[#00bfa5] to-[#00897b] rounded-sm flex items-center justify-center text-white font-semibold hover-lift hover-shine">
                            <span className="flex items-center gap-2">
                                <span className="animate-pulse">🔴</span> Shopee Live
                            </span>
                        </Link>
                        <Link href="/deals/coupons" className="flex-1 bg-gradient-to-r from-[#5c6bc0] to-[#3949ab] rounded-sm flex items-center justify-center text-white font-semibold hover-lift hover-shine">
                            🎁 Voucher Xtra
                        </Link>
                    </div>
                </div>
            </section>

            {/* Categories */}
            <section className="container mx-auto px-4 py-4">
                <div className="bg-white rounded-sm shadow-sm overflow-hidden">
                    <div className="p-4 border-b text-gray-500 font-medium">DANH MỤC</div>
                    <div className="grid grid-cols-5 md:grid-cols-10 lg:grid-cols-15">
                        {categories.map((cat, index) => (
                            <Link
                                key={cat.name}
                                href={`/products?category=${encodeURIComponent(cat.name)}`}
                                className="category-item animate-fade-in-up"
                                style={{ animationDelay: `${index * 30}ms` }}
                                onMouseEnter={() => setHoveredCategory(index)}
                                onMouseLeave={() => setHoveredCategory(null)}
                            >
                                <span
                                    className={`text-3xl mb-2 category-icon transition-transform duration-300 ${hoveredCategory === index ? 'animate-wiggle' : ''
                                        }`}
                                >
                                    {cat.icon}
                                </span>
                                <span className="text-xs text-gray-600 text-center line-clamp-2">{cat.name}</span>
                            </Link>
                        ))}
                    </div>
                </div>
            </section>

            {/* Flash Sale */}
            <section className="container mx-auto px-4 py-2">
                <div className="bg-white rounded-sm shadow-sm overflow-hidden">
                    <div className="flex items-center justify-between p-4 border-b bg-gradient-to-r from-[#fff5f5] to-white">
                        <div className="flex items-center gap-4">
                            <Link href="/deals/flash-sale" className="text-[#ee4d2d] text-xl font-bold flex items-center">
                                <span className="mr-2 animate-float">⚡</span> FLASH SALE
                            </Link>
                            <div className="flex items-center gap-1">
                                <div className="timer-box">{String(countdown.hours).padStart(2, '0')}</div>
                                <span className="text-[#ee4d2d] font-bold animate-pulse">:</span>
                                <div className="timer-box">{String(countdown.minutes).padStart(2, '0')}</div>
                                <span className="text-[#ee4d2d] font-bold animate-pulse">:</span>
                                <div className="timer-box">{String(countdown.seconds).padStart(2, '0')}</div>
                            </div>
                        </div>
                        <Link href="/deals/flash-sale" className="text-[#ee4d2d] text-sm hover:opacity-80 flex items-center gap-1 group">
                            Xem tất cả
                            <svg className="w-4 h-4 transition-transform group-hover:translate-x-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                            </svg>
                        </Link>
                    </div>
                    <div className="grid grid-cols-6 divide-x">
                        {flashSaleProducts.map((product, index) => (
                            <Link
                                key={product.id}
                                href="/deals/flash-sale"
                                className="product-card border-0 rounded-none animate-fade-in-up hover-shine"
                                style={{ animationDelay: `${index * 100}ms` }}
                            >
                                <div className="relative aspect-square bg-gray-50 flex items-center justify-center overflow-hidden">
                                    <span className="text-5xl product-image">{product.image}</span>
                                    <div className="discount-badge">-{product.discount}%</div>
                                </div>
                                <div className="p-2 text-center">
                                    <div className="price-current text-lg font-bold">₫{formatPrice(product.price)}</div>
                                    <div className="flash-progress mt-2">
                                        <div className="flash-progress-bar" style={{ width: `${product.sold}%` }}>
                                            <span>Đã bán {product.sold}</span>
                                        </div>
                                    </div>
                                </div>
                            </Link>
                        ))}
                    </div>
                </div>
            </section>

            {/* Quick Access */}
            <section className="container mx-auto px-4 py-4">
                <div className="grid grid-cols-4 gap-2">
                    {[
                        { href: '/rewards', icon: '🎮', title: 'Shopee Xu', desc: 'Đổi xu lấy quà', gradient: 'from-yellow-400 to-orange-500' },
                        { href: '/deals/coupons', icon: '🎟️', title: 'Mã Giảm Giá', desc: 'Săn voucher hot', gradient: 'from-purple-500 to-pink-500' },
                        { href: '/live', icon: '🔴', title: 'Shopee Live', desc: 'Xem & mua sắm', gradient: 'from-red-500 to-rose-500', live: true },
                        { href: '/products', icon: '🛍️', title: 'Hàng Chính Hãng', desc: '100% authentic', gradient: 'from-blue-500 to-cyan-500' },
                    ].map((item, index) => (
                        <Link
                            key={item.href}
                            href={item.href}
                            className={`bg-white rounded-sm p-4 flex items-center gap-3 hover-lift hover-shine animate-fade-in-up`}
                            style={{ animationDelay: `${index * 100}ms` }}
                        >
                            <span className={`text-3xl ${item.live ? 'animate-pulse' : ''}`}>
                                {item.icon}
                            </span>
                            <div>
                                <div className="font-semibold text-sm">{item.title}</div>
                                <div className="text-xs text-gray-500">{item.desc}</div>
                            </div>
                        </Link>
                    ))}
                </div>
            </section>

            {/* Recommendations */}
            <section className="container mx-auto px-4 py-4">
                <div className="bg-white rounded-sm shadow-sm overflow-hidden">
                    <div className="p-4 border-b text-center sticky top-[120px] bg-white z-10">
                        <span className="text-[#ee4d2d] text-lg font-bold">GỢI Ý HÔM NAY</span>
                    </div>
                    <div className="grid grid-cols-2 md:grid-cols-5 lg:grid-cols-6 gap-[2px] bg-gray-100 p-[2px]">
                        {recommendedProducts.map((product, index) => (
                            <Link
                                key={product.id}
                                href={`/products/${product.id}`}
                                className="product-card bg-white animate-fade-in-up"
                                style={{ animationDelay: `${index * 50}ms` }}
                            >
                                <div className="relative aspect-square bg-gray-50 flex items-center justify-center overflow-hidden">
                                    <span className="text-6xl product-image">{product.image}</span>
                                    <div className="absolute top-0 left-0 bg-[#ee4d2d] text-white text-[10px] px-1 py-0.5">
                                        Yêu thích
                                    </div>
                                </div>
                                <div className="p-2">
                                    <h3 className="text-xs line-clamp-2 h-8 mb-2">{product.name}</h3>
                                    <div className="flex items-center justify-between">
                                        <span className="price-current">₫{formatPrice(product.price)}</span>
                                        <span className="text-xs text-gray-500">Đã bán {product.sold}</span>
                                    </div>
                                    <div className="flex items-center gap-1 mt-1 text-xs text-gray-500">
                                        <span className="star-rating">★</span>
                                        <span>{product.rating}</span>
                                        <span className="mx-1">|</span>
                                        <span className="truncate">{product.location}</span>
                                    </div>
                                </div>
                            </Link>
                        ))}
                    </div>
                    <div className="p-4 text-center">
                        <Link
                            href="/products"
                            className="inline-block px-10 py-2 border border-gray-300 text-gray-600 hover:bg-gray-50 transition-all hover:border-[#ee4d2d] hover:text-[#ee4d2d] text-sm"
                        >
                            Xem Thêm
                        </Link>
                    </div>
                </div>
            </section>
        </div>
    );
}
