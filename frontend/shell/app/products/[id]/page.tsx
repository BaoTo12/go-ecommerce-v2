'use client';

import React, { useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';

const allProducts = [
    { id: 'p1', name: 'iPhone 15 Pro Max 256GB Xanh Titan', price: 29990000, originalPrice: 34990000, discount: 14, sold: 1234, rating: 4.9, reviews: 856, image: '📱', category: 'Điện thoại', description: 'iPhone 15 Pro Max với chip A17 Pro, camera 48MP, màn hình Super Retina XDR 6.7 inch, thời lượng pin cả ngày.' },
    { id: 'p2', name: 'Samsung Galaxy S24 Ultra', price: 25990000, originalPrice: 29990000, discount: 13, sold: 987, rating: 4.8, reviews: 543, image: '📲', category: 'Điện thoại', description: 'Galaxy S24 Ultra với S Pen tích hợp, camera 200MP, màn hình Dynamic AMOLED 2X 6.8 inch.' },
    { id: 'p3', name: 'MacBook Air M3 13"', price: 27990000, originalPrice: 31990000, discount: 12, sold: 654, rating: 4.9, reviews: 432, image: '💻', category: 'Laptop', description: 'MacBook Air với chip M3, 8GB RAM, 256GB SSD, màn hình Liquid Retina 13.6 inch.' },
    { id: 'p4', name: 'Áo Hoodie Unisex Premium', price: 299000, originalPrice: 450000, discount: 34, sold: 5432, rating: 4.7, reviews: 1234, image: '👕', category: 'Thời trang', description: 'Áo hoodie unisex chất liệu cotton dày dặn, form rộng thoải mái, phù hợp mọi dáng người.' },
    { id: 'p5', name: 'Giày Nike Air Max 90', price: 3990000, originalPrice: 4590000, discount: 13, sold: 2345, rating: 4.8, reviews: 876, image: '👟', category: 'Giày dép', description: 'Nike Air Max 90 chính hãng, đệm Air Max êm ái, thiết kế iconic từ năm 1990.' },
    { id: 'p6', name: 'Son Dưỡng Môi Dior', price: 950000, originalPrice: 1200000, discount: 21, sold: 8765, rating: 4.6, reviews: 3456, image: '💄', category: 'Làm đẹp', description: 'Son dưỡng môi Dior Addict Lip Glow, dưỡng ẩm và tạo màu tự nhiên.' },
];

export default function ProductDetailPage() {
    const params = useParams();
    const router = useRouter();
    const productId = params.id as string;

    const product = allProducts.find(p => p.id === productId) || allProducts[0];

    const [quantity, setQuantity] = useState(1);
    const [selectedColor, setSelectedColor] = useState('Đen');
    const [addedToCart, setAddedToCart] = useState(false);

    const colors = ['Đen', 'Trắng', 'Xanh', 'Hồng'];

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    const addToCart = () => {
        setAddedToCart(true);
        setTimeout(() => setAddedToCart(false), 3000);
    };

    const buyNow = () => {
        router.push('/checkout');
    };

    return (
        <div className="min-h-screen bg-[#F5F5F5] animate-fade-in">
            {/* Toast */}
            {addedToCart && (
                <div className="fixed top-24 right-4 z-50 bg-green-500 text-white px-6 py-3 rounded-xl shadow-xl animate-slide-in-right">
                    ✓ Đã thêm vào giỏ hàng!
                </div>
            )}

            {/* Breadcrumb */}
            <div className="bg-white border-b">
                <div className="container mx-auto px-4 py-3">
                    <div className="flex items-center gap-2 text-sm text-gray-500">
                        <Link href="/" className="hover:text-[#EE4D2D]">Trang chủ</Link>
                        <span>›</span>
                        <Link href="/products" className="hover:text-[#EE4D2D]">Sản phẩm</Link>
                        <span>›</span>
                        <span className="text-gray-800">{product.name}</span>
                    </div>
                </div>
            </div>

            <div className="container mx-auto px-4 py-6">
                <div className="bg-white rounded-2xl shadow-sm overflow-hidden">
                    <div className="grid md:grid-cols-2 gap-8 p-6">
                        {/* Product Image */}
                        <div className="space-y-4">
                            <div className="aspect-square bg-gradient-to-br from-gray-50 to-gray-100 rounded-2xl flex items-center justify-center">
                                <span className="text-[200px] animate-float">{product.image}</span>
                            </div>
                            <div className="flex gap-2">
                                {[1, 2, 3, 4].map(i => (
                                    <div key={i} className="w-20 h-20 bg-gray-100 rounded-xl flex items-center justify-center cursor-pointer hover:ring-2 hover:ring-[#EE4D2D] transition-all">
                                        <span className="text-3xl">{product.image}</span>
                                    </div>
                                ))}
                            </div>
                        </div>

                        {/* Product Info */}
                        <div className="space-y-6">
                            <div>
                                <span className="inline-block bg-[#EE4D2D] text-white text-xs px-3 py-1 rounded-full mb-2">
                                    Yêu thích+
                                </span>
                                <h1 className="text-2xl font-bold text-gray-800">{product.name}</h1>
                            </div>

                            {/* Rating */}
                            <div className="flex items-center gap-4">
                                <div className="flex items-center gap-1">
                                    <span className="text-[#EE4D2D] font-bold">{product.rating}</span>
                                    <span className="text-yellow-400">★★★★★</span>
                                </div>
                                <span className="text-gray-400">|</span>
                                <span className="text-gray-500">{product.reviews.toLocaleString()} đánh giá</span>
                                <span className="text-gray-400">|</span>
                                <span className="text-gray-500">{product.sold.toLocaleString()} đã bán</span>
                            </div>

                            {/* Price */}
                            <div className="bg-[#FAFAFA] p-4 rounded-xl">
                                <div className="flex items-baseline gap-3">
                                    <span className="text-3xl font-bold text-[#EE4D2D]">
                                        ₫{formatPrice(product.price)}
                                    </span>
                                    <span className="text-gray-400 line-through">
                                        ₫{formatPrice(product.originalPrice)}
                                    </span>
                                    <span className="bg-[#EE4D2D] text-white text-sm px-2 py-1 rounded-lg">
                                        -{product.discount}%
                                    </span>
                                </div>
                            </div>

                            {/* Color Selection */}
                            <div>
                                <span className="text-gray-600 mb-2 block">Màu sắc</span>
                                <div className="flex gap-2">
                                    {colors.map(color => (
                                        <button
                                            key={color}
                                            onClick={() => setSelectedColor(color)}
                                            className={`px-4 py-2 rounded-xl border-2 transition-all ${selectedColor === color
                                                    ? 'border-[#EE4D2D] bg-[#FFEEE8] text-[#EE4D2D]'
                                                    : 'border-gray-200 hover:border-[#EE4D2D]'
                                                }`}
                                        >
                                            {color}
                                        </button>
                                    ))}
                                </div>
                            </div>

                            {/* Quantity */}
                            <div>
                                <span className="text-gray-600 mb-2 block">Số lượng</span>
                                <div className="flex items-center gap-3">
                                    <div className="flex items-center border-2 rounded-xl overflow-hidden">
                                        <button
                                            onClick={() => setQuantity(Math.max(1, quantity - 1))}
                                            className="px-4 py-2 hover:bg-gray-100"
                                        >
                                            −
                                        </button>
                                        <span className="px-6 py-2 border-x-2">{quantity}</span>
                                        <button
                                            onClick={() => setQuantity(quantity + 1)}
                                            className="px-4 py-2 hover:bg-gray-100"
                                        >
                                            +
                                        </button>
                                    </div>
                                    <span className="text-gray-500">Còn 999 sản phẩm</span>
                                </div>
                            </div>

                            {/* Action Buttons */}
                            <div className="flex gap-4 pt-4">
                                <button
                                    onClick={addToCart}
                                    className="flex-1 py-4 border-2 border-[#EE4D2D] text-[#EE4D2D] rounded-xl font-bold hover:bg-[#FFEEE8] transition-all flex items-center justify-center gap-2"
                                >
                                    <span>🛒</span> Thêm vào giỏ
                                </button>
                                <button
                                    onClick={buyNow}
                                    className="flex-1 py-4 bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] text-white rounded-xl font-bold hover:opacity-90 transition-all"
                                >
                                    Mua ngay
                                </button>
                            </div>

                            {/* Shipping Info */}
                            <div className="border-t pt-4 space-y-3">
                                <div className="flex items-center gap-3 text-sm">
                                    <span className="text-green-500">🚚</span>
                                    <span>Miễn phí vận chuyển cho đơn từ 500K</span>
                                </div>
                                <div className="flex items-center gap-3 text-sm">
                                    <span className="text-blue-500">🔄</span>
                                    <span>Đổi trả trong 15 ngày</span>
                                </div>
                                <div className="flex items-center gap-3 text-sm">
                                    <span className="text-yellow-500">✅</span>
                                    <span>100% hàng chính hãng</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Description */}
                    <div className="border-t p-6">
                        <h2 className="text-xl font-bold mb-4">Mô tả sản phẩm</h2>
                        <p className="text-gray-600 leading-relaxed">{product.description}</p>
                    </div>
                </div>
            </div>
        </div>
    );
}
