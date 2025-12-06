'use client';

import React, { useState } from 'react';
import Link from 'next/link';

interface CartItem {
    id: string;
    name: string;
    price: number;
    quantity: number;
    image: string;
}

export default function CheckoutPage() {
    const [step, setStep] = useState(1);
    const [selectedPayment, setSelectedPayment] = useState('cod');
    const [selectedShipping, setSelectedShipping] = useState('standard');
    const [orderPlaced, setOrderPlaced] = useState(false);
    const [isProcessing, setIsProcessing] = useState(false);

    const [address, setAddress] = useState({
        name: 'Nguyễn Văn A',
        phone: '0901234567',
        address: '123 Đường ABC, Phường XYZ, Quận 1, TP.HCM',
    });

    const cartItems: CartItem[] = [
        { id: 'p1', name: 'iPhone 15 Pro Max 256GB', price: 29990000, quantity: 1, image: '📱' },
        { id: 'p6', name: 'Son Dưỡng Môi Dior', price: 950000, quantity: 2, image: '💄' },
    ];

    const paymentMethods = [
        { id: 'cod', name: 'Thanh toán khi nhận hàng', icon: '💵', desc: 'Thanh toán bằng tiền mặt' },
        { id: 'momo', name: 'Ví MoMo', icon: '🟣', desc: 'Thanh toán qua ví điện tử MoMo' },
        { id: 'vnpay', name: 'VNPay', icon: '🔵', desc: 'Thanh toán qua cổng VNPay' },
        { id: 'zalopay', name: 'ZaloPay', icon: '🔷', desc: 'Thanh toán qua ví ZaloPay' },
        { id: 'card', name: 'Thẻ tín dụng/Ghi nợ', icon: '💳', desc: 'Visa, Mastercard, JCB' },
        { id: 'bank', name: 'Chuyển khoản ngân hàng', icon: '🏦', desc: 'Chuyển khoản trực tiếp' },
    ];

    const shippingMethods = [
        { id: 'standard', name: 'Giao hàng tiêu chuẩn', time: '3-5 ngày', price: 0, icon: '📦' },
        { id: 'fast', name: 'Giao hàng nhanh', time: '1-2 ngày', price: 25000, icon: '🚀' },
        { id: 'express', name: 'Hỏa tốc', time: 'Trong ngày', price: 50000, icon: '⚡' },
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
            <div className="min-h-screen bg-[#F5F5F5] flex items-center justify-center animate-fade-in">
                <div className="bg-white rounded-2xl p-12 text-center max-w-md mx-4 shadow-xl">
                    <div className="text-8xl mb-6 animate-bounce">🎉</div>
                    <h1 className="text-2xl font-bold text-green-600 mb-2">Đặt hàng thành công!</h1>
                    <p className="text-gray-500 mb-6">Cảm ơn bạn đã mua hàng. Đơn hàng của bạn đang được xử lý.</p>
                    <div className="bg-gray-50 rounded-xl p-4 mb-6">
                        <p className="text-sm text-gray-500">Mã đơn hàng</p>
                        <p className="text-xl font-bold text-[#EE4D2D]">#SP{Date.now().toString().slice(-8)}</p>
                    </div>
                    <div className="flex gap-3">
                        <Link href="/products" className="flex-1 py-3 border-2 border-[#EE4D2D] text-[#EE4D2D] rounded-xl font-bold hover:bg-[#FFEEE8]">
                            Tiếp tục mua
                        </Link>
                        <Link href="/" className="flex-1 py-3 bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] text-white rounded-xl font-bold">
                            Trang chủ
                        </Link>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-[#F5F5F5] animate-fade-in">
            {/* Header */}
            <div className="bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] py-6">
                <div className="container mx-auto px-4">
                    <h1 className="text-2xl font-bold text-white flex items-center gap-2">
                        <span>🛒</span> Thanh Toán
                    </h1>
                </div>
            </div>

            {/* Progress Steps */}
            <div className="bg-white border-b">
                <div className="container mx-auto px-4 py-4">
                    <div className="flex items-center justify-center gap-4">
                        {['Địa chỉ', 'Vận chuyển', 'Thanh toán', 'Xác nhận'].map((s, i) => (
                            <div key={s} className="flex items-center">
                                <div className={`w-8 h-8 rounded-full flex items-center justify-center font-bold text-sm ${step > i + 1 ? 'bg-green-500 text-white' :
                                        step === i + 1 ? 'bg-[#EE4D2D] text-white' :
                                            'bg-gray-200 text-gray-500'
                                    }`}>
                                    {step > i + 1 ? '✓' : i + 1}
                                </div>
                                <span className={`ml-2 text-sm hidden sm:block ${step === i + 1 ? 'font-bold text-[#EE4D2D]' : 'text-gray-500'}`}>
                                    {s}
                                </span>
                                {i < 3 && <div className={`w-12 h-1 mx-2 rounded ${step > i + 1 ? 'bg-green-500' : 'bg-gray-200'}`} />}
                            </div>
                        ))}
                    </div>
                </div>
            </div>

            <div className="container mx-auto px-4 py-6">
                <div className="grid lg:grid-cols-3 gap-6">
                    {/* Main Content */}
                    <div className="lg:col-span-2 space-y-6">
                        {/* Step 1: Address */}
                        {step >= 1 && (
                            <div className={`bg-white rounded-2xl shadow-sm overflow-hidden ${step === 1 ? 'ring-2 ring-[#EE4D2D]' : ''}`}>
                                <div className="p-4 border-b flex items-center justify-between bg-gray-50">
                                    <h2 className="font-bold flex items-center gap-2">
                                        <span className="w-6 h-6 bg-[#EE4D2D] text-white rounded-full flex items-center justify-center text-sm">1</span>
                                        Địa chỉ nhận hàng
                                    </h2>
                                    {step > 1 && (
                                        <button onClick={() => setStep(1)} className="text-[#EE4D2D] text-sm hover:underline">Thay đổi</button>
                                    )}
                                </div>
                                <div className="p-4">
                                    <div className="flex items-start gap-4 p-4 border-2 border-[#EE4D2D] rounded-xl bg-[#FFEEE8]">
                                        <span className="text-2xl">📍</span>
                                        <div className="flex-1">
                                            <p className="font-bold">{address.name} | {address.phone}</p>
                                            <p className="text-gray-600">{address.address}</p>
                                        </div>
                                        <span className="px-2 py-1 bg-[#EE4D2D] text-white text-xs rounded-lg">Mặc định</span>
                                    </div>
                                    {step === 1 && (
                                        <button
                                            onClick={() => setStep(2)}
                                            className="w-full mt-4 py-3 bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] text-white rounded-xl font-bold hover:opacity-90"
                                        >
                                            Tiếp tục
                                        </button>
                                    )}
                                </div>
                            </div>
                        )}

                        {/* Step 2: Shipping */}
                        {step >= 2 && (
                            <div className={`bg-white rounded-2xl shadow-sm overflow-hidden ${step === 2 ? 'ring-2 ring-[#EE4D2D]' : ''}`}>
                                <div className="p-4 border-b flex items-center justify-between bg-gray-50">
                                    <h2 className="font-bold flex items-center gap-2">
                                        <span className="w-6 h-6 bg-[#EE4D2D] text-white rounded-full flex items-center justify-center text-sm">2</span>
                                        Phương thức vận chuyển
                                    </h2>
                                    {step > 2 && (
                                        <button onClick={() => setStep(2)} className="text-[#EE4D2D] text-sm hover:underline">Thay đổi</button>
                                    )}
                                </div>
                                <div className="p-4 space-y-3">
                                    {shippingMethods.map(method => (
                                        <button
                                            key={method.id}
                                            onClick={() => setSelectedShipping(method.id)}
                                            className={`w-full p-4 rounded-xl border-2 text-left flex items-center gap-4 transition-all ${selectedShipping === method.id
                                                    ? 'border-[#EE4D2D] bg-[#FFEEE8]'
                                                    : 'border-gray-200 hover:border-[#EE4D2D]'
                                                }`}
                                        >
                                            <span className="text-3xl">{method.icon}</span>
                                            <div className="flex-1">
                                                <p className="font-bold">{method.name}</p>
                                                <p className="text-sm text-gray-500">{method.time}</p>
                                            </div>
                                            <span className={`font-bold ${method.price === 0 ? 'text-green-500' : 'text-[#EE4D2D]'}`}>
                                                {method.price === 0 ? 'Miễn phí' : `₫${formatPrice(method.price)}`}
                                            </span>
                                        </button>
                                    ))}
                                    {step === 2 && (
                                        <button
                                            onClick={() => setStep(3)}
                                            className="w-full mt-4 py-3 bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] text-white rounded-xl font-bold hover:opacity-90"
                                        >
                                            Tiếp tục
                                        </button>
                                    )}
                                </div>
                            </div>
                        )}

                        {/* Step 3: Payment */}
                        {step >= 3 && (
                            <div className={`bg-white rounded-2xl shadow-sm overflow-hidden ${step === 3 ? 'ring-2 ring-[#EE4D2D]' : ''}`}>
                                <div className="p-4 border-b flex items-center justify-between bg-gray-50">
                                    <h2 className="font-bold flex items-center gap-2">
                                        <span className="w-6 h-6 bg-[#EE4D2D] text-white rounded-full flex items-center justify-center text-sm">3</span>
                                        Phương thức thanh toán
                                    </h2>
                                    {step > 3 && (
                                        <button onClick={() => setStep(3)} className="text-[#EE4D2D] text-sm hover:underline">Thay đổi</button>
                                    )}
                                </div>
                                <div className="p-4 grid sm:grid-cols-2 gap-3">
                                    {paymentMethods.map(method => (
                                        <button
                                            key={method.id}
                                            onClick={() => setSelectedPayment(method.id)}
                                            className={`p-4 rounded-xl border-2 text-left transition-all ${selectedPayment === method.id
                                                    ? 'border-[#EE4D2D] bg-[#FFEEE8]'
                                                    : 'border-gray-200 hover:border-[#EE4D2D]'
                                                }`}
                                        >
                                            <div className="flex items-center gap-3">
                                                <span className="text-3xl">{method.icon}</span>
                                                <div>
                                                    <p className="font-bold text-sm">{method.name}</p>
                                                    <p className="text-xs text-gray-500">{method.desc}</p>
                                                </div>
                                            </div>
                                        </button>
                                    ))}
                                </div>
                                {step === 3 && (
                                    <div className="p-4 pt-0">
                                        <button
                                            onClick={() => setStep(4)}
                                            className="w-full py-3 bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] text-white rounded-xl font-bold hover:opacity-90"
                                        >
                                            Tiếp tục
                                        </button>
                                    </div>
                                )}
                            </div>
                        )}

                        {/* Step 4: Confirm */}
                        {step === 4 && (
                            <div className="bg-white rounded-2xl shadow-sm overflow-hidden ring-2 ring-[#EE4D2D]">
                                <div className="p-4 border-b bg-gray-50">
                                    <h2 className="font-bold flex items-center gap-2">
                                        <span className="w-6 h-6 bg-[#EE4D2D] text-white rounded-full flex items-center justify-center text-sm">4</span>
                                        Xác nhận đơn hàng
                                    </h2>
                                </div>
                                <div className="p-4">
                                    {cartItems.map(item => (
                                        <div key={item.id} className="flex items-center gap-4 py-3 border-b last:border-0">
                                            <span className="text-4xl">{item.image}</span>
                                            <div className="flex-1">
                                                <p className="font-semibold">{item.name}</p>
                                                <p className="text-sm text-gray-500">x{item.quantity}</p>
                                            </div>
                                            <span className="font-bold text-[#EE4D2D]">₫{formatPrice(item.price * item.quantity)}</span>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}
                    </div>

                    {/* Order Summary */}
                    <div className="lg:col-span-1">
                        <div className="bg-white rounded-2xl shadow-sm p-6 sticky top-24">
                            <h3 className="font-bold text-lg mb-4">Tóm tắt đơn hàng</h3>

                            <div className="space-y-3 text-sm">
                                <div className="flex justify-between">
                                    <span className="text-gray-500">Tạm tính ({cartItems.length} sản phẩm)</span>
                                    <span>₫{formatPrice(subtotal)}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-gray-500">Phí vận chuyển</span>
                                    <span className={shippingFee === 0 ? 'text-green-500' : ''}>
                                        {shippingFee === 0 ? 'Miễn phí' : `₫${formatPrice(shippingFee)}`}
                                    </span>
                                </div>
                                <div className="flex justify-between text-green-500">
                                    <span>Giảm giá voucher</span>
                                    <span>-₫{formatPrice(discount)}</span>
                                </div>
                                <div className="pt-3 mt-3 border-t flex justify-between font-bold text-lg">
                                    <span>Tổng cộng</span>
                                    <span className="text-[#EE4D2D]">₫{formatPrice(total)}</span>
                                </div>
                            </div>

                            {step === 4 && (
                                <button
                                    onClick={placeOrder}
                                    disabled={isProcessing}
                                    className={`w-full mt-6 py-4 bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] text-white rounded-xl font-bold text-lg transition-all ${isProcessing ? 'opacity-70' : 'hover:opacity-90'
                                        }`}
                                >
                                    {isProcessing ? (
                                        <span className="flex items-center justify-center gap-2">
                                            <span className="animate-spin">⏳</span> Đang xử lý...
                                        </span>
                                    ) : (
                                        `Đặt hàng (₫${formatPrice(total)})`
                                    )}
                                </button>
                            )}

                            <div className="mt-4 p-3 bg-[#FFEEE8] rounded-xl">
                                <p className="text-xs text-gray-600 text-center">
                                    🔒 Thanh toán an toàn & bảo mật
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
