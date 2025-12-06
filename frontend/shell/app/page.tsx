'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { productService, Product, Category } from '@/services/productService';

export default function HomePage() {
    const [countdown, setCountdown] = useState({ hours: 2, minutes: 45, seconds: 30 });
    const [currentBanner, setCurrentBanner] = useState(0);
    const [isLoaded, setIsLoaded] = useState(false);
    const [categories, setCategories] = useState<Category[]>([]);
    const [flashSaleProducts, setFlashSaleProducts] = useState<Product[]>([]);
    const [recommendedProducts, setRecommendedProducts] = useState<Product[]>([]);
    const [hoveredCategory, setHoveredCategory] = useState<number | null>(null);

    const banners = [
        { bg: 'linear-gradient(135deg, #f53d2d 0%, #ff6533 100%)', text: '🔥 Flash Sale 12.12', subtitle: 'Giảm đến 50%!' },
        { bg: 'linear-gradient(135deg, #00bfa5 0%, #00897b 100%)', text: '🚚 Miễn Phí Vận Chuyển', subtitle: 'Đơn từ 0Đ' },
        { bg: 'linear-gradient(135deg, #5c6bc0 0%, #3949ab 100%)', text: '💳 Hoàn Tiền 10%', subtitle: 'Qua ShopeePay' },
    ];

    // Load data from service
    useEffect(() => {
        const loadData = async () => {
            try {
                const [cats, flash, recommended] = await Promise.all([
                    productService.getCategories(),
                    productService.getFlashSaleProducts(),
                    productService.getRecommendedProducts(),
                ]);
                setCategories(cats);
                setFlashSaleProducts(flash);
                setRecommendedProducts(recommended);
                setIsLoaded(true);
            } catch (error) {
                console.error('Failed to load data:', error);
                setIsLoaded(true);
            }
        };
        loadData();
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
                    <div className="grid grid-cols-5 md:grid-cols-10">
                        {categories.map((cat, index) => (
                            <Link
                                key={cat.id}
                                href={`/products?category=${encodeURIComponent(cat.name)}`}
                                className="category-item animate-fade-in-up"
                                style={{ animationDelay: `${index * 30}ms` }}
                                onMouseEnter={() => setHoveredCategory(index)}
                                onMouseLeave={() => setHoveredCategory(null)}
                            >
                                <div className={`w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center mb-2 overflow-hidden category-icon transition-transform duration-300 ${hoveredCategory === index ? 'animate-wiggle' : ''
                                    }`}>
                                    {cat.image ? (
                                        <Image
                                            src={cat.image}
                                            alt={cat.name}
                                            width={48}
                                            height={48}
                                            className="w-full h-full object-cover"
                                            unoptimized
                                        />
                                    ) : (
                                        <span className="text-2xl">{cat.icon}</span>
                                    )}
                                </div>
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
                                href={`/products/${product.id}`}
                                className="product-card border-0 rounded-none animate-fade-in-up hover-shine"
                                style={{ animationDelay: `${index * 100}ms` }}
                            >
                                <div className="relative aspect-square bg-gray-50 flex items-center justify-center overflow-hidden">
                                    <Image
                                        src={product.thumbnail}
                                        alt={product.name}
                                        fill
                                        className="object-cover product-image"
                                        unoptimized
                                    />
                                    <div className="discount-badge">-{product.discount}%</div>
                                </div>
                                <div className="p-2 text-center">
                                    <div className="price-current text-lg font-bold">₫{formatPrice(product.price)}</div>
                                    <div className="flash-progress mt-2">
                                        <div className="flash-progress-bar" style={{ width: `${Math.min(90, (product.sold / 1000) * 5)}%` }}>
                                            <span>Đã bán {product.soldDisplay}</span>
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
                        { href: '/rewards', icon: '🎮', title: 'Shopee Xu', desc: 'Đổi xu lấy quà' },
                        { href: '/deals/coupons', icon: '🎟️', title: 'Mã Giảm Giá', desc: 'Săn voucher hot' },
                        { href: '/live', icon: '🔴', title: 'Shopee Live', desc: 'Xem & mua sắm', live: true },
                        { href: '/products', icon: '🛍️', title: 'Hàng Chính Hãng', desc: '100% authentic' },
                    ].map((item, index) => (
                        <Link
                            key={item.href}
                            href={item.href}
                            className="bg-white rounded-sm p-4 flex items-center gap-3 hover-lift hover-shine animate-fade-in-up"
                            style={{ animationDelay: `${index * 100}ms` }}
                        >
                            <span className={`text-3xl ${item.live ? 'animate-pulse' : ''}`}>{item.icon}</span>
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
                    {/* Fixed header - no longer sticky to fix the bug */}
                    <div className="p-4 border-b text-center bg-white">
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
                                <div className="relative aspect-square bg-gray-50 overflow-hidden">
                                    <Image
                                        src={product.thumbnail}
                                        alt={product.name}
                                        fill
                                        className="object-cover product-image"
                                        unoptimized
                                    />
                                    {product.isOfficial && (
                                        <div className="absolute top-0 left-0 bg-[#ee4d2d] text-white text-[10px] px-1 py-0.5">
                                            Mall
                                        </div>
                                    )}
                                    {product.isFavorite && (
                                        <div className="absolute top-0 left-0 bg-[#ee4d2d] text-white text-[10px] px-1 py-0.5">
                                            Yêu thích
                                        </div>
                                    )}
                                </div>
                                <div className="p-2">
                                    <h3 className="text-xs line-clamp-2 h-8 mb-2">{product.name}</h3>
                                    <div className="flex items-center justify-between">
                                        <span className="price-current">₫{formatPrice(product.price)}</span>
                                        <span className="text-xs text-gray-500">Đã bán {product.soldDisplay}</span>
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
