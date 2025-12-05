'use client';

import React, { useState, useEffect } from 'react';

interface Metrics {
    activeUsers: number;
    ordersInProgress: number;
    todayRevenue: number;
    todayOrders: number;
    revenueChange: number;
    ordersChange: number;
}

export default function AnalyticsDashboard() {
    const [metrics, setMetrics] = useState<Metrics | null>(null);
    const [period, setPeriod] = useState('today');

    useEffect(() => {
        setMetrics({
            activeUsers: 12453,
            ordersInProgress: 847,
            todayRevenue: 5678900000,
            todayOrders: 34521,
            revenueChange: 12.5,
            ordersChange: 8.3,
        });
    }, []);

    const funnelData = [
        { name: 'Truy cập', users: 100000, rate: 100 },
        { name: 'Xem sản phẩm', users: 45000, rate: 45 },
        { name: 'Thêm giỏ hàng', users: 18000, rate: 40 },
        { name: 'Thanh toán', users: 9000, rate: 50 },
        { name: 'Hoàn thành', users: 7200, rate: 80 },
    ];

    const topProducts = [
        { name: 'iPhone 15 Pro Max', sales: 1234, revenue: 37020000000 },
        { name: 'Samsung S24 Ultra', sales: 987, revenue: 19740000000 },
        { name: 'MacBook Air M3', sales: 654, revenue: 16350000000 },
        { name: 'AirPods Pro 2', sales: 2345, revenue: 5862500000 },
    ];

    if (!metrics) {
        return <div className="flex h-96 items-center justify-center">Đang tải...</div>;
    }

    return (
        <div className="min-h-screen bg-[#F5F5F5]">
            {/* Header */}
            <div className="bg-white border-b px-6 py-4">
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-2xl font-bold text-gray-800">📊 Analytics Dashboard</h1>
                        <p className="text-gray-500 text-sm">Số liệu thời gian thực</p>
                    </div>
                    <div className="flex gap-2">
                        {['today', 'week', 'month'].map(p => (
                            <button
                                key={p}
                                onClick={() => setPeriod(p)}
                                className={`px-4 py-2 rounded text-sm font-medium ${period === p
                                        ? 'bg-[#EE4D2D] text-white'
                                        : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                                    }`}
                            >
                                {p === 'today' ? 'Hôm nay' : p === 'week' ? 'Tuần' : 'Tháng'}
                            </button>
                        ))}
                    </div>
                </div>
            </div>

            <div className="p-6">
                {/* Real-time Stats */}
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
                    <StatCard
                        icon="👥"
                        label="Người dùng online"
                        value={metrics.activeUsers.toLocaleString()}
                        live
                    />
                    <StatCard
                        icon="🛒"
                        label="Đơn đang xử lý"
                        value={metrics.ordersInProgress.toLocaleString()}
                        live
                    />
                    <StatCard
                        icon="💰"
                        label="Doanh thu hôm nay"
                        value={`${(metrics.todayRevenue / 1000000000).toFixed(1)}B₫`}
                        change={metrics.revenueChange}
                    />
                    <StatCard
                        icon="📦"
                        label="Đơn hàng hôm nay"
                        value={metrics.todayOrders.toLocaleString()}
                        change={metrics.ordersChange}
                    />
                </div>

                <div className="grid lg:grid-cols-2 gap-6">
                    {/* Conversion Funnel */}
                    <div className="bg-white rounded p-6">
                        <h2 className="font-bold text-lg mb-4">🔄 Phễu chuyển đổi</h2>
                        <div className="space-y-4">
                            {funnelData.map((step, i) => {
                                const width = (step.users / funnelData[0].users) * 100;
                                return (
                                    <div key={step.name}>
                                        <div className="flex justify-between text-sm mb-1">
                                            <span>{step.name}</span>
                                            <span className="text-gray-500">{step.users.toLocaleString()} ({step.rate}%)</span>
                                        </div>
                                        <div className="h-8 bg-gray-100 rounded overflow-hidden">
                                            <div
                                                className="h-full bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] rounded flex items-center justify-end pr-2"
                                                style={{ width: `${width}%` }}
                                            >
                                                {width > 25 && <span className="text-white text-xs font-bold">{step.rate}%</span>}
                                            </div>
                                        </div>
                                    </div>
                                );
                            })}
                            <div className="mt-4 p-4 bg-green-50 rounded text-center">
                                <span className="text-sm text-green-600">Tỷ lệ chuyển đổi tổng:</span>
                                <div className="text-3xl font-bold text-green-600">7.2%</div>
                            </div>
                        </div>
                    </div>

                    {/* Top Products */}
                    <div className="bg-white rounded p-6">
                        <h2 className="font-bold text-lg mb-4">🏆 Sản phẩm bán chạy</h2>
                        <div className="space-y-3">
                            {topProducts.map((product, i) => (
                                <div key={product.name} className="flex items-center gap-3 p-3 bg-gray-50 rounded">
                                    <span className="text-2xl font-bold text-[#EE4D2D]">#{i + 1}</span>
                                    <div className="flex-1">
                                        <div className="font-semibold text-sm">{product.name}</div>
                                        <div className="text-xs text-gray-500">{product.sales.toLocaleString()} đã bán</div>
                                    </div>
                                    <div className="text-right">
                                        <div className="font-bold text-[#EE4D2D]">
                                            {(product.revenue / 1000000000).toFixed(1)}B₫
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>

                {/* Revenue Chart Placeholder */}
                <div className="mt-6 bg-white rounded p-6">
                    <h2 className="font-bold text-lg mb-4">📈 Biểu đồ doanh thu</h2>
                    <div className="h-64 bg-gradient-to-t from-[#EE4D2D]/10 to-transparent rounded flex items-end justify-around p-4">
                        {[65, 82, 45, 91, 78, 95, 88, 72, 85, 90, 75, 88].map((h, i) => (
                            <div key={i} className="flex flex-col items-center gap-1">
                                <div
                                    className="w-8 bg-gradient-to-t from-[#EE4D2D] to-[#FF6633] rounded-t hover:opacity-80 transition-opacity"
                                    style={{ height: `${h}%` }}
                                />
                                <span className="text-xs text-gray-400">{i + 1}</span>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}

function StatCard({ icon, label, value, change, live }: {
    icon: string;
    label: string;
    value: string;
    change?: number;
    live?: boolean;
}) {
    return (
        <div className="bg-white rounded p-4">
            <div className="flex items-start justify-between">
                <span className="text-2xl">{icon}</span>
                {live && <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />}
            </div>
            <div className="mt-2">
                <div className="text-2xl font-bold">{value}</div>
                <div className="text-sm text-gray-500">{label}</div>
                {change !== undefined && (
                    <div className={`text-sm font-semibold ${change >= 0 ? 'text-green-500' : 'text-red-500'}`}>
                        {change >= 0 ? '↑' : '↓'} {Math.abs(change)}%
                    </div>
                )}
            </div>
        </div>
    );
}
