'use client';

import React, { useState } from 'react';
import Link from 'next/link';

export default function SellerDashboard() {
    const [activeTab, setActiveTab] = useState('overview');

    // Mock data
    const stats = {
        todayOrders: 24,
        todayRevenue: 15680000,
        pendingOrders: 8,
        totalProducts: 156,
        lowStock: 12,
        rating: 4.8,
        responseRate: 98,
        onTimeShipping: 95,
    };

    const recentOrders = [
        { id: 'ORD001', customer: 'Nguyễn Văn B', product: 'iPhone 15 Pro Max', amount: 29990000, status: 'pending', time: '10 phút trước' },
        { id: 'ORD002', customer: 'Trần Thị C', product: 'Son Dior x2', amount: 1900000, status: 'processing', time: '25 phút trước' },
        { id: 'ORD003', customer: 'Lê Văn D', product: 'Nike Air Force 1', amount: 2590000, status: 'shipped', time: '1 giờ trước' },
        { id: 'ORD004', customer: 'Phạm Thị E', product: 'MacBook Pro 14"', amount: 52990000, status: 'delivered', time: '2 giờ trước' },
    ];

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    const statusColors: Record<string, string> = {
        pending: 'bg-yellow-100 text-yellow-700',
        processing: 'bg-blue-100 text-blue-700',
        shipped: 'bg-purple-100 text-purple-700',
        delivered: 'bg-green-100 text-green-700',
    };

    const statusLabels: Record<string, string> = {
        pending: 'Chờ xử lý',
        processing: 'Đang xử lý',
        shipped: 'Đang giao',
        delivered: 'Đã giao',
    };

    return (
        <div className="min-h-screen bg-[#f5f5f5]">
            {/* Header */}
            <div className="bg-[#ee4d2d] text-white">
                <div className="container mx-auto px-4 py-4">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-4">
                            <Link href="/" className="text-2xl font-bold">Shopee</Link>
                            <span className="text-white/60">|</span>
                            <span className="text-lg">Kênh Người Bán</span>
                        </div>
                        <div className="flex items-center gap-4">
                            <button className="flex items-center gap-2 px-4 py-2 bg-white/10 rounded-sm hover:bg-white/20">
                                🔔 <span className="bg-white text-[#ee4d2d] text-xs px-1.5 rounded-full">3</span>
                            </button>
                            <div className="flex items-center gap-2">
                                <div className="w-8 h-8 bg-white/20 rounded-full flex items-center justify-center">👤</div>
                                <span>Apple Store Official</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div className="container mx-auto px-4 py-6">
                <div className="flex gap-6">
                    {/* Sidebar */}
                    <aside className="w-56 flex-shrink-0">
                        <nav className="bg-white rounded-sm shadow-sm overflow-hidden">
                            {[
                                { key: 'overview', icon: '📊', label: 'Tổng quan' },
                                { key: 'orders', icon: '📦', label: 'Quản lý đơn hàng' },
                                { key: 'products', icon: '🏷️', label: 'Quản lý sản phẩm' },
                                { key: 'finance', icon: '💰', label: 'Tài chính' },
                                { key: 'marketing', icon: '📣', label: 'Marketing' },
                                { key: 'analytics', icon: '📈', label: 'Phân tích' },
                                { key: 'settings', icon: '⚙️', label: 'Cài đặt Shop' },
                            ].map(item => (
                                <button
                                    key={item.key}
                                    onClick={() => setActiveTab(item.key)}
                                    className={`w-full px-4 py-3 flex items-center gap-3 text-sm transition-colors ${activeTab === item.key
                                            ? 'bg-[#fef6f5] text-[#ee4d2d] border-l-4 border-[#ee4d2d]'
                                            : 'hover:bg-gray-50'
                                        }`}
                                >
                                    <span>{item.icon}</span>
                                    <span>{item.label}</span>
                                </button>
                            ))}
                        </nav>
                    </aside>

                    {/* Main Content */}
                    <main className="flex-1 space-y-4">
                        {/* Stats Cards */}
                        <div className="grid grid-cols-4 gap-4">
                            <div className="bg-white rounded-sm shadow-sm p-4 animate-fade-in-up">
                                <div className="text-sm text-gray-500">Đơn hàng hôm nay</div>
                                <div className="text-2xl font-bold text-[#ee4d2d] mt-1">{stats.todayOrders}</div>
                                <div className="text-xs text-green-600 mt-1">↑ 12% so với hôm qua</div>
                            </div>
                            <div className="bg-white rounded-sm shadow-sm p-4 animate-fade-in-up" style={{ animationDelay: '50ms' }}>
                                <div className="text-sm text-gray-500">Doanh thu hôm nay</div>
                                <div className="text-2xl font-bold text-[#ee4d2d] mt-1">₫{formatPrice(stats.todayRevenue)}</div>
                                <div className="text-xs text-green-600 mt-1">↑ 8% so với hôm qua</div>
                            </div>
                            <div className="bg-white rounded-sm shadow-sm p-4 animate-fade-in-up" style={{ animationDelay: '100ms' }}>
                                <div className="text-sm text-gray-500">Chờ xử lý</div>
                                <div className="text-2xl font-bold text-orange-500 mt-1">{stats.pendingOrders}</div>
                                <div className="text-xs text-gray-500 mt-1">Cần xử lý ngay</div>
                            </div>
                            <div className="bg-white rounded-sm shadow-sm p-4 animate-fade-in-up" style={{ animationDelay: '150ms' }}>
                                <div className="text-sm text-gray-500">Sản phẩm</div>
                                <div className="text-2xl font-bold mt-1">{stats.totalProducts}</div>
                                <div className="text-xs text-orange-500 mt-1">{stats.lowStock} sản phẩm sắp hết hàng</div>
                            </div>
                        </div>

                        {/* Performance */}
                        <div className="bg-white rounded-sm shadow-sm p-4 animate-fade-in-up" style={{ animationDelay: '200ms' }}>
                            <h2 className="font-medium mb-4">Hiệu suất Shop</h2>
                            <div className="grid grid-cols-3 gap-6">
                                <div className="text-center">
                                    <div className="w-20 h-20 mx-auto rounded-full border-4 border-green-500 flex items-center justify-center">
                                        <span className="text-2xl font-bold text-green-500">{stats.rating}</span>
                                    </div>
                                    <p className="text-sm mt-2">Đánh giá Shop</p>
                                    <p className="text-xs text-gray-500">⭐⭐⭐⭐⭐</p>
                                </div>
                                <div className="text-center">
                                    <div className="w-20 h-20 mx-auto rounded-full border-4 border-blue-500 flex items-center justify-center">
                                        <span className="text-2xl font-bold text-blue-500">{stats.responseRate}%</span>
                                    </div>
                                    <p className="text-sm mt-2">Tỉ lệ phản hồi</p>
                                    <p className="text-xs text-gray-500">Trong vòng 12 giờ</p>
                                </div>
                                <div className="text-center">
                                    <div className="w-20 h-20 mx-auto rounded-full border-4 border-purple-500 flex items-center justify-center">
                                        <span className="text-2xl font-bold text-purple-500">{stats.onTimeShipping}%</span>
                                    </div>
                                    <p className="text-sm mt-2">Giao đúng hạn</p>
                                    <p className="text-xs text-gray-500">30 ngày qua</p>
                                </div>
                            </div>
                        </div>

                        {/* Recent Orders */}
                        <div className="bg-white rounded-sm shadow-sm animate-fade-in-up" style={{ animationDelay: '250ms' }}>
                            <div className="p-4 border-b flex items-center justify-between">
                                <h2 className="font-medium">Đơn hàng gần đây</h2>
                                <Link href="/seller/orders" className="text-[#ee4d2d] text-sm hover:underline">
                                    Xem tất cả →
                                </Link>
                            </div>
                            <div className="divide-y">
                                {recentOrders.map((order, i) => (
                                    <div key={order.id} className="p-4 flex items-center gap-4 hover:bg-gray-50 transition-colors">
                                        <div className="flex-1">
                                            <div className="flex items-center gap-2">
                                                <span className="font-medium text-sm">{order.id}</span>
                                                <span className={`text-xs px-2 py-0.5 rounded ${statusColors[order.status]}`}>
                                                    {statusLabels[order.status]}
                                                </span>
                                            </div>
                                            <p className="text-sm text-gray-500 mt-1">{order.customer} • {order.product}</p>
                                        </div>
                                        <div className="text-right">
                                            <p className="font-medium text-[#ee4d2d]">₫{formatPrice(order.amount)}</p>
                                            <p className="text-xs text-gray-400">{order.time}</p>
                                        </div>
                                        <button className="px-3 py-1 border text-sm hover:bg-gray-100 rounded-sm">
                                            Chi tiết
                                        </button>
                                    </div>
                                ))}
                            </div>
                        </div>

                        {/* Quick Actions */}
                        <div className="grid grid-cols-4 gap-4">
                            {[
                                { icon: '➕', label: 'Thêm sản phẩm', color: 'bg-blue-500' },
                                { icon: '🎁', label: 'Tạo voucher', color: 'bg-purple-500' },
                                { icon: '📢', label: 'Chạy quảng cáo', color: 'bg-orange-500' },
                                { icon: '💬', label: 'Trả lời chat', color: 'bg-green-500' },
                            ].map((action, i) => (
                                <button
                                    key={action.label}
                                    className="bg-white rounded-sm shadow-sm p-4 flex items-center gap-3 hover:shadow-md transition-shadow animate-fade-in-up"
                                    style={{ animationDelay: `${300 + i * 50}ms` }}
                                >
                                    <div className={`w-10 h-10 ${action.color} rounded-full flex items-center justify-center text-white text-xl`}>
                                        {action.icon}
                                    </div>
                                    <span className="font-medium">{action.label}</span>
                                </button>
                            ))}
                        </div>
                    </main>
                </div>
            </div>
        </div>
    );
}
