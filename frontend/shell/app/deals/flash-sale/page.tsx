'use client';

import { useState, useEffect, useCallback } from 'react';
import { flashSaleApi, FlashSale } from '../../lib/api';

export default function FlashSalesPage() {
    const [sales, setSales] = useState<FlashSale[]>([]);
    const [loading, setLoading] = useState(true);
    const [selectedSale, setSelectedSale] = useState<FlashSale | null>(null);
    const [purchasing, setPurchasing] = useState(false);
    const [powStatus, setPowStatus] = useState<'idle' | 'solving' | 'solved'>('idle');

    useEffect(() => {
        loadSales();
        const interval = setInterval(loadSales, 5000);
        return () => clearInterval(interval);
    }, []);

    const loadSales = async () => {
        try {
            const data = await flashSaleApi.getActiveSales();
            setSales(data || []);
        } catch {
            // Mock data for demo
            setSales([
                {
                    id: 'fs-001',
                    product_id: 'prod-001',
                    original_price: 999,
                    sale_price: 499,
                    discount_percent: 50,
                    total_quantity: 1000,
                    sold_quantity: 750,
                    max_per_user: 2,
                    status: 'ACTIVE',
                    start_time: new Date().toISOString(),
                    end_time: new Date(Date.now() + 3600000).toISOString(),
                },
                {
                    id: 'fs-002',
                    product_id: 'prod-002',
                    original_price: 2499,
                    sale_price: 999,
                    discount_percent: 60,
                    total_quantity: 500,
                    sold_quantity: 489,
                    max_per_user: 1,
                    status: 'ACTIVE',
                    start_time: new Date().toISOString(),
                    end_time: new Date(Date.now() + 1800000).toISOString(),
                },
            ]);
        } finally {
            setLoading(false);
        }
    };

    const solvePoW = useCallback(async (challenge: string, difficulty: number): Promise<string> => {
        // Simplified PoW solving (real implementation would hash with SHA256)
        let nonce = 0;
        const prefix = '0'.repeat(difficulty);

        return new Promise((resolve) => {
            const solve = () => {
                for (let i = 0; i < 10000; i++) {
                    const testNonce = (nonce++).toString();
                    // Simulated hash check
                    if (Math.random() < 0.0001) {
                        resolve(testNonce);
                        return;
                    }
                }
                setTimeout(solve, 0);
            };
            solve();
        });
    }, []);

    const handlePurchase = async (sale: FlashSale) => {
        setSelectedSale(sale);
        setPurchasing(true);
        setPowStatus('solving');

        try {
            // Get challenge
            const challengeData = await flashSaleApi.getChallenge(sale.id, 'user-123');

            // Solve PoW (simulated)
            await new Promise(r => setTimeout(r, 1500)); // Simulate solving
            setPowStatus('solved');

            // Attempt purchase
            const result = await flashSaleApi.attemptPurchase({
                flash_sale_id: sale.id,
                user_id: 'user-123',
                quantity: 1,
                challenge: challengeData.challenge,
                nonce: 'solved-nonce',
            });

            alert(`🎉 Reservation successful! ID: ${result.reservation_id}`);
        } catch (err) {
            console.error(err);
        } finally {
            setPurchasing(false);
            setSelectedSale(null);
            setPowStatus('idle');
        }
    };

    const getTimeRemaining = (endTime: string) => {
        const diff = new Date(endTime).getTime() - Date.now();
        if (diff <= 0) return '00:00:00';
        const hours = Math.floor(diff / 3600000);
        const mins = Math.floor((diff % 3600000) / 60000);
        const secs = Math.floor((diff % 60000) / 1000);
        return `${hours.toString().padStart(2, '0')}:${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
    };

    if (loading) {
        return (
            <div className="flex h-96 items-center justify-center">
                <div className="animate-pulse text-2xl">⚡ Loading Flash Sales...</div>
            </div>
        );
    }

    return (
        <div className="container mx-auto py-8">
            {/* Header */}
            <div className="mb-8 rounded-xl bg-gradient-to-r from-red-600 to-orange-500 p-8 text-white">
                <div className="flex items-center gap-4">
                    <span className="text-5xl">⚡</span>
                    <div>
                        <h1 className="text-4xl font-bold">Flash Sale</h1>
                        <p className="text-lg opacity-90">11.11 Mega Deals - Up to 90% Off!</p>
                    </div>
                </div>
                <div className="mt-4 flex gap-8 text-sm">
                    <div className="flex items-center gap-2">
                        <span className="h-2 w-2 animate-pulse rounded-full bg-white"></span>
                        1.2M users online
                    </div>
                    <div>🔥 3,456 orders/second</div>
                    <div>🛡️ PoW Protected</div>
                </div>
            </div>

            {/* PoW Modal */}
            {purchasing && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
                    <div className="w-96 rounded-xl bg-white p-8 text-center">
                        <div className="mb-4 text-6xl">
                            {powStatus === 'solving' ? '🧮' : '✅'}
                        </div>
                        <h3 className="mb-2 text-xl font-bold">
                            {powStatus === 'solving' ? 'Solving Challenge...' : 'Challenge Solved!'}
                        </h3>
                        <p className="mb-4 text-muted-foreground">
                            {powStatus === 'solving'
                                ? 'Anti-bot verification in progress'
                                : 'Submitting your purchase...'}
                        </p>
                        <div className="h-2 overflow-hidden rounded-full bg-gray-200">
                            <div
                                className={`h-full bg-gradient-to-r from-orange-500 to-red-500 transition-all duration-1000 ${powStatus === 'solved' ? 'w-full' : 'w-2/3 animate-pulse'
                                    }`}
                            />
                        </div>
                    </div>
                </div>
            )}

            {/* Sales Grid */}
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {sales.map((sale) => {
                    const progress = (sale.sold_quantity / sale.total_quantity) * 100;
                    const remaining = sale.total_quantity - sale.sold_quantity;

                    return (
                        <div
                            key={sale.id}
                            className="group relative overflow-hidden rounded-xl border bg-white shadow-lg transition-transform hover:scale-[1.02]"
                        >
                            {/* Discount Badge */}
                            <div className="absolute right-0 top-0 rounded-bl-xl bg-red-600 px-4 py-2 text-lg font-bold text-white">
                                -{sale.discount_percent}%
                            </div>

                            {/* Product Image Placeholder */}
                            <div className="aspect-square bg-gradient-to-br from-gray-100 to-gray-200 p-8">
                                <div className="flex h-full items-center justify-center text-8xl">
                                    📱
                                </div>
                            </div>

                            <div className="p-4">
                                {/* Price */}
                                <div className="mb-2 flex items-baseline gap-2">
                                    <span className="text-2xl font-bold text-red-600">
                                        ${sale.sale_price}
                                    </span>
                                    <span className="text-lg text-gray-400 line-through">
                                        ${sale.original_price}
                                    </span>
                                </div>

                                {/* Timer */}
                                <div className="mb-3 flex items-center gap-2 text-sm text-orange-600">
                                    <span>⏰</span>
                                    <span className="font-mono font-bold">
                                        {getTimeRemaining(sale.end_time)}
                                    </span>
                                </div>

                                {/* Progress Bar */}
                                <div className="mb-3">
                                    <div className="mb-1 flex justify-between text-xs">
                                        <span className="text-red-600 font-semibold">
                                            🔥 {Math.round(progress)}% Claimed
                                        </span>
                                        <span className="text-gray-500">
                                            {remaining} left
                                        </span>
                                    </div>
                                    <div className="h-3 overflow-hidden rounded-full bg-gray-200">
                                        <div
                                            className="h-full bg-gradient-to-r from-yellow-400 via-red-500 to-red-600 transition-all"
                                            style={{ width: `${progress}%` }}
                                        />
                                    </div>
                                </div>

                                {/* Buy Button */}
                                <button
                                    onClick={() => handlePurchase(sale)}
                                    disabled={remaining === 0 || purchasing}
                                    className={`w-full rounded-lg py-3 font-bold text-white transition-all ${remaining === 0
                                            ? 'bg-gray-400 cursor-not-allowed'
                                            : 'bg-gradient-to-r from-red-600 to-orange-500 hover:from-red-700 hover:to-orange-600 active:scale-95'
                                        }`}
                                >
                                    {remaining === 0 ? '🚫 Sold Out' : '⚡ Flash Buy Now'}
                                </button>
                            </div>
                        </div>
                    );
                })}
            </div>

            {/* How it Works */}
            <div className="mt-12 rounded-xl border bg-gradient-to-r from-slate-50 to-slate-100 p-8">
                <h2 className="mb-6 text-2xl font-bold">How Flash Sale Works</h2>
                <div className="grid gap-4 md:grid-cols-4">
                    {[
                        { icon: '🔐', title: '1. Get Challenge', desc: 'Receive a cryptographic puzzle' },
                        { icon: '🧮', title: '2. Solve PoW', desc: 'Your browser solves the challenge' },
                        { icon: '⚡', title: '3. Submit', desc: 'Atomic stock decrement with Redis' },
                        { icon: '✅', title: '4. Reserve', desc: '5 minutes to complete payment' },
                    ].map((step) => (
                        <div key={step.title} className="text-center">
                            <div className="mb-2 text-4xl">{step.icon}</div>
                            <h3 className="font-semibold">{step.title}</h3>
                            <p className="text-sm text-muted-foreground">{step.desc}</p>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
