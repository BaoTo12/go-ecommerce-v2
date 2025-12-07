'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { productService, Product } from '@/services/productService';

export default function FlashSalePage() {
    const [products, setProducts] = useState<Product[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [timeLeft, setTimeLeft] = useState({ hours: 2, minutes: 30, seconds: 45 });

    useEffect(() => {
        productService.getFlashSaleProducts().then(data => {
            setProducts(data);
            setIsLoading(false);
        });
    }, []);

    // Countdown timer
    useEffect(() => {
        const timer = setInterval(() => {
            setTimeLeft(prev => {
                let { hours, minutes, seconds } = prev;
                seconds--;
                if (seconds < 0) {
                    seconds = 59;
                    minutes--;
                }
                if (minutes < 0) {
                    minutes = 59;
                    hours--;
                }
                if (hours < 0) {
                    hours = 23;
                    minutes = 59;
                    seconds = 59;
                }
                return { hours, minutes, seconds };
            });
        }, 1000);
        return () => clearInterval(timer);
    }, []);

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);
    const padZero = (num: number) => num.toString().padStart(2, '0');

    return (
        <div className="min-h-screen bg-[#f5f5f5]">
            {/* Hero Banner */}
            <div className="bg-gradient-to-r from-[#ee4d2d] to-[#ff6633] text-white">
                <div className="container mx-auto px-4 py-8">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-6">
                            <h1 className="text-3xl font-bold flex items-center gap-2 animate-pulse">
                                ⚡ FLASH SALE
                            </h1>
                            <div className="flex items-center gap-2">
                                <span className="text-sm opacity-80">Kết thúc trong</span>
                                <div className="flex gap-1">
                                    <span className="bg-black/30 px-3 py-1 rounded font-mono text-xl font-bold">
                                        {padZero(timeLeft.hours)}
                                    </span>
                                    <span className="text-xl">:</span>
                                    <span className="bg-black/30 px-3 py-1 rounded font-mono text-xl font-bold">
                                        {padZero(timeLeft.minutes)}
                                    </span>
                                    <span className="text-xl">:</span>
                                    <span className="bg-black/30 px-3 py-1 rounded font-mono text-xl font-bold animate-pulse">
                                        {padZero(timeLeft.seconds)}
                                    </span>
                                </div>
                            </div>
                        </div>
                        <Link href="/deals" className="text-sm hover:underline flex items-center gap-1">
                            Xem tất cả <span>→</span>
                        </Link>
                    </div>
                </div>
            </div>

            {/* Time slots */}
            <div className="bg-white shadow-sm">
                <div className="container mx-auto px-4">
                    <div className="flex overflow-x-auto gap-2 py-4">
                        {['00:00', '09:00', '12:00', '15:00', '18:00', '21:00'].map((time, i) => (
                            <button
                                key={time}
                                className={`flex-shrink-0 px-6 py-3 rounded-sm text-sm font-medium transition-all ${i === 2
                                        ? 'bg-[#ee4d2d] text-white'
                                        : i < 2
                                            ? 'bg-gray-100 text-gray-400'
                                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                                    }`}
                            >
                                <div>{time}</div>
                                <div className="text-xs mt-1">
                                    {i < 2 ? 'Đã kết thúc' : i === 2 ? 'Đang diễn ra' : 'Sắp diễn ra'}
                                </div>
                            </button>
                        ))}
                    </div>
                </div>
            </div>

            {/* Products */}
            <div className="container mx-auto px-4 py-6">
                {isLoading ? (
                    <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-3">
                        {[...Array(12)].map((_, i) => (
                            <div key={i} className="bg-white rounded-sm p-3 animate-pulse">
                                <div className="aspect-square bg-gray-200 mb-3" />
                                <div className="h-4 bg-gray-200 rounded mb-2" />
                                <div className="h-4 bg-gray-200 rounded w-2/3" />
                            </div>
                        ))}
                    </div>
                ) : (
                    <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-3">
                        {products.map((product, index) => {
                            const soldPercent = Math.min(Math.floor(Math.random() * 80 + 20), 100);
                            return (
                                <Link
                                    key={product.id}
                                    href={`/products/${product.id}`}
                                    className="product-card group animate-fade-in-up"
                                    style={{ animationDelay: `${index * 30}ms` }}
                                >
                                    <div className="relative aspect-square bg-gray-100 overflow-hidden">
                                        <Image
                                            src={product.thumbnail}
                                            alt={product.name}
                                            fill
                                            className="object-cover product-image"
                                            unoptimized
                                        />
                                        {/* Flash sale badge */}
                                        <div className="absolute top-0 right-0 bg-[#ee4d2d] text-white text-xs font-bold px-2 py-1">
                                            -{product.discount}%
                                        </div>
                                        {/* Sold bar */}
                                        <div className="absolute bottom-0 left-0 right-0 bg-[#ffe8e3] h-6 flex items-center justify-center">
                                            <div
                                                className="absolute left-0 top-0 bottom-0 bg-gradient-to-r from-[#ff5030] to-[#ff9843]"
                                                style={{ width: `${soldPercent}%` }}
                                            />
                                            <span className="relative z-10 text-xs font-medium text-white">
                                                🔥 Đã bán {soldPercent}%
                                            </span>
                                        </div>
                                    </div>
                                    <div className="p-3 text-center">
                                        <span className="text-lg font-bold text-[#ee4d2d]">
                                            ₫{formatPrice(product.price)}
                                        </span>
                                        <div className="text-xs text-gray-400 line-through mt-1">
                                            ₫{formatPrice(product.originalPrice)}
                                        </div>
                                    </div>
                                </Link>
                            );
                        })}
                    </div>
                )}

                {/* Load more */}
                <div className="text-center mt-8">
                    <button className="px-12 py-3 border border-[#ee4d2d] text-[#ee4d2d] hover:bg-[#fef6f5] transition-all">
                        Xem thêm sản phẩm
                    </button>
                </div>
            </div>
        </div>
    );
}
