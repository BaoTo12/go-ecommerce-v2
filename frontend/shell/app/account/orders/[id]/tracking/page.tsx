'use client';

import React, { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import Image from 'next/image';
import { userService, Order, ShippingUpdate } from '@/services/userService';

export default function TrackingPage() {
    const params = useParams();
    const orderId = params.id as string;

    const [order, setOrder] = useState<Order | null>(null);
    const [updates, setUpdates] = useState<ShippingUpdate[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const loadData = async () => {
            setIsLoading(true);
            try {
                const [orderData, updatesData] = await Promise.all([
                    userService.getOrder(orderId),
                    userService.getShippingUpdates(orderId),
                ]);
                setOrder(orderData);
                setUpdates(updatesData);
            } catch (error) {
                console.error('Failed to load tracking data:', error);
            } finally {
                setIsLoading(false);
            }
        };
        loadData();
    }, [orderId]);

    const formatDate = (date: string) => new Date(date).toLocaleDateString('vi-VN', {
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    });

    if (isLoading) {
        return (
            <div className="bg-white rounded-sm shadow-sm p-6 animate-pulse">
                <div className="h-6 bg-gray-200 rounded w-1/3 mb-4" />
                <div className="space-y-6">
                    {[1, 2, 3, 4].map(i => (
                        <div key={i} className="flex gap-4">
                            <div className="w-4 h-4 bg-gray-200 rounded-full" />
                            <div className="flex-1 space-y-2">
                                <div className="h-4 bg-gray-200 rounded w-1/4" />
                                <div className="h-4 bg-gray-200 rounded w-3/4" />
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        );
    }

    if (!order) {
        return (
            <div className="bg-white rounded-sm shadow-sm p-12 text-center">
                <div className="text-5xl mb-4">😕</div>
                <p className="text-gray-500">Không tìm thấy đơn hàng</p>
                <Link href="/account/orders" className="inline-block mt-4 text-[#ee4d2d] hover:underline">
                    ← Quay lại danh sách đơn hàng
                </Link>
            </div>
        );
    }

    return (
        <div className="animate-fade-in space-y-4">
            {/* Header */}
            <div className="bg-white rounded-sm shadow-sm p-4 flex items-center justify-between">
                <Link href={`/account/orders/${orderId}`} className="text-sm text-gray-500 hover:text-[#ee4d2d]">
                    ← TRỞ LẠI ĐƠN HÀNG
                </Link>
                <span className="text-sm text-gray-500">
                    Mã vận đơn: <span className="text-[#ee4d2d] font-medium">{order.trackingNumber}</span>
                </span>
            </div>

            {/* Shipping Info */}
            <div className="bg-gradient-to-r from-[#00bfa5] to-[#00897b] rounded-sm p-6 text-white">
                <div className="flex items-start justify-between">
                    <div>
                        <p className="text-sm opacity-80">Đơn vị vận chuyển</p>
                        <p className="text-xl font-medium mt-1">{order.shippingMethod}</p>
                        <p className="text-sm opacity-80 mt-2">Mã vận đơn: {order.trackingNumber}</p>
                    </div>
                    <div className="text-right">
                        <p className="text-sm opacity-80">Dự kiến giao</p>
                        {order.estimatedDelivery && (
                            <p className="text-lg font-medium mt-1">
                                {new Date(order.estimatedDelivery).toLocaleDateString('vi-VN', {
                                    weekday: 'long',
                                    month: 'long',
                                    day: 'numeric'
                                })}
                            </p>
                        )}
                    </div>
                </div>

                {/* Map Preview */}
                <div className="mt-4 bg-white/10 rounded-sm p-4 flex items-center gap-4">
                    <div className="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center text-2xl animate-pulse">
                        🚚
                    </div>
                    <div>
                        <p className="font-medium">Đang trên đường giao đến bạn</p>
                        <p className="text-sm opacity-80">Shipper đang di chuyển đến địa chỉ của bạn</p>
                    </div>
                </div>
            </div>

            {/* Delivery Address */}
            <div className="bg-white rounded-sm shadow-sm p-4">
                <h3 className="text-sm font-medium mb-3 flex items-center gap-2 text-[#ee4d2d]">
                    <span>📍</span> Địa Chỉ Giao Hàng
                </h3>
                <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-[#fef6f5] rounded-full flex items-center justify-center text-lg">
                        🏠
                    </div>
                    <div>
                        <p className="font-medium text-sm">
                            {order.shippingAddress.name} | {order.shippingAddress.phone}
                        </p>
                        <p className="text-sm text-gray-500">
                            {order.shippingAddress.address}, {order.shippingAddress.district}, {order.shippingAddress.city}
                        </p>
                    </div>
                </div>
            </div>

            {/* Tracking Timeline */}
            <div className="bg-white rounded-sm shadow-sm p-6">
                <h3 className="text-sm font-medium mb-6">Lịch Sử Vận Chuyển</h3>

                {updates.length === 0 ? (
                    <div className="text-center py-8">
                        <div className="text-4xl mb-2">📦</div>
                        <p className="text-gray-500 text-sm">Chưa có thông tin vận chuyển</p>
                    </div>
                ) : (
                    <div className="relative">
                        {/* Timeline line */}
                        <div className="absolute left-[7px] top-2 bottom-2 w-0.5 bg-gray-200" />

                        {updates.map((update, index) => (
                            <div
                                key={update.id}
                                className="relative flex gap-4 pb-6 last:pb-0 animate-fade-in-left"
                                style={{ animationDelay: `${index * 100}ms` }}
                            >
                                {/* Dot */}
                                <div className={`relative z-10 w-4 h-4 rounded-full flex-shrink-0 mt-0.5 ${index === 0
                                        ? 'bg-[#00bfa5] ring-4 ring-[#00bfa5]/20 animate-pulse'
                                        : 'bg-gray-300'
                                    }`} />

                                {/* Content */}
                                <div className="flex-1">
                                    <div className="flex items-center gap-2 mb-1">
                                        <span className={`font-medium text-sm ${index === 0 ? 'text-[#00bfa5]' : 'text-gray-700'}`}>
                                            {update.status}
                                        </span>
                                        <span className="text-xs text-gray-400">{formatDate(update.timestamp)}</span>
                                    </div>
                                    <p className="text-sm text-gray-600">{update.description}</p>
                                    {update.location && (
                                        <p className="text-xs text-gray-400 mt-1 flex items-center gap-1">
                                            <span>📍</span> {update.location}
                                        </p>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            {/* Product Preview */}
            <div className="bg-white rounded-sm shadow-sm p-4">
                <h3 className="text-sm font-medium mb-3">Sản Phẩm Trong Đơn Hàng</h3>
                <div className="flex gap-2 overflow-x-auto pb-2">
                    {order.items.map(item => (
                        <Link
                            key={item.id}
                            href={`/products/${item.productId}`}
                            className="relative w-16 h-16 bg-gray-100 rounded-sm overflow-hidden flex-shrink-0 hover:scale-105 transition-transform"
                        >
                            <Image
                                src={item.image}
                                alt={item.name}
                                fill
                                className="object-cover"
                                unoptimized
                            />
                            {item.quantity > 1 && (
                                <span className="absolute bottom-0 right-0 bg-black/60 text-white text-xs px-1">
                                    x{item.quantity}
                                </span>
                            )}
                        </Link>
                    ))}
                </div>
            </div>

            {/* Contact Options */}
            <div className="bg-white rounded-sm shadow-sm p-4">
                <h3 className="text-sm font-medium mb-3">Cần Hỗ Trợ?</h3>
                <div className="grid grid-cols-3 gap-2">
                    <button className="p-3 border rounded-sm text-center hover:bg-gray-50 transition-all">
                        <span className="text-2xl block mb-1">📞</span>
                        <span className="text-xs text-gray-600">Gọi Shipper</span>
                    </button>
                    <button className="p-3 border rounded-sm text-center hover:bg-gray-50 transition-all">
                        <span className="text-2xl block mb-1">💬</span>
                        <span className="text-xs text-gray-600">Chat Người Bán</span>
                    </button>
                    <button className="p-3 border rounded-sm text-center hover:bg-gray-50 transition-all">
                        <span className="text-2xl block mb-1">🎧</span>
                        <span className="text-xs text-gray-600">Hỗ Trợ Shopee</span>
                    </button>
                </div>
            </div>
        </div>
    );
}
