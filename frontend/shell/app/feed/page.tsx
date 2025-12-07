'use client';

import React, { useState, useRef, useEffect } from 'react';
import Image from 'next/image';
import Link from 'next/link';

interface FeedItem {
    id: string;
    type: 'video' | 'image';
    url: string;
    thumbnail: string;
    shop: {
        name: string;
        avatar: string;
        isVerified: boolean;
    };
    product: {
        id: string;
        name: string;
        price: number;
    };
    likes: number;
    comments: number;
    shares: number;
    description: string;
    isLiked: boolean;
}

const FEED_DATA: FeedItem[] = [
    {
        id: 'f1',
        type: 'video',
        url: 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=400',
        thumbnail: 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=400',
        shop: { name: 'Apple Store', avatar: 'https://ui-avatars.com/api/?name=Apple&background=000&color=fff', isVerified: true },
        product: { id: 'p1', name: 'iPhone 15 Pro Max', price: 29990000 },
        likes: 12500,
        comments: 890,
        shares: 234,
        description: '✨ Unboxing iPhone 15 Pro Max Titan Xanh! Camera 48MP siêu đỉnh 📸 #iPhone15 #Apple #Shopee',
        isLiked: false,
    },
    {
        id: 'f2',
        type: 'video',
        url: 'https://images.unsplash.com/photo-1556821840-3a63f95609a7?w=400',
        thumbnail: 'https://images.unsplash.com/photo-1556821840-3a63f95609a7?w=400',
        shop: { name: 'Fashion VN', avatar: 'https://ui-avatars.com/api/?name=FV&background=e91e63&color=fff', isVerified: false },
        product: { id: 'p4', name: 'Áo Hoodie Premium', price: 199000 },
        likes: 8900,
        comments: 456,
        shares: 123,
        description: '🔥 Mix & Match hoodie cực chất cho mùa đông! Giảm 50% hôm nay thôi! #Fashion #Hoodie #Shopee',
        isLiked: true,
    },
    {
        id: 'f3',
        type: 'video',
        url: 'https://images.unsplash.com/photo-1586495777744-4413f21062fa?w=400',
        thumbnail: 'https://images.unsplash.com/photo-1586495777744-4413f21062fa?w=400',
        shop: { name: 'Dior Beauty', avatar: 'https://ui-avatars.com/api/?name=Dior&background=9c27b0&color=fff', isVerified: true },
        product: { id: 'p6', name: 'Son Dior Addict', price: 950000 },
        likes: 25600,
        comments: 1234,
        shares: 567,
        description: '💄 Review son Dior Addict Lip Glow - Màu tự nhiên siêu xinh! #Dior #Makeup #Shopee',
        isLiked: false,
    },
    {
        id: 'f4',
        type: 'video',
        url: 'https://images.unsplash.com/photo-1600269452121-4f2416e55c28?w=400',
        thumbnail: 'https://images.unsplash.com/photo-1600269452121-4f2416e55c28?w=400',
        shop: { name: 'Nike Official', avatar: 'https://ui-avatars.com/api/?name=Nike&background=000&color=fff', isVerified: true },
        product: { id: 'p5', name: 'Nike Air Force 1', price: 2590000 },
        likes: 18200,
        comments: 789,
        shares: 345,
        description: '👟 Unbox Nike Air Force 1 - Huyền thoại không bao giờ lỗi mốt! #Nike #Sneaker #Shopee',
        isLiked: false,
    },
];

export default function ShopeeFeed() {
    const [feedItems, setFeedItems] = useState(FEED_DATA);
    const [currentIndex, setCurrentIndex] = useState(0);
    const containerRef = useRef<HTMLDivElement>(null);

    const formatNumber = (num: number) => {
        if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
        if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
        return num.toString();
    };

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    const toggleLike = (id: string) => {
        setFeedItems(items =>
            items.map(item =>
                item.id === id
                    ? { ...item, isLiked: !item.isLiked, likes: item.isLiked ? item.likes - 1 : item.likes + 1 }
                    : item
            )
        );
    };

    // Vertical scroll snap
    useEffect(() => {
        const container = containerRef.current;
        if (!container) return;

        const handleScroll = () => {
            const scrollTop = container.scrollTop;
            const itemHeight = container.clientHeight;
            const index = Math.round(scrollTop / itemHeight);
            setCurrentIndex(index);
        };

        container.addEventListener('scroll', handleScroll);
        return () => container.removeEventListener('scroll', handleScroll);
    }, []);

    return (
        <div className="h-screen bg-black overflow-hidden">
            {/* Header */}
            <div className="absolute top-0 left-0 right-0 z-50 bg-gradient-to-b from-black/80 to-transparent p-4">
                <div className="flex items-center justify-between">
                    <Link href="/" className="text-white text-lg font-bold">
                        ← Shopee Feed
                    </Link>
                    <div className="flex gap-4">
                        <button className="text-white">🔍</button>
                        <button className="text-white">📷</button>
                    </div>
                </div>
                <div className="flex gap-4 mt-4 justify-center">
                    <button className="text-white/60 text-sm">Theo dõi</button>
                    <button className="text-white text-sm font-bold border-b-2 border-white pb-1">Dành cho bạn</button>
                </div>
            </div>

            {/* Feed Container */}
            <div
                ref={containerRef}
                className="h-full overflow-y-auto snap-y snap-mandatory scroll-smooth"
                style={{ scrollSnapType: 'y mandatory' }}
            >
                {feedItems.map((item, index) => (
                    <div
                        key={item.id}
                        className="h-screen w-full relative snap-start"
                        style={{ scrollSnapAlign: 'start' }}
                    >
                        {/* Background */}
                        <div className="absolute inset-0">
                            <Image
                                src={item.thumbnail}
                                alt=""
                                fill
                                className="object-cover"
                                unoptimized
                                priority={index === 0}
                            />
                            <div className="absolute inset-0 bg-black/20" />
                        </div>

                        {/* Play Button */}
                        <div className="absolute inset-0 flex items-center justify-center">
                            <button className="w-20 h-20 bg-white/20 backdrop-blur-sm rounded-full flex items-center justify-center">
                                <span className="text-white text-4xl ml-2">▶</span>
                            </button>
                        </div>

                        {/* Right Actions */}
                        <div className="absolute right-4 bottom-32 flex flex-col items-center gap-6">
                            {/* Shop Avatar */}
                            <div className="relative">
                                <div className="w-12 h-12 rounded-full border-2 border-white overflow-hidden">
                                    <Image src={item.shop.avatar} alt="" fill className="object-cover" unoptimized />
                                </div>
                                <button className="absolute -bottom-2 left-1/2 -translate-x-1/2 w-6 h-6 bg-[#ee4d2d] rounded-full text-white text-xs flex items-center justify-center">
                                    +
                                </button>
                            </div>

                            {/* Like */}
                            <button onClick={() => toggleLike(item.id)} className="flex flex-col items-center">
                                <span className={`text-3xl ${item.isLiked ? 'text-red-500' : 'text-white'}`}>
                                    {item.isLiked ? '❤️' : '🤍'}
                                </span>
                                <span className="text-white text-xs mt-1">{formatNumber(item.likes)}</span>
                            </button>

                            {/* Comment */}
                            <button className="flex flex-col items-center">
                                <span className="text-3xl text-white">💬</span>
                                <span className="text-white text-xs mt-1">{formatNumber(item.comments)}</span>
                            </button>

                            {/* Share */}
                            <button className="flex flex-col items-center">
                                <span className="text-3xl text-white">↗️</span>
                                <span className="text-white text-xs mt-1">{formatNumber(item.shares)}</span>
                            </button>

                            {/* Product */}
                            <Link href={`/products/${item.product.id}`} className="flex flex-col items-center">
                                <span className="text-3xl">🛒</span>
                                <span className="text-white text-xs mt-1">Mua</span>
                            </Link>
                        </div>

                        {/* Bottom Info */}
                        <div className="absolute left-4 right-20 bottom-8">
                            {/* Shop Name */}
                            <div className="flex items-center gap-2 mb-2">
                                <span className="text-white font-bold">{item.shop.name}</span>
                                {item.shop.isVerified && <span className="text-blue-400">✓</span>}
                            </div>

                            {/* Description */}
                            <p className="text-white text-sm mb-3 line-clamp-2">{item.description}</p>

                            {/* Product Card */}
                            <Link
                                href={`/products/${item.product.id}`}
                                className="inline-flex items-center gap-3 bg-white/10 backdrop-blur-sm rounded-lg p-2 pr-4"
                            >
                                <div className="w-12 h-12 bg-white rounded overflow-hidden relative">
                                    <Image src={item.thumbnail} alt="" fill className="object-cover" unoptimized />
                                </div>
                                <div>
                                    <p className="text-white text-sm font-medium line-clamp-1">{item.product.name}</p>
                                    <p className="text-[#ee4d2d] text-sm font-bold">₫{formatPrice(item.product.price)}</p>
                                </div>
                            </Link>
                        </div>

                        {/* Progress Indicator */}
                        <div className="absolute right-2 top-1/2 -translate-y-1/2 flex flex-col gap-1">
                            {feedItems.map((_, i) => (
                                <div
                                    key={i}
                                    className={`w-1 h-6 rounded-full transition-all ${i === index ? 'bg-white' : 'bg-white/30'
                                        }`}
                                />
                            ))}
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}
