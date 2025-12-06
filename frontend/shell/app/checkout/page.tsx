'use client';

import React, { useState, useEffect } from 'react';
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
    const [isLoaded, setIsLoaded] = useState(false);
    const [showConfetti, setShowConfetti] = useState(false);

    useEffect(() => {
        setIsLoaded(true);
    }, []);

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
        { id: 'cod', name: 'Thanh toán khi nhận hàng', icon: '💵', color: 'from-green-400 to-green-600' },
        { id: 'shopee_pay', name: 'Ví ShopeePay', icon: '🟠', desc: 'Ví điện tử', color: 'from-orange-400 to-orange-600' },
        { id: 'momo', name: 'Ví MoMo', icon: '🟣', color: 'from-pink-400 to-pink-600' },
        { id: 'vnpay', name: 'VNPay QR', icon: '🔵', color: 'from-blue-400 to-blue-600' },
        { id: 'zalopay', name: 'ZaloPay', icon: '🔷', color: 'from-cyan-400 to-cyan-600' },
        { id: 'card', name: 'Thẻ tín dụng/ghi nợ', icon: '💳', desc: 'Visa, Mastercard, JCB', color: 'from-purple-400 to-purple-600' },
    ];

    const shippingMethods = [
        { id: 'standard', name: 'Giao Hàng Tiết Kiệm', time: '4-6 ngày', price: 0, icon: '📦' },
        { id: 'fast', name: 'Giao Hàng Nhanh', time: '2-3 ngày', price: 25000, icon: '🚀' },
        { id: 'express', name: 'Hỏa Tốc', time: 'Trong ngày', price: 50000, icon: '⚡' },
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
            setShowConfetti(true);
            setOrderPlaced(true);
        }, 2000);
    };

    if (orderPlaced) {
        return (
            <div className="min-h-screen bg-[#f5f5f5] flex items-center justify-center relative overflow-hidden">
                {/* Confetti */}
                {showConfetti && (
                    <div className="fixed inset-0 pointer-events-none z-40">
                        {[...Array(50)].map((_, i) => (
                            <div
                                key={i}
                                className="confetti"
                                style={{
                                    left: `${Math.random() * 100}%`,
                                    backgroundColor: ['#ee4d2d', '#00bfa5', '#ffc107', '#5c6bc0', '#ff6533'][Math.floor(Math.random() * 5)],
                                    animationDelay: `${Math.random() * 2}s`,
                                }}
                            />
                        ))}
                    </div>
                )}

                <div className="bg-white rounded-sm p-8 text-center max-w-md mx-4 shadow-lg animate-bounce-in relative z-50">
                    <div className="w-20 h-20 bg-gradient-to-br from-[#00bfa5] to-[#00897b] rounded-full flex items-center justify-center mx-auto mb-4 animate-pulse-glow">
                        <svg className="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                        </svg>
                    </div>
                    <h1 className="text-2xl font-bold text-gray-800 mb-2 animate-fade-in-up">Đặt hàng thành công!</h1>
                    <p className="text-gray-500 text-sm mb-4 animate-fade-in-up" style={{ animationDelay: '100ms' }}>Cảm ơn bạn đã mua hàng tại Shopee</p>
                    <div className="bg-gradient-to-r from-[#fef6f5] to-[#fff5f5] border border-[#ee4d2d] rounded-sm p-4 mb-4 animate-fade-in-up" style={{ animationDelay: '200ms' }}>
                        <p className="text-xs text-gray-500">Mã đơn hàng</p>
                        <p className="text-[#ee4d2d] font-bold text-xl">#SP{Date.now().toString().slice(-10)}</p>
                    </div>
                    <p className="text-sm text-gray-600 mb-4 animate-fade-in-up" style={{ animationDelay: '300ms' }}>
                        📧 Thông tin đơn hàng đã được gửi đến email của bạn
                    </p>
                    <div className="flex gap-2 animate-fade-in-up" style={{ animationDelay: '400ms' }}>
                        <Link href="/products" className="flex-1 py-2.5 border border-[#ee4d2d] text-[#ee4d2d] text-sm hover:bg-[#fef6f5] transition-all hover-shrink">
                            Tiếp tục mua sắm
                        </Link>
                        <Link href="/" className="flex-1 py-2.5 bg-gradient-to-r from-[#ee4d2d] to-[#ff6533] text-white text-sm hover:opacity-90 transition-all hover-shrink">
                            Về trang chủ
                        </Link>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className={`min-h-screen bg-[#f5f5f5] ${isLoaded ? 'animate-fade-in' : 'opacity-0'}`}>
            {/* Header */}
            <div className="bg-white border-b animate-fade-in-down">
                <div className="container mx-auto px-4 py-4">
                    <div className="flex items-center gap-3">
                        <Link href="/" className="text-2xl font-bold text-[#ee4d2d] hover:scale-105 transition-transform">Shopee</Link>
                        <span className="text-gray-300">|</span>
                        <h1 className="text-xl text-gray-700">Thanh Toán</h1>
                        <div className="ml-auto flex items-center gap-2 text-sm text-gray-500">
                            <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
                            Giao dịch bảo mật
                        </div>
                    </div>
                </div>
            </div>

            <div className="container mx-auto px-4 py-6">
                {/* Address */}
                <div className="bg-white rounded-sm shadow-sm mb-3 animate-fade-in-up overflow-hidden">
                    <div className="h-1 bg-gradient-to-r from-[#ee4d2d] via-[#ff6533] to-[#ee4d2d] animate-gradient" />
                    <div className="p-4">
                        <div className="flex items-center gap-2 text-[#ee4d2d] text-sm font-medium mb-2">
                            <svg className="w-4 h-4 animate-float" fill="currentColor" viewBox="0 0 24 24">
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
                                    <span className="inline-block mt-1 text-xs text-[#ee4d2d] border border-[#ee4d2d] px-1.5 py-0.5 animate-pulse">Mặc Định</span>
                                )}
                            </div>
                            <button className="text-[#4080ee] text-sm hover:underline transition-all hover-shrink">Thay Đổi</button>
                        </div>
                    </div>
                </div>

                {/* Products */}
                <div className="bg-white rounded-sm shadow-sm mb-3 animate-fade-in-up" style={{ animationDelay: '100ms' }}>
                    <div className="p-4 border-b flex items-center justify-between">
                        <span className="text-sm font-medium">Sản phẩm đặt mua</span>
                        <span className="text-sm text-gray-500">{cartItems.length} sản phẩm</span>
                    </div>
                    {cartItems.map((item, index) => (
                        <div
                            key={item.id}
                            className="p-4 border-b flex items-center gap-4 hover:bg-gray-50 transition-colors animate-fade-in-left"
                            style={{ animationDelay: `${(index + 1) * 100}ms` }}
                        >
                            <div className="w-16 h-16 bg-gray-100 rounded-sm flex items-center justify-center text-3xl flex-shrink-0 hover:scale-110 transition-transform">
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
                            <div className="text-right w-28">
                                <p className="text-sm text-[#ee4d2d] font-medium">₫{formatPrice(item.price * item.quantity)}</p>
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
                            className="flex-1 border px-3 py-2 text-sm outline-none focus:border-[#ee4d2d] rounded-sm transition-all"
                        />
                    </div>

                    {/* Shipping */}
                    <div className="p-4">
                        <div className="flex items-center gap-2 mb-3">
                            <span className="text-sm text-[#00bfa5] font-medium">🚚 Đơn vị vận chuyển:</span>
                        </div>
                        <div className="space-y-2">
                            {shippingMethods.map((method, index) => (
                                <label
                                    key={method.id}
                                    className={`flex items-center gap-3 p-3 border rounded-sm cursor-pointer transition-all hover-shrink animate-fade-in-right ${selectedShipping === method.id ? 'border-[#ee4d2d] bg-[#fef6f5]' : 'hover:border-[#ee4d2d]'
                                        }`}
                                    style={{ animationDelay: `${index * 100}ms` }}
                                >
                                    <input
                                        type="radio"
                                        name="shipping"
                                        value={method.id}
                                        checked={selectedShipping === method.id}
                                        onChange={(e) => setSelectedShipping(e.target.value)}
                                        className="accent-[#ee4d2d]"
                                    />
                                    <span className="text-xl">{method.icon}</span>
                                    <div className="flex-1">
                                        <span className="text-sm font-medium">{method.name}</span>
                                        <span className="text-xs text-gray-400 ml-2">({method.time})</span>
                                    </div>
                                    <span className={`text-sm font-medium ${method.price === 0 ? 'text-[#00bfa5]' : 'text-gray-700'}`}>
                                        {method.price === 0 ? 'Miễn phí' : `₫${formatPrice(method.price)}`}
                                    </span>
                                </label>
                            ))}
                        </div>
                    </div>
                </div>

                {/* Payment Methods */}
                <div className="bg-white rounded-sm shadow-sm mb-3 animate-fade-in-up" style={{ animationDelay: '200ms' }}>
                    <div className="p-4 border-b">
                        <span className="text-sm font-medium">💳 Phương thức thanh toán</span>
                    </div>
                    <div className="p-4 grid grid-cols-2 md:grid-cols-3 gap-3">
                        {paymentMethods.map((method, index) => (
                            <button
                                key={method.id}
                                onClick={() => setSelectedPayment(method.id)}
                                className={`p-4 border rounded-sm text-left transition-all hover-lift animate-fade-in-up ${selectedPayment === method.id
                                        ? 'border-[#ee4d2d] bg-[#fef6f5] ring-1 ring-[#ee4d2d]'
                                        : 'hover:border-[#ee4d2d]'
                                    }`}
                                style={{ animationDelay: `${index * 50}ms` }}
                            >
                                <div className="flex items-center gap-3">
                                    <span className={`text-2xl ${selectedPayment === method.id ? 'animate-wiggle' : ''}`}>
                                        {method.icon}
                                    </span>
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
                <div className="bg-white rounded-sm shadow-sm animate-fade-in-up" style={{ animationDelay: '300ms' }}>
                    <div className="p-4 border-b space-y-3">
                        <div className="flex justify-between text-sm">
                            <span className="text-gray-500">Tổng tiền hàng</span>
                            <span>₫{formatPrice(subtotal)}</span>
                        </div>
                        <div className="flex justify-between text-sm">
                            <span className="text-gray-500">Phí vận chuyển</span>
                            <span className={shippingFee === 0 ? 'text-[#00bfa5]' : ''}>
                                {shippingFee === 0 ? 'Miễn phí' : `₫${formatPrice(shippingFee)}`}
                            </span>
                        </div>
                        <div className="flex justify-between text-sm">
                            <span className="text-gray-500">Voucher Shopee</span>
                            <span className="text-[#ee4d2d]">-₫{formatPrice(discount)}</span>
                        </div>
                    </div>
                    <div className="p-4 flex items-center justify-between bg-gradient-to-r from-[#fef6f5] to-white">
                        <div>
                            <span className="text-gray-500 text-sm">Tổng thanh toán:</span>
                            <span className="text-3xl text-[#ee4d2d] font-bold ml-2 animate-pulse">₫{formatPrice(total)}</span>
                        </div>
                        <button
                            onClick={placeOrder}
                            disabled={isProcessing}
                            className={`px-12 py-3 bg-gradient-to-r from-[#ee4d2d] to-[#ff6533] text-white font-medium 
                         hover:opacity-90 transition-all hover-shrink ripple relative overflow-hidden
                         ${isProcessing ? 'cursor-wait' : ''}`}
                        >
                            {isProcessing ? (
                                <span className="flex items-center gap-2">
                                    <span className="loading-spinner" />
                                    Đang xử lý...
                                </span>
                            ) : (
                                'Đặt hàng'
                            )}
                        </button>
                    </div>

                    {/* Trust badges */}
                    <div className="p-4 border-t flex items-center justify-center gap-6 text-xs text-gray-500">
                        <span className="flex items-center gap-1">🔒 Thanh toán an toàn</span>
                        <span className="flex items-center gap-1">✅ Đảm bảo chính hãng</span>
                        <span className="flex items-center gap-1">🔄 Đổi trả dễ dàng</span>
                    </div>
                </div>
            </div>
        </div>
    );
}
