'use client';

import React, { useState, useEffect, Suspense } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useSearchParams } from 'next/navigation';
import { userService, Order, OrderStatus } from '@/services/userService';

const STATUS_LABELS: Record<OrderStatus, { label: string; color: string; icon: string }> = {
    pending_payment: { label: 'Chờ thanh toán', color: 'text-yellow-600 bg-yellow-50', icon: '💳' },
    processing: { label: 'Đang xử lý', color: 'text-blue-600 bg-blue-50', icon: '⏳' },
    shipped: { label: 'Đang vận chuyển', color: 'text-purple-600 bg-purple-50', icon: '📦' },
    out_for_delivery: { label: 'Đang giao hàng', color: 'text-orange-600 bg-orange-50', icon: '🚚' },
    delivered: { label: 'Đã giao', color: 'text-green-600 bg-green-50', icon: '✅' },
    cancelled: { label: 'Đã hủy', color: 'text-gray-600 bg-gray-50', icon: '❌' },
    refunded: { label: 'Hoàn tiền', color: 'text-red-600 bg-red-50', icon: '↩️' },
};

function OrdersContent() {
    const searchParams = useSearchParams();
    const statusFilter = searchParams.get('status') as OrderStatus | null;

    const [orders, setOrders] = useState<Order[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [activeTab, setActiveTab] = useState<OrderStatus | 'all'>('all');

    const tabs = [
        { key: 'all', label: 'Tất cả' },
        { key: 'pending_payment', label: 'Chờ thanh toán' },
        { key: 'processing', label: 'Đang xử lý' },
        { key: 'shipped', label: 'Đang vận chuyển' },
        { key: 'out_for_delivery', label: 'Đang giao' },
        { key: 'delivered', label: 'Đã giao' },
        { key: 'cancelled', label: 'Đã hủy' },
    ];

    useEffect(() => {
        const loadOrders = async () => {
            setIsLoading(true);
            try {
                const data = await userService.getOrders(statusFilter || undefined);
                setOrders(data);
                setActiveTab(statusFilter || 'all');
            } catch (error) {
                console.error('Failed to load orders:', error);
            } finally {
                setIsLoading(false);
            }
        };
        loadOrders();
    }, [statusFilter]);

    const filteredOrders = activeTab === 'all'
        ? orders
        : orders.filter(o => o.status === activeTab);

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);
    const formatDate = (date: string) => new Date(date).toLocaleDateString('vi-VN', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    });

    return (
        <div className="animate-fade-in">
            {/* Tabs */}
            <div className="bg-white rounded-sm shadow-sm mb-4 overflow-x-auto">
                <div className="flex min-w-max">
                    {tabs.map(tab => (
                        <button
                            key={tab.key}
                            onClick={() => setActiveTab(tab.key as OrderStatus | 'all')}
                            className={`px-6 py-4 text-sm font-medium border-b-2 transition-all whitespace-nowrap ${activeTab === tab.key
                                    ? 'text-[#ee4d2d] border-[#ee4d2d]'
                                    : 'text-gray-500 border-transparent hover:text-[#ee4d2d]'
                                }`}
                        >
                            {tab.label}
                        </button>
                    ))}
                </div>
            </div>

            {/* Search */}
            <div className="bg-white rounded-sm shadow-sm p-4 mb-4">
                <div className="flex gap-2">
                    <input
                        type="text"
                        placeholder="Tìm kiếm theo Tên Shop, ID đơn hàng hoặc Tên Sản phẩm"
                        className="flex-1 border px-3 py-2 text-sm outline-none focus:border-[#ee4d2d] rounded-sm"
                    />
                    <button className="px-4 py-2 bg-[#ee4d2d] text-white text-sm hover:opacity-90 transition-all">
                        🔍 Tìm
                    </button>
                </div>
            </div>

            {/* Orders List */}
            {isLoading ? (
                <div className="space-y-4">
                    {[1, 2, 3].map(i => (
                        <div key={i} className="bg-white rounded-sm shadow-sm p-4 animate-pulse">
                            <div className="h-4 bg-gray-200 rounded w-1/4 mb-4" />
                            <div className="flex gap-4">
                                <div className="w-20 h-20 bg-gray-200 rounded" />
                                <div className="flex-1 space-y-2">
                                    <div className="h-4 bg-gray-200 rounded w-3/4" />
                                    <div className="h-4 bg-gray-200 rounded w-1/2" />
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            ) : filteredOrders.length === 0 ? (
                <div className="bg-white rounded-sm shadow-sm p-12 text-center">
                    <div className="text-5xl mb-4">📦</div>
                    <p className="text-gray-500">Chưa có đơn hàng nào</p>
                    <Link href="/products" className="inline-block mt-4 px-6 py-2 bg-[#ee4d2d] text-white text-sm hover:opacity-90">
                        Mua sắm ngay
                    </Link>
                </div>
            ) : (
                <div className="space-y-4">
                    {filteredOrders.map((order, index) => {
                        const statusInfo = STATUS_LABELS[order.status];
                        return (
                            <div
                                key={order.id}
                                className="bg-white rounded-sm shadow-sm overflow-hidden animate-fade-in-up"
                                style={{ animationDelay: `${index * 50}ms` }}
                            >
                                {/* Order Header */}
                                <div className="p-4 border-b flex items-center justify-between bg-gray-50">
                                    <div className="flex items-center gap-3">
                                        <span className="text-sm font-medium">{order.items[0].shopName}</span>
                                        <button className="text-xs text-[#ee4d2d] border border-[#ee4d2d] px-2 py-0.5 hover:bg-[#fef6f5]">
                                            💬 Chat
                                        </button>
                                        <button className="text-xs text-gray-500 border px-2 py-0.5 hover:bg-gray-100">
                                            🏪 Xem Shop
                                        </button>
                                    </div>
                                    <div className="flex items-center gap-2">
                                        {order.trackingNumber && (
                                            <Link
                                                href={`/account/orders/${order.id}/tracking`}
                                                className="text-xs text-[#ee4d2d] hover:underline"
                                            >
                                                🚚 {order.trackingNumber}
                                            </Link>
                                        )}
                                        <span className={`text-xs px-2 py-1 rounded ${statusInfo.color}`}>
                                            {statusInfo.icon} {statusInfo.label}
                                        </span>
                                    </div>
                                </div>

                                {/* Order Items */}
                                {order.items.map(item => (
                                    <Link
                                        key={item.id}
                                        href={`/products/${item.productId}`}
                                        className="p-4 border-b flex gap-4 hover:bg-gray-50 transition-colors"
                                    >
                                        <div className="relative w-20 h-20 bg-gray-100 rounded-sm overflow-hidden flex-shrink-0">
                                            <Image
                                                src={item.image}
                                                alt={item.name}
                                                fill
                                                className="object-cover"
                                                unoptimized
                                            />
                                        </div>
                                        <div className="flex-1">
                                            <h3 className="text-sm line-clamp-2">{item.name}</h3>
                                            {item.variant && (
                                                <p className="text-xs text-gray-500 mt-1">Phân loại: {item.variant}</p>
                                            )}
                                            <p className="text-xs text-gray-500">x{item.quantity}</p>
                                        </div>
                                        <div className="text-right">
                                            <p className="text-sm text-[#ee4d2d]">₫{formatPrice(item.price)}</p>
                                        </div>
                                    </Link>
                                ))}

                                {/* Order Footer */}
                                <div className="p-4 flex items-center justify-between bg-[#fffefb]">
                                    <div className="text-xs text-gray-500">
                                        Đặt lúc: {formatDate(order.createdAt)}
                                    </div>
                                    <div className="flex items-center gap-4">
                                        <div className="text-right">
                                            <span className="text-xs text-gray-500">Thành tiền: </span>
                                            <span className="text-lg text-[#ee4d2d] font-medium">₫{formatPrice(order.total)}</span>
                                        </div>
                                        <div className="flex gap-2">
                                            {order.status === 'delivered' && (
                                                <button className="px-4 py-2 bg-[#ee4d2d] text-white text-sm hover:opacity-90 transition-all">
                                                    Mua Lại
                                                </button>
                                            )}
                                            {['shipped', 'out_for_delivery'].includes(order.status) && (
                                                <Link
                                                    href={`/account/orders/${order.id}/tracking`}
                                                    className="px-4 py-2 border border-[#ee4d2d] text-[#ee4d2d] text-sm hover:bg-[#fef6f5] transition-all"
                                                >
                                                    Theo Dõi Đơn Hàng
                                                </Link>
                                            )}
                                            {['pending_payment', 'processing'].includes(order.status) && (
                                                <button className="px-4 py-2 border text-sm text-gray-600 hover:bg-gray-50 transition-all">
                                                    Hủy Đơn Hàng
                                                </button>
                                            )}
                                            <Link
                                                href={`/account/orders/${order.id}`}
                                                className="px-4 py-2 border text-sm hover:bg-gray-50 transition-all"
                                            >
                                                Xem Chi Tiết
                                            </Link>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        );
                    })}
                </div>
            )}
        </div>
    );
}

export default function OrdersPage() {
    return (
        <Suspense fallback={<div className="animate-pulse">Loading...</div>}>
            <OrdersContent />
        </Suspense>
    );
}
