'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

interface CartItem {
    id: string;
    name: string;
    price: number;
    originalPrice: number;
    quantity: number;
    image: string;
    shop: string;
    variant?: string;
    selected: boolean;
}

export default function CartPage() {
    const router = useRouter();
    const [items, setItems] = useState<CartItem[]>([]);
    const [voucher, setVoucher] = useState('');
    const [voucherApplied, setVoucherApplied] = useState(false);
    const [discount, setDiscount] = useState(0);

    useEffect(() => {
        setItems([
            { id: 'p1', name: 'iPhone 15 Pro Max 256GB Titan Xanh Chính Hãng VN/A', price: 29990000, originalPrice: 34990000, quantity: 1, image: '📱', shop: 'Apple Store Official', variant: 'Xanh Titan', selected: true },
            { id: 'p6', name: 'Son Dưỡng Môi Dior Addict Lip Glow', price: 950000, originalPrice: 1200000, quantity: 2, image: '💄', shop: 'Dior Beauty Official', variant: '001 Pink', selected: true },
            { id: 'p4', name: 'Áo Hoodie Unisex Form Rộng Nỉ Cotton', price: 199000, originalPrice: 350000, quantity: 1, image: '👕', shop: 'Fashion Store', variant: 'Đen - L', selected: false },
        ]);
    }, []);

    const updateQuantity = (id: string, delta: number) => {
        setItems(prev => prev.map(item => {
            if (item.id === id) {
                return { ...item, quantity: Math.max(1, item.quantity + delta) };
            }
            return item;
        }));
    };

    const toggleSelect = (id: string) => {
        setItems(prev => prev.map(item =>
            item.id === id ? { ...item, selected: !item.selected } : item
        ));
    };

    const selectAll = (checked: boolean) => {
        setItems(prev => prev.map(item => ({ ...item, selected: checked })));
    };

    const removeItem = (id: string) => {
        setItems(prev => prev.filter(item => item.id !== id));
    };

    const applyVoucher = () => {
        if (voucher.toUpperCase() === 'GIAM50K') {
            setDiscount(50000);
            setVoucherApplied(true);
        } else {
            alert('Mã giảm giá không hợp lệ');
        }
    };

    const selectedItems = items.filter(item => item.selected);
    const subtotal = selectedItems.reduce((sum, item) => sum + item.price * item.quantity, 0);
    const total = subtotal - discount;

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    const proceedToCheckout = () => {
        if (selectedItems.length > 0) {
            router.push('/checkout');
        }
    };

    // Group items by shop
    const groupedItems = items.reduce((acc, item) => {
        if (!acc[item.shop]) acc[item.shop] = [];
        acc[item.shop].push(item);
        return acc;
    }, {} as Record<string, CartItem[]>);

    return (
        <div className="min-h-screen bg-[#f5f5f5]">
            {/* Header */}
            <div className="bg-white border-b sticky top-0 z-30">
                <div className="container mx-auto px-4 py-4">
                    <div className="flex items-center gap-3">
                        <Link href="/" className="text-2xl font-bold text-[#ee4d2d]">Shopee</Link>
                        <span className="text-gray-300">|</span>
                        <h1 className="text-xl text-gray-700">Giỏ Hàng</h1>
                    </div>
                </div>
            </div>

            <div className="container mx-auto px-4 py-4">
                {items.length === 0 ? (
                    <div className="bg-white rounded-sm p-12 text-center">
                        <div className="w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                            <svg className="w-12 h-12 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
                            </svg>
                        </div>
                        <p className="text-gray-500 mb-4">Giỏ hàng của bạn còn trống</p>
                        <Link href="/products" className="inline-block px-8 py-2 bg-[#ee4d2d] text-white hover:opacity-90">
                            Mua Sắm Ngay
                        </Link>
                    </div>
                ) : (
                    <>
                        {/* Cart Header */}
                        <div className="bg-white rounded-sm shadow-sm p-4 mb-3 hidden md:grid grid-cols-12 gap-4 text-sm text-gray-500">
                            <div className="col-span-5 flex items-center gap-3">
                                <input
                                    type="checkbox"
                                    checked={items.every(item => item.selected)}
                                    onChange={(e) => selectAll(e.target.checked)}
                                    className="w-4 h-4 accent-[#ee4d2d]"
                                />
                                <span>Sản Phẩm</span>
                            </div>
                            <div className="col-span-2 text-center">Đơn Giá</div>
                            <div className="col-span-2 text-center">Số Lượng</div>
                            <div className="col-span-2 text-center">Số Tiền</div>
                            <div className="col-span-1 text-center">Thao Tác</div>
                        </div>

                        {/* Cart Items by Shop */}
                        {Object.entries(groupedItems).map(([shop, shopItems]) => (
                            <div key={shop} className="bg-white rounded-sm shadow-sm mb-3">
                                {/* Shop Header */}
                                <div className="p-4 border-b flex items-center gap-3">
                                    <input
                                        type="checkbox"
                                        checked={shopItems.every(item => item.selected)}
                                        onChange={() => {
                                            const allSelected = shopItems.every(item => item.selected);
                                            setItems(prev => prev.map(item =>
                                                item.shop === shop ? { ...item, selected: !allSelected } : item
                                            ));
                                        }}
                                        className="w-4 h-4 accent-[#ee4d2d]"
                                    />
                                    <span className="bg-[#ee4d2d] text-white text-[10px] px-1">Mall</span>
                                    <span className="font-medium text-sm">{shop}</span>
                                </div>

                                {/* Items */}
                                {shopItems.map(item => (
                                    <div key={item.id} className="p-4 border-b last:border-0 grid grid-cols-12 gap-4 items-center">
                                        <div className="col-span-12 md:col-span-5 flex items-center gap-3">
                                            <input
                                                type="checkbox"
                                                checked={item.selected}
                                                onChange={() => toggleSelect(item.id)}
                                                className="w-4 h-4 accent-[#ee4d2d]"
                                            />
                                            <Link href={`/products/${item.id}`} className="w-20 h-20 bg-gray-100 rounded-sm flex items-center justify-center text-4xl flex-shrink-0 hover:opacity-80">
                                                {item.image}
                                            </Link>
                                            <div className="flex-1 min-w-0">
                                                <Link href={`/products/${item.id}`} className="text-sm line-clamp-2 hover:text-[#ee4d2d]">
                                                    {item.name}
                                                </Link>
                                                {item.variant && (
                                                    <p className="text-xs text-gray-400 mt-1">Phân loại: {item.variant}</p>
                                                )}
                                            </div>
                                        </div>

                                        <div className="col-span-4 md:col-span-2 text-center">
                                            <span className="text-gray-400 line-through text-xs block">₫{formatPrice(item.originalPrice)}</span>
                                            <span className="text-sm">₫{formatPrice(item.price)}</span>
                                        </div>

                                        <div className="col-span-4 md:col-span-2 flex justify-center">
                                            <div className="flex items-center">
                                                <button
                                                    onClick={() => updateQuantity(item.id, -1)}
                                                    className="w-7 h-7 border flex items-center justify-center text-lg hover:bg-gray-50"
                                                >
                                                    −
                                                </button>
                                                <span className="w-10 h-7 border-y flex items-center justify-center text-sm">{item.quantity}</span>
                                                <button
                                                    onClick={() => updateQuantity(item.id, 1)}
                                                    className="w-7 h-7 border flex items-center justify-center text-lg hover:bg-gray-50"
                                                >
                                                    +
                                                </button>
                                            </div>
                                        </div>

                                        <div className="col-span-2 md:col-span-2 text-center">
                                            <span className="text-[#ee4d2d] text-sm">₫{formatPrice(item.price * item.quantity)}</span>
                                        </div>

                                        <div className="col-span-2 md:col-span-1 text-center">
                                            <button
                                                onClick={() => removeItem(item.id)}
                                                className="text-gray-500 hover:text-[#ee4d2d] text-sm"
                                            >
                                                Xóa
                                            </button>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        ))}

                        {/* Voucher */}
                        <div className="bg-white rounded-sm shadow-sm p-4 mb-3">
                            <div className="flex items-center gap-4">
                                <span className="text-[#ee4d2d]">🎟️</span>
                                <span className="text-sm">Shopee Voucher</span>
                                <div className="flex-1 flex gap-2">
                                    <input
                                        type="text"
                                        value={voucher}
                                        onChange={(e) => setVoucher(e.target.value)}
                                        placeholder="Nhập mã voucher"
                                        className="border px-3 py-1.5 text-sm outline-none focus:border-[#ee4d2d] flex-1 max-w-xs"
                                        disabled={voucherApplied}
                                    />
                                    <button
                                        onClick={applyVoucher}
                                        disabled={voucherApplied}
                                        className={`px-4 py-1.5 text-sm transition-colors ${voucherApplied
                                                ? 'bg-gray-100 text-gray-400'
                                                : 'bg-[#ee4d2d] text-white hover:opacity-90'
                                            }`}
                                    >
                                        {voucherApplied ? 'Đã áp dụng' : 'Áp dụng'}
                                    </button>
                                </div>
                            </div>
                        </div>

                        {/* Footer */}
                        <div className="bg-white rounded-sm shadow-sm p-4 sticky bottom-0">
                            <div className="flex items-center justify-between flex-wrap gap-4">
                                <div className="flex items-center gap-4">
                                    <label className="flex items-center gap-2 cursor-pointer">
                                        <input
                                            type="checkbox"
                                            checked={items.every(item => item.selected)}
                                            onChange={(e) => selectAll(e.target.checked)}
                                            className="w-4 h-4 accent-[#ee4d2d]"
                                        />
                                        <span className="text-sm">Chọn Tất Cả ({items.length})</span>
                                    </label>
                                    <button className="text-sm text-gray-500 hover:text-[#ee4d2d]">Xóa</button>
                                </div>
                                <div className="flex items-center gap-4">
                                    <div className="text-right">
                                        <span className="text-sm text-gray-500">Tổng thanh toán ({selectedItems.length} Sản phẩm): </span>
                                        <span className="text-2xl text-[#ee4d2d] font-medium">₫{formatPrice(total)}</span>
                                        {discount > 0 && (
                                            <div className="text-xs text-gray-400">Tiết kiệm: ₫{formatPrice(discount)}</div>
                                        )}
                                    </div>
                                    <button
                                        onClick={proceedToCheckout}
                                        disabled={selectedItems.length === 0}
                                        className={`px-12 py-3 text-sm transition-colors ${selectedItems.length > 0
                                                ? 'bg-[#ee4d2d] text-white hover:opacity-90'
                                                : 'bg-gray-200 text-gray-400 cursor-not-allowed'
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
