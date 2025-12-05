'use client';

import React, { useState } from 'react';

interface Product {
    id: string;
    name: string;
    category: string;
    basePrice: number;
    thumbnail: string;
}

interface PriceAnalysis {
    productId: string;
    productName: string;
    basePrice: number;
    optimizedPrice: number;
    demandLevel: string;
    competitorPrices: number[];
    revenueIncrease: number;
    recommendation: string;
}

export default function DynamicPricingPage() {
    const [selectedProduct, setSelectedProduct] = useState<string | null>(null);
    const [analysis, setAnalysis] = useState<PriceAnalysis | null>(null);
    const [loading, setLoading] = useState(false);
    const [optimizing, setOptimizing] = useState(false);

    const products: Product[] = [
        { id: 'p1', name: 'iPhone 15 Pro Max', category: 'Điện thoại', basePrice: 34990000, thumbnail: '📱' },
        { id: 'p2', name: 'MacBook Air M3', category: 'Laptop', basePrice: 27990000, thumbnail: '💻' },
        { id: 'p3', name: 'Sony WH-1000XM5', category: 'Tai nghe', basePrice: 8990000, thumbnail: '🎧' },
        { id: 'p4', name: 'Nike Air Max', category: 'Giày dép', basePrice: 4590000, thumbnail: '👟' },
    ];

    const analyzePrice = (productId: string) => {
        setSelectedProduct(productId);
        setLoading(true);

        setTimeout(() => {
            const product = products.find(p => p.id === productId);
            if (!product) return;

            const demandMultiplier = 0.9 + Math.random() * 0.2;
            const optimizedPrice = Math.round(product.basePrice * demandMultiplier);

            setAnalysis({
                productId,
                productName: product.name,
                basePrice: product.basePrice,
                optimizedPrice,
                demandLevel: Math.random() > 0.5 ? 'CAO' : 'TRUNG BÌNH',
                competitorPrices: [
                    Math.round(product.basePrice * (0.92 + Math.random() * 0.1)),
                    Math.round(product.basePrice * (0.95 + Math.random() * 0.12)),
                    Math.round(product.basePrice * (0.88 + Math.random() * 0.15)),
                ],
                revenueIncrease: Math.round(Math.random() * 15 * 10) / 10,
                recommendation: optimizedPrice < product.basePrice ? 'GIẢM' : 'TĂNG',
            });
            setLoading(false);
        }, 1000);
    };

    const applyOptimizedPrice = () => {
        setOptimizing(true);
        setTimeout(() => {
            setOptimizing(false);
            alert('✅ Đã áp dụng giá tối ưu!');
        }, 1500);
    };

    const formatPrice = (price: number) => {
        return new Intl.NumberFormat('vi-VN').format(price) + '₫';
    };

    return (
        <div className="min-h-screen bg-[#F5F5F5]">
            {/* Header */}
            <div className="bg-gradient-to-r from-emerald-600 to-teal-600 text-white px-6 py-6">
                <h1 className="text-2xl font-bold">💹 Dynamic Pricing Engine</h1>
                <p className="text-emerald-100 text-sm">Tối ưu giá với AI để tăng doanh thu</p>
            </div>

            <div className="p-6">
                <div className="grid lg:grid-cols-3 gap-6">
                    {/* Product Selection */}
                    <div className="bg-white rounded p-6">
                        <h2 className="font-bold text-lg mb-4">📦 Chọn sản phẩm</h2>
                        <div className="space-y-3">
                            {products.map(product => (
                                <button
                                    key={product.id}
                                    onClick={() => analyzePrice(product.id)}
                                    className={`w-full p-4 rounded border text-left transition-all ${selectedProduct === product.id
                                            ? 'border-emerald-500 bg-emerald-50'
                                            : 'border-gray-200 hover:border-gray-300'
                                        }`}
                                >
                                    <div className="flex items-center gap-3">
                                        <span className="text-3xl">{product.thumbnail}</span>
                                        <div>
                                            <div className="font-semibold text-sm">{product.name}</div>
                                            <div className="text-xs text-gray-500">{product.category}</div>
                                            <div className="text-sm text-emerald-600 font-bold mt-1">
                                                {formatPrice(product.basePrice)}
                                            </div>
                                        </div>
                                    </div>
                                </button>
                            ))}
                        </div>
                    </div>

                    {/* Analysis Results */}
                    <div className="lg:col-span-2 space-y-6">
                        {loading && (
                            <div className="bg-white rounded p-12 text-center">
                                <div className="text-3xl animate-bounce">📊</div>
                                <p className="mt-2 text-gray-500">Đang phân tích...</p>
                            </div>
                        )}

                        {analysis && !loading && (
                            <>
                                {/* Price Comparison */}
                                <div className="bg-white rounded p-6">
                                    <h2 className="font-bold text-lg mb-4">💰 Phân tích giá</h2>

                                    <div className="grid md:grid-cols-3 gap-4 mb-6">
                                        <div className="p-4 bg-gray-100 rounded text-center">
                                            <div className="text-sm text-gray-500">Giá gốc</div>
                                            <div className="text-xl font-bold">{formatPrice(analysis.basePrice)}</div>
                                        </div>
                                        <div className="p-4 bg-emerald-100 rounded text-center">
                                            <div className="text-sm text-emerald-700">Giá tối ưu</div>
                                            <div className="text-xl font-bold text-emerald-600">
                                                {formatPrice(analysis.optimizedPrice)}
                                            </div>
                                            <div className={`text-sm font-semibold ${analysis.optimizedPrice < analysis.basePrice ? 'text-red-500' : 'text-green-500'
                                                }`}>
                                                {analysis.optimizedPrice < analysis.basePrice ? '↓' : '↑'}
                                                {Math.abs(((analysis.optimizedPrice - analysis.basePrice) / analysis.basePrice) * 100).toFixed(1)}%
                                            </div>
                                        </div>
                                        <div className="p-4 bg-blue-100 rounded text-center">
                                            <div className="text-sm text-blue-700">Tăng doanh thu dự kiến</div>
                                            <div className="text-xl font-bold text-blue-600">+{analysis.revenueIncrease}%</div>
                                        </div>
                                    </div>

                                    {/* Competitor Prices */}
                                    <div className="mb-6">
                                        <h3 className="font-semibold text-sm mb-2">🏪 Giá đối thủ</h3>
                                        <div className="flex gap-3">
                                            {analysis.competitorPrices.map((price, i) => (
                                                <div key={i} className="px-4 py-2 bg-gray-50 rounded text-center">
                                                    <div className="text-xs text-gray-400">Đối thủ {i + 1}</div>
                                                    <div className="font-semibold">{formatPrice(price)}</div>
                                                </div>
                                            ))}
                                        </div>
                                    </div>

                                    <button
                                        onClick={applyOptimizedPrice}
                                        disabled={optimizing}
                                        className={`w-full py-3 rounded font-bold text-white ${optimizing
                                                ? 'bg-gray-400'
                                                : 'bg-gradient-to-r from-emerald-600 to-teal-600 hover:opacity-90'
                                            }`}
                                    >
                                        {optimizing ? '⏳ Đang áp dụng...' : '🚀 Áp dụng giá tối ưu'}
                                    </button>
                                </div>

                                {/* Pricing Factors */}
                                <div className="bg-white rounded p-6">
                                    <h2 className="font-bold text-lg mb-4">📈 Yếu tố định giá</h2>
                                    <div className="grid md:grid-cols-2 gap-4">
                                        <div className="flex items-center gap-3 p-3 bg-gray-50 rounded">
                                            <span className={`px-3 py-1 rounded text-sm font-bold ${analysis.demandLevel === 'CAO' ? 'bg-red-100 text-red-700' : 'bg-yellow-100 text-yellow-700'
                                                }`}>
                                                {analysis.demandLevel}
                                            </span>
                                            <div>
                                                <div className="font-semibold text-sm">Mức độ nhu cầu</div>
                                                <div className="text-xs text-gray-500">Dựa trên xu hướng thị trường</div>
                                            </div>
                                        </div>

                                        <div className="flex items-center gap-3 p-3 bg-gray-50 rounded">
                                            <span className="text-2xl">{analysis.recommendation === 'GIẢM' ? '⬇️' : '⬆️'}</span>
                                            <div>
                                                <div className="font-semibold text-sm">{analysis.recommendation} giá</div>
                                                <div className="text-xs text-gray-500">Khuyến nghị AI</div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </>
                        )}

                        {!analysis && !loading && (
                            <div className="bg-white rounded p-12 text-center">
                                <div className="text-4xl mb-2">👈</div>
                                <p className="text-gray-500">Chọn sản phẩm để phân tích</p>
                            </div>
                        )}
                    </div>
                </div>

                {/* Algorithm Info */}
                <div className="mt-6 bg-white rounded p-6">
                    <h2 className="font-bold text-lg mb-4">🤖 Thuật toán định giá</h2>
                    <div className="grid grid-cols-2 md:grid-cols-5 gap-4 text-center">
                        {[
                            { icon: '📈', name: 'Dự báo nhu cầu' },
                            { icon: '🏪', name: 'Phân tích đối thủ' },
                            { icon: '📦', name: 'Mức tồn kho' },
                            { icon: '⏰', name: 'Mẫu theo thời gian' },
                            { icon: '👥', name: 'Phân khúc KH' },
                        ].map(f => (
                            <div key={f.name} className="p-4 bg-gray-50 rounded">
                                <div className="text-3xl mb-2">{f.icon}</div>
                                <div className="text-sm font-semibold">{f.name}</div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}
