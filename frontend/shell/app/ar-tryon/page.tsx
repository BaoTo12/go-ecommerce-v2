'use client';

import React, { useState, useRef, useEffect } from 'react';
import Image from 'next/image';
import Link from 'next/link';

interface TryOnProduct {
    id: string;
    name: string;
    image: string;
    price: number;
    type: 'glasses' | 'earrings' | 'necklace' | 'hat' | 'watch';
    overlay: string;
    position: { x: number; y: number; scale: number };
}

const TRYON_PRODUCTS: TryOnProduct[] = [
    {
        id: 'ar1',
        name: 'Kính Mắt Ray-Ban Aviator',
        image: 'https://images.unsplash.com/photo-1572635196237-14b3f281503f?w=200',
        price: 3990000,
        type: 'glasses',
        overlay: 'https://images.unsplash.com/photo-1572635196237-14b3f281503f?w=200',
        position: { x: 50, y: 35, scale: 0.4 },
    },
    {
        id: 'ar2',
        name: 'Đồng Hồ Casio G-Shock',
        image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=200',
        price: 2890000,
        type: 'watch',
        overlay: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=200',
        position: { x: 15, y: 75, scale: 0.25 },
    },
    {
        id: 'ar3',
        name: 'Mũ Bucket Hat Unisex',
        image: 'https://images.unsplash.com/photo-1588850561407-ed78c282e89b?w=200',
        price: 250000,
        type: 'hat',
        overlay: 'https://images.unsplash.com/photo-1588850561407-ed78c282e89b?w=200',
        position: { x: 50, y: 10, scale: 0.5 },
    },
];

export default function ARTryOnPage() {
    const [selectedProduct, setSelectedProduct] = useState<TryOnProduct | null>(null);
    const [isCameraActive, setIsCameraActive] = useState(false);
    const [capturedImage, setCapturedImage] = useState<string | null>(null);
    const videoRef = useRef<HTMLVideoElement>(null);
    const canvasRef = useRef<HTMLCanvasElement>(null);

    const startCamera = async () => {
        try {
            const stream = await navigator.mediaDevices.getUserMedia({
                video: { facingMode: 'user', width: 640, height: 480 }
            });
            if (videoRef.current) {
                videoRef.current.srcObject = stream;
                setIsCameraActive(true);
            }
        } catch (err) {
            console.error('Camera error:', err);
            alert('Không thể truy cập camera. Vui lòng cho phép quyền camera.');
        }
    };

    const stopCamera = () => {
        if (videoRef.current?.srcObject) {
            const stream = videoRef.current.srcObject as MediaStream;
            stream.getTracks().forEach(track => track.stop());
            setIsCameraActive(false);
        }
    };

    const capturePhoto = () => {
        if (videoRef.current && canvasRef.current) {
            const canvas = canvasRef.current;
            const video = videoRef.current;
            canvas.width = video.videoWidth;
            canvas.height = video.videoHeight;

            const ctx = canvas.getContext('2d');
            if (ctx) {
                // Mirror the image
                ctx.translate(canvas.width, 0);
                ctx.scale(-1, 1);
                ctx.drawImage(video, 0, 0);

                // Draw overlay if product selected
                if (selectedProduct) {
                    const img = document.createElement('img');
                    img.crossOrigin = 'anonymous';
                    img.onload = () => {
                        ctx.setTransform(1, 0, 0, 1, 0, 0); // Reset transform
                        const size = canvas.width * selectedProduct.position.scale;
                        const x = (selectedProduct.position.x / 100) * canvas.width - size / 2;
                        const y = (selectedProduct.position.y / 100) * canvas.height - size / 2;
                        ctx.drawImage(img, x, y, size, size);
                        setCapturedImage(canvas.toDataURL('image/png'));
                    };
                    img.src = selectedProduct.overlay;
                } else {
                    setCapturedImage(canvas.toDataURL('image/png'));
                }
            }
        }
    };

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    useEffect(() => {
        return () => stopCamera();
    }, []);

    return (
        <div className="min-h-screen bg-gray-900">
            {/* Header */}
            <div className="bg-gradient-to-r from-pink-500 to-violet-500 p-4">
                <h1 className="text-xl font-bold text-white flex items-center gap-2">
                    ✨ AR Try-On
                </h1>
                <p className="text-white/80 text-sm">Thử đồ ảo trước khi mua</p>
            </div>

            <div className="container mx-auto px-4 py-6">
                {/* Camera / Preview Area */}
                <div className="relative aspect-[4/3] bg-black rounded-2xl overflow-hidden mb-6">
                    {capturedImage ? (
                        <img src={capturedImage} alt="Captured" className="w-full h-full object-cover" />
                    ) : isCameraActive ? (
                        <>
                            <video
                                ref={videoRef}
                                autoPlay
                                playsInline
                                muted
                                className="w-full h-full object-cover"
                                style={{ transform: 'scaleX(-1)' }}
                            />
                            {/* AR Overlay */}
                            {selectedProduct && (
                                <div
                                    className="absolute pointer-events-none"
                                    style={{
                                        left: `${selectedProduct.position.x}%`,
                                        top: `${selectedProduct.position.y}%`,
                                        transform: `translate(-50%, -50%) scale(${selectedProduct.position.scale * 2})`,
                                    }}
                                >
                                    <img
                                        src={selectedProduct.overlay}
                                        alt=""
                                        className="w-32 h-32 object-contain opacity-80"
                                    />
                                </div>
                            )}
                        </>
                    ) : (
                        <div className="w-full h-full flex flex-col items-center justify-center">
                            <div className="text-6xl mb-4">📷</div>
                            <p className="text-white mb-4">Bật camera để thử đồ ảo</p>
                            <button
                                onClick={startCamera}
                                className="px-8 py-3 bg-gradient-to-r from-pink-500 to-violet-500 text-white rounded-full font-medium"
                            >
                                🎥 Bật Camera
                            </button>
                        </div>
                    )}

                    <canvas ref={canvasRef} className="hidden" />
                </div>

                {/* Camera Controls */}
                {isCameraActive && !capturedImage && (
                    <div className="flex justify-center gap-4 mb-6">
                        <button
                            onClick={capturePhoto}
                            className="w-16 h-16 bg-white rounded-full flex items-center justify-center text-3xl shadow-lg"
                        >
                            📸
                        </button>
                        <button
                            onClick={stopCamera}
                            className="w-12 h-12 bg-red-500 text-white rounded-full flex items-center justify-center"
                        >
                            ✕
                        </button>
                    </div>
                )}

                {capturedImage && (
                    <div className="flex justify-center gap-4 mb-6">
                        <button
                            onClick={() => setCapturedImage(null)}
                            className="px-6 py-3 bg-gray-700 text-white rounded-full"
                        >
                            ← Thử lại
                        </button>
                        <button className="px-6 py-3 bg-green-500 text-white rounded-full">
                            💾 Lưu ảnh
                        </button>
                        <button className="px-6 py-3 bg-blue-500 text-white rounded-full">
                            ↗️ Chia sẻ
                        </button>
                    </div>
                )}

                {/* Products to try */}
                <div className="bg-gray-800 rounded-2xl p-4">
                    <h3 className="text-white font-medium mb-4">👗 Chọn sản phẩm để thử</h3>
                    <div className="flex gap-4 overflow-x-auto pb-2">
                        {TRYON_PRODUCTS.map(product => (
                            <button
                                key={product.id}
                                onClick={() => setSelectedProduct(product)}
                                className={`flex-shrink-0 w-24 rounded-xl overflow-hidden transition-all ${selectedProduct?.id === product.id
                                        ? 'ring-2 ring-pink-500 scale-105'
                                        : 'opacity-70'
                                    }`}
                            >
                                <div className="aspect-square bg-white relative">
                                    <Image src={product.image} alt="" fill className="object-cover" unoptimized />
                                </div>
                                <div className="p-2 bg-gray-700">
                                    <p className="text-white text-xs truncate">{product.name}</p>
                                    <p className="text-pink-400 text-xs font-medium">₫{formatPrice(product.price)}</p>
                                </div>
                            </button>
                        ))}
                    </div>
                </div>

                {/* Selected Product Details */}
                {selectedProduct && (
                    <div className="mt-6 bg-white dark:bg-gray-800 rounded-2xl p-4 animate-fade-in-up">
                        <div className="flex gap-4 items-center">
                            <div className="w-20 h-20 bg-gray-100 rounded-lg overflow-hidden relative">
                                <Image src={selectedProduct.image} alt="" fill className="object-cover" unoptimized />
                            </div>
                            <div className="flex-1">
                                <h4 className="font-bold dark:text-white">{selectedProduct.name}</h4>
                                <p className="text-xl font-bold text-[#ee4d2d]">₫{formatPrice(selectedProduct.price)}</p>
                            </div>
                            <Link
                                href={`/products/${selectedProduct.id}`}
                                className="px-6 py-3 bg-[#ee4d2d] text-white rounded-full font-medium"
                            >
                                Mua ngay
                            </Link>
                        </div>
                    </div>
                )}

                {/* Tips */}
                <div className="mt-6 text-center text-gray-400 text-sm">
                    <p>💡 Tip: Đặt khuôn mặt ở giữa camera và di chuyển chậm để AR tracking tốt hơn</p>
                </div>
            </div>
        </div>
    );
}
