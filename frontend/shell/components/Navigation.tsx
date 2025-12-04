'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useState } from 'react';
import { useGamificationStore, useCartStore, useAuthStore } from '../lib/store';

const navigation = [
    {
        label: 'Shop',
        items: [
            { href: '/', label: 'Home', icon: '🏠' },
            { href: '/products', label: 'Products', icon: '🛍️' },
            { href: '/live', label: 'Live Shopping', icon: '🔴' },
        ],
    },
    {
        label: 'Deals',
        items: [
            { href: '/deals/flash-sale', label: 'Flash Sale', icon: '⚡' },
            { href: '/deals/coupons', label: 'Coupons', icon: '🎟️' },
        ],
    },
    {
        label: 'Rewards',
        items: [
            { href: '/rewards', label: 'Gamification', icon: '🎮' },
        ],
    },
    {
        label: 'Admin',
        items: [
            { href: '/admin/analytics', label: 'Analytics', icon: '📊' },
            { href: '/admin/fraud', label: 'Fraud Detection', icon: '🛡️' },
            { href: '/admin/pricing', label: 'Dynamic Pricing', icon: '💹' },
        ],
    },
];

export default function Navigation() {
    const pathname = usePathname();
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
    const { balance } = useGamificationStore();
    const { items } = useCartStore();
    const { isAuthenticated, user, login } = useAuthStore();

    const handleMockLogin = () => {
        login({
            id: 'user-123',
            name: 'Demo User',
            email: 'demo@titan.com',
        });
    };

    return (
        <header className="sticky top-0 z-50 border-b bg-white/80 backdrop-blur-lg">
            <div className="container mx-auto">
                <div className="flex h-16 items-center justify-between px-4">
                    {/* Logo */}
                    <Link href="/" className="flex items-center gap-2">
                        <span className="text-2xl">🚀</span>
                        <span className="text-xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                            Titan Commerce
                        </span>
                    </Link>

                    {/* Desktop Navigation */}
                    <nav className="hidden md:flex items-center gap-6">
                        {navigation.map((group) => (
                            <div key={group.label} className="relative group">
                                <button className="text-sm font-medium text-gray-700 hover:text-gray-900 py-2">
                                    {group.label}
                                </button>
                                <div className="absolute top-full left-0 pt-2 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all">
                                    <div className="rounded-lg border bg-white shadow-lg p-2 min-w-[180px]">
                                        {group.items.map((item) => (
                                            <Link
                                                key={item.href}
                                                href={item.href}
                                                className={`flex items-center gap-2 rounded-md px-3 py-2 text-sm hover:bg-gray-100 ${pathname === item.href ? 'bg-blue-50 text-blue-600' : ''
                                                    }`}
                                            >
                                                <span>{item.icon}</span>
                                                {item.label}
                                            </Link>
                                        ))}
                                    </div>
                                </div>
                            </div>
                        ))}
                    </nav>

                    {/* Right side */}
                    <div className="flex items-center gap-4">
                        {/* Coins */}
                        <div className="hidden sm:flex items-center gap-1 rounded-full bg-yellow-100 px-3 py-1">
                            <span>🪙</span>
                            <span className="font-semibold text-yellow-700">{balance}</span>
                        </div>

                        {/* Cart */}
                        <Link
                            href="/cart"
                            className="relative rounded-full bg-gray-100 p-2 hover:bg-gray-200"
                        >
                            <span className="text-xl">🛒</span>
                            {items.length > 0 && (
                                <span className="absolute -right-1 -top-1 flex h-5 w-5 items-center justify-center rounded-full bg-red-500 text-xs font-bold text-white">
                                    {items.length}
                                </span>
                            )}
                        </Link>

                        {/* User */}
                        {isAuthenticated ? (
                            <div className="flex items-center gap-2 rounded-full bg-gray-100 px-3 py-1">
                                <span>👤</span>
                                <span className="text-sm font-medium">{user?.name}</span>
                            </div>
                        ) : (
                            <button
                                onClick={handleMockLogin}
                                className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-700"
                            >
                                Login
                            </button>
                        )}

                        {/* Mobile menu button */}
                        <button
                            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
                            className="md:hidden rounded-lg p-2 hover:bg-gray-100"
                        >
                            <span className="text-xl">{mobileMenuOpen ? '✕' : '☰'}</span>
                        </button>
                    </div>
                </div>

                {/* Mobile Navigation */}
                {mobileMenuOpen && (
                    <div className="md:hidden border-t py-4 px-4">
                        {navigation.map((group) => (
                            <div key={group.label} className="mb-4">
                                <div className="text-xs font-semibold text-gray-500 uppercase mb-2">
                                    {group.label}
                                </div>
                                <div className="space-y-1">
                                    {group.items.map((item) => (
                                        <Link
                                            key={item.href}
                                            href={item.href}
                                            onClick={() => setMobileMenuOpen(false)}
                                            className={`flex items-center gap-2 rounded-md px-3 py-2 hover:bg-gray-100 ${pathname === item.href ? 'bg-blue-50 text-blue-600' : ''
                                                }`}
                                        >
                                            <span>{item.icon}</span>
                                            {item.label}
                                        </Link>
                                    ))}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </header>
    );
}
