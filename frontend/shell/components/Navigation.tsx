'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const navigation = [
    { href: '/', label: 'Trang chủ', icon: '🏠' },
    { href: '/deals/flash-sale', label: 'Flash Sale', icon: '⚡' },
    { href: '/live', label: 'Shopee Live', icon: '🔴' },
    { href: '/rewards', label: 'Xu & Game', icon: '🎮' },
    { href: '/deals/coupons', label: 'Mã giảm giá', icon: '🎟️' },
];

const adminNav = [
    { href: '/admin/analytics', label: 'Analytics', icon: '📊' },
    { href: '/admin/fraud', label: 'Fraud', icon: '🛡️' },
    { href: '/admin/pricing', label: 'Pricing', icon: '💹' },
];

export default function Navigation() {
    const pathname = usePathname();
    const [mobileOpen, setMobileOpen] = useState(false);
    const [coins] = useState(1250);

    return (
        <header className="sticky top-0 z-50 bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] shadow-md">
            <div className="container mx-auto">
                {/* Main nav */}
                <div className="flex items-center justify-between h-14 px-4">
                    {/* Logo */}
                    <Link href="/" className="flex items-center gap-2">
                        <span className="text-white text-2xl font-bold tracking-tight">Shopee</span>
                    </Link>

                    {/* Search Bar */}
                    <div className="hidden md:flex flex-1 max-w-xl mx-8">
                        <div className="relative w-full">
                            <input
                                type="text"
                                placeholder="Tìm kiếm sản phẩm..."
                                className="w-full py-2 px-4 pr-12 rounded-sm text-sm focus:outline-none"
                            />
                            <button className="absolute right-0 top-0 h-full px-4 bg-[#FB6445] text-white rounded-r-sm hover:bg-[#EE4D2D]">
                                🔍
                            </button>
                        </div>
                    </div>

                    {/* Right side */}
                    <div className="flex items-center gap-4">
                        {/* Coins */}
                        <div className="hidden sm:flex items-center gap-1 bg-white/10 text-white px-3 py-1 rounded">
                            <span className="text-yellow-300">🪙</span>
                            <span className="font-semibold">{coins.toLocaleString()}</span>
                        </div>

                        {/* Cart */}
                        <Link href="/cart" className="relative text-white hover:opacity-80">
                            <span className="text-xl">🛒</span>
                            <span className="absolute -top-1 -right-1 bg-white text-[#EE4D2D] text-xs font-bold rounded-full w-4 h-4 flex items-center justify-center">
                                3
                            </span>
                        </Link>

                        {/* User */}
                        <div className="hidden sm:flex items-center gap-2 text-white">
                            <span>👤</span>
                            <span className="text-sm">Đăng nhập</span>
                        </div>

                        {/* Mobile menu */}
                        <button
                            onClick={() => setMobileOpen(!mobileOpen)}
                            className="md:hidden text-white text-xl"
                        >
                            {mobileOpen ? '✕' : '☰'}
                        </button>
                    </div>
                </div>

                {/* Desktop Navigation */}
                <div className="hidden md:flex items-center h-10 px-4 bg-white/10">
                    {navigation.map(item => (
                        <Link
                            key={item.href}
                            href={item.href}
                            className={`flex items-center gap-1 px-4 text-sm text-white hover:opacity-80 ${pathname === item.href ? 'font-bold' : ''
                                }`}
                        >
                            <span>{item.icon}</span>
                            {item.label}
                        </Link>
                    ))}
                    <div className="ml-auto flex items-center">
                        <span className="text-white/50 text-sm mr-2">Admin:</span>
                        {adminNav.map(item => (
                            <Link
                                key={item.href}
                                href={item.href}
                                className={`flex items-center gap-1 px-3 text-sm text-white hover:opacity-80 ${pathname === item.href ? 'font-bold' : ''
                                    }`}
                            >
                                <span>{item.icon}</span>
                                {item.label}
                            </Link>
                        ))}
                    </div>
                </div>

                {/* Mobile Navigation */}
                {mobileOpen && (
                    <div className="md:hidden bg-white border-t">
                        <div className="p-4 space-y-2">
                            {[...navigation, ...adminNav].map(item => (
                                <Link
                                    key={item.href}
                                    href={item.href}
                                    onClick={() => setMobileOpen(false)}
                                    className={`flex items-center gap-2 px-4 py-2 rounded ${pathname === item.href ? 'bg-[#FFEEE8] text-[#EE4D2D]' : 'hover:bg-gray-100'
                                        }`}
                                >
                                    <span>{item.icon}</span>
                                    {item.label}
                                </Link>
                            ))}
                        </div>
                    </div>
                )}
            </div>
        </header>
    );
}
