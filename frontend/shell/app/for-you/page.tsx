'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { recommendationService, RecommendationResult } from '@/services/recommendationService';
import { recentlyViewedService } from '@/services/recentlyViewedService';

export default function PersonalizedHome() {
    const [recommendations, setRecommendations] = useState<{
        personalized: RecommendationResult | null;
        trending: RecommendationResult | null;
        becauseViewed: RecommendationResult | null;
    }>({ personalized: null, trending: null, becauseViewed: null });
    const [isLoading, setIsLoading] = useState(true);
    const [userName] = useState('Nguyễn Văn A');

    useEffect(() => {
        const loadRecommendations = async () => {
            setIsLoading(true);
            const history = await recentlyViewedService.getRecentlyViewed();

            const [personalized, trending] = await Promise.all([
                recommendationService.getPersonalized(),
                recommendationService.getTrendingNearYou('TP. Hồ Chí Minh'),
            ]);

            let becauseViewed = null;
            if (history.length > 0) {
                becauseViewed = await recommendationService.getBecauseYouViewed(history[0].product.id);
            }

            setRecommendations({ personalized, trending, becauseViewed });
            setIsLoading(false);
        };

        loadRecommendations();
    }, []);

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    const ProductSection = ({ result, icon }: { result: RecommendationResult | null; icon: string }) => {
        if (!result || result.products.length === 0) return null;

        return (
            <div className="mb-8 animate-fade-in-up">
                <div className="flex items-center justify-between mb-4">
                    <h2 className="text-lg font-bold dark:text-white flex items-center gap-2">
                        {icon} {result.reason}
                    </h2>
                    <Link href="/products" className="text-[#ee4d2d] text-sm">Xem thêm →</Link>
                </div>
                <div className="flex gap-3 overflow-x-auto pb-2">
                    {result.products.map((product, index) => (
                        <Link
                            key={product.id}
                            href={`/products/${product.id}`}
                            className="flex-shrink-0 w-36 bg-white dark:bg-gray-800 rounded-lg shadow-sm overflow-hidden animate-fade-in-up"
                            style={{ animationDelay: `${index * 50}ms` }}
                        >
                            <div className="aspect-square relative bg-gray-100 dark:bg-gray-700">
                                <Image
                                    src={product.thumbnail}
                                    alt={product.name}
                                    fill
                                    className="object-cover"
                                    unoptimized
                                />
                                {product.discount > 0 && (
                                    <span className="absolute top-1 left-1 px-1.5 py-0.5 bg-[#ee4d2d] text-white text-xs rounded">
                                        -{product.discount}%
                                    </span>
                                )}
                            </div>
                            <div className="p-2">
                                <h3 className="text-xs line-clamp-2 h-8 dark:text-white">{product.name}</h3>
                                <p className="text-sm font-bold text-[#ee4d2d] mt-1">₫{formatPrice(product.price)}</p>
                            </div>
                        </Link>
                    ))}
                </div>
            </div>
        );
    };

    if (isLoading) {
        return (
            <div className="min-h-screen bg-[#f5f5f5] dark:bg-gray-900 p-4">
                <div className="animate-pulse space-y-8">
                    {[1, 2, 3].map(i => (
                        <div key={i}>
                            <div className="h-6 bg-gray-200 dark:bg-gray-700 rounded w-1/3 mb-4" />
                            <div className="flex gap-3">
                                {[1, 2, 3, 4].map(j => (
                                    <div key={j} className="w-36 h-48 bg-gray-200 dark:bg-gray-700 rounded-lg" />
                                ))}
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-[#f5f5f5] dark:bg-gray-900">
            {/* Welcome Banner */}
            <div className="bg-gradient-to-r from-[#ee4d2d] to-[#ff6633] p-6 text-white">
                <h1 className="text-2xl font-bold mb-1">Xin chào, {userName}! 👋</h1>
                <p className="opacity-90">Khám phá những gợi ý dành riêng cho bạn</p>
            </div>

            <div className="container mx-auto px-4 py-6">
                {/* Quick Actions */}
                <div className="grid grid-cols-4 gap-3 mb-8">
                    {[
                        { icon: '📷', label: 'Tìm bằng ảnh', href: '/image-search' },
                        { icon: '🎮', label: 'Mini Games', href: '/games' },
                        { icon: '🔥', label: 'Flash Sale', href: '/deals/flash-sale' },
                        { icon: '🎁', label: 'Shopee Xu', href: '/rewards' },
                    ].map((action, i) => (
                        <Link
                            key={i}
                            href={action.href}
                            className="bg-white dark:bg-gray-800 rounded-xl p-3 text-center shadow-sm hover:shadow-md transition-shadow"
                        >
                            <div className="text-2xl mb-1">{action.icon}</div>
                            <div className="text-xs dark:text-white">{action.label}</div>
                        </Link>
                    ))}
                </div>

                {/* Personalized Recommendations */}
                <ProductSection result={recommendations.personalized} icon="✨" />

                {/* Because you viewed */}
                <ProductSection result={recommendations.becauseViewed} icon="👀" />

                {/* Trending near you */}
                <ProductSection result={recommendations.trending} icon="📍" />

                {/* AI Confidence Indicator */}
                <div className="bg-gradient-to-r from-purple-100 to-pink-100 dark:from-purple-900 dark:to-pink-900 rounded-xl p-4 text-center">
                    <p className="text-sm text-gray-600 dark:text-gray-300">
                        🤖 Gợi ý được cá nhân hóa bởi AI dựa trên lịch sử duyệt web của bạn
                    </p>
                    <Link href="/account/settings/privacy" className="text-xs text-[#ee4d2d] mt-2 inline-block">
                        Quản lý cài đặt riêng tư →
                    </Link>
                </div>
            </div>
        </div>
    );
}
