'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { productService, Product } from '@/services/productService';

export default function ProductComparison() {
    const [products, setProducts] = useState<Product[]>([]);
    const [compareList, setCompareList] = useState<Product[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        productService.getProducts({ limit: 12 }).then(response => {
            const productList = Array.isArray(response) ? response : response.products || [];
            setProducts(productList);
            setIsLoading(false);
        });
    }, []);

    const addToCompare = (product: Product) => {
        if (compareList.length >= 4) {
            alert('Chỉ so sánh tối đa 4 sản phẩm');
            return;
        }
        if (compareList.find(p => p.id === product.id)) {
            return;
        }
        setCompareList([...compareList, product]);
    };

    const removeFromCompare = (productId: string) => {
        setCompareList(compareList.filter(p => p.id !== productId));
    };

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    return (
        <div className="min-h-screen bg-[#f5f5f5] dark:bg-gray-900">
            <div className="container mx-auto px-4 py-6">
                <h1 className="text-xl font-bold mb-6 dark:text-white">So Sánh Sản Phẩm</h1>

                {/* Compare Bar */}
                {compareList.length > 0 && (
                    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-4 mb-6 sticky top-[120px] z-40">
                        <div className="flex items-center gap-4">
                            <span className="font-medium dark:text-white">Đang so sánh ({compareList.length}/4):</span>
                            <div className="flex gap-2 flex-1 overflow-x-auto">
                                {compareList.map(product => (
                                    <div
                                        key={product.id}
                                        className="flex items-center gap-2 bg-gray-100 dark:bg-gray-700 rounded-full px-3 py-1"
                                    >
                                        <span className="text-sm truncate max-w-[100px] dark:text-white">{product.name}</span>
                                        <button
                                            onClick={() => removeFromCompare(product.id)}
                                            className="text-gray-400 hover:text-red-500"
                                        >
                                            ✕
                                        </button>
                                    </div>
                                ))}
                            </div>
                            {compareList.length >= 2 && (
                                <button
                                    onClick={() => { }}
                                    className="px-4 py-2 bg-[#ee4d2d] text-white rounded-sm hover:opacity-90"
                                >
                                    So sánh ngay
                                </button>
                            )}
                        </div>
                    </div>
                )}

                {/* Product Grid */}
                <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-3 mb-6">
                    {isLoading ? (
                        [...Array(12)].map((_, i) => (
                            <div key={i} className="bg-white dark:bg-gray-800 rounded-sm p-3 animate-pulse">
                                <div className="aspect-square bg-gray-200 dark:bg-gray-700 mb-3" />
                                <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded mb-2" />
                                <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-2/3" />
                            </div>
                        ))
                    ) : (
                        products.map(product => (
                            <div
                                key={product.id}
                                className="bg-white dark:bg-gray-800 rounded-sm shadow-sm overflow-hidden group"
                            >
                                <div className="relative aspect-square bg-gray-100 dark:bg-gray-700">
                                    <Link href={`/products/${product.id}`}>
                                        <Image
                                            src={product.thumbnail}
                                            alt={product.name}
                                            fill
                                            className="object-cover"
                                            unoptimized
                                        />
                                    </Link>
                                    <button
                                        onClick={() => addToCompare(product)}
                                        className={`absolute bottom-2 right-2 w-8 h-8 rounded-full flex items-center justify-center shadow-md transition-all ${compareList.find(p => p.id === product.id)
                                                ? 'bg-[#ee4d2d] text-white'
                                                : 'bg-white text-gray-600 opacity-0 group-hover:opacity-100'
                                            }`}
                                    >
                                        ⚖️
                                    </button>
                                </div>
                                <div className="p-2">
                                    <h3 className="text-xs line-clamp-2 dark:text-white">{product.name}</h3>
                                    <div className="text-sm font-medium text-[#ee4d2d] mt-1">
                                        ₫{formatPrice(product.price)}
                                    </div>
                                </div>
                            </div>
                        ))
                    )}
                </div>

                {/* Comparison Table */}
                {compareList.length >= 2 && (
                    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden animate-fade-in-up">
                        <div className="p-4 border-b dark:border-gray-700">
                            <h2 className="font-bold dark:text-white">Bảng So Sánh</h2>
                        </div>
                        <div className="overflow-x-auto">
                            <table className="w-full">
                                <thead>
                                    <tr className="bg-gray-50 dark:bg-gray-700">
                                        <th className="p-4 text-left dark:text-white">Thông số</th>
                                        {compareList.map(product => (
                                            <th key={product.id} className="p-4 text-center min-w-[200px]">
                                                <Image
                                                    src={product.thumbnail}
                                                    alt={product.name}
                                                    width={100}
                                                    height={100}
                                                    className="mx-auto mb-2"
                                                    unoptimized
                                                />
                                                <div className="text-sm font-medium dark:text-white">{product.name}</div>
                                            </th>
                                        ))}
                                    </tr>
                                </thead>
                                <tbody className="divide-y dark:divide-gray-700">
                                    <tr>
                                        <td className="p-4 font-medium dark:text-white">Giá</td>
                                        {compareList.map(product => (
                                            <td key={product.id} className="p-4 text-center text-[#ee4d2d] font-bold">
                                                ₫{formatPrice(product.price)}
                                            </td>
                                        ))}
                                    </tr>
                                    <tr>
                                        <td className="p-4 font-medium dark:text-white">Đánh giá</td>
                                        {compareList.map(product => (
                                            <td key={product.id} className="p-4 text-center dark:text-white">
                                                ⭐ {product.rating} ({product.soldCount} đã bán)
                                            </td>
                                        ))}
                                    </tr>
                                    <tr>
                                        <td className="p-4 font-medium dark:text-white">Shop</td>
                                        {compareList.map(product => (
                                            <td key={product.id} className="p-4 text-center dark:text-white">
                                                {product.shop?.name || 'N/A'}
                                            </td>
                                        ))}
                                    </tr>
                                    <tr>
                                        <td className="p-4 font-medium dark:text-white">Giảm giá</td>
                                        {compareList.map(product => (
                                            <td key={product.id} className="p-4 text-center">
                                                <span className="text-[#ee4d2d]">-{product.discount}%</span>
                                            </td>
                                        ))}
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
