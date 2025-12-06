'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

export default function Navigation() {
    const router = useRouter();
    const [searchQuery, setSearchQuery] = useState('');
    const [cartCount] = useState(3);

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        if (searchQuery.trim()) {
            router.push(`/products?search=${encodeURIComponent(searchQuery)}`);
        }
    };

    const searchSuggestions = ['iPhone 15', 'Áo hoodie', 'Tai nghe bluetooth', 'Son môi', 'Giày sneaker'];

    return (
        <header className="sticky top-0 z-50">
            {/* Top Navigation Bar */}
            <div className="bg-gradient-to-b from-[#f53d2d] to-[#f63] text-white">
                <div className="container mx-auto px-4">
                    {/* Top links row */}
                    <div className="flex items-center justify-between h-[34px] text-[13px]">
                        <div className="flex items-center gap-4">
                            <Link href="#" className="hover:opacity-80">Kênh Người Bán</Link>
                            <Link href="#" className="hover:opacity-80">Trở thành Người bán Shopee</Link>
                            <span className="border-l border-white/30 h-3" />
                            <Link href="#" className="hover:opacity-80 flex items-center gap-1">
                                <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                                    <path d="M17.05 20.28c-.98.95-2.05.8-3.08.35-1.09-.46-2.09-.48-3.24 0-1.44.62-2.2.44-3.06-.35C2.79 15.25 3.51 7.59 9.05 7.31c1.35.07 2.29.74 3.08.8 1.18-.24 2.31-.93 3.57-.84 1.51.12 2.65.72 3.4 1.8-3.12 1.87-2.38 5.98.48 7.13-.57 1.5-1.31 2.99-2.54 4.09l.01-.01zM12.03 7.25c-.15-2.23 1.66-4.07 3.74-4.25.29 2.58-2.34 4.5-3.74 4.25z" />
                                </svg>
                                Tải ứng dụng
                            </Link>
                            <span className="border-l border-white/30 h-3" />
                            <span className="flex items-center gap-2">
                                Kết nối
                                <Link href="#" className="hover:opacity-80">
                                    <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                                        <path d="M12 2C6.477 2 2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.879V14.89h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.89h-2.33v6.989C18.343 21.129 22 16.99 22 12c0-5.523-4.477-10-10-10z" />
                                    </svg>
                                </Link>
                                <Link href="#" className="hover:opacity-80">
                                    <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                                        <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zM12 0C8.741 0 8.333.014 7.053.072 2.695.272.273 2.69.073 7.052.014 8.333 0 8.741 0 12c0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98C8.333 23.986 8.741 24 12 24c3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98C15.668.014 15.259 0 12 0zm0 5.838a6.162 6.162 0 100 12.324 6.162 6.162 0 000-12.324zM12 16a4 4 0 110-8 4 4 0 010 8zm6.406-11.845a1.44 1.44 0 100 2.881 1.44 1.44 0 000-2.881z" />
                                    </svg>
                                </Link>
                            </span>
                        </div>
                        <div className="flex items-center gap-4">
                            <Link href="#" className="hover:opacity-80 flex items-center gap-1">
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                                </svg>
                                Thông Báo
                            </Link>
                            <Link href="#" className="hover:opacity-80 flex items-center gap-1">
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                </svg>
                                Hỗ Trợ
                            </Link>
                            <span className="border-l border-white/30 h-3" />
                            <Link href="#" className="hover:opacity-80">Đăng Ký</Link>
                            <span className="border-l border-white/30 h-3" />
                            <Link href="#" className="hover:opacity-80">Đăng Nhập</Link>
                        </div>
                    </div>

                    {/* Main nav row */}
                    <div className="flex items-center gap-8 h-[85px]">
                        {/* Logo */}
                        <Link href="/" className="flex items-center gap-2 flex-shrink-0">
                            <svg className="h-[50px] w-auto" viewBox="0 0 200 60" fill="white">
                                <text x="0" y="45" fontFamily="Helvetica, Arial" fontSize="42" fontWeight="bold">Shopee</text>
                            </svg>
                        </Link>

                        {/* Search */}
                        <div className="flex-1">
                            <form onSubmit={handleSearch} className="relative">
                                <div className="flex bg-white rounded-sm overflow-hidden">
                                    <input
                                        type="text"
                                        value={searchQuery}
                                        onChange={(e) => setSearchQuery(e.target.value)}
                                        placeholder="Đăng kí và nhận voucher bạn mới đến 70k!"
                                        className="flex-1 px-4 py-[10px] text-[14px] text-gray-700 outline-none"
                                    />
                                    <button
                                        type="submit"
                                        className="px-5 bg-[#fb5533] hover:opacity-90 flex items-center justify-center"
                                    >
                                        <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                                        </svg>
                                    </button>
                                </div>
                                {/* Search suggestions */}
                                <div className="absolute top-full left-0 right-0 mt-1 flex gap-2 overflow-x-auto text-xs">
                                    {searchSuggestions.map(s => (
                                        <Link key={s} href={`/products?search=${s}`} className="text-white/80 hover:text-white whitespace-nowrap">
                                            {s}
                                        </Link>
                                    ))}
                                </div>
                            </form>
                        </div>

                        {/* Cart */}
                        <Link href="/cart" className="relative p-2 hover:opacity-80">
                            <svg className="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
                            </svg>
                            {cartCount > 0 && (
                                <span className="absolute top-0 right-0 bg-white text-[#ee4d2d] text-xs font-bold rounded-full min-w-[18px] h-[18px] flex items-center justify-center px-1">
                                    {cartCount}
                                </span>
                            )}
                        </Link>
                    </div>
                </div>
            </div>

            {/* Secondary Navigation */}
            <div className="bg-white border-b shadow-sm">
                <div className="container mx-auto px-4">
                    <nav className="flex items-center gap-1 h-10 text-[13px]">
                        <Link href="/" className="px-3 py-2 hover:text-[#ee4d2d] transition-colors">
                            Trang Chủ
                        </Link>
                        <Link href="/products" className="px-3 py-2 hover:text-[#ee4d2d] transition-colors">
                            Sản Phẩm
                        </Link>
                        <Link href="/deals/flash-sale" className="px-3 py-2 hover:text-[#ee4d2d] transition-colors flex items-center gap-1">
                            <span className="text-[#ee4d2d]">⚡</span> Flash Sale
                        </Link>
                        <Link href="/live" className="px-3 py-2 hover:text-[#ee4d2d] transition-colors flex items-center gap-1">
                            <span className="w-2 h-2 bg-red-500 rounded-full animate-pulse" /> Shopee Live
                        </Link>
                        <Link href="/rewards" className="px-3 py-2 hover:text-[#ee4d2d] transition-colors">
                            🎮 Shopee Xu
                        </Link>
                        <Link href="/deals/coupons" className="px-3 py-2 hover:text-[#ee4d2d] transition-colors">
                            Mã Giảm Giá
                        </Link>
                        <div className="ml-auto flex items-center gap-1 text-gray-500">
                            <Link href="/admin/analytics" className="px-3 py-2 hover:text-[#ee4d2d] transition-colors">
                                📊 Analytics
                            </Link>
                            <Link href="/admin/fraud" className="px-3 py-2 hover:text-[#ee4d2d] transition-colors">
                                🛡️ Fraud
                            </Link>
                            <Link href="/admin/pricing" className="px-3 py-2 hover:text-[#ee4d2d] transition-colors">
                                💹 Pricing
                            </Link>
                        </div>
                    </nav>
                </div>
            </div>
        </header>
    );
}
