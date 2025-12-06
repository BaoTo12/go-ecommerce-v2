'use client';

import React, { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import Image from 'next/image';
import { userService, Order, OrderStatus } from '@/services/userService';

const STATUS_LABELS: Record<OrderStatus, { label: string; color: string; icon: string }> = {
    pending_payment: { label: 'Chờ thanh toán', color: 'text-yellow-600', icon: '💳' },
    processing: { label: 'Đang xử lý', color: 'text-blue-600', icon: '⏳' },
    shipped: { label: 'Đang vận chuyển', color: 'text-purple-600', icon: '📦' },
    out_for_delivery: { label: 'Đang giao hàng', color: 'text-orange-600', icon: '🚚' },
    delivered: { label: 'Đã giao', color: 'text-green-600', icon: '✅' },
    cancelled: { label: 'Đã hủy', color: 'text-gray-600', icon: '❌' },
    refunded: { label: 'Hoàn tiền', color: 'text-red-600', icon: '↩️' },
};

export default function OrderDetailPage() {
    const params = useParams();
    const orderId = params.id as string;

    const [order, setOrder] = useState<Order | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const loadOrder = async () => {
            setIsLoading(true);
            try {
                const data = await userService.getOrder(orderId);
                setOrder(data);
            } catch (error) {
                console.error('Failed to load order:', error);
            } finally {
                setIsLoading(false);
            }
        };
        loadOrder();
    }, [orderId]);

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);
    const formatDate = (date: string) => new Date(date).toLocaleDateString('vi-VN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    });

    if (isLoading) {
        return (
            <div className="bg-white rounded-sm shadow-sm p-6 animate-pulse">
                <div className="h-6 bg-gray-200 rounded w-1/3 mb-4" />
                <div className="h-4 bg-gray-200 rounded w-1/2 mb-8" />
                <div className="space-y-4">
                    {[1, 2, 3].map(i => (
                        <div key={i} className="h-20 bg-gray-200 rounded" />
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

    const statusInfo = STATUS_LABELS[order.status];

    return (
        <div className="animate-fade-in space-y-4">
            {/* Header */}
            <div className="bg-white rounded-sm shadow-sm p-4 flex items-center justify-between">
                <Link href="/account/orders" className="text-sm text-gray-500 hover:text-[#ee4d2d]">
                    ← TRỞ LẠI
                </Link>
                <div className="flex items-center gap-4">
                    <span className="text-sm text-gray-500">MÃ ĐƠN HÀNG: {order.orderNumber}</span>
                    <span className="text-sm">|</span>
                    <span className={`text-sm font-medium ${statusInfo.color}`}>
                        {statusInfo.icon} {statusInfo.label.toUpperCase()}
                    </span>
                </div>
            </div>

            {/* Shipping Progress */}
            {['shipped', 'out_for_delivery', 'delivered'].includes(order.status) && (
                <div className="bg-gradient-to-r from-[#00bfa5] to-[#00897b] rounded-sm p-6 text-white">
                    <div className="flex items-center justify-between mb-4">
                        <div>
                            <p className="text-lg font-medium">{statusInfo.label}</p>
                            {order.estimatedDelivery && order.status !== 'delivered' && (
                                <p className="text-sm opacity-90">
                                    Dự kiến giao: {new Date(order.estimatedDelivery).toLocaleDateString('vi-VN', { weekday: 'long', month: 'long', day: 'numeric' })}
                                </p>
                            )}
                            {order.deliveredAt && (
                                <p className="text-sm opacity-90">
                                    Đã giao lúc: {formatDate(order.deliveredAt)}
                                </p>
                            )}
                        </div>
                        <Link
                            href={`/account/orders/${orderId}/tracking`}
                            className="px-4 py-2 bg-white text-[#00bfa5] text-sm font-medium rounded-sm hover:bg-gray-100 transition-all"
                        >
                            📍 Xem Chi Tiết Vận Chuyển
                        </Link>
                    </div>

                    {/* Progress Bar */}
                    <div className="relative">
                        <div className="flex justify-between">
                            {['Đã đặt', 'Đã lấy hàng', 'Đang vận chuyển', 'Đang giao', 'Đã giao'].map((step, i) => {
                                const isCompleted = i <= ['pending_payment', 'processing', 'shipped', 'out_for_delivery', 'delivered'].indexOf(order.status);
                                return (
                                    <div key={step} className="flex flex-col items-center z-10">
                                        <div className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold ${isCompleted ? 'bg-white text-[#00bfa5]' : 'bg-white/30 text-white/50'
                                            }`}>
                                            {isCompleted ? '✓' : i + 1}
                                        </div>
                                        <span className={`text-xs mt-1 ${isCompleted ? 'opacity-100' : 'opacity-50'}`}>{step}</span>
                                    </div>
                                );
                            })}
                        </div>
                        <div className="absolute top-3 left-0 right-0 h-0.5 bg-white/30 -z-0" />
                        <div
                            className="absolute top-3 left-0 h-0.5 bg-white transition-all duration-500"
                            style={{ width: `${(['pending_payment', 'processing', 'shipped', 'out_for_delivery', 'delivered'].indexOf(order.status) / 4) * 100}%` }}
                        />
                    </div>
                </div>
            )}

            {/* Shipping Address */}
            <div className="bg-white rounded-sm shadow-sm p-4">
                <h3 className="text-sm font-medium mb-3 flex items-center gap-2 text-[#ee4d2d]">
                    <span>📍</span> Địa Chỉ Nhận Hàng
                </h3>
                <div className="flex items-start gap-3">
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

            {/* Products */}
            <div className="bg-white rounded-sm shadow-sm overflow-hidden">
                <div className="p-4 border-b bg-gray-50 flex items-center justify-between">
                    <div className="flex items-center gap-2">
                        <span className="font-medium text-sm">{order.items[0].shopName}</span>
                        <button className="text-xs text-[#ee4d2d] border border-[#ee4d2d] px-2 py-0.5 hover:bg-[#fef6f5]">
                            💬 Chat
                        </button>
                    </div>
                </div>

                {order.items.map(item => (
                    <div key={item.id} className="p-4 border-b flex gap-4">
                        <div className="relative w-20 h-20 bg-gray-100 rounded-sm overflow-hidden">
                            <Image
                                src={item.image}
                                alt={item.name}
                                fill
                                className="object-cover"
                                unoptimized
                            />
                        </div>
                        <div className="flex-1">
                            <Link href={`/products/${item.productId}`} className="text-sm hover:text-[#ee4d2d]">
                                {item.name}
                            </Link>
                            {item.variant && (
                                <p className="text-xs text-gray-500 mt-1">Phân loại: {item.variant}</p>
                            )}
                            <p className="text-xs text-gray-500">x{item.quantity}</p>
                        </div>
                        <div className="text-right">
                            <p className="text-sm text-[#ee4d2d]">₫{formatPrice(item.price)}</p>
                        </div>
                    </div>
                ))}

                {/* Order Summary */}
                <div className="p-4 bg-[#fafafa] space-y-2">
                    <div className="flex justify-end gap-20 text-sm">
                        <span className="text-gray-500">Tổng tiền hàng</span>
                        <span>₫{formatPrice(order.subtotal)}</span>
                    </div>
                    <div className="flex justify-end gap-20 text-sm">
                        <span className="text-gray-500">Phí vận chuyển</span>
                        <span className={order.shippingFee === 0 ? 'text-[#00bfa5]' : ''}>
                            {order.shippingFee === 0 ? 'Miễn phí' : `₫${formatPrice(order.shippingFee)}`}
                        </span>
                    </div>
                    {order.discount > 0 && (
                        <div className="flex justify-end gap-20 text-sm">
                            <span className="text-gray-500">Giảm giá</span>
                            <span className="text-[#ee4d2d]">-₫{formatPrice(order.discount)}</span>
                        </div>
                    )}
                    <div className="flex justify-end gap-20 text-lg pt-2 border-t">
                        <span className="text-gray-500">Thành tiền</span>
                        <span className="text-[#ee4d2d] font-medium">₫{formatPrice(order.total)}</span>
                    </div>
                </div>
            </div>

            {/* Payment Info */}
            <div className="bg-white rounded-sm shadow-sm p-4">
                <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                        <span className="text-gray-500">Phương thức thanh toán:</span>
                        <span className="ml-2">{order.paymentMethod}</span>
                    </div>
                    <div>
                        <span className="text-gray-500">Đơn vị vận chuyển:</span>
                        <span className="ml-2">{order.shippingMethod}</span>
                    </div>
                    <div>
                        <span className="text-gray-500">Thời gian đặt hàng:</span>
                        <span className="ml-2">{formatDate(order.createdAt)}</span>
                    </div>
                    {order.trackingNumber && (
                        <div>
                            <span className="text-gray-500">Mã vận đơn:</span>
                            <span className="ml-2 text-[#ee4d2d]">{order.trackingNumber}</span>
                        </div>
                    )}
                </div>
            </div>

            {/* Actions */}
            <div className="bg-white rounded-sm shadow-sm p-4 flex justify-end gap-2">
                {order.status === 'delivered' && (
                    <>
                        <button className="px-4 py-2 border text-sm hover:bg-gray-50">Yêu Cầu Trả Hàng/Hoàn Tiền</button>
                        <button className="px-4 py-2 border text-sm hover:bg-gray-50">Liên Hệ Người Bán</button>
                        <button className="px-4 py-2 bg-[#ee4d2d] text-white text-sm hover:opacity-90">Mua Lại</button>
                    </>
                )}
                {['shipped', 'out_for_delivery'].includes(order.status) && (
                    <Link
                        href={`/account/orders/${orderId}/tracking`}
                        className="px-4 py-2 bg-[#ee4d2d] text-white text-sm hover:opacity-90"
                    >
                        Theo Dõi Đơn Hàng
                    </Link>
                )}
            </div>
        </div>
    );
}
