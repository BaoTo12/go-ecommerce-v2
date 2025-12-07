'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import Image from 'next/image';

interface LiveStream {
    id: string;
    shopName: string;
    shopAvatar: string;
    title: string;
    thumbnail: string;
    viewers: number;
    isLive: boolean;
    products: { name: string; price: number; discount: number }[];
}

const LIVE_STREAMS: LiveStream[] = [
    {
        id: 'live1',
        shopName: 'Apple Store Official',
        shopAvatar: 'https://ui-avatars.com/api/?name=Apple&background=000&color=fff',
        title: '⚡ Flash Sale iPhone 15 Series - Giảm đến 3 TRIỆU',
        thumbnail: 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=400',
        viewers: 12453,
        isLive: true,
        products: [
            { name: 'iPhone 15 Pro Max', price: 29990000, discount: 10 },
            { name: 'iPhone 15 Pro', price: 26990000, discount: 8 },
        ],
    },
    {
        id: 'live2',
        shopName: 'Nike Official Store',
        shopAvatar: 'https://ui-avatars.com/api/?name=Nike&background=000&color=fff',
        title: '🔥 Giày Nike chính hãng - Săn deal 50%',
        thumbnail: 'https://images.unsplash.com/photo-1600269452121-4f2416e55c28?w=400',
        viewers: 8234,
        isLive: true,
        products: [
            { name: 'Nike Air Force 1', price: 2590000, discount: 20 },
            { name: 'Nike Air Max', price: 3190000, discount: 15 },
        ],
    },
    {
        id: 'live3',
        shopName: 'Dior Beauty Official',
        shopAvatar: 'https://ui-avatars.com/api/?name=Dior&background=9c27b0&color=fff',
        title: '💄 Makeup Tutorial & Giveaway Son Dior',
        thumbnail: 'https://images.unsplash.com/photo-1586495777744-4413f21062fa?w=400',
        viewers: 5621,
        isLive: true,
        products: [
            { name: 'Son Dior Addict', price: 950000, discount: 5 },
        ],
    },
    {
        id: 'live4',
        shopName: 'Samsung Official',
        shopAvatar: 'https://ui-avatars.com/api/?name=Samsung&background=1428a0&color=fff',
        title: 'Ra mắt Galaxy S24 Ultra - Đặt trước giảm 5 triệu',
        thumbnail: 'https://images.unsplash.com/photo-1610945415295-d9bbf067e59c?w=400',
        viewers: 15234,
        isLive: false,
        products: [],
    },
];

export default function ShopLivePage() {
    const [selectedStream, setSelectedStream] = useState<LiveStream | null>(null);

    const formatViewers = (num: number) => {
        if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
        return num.toString();
    };

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    return (
        <div className="min-h-screen bg-[#1a1a1a]">
            {/* Header */}
            <div className="bg-gradient-to-r from-[#ee4d2d] to-[#ff6633]">
                <div className="container mx-auto px-4 py-4">
                    <div className="flex items-center justify-between">
                        <h1 className="text-white text-2xl font-bold flex items-center gap-3">
                            <span className="w-3 h-3 bg-white rounded-full animate-pulse" />
                            Shopee Live
                        </h1>
                        <div className="flex gap-4 text-white text-sm">
                            <Link href="/live" className="hover:opacity-80">Đang phát</Link>
                            <Link href="#" className="hover:opacity-80">Sắp phát</Link>
                            <Link href="#" className="hover:opacity-80">Phát lại</Link>
                        </div>
                    </div>
                </div>
            </div>

            <div className="container mx-auto px-4 py-6">
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                    {/* Main Video Player */}
                    <div className="lg:col-span-2">
                        {selectedStream ? (
                            <div className="bg-black rounded-lg overflow-hidden aspect-video relative animate-fade-in">
                                <Image
                                    src={selectedStream.thumbnail}
                                    alt={selectedStream.title}
                                    fill
                                    className="object-cover"
                                    unoptimized
                                />
                                <div className="absolute inset-0 bg-black/30 flex items-center justify-center">
                                    <button className="w-20 h-20 bg-white/20 backdrop-blur-sm rounded-full flex items-center justify-center hover:bg-white/30 transition-all">
                                        <span className="text-white text-4xl">▶</span>
                                    </button>
                                </div>
                                {/* Live badge */}
                                <div className="absolute top-4 left-4 flex items-center gap-2">
                                    <span className="bg-red-600 text-white text-xs font-bold px-2 py-1 rounded flex items-center gap-1">
                                        <span className="w-2 h-2 bg-white rounded-full animate-pulse" />
                                        LIVE
                                    </span>
                                    <span className="bg-black/50 text-white text-xs px-2 py-1 rounded">
                                        👁 {formatViewers(selectedStream.viewers)}
                                    </span>
                                </div>
                                {/* Stream info */}
                                <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 p-4">
                                    <div className="flex items-center gap-3">
                                        <div className="w-10 h-10 rounded-full overflow-hidden">
                                            <Image
                                                src={selectedStream.shopAvatar}
                                                alt={selectedStream.shopName}
                                                width={40}
                                                height={40}
                                                className="object-cover"
                                                unoptimized
                                            />
                                        </div>
                                        <div>
                                            <h3 className="text-white font-medium">{selectedStream.shopName}</h3>
                                            <p className="text-white/80 text-sm">{selectedStream.title}</p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        ) : (
                            <div className="bg-gray-800 rounded-lg aspect-video flex items-center justify-center">
                                <div className="text-center text-gray-400">
                                    <div className="text-5xl mb-4">📺</div>
                                    <p>Chọn một live stream để xem</p>
                                </div>
                            </div>
                        )}

                        {/* Products in stream */}
                        {selectedStream && selectedStream.products.length > 0 && (
                            <div className="mt-4 bg-gray-800 rounded-lg p-4">
                                <h3 className="text-white font-medium mb-3">🛒 Sản phẩm đang bán</h3>
                                <div className="flex gap-4 overflow-x-auto pb-2">
                                    {selectedStream.products.map((product, i) => (
                                        <div key={i} className="flex-shrink-0 bg-gray-700 rounded-lg p-3 w-48">
                                            <h4 className="text-white text-sm truncate">{product.name}</h4>
                                            <div className="flex items-center gap-2 mt-2">
                                                <span className="text-[#ee4d2d] font-bold">₫{formatPrice(product.price)}</span>
                                                <span className="text-xs text-white bg-[#ee4d2d] px-1 rounded">-{product.discount}%</span>
                                            </div>
                                            <button className="w-full mt-2 py-2 bg-[#ee4d2d] text-white text-sm rounded hover:opacity-90">
                                                Mua ngay
                                            </button>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}
                    </div>

                    {/* Stream List */}
                    <div className="space-y-4">
                        <h2 className="text-white font-medium">Đang phát trực tiếp</h2>
                        {LIVE_STREAMS.filter(s => s.isLive).map((stream, index) => (
                            <button
                                key={stream.id}
                                onClick={() => setSelectedStream(stream)}
                                className={`w-full bg-gray-800 rounded-lg overflow-hidden hover:ring-2 hover:ring-[#ee4d2d] transition-all animate-fade-in-up ${selectedStream?.id === stream.id ? 'ring-2 ring-[#ee4d2d]' : ''
                                    }`}
                                style={{ animationDelay: `${index * 50}ms` }}
                            >
                                <div className="relative aspect-video">
                                    <Image
                                        src={stream.thumbnail}
                                        alt={stream.title}
                                        fill
                                        className="object-cover"
                                        unoptimized
                                    />
                                    <div className="absolute top-2 left-2 flex items-center gap-1">
                                        <span className="bg-red-600 text-white text-xs px-1.5 py-0.5 rounded">LIVE</span>
                                        <span className="bg-black/50 text-white text-xs px-1.5 py-0.5 rounded">
                                            {formatViewers(stream.viewers)}
                                        </span>
                                    </div>
                                </div>
                                <div className="p-3 flex items-center gap-2">
                                    <div className="w-8 h-8 rounded-full overflow-hidden flex-shrink-0">
                                        <Image
                                            src={stream.shopAvatar}
                                            alt={stream.shopName}
                                            width={32}
                                            height={32}
                                            className="object-cover"
                                            unoptimized
                                        />
                                    </div>
                                    <div className="text-left">
                                        <h3 className="text-white text-sm font-medium truncate">{stream.shopName}</h3>
                                        <p className="text-gray-400 text-xs truncate">{stream.title}</p>
                                    </div>
                                </div>
                            </button>
                        ))}

                        <h2 className="text-white font-medium pt-4">Sắp phát</h2>
                        {LIVE_STREAMS.filter(s => !s.isLive).map(stream => (
                            <div key={stream.id} className="bg-gray-800 rounded-lg overflow-hidden opacity-60">
                                <div className="relative aspect-video">
                                    <Image
                                        src={stream.thumbnail}
                                        alt={stream.title}
                                        fill
                                        className="object-cover grayscale"
                                        unoptimized
                                    />
                                    <div className="absolute inset-0 flex items-center justify-center bg-black/50">
                                        <span className="text-white text-sm">Lịch: 15:00 hôm nay</span>
                                    </div>
                                </div>
                                <div className="p-3 flex items-center gap-2">
                                    <div className="w-8 h-8 rounded-full overflow-hidden flex-shrink-0">
                                        <Image
                                            src={stream.shopAvatar}
                                            alt={stream.shopName}
                                            width={32}
                                            height={32}
                                            className="object-cover"
                                            unoptimized
                                        />
                                    </div>
                                    <div className="text-left">
                                        <h3 className="text-white text-sm font-medium truncate">{stream.shopName}</h3>
                                        <p className="text-gray-400 text-xs truncate">{stream.title}</p>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}
