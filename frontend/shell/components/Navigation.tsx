'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const navigation = [
    { href: '/', label: 'Trang chủ', icon: '🏠' },
    { href: '/products', label: 'Sản phẩm', icon: '🛍️' },
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
    const [coins, setCoins] = useState(1250);
    const [cartCount, setCartCount] = useState(3);
    const [searchQuery, setSearchQuery] = useState('');
    const [showSearch, setShowSearch] = useState(false);
    const [notifications, setNotifications] = useState(2);

    // Simulate coin updates
    useEffect(() => {
        const timer = setInterval(() => {
            setCoins(prev => prev + Math.floor(Math.random() * 3));
        }, 10000);
        return () => clearInterval(timer);
    }, []);

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        if (searchQuery.trim()) {
            window.location.href = `/products?search=${encodeURIComponent(searchQuery)}`;
        }
    };

    return (
        <header className="sticky top-0 z-50">
            {/* Top bar */}
            <div className="bg-gradient-to-r from-[#F63] to-[#EE4D2D] text-white text-xs">
                <div className="container mx-auto px-4">
                    <div className="flex items-center justify-between h-8">
                        <div className="flex items-center gap-4">
                            <span>📱 Tải ứng dụng</span>
                            <span className="hidden md:inline">|</span>
                            <span className="hidden md:inline">Kết nối: 📘 📸 🐦</span>
                        </div>
                        <div className="flex items-center gap-4">
                            <button className="flex items-center gap-1 hover:opacity-80">
                                <span>🔔</span>
                                Thông báo
                                {notifications > 0 && (
                                    <span className="bg-yellow-400 text-[#EE4D2D] text-[10px] font-bold px-1.5 rounded-full">
                                        {notifications}
                                    </span>
                                )}
                            </button>
                            <span className="hidden md:inline">|</span>
                            <button className="hidden md:flex items-center gap-1 hover:opacity-80">
                                <span>❓</span>
                                Hỗ trợ
                            </button>
                            <span>|</span>
                            <div className="flex items-center gap-2">
                                <span className="w-6 h-6 bg-white/20 rounded-full flex items-center justify-center">👤</span>
                                <span>Đăng nhập</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Main nav */}
            <div className="bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] shadow-md">
                <div className="container mx-auto px-4">
                    <div className="flex items-center justify-between h-16 gap-4">
                        {/* Logo */}
                        <Link href="/" className="flex items-center gap-2 flex-shrink-0">
                            <span className="text-white text-3xl font-bold tracking-tight">Shopee</span>
                        </Link>

                        {/* Search Bar */}
                        <form onSubmit={handleSearch} className="hidden md:flex flex-1 max-w-2xl">
                            <div className="relative w-full">
                                <input
                                    type="text"
                                    value={searchQuery}
                                    onChange={(e) => setSearchQuery(e.target.value)}
                                    placeholder="Tìm kiếm sản phẩm, thương hiệu, và danh mục..."
                                    className="w-full py-2.5 px-4 pr-24 rounded-sm text-sm focus:outline-none"
                                />
                                <button
                                    type="submit"
                                    className="absolute right-0 top-0 h-full px-6 bg-[#FB6445] text-white rounded-r-sm hover:bg-[#EE4D2D] flex items-center gap-1"
                                >
                                    <span>🔍</span>
                                </button>
                            </div>
                        </form>

                        {/* Right side */}
                        <div className="flex items-center gap-4">
                            {/* Mobile search toggle */}
                            <button
                                onClick={() => setShowSearch(!showSearch)}
                                className="md:hidden text-white text-xl hover:opacity-80"
                            >
                                🔍
                            </button>

                            {/* Coins */}
                            <Link
                                href="/rewards"
                                className="hidden sm:flex items-center gap-1 bg-white/10 text-white px-3 py-1.5 rounded hover:bg-white/20 transition-colors"
                            >
                                <span className="text-yellow-300">🪙</span>
                                <span className="font-semibold">{coins.toLocaleString()}</span>
                            </Link>

                            {/* Cart */}
                            <Link href="/cart" className="relative text-white hover:opacity-80">
                                <span className="text-2xl">🛒</span>
                                {cartCount > 0 && (
                                    <span className="absolute -top-1 -right-1 bg-white text-[#EE4D2D] text-xs font-bold rounded-full min-w-[18px] h-[18px] flex items-center justify-center px-1">
                                        {cartCount}
                                    </span>
                                )}
                            </Link>

                            {/* Mobile menu */}
                            <button
                                onClick={() => setMobileOpen(!mobileOpen)}
                                className="md:hidden text-white text-2xl hover:opacity-80"
                            >
                                {mobileOpen ? '✕' : '☰'}
                            </button>
                        </div>
                    </div>

                    {/* Mobile search bar */}
                    {showSearch && (
                        <form onSubmit={handleSearch} className="md:hidden pb-3">
                            <div className="relative">
                                <input
                                    type="text"
                                    value={searchQuery}
                                    onChange={(e) => setSearchQuery(e.target.value)}
                                    placeholder="Tìm kiếm..."
                                    className="w-full py-2 px-4 rounded-sm text-sm focus:outline-none"
                                    autoFocus
                                />
                                <button type="submit" className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400">
                                    🔍
                                </button>
                            </div>
                        </form>
                    )}

                    {/* Desktop Navigation */}
                    <div className="hidden md:flex items-center h-10">
                        {navigation.map(item => (
                            <Link
                                key={item.href}
                                href={item.href}
                                className={`flex items-center gap-1.5 px-4 text-sm text-white hover:bg-white/10 h-full transition-colors ${pathname === item.href ? 'font-bold bg-white/10' : ''
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
                                    className={`flex items-center gap-1 px-3 text-sm text-white hover:bg-white/10 h-full transition-colors ${pathname === item.href ? 'font-bold bg-white/10' : ''
                                        }`}
                                >
                                    <span>{item.icon}</span>
                                    {item.label}
                                </Link>
                            ))}
                        </div>
                    </div>
                </div>
            </div>

            {/* Mobile Navigation */}
            {mobileOpen && (
                <div className="md:hidden bg-white border-t shadow-lg">
                    <div className="p-4 space-y-1">
                        {/* User info */}
                        <div className="flex items-center gap-3 p-3 bg-[#FFEEE8] rounded mb-3">
                            <span className="w-10 h-10 bg-[#EE4D2D] text-white rounded-full flex items-center justify-center text-xl">👤</span>
                            <div>
                                <p className="font-semibold">Xin chào!</p>
                                <p className="text-sm text-gray-500">Đăng nhập để mua sắm</p>
                            </div>
                        </div>

                        {/* Coins */}
                        <Link
                            href="/rewards"
                            onClick={() => setMobileOpen(false)}
                            className="flex items-center justify-between p-3 bg-yellow-50 rounded"
                        >
                            <span className="flex items-center gap-2">
                                <span className="text-yellow-500">🪙</span>
                                Shopee Xu
                            </span>
                            <span className="font-bold text-[#EE4D2D]">{coins.toLocaleString()}</span>
                        </Link>

                        {/* Nav items */}
                        {[...navigation, ...adminNav].map(item => (
                            <Link
                                key={item.href}
                                href={item.href}
                                onClick={() => setMobileOpen(false)}
                                className={`flex items-center gap-3 px-3 py-3 rounded transition-colors ${pathname === item.href
                                        ? 'bg-[#FFEEE8] text-[#EE4D2D] font-semibold'
                                        : 'hover:bg-gray-100'
                                    }`}
                            >
                                <span className="text-xl">{item.icon}</span>
                                {item.label}
                            </Link>
                        ))}
                    </div>
                </div>
            )}
        </header>
    );
}
