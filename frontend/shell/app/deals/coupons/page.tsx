'use client';

import { useState, useEffect } from 'react';
import { couponApi, Coupon } from '../../lib/api';

export default function CouponsPage() {
    const [coupons, setCoupons] = useState<Coupon[]>([]);
    const [loading, setLoading] = useState(true);
    const [claimedCoupons, setClaimedCoupons] = useState<string[]>([]);
    const [validationResult, setValidationResult] = useState<{ valid: boolean; discount: number; message: string } | null>(null);
    const [couponCode, setCouponCode] = useState('');

    useEffect(() => {
        loadCoupons();
    }, []);

    const loadCoupons = async () => {
        try {
            const data = await couponApi.getActiveCoupons();
            setCoupons(data || []);
        } catch {
            // Mock data
            setCoupons([
                {
                    id: 'c1',
                    code: 'WELCOME50',
                    discount_type: 'percentage',
                    discount_value: 50,
                    min_purchase: 100,
                    max_discount: 50,
                },
                {
                    id: 'c2',
                    code: 'FREESHIP',
                    discount_type: 'fixed',
                    discount_value: 10,
                    min_purchase: 50,
                    max_discount: 10,
                },
                {
                    id: 'c3',
                    code: 'MEGA20',
                    discount_type: 'percentage',
                    discount_value: 20,
                    min_purchase: 200,
                    max_discount: 100,
                },
                {
                    id: 'c4',
                    code: 'FLASH99',
                    discount_type: 'fixed',
                    discount_value: 99,
                    min_purchase: 500,
                    max_discount: 99,
                },
                {
                    id: 'c5',
                    code: 'VIP30',
                    discount_type: 'percentage',
                    discount_value: 30,
                    min_purchase: 300,
                    max_discount: 150,
                },
                {
                    id: 'c6',
                    code: 'NEWUSER',
                    discount_type: 'percentage',
                    discount_value: 15,
                    min_purchase: 0,
                    max_discount: 30,
                },
            ]);
        } finally {
            setLoading(false);
        }
    };

    const handleClaim = (couponId: string) => {
        if (!claimedCoupons.includes(couponId)) {
            setClaimedCoupons([...claimedCoupons, couponId]);
        }
    };

    const validateCoupon = async () => {
        if (!couponCode.trim()) return;

        try {
            const result = await couponApi.validate(couponCode, 'user-123', 150);
            setValidationResult(result);
        } catch {
            // Mock validation
            const coupon = coupons.find(c => c.code === couponCode.toUpperCase());
            if (coupon) {
                const discount = coupon.discount_type === 'percentage'
                    ? Math.min((150 * coupon.discount_value) / 100, coupon.max_discount)
                    : coupon.discount_value;
                setValidationResult({
                    valid: true,
                    discount,
                    message: `Coupon valid! You save $${discount.toFixed(2)}`,
                });
            } else {
                setValidationResult({
                    valid: false,
                    discount: 0,
                    message: 'Invalid coupon code',
                });
            }
        }
    };

    if (loading) {
        return (
            <div className="flex h-96 items-center justify-center">
                <div className="animate-pulse text-2xl">🎟️ Loading Coupons...</div>
            </div>
        );
    }

    return (
        <div className="container mx-auto py-8">
            {/* Header */}
            <div className="mb-8 rounded-xl bg-gradient-to-r from-purple-600 via-pink-600 to-red-600 p-8 text-white">
                <h1 className="text-4xl font-bold">🎟️ Coupon Center</h1>
                <p className="mt-2 text-lg opacity-90">
                    Claim exclusive coupons and save on your purchases
                </p>
            </div>

            {/* Coupon Validator */}
            <div className="mb-8 rounded-xl border bg-white p-6">
                <h2 className="mb-4 text-xl font-bold">🔍 Validate Coupon</h2>
                <div className="flex gap-4">
                    <input
                        type="text"
                        value={couponCode}
                        onChange={(e) => setCouponCode(e.target.value.toUpperCase())}
                        placeholder="Enter coupon code"
                        className="flex-1 rounded-lg border px-4 py-3 font-mono uppercase focus:outline-none focus:ring-2 focus:ring-purple-500"
                    />
                    <button
                        onClick={validateCoupon}
                        className="rounded-lg bg-gradient-to-r from-purple-600 to-pink-600 px-6 py-3 font-bold text-white hover:from-purple-700 hover:to-pink-700"
                    >
                        Validate
                    </button>
                </div>
                {validationResult && (
                    <div className={`mt-4 rounded-lg p-4 ${validationResult.valid ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                        <span className="text-xl mr-2">{validationResult.valid ? '✅' : '❌'}</span>
                        {validationResult.message}
                    </div>
                )}
            </div>

            {/* Coupons Grid */}
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {coupons.map((coupon) => {
                    const claimed = claimedCoupons.includes(coupon.id);

                    return (
                        <div
                            key={coupon.id}
                            className="relative overflow-hidden rounded-xl border-2 border-dashed border-gray-300 bg-gradient-to-br from-white to-gray-50"
                        >
                            {/* Discount badge */}
                            <div className="absolute -right-8 top-4 rotate-45 bg-red-600 px-12 py-1 text-sm font-bold text-white">
                                {coupon.discount_type === 'percentage'
                                    ? `${coupon.discount_value}% OFF`
                                    : `$${coupon.discount_value} OFF`}
                            </div>

                            <div className="p-6">
                                {/* Coupon code */}
                                <div className="mb-4 flex items-center gap-3">
                                    <div className="text-4xl">🎟️</div>
                                    <div>
                                        <div className="font-mono text-2xl font-bold tracking-wider">
                                            {coupon.code}
                                        </div>
                                        <div className="text-sm text-muted-foreground">
                                            {coupon.discount_type === 'percentage'
                                                ? `${coupon.discount_value}% off`
                                                : `$${coupon.discount_value} off`}
                                            {coupon.max_discount > 0 && ` (max $${coupon.max_discount})`}
                                        </div>
                                    </div>
                                </div>

                                {/* Terms */}
                                <div className="mb-4 text-sm text-muted-foreground">
                                    {coupon.min_purchase > 0
                                        ? `Min. purchase: $${coupon.min_purchase}`
                                        : 'No minimum purchase'}
                                </div>

                                {/* Claim button */}
                                <button
                                    onClick={() => handleClaim(coupon.id)}
                                    disabled={claimed}
                                    className={`w-full rounded-lg py-3 font-bold transition-all ${claimed
                                            ? 'bg-gray-200 text-gray-500 cursor-not-allowed'
                                            : 'bg-gradient-to-r from-purple-600 to-pink-600 text-white hover:from-purple-700 hover:to-pink-700 active:scale-95'
                                        }`}
                                >
                                    {claimed ? '✓ Claimed' : 'Claim Coupon'}
                                </button>
                            </div>

                            {/* Perforated edge */}
                            <div className="absolute left-0 top-1/2 h-6 w-3 -translate-y-1/2 rounded-r-full bg-gray-100"></div>
                            <div className="absolute right-0 top-1/2 h-6 w-3 -translate-y-1/2 rounded-l-full bg-gray-100"></div>
                        </div>
                    );
                })}
            </div>

            {/* My Coupons */}
            {claimedCoupons.length > 0 && (
                <div className="mt-12">
                    <h2 className="mb-4 text-2xl font-bold">📋 My Coupons</h2>
                    <div className="flex flex-wrap gap-3">
                        {claimedCoupons.map((id) => {
                            const coupon = coupons.find(c => c.id === id);
                            return coupon ? (
                                <div
                                    key={id}
                                    className="rounded-lg bg-gradient-to-r from-purple-100 to-pink-100 px-4 py-2 font-mono font-bold"
                                >
                                    {coupon.code}
                                </div>
                            ) : null;
                        })}
                    </div>
                </div>
            )}

            {/* How it works */}
            <div className="mt-12 rounded-xl border bg-gradient-to-r from-slate-50 to-slate-100 p-8">
                <h2 className="mb-6 text-2xl font-bold">How to Use Coupons</h2>
                <div className="grid gap-4 md:grid-cols-4">
                    {[
                        { step: '1', icon: '🎟️', title: 'Claim', desc: 'Get your coupon' },
                        { step: '2', icon: '🛒', title: 'Shop', desc: 'Add items to cart' },
                        { step: '3', icon: '✂️', title: 'Apply', desc: 'Enter code at checkout' },
                        { step: '4', icon: '💰', title: 'Save', desc: 'Enjoy the discount!' },
                    ].map((item) => (
                        <div key={item.step} className="text-center">
                            <div className="mx-auto mb-2 h-12 w-12 rounded-full bg-purple-600 flex items-center justify-center text-xl font-bold text-white">
                                {item.step}
                            </div>
                            <div className="text-3xl">{item.icon}</div>
                            <h3 className="font-semibold">{item.title}</h3>
                            <p className="text-sm text-muted-foreground">{item.desc}</p>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
