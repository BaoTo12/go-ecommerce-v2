'use client';

import React, { useState } from 'react';

interface Coupon {
    id: string;
    code: string;
    discount: string;
    minSpend: number;
    maxDiscount: number;
    expiry: string;
    type: 'percent' | 'fixed';
    category: string;
}

export default function CouponsPage() {
    const [claimedCoupons, setClaimedCoupons] = useState<string[]>([]);
    const [couponCode, setCouponCode] = useState('');
    const [validationResult, setValidationResult] = useState<{ valid: boolean; message: string } | null>(null);

    const coupons: Coupon[] = [
        { id: 'c1', code: 'GIẢM50K', discount: '50.000₫', minSpend: 500000, maxDiscount: 50000, expiry: '31/12/2024', type: 'fixed', category: 'Tất cả' },
        { id: 'c2', code: 'FREESHIP', discount: 'Freeship', minSpend: 0, maxDiscount: 30000, expiry: '31/12/2024', type: 'fixed', category: 'Vận chuyển' },
        { id: 'c3', code: 'THUDONG20', discount: '20%', minSpend: 200000, maxDiscount: 100000, expiry: '31/01/2025', type: 'percent', category: 'Thời trang' },
        { id: 'c4', code: 'TECH15', discount: '15%', minSpend: 1000000, maxDiscount: 500000, expiry: '15/01/2025', type: 'percent', category: 'Điện tử' },
        { id: 'c5', code: 'NEWUSER', discount: '100.000₫', minSpend: 0, maxDiscount: 100000, expiry: '31/12/2024', type: 'fixed', category: 'Người mới' },
        { id: 'c6', code: 'BEAUTY30', discount: '30%', minSpend: 300000, maxDiscount: 150000, expiry: '28/02/2025', type: 'percent', category: 'Làm đẹp' },
    ];

    const handleClaim = (couponId: string) => {
        if (!claimedCoupons.includes(couponId)) {
            setClaimedCoupons([...claimedCoupons, couponId]);
        }
    };

    const validateCoupon = () => {
        if (!couponCode.trim()) return;
        const coupon = coupons.find(c => c.code === couponCode.toUpperCase());
        if (coupon) {
            setValidationResult({ valid: true, message: `Mã hợp lệ! Giảm ${coupon.discount}` });
        } else {
            setValidationResult({ valid: false, message: 'Mã không hợp lệ' });
        }
    };

    return (
        <div className="min-h-screen bg-[#F5F5F5]">
            {/* Header */}
            <div className="bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] py-6">
                <div className="container mx-auto px-4">
                    <h1 className="text-2xl font-bold text-white">🎟️ Mã Giảm Giá</h1>
                    <p className="text-white/80 text-sm">Săn voucher, mua sắm tiết kiệm!</p>
                </div>
            </div>

            <div className="container mx-auto px-4 py-6">
                {/* Validate Coupon */}
                <div className="bg-white rounded-sm p-4 mb-6">
                    <h2 className="font-bold mb-3 text-[#EE4D2D]">🔍 Kiểm tra mã giảm giá</h2>
                    <div className="flex gap-2">
                        <input
                            type="text"
                            value={couponCode}
                            onChange={e => setCouponCode(e.target.value.toUpperCase())}
                            placeholder="Nhập mã giảm giá"
                            className="flex-1 border border-gray-300 rounded px-4 py-2 focus:outline-none focus:border-[#EE4D2D]"
                        />
                        <button
                            onClick={validateCoupon}
                            className="bg-[#EE4D2D] text-white px-6 py-2 rounded font-semibold hover:bg-[#D73211]"
                        >
                            Kiểm tra
                        </button>
                    </div>
                    {validationResult && (
                        <div className={`mt-3 p-3 rounded ${validationResult.valid ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                            {validationResult.valid ? '✅' : '❌'} {validationResult.message}
                        </div>
                    )}
                </div>

                {/* Coupons Grid */}
                <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {coupons.map(coupon => {
                        const claimed = claimedCoupons.includes(coupon.id);

                        return (
                            <div
                                key={coupon.id}
                                className="bg-white rounded-sm overflow-hidden flex"
                            >
                                {/* Left Part - Discount */}
                                <div className="w-32 bg-gradient-to-br from-[#EE4D2D] to-[#FF6633] text-white p-4 flex flex-col items-center justify-center relative">
                                    <div className="text-2xl font-bold">{coupon.discount}</div>
                                    <div className="text-xs opacity-80 mt-1">{coupon.category}</div>

                                    {/* Perforated edge */}
                                    <div className="absolute right-0 top-0 bottom-0 w-3">
                                        {[...Array(8)].map((_, i) => (
                                            <div key={i} className="w-3 h-3 bg-[#F5F5F5] rounded-full -mr-1.5" style={{ marginTop: i === 0 ? 0 : 8 }} />
                                        ))}
                                    </div>
                                </div>

                                {/* Right Part - Details */}
                                <div className="flex-1 p-4 border-t border-b border-r border-dashed border-gray-300">
                                    <div className="font-mono text-lg font-bold text-[#EE4D2D] mb-2">{coupon.code}</div>

                                    <div className="text-xs text-gray-500 space-y-1">
                                        {coupon.minSpend > 0 && (
                                            <p>Đơn tối thiểu: {coupon.minSpend.toLocaleString()}₫</p>
                                        )}
                                        <p>Giảm tối đa: {coupon.maxDiscount.toLocaleString()}₫</p>
                                        <p>HSD: {coupon.expiry}</p>
                                    </div>

                                    <button
                                        onClick={() => handleClaim(coupon.id)}
                                        disabled={claimed}
                                        className={`mt-3 w-full py-2 rounded text-sm font-semibold transition-all ${claimed
                                                ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                                                : 'bg-[#EE4D2D] text-white hover:bg-[#D73211]'
                                            }`}
                                    >
                                        {claimed ? '✓ Đã lưu' : 'Lưu mã'}
                                    </button>
                                </div>
                            </div>
                        );
                    })}
                </div>

                {/* My Coupons */}
                {claimedCoupons.length > 0 && (
                    <div className="mt-8 bg-white rounded-sm p-4">
                        <h2 className="font-bold mb-3 text-[#EE4D2D]">📋 Mã của tôi ({claimedCoupons.length})</h2>
                        <div className="flex flex-wrap gap-2">
                            {claimedCoupons.map(id => {
                                const coupon = coupons.find(c => c.id === id);
                                return coupon ? (
                                    <span key={id} className="bg-[#FFEEE8] text-[#EE4D2D] px-3 py-1 rounded font-mono font-bold">
                                        {coupon.code}
                                    </span>
                                ) : null;
                            })}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
