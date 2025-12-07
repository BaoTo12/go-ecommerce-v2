'use client';

import React, { useState, useEffect } from 'react';
import Image from 'next/image';
import Link from 'next/link';

interface GroupBuyDeal {
    id: string;
    product: {
        id: string;
        name: string;
        image: string;
        originalPrice: number;
    };
    tiers: { count: number; price: number }[];
    currentCount: number;
    targetCount: number;
    endTime: string;
    participants: { name: string; avatar: string }[];
}

const GROUP_DEALS: GroupBuyDeal[] = [
    {
        id: 'g1',
        product: {
            id: 'p1',
            name: 'iPhone 15 Pro Max 256GB',
            image: 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=400',
            originalPrice: 29990000,
        },
        tiers: [
            { count: 5, price: 28990000 },
            { count: 10, price: 27990000 },
            { count: 20, price: 26990000 },
        ],
        currentCount: 12,
        targetCount: 20,
        endTime: new Date(Date.now() + 12 * 60 * 60 * 1000).toISOString(),
        participants: [
            { name: 'Nguyễn A', avatar: 'https://ui-avatars.com/api/?name=NA' },
            { name: 'Trần B', avatar: 'https://ui-avatars.com/api/?name=TB' },
            { name: 'Lê C', avatar: 'https://ui-avatars.com/api/?name=LC' },
        ],
    },
    {
        id: 'g2',
        product: {
            id: 'p5',
            name: 'Nike Air Force 1 Low',
            image: 'https://images.unsplash.com/photo-1600269452121-4f2416e55c28?w=400',
            originalPrice: 2590000,
        },
        tiers: [
            { count: 10, price: 2390000 },
            { count: 20, price: 2190000 },
            { count: 50, price: 1990000 },
        ],
        currentCount: 35,
        targetCount: 50,
        endTime: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString(),
        participants: [
            { name: 'Hoàng D', avatar: 'https://ui-avatars.com/api/?name=HD' },
            { name: 'Mai E', avatar: 'https://ui-avatars.com/api/?name=ME' },
        ],
    },
    {
        id: 'g3',
        product: {
            id: 'p6',
            name: 'Son Dior Addict Lip Glow',
            image: 'https://images.unsplash.com/photo-1586495777744-4413f21062fa?w=400',
            originalPrice: 950000,
        },
        tiers: [
            { count: 20, price: 890000 },
            { count: 50, price: 850000 },
            { count: 100, price: 790000 },
        ],
        currentCount: 67,
        targetCount: 100,
        endTime: new Date(Date.now() + 6 * 60 * 60 * 1000).toISOString(),
        participants: [
            { name: 'Linh F', avatar: 'https://ui-avatars.com/api/?name=LF' },
            { name: 'Thu G', avatar: 'https://ui-avatars.com/api/?name=TG' },
            { name: 'Hà H', avatar: 'https://ui-avatars.com/api/?name=HH' },
        ],
    },
];

export default function GroupBuyPage() {
    const [deals, setDeals] = useState(GROUP_DEALS);
    const [timeLeft, setTimeLeft] = useState<Record<string, string>>({});
    const [notification, setNotification] = useState<string | null>(null);

    useEffect(() => {
        const timer = setInterval(() => {
            const newTimeLeft: Record<string, string> = {};
            deals.forEach(deal => {
                const end = new Date(deal.endTime).getTime();
                const diff = end - Date.now();
                if (diff <= 0) {
                    newTimeLeft[deal.id] = 'Đã kết thúc';
                } else {
                    const hours = Math.floor(diff / (1000 * 60 * 60));
                    const mins = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
                    newTimeLeft[deal.id] = `${hours}h ${mins}m còn lại`;
                }
            });
            setTimeLeft(newTimeLeft);
        }, 1000);
        return () => clearInterval(timer);
    }, [deals]);

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    const getCurrentPrice = (deal: GroupBuyDeal) => {
        for (let i = deal.tiers.length - 1; i >= 0; i--) {
            if (deal.currentCount >= deal.tiers[i].count) {
                return deal.tiers[i].price;
            }
        }
        return deal.product.originalPrice;
    };

    const joinDeal = (dealId: string) => {
        setDeals(prev => prev.map(d =>
            d.id === dealId ? { ...d, currentCount: d.currentCount + 1 } : d
        ));
        setNotification('🎉 Tham gia mua chung thành công!');
        setTimeout(() => setNotification(null), 3000);
    };

    return (
        <div className="min-h-screen bg-gradient-to-b from-green-500 to-teal-600">
            {notification && (
                <div className="fixed top-20 left-1/2 -translate-x-1/2 bg-white rounded-lg shadow-lg px-6 py-3 z-50 animate-fade-in">
                    {notification}
                </div>
            )}

            {/* Header */}
            <div className="p-6 text-white text-center">
                <h1 className="text-3xl font-bold mb-2">👥 Mua Chung Giá Sốc</h1>
                <p className="opacity-90">Rủ bạn bè mua chung, giá càng rẻ!</p>
            </div>

            <div className="container mx-auto px-4 pb-8">
                {/* How it works */}
                <div className="bg-white/20 backdrop-blur-sm rounded-2xl p-4 mb-6">
                    <div className="flex justify-around text-white text-center text-sm">
                        <div>
                            <div className="text-2xl mb-1">1️⃣</div>
                            <div>Chọn deal</div>
                        </div>
                        <div className="text-2xl">→</div>
                        <div>
                            <div className="text-2xl mb-1">2️⃣</div>
                            <div>Rủ bạn bè</div>
                        </div>
                        <div className="text-2xl">→</div>
                        <div>
                            <div className="text-2xl mb-1">3️⃣</div>
                            <div>Đủ người = Giảm giá!</div>
                        </div>
                    </div>
                </div>

                {/* Deals Grid */}
                <div className="space-y-6">
                    {deals.map((deal, index) => {
                        const currentPrice = getCurrentPrice(deal);
                        const discount = Math.round((1 - currentPrice / deal.product.originalPrice) * 100);
                        const progress = (deal.currentCount / deal.targetCount) * 100;

                        return (
                            <div
                                key={deal.id}
                                className="bg-white rounded-2xl overflow-hidden shadow-lg animate-fade-in-up"
                                style={{ animationDelay: `${index * 100}ms` }}
                            >
                                <div className="flex">
                                    {/* Product Image */}
                                    <div className="w-40 h-40 relative flex-shrink-0">
                                        <Image
                                            src={deal.product.image}
                                            alt={deal.product.name}
                                            fill
                                            className="object-cover"
                                            unoptimized
                                        />
                                        <div className="absolute top-2 left-2 bg-red-500 text-white text-xs px-2 py-1 rounded-full font-bold">
                                            -{discount}%
                                        </div>
                                    </div>

                                    {/* Info */}
                                    <div className="flex-1 p-4">
                                        <h3 className="font-bold line-clamp-2 mb-2">{deal.product.name}</h3>

                                        <div className="flex items-baseline gap-2 mb-2">
                                            <span className="text-xl font-bold text-[#ee4d2d]">₫{formatPrice(currentPrice)}</span>
                                            <span className="text-sm text-gray-400 line-through">₫{formatPrice(deal.product.originalPrice)}</span>
                                        </div>

                                        <div className="text-xs text-gray-500 mb-2">
                                            ⏰ {timeLeft[deal.id] || '...'}
                                        </div>

                                        {/* Price Tiers */}
                                        <div className="flex gap-1 mb-2">
                                            {deal.tiers.map((tier, i) => (
                                                <div
                                                    key={i}
                                                    className={`flex-1 text-center text-xs py-1 rounded ${deal.currentCount >= tier.count
                                                            ? 'bg-green-100 text-green-700'
                                                            : 'bg-gray-100 text-gray-500'
                                                        }`}
                                                >
                                                    {tier.count}+ = ₫{(tier.price / 1000000).toFixed(1)}M
                                                </div>
                                            ))}
                                        </div>

                                        {/* Progress */}
                                        <div className="relative h-2 bg-gray-200 rounded-full overflow-hidden mb-2">
                                            <div
                                                className="absolute h-full bg-gradient-to-r from-green-400 to-green-600 rounded-full transition-all"
                                                style={{ width: `${progress}%` }}
                                            />
                                        </div>

                                        <div className="flex items-center justify-between">
                                            <div className="flex items-center">
                                                <div className="flex -space-x-2">
                                                    {deal.participants.slice(0, 3).map((p, i) => (
                                                        <div key={i} className="w-6 h-6 rounded-full border-2 border-white overflow-hidden relative">
                                                            <img src={p.avatar} alt="" className="w-full h-full object-cover" />
                                                        </div>
                                                    ))}
                                                </div>
                                                <span className="text-xs text-gray-500 ml-2">
                                                    {deal.currentCount}/{deal.targetCount} người
                                                </span>
                                            </div>
                                            <button
                                                onClick={() => joinDeal(deal.id)}
                                                className="px-4 py-2 bg-gradient-to-r from-green-500 to-teal-500 text-white text-sm font-medium rounded-full hover:opacity-90"
                                            >
                                                Tham gia
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        );
                    })}
                </div>
            </div>
        </div>
    );
}
