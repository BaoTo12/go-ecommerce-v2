'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';

export default function HomePage() {
    const [countdown, setCountdown] = useState({ hours: 2, minutes: 45, seconds: 30 });
    const [stats, setStats] = useState({ users: 12453, orders: 847, revenue: 5.6 });
    const [currentBanner, setCurrentBanner] = useState(0);
    const [isLoaded, setIsLoaded] = useState(false);
    const [hoveredCategory, setHoveredCategory] = useState<string | null>(null);

    const categories = [
        { icon: '📱', name: 'Điện Thoại', count: 1234, color: 'from-blue-400 to-blue-600' },
        { icon: '💻', name: 'Laptop', count: 567, color: 'from-gray-400 to-gray-600' },
        { icon: '👗', name: 'Thời Trang', count: 4567, color: 'from-pink-400 to-pink-600' },
        { icon: '💄', name: 'Làm Đẹp', count: 2345, color: 'from-rose-400 to-rose-600' },
        { icon: '🏠', name: 'Nhà Cửa', count: 1234, color: 'from-amber-400 to-amber-600' },
        { icon: '🎮', name: 'Gaming', count: 876, color: 'from-purple-400 to-purple-600' },
        { icon: '👟', name: 'Giày Dép', count: 1543, color: 'from-emerald-400 to-emerald-600' },
        { icon: '⌚', name: 'Đồng Hồ', count: 432, color: 'from-yellow-400 to-yellow-600' },
        { icon: '📚', name: 'Sách', count: 2134, color: 'from-indigo-400 to-indigo-600' },
        { icon: '🧸', name: 'Đồ Chơi', count: 987, color: 'from-red-400 to-red-600' },
    ];

    const flashSaleProducts = [
        { id: 1, name: 'iPhone 15 Pro', price: 29990000, originalPrice: 34990000, discount: 14, sold: 87, image: '📱' },
        { id: 2, name: 'AirPods Pro 2', price: 4990000, originalPrice: 6990000, discount: 29, sold: 156, image: '🎧' },
        { id: 3, name: 'MacBook Air', price: 24990000, originalPrice: 27990000, discount: 11, sold: 45, image: '💻' },
        { id: 4, name: 'Galaxy Watch', price: 6990000, originalPrice: 8990000, discount: 22, sold: 78, image: '⌚' },
        { id: 5, name: 'iPad Pro', price: 21990000, originalPrice: 25990000, discount: 15, sold: 34, image: '📟' },
        { id: 6, name: 'Sony WH-1000', price: 6990000, originalPrice: 8990000, discount: 22, sold: 92, image: '🎧' },
    ];

    const banners = [
        { id: 1, title: '🔥 Flash Sale Đặc Biệt', subtitle: 'Giảm đến 90% - Chỉ hôm nay!', emoji: '⚡', color: 'from-[#EE4D2D] via-[#FF6633] to-[#FFAB91]' },
        { id: 2, title: '🎮 Shopee Rewards', subtitle: 'Quay số trúng 1 triệu xu!', emoji: '🎰', color: 'from-yellow-400 via-orange-500 to-red-500' },
        { id: 3, title: '🔴 Shopee Live', subtitle: 'Xem live, mua giá sốc', emoji: '📺', color: 'from-pink-500 via-red-500 to-rose-600' },
    ];

    useEffect(() => {
        setIsLoaded(true);
    }, []);

    // Countdown timer
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

    // Live stats
    useEffect(() => {
        const timer = setInterval(() => {
            setStats(prev => ({
                users: prev.users + Math.floor(Math.random() * 10) - 3,
                orders: prev.orders + Math.floor(Math.random() * 3),
                revenue: prev.revenue + Math.random() * 0.1,
            }));
        }, 2000);
        return () => clearInterval(timer);
    }, []);

    // Banner carousel
    useEffect(() => {
        const timer = setInterval(() => {
            setCurrentBanner(prev => (prev + 1) % banners.length);
        }, 4000);
        return () => clearInterval(timer);
    }, [banners.length]);

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    return (
        <div className={`min-h-screen bg-[#F5F5F5] ${isLoaded ? 'animate-fade-in' : 'opacity-0'}`}>
            {/* Hero Banner Carousel */}
            <section className="relative overflow-hidden h-[300px] md:h-[400px]">
                {banners.map((banner, index) => (
                    <div
                        key={banner.id}
                        className={`absolute inset-0 bg-gradient-to-r ${banner.color} transition-all duration-700 ease-in-out ${currentBanner === index ? 'opacity-100 scale-100' : 'opacity-0 scale-105'
                            }`}
                    >
                        <div className="container mx-auto px-4 h-full flex items-center">
                            <div className="grid md:grid-cols-2 gap-8 items-center w-full">
                                <div className="text-white space-y-4 animate-slide-up">
                                    <h1 className="text-4xl md:text-6xl font-black drop-shadow-lg">
                                        {banner.title}
                                    </h1>
                                    <p className="text-xl md:text-2xl opacity-90 font-medium">
                                        {banner.subtitle}
                                    </p>
                                    <div className="flex gap-4 pt-4">
                                        <Link
                                            href="/deals/flash-sale"
                                            className="btn-primary px-8 py-3 text-lg flex items-center gap-2 hover-shine"
                                        >
                                            <span className="animate-bounce">⚡</span> Mua ngay
                                        </Link>
                                        <Link
                                            href="/products"
                                            className="btn-outline bg-white/10 border-white text-white px-8 py-3 text-lg backdrop-blur-sm hover:bg-white/20"
                                        >
                                            Khám phá
                                        </Link>
                                    </div>
                                </div>
                                <div className="hidden md:flex justify-center">
                                    <span className="text-[150px] animate-float drop-shadow-2xl filter">
                                        {banner.emoji}
                                    </span>
                                </div>
                            </div>
                        </div>

                        {/* Decorative elements */}
                        <div className="absolute top-10 right-10 w-32 h-32 bg-white/10 rounded-full blur-2xl animate-pulse" />
                        <div className="absolute bottom-10 left-10 w-40 h-40 bg-white/10 rounded-full blur-3xl animate-pulse-slow" />
                    </div>
                ))}

                {/* Banner dots */}
                <div className="absolute bottom-6 left-1/2 -translate-x-1/2 flex gap-3 z-10">
                    {banners.map((_, index) => (
                        <button
                            key={index}
                            onClick={() => setCurrentBanner(index)}
                            className={`transition-all duration-300 ${currentBanner === index
                                    ? 'w-8 h-3 bg-white rounded-full'
                                    : 'w-3 h-3 bg-white/50 rounded-full hover:bg-white/80'
                                }`}
                        />
                    ))}
                </div>
            </section>

            {/* Live Stats - Glassmorphism */}
            <section className="bg-gradient-to-r from-slate-900 via-slate-800 to-slate-900 py-4 relative overflow-hidden">
                <div className="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGNpcmNsZSBjeD0iMjAiIGN5PSIyMCIgcj0iMSIgZmlsbD0icmdiYSgyNTUsMjU1LDI1NSwwLjEpIi8+PC9zdmc+')] opacity-30" />
                <div className="container mx-auto px-4 relative">
                    <div className="flex items-center justify-center gap-8 md:gap-16 text-white">
                        <div className="flex items-center gap-3 group">
                            <span className="w-3 h-3 bg-green-400 rounded-full animate-pulse shadow-lg shadow-green-400/50" />
                            <div className="text-center">
                                <span className="text-2xl md:text-3xl font-bold text-green-400 tabular-nums group-hover:scale-110 inline-block transition-transform">
                                    {stats.users.toLocaleString()}
                                </span>
                                <p className="text-xs opacity-75">đang online</p>
                            </div>
                        </div>
                        <div className="w-px h-10 bg-white/20" />
                        <div className="text-center group">
                            <span className="text-2xl md:text-3xl font-bold text-yellow-400 tabular-nums group-hover:scale-110 inline-block transition-transform">
                                {stats.orders.toLocaleString()}
                            </span>
                            <p className="text-xs opacity-75">đơn hôm nay</p>
                        </div>
                        <div className="w-px h-10 bg-white/20" />
                        <div className="text-center group">
                            <span className="text-2xl md:text-3xl font-bold text-blue-400 tabular-nums group-hover:scale-110 inline-block transition-transform">
                                {stats.revenue.toFixed(1)}B₫
                            </span>
                            <p className="text-xs opacity-75">doanh thu</p>
                        </div>
                    </div>
                </div>
            </section>

            {/* Categories - Interactive Grid */}
            <section className="container mx-auto px-4 py-8">
                <div className="bg-white rounded-lg shadow-sm overflow-hidden">
                    <div className="p-4 border-b flex items-center justify-between">
                        <h2 className="font-bold text-xl text-gradient">DANH MỤC</h2>
                        <Link href="/products" className="text-[#EE4D2D] text-sm hover:underline flex items-center gap-1 group">
                            Xem tất cả
                            <span className="group-hover:translate-x-1 transition-transform">→</span>
                        </Link>
                    </div>
                    <div className="grid grid-cols-5 md:grid-cols-10 gap-1 p-4">
                        {categories.map((cat, index) => (
                            <Link
                                key={cat.name}
                                href={`/products?category=${encodeURIComponent(cat.name)}`}
                                onMouseEnter={() => setHoveredCategory(cat.name)}
                                onMouseLeave={() => setHoveredCategory(null)}
                                className={`flex flex-col items-center p-3 rounded-lg transition-all duration-300 text-center group animate-slide-up`}
                                style={{ animationDelay: `${index * 50}ms` }}
                            >
                                <div className={`text-4xl mb-2 transition-all duration-300 ${hoveredCategory === cat.name
                                        ? 'scale-125 transform'
                                        : 'group-hover:scale-110'
                                    }`}>
                                    {cat.icon}
                                </div>
                                <span className="text-xs text-gray-700 font-medium">{cat.name}</span>
                                <span className={`text-[10px] transition-all duration-300 ${hoveredCategory === cat.name
                                        ? 'text-[#EE4D2D] font-bold'
                                        : 'text-gray-400'
                                    }`}>
                                    {cat.count.toLocaleString()}
                                </span>

                                {/* Hover effect background */}
                                <div className={`absolute inset-0 rounded-lg bg-gradient-to-br ${cat.color} opacity-0 group-hover:opacity-10 transition-opacity duration-300 -z-10`} />
                            </Link>
                        ))}
                    </div>
                </div>
            </section>

            {/* Flash Sale - Animated Section */}
            <section className="container mx-auto px-4 py-4">
                <div className="bg-white rounded-lg shadow-sm overflow-hidden relative">
                    {/* Header with pulsing effect */}
                    <div className="bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] p-4 flex items-center justify-between relative overflow-hidden">
                        <div className="absolute inset-0 opacity-30">
                            <div className="absolute inset-0 animate-shimmer bg-gradient-to-r from-transparent via-white/20 to-transparent"
                                style={{ backgroundSize: '200% 100%' }} />
                        </div>

                        <div className="flex items-center gap-4 relative z-10">
                            <span className="text-white text-2xl font-black flex items-center gap-2">
                                <span className="animate-bounce">⚡</span> FLASH SALE
                            </span>
                            <div className="flex items-center gap-1">
                                <div className={`countdown-box ${countdown.hours === 0 && countdown.minutes < 5 ? 'urgent' : ''}`}>
                                    {String(countdown.hours).padStart(2, '0')}
                                </div>
                                <span className="text-white font-bold text-xl animate-pulse">:</span>
                                <div className={`countdown-box ${countdown.hours === 0 && countdown.minutes < 5 ? 'urgent' : ''}`}>
                                    {String(countdown.minutes).padStart(2, '0')}
                                </div>
                                <span className="text-white font-bold text-xl animate-pulse">:</span>
                                <div className={`countdown-box ${countdown.hours === 0 && countdown.minutes < 5 ? 'urgent' : ''}`}>
                                    {String(countdown.seconds).padStart(2, '0')}
                                </div>
                            </div>
                        </div>
                        <Link
                            href="/deals/flash-sale"
                            className="text-white text-sm font-semibold flex items-center gap-1 hover:underline relative z-10 group"
                        >
                            Xem tất cả
                            <span className="group-hover:translate-x-1 transition-transform">→</span>
                        </Link>
                    </div>

                    {/* Products Grid */}
                    <div className="grid grid-cols-3 md:grid-cols-6 gap-3 p-4">
                        {flashSaleProducts.map((product, index) => (
                            <Link
                                key={product.id}
                                href="/deals/flash-sale"
                                className="card-product group animate-slide-up"
                                style={{ animationDelay: `${index * 100}ms` }}
                            >
                                <div className="relative aspect-square bg-gray-50 flex items-center justify-center overflow-hidden">
                                    <span className="text-5xl product-image">{product.image}</span>

                                    {/* Discount badge */}
                                    <span className="absolute top-0 right-0 badge-sale">
                                        -{product.discount}%
                                    </span>

                                    {/* Hot badge */}
                                    {product.sold > 80 && (
                                        <span className="absolute top-0 left-0 badge-hot">
                                            🔥 HOT
                                        </span>
                                    )}
                                </div>
                                <div className="p-3">
                                    <div className="text-[#EE4D2D] font-bold">₫{formatPrice(product.price)}</div>
                                    <div className="text-gray-400 text-xs line-through">₫{formatPrice(product.originalPrice)}</div>

                                    {/* Progress bar */}
                                    <div className="flash-sale-progress mt-2">
                                        <div
                                            className="flash-sale-progress-bar flex items-center justify-center"
                                            style={{ width: `${Math.min(product.sold, 100)}%` }}
                                        >
                                            <span className="text-[10px] text-white font-bold whitespace-nowrap px-1">
                                                {product.sold > 90 ? '🔥 SẮP HẾT!' : `ĐÃ BÁN ${product.sold}`}
                                            </span>
                                        </div>
                                    </div>
                                </div>
                            </Link>
                        ))}
                    </div>
                </div>
            </section>

            {/* Quick Features - Animated Cards */}
            <section className="container mx-auto px-4 py-6">
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    {[
                        { href: '/rewards', icon: '🎮', title: 'Shopee Xu', desc: 'Chơi game nhận xu', color: 'from-yellow-400 via-orange-500 to-red-500' },
                        { href: '/deals/coupons', icon: '🎟️', title: 'Mã Giảm Giá', desc: 'Voucher hot', color: 'from-purple-500 via-pink-500 to-rose-500' },
                        { href: '/live', icon: '🔴', title: 'Shopee Live', desc: 'Xem & mua sắm', color: 'from-red-500 via-rose-500 to-pink-500', live: true },
                        { href: '/deals/flash-sale', icon: '⚡', title: 'Flash Sale', desc: 'Giảm đến 90%', color: 'from-[#EE4D2D] via-[#FF6633] to-[#FFAB91]' },
                    ].map((item, index) => (
                        <Link
                            key={item.href}
                            href={item.href}
                            className={`bg-gradient-to-br ${item.color} rounded-xl p-5 text-white hover-lift hover-shine relative overflow-hidden animate-slide-up`}
                            style={{ animationDelay: `${index * 100}ms` }}
                        >
                            <div className="relative z-10">
                                <div className="text-4xl mb-3 flex items-center gap-2">
                                    {item.icon}
                                    {item.live && <span className="w-2 h-2 bg-white rounded-full animate-pulse" />}
                                </div>
                                <h3 className="font-bold text-lg">{item.title}</h3>
                                <p className="text-sm opacity-90">{item.desc}</p>
                            </div>

                            {/* Decorative blob */}
                            <div className="absolute -bottom-6 -right-6 w-24 h-24 bg-white/10 rounded-full blur-xl" />
                        </Link>
                    ))}
                </div>
            </section>

            {/* Product Suggestions */}
            <section className="container mx-auto px-4 py-6">
                <div className="bg-white rounded-lg shadow-sm overflow-hidden">
                    <div className="p-4 border-b flex items-center justify-between">
                        <h2 className="font-bold text-xl">
                            <span className="text-gradient">GỢI Ý HÔM NAY</span>
                        </h2>
                        <Link href="/products" className="btn-outline px-4 py-1 text-sm group flex items-center gap-1">
                            Xem thêm
                            <span className="group-hover:translate-x-1 transition-transform">→</span>
                        </Link>
                    </div>
                    <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-3 p-4">
                        {[...flashSaleProducts, ...flashSaleProducts].slice(0, 12).map((product, i) => (
                            <Link
                                key={`${product.id}-${i}`}
                                href="/products"
                                className="card-product group animate-slide-up"
                                style={{ animationDelay: `${i * 50}ms` }}
                            >
                                <div className="aspect-square bg-gray-50 flex items-center justify-center relative overflow-hidden">
                                    <span className="text-5xl product-image">{product.image}</span>

                                    {/* Quick add button */}
                                    <button className="absolute bottom-2 right-2 w-8 h-8 bg-[#EE4D2D] text-white rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all duration-300 transform translate-y-2 group-hover:translate-y-0 hover:scale-110">
                                        +
                                    </button>
                                </div>
                                <div className="p-2">
                                    <h3 className="text-xs line-clamp-2 h-8">{product.name}</h3>
                                    <div className="flex items-baseline gap-1 mt-1">
                                        <span className="text-[#EE4D2D] font-bold">₫{formatPrice(product.price)}</span>
                                    </div>
                                    <div className="flex items-center justify-between text-[10px] text-gray-400 mt-1">
                                        <span className="flex items-center gap-0.5">
                                            <span className="text-yellow-400">★</span> 4.9
                                        </span>
                                        <span>Đã bán {product.sold}</span>
                                    </div>
                                </div>
                            </Link>
                        ))}
                    </div>
                </div>
            </section>

            {/* Admin Tools */}
            <section className="container mx-auto px-4 py-6 pb-12">
                <div className="bg-gradient-to-br from-slate-800 to-slate-900 rounded-xl p-6 text-white">
                    <h2 className="font-bold text-xl mb-6 flex items-center gap-2">
                        <span className="animate-pulse">🔧</span> Admin Tools
                    </h2>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        {[
                            { href: '/admin/analytics', icon: '📊', title: 'Analytics', desc: 'Số liệu thời gian thực', color: 'from-blue-500 to-blue-600' },
                            { href: '/admin/fraud', icon: '🛡️', title: 'Fraud Detection', desc: 'ML 99.7% accuracy', color: 'from-slate-500 to-slate-600' },
                            { href: '/admin/pricing', icon: '💹', title: 'Dynamic Pricing', desc: 'AI price optimization', color: 'from-emerald-500 to-emerald-600' },
                        ].map((item, index) => (
                            <Link
                                key={item.href}
                                href={item.href}
                                className={`flex items-center gap-4 p-4 rounded-lg bg-gradient-to-r ${item.color} hover-lift hover-shine relative overflow-hidden animate-slide-up`}
                                style={{ animationDelay: `${index * 100}ms` }}
                            >
                                <span className="text-4xl">{item.icon}</span>
                                <div>
                                    <h3 className="font-bold">{item.title}</h3>
                                    <p className="text-sm opacity-80">{item.desc}</p>
                                </div>
                            </Link>
                        ))}
                    </div>
                </div>
            </section>
        </div>
    );
}
