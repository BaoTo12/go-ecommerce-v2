'use client';

import React, { useState, useRef, useEffect, useCallback } from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { Product, productService } from '@/services/productService';

const ITEM_HEIGHT = 200;
const OVERSCAN = 5;

export default function VirtualScroll() {
    const [products, setProducts] = useState<Product[]>([]);
    const [visibleRange, setVisibleRange] = useState({ start: 0, end: 20 });
    const [totalHeight, setTotalHeight] = useState(0);
    const containerRef = useRef<HTMLDivElement>(null);
    const [isLoading, setIsLoading] = useState(true);

    // Load all products (in real app, this would be paginated)
    useEffect(() => {
        const loadProducts = async () => {
            setIsLoading(true);
            const response = await productService.getProducts({ limit: 100 });
            const productList = Array.isArray(response) ? response : response.products || [];

            // Duplicate to simulate large dataset
            const largeList = [...productList, ...productList, ...productList, ...productList, ...productList];
            setProducts(largeList.map((p, i) => ({ ...p, id: `${p.id}-${i}` })));
            setTotalHeight(largeList.length * ITEM_HEIGHT);
            setIsLoading(false);
        };
        loadProducts();
    }, [ITEM_HEIGHT]);

    // Calculate visible items on scroll
    const handleScroll = useCallback(() => {
        const container = containerRef.current;
        if (!container) return;

        const scrollTop = container.scrollTop;
        const viewportHeight = container.clientHeight;

        const start = Math.max(0, Math.floor(scrollTop / ITEM_HEIGHT) - OVERSCAN);
        const end = Math.min(
            products.length,
            Math.ceil((scrollTop + viewportHeight) / ITEM_HEIGHT) + OVERSCAN
        );

        setVisibleRange({ start, end });
    }, [products.length, ITEM_HEIGHT, OVERSCAN]);

    useEffect(() => {
        const container = containerRef.current;
        if (!container) return;

        container.addEventListener('scroll', handleScroll);
        handleScroll(); // Initial calculation

        return () => container.removeEventListener('scroll', handleScroll);
    }, [handleScroll]);

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    const visibleProducts = products.slice(visibleRange.start, visibleRange.end);

    return (
        <div className="min-h-screen bg-gray-100 dark:bg-gray-900">
            {/* Header */}
            <div className="bg-white dark:bg-gray-800 p-4 shadow-sm sticky top-0 z-10">
                <h1 className="text-xl font-bold dark:text-white">
                    ⚡ Virtual Scroll Demo ({products.length} sản phẩm)
                </h1>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Hiển thị: {visibleRange.start + 1} - {Math.min(visibleRange.end, products.length)} / {products.length}
                </p>
            </div>

            {isLoading ? (
                <div className="p-8 text-center">
                    <div className="loading-spinner mx-auto mb-4" />
                    <p className="text-gray-500">Đang tải...</p>
                </div>
            ) : (
                <div
                    ref={containerRef}
                    className="h-[calc(100vh-80px)] overflow-y-auto"
                    style={{ willChange: 'scroll-position' }}
                >
                    {/* Spacer for total height */}
                    <div style={{ height: totalHeight, position: 'relative' }}>
                        {/* Visible items */}
                        {visibleProducts.map((product, index) => {
                            const actualIndex = visibleRange.start + index;
                            const top = actualIndex * ITEM_HEIGHT;

                            return (
                                <div
                                    key={product.id}
                                    className="absolute left-0 right-0 px-4"
                                    style={{ top, height: ITEM_HEIGHT }}
                                >
                                    <Link
                                        href={`/products/${product.id.split('-')[0]}`}
                                        className="flex gap-4 p-4 bg-white dark:bg-gray-800 rounded-lg shadow-sm h-[calc(100%-8px)] hover:shadow-md transition-shadow"
                                    >
                                        <div className="w-32 h-32 bg-gray-100 dark:bg-gray-700 rounded-lg overflow-hidden relative flex-shrink-0">
                                            <Image
                                                src={product.thumbnail}
                                                alt={product.name}
                                                fill
                                                className="object-cover"
                                                unoptimized
                                                loading="lazy"
                                            />
                                        </div>
                                        <div className="flex-1 flex flex-col justify-between py-1">
                                            <div>
                                                <h3 className="font-medium line-clamp-2 dark:text-white">{product.name}</h3>
                                                <div className="flex items-center gap-2 mt-1">
                                                    <span className="text-yellow-500">⭐ {product.rating}</span>
                                                    <span className="text-gray-400 text-sm">|</span>
                                                    <span className="text-gray-500 text-sm">{product.soldDisplay} đã bán</span>
                                                </div>
                                            </div>
                                            <div className="flex items-center justify-between">
                                                <div>
                                                    <span className="text-xl font-bold text-[#ee4d2d]">₫{formatPrice(product.price)}</span>
                                                    {product.discount > 0 && (
                                                        <span className="ml-2 text-sm text-gray-400 line-through">
                                                            ₫{formatPrice(product.originalPrice)}
                                                        </span>
                                                    )}
                                                </div>
                                                {product.discount > 0 && (
                                                    <span className="px-2 py-1 bg-red-100 text-red-600 text-xs rounded">
                                                        -{product.discount}%
                                                    </span>
                                                )}
                                            </div>
                                        </div>
                                    </Link>
                                </div>
                            );
                        })}
                    </div>
                </div>
            )}

            {/* Performance Stats */}
            <div className="fixed bottom-4 right-4 bg-black/80 text-white text-xs p-3 rounded-lg">
                <div>Total items: {products.length}</div>
                <div>Rendered: {visibleProducts.length}</div>
                <div>DOM nodes saved: {products.length - visibleProducts.length}</div>
            </div>
        </div>
    );
}
