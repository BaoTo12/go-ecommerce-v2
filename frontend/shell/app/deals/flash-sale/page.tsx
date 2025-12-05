'use client';

import React, { useState, useEffect } from 'react';

interface FlashSale {
    id: string;
    product_id: string;
    name: string;
    image: string;
    original_price: number;
    sale_price: number;
    discount_percent: number;
    total_quantity: number;
    sold_quantity: number;
    end_time: string;
}

export default function FlashSalePage() {
    const [sales, setSales] = useState<FlashSale[]>([]);
    const [loading, setLoading] = useState(true);
    const [timeLeft, setTimeLeft] = useState<Record<string, string>>({});

    useEffect(() => {
        // Mock Shopee flash sale data
        setSales([
            {
                id: 'fs-001',
                product_id: 'p-001',
                name: 'iPhone 15 Pro Max 256GB',
                image: '📱',
                original_price: 34990000,
                sale_price: 29990000,
                discount_percent: 14,
                total_quantity: 100,
                sold_quantity: 87,
                end_time: new Date(Date.now() + 7200000).toISOString(),
            },
            {
                id: 'fs-002',
                product_id: 'p-002',
                name: 'Samsung Galaxy S24 Ultra',
                image: '📲',
                original_price: 25990000,
                sale_price: 19990000,
                discount_percent: 23,
                total_quantity: 50,
                sold_quantity: 48,
                end_time: new Date(Date.now() + 7200000).toISOString(),
            },
            {
                id: 'fs-003',
                product_id: 'p-003',
                name: 'MacBook Air M3 13"',
                image: '💻',
                original_price: 27990000,
                sale_price: 24990000,
                discount_percent: 11,
                total_quantity: 30,
                sold_quantity: 21,
                end_time: new Date(Date.now() + 7200000).toISOString(),
            },
            {
                id: 'fs-004',
                product_id: 'p-004',
                name: 'Sony WH-1000XM5',
                image: '🎧',
                original_price: 8990000,
                sale_price: 6990000,
                discount_percent: 22,
                total_quantity: 200,
                sold_quantity: 156,
                end_time: new Date(Date.now() + 7200000).toISOString(),
            },
            {
                id: 'fs-005',
                product_id: 'p-005',
                name: 'iPad Pro M4 11"',
                image: '📟',
                original_price: 23990000,
                sale_price: 21990000,
                discount_percent: 8,
                total_quantity: 80,
                sold_quantity: 45,
                end_time: new Date(Date.now() + 7200000).toISOString(),
            },
            {
                id: 'fs-006',
                product_id: 'p-006',
                name: 'Apple Watch Ultra 2',
                image: '⌚',
                original_price: 21990000,
                sale_price: 18990000,
                discount_percent: 14,
                total_quantity: 60,
                sold_quantity: 59,
                end_time: new Date(Date.now() + 7200000).toISOString(),
            },
        ]);
        setLoading(false);
    }, []);

    useEffect(() => {
        const timer = setInterval(() => {
            const newTimeLeft: Record<string, string> = {};
            sales.forEach(sale => {
                const diff = new Date(sale.end_time).getTime() - Date.now();
                if (diff > 0) {
                    const h = Math.floor(diff / 3600000);
                    const m = Math.floor((diff % 3600000) / 60000);
                    const s = Math.floor((diff % 60000) / 1000);
                    newTimeLeft[sale.id] = `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
                } else {
                    newTimeLeft[sale.id] = '00:00:00';
                }
            });
            setTimeLeft(newTimeLeft);
        }, 1000);
        return () => clearInterval(timer);
    }, [sales]);

    const formatPrice = (price: number) => {
        return new Intl.NumberFormat('vi-VN').format(price);
    };

    if (loading) {
        return (
            <div className="flex h-96 items-center justify-center">
                <div className="text-2xl text-[#EE4D2D]">⚡ Đang tải...</div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-[#F5F5F5]">
            {/* Shopee Flash Sale Header */}
            <div className="bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] py-4">
                <div className="container mx-auto px-4">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-4">
                            <div className="flex items-center gap-2">
                                <span className="text-white text-2xl font-bold">⚡</span>
                                <span className="text-white text-2xl font-bold uppercase tracking-wider">Flash Sale</span>
                            </div>
                            <div className="flex items-center gap-1 ml-4">
                                {timeLeft[sales[0]?.id]?.split(':').map((unit, i) => (
                                    <React.Fragment key={i}>
                                        <span className="bg-[#333] text-white px-2 py-1 rounded text-lg font-bold font-mono">
                                            {unit}
                                        </span>
                                        {i < 2 && <span className="text-white">:</span>}
                                    </React.Fragment>
                                ))}
                            </div>
                        </div>
                        <a href="#" className="text-white hover:underline text-sm">
                            Xem tất cả &gt;
                        </a>
                    </div>
                </div>
            </div>

            {/* Products Grid */}
            <div className="container mx-auto px-4 py-6">
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-2">
                    {sales.map(sale => {
                        const soldPercent = Math.round((sale.sold_quantity / sale.total_quantity) * 100);
                        const isAlmostGone = soldPercent >= 90;

                        return (
                            <div
                                key={sale.id}
                                className="bg-white rounded-sm overflow-hidden hover:shadow-lg transition-shadow cursor-pointer border border-transparent hover:border-[#EE4D2D]"
                            >
                                {/* Discount Badge */}
                                <div className="relative">
                                    <div className="absolute top-0 right-0 bg-[#FFEB3B] text-[#EE4D2D] px-2 py-1 text-xs font-bold">
                                        -{sale.discount_percent}%
                                    </div>
                                    <div className="aspect-square bg-gradient-to-br from-gray-50 to-gray-100 flex items-center justify-center text-6xl p-4">
                                        {sale.image}
                                    </div>
                                </div>

                                {/* Product Info */}
                                <div className="p-2">
                                    <h3 className="text-sm line-clamp-2 h-10 mb-2">{sale.name}</h3>

                                    {/* Price */}
                                    <div className="flex items-baseline gap-1 mb-2">
                                        <span className="text-[#EE4D2D] font-bold">₫</span>
                                        <span className="text-[#EE4D2D] text-lg font-bold">
                                            {formatPrice(sale.sale_price)}
                                        </span>
                                    </div>

                                    {/* Progress Bar */}
                                    <div className="relative h-4 bg-[#FFE0DB] rounded-full overflow-hidden">
                                        <div
                                            className="absolute inset-0 bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] rounded-full transition-all"
                                            style={{ width: `${soldPercent}%` }}
                                        >
                                            {/* Shimmer effect */}
                                            <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/30 to-transparent animate-[shimmer_2s_infinite]" />
                                        </div>
                                        <span className="absolute inset-0 flex items-center justify-center text-[10px] font-bold text-white drop-shadow">
                                            {isAlmostGone ? '🔥 SẮP HẾT' : `ĐÃ BÁN ${sale.sold_quantity}`}
                                        </span>
                                    </div>
                                </div>
                            </div>
                        );
                    })}
                </div>
            </div>

            {/* Why Shopee Flash Sale */}
            <div className="container mx-auto px-4 py-8">
                <div className="bg-white rounded-sm p-6">
                    <h2 className="text-lg font-bold text-center mb-6">Ưu đãi Flash Sale</h2>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
                        <div className="p-4">
                            <div className="text-4xl mb-2">⚡</div>
                            <h3 className="font-semibold text-sm">Giá Shock</h3>
                            <p className="text-xs text-gray-500">Giảm đến 90%</p>
                        </div>
                        <div className="p-4">
                            <div className="text-4xl mb-2">🛡️</div>
                            <h3 className="font-semibold text-sm">Chống BOT</h3>
                            <p className="text-xs text-gray-500">Bảo vệ công bằng</p>
                        </div>
                        <div className="p-4">
                            <div className="text-4xl mb-2">🚚</div>
                            <h3 className="font-semibold text-sm">Freeship</h3>
                            <p className="text-xs text-gray-500">Miễn phí vận chuyển</p>
                        </div>
                        <div className="p-4">
                            <div className="text-4xl mb-2">✅</div>
                            <h3 className="font-semibold text-sm">Chính Hãng</h3>
                            <p className="text-xs text-gray-500">100% authentic</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
