'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import Image from 'next/image';

interface AffiliateStats {
    totalClicks: number;
    totalOrders: number;
    totalEarnings: number;
    pendingEarnings: number;
    conversionRate: number;
}

interface AffiliateProduct {
    id: string;
    name: string;
    image: string;
    price: number;
    commission: number;
    clicks: number;
    orders: number;
    earnings: number;
}

const MOCK_STATS: AffiliateStats = {
    totalClicks: 12580,
    totalOrders: 342,
    totalEarnings: 15680000,
    pendingEarnings: 2340000,
    conversionRate: 2.72,
};

const MOCK_PRODUCTS: AffiliateProduct[] = [
    {
        id: 'aff1',
        name: 'iPhone 15 Pro Max 256GB',
        image: 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=200',
        price: 29990000,
        commission: 3,
        clicks: 2340,
        orders: 45,
        earnings: 4048650,
    },
    {
        id: 'aff2',
        name: 'Son Dior Addict Lip Glow',
        image: 'https://images.unsplash.com/photo-1586495777744-4413f21062fa?w=200',
        price: 950000,
        commission: 8,
        clicks: 3456,
        orders: 89,
        earnings: 676400,
    },
    {
        id: 'aff3',
        name: 'Nike Air Force 1 Low',
        image: 'https://images.unsplash.com/photo-1600269452121-4f2416e55c28?w=200',
        price: 2590000,
        commission: 5,
        clicks: 1890,
        orders: 34,
        earnings: 440300,
    },
];

export default function AffiliateDashboard() {
    const [stats] = useState(MOCK_STATS);
    const [products] = useState(MOCK_PRODUCTS);
    const [selectedProduct, setSelectedProduct] = useState<string | null>(null);
    const [copied, setCopied] = useState(false);

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);
    const formatNumber = (num: number) => new Intl.NumberFormat('vi-VN').format(num);

    const generateAffiliateLink = (productId: string) => {
        return `https://shopee.vn/product/${productId}?aff=USER123`;
    };

    const copyLink = (productId: string) => {
        navigator.clipboard.writeText(generateAffiliateLink(productId));
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    return (
        <div className="min-h-screen bg-gray-100 dark:bg-gray-900">
            {/* Header */}
            <div className="bg-gradient-to-r from-indigo-600 to-purple-600 p-6">
                <div className="container mx-auto">
                    <h1 className="text-2xl font-bold text-white">💰 Affiliate Dashboard</h1>
                    <p className="text-white/80">Kiếm tiền bằng cách chia sẻ sản phẩm</p>
                </div>
            </div>

            <div className="container mx-auto px-4 py-6">
                {/* Stats Grid */}
                <div className="grid grid-cols-2 md:grid-cols-5 gap-4 mb-6">
                    <div className="bg-white dark:bg-gray-800 rounded-xl p-4 shadow-sm">
                        <div className="text-gray-500 dark:text-gray-400 text-sm">Tổng clicks</div>
                        <div className="text-2xl font-bold dark:text-white">{formatNumber(stats.totalClicks)}</div>
                    </div>
                    <div className="bg-white dark:bg-gray-800 rounded-xl p-4 shadow-sm">
                        <div className="text-gray-500 dark:text-gray-400 text-sm">Đơn hàng</div>
                        <div className="text-2xl font-bold dark:text-white">{formatNumber(stats.totalOrders)}</div>
                    </div>
                    <div className="bg-white dark:bg-gray-800 rounded-xl p-4 shadow-sm">
                        <div className="text-gray-500 dark:text-gray-400 text-sm">Tỷ lệ chuyển đổi</div>
                        <div className="text-2xl font-bold text-green-600">{stats.conversionRate}%</div>
                    </div>
                    <div className="bg-white dark:bg-gray-800 rounded-xl p-4 shadow-sm">
                        <div className="text-gray-500 dark:text-gray-400 text-sm">Đã nhận</div>
                        <div className="text-2xl font-bold text-[#ee4d2d]">₫{formatPrice(stats.totalEarnings)}</div>
                    </div>
                    <div className="bg-white dark:bg-gray-800 rounded-xl p-4 shadow-sm">
                        <div className="text-gray-500 dark:text-gray-400 text-sm">Đang chờ</div>
                        <div className="text-2xl font-bold text-yellow-600">₫{formatPrice(stats.pendingEarnings)}</div>
                    </div>
                </div>

                {/* Revenue Chart Placeholder */}
                <div className="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm mb-6">
                    <h3 className="font-bold mb-4 dark:text-white">📊 Thu nhập theo ngày</h3>
                    <div className="h-48 flex items-end justify-between gap-2">
                        {[65, 45, 78, 52, 90, 68, 85, 72, 95, 60, 82, 75].map((height, i) => (
                            <div key={i} className="flex-1 flex flex-col items-center">
                                <div
                                    className="w-full bg-gradient-to-t from-indigo-600 to-purple-500 rounded-t"
                                    style={{ height: `${height}%` }}
                                />
                                <span className="text-xs text-gray-400 mt-1">{i + 1}</span>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Top Products */}
                <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm overflow-hidden">
                    <div className="p-4 border-b dark:border-gray-700">
                        <h3 className="font-bold dark:text-white">🏆 Sản phẩm hot nhất</h3>
                    </div>

                    <div className="divide-y dark:divide-gray-700">
                        {products.map((product, index) => (
                            <div key={product.id} className="p-4">
                                <div className="flex items-center gap-4">
                                    <span className="text-2xl font-bold text-gray-300">#{index + 1}</span>
                                    <div className="w-16 h-16 bg-gray-100 dark:bg-gray-700 rounded-lg overflow-hidden relative">
                                        <Image src={product.image} alt="" fill className="object-cover" unoptimized />
                                    </div>
                                    <div className="flex-1">
                                        <h4 className="font-medium dark:text-white line-clamp-1">{product.name}</h4>
                                        <div className="text-sm text-gray-500 dark:text-gray-400">
                                            Hoa hồng: <span className="text-green-600 font-medium">{product.commission}%</span>
                                        </div>
                                    </div>
                                    <div className="text-right">
                                        <div className="font-bold text-[#ee4d2d]">₫{formatPrice(product.earnings)}</div>
                                        <div className="text-xs text-gray-500">{product.clicks} clicks • {product.orders} orders</div>
                                    </div>
                                </div>

                                {/* Affiliate Link */}
                                <div className="mt-3 flex gap-2">
                                    <input
                                        type="text"
                                        value={generateAffiliateLink(product.id)}
                                        readOnly
                                        className="flex-1 bg-gray-100 dark:bg-gray-700 px-3 py-2 rounded text-xs dark:text-gray-300"
                                    />
                                    <button
                                        onClick={() => copyLink(product.id)}
                                        className="px-4 py-2 bg-indigo-600 text-white text-sm rounded hover:opacity-90"
                                    >
                                        📋 Copy
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Find Products */}
                <div className="mt-6 bg-gradient-to-r from-indigo-600 to-purple-600 rounded-xl p-6 text-white">
                    <h3 className="font-bold text-xl mb-2">🔍 Tìm sản phẩm để chia sẻ</h3>
                    <p className="opacity-80 mb-4">Duyệt qua hàng triệu sản phẩm và tạo link affiliate</p>
                    <Link href="/products" className="inline-block px-6 py-3 bg-white text-indigo-600 rounded-full font-medium hover:opacity-90">
                        Khám phá ngay →
                    </Link>
                </div>

                {/* Tips */}
                <div className="mt-6 bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
                    <h3 className="font-bold mb-4 dark:text-white">💡 Mẹo tăng thu nhập</h3>
                    <ul className="space-y-2 text-sm text-gray-600 dark:text-gray-400">
                        <li>• Chia sẻ sản phẩm phù hợp với audience của bạn</li>
                        <li>• Tạo review chi tiết kèm link affiliate</li>
                        <li>• Tận dụng các chương trình Flash Sale để tăng tỷ lệ chuyển đổi</li>
                        <li>• Chia sẻ qua nhiều kênh: Facebook, TikTok, Blog, YouTube</li>
                        <li>• Theo dõi analytics để tối ưu chiến lược</li>
                    </ul>
                </div>
            </div>
        </div>
    );
}
