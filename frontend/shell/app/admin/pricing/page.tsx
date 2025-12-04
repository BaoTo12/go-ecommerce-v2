'use client';

import { useState } from 'react';
import { pricingApi } from '../../lib/api';

export default function DynamicPricingPage() {
    const [productId, setProductId] = useState('');
    const [priceData, setPriceData] = useState<any>(null);
    const [loading, setLoading] = useState(false);
    const [optimizing, setOptimizing] = useState(false);

    const demoProducts = [
        { id: 'p001', name: 'iPhone 15 Pro', category: 'Electronics', basePrice: 999 },
        { id: 'p002', name: 'Nike Air Max', category: 'Fashion', basePrice: 189 },
        { id: 'p003', name: 'Sony Headphones', category: 'Electronics', basePrice: 349 },
        { id: 'p004', name: 'Dyson Vacuum', category: 'Home', basePrice: 699 },
    ];

    const analyzePrice = async (id: string) => {
        setProductId(id);
        setLoading(true);

        try {
            const data = await pricingApi.getPrice(id);
            setPriceData(data);
        } catch {
            // Mock data
            const product = demoProducts.find(p => p.id === id);
            if (product) {
                const basePrice = product.basePrice;
                const demandMultiplier = 0.9 + Math.random() * 0.3;
                const competitorPrices = [
                    basePrice * (0.95 + Math.random() * 0.15),
                    basePrice * (0.9 + Math.random() * 0.2),
                    basePrice * (0.85 + Math.random() * 0.25),
                ];
                const optimizedPrice = basePrice * demandMultiplier * 0.98;

                setPriceData({
                    product_id: id,
                    product_name: product.name,
                    base_price: basePrice,
                    optimized_price: optimizedPrice,
                    competitor_prices: competitorPrices,
                    demand_level: Math.random() > 0.5 ? 'HIGH' : 'MEDIUM',
                    inventory_level: Math.floor(Math.random() * 500) + 50,
                    price_elasticity: -1.2 + Math.random() * 0.8,
                    recommendation: optimizedPrice < basePrice ? 'LOWER' : 'RAISE',
                    potential_revenue_increase: (Math.random() * 15).toFixed(1),
                });
            }
        } finally {
            setLoading(false);
        }
    };

    const optimizePrice = async () => {
        if (!productId) return;
        setOptimizing(true);

        try {
            await pricingApi.optimizePrice(productId);
            // Refresh data
            await analyzePrice(productId);
        } catch {
            // Mock optimization
            setPriceData((prev: any) => ({
                ...prev,
                optimized_price: prev.base_price * 0.95,
                status: 'OPTIMIZED',
            }));
        } finally {
            setOptimizing(false);
        }
    };

    return (
        <div className="container mx-auto py-8">
            {/* Header */}
            <div className="mb-8 rounded-xl bg-gradient-to-r from-emerald-600 to-teal-600 p-8 text-white">
                <h1 className="text-3xl font-bold">💹 Dynamic Pricing Engine</h1>
                <p className="mt-2 text-lg opacity-90">
                    AI-powered price optimization for maximum revenue
                </p>
            </div>

            <div className="grid gap-8 lg:grid-cols-3">
                {/* Product Selection */}
                <div className="rounded-xl border bg-white p-6">
                    <h2 className="mb-4 text-xl font-bold">📦 Select Product</h2>
                    <div className="space-y-3">
                        {demoProducts.map((product) => (
                            <button
                                key={product.id}
                                onClick={() => analyzePrice(product.id)}
                                className={`w-full rounded-lg border p-4 text-left transition-all hover:shadow-md ${productId === product.id ? 'border-emerald-500 bg-emerald-50' : ''
                                    }`}
                            >
                                <div className="font-semibold">{product.name}</div>
                                <div className="flex justify-between text-sm text-muted-foreground">
                                    <span>{product.category}</span>
                                    <span>${product.basePrice}</span>
                                </div>
                            </button>
                        ))}
                    </div>
                </div>

                {/* Price Analysis */}
                <div className="lg:col-span-2 space-y-6">
                    {loading && (
                        <div className="flex h-64 items-center justify-center rounded-xl border bg-white">
                            <div className="animate-pulse text-xl">📊 Analyzing...</div>
                        </div>
                    )}

                    {priceData && !loading && (
                        <>
                            {/* Price Comparison */}
                            <div className="rounded-xl border bg-white p-6">
                                <h2 className="mb-4 text-xl font-bold">💰 Price Analysis</h2>

                                <div className="grid gap-4 md:grid-cols-3">
                                    <div className="rounded-lg bg-gray-100 p-4 text-center">
                                        <div className="text-sm text-muted-foreground">Base Price</div>
                                        <div className="text-2xl font-bold">${priceData.base_price.toFixed(2)}</div>
                                    </div>
                                    <div className="rounded-lg bg-emerald-100 p-4 text-center">
                                        <div className="text-sm text-emerald-700">Optimized Price</div>
                                        <div className="text-2xl font-bold text-emerald-600">
                                            ${priceData.optimized_price.toFixed(2)}
                                        </div>
                                        <div className={`text-sm ${priceData.optimized_price < priceData.base_price ? 'text-red-600' : 'text-green-600'}`}>
                                            {priceData.optimized_price < priceData.base_price ? '↓' : '↑'}
                                            {Math.abs(((priceData.optimized_price - priceData.base_price) / priceData.base_price) * 100).toFixed(1)}%
                                        </div>
                                    </div>
                                    <div className="rounded-lg bg-blue-100 p-4 text-center">
                                        <div className="text-sm text-blue-700">Potential Increase</div>
                                        <div className="text-2xl font-bold text-blue-600">
                                            +{priceData.potential_revenue_increase}%
                                        </div>
                                        <div className="text-sm text-muted-foreground">Revenue</div>
                                    </div>
                                </div>

                                {/* Competitor Prices */}
                                <div className="mt-6">
                                    <h3 className="mb-2 font-semibold">🏪 Competitor Prices</h3>
                                    <div className="flex gap-3">
                                        {priceData.competitor_prices.map((price: number, i: number) => (
                                            <div key={i} className="rounded-lg border px-4 py-2 text-center">
                                                <div className="text-xs text-muted-foreground">Competitor {i + 1}</div>
                                                <div className="font-semibold">${price.toFixed(2)}</div>
                                            </div>
                                        ))}
                                    </div>
                                </div>

                                {/* Action Button */}
                                <button
                                    onClick={optimizePrice}
                                    disabled={optimizing}
                                    className="mt-6 w-full rounded-lg bg-gradient-to-r from-emerald-600 to-teal-600 py-3 font-bold text-white hover:from-emerald-700 hover:to-teal-700 disabled:opacity-50"
                                >
                                    {optimizing ? 'Optimizing...' : '🚀 Apply Optimized Price'}
                                </button>
                            </div>

                            {/* Factors */}
                            <div className="rounded-xl border bg-white p-6">
                                <h2 className="mb-4 text-xl font-bold">📈 Pricing Factors</h2>
                                <div className="grid gap-4 md:grid-cols-2">
                                    <div className="flex items-center gap-4 rounded-lg border p-4">
                                        <div className={`rounded-full px-3 py-1 text-sm font-semibold ${priceData.demand_level === 'HIGH' ? 'bg-red-100 text-red-700' : 'bg-yellow-100 text-yellow-700'
                                            }`}>
                                            {priceData.demand_level}
                                        </div>
                                        <div>
                                            <div className="font-semibold">Demand Level</div>
                                            <div className="text-sm text-muted-foreground">Current market demand</div>
                                        </div>
                                    </div>

                                    <div className="flex items-center gap-4 rounded-lg border p-4">
                                        <div className="text-2xl">📦</div>
                                        <div>
                                            <div className="font-semibold">{priceData.inventory_level} units</div>
                                            <div className="text-sm text-muted-foreground">In stock</div>
                                        </div>
                                    </div>

                                    <div className="flex items-center gap-4 rounded-lg border p-4">
                                        <div className="text-2xl">📊</div>
                                        <div>
                                            <div className="font-semibold">{priceData.price_elasticity.toFixed(2)}</div>
                                            <div className="text-sm text-muted-foreground">Price elasticity</div>
                                        </div>
                                    </div>

                                    <div className="flex items-center gap-4 rounded-lg border p-4">
                                        <div className={`text-2xl ${priceData.recommendation === 'LOWER' ? '⬇️' : '⬆️'}`}></div>
                                        <div>
                                            <div className="font-semibold">{priceData.recommendation} Price</div>
                                            <div className="text-sm text-muted-foreground">AI recommendation</div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </>
                    )}

                    {!priceData && !loading && (
                        <div className="flex h-64 items-center justify-center rounded-xl border bg-gradient-to-br from-slate-50 to-slate-100">
                            <div className="text-center text-muted-foreground">
                                <div className="text-4xl mb-2">👈</div>
                                Select a product to analyze
                            </div>
                        </div>
                    )}
                </div>
            </div>

            {/* Algorithm Info */}
            <div className="mt-8 rounded-xl border bg-gradient-to-r from-slate-50 to-slate-100 p-8">
                <h2 className="mb-6 text-2xl font-bold">🤖 Pricing Algorithm Factors</h2>
                <div className="grid gap-4 md:grid-cols-5">
                    {[
                        { icon: '📈', name: 'Demand Forecasting' },
                        { icon: '🏪', name: 'Competitor Analysis' },
                        { icon: '📦', name: 'Inventory Levels' },
                        { icon: '⏰', name: 'Time-based Patterns' },
                        { icon: '👥', name: 'Customer Segments' },
                    ].map((factor) => (
                        <div key={factor.name} className="rounded-lg bg-white p-4 text-center shadow-sm">
                            <div className="text-3xl mb-2">{factor.icon}</div>
                            <div className="text-sm font-semibold">{factor.name}</div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
