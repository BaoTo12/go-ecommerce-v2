'use client';

import React, { useState, useEffect, useCallback } from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { productService, Product } from '@/services/productService';

interface InfiniteProductGridProps {
    initialProducts?: Product[];
    category?: string;
    search?: string;
}

export default function InfiniteProductGrid({
    initialProducts = [],
    category,
    search
}: InfiniteProductGridProps) {
    const [products, setProducts] = useState<Product[]>(initialProducts);
    const [page, setPage] = useState(1);
    const [hasMore, setHasMore] = useState(true);
    const [isLoading, setIsLoading] = useState(false);

    const loadMore = useCallback(async () => {
        if (isLoading || !hasMore) return;

        setIsLoading(true);
        try {
            const response = await productService.getProducts({
                category,
                search,
                page: page + 1,
                limit: 20
            });

            const newProducts = Array.isArray(response) ? response : response.products || [];

            if (newProducts.length === 0) {
                setHasMore(false);
            } else {
                setProducts(prev => [...prev, ...newProducts]);
                setPage(prev => prev + 1);
            }
        } catch (error) {
            console.error('Error loading more products:', error);
        } finally {
            setIsLoading(false);
        }
    }, [page, category, search, isLoading, hasMore]);

    // Intersection Observer for infinite scroll
    useEffect(() => {
        const observer = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting && hasMore && !isLoading) {
                    loadMore();
                }
            },
            { threshold: 0.1 }
        );

        const sentinel = document.getElementById('scroll-sentinel');
        if (sentinel) observer.observe(sentinel);

        return () => observer.disconnect();
    }, [loadMore, hasMore, isLoading]);

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    return (
        <div>
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-[10px]">
                {products.map((product, index) => (
                    <Link
                        key={`${product.id}-${index}`}
                        href={`/products/${product.id}`}
                        className="product-card group animate-fade-in-up"
                        style={{ animationDelay: `${(index % 20) * 30}ms` }}
                    >
                        <div className="relative aspect-square bg-gray-100 dark:bg-gray-700 overflow-hidden">
                            <Image
                                src={product.thumbnail}
                                alt={product.name}
                                fill
                                className="object-cover product-image"
                                unoptimized
                                loading="lazy"
                            />
                            {product.discount > 0 && (
                                <div className="discount-badge">-{product.discount}%</div>
                            )}
                        </div>
                        <div className="p-2">
                            <h3 className="text-xs line-clamp-2 h-8 mb-1 dark:text-white">{product.name}</h3>
                            <div className="flex items-end justify-between">
                                <div>
                                    <span className="price-current text-sm">₫{formatPrice(product.price)}</span>
                                    {product.originalPrice > product.price && (
                                        <span className="price-original block">₫{formatPrice(product.originalPrice)}</span>
                                    )}
                                </div>
                                <span className="text-[10px] text-gray-400 dark:text-gray-500">
                                    Đã bán {product.sold > 1000 ? `${(product.sold / 1000).toFixed(1)}k` : product.sold}
                                </span>
                            </div>
                        </div>
                    </Link>
                ))}
            </div>

            {/* Loading Skeleton */}
            {isLoading && (
                <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-[10px] mt-4">
                    {[...Array(12)].map((_, i) => (
                        <div key={i} className="bg-white dark:bg-gray-800 rounded-sm shadow-sm animate-pulse">
                            <div className="aspect-square bg-gray-200 dark:bg-gray-700" />
                            <div className="p-2">
                                <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded mb-2" />
                                <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-2/3" />
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* Scroll Sentinel */}
            <div id="scroll-sentinel" className="h-10" />

            {/* No more products */}
            {!hasMore && products.length > 0 && (
                <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                    <p>Bạn đã xem hết tất cả sản phẩm 🎉</p>
                </div>
            )}
        </div>
    );
}
