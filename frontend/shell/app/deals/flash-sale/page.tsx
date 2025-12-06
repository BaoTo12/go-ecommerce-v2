'use client';

import React, { useState, useEffect, useRef } from 'react';
import Link from 'next/link';

interface FlashSaleProduct {
    id: number;
    name: string;
    price: number;
    originalPrice: number;
    discount: number;
    sold: number;
    total: number;
    image: string;
    isHot?: boolean;
}

export default function FlashSalePage() {
    const [countdown, setCountdown] = useState({ hours: 2, minutes: 45, seconds: 30 });
    const [products, setProducts] = useState<FlashSaleProduct[]>([]);
    const [addedToCart, setAddedToCart] = useState<number | null>(null);
    const [isLoaded, setIsLoaded] = useState(false);
    const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 });

    // Create ref for the header section
    const headerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        setIsLoaded(true);
        setProducts([
            { id: 1, name: 'iPhone 15 Pro Max 256GB Xanh Titan', price: 29990000, originalPrice: 34990000, discount: 14, sold: 87, total: 100, image: '📱', isHot: true },
            { id: 2, name: 'AirPods Pro 2 USB-C Wireless', price: 4990000, originalPrice: 6990000, discount: 29, sold: 156, total: 200, image: '🎧', isHot: true },
            { id: 3, name: 'MacBook Air M3 13" 256GB', price: 24990000, originalPrice: 27990000, discount: 11, sold: 45, total: 80, image: '💻' },
            { id: 4, name: 'Samsung Galaxy Watch 6 Classic', price: 6990000, originalPrice: 8990000, discount: 22, sold: 78, total: 100, image: '⌚' },
            { id: 5, name: 'iPad Pro M2 11" 128GB WiFi', price: 21990000, originalPrice: 25990000, discount: 15, sold: 34, total: 50, image: '📟' },
            { id: 6, name: 'Sony WH-1000XM5 Wireless', price: 6990000, originalPrice: 8990000, discount: 22, sold: 92, total: 100, image: '🎧', isHot: true },
            { id: 7, name: 'Nintendo Switch OLED White', price: 7490000, originalPrice: 9490000, discount: 21, sold: 67, total: 80, image: '🎮' },
            { id: 8, name: 'DJI Mini 3 Pro Drone', price: 18990000, originalPrice: 22990000, discount: 17, sold: 23, total: 30, image: '🚁' },
            { id: 9, name: 'GoPro Hero 12 Black', price: 10990000, originalPrice: 12990000, discount: 15, sold: 45, total: 60, image: '📷' },
            { id: 10, name: 'Dyson V15 Detect Vacuum', price: 15990000, originalPrice: 19990000, discount: 20, sold: 55, total: 70, image: '🧹', isHot: true },
            { id: 11, name: 'LG OLED C3 55" 4K TV', price: 32990000, originalPrice: 39990000, discount: 18, sold: 28, total: 40, image: '📺' },
            { id: 12, name: 'Bose QuietComfort Ultra', price: 8990000, originalPrice: 10990000, discount: 18, sold: 72, total: 90, image: '🎧' },
        ]);
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

    // Mouse parallax effect
    useEffect(() => {
        const handleMouseMove = (e: MouseEvent) => {
            if (headerRef.current) {
                const rect = headerRef.current.getBoundingClientRect();
                setMousePosition({
                    x: (e.clientX - rect.left - rect.width / 2) / 20,
                    y: (e.clientY - rect.top - rect.height / 2) / 20,
                });
            }
        };

        window.addEventListener('mousemove', handleMouseMove);
        return () => window.removeEventListener('mousemove', handleMouseMove);
    }, []);

    const addToCart = (productId: number) => {
        setAddedToCart(productId);
        setTimeout(() => setAddedToCart(null), 2000);
    };

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    const getSoldPercentage = (sold: number, total: number) => Math.min((sold / total) * 100, 100);

    return (
        <div className={`min-h-screen bg-[#F5F5F5] ${isLoaded ? 'animate-fade-in' : 'opacity-0'}`}>
            {/* Toast notification */}
            {addedToCart && (
                <div className="toast toast-success">
                    <span className="text-xl mr-2">✓</span>
                    Đã thêm vào giỏ hàng!
                </div>
            )}

            {/* Animated Header */}
            <div
                ref={headerRef}
                className="relative overflow-hidden bg-gradient-to-r from-[#EE4D2D] via-[#FF6633] to-[#FFAB91]"
            >
                {/* Animated background elements */}
                <div className="absolute inset-0 overflow-hidden">
                    {/* Floating circles */}
                    {[...Array(20)].map((_, i) => (
                        <div
                            key={i}
                            className="absolute rounded-full bg-white/10 animate-float"
                            style={{
                                width: Math.random() * 60 + 20,
                                height: Math.random() * 60 + 20,
                                left: `${Math.random() * 100}%`,
                                top: `${Math.random() * 100}%`,
                                animationDelay: `${Math.random() * 3}s`,
                                animationDuration: `${3 + Math.random() * 2}s`,
                            }}
                        />
                    ))}

                    {/* Lightning bolts */}
                    <div className="absolute top-4 left-[10%] text-6xl animate-bounce opacity-50" style={{ animationDelay: '0s' }}>⚡</div>
                    <div className="absolute top-8 right-[15%] text-4xl animate-bounce opacity-50" style={{ animationDelay: '0.5s' }}>⚡</div>
                    <div className="absolute bottom-4 left-[30%] text-5xl animate-bounce opacity-50" style={{ animationDelay: '1s' }}>⚡</div>
                </div>

                <div className="container mx-auto px-4 py-12 relative z-10">
                    <div className="text-center text-white">
                        <h1
                            className="text-5xl md:text-7xl font-black mb-4 drop-shadow-lg"
                            style={{ transform: `translate(${mousePosition.x}px, ${mousePosition.y}px)` }}
                        >
                            <span className="inline-block animate-bounce mr-2">⚡</span>
                            FLASH SALE
                            <span className="inline-block animate-bounce ml-2" style={{ animationDelay: '0.2s' }}>⚡</span>
                        </h1>
                        <p className="text-xl md:text-2xl opacity-90 mb-8 font-medium">
                            Săn deal shock - Giảm đến 90%!
                        </p>

                        {/* Countdown */}
                        <div className="flex items-center justify-center gap-3">
                            <span className="text-lg font-semibold">Kết thúc sau:</span>
                            <div className="flex items-center gap-2">
                                <div className="countdown-box text-2xl md:text-3xl">
                                    {String(countdown.hours).padStart(2, '0')}
                                </div>
                                <span className="text-3xl font-bold animate-pulse">:</span>
                                <div className="countdown-box text-2xl md:text-3xl">
                                    {String(countdown.minutes).padStart(2, '0')}
                                </div>
                                <span className="text-3xl font-bold animate-pulse">:</span>
                                <div className="countdown-box text-2xl md:text-3xl animate-pulse">
                                    {String(countdown.seconds).padStart(2, '0')}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Wave divider */}
                <div className="absolute bottom-0 left-0 right-0">
                    <svg viewBox="0 0 1440 60" fill="none" className="w-full">
                        <path d="M0 60L48 55C96 50 192 40 288 35C384 30 480 30 576 32.5C672 35 768 40 864 42.5C960 45 1056 45 1152 42.5C1248 40 1344 35 1392 32.5L1440 30V60H0Z" fill="#F5F5F5" />
                    </svg>
                </div>
            </div>

            {/* Products Grid */}
            <div className="container mx-auto px-4 py-8">
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
                    {products.map((product, index) => {
                        const soldPercent = getSoldPercentage(product.sold, product.total);
                        const isAlmostGone = soldPercent > 80;

                        return (
                            <div
                                key={product.id}
                                className="card-product group animate-slide-up relative"
                                style={{ animationDelay: `${index * 50}ms` }}
                            >
                                {/* Image */}
                                <div className="relative aspect-square bg-gradient-to-br from-gray-50 to-gray-100 flex items-center justify-center overflow-hidden">
                                    <span className="text-7xl product-image drop-shadow-lg">{product.image}</span>

                                    {/* Discount badge */}
                                    <div className="absolute top-0 right-0">
                                        <div className="badge-sale py-1 px-2 text-sm font-black">
                                            -{product.discount}%
                                        </div>
                                    </div>

                                    {/* Hot badge */}
                                    {product.isHot && (
                                        <div className="absolute top-0 left-0">
                                            <div className="badge-hot py-1 px-2 text-sm flex items-center gap-1">
                                                <span className="animate-bounce">🔥</span> HOT
                                            </div>
                                        </div>
                                    )}

                                    {/* Quick add button */}
                                    <button
                                        onClick={() => addToCart(product.id)}
                                        className="absolute bottom-3 right-3 w-10 h-10 bg-[#EE4D2D] text-white rounded-full flex items-center justify-center 
                               opacity-0 group-hover:opacity-100 transition-all duration-300 
                               transform translate-y-4 group-hover:translate-y-0 
                               hover:scale-110 hover:bg-[#D73211] shadow-lg ripple"
                                    >
                                        🛒
                                    </button>
                                </div>

                                {/* Info */}
                                <div className="p-3">
                                    <h3 className="text-sm line-clamp-2 h-10 text-gray-800 group-hover:text-[#EE4D2D] transition-colors">
                                        {product.name}
                                    </h3>

                                    <div className="flex items-baseline gap-2 mt-2">
                                        <span className="text-[#EE4D2D] text-lg font-bold">
                                            ₫{formatPrice(product.price)}
                                        </span>
                                    </div>
                                    <div className="text-gray-400 text-xs line-through">
                                        ₫{formatPrice(product.originalPrice)}
                                    </div>

                                    {/* Progress bar */}
                                    <div className="mt-3">
                                        <div className={`flash-sale-progress ${isAlmostGone ? 'animate-pulse-glow' : ''}`}>
                                            <div
                                                className="flash-sale-progress-bar flex items-center justify-center"
                                                style={{ width: `${soldPercent}%` }}
                                            >
                                                <span className="text-[10px] text-white font-bold whitespace-nowrap drop-shadow">
                                                    {isAlmostGone ? '🔥 SẮP HẾT!' : `Đã bán ${product.sold}`}
                                                </span>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        );
                    })}
                </div>
            </div>

            {/* Features Section */}
            <div className="container mx-auto px-4 py-12">
                <div className="bg-white rounded-xl p-8 shadow-sm">
                    <h2 className="text-2xl font-bold text-center mb-8">
                        <span className="text-gradient">Tại sao chọn Flash Sale?</span>
                    </h2>
                    <div className="grid md:grid-cols-4 gap-6">
                        {[
                            { icon: '💰', title: 'Giá Shock', desc: 'Giảm đến 90% mỗi ngày', color: 'from-yellow-400 to-orange-500' },
                            { icon: '🤖', title: 'Chống BOT', desc: 'Công nghệ chống gian lận', color: 'from-blue-400 to-indigo-500' },
                            { icon: '🚚', title: 'Freeship', desc: 'Miễn phí vận chuyển', color: 'from-green-400 to-emerald-500' },
                            { icon: '✅', title: 'Chính Hãng', desc: '100% hàng chính hãng', color: 'from-purple-400 to-pink-500' },
                        ].map((feature, index) => (
                            <div
                                key={feature.title}
                                className="text-center group hover-lift animate-slide-up"
                                style={{ animationDelay: `${index * 100}ms` }}
                            >
                                <div className={`w-20 h-20 mx-auto rounded-2xl bg-gradient-to-br ${feature.color} flex items-center justify-center text-4xl shadow-lg group-hover:scale-110 transition-transform`}>
                                    {feature.icon}
                                </div>
                                <h3 className="font-bold mt-4 text-gray-800">{feature.title}</h3>
                                <p className="text-sm text-gray-500 mt-1">{feature.desc}</p>
                            </div>
                        ))}
                    </div>
                </div>
            </div>

            {/* Floating Action Button */}
            <Link
                href="/cart"
                className="fixed bottom-6 right-6 z-50 w-14 h-14 bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] text-white rounded-full shadow-xl flex items-center justify-center text-2xl hover:scale-110 transition-transform animate-bounce-slow neon-orange"
            >
                🛒
            </Link>
        </div>
    );
}
