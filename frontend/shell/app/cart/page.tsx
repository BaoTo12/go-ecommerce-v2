'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';

interface CartItem {
    id: string;
    name: string;
    price: number;
    originalPrice: number;
    quantity: number;
    image: string;
    selected: boolean;
}

export default function CartPage() {
    const [items, setItems] = useState<CartItem[]>([]);
    const [voucher, setVoucher] = useState('');
    const [voucherApplied, setVoucherApplied] = useState(false);
    const [discount, setDiscount] = useState(0);

    useEffect(() => {
        // Load demo cart items
        setItems([
            { id: 'p1', name: 'iPhone 15 Pro Max 256GB', price: 29990000, originalPrice: 34990000, quantity: 1, image: '📱', selected: true },
            { id: 'p6', name: 'Son Dưỡng Môi', price: 150000, originalPrice: 250000, quantity: 2, image: '💄', selected: true },
            { id: 'p4', name: 'Áo Hoodie Unisex', price: 299000, originalPrice: 450000, quantity: 1, image: '👕', selected: false },
        ]);
    }, []);

    const updateQuantity = (id: string, delta: number) => {
        setItems(prev => prev.map(item => {
            if (item.id === id) {
                const newQty = Math.max(1, item.quantity + delta);
                return { ...item, quantity: newQty };
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
        if (voucher.toUpperCase() === 'GIẢM50K') {
            setDiscount(50000);
            setVoucherApplied(true);
        } else if (voucher.toUpperCase() === 'SALE10') {
            setDiscount(subtotal * 0.1);
            setVoucherApplied(true);
        } else {
            alert('Mã giảm giá không hợp lệ!');
        }
    };

    const selectedItems = items.filter(item => item.selected);
    const subtotal = selectedItems.reduce((sum, item) => sum + item.price * item.quantity, 0);
    const totalSavings = selectedItems.reduce((sum, item) => sum + (item.originalPrice - item.price) * item.quantity, 0);
    const total = subtotal - discount;

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    return (
        <div className="min-h-screen bg-[#F5F5F5]">
            {/* Header */}
            <div className="bg-white border-b">
                <div className="container mx-auto px-4 py-4">
                    <div className="flex items-center gap-4">
                        <h1 className="text-xl font-bold text-[#EE4D2D]">🛒 Giỏ Hàng</h1>
                        <span className="text-gray-500">({items.length} sản phẩm)</span>
                    </div>
                </div>
            </div>

            <div className="container mx-auto px-4 py-6">
                {items.length === 0 ? (
                    <div className="bg-white rounded p-12 text-center">
                        <div className="text-8xl mb-4">🛒</div>
                        <p className="text-gray-500 mb-4">Giỏ hàng của bạn đang trống</p>
                        <Link href="/products" className="bg-[#EE4D2D] text-white px-8 py-3 rounded hover:bg-[#D73211]">
                            Mua sắm ngay
                        </Link>
                    </div>
                ) : (
                    <div className="grid lg:grid-cols-3 gap-6">
                        {/* Cart Items */}
                        <div className="lg:col-span-2 space-y-3">
                            {/* Select All */}
                            <div className="bg-white rounded p-4 flex items-center gap-4">
                                <input
                                    type="checkbox"
                                    checked={items.every(item => item.selected)}
                                    onChange={(e) => selectAll(e.target.checked)}
                                    className="w-5 h-5 accent-[#EE4D2D]"
                                />
                                <span className="font-semibold">Chọn tất cả ({items.length})</span>
                            </div>

                            {/* Items */}
                            {items.map(item => (
                                <div key={item.id} className="bg-white rounded p-4">
                                    <div className="flex items-start gap-4">
                                        <input
                                            type="checkbox"
                                            checked={item.selected}
                                            onChange={() => toggleSelect(item.id)}
                                            className="w-5 h-5 accent-[#EE4D2D] mt-2"
                                        />

                                        {/* Image */}
                                        <div className="w-20 h-20 bg-gray-100 rounded flex items-center justify-center text-4xl flex-shrink-0">
                                            {item.image}
                                        </div>

                                        {/* Info */}
                                        <div className="flex-1">
                                            <h3 className="font-semibold text-sm mb-1">{item.name}</h3>
                                            <div className="flex items-baseline gap-2 mb-3">
                                                <span className="text-[#EE4D2D] font-bold">₫{formatPrice(item.price)}</span>
                                                <span className="text-gray-400 text-xs line-through">₫{formatPrice(item.originalPrice)}</span>
                                            </div>

                                            {/* Quantity */}
                                            <div className="flex items-center gap-3">
                                                <div className="flex items-center border rounded">
                                                    <button
                                                        onClick={() => updateQuantity(item.id, -1)}
                                                        className="px-3 py-1 hover:bg-gray-100"
                                                        disabled={item.quantity <= 1}
                                                    >
                                                        −
                                                    </button>
                                                    <span className="px-4 py-1 border-x">{item.quantity}</span>
                                                    <button
                                                        onClick={() => updateQuantity(item.id, 1)}
                                                        className="px-3 py-1 hover:bg-gray-100"
                                                    >
                                                        +
                                                    </button>
                                                </div>
                                                <span className="text-sm text-gray-500">
                                                    Thành tiền: <span className="text-[#EE4D2D] font-semibold">₫{formatPrice(item.price * item.quantity)}</span>
                                                </span>
                                            </div>
                                        </div>

                                        {/* Remove */}
                                        <button
                                            onClick={() => removeItem(item.id)}
                                            className="text-gray-400 hover:text-red-500 text-xl"
                                        >
                                            🗑️
                                        </button>
                                    </div>
                                </div>
                            ))}
                        </div>

                        {/* Summary */}
                        <div className="space-y-4">
                            {/* Voucher */}
                            <div className="bg-white rounded p-4">
                                <h3 className="font-semibold mb-3">🎟️ Mã giảm giá</h3>
                                <div className="flex gap-2">
                                    <input
                                        type="text"
                                        value={voucher}
                                        onChange={(e) => setVoucher(e.target.value)}
                                        placeholder="Nhập mã..."
                                        className="flex-1 border rounded px-3 py-2 text-sm focus:outline-none focus:border-[#EE4D2D]"
                                        disabled={voucherApplied}
                                    />
                                    <button
                                        onClick={applyVoucher}
                                        disabled={voucherApplied}
                                        className={`px-4 py-2 rounded text-sm font-semibold ${voucherApplied
                                                ? 'bg-green-100 text-green-600'
                                                : 'bg-[#EE4D2D] text-white hover:bg-[#D73211]'
                                            }`}
                                    >
                                        {voucherApplied ? '✓ Đã áp dụng' : 'Áp dụng'}
                                    </button>
                                </div>
                                {voucherApplied && (
                                    <p className="text-green-600 text-sm mt-2">Giảm ₫{formatPrice(discount)}</p>
                                )}
                            </div>

                            {/* Total */}
                            <div className="bg-white rounded p-4">
                                <h3 className="font-semibold mb-3">📦 Tổng đơn hàng</h3>
                                <div className="space-y-2 text-sm">
                                    <div className="flex justify-between">
                                        <span className="text-gray-500">Tạm tính ({selectedItems.length} sản phẩm)</span>
                                        <span>₫{formatPrice(subtotal)}</span>
                                    </div>
                                    <div className="flex justify-between text-green-600">
                                        <span>Tiết kiệm</span>
                                        <span>-₫{formatPrice(totalSavings)}</span>
                                    </div>
                                    {discount > 0 && (
                                        <div className="flex justify-between text-green-600">
                                            <span>Voucher</span>
                                            <span>-₫{formatPrice(discount)}</span>
                                        </div>
                                    )}
                                    <div className="flex justify-between">
                                        <span className="text-gray-500">Vận chuyển</span>
                                        <span className="text-green-600">Miễn phí</span>
                                    </div>
                                    <div className="pt-3 mt-3 border-t flex justify-between font-bold text-lg">
                                        <span>Tổng cộng</span>
                                        <span className="text-[#EE4D2D]">₫{formatPrice(total)}</span>
                                    </div>
                                </div>

                                <button
                                    disabled={selectedItems.length === 0}
                                    className={`w-full mt-4 py-3 rounded font-bold ${selectedItems.length > 0
                                            ? 'bg-[#EE4D2D] text-white hover:bg-[#D73211]'
                                            : 'bg-gray-200 text-gray-400 cursor-not-allowed'
                                        }`}
                                >
                                    Đặt hàng ({selectedItems.length})
                                </button>
                            </div>

                            {/* Promo */}
                            <div className="bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] rounded p-4 text-white text-center">
                                <p className="text-sm">Dùng mã <span className="font-mono font-bold">GIẢM50K</span> giảm 50K!</p>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
