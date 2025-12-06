'use client';

import React, { useState } from 'react';
import Link from 'next/link';

interface CartItem {
    id: string;
    name: string;
    price: number;
    quantity: number;
    image: string;
    variant?: string;
}

export default function CheckoutPage() {
    const [selectedPayment, setSelectedPayment] = useState('cod');
    const [selectedShipping, setSelectedShipping] = useState('standard');
    const [orderPlaced, setOrderPlaced] = useState(false);
    const [isProcessing, setIsProcessing] = useState(false);
    const [note, setNote] = useState('');

    const address = {
        name: 'Nguyễn Văn A',
        phone: '(+84) 901 234 567',
        address: '123 Đường ABC, Phường XYZ, Quận 1, Thành phố Hồ Chí Minh',
        isDefault: true,
    };

    const cartItems: CartItem[] = [
        { id: 'p1', name: 'iPhone 15 Pro Max 256GB Titan Xanh Chính Hãng VN/A', price: 29990000, quantity: 1, image: '📱', variant: 'Xanh Titan, 256GB' },
        { id: 'p6', name: 'Son Dưỡng Môi Dior Addict Lip Glow', price: 950000, quantity: 2, image: '💄', variant: 'Màu 001 Pink' },
    ];

    const paymentMethods = [
        { id: 'cod', name: 'Thanh toán khi nhận hàng', icon: '💵' },
        { id: 'shopee_pay', name: 'Ví ShopeePay', icon: '🟠', desc: 'Ví điện tử' },
        { id: 'momo', name: 'Ví MoMo', icon: '🟣' },
        { id: 'vnpay', name: 'VNPay QR', icon: '🔵' },
        { id: 'zalopay', name: 'ZaloPay', icon: '🔷' },
        { id: 'card', name: 'Thẻ tín dụng/ghi nợ', icon: '💳', desc: 'Visa, Mastercard, JCB' },
    ];

    const shippingMethods = [
        { id: 'standard', name: 'Giao Hàng Tiết Kiệm', time: '4-6 ngày', price: 0 },
        { id: 'fast', name: 'Giao Hàng Nhanh', time: '2-3 ngày', price: 25000 },
        { id: 'express', name: 'Hỏa Tốc', time: 'Trong ngày', price: 50000 },
    ];

    const subtotal = cartItems.reduce((sum, item) => sum + item.price * item.quantity, 0);
    const shippingFee = shippingMethods.find(s => s.id === selectedShipping)?.price || 0;
    const discount = 50000;
    const total = subtotal + shippingFee - discount;

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    const placeOrder = () => {
        setIsProcessing(true);
        setTimeout(() => {
            setIsProcessing(false);
            setOrderPlaced(true);
        }, 2000);
    };

    if (orderPlaced) {
        return (
            <div className="min-h-screen bg-[#f5f5f5] flex items-center justify-center">
                <div className="bg-white rounded-sm p-8 text-center max-w-md mx-4 shadow-sm">
                    <div className="w-16 h-16 bg-[#00bfa5] rounded-full flex items-center justify-center mx-auto mb-4">
                        <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                        </svg>
                    </div>
                    <h1 className="text-xl font-semibold text-gray-800 mb-2">Đặt hàng thành công!</h1>
                    <p className="text-gray-500 text-sm mb-4">Cảm ơn bạn đã mua hàng tại Shopee</p>
                    <div className="bg-[#fef6f5] border border-[#ee4d2d] rounded-sm p-3 mb-4">
                        <p className="text-xs text-gray-500">Mã đơn hàng</p>
                        <p className="text-[#ee4d2d] font-bold">#SP{Date.now().toString().slice(-10)}</p>
                    </div>
                    <div className="flex gap-2">
                        <Link href="/products" className="flex-1 py-2 border border-[#ee4d2d] text-[#ee4d2d] text-sm hover:bg-[#fef6f5] transition-colors">
                            Tiếp tục mua sắm
                        </Link>
                        <Link href="/" className="flex-1 py-2 bg-[#ee4d2d] text-white text-sm hover:opacity-90">
                            Về trang chủ
                        </Link>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-[#f5f5f5]">
            {/* Header */}
            <div className="bg-white border-b">
                <div className="container mx-auto px-4 py-4">
                    <div className="flex items-center gap-3">
                        <Link href="/" className="text-2xl font-bold text-[#ee4d2d]">Shopee</Link>
                        <span className="text-gray-300">|</span>
                        <h1 className="text-xl text-gray-700">Thanh Toán</h1>
                    </div>
                </div>
            </div>

            <div className="container mx-auto px-4 py-6">
                {/* Address */}
                <div className="bg-white rounded-sm shadow-sm mb-3">
                    <div className="p-4 border-b border-dashed">
                        <div className="flex items-center gap-2 text-[#ee4d2d] text-sm font-medium mb-2">
                            <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                                <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z" />
                            </svg>
                            Địa Chỉ Nhận Hàng
                        </div>
                        <div className="flex items-start justify-between">
                            <div>
                                <div className="flex items-center gap-2 mb-1">
                                    <span className="font-semibold">{address.name}</span>
                                    <span className="text-gray-400">|</span>
                                    <span className="text-gray-500">{address.phone}</span>
                                </div>
                                <p className="text-gray-600 text-sm">{address.address}</p>
                                {address.isDefault && (
                                    <span className="inline-block mt-1 text-xs text-[#ee4d2d] border border-[#ee4d2d] px-1">Mặc Định</span>
                                )}
                            </div>
                            <button className="text-[#4080ee] text-sm hover:underline">Thay Đổi</button>
                        </div>
                    </div>
                </div>

                {/* Products */}
                <div className="bg-white rounded-sm shadow-sm mb-3">
                    <div className="p-4 border-b">
                        <span className="text-sm">Sản phẩm</span>
                    </div>
                    {cartItems.map(item => (
                        <div key={item.id} className="p-4 border-b flex items-center gap-4">
                            <div className="w-16 h-16 bg-gray-100 rounded-sm flex items-center justify-center text-3xl flex-shrink-0">
                                {item.image}
                            </div>
                            <div className="flex-1">
                                <h3 className="text-sm line-clamp-2">{item.name}</h3>
                                {item.variant && <p className="text-xs text-gray-400 mt-1">Phân loại: {item.variant}</p>}
                            </div>
                            <div className="text-right">
                                <p className="text-sm">₫{formatPrice(item.price)}</p>
                                <p className="text-sm text-gray-400">x{item.quantity}</p>
                            </div>
                            <div className="text-right w-24">
                                <p className="text-sm text-[#ee4d2d]">₫{formatPrice(item.price * item.quantity)}</p>
                            </div>
                        </div>
                    ))}

                    {/* Note */}
                    <div className="p-4 flex items-center gap-4 border-b">
                        <span className="text-sm text-gray-500">Lời nhắn:</span>
                        <input
                            type="text"
                            value={note}
                            onChange={(e) => setNote(e.target.value)}
                            placeholder="Lưu ý cho Người bán..."
                            className="flex-1 border px-3 py-1.5 text-sm outline-none focus:border-[#ee4d2d]"
                        />
                    </div>

                    {/* Shipping */}
                    <div className="p-4">
                        <div className="flex items-center justify-between mb-3">
                            <span className="text-sm text-[#00bfa5] font-medium">Đơn vị vận chuyển:</span>
                        </div>
                        <div className="space-y-2">
                            {shippingMethods.map(method => (
                                <label
                                    key={method.id}
                                    className={`flex items-center gap-3 p-3 border rounded-sm cursor-pointer transition-colors ${selectedShipping === method.id ? 'border-[#ee4d2d] bg-[#fef6f5]' : 'hover:border-gray-400'
                                        }`}
                                >
                                    <input
                                        type="radio"
                                        name="shipping"
                                        value={method.id}
                                        checked={selectedShipping === method.id}
                                        onChange={(e) => setSelectedShipping(e.target.value)}
                                        className="accent-[#ee4d2d]"
                                    />
                                    <div className="flex-1">
                                        <span className="text-sm font-medium">{method.name}</span>
                                        <span className="text-xs text-gray-400 ml-2">({method.time})</span>
                                    </div>
                                    <span className={`text-sm ${method.price === 0 ? 'text-[#00bfa5]' : ''}`}>
                                        {method.price === 0 ? 'Miễn phí' : `₫${formatPrice(method.price)}`}
                                    </span>
                                </label>
                            ))}
                        </div>
                    </div>
                </div>

                {/* Payment Methods */}
                <div className="bg-white rounded-sm shadow-sm mb-3">
                    <div className="p-4 border-b">
                        <span className="text-sm">Phương thức thanh toán</span>
                    </div>
                    <div className="p-4 grid grid-cols-2 md:grid-cols-3 gap-2">
                        {paymentMethods.map(method => (
                            <button
                                key={method.id}
                                onClick={() => setSelectedPayment(method.id)}
                                className={`p-3 border rounded-sm text-left transition-colors ${selectedPayment === method.id
                                        ? 'border-[#ee4d2d] bg-[#fef6f5]'
                                        : 'hover:border-gray-400'
                                    }`}
                            >
                                <div className="flex items-center gap-2">
                                    <span className="text-xl">{method.icon}</span>
                                    <div>
                                        <p className="text-sm font-medium">{method.name}</p>
                                        {method.desc && <p className="text-[10px] text-gray-400">{method.desc}</p>}
                                    </div>
                                </div>
                            </button>
                        ))}
                    </div>
                </div>

                {/* Order Summary */}
                <div className="bg-white rounded-sm shadow-sm">
                    <div className="p-4 border-b">
                        <div className="flex justify-between text-sm mb-2">
                            <span className="text-gray-500">Tổng tiền hàng</span>
                            <span>₫{formatPrice(subtotal)}</span>
                        </div>
                        <div className="flex justify-between text-sm mb-2">
                            <span className="text-gray-500">Phí vận chuyển</span>
                            <span className={shippingFee === 0 ? 'text-[#00bfa5]' : ''}>
                                {shippingFee === 0 ? 'Miễn phí' : `₫${formatPrice(shippingFee)}`}
                            </span>
                        </div>
                        <div className="flex justify-between text-sm mb-2">
                            <span className="text-gray-500">Voucher Shopee</span>
                            <span className="text-[#ee4d2d]">-₫{formatPrice(discount)}</span>
                        </div>
                    </div>
                    <div className="p-4 flex items-center justify-between">
                        <div>
                            <span className="text-gray-500 text-sm">Tổng thanh toán:</span>
                            <span className="text-2xl text-[#ee4d2d] font-medium ml-2">₫{formatPrice(total)}</span>
                        </div>
                        <button
                            onClick={placeOrder}
                            disabled={isProcessing}
                            className={`px-12 py-3 bg-[#ee4d2d] text-white font-medium hover:opacity-90 transition-opacity ${isProcessing ? 'opacity-70 cursor-wait' : ''
                                }`}
                        >
                            {isProcessing ? 'Đang xử lý...' : 'Đặt hàng'}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
