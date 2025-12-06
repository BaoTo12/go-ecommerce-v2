'use client';

import React, { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';

const allProducts = [
    { id: 'p1', name: 'iPhone 15 Pro Max 256GB Titan Xanh Chính Hãng VN/A Bảo Hành 12 Tháng', price: 29990000, originalPrice: 34990000, sold: '12.3k', rating: 4.9, reviews: 8560, image: '📱', category: 'Điện thoại', description: 'iPhone 15 Pro Max với chip A17 Pro mạnh mẽ nhất, camera 48MP, màn hình Super Retina XDR 6.7 inch, thời lượng pin cả ngày. Thiết kế titan cao cấp, nhẹ hơn và bền hơn.', shop: 'Apple Store Official', shopLocation: 'TP. Hồ Chí Minh', shopRating: 4.9, shopProducts: 156, shopResponse: '95%' },
    { id: 'p2', name: 'Samsung Galaxy S24 Ultra 512GB Xám Titan Chính Hãng', price: 25990000, originalPrice: 29990000, sold: '8.7k', rating: 4.8, reviews: 5430, image: '📲', category: 'Điện thoại', description: 'Galaxy S24 Ultra với S Pen tích hợp, camera 200MP, màn hình Dynamic AMOLED 2X 6.8 inch, chip Snapdragon 8 Gen 3.', shop: 'Samsung Official', shopLocation: 'Hà Nội', shopRating: 4.9, shopProducts: 234, shopResponse: '97%' },
    { id: 'p3', name: 'MacBook Air M3 13 inch 256GB Space Gray 2024', price: 27990000, originalPrice: 31990000, sold: '3.2k', rating: 4.9, reviews: 2340, image: '💻', category: 'Laptop', description: 'MacBook Air với chip M3 thế hệ mới, 8GB RAM, 256GB SSD, màn hình Liquid Retina 13.6 inch sắc nét.', shop: 'Apple Store Official', shopLocation: 'TP. Hồ Chí Minh', shopRating: 4.9, shopProducts: 156, shopResponse: '95%' },
    { id: 'p4', name: 'Áo Hoodie Unisex Form Rộng Nỉ Cotton Dày Dặn Premium', price: 199000, originalPrice: 350000, sold: '45.2k', rating: 4.7, reviews: 12340, image: '👕', category: 'Thời trang', description: 'Áo hoodie unisex chất liệu cotton dày dặn, form rộng thoải mái, phù hợp mọi dáng người. Nhiều màu sắc để lựa chọn.', shop: 'Fashion Store', shopLocation: 'Hà Nội', shopRating: 4.7, shopProducts: 567, shopResponse: '92%' },
    { id: 'p5', name: 'Giày Nike Air Force 1 07 Low White Chính Hãng', price: 2590000, originalPrice: 3200000, sold: '5.2k', rating: 4.8, reviews: 3456, image: '👟', category: 'Giày dép', description: 'Nike Air Force 1 chính hãng, đệm Air êm ái, thiết kế iconic từ năm 1982, phù hợp mọi phong cách.', shop: 'Nike Official', shopLocation: 'TP. Hồ Chí Minh', shopRating: 4.8, shopProducts: 234, shopResponse: '96%' },
    { id: 'p6', name: 'Son Dưỡng Môi Dior Addict Lip Glow Fullsize', price: 950000, originalPrice: 1200000, sold: '18.7k', rating: 4.9, reviews: 8765, image: '💄', category: 'Làm đẹp', description: 'Son dưỡng môi Dior Addict Lip Glow, dưỡng ẩm và tạo màu tự nhiên, công nghệ Color Reviver phản ứng với độ pH của môi.', shop: 'Dior Beauty Official', shopLocation: 'TP. Hồ Chí Minh', shopRating: 4.9, shopProducts: 89, shopResponse: '98%' },
];

export default function ProductDetailPage() {
    const params = useParams();
    const router = useRouter();
    const productId = params.id as string;

    const product = allProducts.find(p => p.id === productId) || allProducts[0];

    const [quantity, setQuantity] = useState(1);
    const [selectedVariant, setSelectedVariant] = useState(0);
    const [selectedImage, setSelectedImage] = useState(0);
    const [notification, setNotification] = useState<string | null>(null);
    const [isLoaded, setIsLoaded] = useState(false);
    const [isAddingToCart, setIsAddingToCart] = useState(false);
    const [isBuying, setIsBuying] = useState(false);
    const [showHearts, setShowHearts] = useState(false);

    const variants = ['Đen', 'Trắng', 'Xanh', 'Hồng'];

    useEffect(() => {
        setIsLoaded(true);
    }, []);

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);
    const getDiscount = () => Math.round((1 - product.price / product.originalPrice) * 100);

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

    return (
        <div className={`min-h-screen bg-[#f5f5f5] ${isLoaded ? 'animate-fade-in' : 'opacity-0'}`}>
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
                            <div className="aspect-square bg-gray-50 rounded-sm flex items-center justify-center mb-2 overflow-hidden">
                                <span className={`text-[180px] transition-all duration-500 ${selectedImage === 0 ? 'animate-scale-in' : ''}`}>
                                    {product.image}
                                </span>
                            </div>
                            <div className="flex gap-2">
                                {[0, 1, 2, 3, 4].map(i => (
                                    <button
                                        key={i}
                                        onClick={() => setSelectedImage(i)}
                                        className={`w-16 h-16 bg-gray-100 rounded-sm flex items-center justify-center text-2xl 
                               cursor-pointer border-2 transition-all duration-300 hover:scale-105
                               ${selectedImage === i ? 'border-[#ee4d2d]' : 'border-transparent hover:border-gray-300'}`}
                                    >
                                        {product.image}
                                    </button>
                                ))}
                            </div>
                        </div>

                        {/* Details */}
                        <div className="md:col-span-3 animate-fade-in-right">
                            <div className="flex items-start gap-2 mb-2">
                                <span className="bg-[#ee4d2d] text-white text-[10px] px-1 py-0.5 animate-pulse">Mall</span>
                                <h1 className="text-lg leading-tight flex-1">{product.name}</h1>
                            </div>

                            {/* Rating & Sold */}
                            <div className="flex items-center gap-4 text-sm py-3 border-b">
                                <div className="flex items-center gap-1">
                                    <span className="text-[#ee4d2d] font-medium border-b border-[#ee4d2d]">{product.rating}</span>
                                    <span className="star-rating">★★★★★</span>
                                </div>
                                <span className="text-gray-300">|</span>
                                <div className="hover:scale-105 transition-transform cursor-pointer">
                                    <span className="font-medium border-b border-gray-600">{(product.reviews / 1000).toFixed(1)}k</span>
                                    <span className="text-gray-500 ml-1">Đánh Giá</span>
                                </div>
                                <span className="text-gray-300">|</span>
                                <div className="hover:scale-105 transition-transform cursor-pointer">
                                    <span className="font-medium">{product.sold}</span>
                                    <span className="text-gray-500 ml-1">Đã Bán</span>
                                </div>
                            </div>

                            {/* Price */}
                            <div className="bg-[#fafafa] p-4 my-3 animate-fade-in">
                                <div className="flex items-center gap-3">
                                    <span className="text-gray-400 line-through text-sm">₫{formatPrice(product.originalPrice)}</span>
                                    <span className="text-[#ee4d2d] text-3xl font-medium animate-pulse-glow rounded px-2">₫{formatPrice(product.price)}</span>
                                    <span className="bg-[#ee4d2d] text-white text-xs px-2 py-0.5 rounded-sm animate-bounce">{getDiscount()}% GIẢM</span>
                                </div>
                            </div>

                            {/* Vouchers */}
                            <div className="flex items-center gap-4 py-3 text-sm">
                                <span className="text-gray-500 w-24">Mã Giảm Giá</span>
                                <div className="flex gap-2">
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
                                    <span className="text-[#00bfa5] font-medium">Miễn phí vận chuyển</span>
                                </div>
                            </div>

                            {/* Variants */}
                            <div className="flex items-start gap-4 py-3 text-sm">
                                <span className="text-gray-500 w-24 pt-2">Màu Sắc</span>
                                <div className="flex flex-wrap gap-2">
                                    {variants.map((v, i) => (
                                        <button
                                            key={v}
                                            onClick={() => setSelectedVariant(i)}
                                            className={`px-4 py-2 border rounded-sm transition-all duration-200 hover-shrink ${selectedVariant === i
                                                    ? 'border-[#ee4d2d] text-[#ee4d2d] bg-[#fef6f5] scale-105'
                                                    : 'border-gray-300 hover:border-[#ee4d2d] hover:text-[#ee4d2d]'
                                                }`}
                                        >
                                            {v}
                                        </button>
                                    ))}
                                </div>
                            </div>

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
                                <span className="text-gray-400">999 sản phẩm có sẵn</span>
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
                        <div className="w-16 h-16 bg-gradient-to-br from-[#ee4d2d] to-[#ff6533] rounded-full flex items-center justify-center text-white text-2xl font-bold hover:scale-110 transition-transform cursor-pointer">
                            {product.shop.charAt(0)}
                        </div>
                        <div className="flex-1">
                            <h3 className="font-medium">{product.shop}</h3>
                            <p className="text-xs text-gray-500 flex items-center gap-1">
                                <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
                                Online {product.shopLocation}
                            </p>
                        </div>
                        <button className="px-4 py-1.5 border border-[#ee4d2d] text-[#ee4d2d] text-sm hover:bg-[#fef6f5] transition-colors hover-shrink">
                            Xem Shop
                        </button>
                        <button className="px-4 py-1.5 border border-gray-300 text-gray-600 text-sm hover:border-[#ee4d2d] hover:text-[#ee4d2d] transition-colors hover-shrink">
                            Chat Ngay
                        </button>
                    </div>
                    <div className="grid grid-cols-3 gap-4 mt-4 pt-4 border-t text-sm">
                        <div className="hover:bg-gray-50 p-2 rounded transition-colors cursor-pointer">
                            <span className="text-gray-500">Đánh Giá: </span>
                            <span className="text-[#ee4d2d] font-medium">{product.shopRating}</span>
                        </div>
                        <div className="hover:bg-gray-50 p-2 rounded transition-colors cursor-pointer">
                            <span className="text-gray-500">Sản Phẩm: </span>
                            <span className="text-[#ee4d2d] font-medium">{product.shopProducts}</span>
                        </div>
                        <div className="hover:bg-gray-50 p-2 rounded transition-colors cursor-pointer">
                            <span className="text-gray-500">Tỉ Lệ Phản Hồi: </span>
                            <span className="text-[#ee4d2d] font-medium">{product.shopResponse}</span>
                        </div>
                    </div>
                </div>

                {/* Description */}
                <div className="bg-white rounded-sm shadow-sm p-4 animate-fade-in-up" style={{ animationDelay: '200ms' }}>
                    <h2 className="bg-[#fafafa] px-4 py-2 text-sm font-medium mb-4">CHI TIẾT SẢN PHẨM</h2>
                    <p className="text-sm text-gray-600 leading-relaxed px-4">{product.description}</p>

                    <h2 className="bg-[#fafafa] px-4 py-2 text-sm font-medium mb-4 mt-6">MÔ TẢ SẢN PHẨM</h2>
                    <div className="text-sm text-gray-600 leading-relaxed px-4 space-y-2">
                        <p>✅ Sản phẩm chính hãng 100%</p>
                        <p>✅ Bảo hành 12 tháng tại trung tâm ủy quyền</p>
                        <p>✅ Hỗ trợ đổi trả trong 7 ngày</p>
                        <p>✅ Giao hàng toàn quốc</p>
                        <p>✅ Thanh toán khi nhận hàng (COD)</p>
                    </div>
                </div>
            </div>
        </div>
    );
}
