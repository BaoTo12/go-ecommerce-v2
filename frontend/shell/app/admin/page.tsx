'use client';

import React, { useState } from 'react';
import Link from 'next/link';

export default function AdminDashboard() {
    const [activeTab, setActiveTab] = useState('overview');

    // Mock stats
    const stats = {
        totalUsers: 125430,
        newUsersToday: 342,
        totalOrders: 45678,
        ordersToday: 1234,
        totalRevenue: 15800000000,
        revenueToday: 156800000,
        totalProducts: 89456,
        activeShops: 12345,
        flaggedTransactions: 23,
        pendingReviews: 156,
    };

    const formatNumber = (num: number) => new Intl.NumberFormat('vi-VN').format(num);
    const formatCurrency = (num: number) => {
        if (num >= 1000000000) return (num / 1000000000).toFixed(1) + 'B';
        if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
        return formatNumber(num);
    };

    const recentActivities = [
        { type: 'order', message: 'Đơn hàng mới #SP241206012345 - ₫29,990,000', time: '2 phút trước' },
        { type: 'user', message: 'Người dùng mới đăng ký: user12345@email.com', time: '5 phút trước' },
        { type: 'fraud', message: '⚠️ Giao dịch đáng ngờ cần xem xét', time: '10 phút trước' },
        { type: 'shop', message: 'Shop mới đăng ký: Tech Store VN', time: '15 phút trước' },
        { type: 'review', message: 'Đánh giá mới cần duyệt (5 sao)', time: '20 phút trước' },
    ];

    return (
        <div className="min-h-screen bg-gray-100">
            {/* Header */}
            <div className="bg-[#1e293b] text-white">
                <div className="container mx-auto px-4 py-4">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-4">
                            <Link href="/" className="text-2xl font-bold text-[#ee4d2d]">Shopee</Link>
                            <span className="text-white/60">|</span>
                            <span className="text-lg">Admin Dashboard</span>
                        </div>
                        <div className="flex items-center gap-4">
                            <button className="relative p-2 hover:bg-white/10 rounded-full">
                                🔔
                                <span className="absolute top-0 right-0 w-2 h-2 bg-red-500 rounded-full" />
                            </button>
                            <div className="flex items-center gap-2">
                                <div className="w-8 h-8 bg-[#ee4d2d] rounded-full flex items-center justify-center">👤</div>
                                <span>Admin</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div className="container mx-auto px-4 py-6">
                <div className="flex gap-6">
                    {/* Sidebar */}
                    <aside className="w-56 flex-shrink-0">
                        <nav className="bg-white rounded-lg shadow-sm overflow-hidden">
                            {[
                                { key: 'overview', icon: '📊', label: 'Tổng quan' },
                                { key: 'users', icon: '👥', label: 'Quản lý người dùng' },
                                { key: 'shops', icon: '🏪', label: 'Quản lý Shop' },
                                { key: 'products', icon: '📦', label: 'Quản lý sản phẩm' },
                                { key: 'orders', icon: '🛒', label: 'Quản lý đơn hàng' },
                                { key: 'finance', icon: '💰', label: 'Tài chính' },
                                { key: 'fraud', icon: '🛡️', label: 'Phát hiện gian lận', badge: stats.flaggedTransactions },
                                { key: 'reviews', icon: '⭐', label: 'Đánh giá', badge: stats.pendingReviews },
                                { key: 'marketing', icon: '📣', label: 'Marketing' },
                                { key: 'settings', icon: '⚙️', label: 'Cài đặt hệ thống' },
                            ].map(item => (
                                <button
                                    key={item.key}
                                    onClick={() => setActiveTab(item.key)}
                                    className={`w-full px-4 py-3 flex items-center justify-between text-sm transition-colors ${activeTab === item.key
                                            ? 'bg-[#ee4d2d]/10 text-[#ee4d2d] border-l-4 border-[#ee4d2d]'
                                            : 'hover:bg-gray-50'
                                        }`}
                                >
                                    <span className="flex items-center gap-3">
                                        <span>{item.icon}</span>
                                        <span>{item.label}</span>
                                    </span>
                                    {item.badge && (
                                        <span className="bg-red-500 text-white text-xs px-2 py-0.5 rounded-full">
                                            {item.badge}
                                        </span>
                                    )}
                                </button>
                            ))}
                        </nav>
                    </aside>

                    {/* Main Content */}
                    <main className="flex-1 space-y-6">
                        {/* Stats Grid */}
                        <div className="grid grid-cols-4 gap-4">
                            {[
                                { label: 'Tổng người dùng', value: formatNumber(stats.totalUsers), sub: `+${stats.newUsersToday} hôm nay`, color: 'text-blue-600', bg: 'bg-blue-50' },
                                { label: 'Tổng đơn hàng', value: formatNumber(stats.totalOrders), sub: `+${stats.ordersToday} hôm nay`, color: 'text-green-600', bg: 'bg-green-50' },
                                { label: 'Tổng doanh thu', value: `₫${formatCurrency(stats.totalRevenue)}`, sub: `+₫${formatCurrency(stats.revenueToday)} hôm nay`, color: 'text-purple-600', bg: 'bg-purple-50' },
                                { label: 'Shop đang hoạt động', value: formatNumber(stats.activeShops), sub: `${formatNumber(stats.totalProducts)} sản phẩm`, color: 'text-orange-600', bg: 'bg-orange-50' },
                            ].map((stat, i) => (
                                <div
                                    key={stat.label}
                                    className={`${stat.bg} rounded-lg p-5 animate-fade-in-up`}
                                    style={{ animationDelay: `${i * 50}ms` }}
                                >
                                    <div className="text-sm text-gray-600">{stat.label}</div>
                                    <div className={`text-2xl font-bold ${stat.color} mt-1`}>{stat.value}</div>
                                    <div className="text-xs text-gray-500 mt-1">{stat.sub}</div>
                                </div>
                            ))}
                        </div>

                        {/* Charts Row */}
                        <div className="grid grid-cols-2 gap-6">
                            {/* Revenue Chart */}
                            <div className="bg-white rounded-lg shadow-sm p-5 animate-fade-in-up" style={{ animationDelay: '200ms' }}>
                                <h3 className="font-medium mb-4">Doanh thu 7 ngày qua</h3>
                                <div className="h-48 flex items-end gap-2">
                                    {[65, 78, 52, 89, 95, 72, 88].map((height, i) => (
                                        <div key={i} className="flex-1 flex flex-col items-center">
                                            <div
                                                className="w-full bg-gradient-to-t from-[#ee4d2d] to-[#ff8f70] rounded-t transition-all hover:opacity-80"
                                                style={{ height: `${height}%` }}
                                            />
                                            <span className="text-xs text-gray-400 mt-2">
                                                {['T2', 'T3', 'T4', 'T5', 'T6', 'T7', 'CN'][i]}
                                            </span>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            {/* Orders Chart */}
                            <div className="bg-white rounded-lg shadow-sm p-5 animate-fade-in-up" style={{ animationDelay: '250ms' }}>
                                <h3 className="font-medium mb-4">Đơn hàng theo trạng thái</h3>
                                <div className="space-y-3">
                                    {[
                                        { label: 'Hoàn thành', value: 75, color: 'bg-green-500' },
                                        { label: 'Đang giao', value: 15, color: 'bg-blue-500' },
                                        { label: 'Đang xử lý', value: 7, color: 'bg-yellow-500' },
                                        { label: 'Đã hủy', value: 3, color: 'bg-red-500' },
                                    ].map(item => (
                                        <div key={item.label} className="flex items-center gap-3">
                                            <span className="w-24 text-sm text-gray-600">{item.label}</span>
                                            <div className="flex-1 h-6 bg-gray-100 rounded-full overflow-hidden">
                                                <div
                                                    className={`h-full ${item.color} transition-all`}
                                                    style={{ width: `${item.value}%` }}
                                                />
                                            </div>
                                            <span className="w-12 text-sm text-gray-600 text-right">{item.value}%</span>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        </div>

                        {/* Recent Activity and Alerts */}
                        <div className="grid grid-cols-2 gap-6">
                            {/* Recent Activity */}
                            <div className="bg-white rounded-lg shadow-sm animate-fade-in-up" style={{ animationDelay: '300ms' }}>
                                <div className="p-4 border-b">
                                    <h3 className="font-medium">Hoạt động gần đây</h3>
                                </div>
                                <div className="divide-y">
                                    {recentActivities.map((activity, i) => (
                                        <div key={i} className="p-4 flex items-center gap-3 hover:bg-gray-50">
                                            <span className={`w-8 h-8 rounded-full flex items-center justify-center text-sm ${activity.type === 'order' ? 'bg-green-100' :
                                                    activity.type === 'user' ? 'bg-blue-100' :
                                                        activity.type === 'fraud' ? 'bg-red-100' :
                                                            activity.type === 'shop' ? 'bg-purple-100' :
                                                                'bg-yellow-100'
                                                }`}>
                                                {activity.type === 'order' ? '🛒' :
                                                    activity.type === 'user' ? '👤' :
                                                        activity.type === 'fraud' ? '⚠️' :
                                                            activity.type === 'shop' ? '🏪' :
                                                                '⭐'}
                                            </span>
                                            <div className="flex-1">
                                                <p className="text-sm">{activity.message}</p>
                                                <p className="text-xs text-gray-400">{activity.time}</p>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            {/* Alerts */}
                            <div className="bg-white rounded-lg shadow-sm animate-fade-in-up" style={{ animationDelay: '350ms' }}>
                                <div className="p-4 border-b">
                                    <h3 className="font-medium">Cảnh báo hệ thống</h3>
                                </div>
                                <div className="p-4 space-y-3">
                                    <div className="p-3 bg-red-50 border border-red-200 rounded-lg">
                                        <div className="flex items-center gap-2 text-red-700 font-medium">
                                            <span>🚨</span> Phát hiện giao dịch gian lận
                                        </div>
                                        <p className="text-sm text-red-600 mt-1">{stats.flaggedTransactions} giao dịch cần xem xét</p>
                                        <button className="mt-2 text-sm text-red-700 hover:underline">Xem chi tiết →</button>
                                    </div>
                                    <div className="p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
                                        <div className="flex items-center gap-2 text-yellow-700 font-medium">
                                            <span>⚠️</span> Đánh giá chờ duyệt
                                        </div>
                                        <p className="text-sm text-yellow-600 mt-1">{stats.pendingReviews} đánh giá cần xem xét</p>
                                        <button className="mt-2 text-sm text-yellow-700 hover:underline">Xem chi tiết →</button>
                                    </div>
                                    <div className="p-3 bg-blue-50 border border-blue-200 rounded-lg">
                                        <div className="flex items-center gap-2 text-blue-700 font-medium">
                                            <span>💡</span> Gợi ý tối ưu
                                        </div>
                                        <p className="text-sm text-blue-600 mt-1">Tỷ lệ chuyển đổi có thể tăng 15% với A/B testing</p>
                                        <button className="mt-2 text-sm text-blue-700 hover:underline">Tìm hiểu thêm →</button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </main>
                </div>
            </div>
        </div>
    );
}
