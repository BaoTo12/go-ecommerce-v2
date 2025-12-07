'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import Image from 'next/image';

interface Notification {
    id: string;
    type: 'order' | 'promo' | 'system' | 'update';
    title: string;
    message: string;
    image?: string;
    link?: string;
    read: boolean;
    timestamp: string;
}

const NOTIFICATIONS: Notification[] = [
    {
        id: 'n1',
        type: 'order',
        title: 'Đơn hàng đang giao',
        message: 'Đơn hàng #SP241206001234 đang được giao đến bạn. Dự kiến nhận hàng hôm nay.',
        image: 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=100',
        link: '/account/orders/ord1/tracking',
        read: false,
        timestamp: '2024-12-06T08:30:00',
    },
    {
        id: 'n2',
        type: 'promo',
        title: '🔥 Flash Sale 12.12',
        message: 'Săn deal khủng giảm đến 50%! Voucher freeship x3 đang chờ bạn.',
        link: '/deals/flash-sale',
        read: false,
        timestamp: '2024-12-06T07:00:00',
    },
    {
        id: 'n3',
        type: 'order',
        title: 'Đơn hàng đã giao',
        message: 'Đơn hàng #SP241205005678 đã được giao thành công. Đánh giá để nhận xu!',
        image: 'https://images.unsplash.com/photo-1586495777744-4413f21062fa?w=100',
        link: '/account/orders/ord2',
        read: true,
        timestamp: '2024-12-05T16:45:00',
    },
    {
        id: 'n4',
        type: 'promo',
        title: '🎁 Voucher mới cho bạn',
        message: 'Giảm 50K cho đơn từ 200K. Áp dụng đến 31/12.',
        link: '/account/vouchers',
        read: true,
        timestamp: '2024-12-05T10:00:00',
    },
    {
        id: 'n5',
        type: 'system',
        title: 'Cập nhật bảo mật',
        message: 'Vui lòng xác minh số điện thoại để bảo vệ tài khoản của bạn.',
        link: '/account/profile',
        read: true,
        timestamp: '2024-12-04T14:00:00',
    },
    {
        id: 'n6',
        type: 'update',
        title: 'Shopee Live đang phát!',
        message: 'Apple Store Official đang livestream với nhiều ưu đãi hấp dẫn.',
        link: '/live',
        read: true,
        timestamp: '2024-12-04T12:00:00',
    },
];

export default function NotificationsPage() {
    const [notifications, setNotifications] = useState<Notification[]>(NOTIFICATIONS);
    const [activeTab, setActiveTab] = useState<'all' | 'order' | 'promo' | 'system'>('all');

    const tabs = [
        { key: 'all', label: 'Tất cả' },
        { key: 'order', label: 'Cập nhật đơn hàng' },
        { key: 'promo', label: 'Khuyến mãi' },
        { key: 'system', label: 'Thông báo hệ thống' },
    ];

    const filteredNotifications = activeTab === 'all'
        ? notifications
        : notifications.filter(n => n.type === activeTab);

    const unreadCount = notifications.filter(n => !n.read).length;

    const markAsRead = (id: string) => {
        setNotifications(prev => prev.map(n => n.id === id ? { ...n, read: true } : n));
    };

    const markAllAsRead = () => {
        setNotifications(prev => prev.map(n => ({ ...n, read: true })));
    };

    const formatTime = (timestamp: string) => {
        const date = new Date(timestamp);
        const now = new Date();
        const diff = now.getTime() - date.getTime();
        const hours = Math.floor(diff / (1000 * 60 * 60));
        const days = Math.floor(hours / 24);

        if (hours < 1) return 'Vừa xong';
        if (hours < 24) return `${hours} giờ trước`;
        if (days < 7) return `${days} ngày trước`;
        return date.toLocaleDateString('vi-VN');
    };

    const getTypeIcon = (type: string) => {
        switch (type) {
            case 'order': return '📦';
            case 'promo': return '🎁';
            case 'system': return '⚙️';
            case 'update': return '📢';
            default: return '🔔';
        }
    };

    return (
        <div className="animate-fade-in">
            {/* Header */}
            <div className="bg-white rounded-sm shadow-sm p-4 mb-4 flex items-center justify-between">
                <div>
                    <h1 className="text-lg font-medium">Thông Báo</h1>
                    <p className="text-sm text-gray-500">Bạn có {unreadCount} thông báo chưa đọc</p>
                </div>
                {unreadCount > 0 && (
                    <button
                        onClick={markAllAsRead}
                        className="text-sm text-[#ee4d2d] hover:underline"
                    >
                        Đánh dấu tất cả đã đọc
                    </button>
                )}
            </div>

            {/* Tabs */}
            <div className="bg-white rounded-sm shadow-sm mb-4 overflow-x-auto">
                <div className="flex">
                    {tabs.map(tab => (
                        <button
                            key={tab.key}
                            onClick={() => setActiveTab(tab.key as typeof activeTab)}
                            className={`px-6 py-3 text-sm font-medium border-b-2 transition-colors whitespace-nowrap ${activeTab === tab.key
                                    ? 'text-[#ee4d2d] border-[#ee4d2d]'
                                    : 'text-gray-500 border-transparent hover:text-[#ee4d2d]'
                                }`}
                        >
                            {tab.label}
                        </button>
                    ))}
                </div>
            </div>

            {/* Notifications List */}
            <div className="bg-white rounded-sm shadow-sm overflow-hidden">
                {filteredNotifications.length === 0 ? (
                    <div className="p-12 text-center">
                        <div className="text-5xl mb-4">🔔</div>
                        <p className="text-gray-500">Không có thông báo nào</p>
                    </div>
                ) : (
                    <div className="divide-y">
                        {filteredNotifications.map((notif, index) => (
                            <Link
                                key={notif.id}
                                href={notif.link || '#'}
                                onClick={() => markAsRead(notif.id)}
                                className={`block p-4 hover:bg-gray-50 transition-colors animate-fade-in-up ${!notif.read ? 'bg-[#fef6f5]' : ''
                                    }`}
                                style={{ animationDelay: `${index * 30}ms` }}
                            >
                                <div className="flex gap-4">
                                    {notif.image ? (
                                        <div className="relative w-16 h-16 rounded-sm overflow-hidden flex-shrink-0">
                                            <Image
                                                src={notif.image}
                                                alt=""
                                                fill
                                                className="object-cover"
                                                unoptimized
                                            />
                                        </div>
                                    ) : (
                                        <div className="w-16 h-16 bg-[#fef6f5] rounded-sm flex items-center justify-center text-2xl flex-shrink-0">
                                            {getTypeIcon(notif.type)}
                                        </div>
                                    )}
                                    <div className="flex-1">
                                        <div className="flex items-start justify-between gap-4">
                                            <div>
                                                <h3 className={`text-sm ${!notif.read ? 'font-medium' : ''}`}>
                                                    {notif.title}
                                                </h3>
                                                <p className="text-sm text-gray-500 mt-1">{notif.message}</p>
                                            </div>
                                            {!notif.read && (
                                                <span className="w-2 h-2 bg-[#ee4d2d] rounded-full flex-shrink-0 mt-2" />
                                            )}
                                        </div>
                                        <p className="text-xs text-gray-400 mt-2">{formatTime(notif.timestamp)}</p>
                                    </div>
                                </div>
                            </Link>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
