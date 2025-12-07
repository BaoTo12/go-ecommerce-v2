'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { wishlistService, Wishlist } from '@/services/wishlistService';
import { cartService } from '@/services/cartService';

export default function WishlistPage() {
    const [wishlist, setWishlist] = useState<Wishlist | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [notification, setNotification] = useState<string | null>(null);

    const loadWishlist = async () => {
        const data = await wishlistService.getWishlist();
        setWishlist(data);
        setIsLoading(false);
    };

    useEffect(() => {
        loadWishlist();
    }, []);

    const removeFromWishlist = async (productId: string) => {
        await wishlistService.removeItem(productId);
        loadWishlist();
        setNotification('✓ Đã xóa khỏi danh sách yêu thích');
        setTimeout(() => setNotification(null), 2000);
    };

    const addToCart = async (item: Wishlist['items'][0]) => {
        await cartService.addItem(item.product);
        setNotification('🛒 Đã thêm vào giỏ hàng');
        setTimeout(() => setNotification(null), 2000);
    };

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    if (isLoading) {
        return (
            <div className="min-h-screen bg-[#f5f5f5] flex items-center justify-center">
                <div className="loading-spinner" />
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-[#f5f5f5] animate-fade-in">
            {/* Toast */}
            {notification && <div className="toast toast-success">{notification}</div>}

            <div className="container mx-auto px-4 py-6">
                <div className="flex items-center justify-between mb-4">
                    <h1 className="text-xl font-medium">Sản Phẩm Yêu Thích</h1>
                    <span className="text-sm text-gray-500">{wishlist?.items.length || 0} sản phẩm</span>
                </div>

                {!wishlist || wishlist.items.length === 0 ? (
                    <div className="bg-white rounded-sm shadow-sm p-12 text-center">
                        <div className="text-6xl mb-4 animate-float">❤️</div>
                        <p className="text-gray-500 mb-4">Bạn chưa có sản phẩm yêu thích nào</p>
                        <Link
                            href="/products"
                            className="inline-block px-8 py-3 bg-[#ee4d2d] text-white hover:opacity-90 transition-all"
                        >
                            Khám Phá Ngay
                        </Link>
                    </div>
                ) : (
                    <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-[10px]">
                        {wishlist.items.map((item, index) => (
                            <div
                                key={item.id}
                                className="product-card group animate-fade-in-up"
                                style={{ animationDelay: `${index * 50}ms` }}
                            >
                                <div className="relative aspect-square bg-gray-100 overflow-hidden">
                                    <Link href={`/products/${item.productId}`}>
                                        <Image
                                            src={item.product.thumbnail}
                                            alt={item.product.name}
                                            fill
                                            className="object-cover product-image"
                                            unoptimized
                                        />
                                    </Link>

                                    {/* Remove from wishlist */}
                                    <button
                                        onClick={() => removeFromWishlist(item.productId)}
                                        className="absolute top-2 right-2 w-8 h-8 bg-white rounded-full flex items-center justify-center shadow-md hover:bg-red-50 transition-colors group/btn"
                                    >
                                        <span className="text-red-500 group-hover/btn:scale-125 transition-transform">❤️</span>
                                    </button>

                                    {/* Discount */}
                                    {item.product.discount > 0 && (
                                        <div className="discount-badge">-{item.product.discount}%</div>
                                    )}
                                </div>

                                <div className="p-2">
                                    <Link href={`/products/${item.productId}`}>
                                        <h3 className="text-xs line-clamp-2 h-8 mb-1 hover:text-[#ee4d2d]">{item.product.name}</h3>
                                    </Link>

                                    <div className="flex items-end justify-between mb-2">
                                        <div>
                                            <span className="price-current text-sm font-medium">₫{formatPrice(item.product.price)}</span>
                                            {item.product.originalPrice > item.product.price && (
                                                <span className="price-original block">₫{formatPrice(item.product.originalPrice)}</span>
                                            )}
                                        </div>
                                    </div>

                                    <button
                                        onClick={() => addToCart(item)}
                                        className="w-full py-2 border border-[#ee4d2d] text-[#ee4d2d] text-xs hover:bg-[#fef6f5] transition-colors"
                                    >
                                        🛒 Thêm vào giỏ
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
