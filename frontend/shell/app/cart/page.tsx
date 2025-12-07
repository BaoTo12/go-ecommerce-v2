'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useRouter } from 'next/navigation';
import { cartService, CartItem, Cart } from '@/services/cartService';

export default function CartPage() {
    const router = useRouter();
    const [cart, setCart] = useState<Cart | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [notification, setNotification] = useState<string | null>(null);
    const [selectAll, setSelectAll] = useState(true);

    const loadCart = async () => {
        const data = await cartService.getCart();
        setCart(data);
        setIsLoading(false);
    };

    useEffect(() => {
        loadCart();
    }, []);

    const updateQuantity = async (itemId: string, quantity: number) => {
        await cartService.updateQuantity(itemId, quantity);
        loadCart();
    };

    const removeItem = async (itemId: string) => {
        await cartService.removeItem(itemId);
        loadCart();
        setNotification('✓ Đã xóa sản phẩm');
        setTimeout(() => setNotification(null), 2000);
    };

    const toggleSelection = async (itemId: string) => {
        await cartService.toggleSelection(itemId);
        loadCart();
    };

    const toggleSelectAll = async () => {
        await cartService.selectAll(!selectAll);
        setSelectAll(!selectAll);
        loadCart();
    };

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    const totals = cart ? cartService.calculateTotals() : { subtotal: 0, itemCount: 0, selectedCount: 0 };

    if (isLoading) {
        return (
            <div className="min-h-screen bg-[#f5f5f5] flex items-center justify-center">
                <div className="loading-spinner" />
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-[#f5f5f5] animate-fade-in">
            {/* Toast */}
            {notification && <div className="toast toast-success">{notification}</div>}

            <div className="container mx-auto px-4 py-6">
                <h1 className="text-xl font-medium mb-4">Giỏ Hàng</h1>

                {!cart || cart.items.length === 0 ? (
                    <div className="bg-white rounded-sm shadow-sm p-12 text-center">
                        <div className="text-6xl mb-4 animate-float">🛒</div>
                        <p className="text-gray-500 mb-4">Giỏ hàng của bạn còn trống</p>
                        <Link
                            href="/products"
                            className="inline-block px-8 py-3 bg-[#ee4d2d] text-white hover:opacity-90 transition-all"
                        >
                            Mua Sắm Ngay
                        </Link>
                    </div>
                ) : (
                    <>
                        {/* Cart Header */}
                        <div className="bg-white rounded-sm shadow-sm p-4 mb-3 hidden md:grid grid-cols-12 gap-4 text-sm text-gray-500">
                            <div className="col-span-6 flex items-center gap-4">
                                <input
                                    type="checkbox"
                                    checked={selectAll}
                                    onChange={toggleSelectAll}
                                    className="w-4 h-4 accent-[#ee4d2d]"
                                />
                                <span>Sản Phẩm</span>
                            </div>
                            <div className="col-span-2 text-center">Đơn Giá</div>
                            <div className="col-span-2 text-center">Số Lượng</div>
                            <div className="col-span-1 text-center">Số Tiền</div>
                            <div className="col-span-1 text-center">Thao Tác</div>
                        </div>

                        {/* Cart Items */}
                        <div className="space-y-3 mb-4">
                            {cart.items.map((item, index) => (
                                <div
                                    key={item.id}
                                    className="bg-white rounded-sm shadow-sm p-4 grid grid-cols-12 gap-4 items-center animate-fade-in-up"
                                    style={{ animationDelay: `${index * 50}ms` }}
                                >
                                    <div className="col-span-6 flex items-center gap-4">
                                        <input
                                            type="checkbox"
                                            checked={item.selected}
                                            onChange={() => toggleSelection(item.id)}
                                            className="w-4 h-4 accent-[#ee4d2d] flex-shrink-0"
                                        />
                                        <div className="relative w-20 h-20 bg-gray-100 rounded-sm overflow-hidden flex-shrink-0">
                                            <Image
                                                src={item.product.thumbnail}
                                                alt={item.product.name}
                                                fill
                                                className="object-cover"
                                                unoptimized
                                            />
                                        </div>
                                        <div className="flex-1 min-w-0">
                                            <Link href={`/products/${item.productId}`} className="text-sm line-clamp-2 hover:text-[#ee4d2d]">
                                                {item.product.name}
                                            </Link>
                                            {item.variant && (
                                                <p className="text-xs text-gray-400 mt-1">Phân loại: {item.variant}</p>
                                            )}
                                        </div>
                                    </div>

                                    <div className="col-span-2 text-center">
                                        <span className="text-sm text-[#ee4d2d]">₫{formatPrice(item.product.price)}</span>
                                    </div>

                                    <div className="col-span-2 flex justify-center">
                                        <div className="flex items-center">
                                            <button
                                                onClick={() => updateQuantity(item.id, item.quantity - 1)}
                                                className="w-8 h-8 border flex items-center justify-center hover:bg-gray-50"
                                            >
                                                −
                                            </button>
                                            <input
                                                type="number"
                                                value={item.quantity}
                                                onChange={(e) => updateQuantity(item.id, parseInt(e.target.value) || 1)}
                                                className="w-12 h-8 border-y text-center text-sm outline-none"
                                            />
                                            <button
                                                onClick={() => updateQuantity(item.id, item.quantity + 1)}
                                                className="w-8 h-8 border flex items-center justify-center hover:bg-gray-50"
                                            >
                                                +
                                            </button>
                                        </div>
                                    </div>

                                    <div className="col-span-1 text-center">
                                        <span className="text-sm text-[#ee4d2d] font-medium">
                                            ₫{formatPrice(item.product.price * item.quantity)}
                                        </span>
                                    </div>

                                    <div className="col-span-1 text-center">
                                        <button
                                            onClick={() => removeItem(item.id)}
                                            className="text-sm text-gray-500 hover:text-[#ee4d2d] transition-colors"
                                        >
                                            Xóa
                                        </button>
                                    </div>
                                </div>
                            ))}
                        </div>

                        {/* Cart Footer */}
                        <div className="bg-white rounded-sm shadow-sm p-4 sticky bottom-0 z-40 border-t">
                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-4">
                                    <label className="flex items-center gap-2 cursor-pointer">
                                        <input
                                            type="checkbox"
                                            checked={selectAll}
                                            onChange={toggleSelectAll}
                                            className="w-4 h-4 accent-[#ee4d2d]"
                                        />
                                        <span className="text-sm">Chọn Tất Cả ({cart.items.length})</span>
                                    </label>
                                    <button className="text-sm text-gray-500 hover:text-[#ee4d2d]">
                                        Xóa
                                    </button>
                                </div>

                                <div className="flex items-center gap-6">
                                    <div className="text-right">
                                        <span className="text-sm text-gray-500">
                                            Tổng thanh toán ({totals.selectedCount} Sản phẩm):
                                        </span>
                                        <span className="text-2xl text-[#ee4d2d] font-medium ml-2">
                                            ₫{formatPrice(totals.subtotal)}
                                        </span>
                                    </div>
                                    <button
                                        onClick={() => router.push('/checkout')}
                                        disabled={totals.selectedCount === 0}
                                        className={`px-12 py-3 bg-[#ee4d2d] text-white font-medium hover:opacity-90 transition-all ${totals.selectedCount === 0 ? 'opacity-50 cursor-not-allowed' : ''
                                            }`}
                                    >
                                        Mua Hàng
                                    </button>
                                </div>
                            </div>
                        </div>
                    </>
                )}
            </div>
        </div>
    );
}
