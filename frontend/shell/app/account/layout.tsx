'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { usePathname } from 'next/navigation';
import { userService, User } from '@/services/userService';

export default function AccountLayout({ children }: { children: React.ReactNode }) {
    const pathname = usePathname();
    const [user, setUser] = useState<User | null>(null);

    useEffect(() => {
        userService.getCurrentUser().then(setUser);
    }, []);

    const menuItems = [
        {
            section: 'Tài Khoản Của Tôi',
            items: [
                { href: '/account/profile', label: 'Hồ Sơ', icon: '👤' },
                { href: '/account/addresses', label: 'Địa Chỉ', icon: '📍' },
                { href: '/account/password', label: 'Đổi Mật Khẩu', icon: '🔒' },
                { href: '/account/notifications', label: 'Cài Đặt Thông Báo', icon: '🔔' },
            ]
        },
        {
            section: 'Đơn Mua',
            items: [
                { href: '/account/orders', label: 'Tất Cả Đơn', icon: '📦' },
                { href: '/account/orders?status=processing', label: 'Đang Xử Lý', icon: '⏳' },
                { href: '/account/orders?status=shipped', label: 'Đang Giao', icon: '🚚' },
                { href: '/account/orders?status=delivered', label: 'Đã Giao', icon: '✅' },
            ]
        },
        {
            section: 'Khác',
            items: [
                { href: '/account/vouchers', label: 'Kho Voucher', icon: '🎟️' },
                { href: '/account/coins', label: 'Shopee Xu', icon: '🪙' },
                { href: '/account/favorites', label: 'Yêu Thích', icon: '❤️' },
                { href: '/account/reviews', label: 'Đánh Giá', icon: '⭐' },
            ]
        },
    ];

    return (
        <div className="min-h-screen bg-[#f5f5f5]">
            <div className="container mx-auto px-4 py-6">
                <div className="flex gap-6">
                    {/* Sidebar */}
                    <aside className="w-[200px] flex-shrink-0 hidden lg:block">
                        {/* User Info */}
                        {user && (
                            <div className="flex items-center gap-3 mb-6 animate-fade-in">
                                <div className="relative w-12 h-12 rounded-full overflow-hidden border-2 border-[#ee4d2d]">
                                    <Image
                                        src={user.avatar}
                                        alt={user.name}
                                        fill
                                        className="object-cover"
                                        unoptimized
                                    />
                                </div>
                                <div>
                                    <p className="font-semibold text-sm">{user.name}</p>
                                    <Link href="/account/profile" className="text-xs text-gray-500 hover:text-[#ee4d2d] flex items-center gap-1">
                                        ✏️ Sửa Hồ Sơ
                                    </Link>
                                </div>
                            </div>
                        )}

                        {/* Menu */}
                        <nav className="space-y-4">
                            {menuItems.map((section, sIndex) => (
                                <div key={section.section} className="animate-fade-in-left" style={{ animationDelay: `${sIndex * 100}ms` }}>
                                    <h3 className="text-sm font-medium mb-2 flex items-center gap-2">
                                        {section.section}
                                    </h3>
                                    <ul className="space-y-1">
                                        {section.items.map(item => {
                                            const isActive = pathname === item.href || (pathname === '/account/orders' && item.href.startsWith('/account/orders'));
                                            return (
                                                <li key={item.href}>
                                                    <Link
                                                        href={item.href}
                                                        className={`flex items-center gap-2 py-2 px-3 text-sm rounded-sm transition-all ${isActive
                                                                ? 'text-[#ee4d2d] bg-[#fef6f5]'
                                                                : 'text-gray-600 hover:text-[#ee4d2d] hover:bg-gray-50'
                                                            }`}
                                                    >
                                                        <span>{item.icon}</span>
                                                        <span>{item.label}</span>
                                                    </Link>
                                                </li>
                                            );
                                        })}
                                    </ul>
                                </div>
                            ))}
                        </nav>
                    </aside>

                    {/* Main Content */}
                    <main className="flex-1 min-w-0">
                        {children}
                    </main>
                </div>
            </div>
        </div>
    );
}
