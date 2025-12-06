'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';

export default function HomePage() {
    const [countdown, setCountdown] = useState({ hours: 2, minutes: 45, seconds: 30 });
    const [stats, setStats] = useState({ users: 12453, orders: 847, revenue: 5.6 });
    const [currentBanner, setCurrentBanner] = useState(0);

    const categories = [
        { icon: '📱', name: 'Điện Thoại', count: 1234 },
        { icon: '💻', name: 'Laptop', count: 567 },
        { icon: '👗', name: 'Thời Trang', count: 4567 },
        { icon: '💄', name: 'Làm Đẹp', count: 2345 },
        { icon: '🏠', name: 'Nhà Cửa', count: 1234 },
        { icon: '🎮', name: 'Gaming', count: 876 },
        { icon: '👟', name: 'Giày Dép', count: 1543 },
        { icon: '⌚', name: 'Đồng Hồ', count: 432 },
        { icon: '📚', name: 'Sách', count: 2134 },
        { icon: '🧸', name: 'Đồ Chơi', count: 987 },
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
        { id: 1, title: '🔥 Flash Sale Đặc Biệt', subtitle: 'Giảm đến 90% - Chỉ hôm nay!', color: 'from-[#EE4D2D] to-[#FF7337]' },
        { id: 2, title: '🎮 Shopee Rewards', subtitle: 'Quay số trúng 1 triệu xu!', color: 'from-yellow-500 to-orange-500' },
        { id: 3, title: '🔴 Shopee Live', subtitle: 'Xem live, mua giá sốc', color: 'from-pink-500 to-red-500' },
    ];

    // Countdown timer
    useEffect(() => {
        const timer = setInterval(() => {
            setCountdown(prev => {
                if (prev.seconds > 0) return { ...prev, seconds: prev.seconds - 1 };
                if (prev.minutes > 0) return { ...prev, minutes: prev.minutes - 1, seconds: 59 };
                if (prev.hours > 0) return { hours: prev.hours - 1, minutes: 59, seconds: 59 };
                return { hours: 2, minutes: 0, seconds: 0 }; // Reset
            });
        }, 1000);
        return () => clearInterval(timer);
    }, []);

    // Live stats update
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
        <div className="min-h-screen bg-[#F5F5F5]">
            {/* Hero Banner Carousel */}
            <section className="relative overflow-hidden">
                <div
                    className="flex transition-transform duration-500"
                    style={{ transform: `translateX(-${currentBanner * 100}%)` }}
                >
                    {banners.map((banner, index) => (
                        <div
                            key={banner.id}
                            className={`min-w-full bg-gradient-to-r ${banner.color} py-12`}
                        >
                            <div className="container mx-auto px-4 text-center text-white">
                                <h1 className="text-4xl md:text-5xl font-bold mb-4">{banner.title}</h1>
                                <p className="text-xl opacity-90">{banner.subtitle}</p>
                            </div>
                        </div>
                    ))}
                </div>

                {/* Banner dots */}
                <div className="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2">
                    {banners.map((_, index) => (
                        <button
                            key={index}
                            onClick={() => setCurrentBanner(index)}
                            className={`w-3 h-3 rounded-full transition-colors ${currentBanner === index ? 'bg-white' : 'bg-white/50'
                                }`}
                        />
                    ))}
                </div>
            </section>

            {/* Live Stats */}
            <section className="bg-gradient-to-r from-slate-800 to-slate-900 py-4">
                <div className="container mx-auto px-4">
                    <div className="flex items-center justify-center gap-8 text-white text-sm">
                        <div className="flex items-center gap-2">
                            <span className="w-2 h-2 bg-green-400 rounded-full animate-pulse" />
                            <span className="text-green-400 font-bold">{stats.users.toLocaleString()}</span>
                            <span className="opacity-75">đang online</span>
                        </div>
                        <div className="flex items-center gap-2">
                            <span className="text-yellow-400 font-bold">{stats.orders.toLocaleString()}</span>
                            <span className="opacity-75">đơn hàng hôm nay</span>
                        </div>
                        <div className="flex items-center gap-2">
                            <span className="text-blue-400 font-bold">{stats.revenue.toFixed(1)}B₫</span>
                            <span className="opacity-75">doanh thu</span>
                        </div>
                    </div>
                </div>
            </section>

            {/* Categories */}
            <section className="container mx-auto px-4 py-6">
                <div className="bg-white rounded">
                    <div className="p-4 border-b flex items-center justify-between">
                        <h2 className="font-bold text-[#EE4D2D] uppercase">Danh mục</h2>
                        <Link href="/products" className="text-[#EE4D2D] text-sm hover:underline">Xem tất cả →</Link>
                    </div>
                    <div className="grid grid-cols-5 md:grid-cols-10 gap-2 p-4">
                        {categories.map(cat => (
                            <Link
                                key={cat.name}
                                href={`/products?category=${encodeURIComponent(cat.name)}`}
                                className="flex flex-col items-center p-2 hover:bg-[#FFEEE8] rounded transition-colors text-center group"
                            >
                                <span className="text-3xl mb-2 group-hover:scale-125 transition-transform">{cat.icon}</span>
                                <span className="text-xs text-gray-700">{cat.name}</span>
                                <span className="text-[10px] text-gray-400">{cat.count}</span>
                            </Link>
                        ))}
                    </div>
                </div>
            </section>

            {/* Flash Sale */}
            <section className="container mx-auto px-4 py-4">
                <div className="bg-white rounded">
                    <div className="p-4 border-b flex items-center justify-between">
                        <div className="flex items-center gap-3">
                            <span className="text-[#EE4D2D] text-xl font-bold uppercase">⚡ Flash Sale</span>
                            <div className="flex items-center gap-1">
                                <span className="bg-[#333] text-white px-2 py-1 rounded text-sm font-mono font-bold">
                                    {String(countdown.hours).padStart(2, '0')}
                                </span>
                                <span className="text-[#333] font-bold">:</span>
                                <span className="bg-[#333] text-white px-2 py-1 rounded text-sm font-mono font-bold">
                                    {String(countdown.minutes).padStart(2, '0')}
                                </span>
                                <span className="text-[#333] font-bold">:</span>
                                <span className="bg-[#333] text-white px-2 py-1 rounded text-sm font-mono font-bold animate-pulse">
                                    {String(countdown.seconds).padStart(2, '0')}
                                </span>
                            </div>
                        </div>
                        <Link href="/deals/flash-sale" className="text-[#EE4D2D] text-sm hover:underline">
                            Xem tất cả →
                        </Link>
                    </div>
                    <div className="grid grid-cols-3 md:grid-cols-6 gap-2 p-4">
                        {flashSaleProducts.map(product => (
                            <Link
                                key={product.id}
                                href={`/deals/flash-sale`}
                                className="group cursor-pointer"
                            >
                                <div className="relative aspect-square bg-gray-100 rounded flex items-center justify-center text-5xl overflow-hidden">
                                    <span className="group-hover:scale-110 transition-transform">{product.image}</span>
                                    <span className="absolute top-0 right-0 bg-[#FFEB3B] text-[#EE4D2D] text-xs font-bold px-1">
                                        -{product.discount}%
                                    </span>
                                </div>
                                <div className="mt-2">
                                    <div className="text-[#EE4D2D] font-bold text-sm">₫{formatPrice(product.price)}</div>
                                    <div className="h-3 bg-[#FFE0DB] rounded-full overflow-hidden">
                                        <div
                                            className="h-full bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] rounded-full relative"
                                            style={{ width: `${Math.min(product.sold, 100)}%` }}
                                        >
                                            <span className="absolute inset-0 flex items-center justify-center text-[8px] text-white font-bold">
                                                {product.sold > 90 ? '🔥 SẮP HẾT' : `ĐÃ BÁN ${product.sold}`}
                                            </span>
                                        </div>
                                    </div>
                                </div>
                            </Link>
                        ))}
                    </div>
                </div>
            </section>

            {/* Quick Features */}
            <section className="container mx-auto px-4 py-4">
                <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                    <Link href="/rewards" className="bg-gradient-to-br from-yellow-400 to-orange-500 rounded p-4 text-white hover:scale-105 transition-transform">
                        <div className="text-3xl mb-2">🎮</div>
                        <h3 className="font-bold">Shopee Xu</h3>
                        <p className="text-xs opacity-90">Chơi game nhận xu</p>
                    </Link>
                    <Link href="/deals/coupons" className="bg-gradient-to-br from-purple-500 to-pink-500 rounded p-4 text-white hover:scale-105 transition-transform">
                        <div className="text-3xl mb-2">🎟️</div>
                        <h3 className="font-bold">Mã Giảm Giá</h3>
                        <p className="text-xs opacity-90">Voucher hot</p>
                    </Link>
                    <Link href="/live" className="bg-gradient-to-br from-red-500 to-pink-600 rounded p-4 text-white hover:scale-105 transition-transform">
                        <div className="text-3xl mb-2 flex items-center gap-2">
                            🔴 <span className="w-2 h-2 bg-white rounded-full animate-pulse" />
                        </div>
                        <h3 className="font-bold">Shopee Live</h3>
                        <p className="text-xs opacity-90">Xem & mua sắm</p>
                    </Link>
                    <Link href="/deals/flash-sale" className="bg-gradient-to-br from-[#EE4D2D] to-[#FF6633] rounded p-4 text-white hover:scale-105 transition-transform">
                        <div className="text-3xl mb-2">⚡</div>
                        <h3 className="font-bold">Flash Sale</h3>
                        <p className="text-xs opacity-90">Giảm đến 90%</p>
                    </Link>
                </div>
            </section>

            {/* Products Grid */}
            <section className="container mx-auto px-4 py-6">
                <div className="bg-white rounded">
                    <div className="p-4 border-b flex items-center justify-between">
                        <h2 className="font-bold text-[#EE4D2D] uppercase">Gợi ý hôm nay</h2>
                        <Link href="/products" className="text-[#EE4D2D] text-sm hover:underline">Xem thêm →</Link>
                    </div>
                    <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-2 p-4">
                        {[...flashSaleProducts, ...flashSaleProducts].slice(0, 12).map((product, i) => (
                            <Link
                                key={`${product.id}-${i}`}
                                href="/products"
                                className="bg-white border rounded-sm overflow-hidden hover:shadow-lg transition-shadow group"
                            >
                                <div className="aspect-square bg-gray-100 flex items-center justify-center text-5xl">
                                    <span className="group-hover:scale-110 transition-transform">{product.image}</span>
                                </div>
                                <div className="p-2">
                                    <h3 className="text-xs line-clamp-2">{product.name}</h3>
                                    <div className="flex items-baseline gap-1 mt-1">
                                        <span className="text-[#EE4D2D] text-sm font-bold">₫{formatPrice(product.price)}</span>
                                    </div>
                                    <div className="flex items-center justify-between text-[10px] text-gray-400 mt-1">
                                        <span>⭐ 4.9</span>
                                        <span>Đã bán {product.sold}</span>
                                    </div>
                                </div>
                            </Link>
                        ))}
                    </div>
                </div>
            </section>

            {/* Admin Quick Access */}
            <section className="container mx-auto px-4 py-4">
                <div className="bg-white rounded p-4">
                    <h2 className="font-bold text-[#EE4D2D] uppercase mb-4">🔧 Admin Tools</h2>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <Link href="/admin/analytics" className="flex items-center gap-3 p-4 bg-blue-50 rounded hover:bg-blue-100 transition-colors">
                            <span className="text-3xl">📊</span>
                            <div>
                                <h3 className="font-bold">Analytics</h3>
                                <p className="text-xs text-gray-500">Số liệu thời gian thực</p>
                            </div>
                        </Link>
                        <Link href="/admin/fraud" className="flex items-center gap-3 p-4 bg-slate-100 rounded hover:bg-slate-200 transition-colors">
                            <span className="text-3xl">🛡️</span>
                            <div>
                                <h3 className="font-bold">Fraud Detection</h3>
                                <p className="text-xs text-gray-500">ML 99.7% accuracy</p>
                            </div>
                        </Link>
                        <Link href="/admin/pricing" className="flex items-center gap-3 p-4 bg-emerald-50 rounded hover:bg-emerald-100 transition-colors">
                            <span className="text-3xl">💹</span>
                            <div>
                                <h3 className="font-bold">Dynamic Pricing</h3>
                                <p className="text-xs text-gray-500">AI price optimization</p>
                            </div>
                        </Link>
                    </div>
                </div>
            </section>
        </div>
    );
}
