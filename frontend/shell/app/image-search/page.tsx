'use client';

import React, { useState, useRef } from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { Product } from '@/services/productService';

// Mock similar products based on "image analysis"
const SIMILAR_PRODUCTS: Product[] = [
    {
        id: 'p1',
        name: 'iPhone 15 Pro Max 256GB Titan Xanh',
        description: 'iPhone 15 Pro Max chính hãng',
        price: 29990000,
        originalPrice: 34990000,
        discount: 14,
        currency: 'VND',
        category: 'Điện thoại',
        categoryId: 'phones',
        images: ['https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=400'],
        thumbnail: 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=300',
        rating: 4.9,
        reviews: 8560,
        sold: 12300,
        soldDisplay: '12.3k',
        stock: 156,
        location: 'TP. Hồ Chí Minh',
        shop: { id: 'shop1', name: 'Apple Store', avatar: '', rating: 4.9, products: 156, responseRate: '95%', location: 'HCM', isOfficial: true },
        createdAt: '2024-01-15',
    },
    {
        id: 'p2',
        name: 'Samsung Galaxy S24 Ultra 512GB',
        description: 'Samsung Galaxy S24 Ultra chính hãng',
        price: 25990000,
        originalPrice: 29990000,
        discount: 13,
        currency: 'VND',
        category: 'Điện thoại',
        categoryId: 'phones',
        images: ['https://images.unsplash.com/photo-1610945415295-d9bbf067e59c?w=400'],
        thumbnail: 'https://images.unsplash.com/photo-1610945415295-d9bbf067e59c?w=300',
        rating: 4.8,
        reviews: 5430,
        sold: 8700,
        soldDisplay: '8.7k',
        stock: 89,
        location: 'Hà Nội',
        shop: { id: 'shop2', name: 'Samsung Store', avatar: '', rating: 4.9, products: 234, responseRate: '97%', location: 'HN', isOfficial: true },
        createdAt: '2024-02-10',
    },
];

export default function ImageSearchPage() {
    const [selectedImage, setSelectedImage] = useState<string | null>(null);
    const [isAnalyzing, setIsAnalyzing] = useState(false);
    const [results, setResults] = useState<Product[]>([]);
    const [detectedObjects, setDetectedObjects] = useState<string[]>([]);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;

        const reader = new FileReader();
        reader.onload = (event) => {
            setSelectedImage(event.target?.result as string);
            analyzeImage();
        };
        reader.readAsDataURL(file);
    };

    const analyzeImage = () => {
        setIsAnalyzing(true);
        setResults([]);
        setDetectedObjects([]);

        // Simulate AI analysis
        setTimeout(() => {
            setDetectedObjects(['Điện thoại', 'Màn hình', 'Vỏ kim loại']);
            setResults(SIMILAR_PRODUCTS);
            setIsAnalyzing(false);
        }, 2500);
    };

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    return (
        <div className="min-h-screen bg-gray-100 dark:bg-gray-900">
            {/* Header */}
            <div className="bg-gradient-to-r from-blue-600 to-purple-600 p-6">
                <div className="container mx-auto">
                    <h1 className="text-2xl font-bold text-white flex items-center gap-2">
                        📷 Tìm Kiếm Bằng Hình Ảnh
                    </h1>
                    <p className="text-white/80">Chụp hoặc tải ảnh lên để tìm sản phẩm tương tự</p>
                </div>
            </div>

            <div className="container mx-auto px-4 py-6">
                {/* Upload Area */}
                <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6 mb-6">
                    <div
                        onClick={() => fileInputRef.current?.click()}
                        className={`border-2 border-dashed rounded-xl p-8 text-center cursor-pointer transition-all ${selectedImage
                                ? 'border-green-500 bg-green-50 dark:bg-green-900/20'
                                : 'border-gray-300 dark:border-gray-600 hover:border-blue-500'
                            }`}
                    >
                        <input
                            ref={fileInputRef}
                            type="file"
                            accept="image/*"
                            onChange={handleImageUpload}
                            className="hidden"
                        />

                        {selectedImage ? (
                            <div className="relative inline-block">
                                <img
                                    src={selectedImage}
                                    alt="Uploaded"
                                    className="max-h-64 rounded-lg mx-auto"
                                />
                                {isAnalyzing && (
                                    <div className="absolute inset-0 bg-black/50 rounded-lg flex items-center justify-center">
                                        <div className="text-white text-center">
                                            <div className="loading-spinner mx-auto mb-2" style={{ borderTopColor: 'white' }} />
                                            <p>Đang phân tích hình ảnh...</p>
                                        </div>
                                    </div>
                                )}
                            </div>
                        ) : (
                            <>
                                <div className="text-6xl mb-4">📷</div>
                                <p className="text-gray-600 dark:text-gray-400 mb-2">Kéo thả hoặc click để tải ảnh lên</p>
                                <p className="text-sm text-gray-400">Hỗ trợ: JPG, PNG, WEBP (tối đa 10MB)</p>
                            </>
                        )}
                    </div>

                    <div className="flex gap-4 mt-4 justify-center">
                        <button
                            onClick={() => fileInputRef.current?.click()}
                            className="px-6 py-3 bg-blue-600 text-white rounded-full font-medium hover:opacity-90"
                        >
                            📁 Tải ảnh lên
                        </button>
                        <button className="px-6 py-3 bg-gray-200 dark:bg-gray-700 rounded-full font-medium dark:text-white">
                            📸 Chụp ảnh
                        </button>
                    </div>
                </div>

                {/* Detected Objects */}
                {detectedObjects.length > 0 && (
                    <div className="bg-white dark:bg-gray-800 rounded-xl p-4 mb-6 animate-fade-in-up">
                        <h3 className="font-medium mb-2 dark:text-white">🔍 Phát hiện trong ảnh:</h3>
                        <div className="flex gap-2 flex-wrap">
                            {detectedObjects.map((obj, i) => (
                                <span key={i} className="px-3 py-1 bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300 rounded-full text-sm">
                                    {obj}
                                </span>
                            ))}
                        </div>
                    </div>
                )}

                {/* Results */}
                {results.length > 0 && (
                    <div className="animate-fade-in-up">
                        <h3 className="font-bold mb-4 dark:text-white">🛒 Sản phẩm tương tự ({results.length})</h3>
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                            {results.map((product, index) => (
                                <Link
                                    key={product.id}
                                    href={`/products/${product.id}`}
                                    className="bg-white dark:bg-gray-800 rounded-lg shadow-sm overflow-hidden animate-fade-in-up"
                                    style={{ animationDelay: `${index * 100}ms` }}
                                >
                                    <div className="aspect-square relative bg-gray-100 dark:bg-gray-700">
                                        <Image
                                            src={product.thumbnail}
                                            alt={product.name}
                                            fill
                                            className="object-cover"
                                            unoptimized
                                        />
                                        <div className="absolute top-2 left-2 px-2 py-1 bg-green-500 text-white text-xs rounded">
                                            98% match
                                        </div>
                                    </div>
                                    <div className="p-3">
                                        <h4 className="text-sm font-medium line-clamp-2 dark:text-white">{product.name}</h4>
                                        <p className="text-lg font-bold text-[#ee4d2d] mt-1">₫{formatPrice(product.price)}</p>
                                        <p className="text-xs text-gray-500">⭐ {product.rating} • {product.soldDisplay} đã bán</p>
                                    </div>
                                </Link>
                            ))}
                        </div>
                    </div>
                )}

                {/* Tips */}
                <div className="mt-8 bg-blue-50 dark:bg-blue-900/20 rounded-xl p-6">
                    <h3 className="font-bold mb-3 dark:text-white">💡 Mẹo tìm kiếm tốt hơn</h3>
                    <ul className="space-y-2 text-sm text-gray-600 dark:text-gray-400">
                        <li>• Chụp ảnh rõ nét, đủ ánh sáng</li>
                        <li>• Focus vào sản phẩm chính trong ảnh</li>
                        <li>• Tránh ảnh bị mờ hoặc nhiều vật thể</li>
                        <li>• Có thể chụp từ nhiều góc độ để tìm chính xác hơn</li>
                    </ul>
                </div>
            </div>
        </div>
    );
}
