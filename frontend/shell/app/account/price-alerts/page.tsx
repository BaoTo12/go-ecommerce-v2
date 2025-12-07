'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';

interface PriceAlert {
    id: string;
    productId: string;
    productName: string;
    productImage: string;
    currentPrice: number;
    targetPrice: number;
    createdAt: string;
    triggered: boolean;
}

// Mock alerts
const MOCK_ALERTS: PriceAlert[] = [
    {
        id: 'alert1',
        productId: 'p1',
        productName: 'iPhone 15 Pro Max 256GB Titan Xanh',
        productImage: 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=200',
        currentPrice: 29990000,
        targetPrice: 27000000,
        createdAt: '2024-12-01',
        triggered: false,
    },
    {
        id: 'alert2',
        productId: 'p7',
        productName: 'Nike Air Force 1 Low White',
        productImage: 'https://images.unsplash.com/photo-1600269452121-4f2416e55c28?w=200',
        currentPrice: 2590000,
        targetPrice: 2000000,
        createdAt: '2024-12-03',
        triggered: true,
    },
];

export default function PriceAlertsPage() {
    const [alerts, setAlerts] = useState<PriceAlert[]>(MOCK_ALERTS);
    const [notification, setNotification] = useState<string | null>(null);

    const deleteAlert = (id: string) => {
        setAlerts(alerts.filter(a => a.id !== id));
        setNotification('✓ Đã xóa thông báo giá');
        setTimeout(() => setNotification(null), 2000);
    };

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    return (
        <div className="animate-fade-in">
            {notification && <div className="toast toast-success">{notification}</div>}

            <div className="bg-white dark:bg-gray-800 rounded-sm shadow-sm p-4 mb-4">
                <h1 className="text-lg font-medium dark:text-white">🔔 Thông Báo Giảm Giá</h1>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Nhận thông báo khi sản phẩm yêu thích giảm giá
                </p>
            </div>

            {alerts.length === 0 ? (
                <div className="bg-white dark:bg-gray-800 rounded-sm shadow-sm p-12 text-center">
                    <div className="text-5xl mb-4">🔔</div>
                    <p className="text-gray-500 dark:text-gray-400 mb-4">Bạn chưa thiết lập thông báo giảm giá nào</p>
                    <Link
                        href="/products"
                        className="inline-block px-6 py-2 bg-[#ee4d2d] text-white rounded-sm hover:opacity-90"
                    >
                        Khám phá sản phẩm
                    </Link>
                </div>
            ) : (
                <div className="space-y-3">
                    {alerts.map((alert, index) => (
                        <div
                            key={alert.id}
                            className={`bg-white dark:bg-gray-800 rounded-sm shadow-sm p-4 flex items-center gap-4 animate-fade-in-up ${alert.triggered ? 'border-l-4 border-green-500' : ''
                                }`}
                            style={{ animationDelay: `${index * 50}ms` }}
                        >
                            <div className="w-20 h-20 bg-gray-100 dark:bg-gray-700 rounded-sm overflow-hidden flex-shrink-0">
                                <img src={alert.productImage} alt="" className="w-full h-full object-cover" />
                            </div>

                            <div className="flex-1">
                                <Link
                                    href={`/products/${alert.productId}`}
                                    className="font-medium text-sm hover:text-[#ee4d2d] dark:text-white"
                                >
                                    {alert.productName}
                                </Link>

                                <div className="flex items-center gap-4 mt-2 text-sm">
                                    <span className="text-gray-500 dark:text-gray-400">
                                        Giá hiện tại: <span className="text-[#ee4d2d] font-medium">₫{formatPrice(alert.currentPrice)}</span>
                                    </span>
                                    <span className="text-gray-500 dark:text-gray-400">
                                        Giá mong muốn: <span className="font-medium dark:text-white">₫{formatPrice(alert.targetPrice)}</span>
                                    </span>
                                </div>

                                {alert.triggered && (
                                    <div className="mt-2">
                                        <span className="inline-flex items-center gap-1 text-xs text-green-600 bg-green-50 dark:bg-green-900 dark:text-green-300 px-2 py-1 rounded">
                                            ✓ Giá đã đạt mục tiêu!
                                        </span>
                                    </div>
                                )}
                            </div>

                            <div className="flex flex-col gap-2">
                                <Link
                                    href={`/products/${alert.productId}`}
                                    className="px-4 py-2 bg-[#ee4d2d] text-white text-sm rounded-sm hover:opacity-90"
                                >
                                    Xem ngay
                                </Link>
                                <button
                                    onClick={() => deleteAlert(alert.id)}
                                    className="px-4 py-2 border text-sm rounded-sm hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
                                >
                                    Xóa
                                </button>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* How it works */}
            <div className="bg-white dark:bg-gray-800 rounded-sm shadow-sm p-4 mt-6">
                <h3 className="font-medium mb-4 dark:text-white">💡 Cách hoạt động</h3>
                <div className="grid grid-cols-3 gap-4 text-center">
                    <div>
                        <div className="w-12 h-12 bg-[#fef6f5] dark:bg-gray-700 rounded-full flex items-center justify-center mx-auto text-xl">
                            1️⃣
                        </div>
                        <p className="text-sm mt-2 dark:text-gray-300">Chọn sản phẩm và thiết lập giá mong muốn</p>
                    </div>
                    <div>
                        <div className="w-12 h-12 bg-[#fef6f5] dark:bg-gray-700 rounded-full flex items-center justify-center mx-auto text-xl">
                            2️⃣
                        </div>
                        <p className="text-sm mt-2 dark:text-gray-300">Shopee theo dõi giá 24/7</p>
                    </div>
                    <div>
                        <div className="w-12 h-12 bg-[#fef6f5] dark:bg-gray-700 rounded-full flex items-center justify-center mx-auto text-xl">
                            3️⃣
                        </div>
                        <p className="text-sm mt-2 dark:text-gray-300">Nhận thông báo ngay khi giá giảm</p>
                    </div>
                </div>
            </div>
        </div>
    );
}
