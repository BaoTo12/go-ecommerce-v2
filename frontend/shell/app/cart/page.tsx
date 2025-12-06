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
            { id: 'p1', name: 'iPhone 15 Pro Max 256GB', price: 29990000, originalPrice: 34990000, quantity: 1, image: '📱', selected: true },
            { id: 'p6', name: 'Son Dưỡng Môi Dior', price: 950000, originalPrice: 1200000, quantity: 2, image: '💄', selected: true },
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

    const proceedToCheckout = () => {
        if (selectedItems.length > 0) {
            router.push('/checkout');
        }
    };

    return (
        <div className="min-h-screen bg-[#F5F5F5] animate-fade-in">
            {/* Header */}
            <div className="bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] py-6">
                <div className="container mx-auto px-4">
                    <h1 className="text-2xl font-bold text-white flex items-center gap-2">
                        <span>🛒</span> Giỏ Hàng
                        <span className="text-white/70 text-lg font-normal">({items.length} sản phẩm)</span>
                    </h1>
                </div>
            </div>

            <div className="container mx-auto px-4 py-6">
                {items.length === 0 ? (
                    <div className="bg-white rounded-2xl p-12 text-center">
                        <div className="text-8xl mb-4">🛒</div>
                        <p className="text-gray-500 mb-4">Giỏ hàng của bạn đang trống</p>
                        <Link href="/products" className="inline-block bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] text-white px-8 py-3 rounded-full font-bold hover:opacity-90">
                            Mua sắm ngay
                        </Link>
                    </div>
                ) : (
                    <div className="grid lg:grid-cols-3 gap-6">
                        {/* Cart Items */}
                        <div className="lg:col-span-2 space-y-4">
                            {/* Select All */}
                            <div className="bg-white rounded-2xl p-4 flex items-center gap-4 shadow-sm">
                                <input
                                    type="checkbox"
                                    checked={items.every(item => item.selected)}
                                    onChange={(e) => selectAll(e.target.checked)}
                                    className="w-5 h-5 accent-[#EE4D2D] rounded"
                                />
                                <span className="font-semibold">Chọn tất cả ({items.length})</span>
                            </div>

                            {/* Items */}
                            {items.map((item, index) => (
                                <div
                                    key={item.id}
                                    className="bg-white rounded-2xl p-4 shadow-sm animate-slide-up"
                                    style={{ animationDelay: `${index * 50}ms` }}
                                >
                                    <div className="flex items-start gap-4">
                                        <input
                                            type="checkbox"
                                            checked={item.selected}
                                            onChange={() => toggleSelect(item.id)}
                                            className="w-5 h-5 accent-[#EE4D2D] mt-4 rounded"
                                        />

                                        {/* Image */}
                                        <Link href={`/products/${item.id}`} className="w-24 h-24 bg-gray-100 rounded-xl flex items-center justify-center text-5xl flex-shrink-0 hover:scale-105 transition-transform">
                                            {item.image}
                                        </Link>

                                        {/* Info */}
                                        <div className="flex-1">
                                            <Link href={`/products/${item.id}`}>
                                                <h3 className="font-semibold hover:text-[#EE4D2D] transition-colors">{item.name}</h3>
                                            </Link>
                                            <div className="flex items-baseline gap-2 my-2">
                                                <span className="text-[#EE4D2D] font-bold text-lg">₫{formatPrice(item.price)}</span>
                                                <span className="text-gray-400 text-sm line-through">₫{formatPrice(item.originalPrice)}</span>
                                            </div>

                                            {/* Quantity */}
                                            <div className="flex items-center gap-3">
                                                <div className="flex items-center border-2 rounded-full overflow-hidden">
                                                    <button
                                                        onClick={() => updateQuantity(item.id, -1)}
                                                        className="px-4 py-2 hover:bg-gray-100 transition-colors"
                                                        disabled={item.quantity <= 1}
                                                    >
                                                        −
                                                    </button>
                                                    <span className="px-4 py-2 border-x-2 font-medium">{item.quantity}</span>
                                                    <button
                                                        onClick={() => updateQuantity(item.id, 1)}
                                                        className="px-4 py-2 hover:bg-gray-100 transition-colors"
                                                    >
                                                        +
                                                    </button>
                                                </div>
                                                <span className="text-sm text-gray-500">
                                                    Thành tiền: <span className="text-[#EE4D2D] font-bold">₫{formatPrice(item.price * item.quantity)}</span>
                                                </span>
                                            </div>
                                        </div>

                                        {/* Remove */}
                                        <button
                                            onClick={() => removeItem(item.id)}
                                            className="text-gray-400 hover:text-red-500 text-xl p-2 hover:bg-red-50 rounded-full transition-all"
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
                            <div className="bg-white rounded-2xl p-5 shadow-sm">
                                <h3 className="font-bold mb-3 flex items-center gap-2">🎟️ Mã giảm giá</h3>
                                <div className="flex gap-2">
                                    <input
                                        type="text"
                                        value={voucher}
                                        onChange={(e) => setVoucher(e.target.value)}
                                        placeholder="Nhập mã..."
                                        className="flex-1 border-2 rounded-full px-4 py-2 text-sm focus:outline-none focus:border-[#EE4D2D]"
                                        disabled={voucherApplied}
                                    />
                                    <button
                                        onClick={applyVoucher}
                                        disabled={voucherApplied}
                                        className={`px-5 py-2 rounded-full text-sm font-bold transition-all ${voucherApplied
                                                ? 'bg-green-100 text-green-600'
                                                : 'bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] text-white hover:opacity-90'
                                            }`}
                                    >
                                        {voucherApplied ? '✓' : 'Áp dụng'}
                                    </button>
                                </div>
                                {voucherApplied && (
                                    <p className="text-green-600 text-sm mt-2">✓ Giảm ₫{formatPrice(discount)}</p>
                                )}
                            </div>

                            {/* Total */}
                            <div className="bg-white rounded-2xl p-5 shadow-sm">
                                <h3 className="font-bold mb-4 flex items-center gap-2">📦 Tổng đơn hàng</h3>
                                <div className="space-y-3 text-sm">
                                    <div className="flex justify-between">
                                        <span className="text-gray-500">Tạm tính ({selectedItems.length} sản phẩm)</span>
                                        <span className="font-medium">₫{formatPrice(subtotal)}</span>
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
                                        <span className="text-green-600 font-medium">Miễn phí</span>
                                    </div>
                                    <div className="pt-4 mt-4 border-t-2 flex justify-between font-bold text-xl">
                                        <span>Tổng cộng</span>
                                        <span className="text-[#EE4D2D]">₫{formatPrice(total)}</span>
                                    </div>
                                </div>

                                <button
                                    onClick={proceedToCheckout}
                                    disabled={selectedItems.length === 0}
                                    className={`w-full mt-6 py-4 rounded-full font-bold text-lg transition-all ${selectedItems.length > 0
                                            ? 'bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] text-white hover:opacity-90 hover:shadow-lg'
                                            : 'bg-gray-200 text-gray-400 cursor-not-allowed'
                                        }`}
                                >
                                    Thanh toán ({selectedItems.length})
                                </button>
                            </div>

                            {/* Promo */}
                            <div className="bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] rounded-2xl p-4 text-white text-center">
                                <p className="text-sm">Dùng mã <span className="font-mono font-bold bg-white/20 px-2 py-1 rounded">GIẢM50K</span> giảm 50K!</p>
                            </div>

                            {/* Trust badges */}
                            <div className="bg-white rounded-2xl p-4 shadow-sm">
                                <div className="flex items-center justify-center gap-6 text-sm text-gray-500">
                                    <span className="flex items-center gap-1">🔒 An toàn</span>
                                    <span className="flex items-center gap-1">🚚 Freeship</span>
                                    <span className="flex items-center gap-1">✅ Chính hãng</span>
                                </div>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
