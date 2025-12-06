'use client';

import React, { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import Image from 'next/image';
import { productService, Product } from '@/services/productService';

export default function ProductDetailPage() {
    const params = useParams();
    const router = useRouter();
    const productId = params.id as string;

    const [product, setProduct] = useState<Product | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [quantity, setQuantity] = useState(1);
    const [selectedVariants, setSelectedVariants] = useState<Record<string, number>>({});
    const [selectedImage, setSelectedImage] = useState(0);
    const [notification, setNotification] = useState<string | null>(null);
    const [isAddingToCart, setIsAddingToCart] = useState(false);
    const [isBuying, setIsBuying] = useState(false);
    const [showHearts, setShowHearts] = useState(false);

    // Load product data
    useEffect(() => {
        const loadProduct = async () => {
            setIsLoading(true);
            try {
                const data = await productService.getProduct(productId);
                if (data) {
                    setProduct(data);
                    // Initialize variant selections
                    const initVariants: Record<string, number> = {};
                    data.variants?.forEach(v => {
                        initVariants[v.id] = 0;
                    });
                    setSelectedVariants(initVariants);
                }
            } catch (error) {
                console.error('Failed to load product:', error);
            } finally {
                setIsLoading(false);
            }
        };
        loadProduct();
    }, [productId]);

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    const addToCart = () => {
        setIsAddingToCart(true);
        setShowHearts(true);

        setTimeout(() => {
            setIsAddingToCart(false);
            setNotification('🛒 Đã thêm sản phẩm vào Giỏ hàng');
            setTimeout(() => {
                setNotification(null);
                setShowHearts(false);
            }, 2500);
        }, 800);
    };

    const buyNow = () => {
        setIsBuying(true);
        setTimeout(() => {
            router.push('/checkout');
        }, 500);
    };

    if (isLoading) {
        return (
            <div className="min-h-screen bg-[#f5f5f5]">
                <div className="container mx-auto px-4 py-4">
                    <div className="bg-white rounded-sm shadow-sm p-6">
                        <div className="grid md:grid-cols-5 gap-6">
                            <div className="md:col-span-2">
                                <div className="aspect-square bg-gray-200 rounded-sm animate-pulse" />
                            </div>
                            <div className="md:col-span-3 space-y-4">
                                <div className="h-6 bg-gray-200 rounded animate-pulse" />
                                <div className="h-4 bg-gray-200 rounded w-1/2 animate-pulse" />
                                <div className="h-12 bg-gray-200 rounded animate-pulse" />
                                <div className="h-4 bg-gray-200 rounded w-3/4 animate-pulse" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }

    if (!product) {
        return (
            <div className="min-h-screen bg-[#f5f5f5] flex items-center justify-center">
                <div className="text-center">
                    <div className="text-6xl mb-4">😕</div>
                    <h2 className="text-xl font-medium mb-2">Sản phẩm không tồn tại</h2>
                    <Link href="/products" className="text-[#ee4d2d] hover:underline">
                        ← Quay lại danh sách sản phẩm
                    </Link>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-[#f5f5f5] animate-fade-in">
            {/* Toast */}
            {notification && <div className="toast toast-success">{notification}</div>}

            {/* Floating hearts animation */}
            {showHearts && (
                <div className="fixed inset-0 pointer-events-none z-50 overflow-hidden">
                    {[...Array(12)].map((_, i) => (
                        <span
                            key={i}
                            className="absolute text-2xl animate-float"
                            style={{
                                left: `${20 + Math.random() * 60}%`,
                                top: `${30 + Math.random() * 40}%`,
                                animationDelay: `${i * 100}ms`,
                                animationDuration: `${1 + Math.random()}s`,
                            }}
                        >
                            ❤️
                        </span>
                    ))}
                </div>
            )}

            {/* Breadcrumb */}
            <div className="bg-white animate-fade-in-down">
                <div className="container mx-auto px-4 py-2">
                    <div className="flex items-center gap-2 text-xs text-gray-500">
                        <Link href="/" className="hover:text-[#ee4d2d] transition-colors">Shopee</Link>
                        <span>›</span>
                        <Link href="/products" className="hover:text-[#ee4d2d] transition-colors">{product.category}</Link>
                        <span>›</span>
                        <span className="text-gray-700 line-clamp-1">{product.name}</span>
                    </div>
                </div>
            </div>

            <div className="container mx-auto px-4 py-4">
                {/* Product Info */}
                <div className="bg-white rounded-sm shadow-sm mb-4 animate-fade-in-up">
                    <div className="grid md:grid-cols-5 gap-6 p-4">
                        {/* Images */}
                        <div className="md:col-span-2 animate-fade-in-left">
                            <div className="relative aspect-square bg-gray-50 rounded-sm overflow-hidden mb-2">
                                <Image
                                    src={product.images[selectedImage] || product.thumbnail}
                                    alt={product.name}
                                    fill
                                    className="object-cover transition-all duration-500"
                                    unoptimized
                                />
                            </div>
                            <div className="flex gap-2 overflow-x-auto">
                                {product.images.map((img, i) => (
                                    <button
                                        key={i}
                                        onClick={() => setSelectedImage(i)}
                                        className={`relative w-16 h-16 flex-shrink-0 rounded-sm overflow-hidden border-2 transition-all duration-300 hover:scale-105
                               ${selectedImage === i ? 'border-[#ee4d2d]' : 'border-transparent hover:border-gray-300'}`}
                                    >
                                        <Image
                                            src={img}
                                            alt={`${product.name} ${i + 1}`}
                                            fill
                                            className="object-cover"
                                            unoptimized
                                        />
                                    </button>
                                ))}
                            </div>
                        </div>

                        {/* Details */}
                        <div className="md:col-span-3 animate-fade-in-right">
                            <div className="flex items-start gap-2 mb-2">
                                {product.isOfficial && (
                                    <span className="bg-[#ee4d2d] text-white text-[10px] px-1 py-0.5 animate-pulse">Mall</span>
                                )}
                                <h1 className="text-lg leading-tight flex-1">{product.name}</h1>
                            </div>

                            {/* Rating & Sold */}
                            <div className="flex items-center gap-4 text-sm py-3 border-b">
                                <div className="flex items-center gap-1">
                                    <span className="text-[#ee4d2d] font-medium border-b border-[#ee4d2d]">{product.rating}</span>
                                    <span className="star-rating">{'★'.repeat(Math.floor(product.rating))}{'☆'.repeat(5 - Math.floor(product.rating))}</span>
                                </div>
                                <span className="text-gray-300">|</span>
                                <div className="hover:scale-105 transition-transform cursor-pointer">
                                    <span className="font-medium border-b border-gray-600">{(product.reviews / 1000).toFixed(1)}k</span>
                                    <span className="text-gray-500 ml-1">Đánh Giá</span>
                                </div>
                                <span className="text-gray-300">|</span>
                                <div className="hover:scale-105 transition-transform cursor-pointer">
                                    <span className="font-medium">{product.soldDisplay}</span>
                                    <span className="text-gray-500 ml-1">Đã Bán</span>
                                </div>
                            </div>

                            {/* Price */}
                            <div className="bg-[#fafafa] p-4 my-3 animate-fade-in">
                                <div className="flex items-center gap-3">
                                    <span className="text-gray-400 line-through text-sm">₫{formatPrice(product.originalPrice)}</span>
                                    <span className="text-[#ee4d2d] text-3xl font-medium animate-pulse-glow rounded px-2">₫{formatPrice(product.price)}</span>
                                    {product.discount > 0 && (
                                        <span className="bg-[#ee4d2d] text-white text-xs px-2 py-0.5 rounded-sm animate-bounce">{product.discount}% GIẢM</span>
                                    )}
                                </div>
                            </div>

                            {/* Vouchers */}
                            <div className="flex items-center gap-4 py-3 text-sm">
                                <span className="text-gray-500 w-24">Mã Giảm Giá</span>
                                <div className="flex gap-2 flex-wrap">
                                    <span className="bg-[#fef6f5] text-[#ee4d2d] border border-[#ee4d2d] px-2 py-0.5 text-xs cursor-pointer hover:bg-[#ffeee8] transition-colors hover-shine">Giảm ₫50k</span>
                                    <span className="bg-[#fef6f5] text-[#ee4d2d] border border-[#ee4d2d] px-2 py-0.5 text-xs cursor-pointer hover:bg-[#ffeee8] transition-colors hover-shine">Giảm 10%</span>
                                    <span className="bg-[#fef6f5] text-[#ee4d2d] border border-[#ee4d2d] px-2 py-0.5 text-xs cursor-pointer hover:bg-[#ffeee8] transition-colors hover-shine">Freeship</span>
                                </div>
                            </div>

                            {/* Shipping */}
                            <div className="flex items-center gap-4 py-3 text-sm">
                                <span className="text-gray-500 w-24">Vận Chuyển</span>
                                <div className="flex items-center gap-2">
                                    <svg className="w-5 h-5 text-[#00bfa5] animate-float" fill="currentColor" viewBox="0 0 24 24">
                                        <path d="M20 8h-3V4H3c-1.1 0-2 .9-2 2v11h2c0 1.66 1.34 3 3 3s3-1.34 3-3h6c0 1.66 1.34 3 3 3s3-1.34 3-3h2v-5l-3-4z" />
                                    </svg>
                                    <span className={product.freeShip ? 'text-[#00bfa5] font-medium' : ''}>
                                        {product.freeShip ? 'Miễn phí vận chuyển' : 'Phí vận chuyển: ₫25,000'}
                                    </span>
                                </div>
                            </div>

                            {/* Variants */}
                            {product.variants?.map(variant => (
                                <div key={variant.id} className="flex items-start gap-4 py-3 text-sm">
                                    <span className="text-gray-500 w-24 pt-2">{variant.name}</span>
                                    <div className="flex flex-wrap gap-2">
                                        {variant.options.map((option, i) => (
                                            <button
                                                key={option}
                                                onClick={() => setSelectedVariants(prev => ({ ...prev, [variant.id]: i }))}
                                                className={`px-4 py-2 border rounded-sm transition-all duration-200 hover-shrink ${selectedVariants[variant.id] === i
                                                        ? 'border-[#ee4d2d] text-[#ee4d2d] bg-[#fef6f5] scale-105'
                                                        : 'border-gray-300 hover:border-[#ee4d2d] hover:text-[#ee4d2d]'
                                                    }`}
                                            >
                                                {option}
                                            </button>
                                        ))}
                                    </div>
                                </div>
                            ))}

                            {/* Quantity */}
                            <div className="flex items-center gap-4 py-3 text-sm">
                                <span className="text-gray-500 w-24">Số Lượng</span>
                                <div className="flex items-center">
                                    <button
                                        onClick={() => setQuantity(Math.max(1, quantity - 1))}
                                        className="w-8 h-8 border flex items-center justify-center text-lg hover:bg-gray-50 transition-colors hover-shrink"
                                    >
                                        −
                                    </button>
                                    <input
                                        type="number"
                                        value={quantity}
                                        onChange={(e) => setQuantity(Math.max(1, parseInt(e.target.value) || 1))}
                                        className="w-14 h-8 border-y text-center text-sm outline-none"
                                    />
                                    <button
                                        onClick={() => setQuantity(quantity + 1)}
                                        className="w-8 h-8 border flex items-center justify-center text-lg hover:bg-gray-50 transition-colors hover-shrink"
                                    >
                                        +
                                    </button>
                                </div>
                                <span className="text-gray-400">{product.stock} sản phẩm có sẵn</span>
                            </div>

                            {/* Actions */}
                            <div className="flex gap-4 pt-4">
                                <button
                                    onClick={addToCart}
                                    disabled={isAddingToCart}
                                    className={`flex-1 py-3 border border-[#ee4d2d] text-[#ee4d2d] bg-[#fef6f5] 
                             hover:bg-[#ffeee8] transition-all flex items-center justify-center gap-2
                             hover-shrink ripple ${isAddingToCart ? 'animate-shake' : ''}`}
                                >
                                    {isAddingToCart ? (
                                        <>
                                            <span className="loading-spinner" />
                                            Đang thêm...
                                        </>
                                    ) : (
                                        <>
                                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
                                            </svg>
                                            Thêm Vào Giỏ Hàng
                                        </>
                                    )}
                                </button>
                                <button
                                    onClick={buyNow}
                                    disabled={isBuying}
                                    className={`flex-1 py-3 bg-[#ee4d2d] text-white hover:opacity-90 transition-all 
                             hover-shrink ripple ${isBuying ? 'animate-pulse' : ''}`}
                                >
                                    {isBuying ? 'Đang chuyển...' : 'Mua Ngay'}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Shop Info */}
                <div className="bg-white rounded-sm shadow-sm mb-4 p-4 animate-fade-in-up" style={{ animationDelay: '100ms' }}>
                    <div className="flex items-center gap-4">
                        <div className="relative w-16 h-16 rounded-full overflow-hidden hover:scale-110 transition-transform cursor-pointer">
                            <Image
                                src={product.shop.avatar}
                                alt={product.shop.name}
                                fill
                                className="object-cover"
                                unoptimized
                            />
                        </div>
                        <div className="flex-1">
                            <h3 className="font-medium">{product.shop.name}</h3>
                            <p className="text-xs text-gray-500 flex items-center gap-1">
                                <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
                                Online • {product.shop.location}
                            </p>
                        </div>
                        <button className="px-4 py-1.5 border border-[#ee4d2d] text-[#ee4d2d] text-sm hover:bg-[#fef6f5] transition-colors hover-shrink">
                            Xem Shop
                        </button>
                        <button className="px-4 py-1.5 border border-gray-300 text-gray-600 text-sm hover:border-[#ee4d2d] hover:text-[#ee4d2d] transition-colors hover-shrink">
                            💬 Chat Ngay
                        </button>
                    </div>
                    <div className="grid grid-cols-3 gap-4 mt-4 pt-4 border-t text-sm">
                        <div className="hover:bg-gray-50 p-2 rounded transition-colors cursor-pointer">
                            <span className="text-gray-500">Đánh Giá: </span>
                            <span className="text-[#ee4d2d] font-medium">{product.shop.rating}</span>
                        </div>
                        <div className="hover:bg-gray-50 p-2 rounded transition-colors cursor-pointer">
                            <span className="text-gray-500">Sản Phẩm: </span>
                            <span className="text-[#ee4d2d] font-medium">{product.shop.products}</span>
                        </div>
                        <div className="hover:bg-gray-50 p-2 rounded transition-colors cursor-pointer">
                            <span className="text-gray-500">Tỉ Lệ Phản Hồi: </span>
                            <span className="text-[#ee4d2d] font-medium">{product.shop.responseRate}</span>
                        </div>
                    </div>
                </div>

                {/* Description */}
                <div className="bg-white rounded-sm shadow-sm p-4 animate-fade-in-up" style={{ animationDelay: '200ms' }}>
                    <h2 className="bg-[#fafafa] px-4 py-2 text-sm font-medium mb-4">CHI TIẾT SẢN PHẨM</h2>
                    <div className="px-4 text-sm text-gray-600 leading-relaxed">
                        <p>{product.description}</p>
                    </div>

                    <h2 className="bg-[#fafafa] px-4 py-2 text-sm font-medium mb-4 mt-6">MÔ TẢ SẢN PHẨM</h2>
                    <div className="text-sm text-gray-600 leading-relaxed px-4 space-y-2">
                        <p>✅ Sản phẩm chính hãng 100%</p>
                        <p>✅ Bảo hành 12 tháng tại trung tâm ủy quyền</p>
                        <p>✅ Hỗ trợ đổi trả trong 7 ngày</p>
                        <p>✅ Giao hàng toàn quốc</p>
                        <p>✅ Thanh toán khi nhận hàng (COD)</p>
                        {product.tags && (
                            <div className="flex gap-2 pt-4">
                                {product.tags.map(tag => (
                                    <span key={tag} className="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded">
                                        {tag}
                                    </span>
                                ))}
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}
