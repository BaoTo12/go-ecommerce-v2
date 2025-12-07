'use client';

import React, { useState, useEffect } from 'react';
import Image from 'next/image';
import Link from 'next/link';

interface AuctionItem {
    id: string;
    product: {
        name: string;
        image: string;
        originalPrice: number;
    };
    currentBid: number;
    startingBid: number;
    bidCount: number;
    endTime: string;
    highestBidder: {
        name: string;
        avatar: string;
    };
    status: 'active' | 'ending' | 'ended';
}

const AUCTION_ITEMS: AuctionItem[] = [
    {
        id: 'a1',
        product: {
            name: 'iPhone 15 Pro Max 512GB Limited Edition',
            image: 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=400',
            originalPrice: 35990000,
        },
        currentBid: 25500000,
        startingBid: 20000000,
        bidCount: 47,
        endTime: new Date(Date.now() + 2 * 60 * 60 * 1000).toISOString(), // 2 hours
        highestBidder: { name: 'Nguyễn V***', avatar: 'https://ui-avatars.com/api/?name=NV' },
        status: 'active',
    },
    {
        id: 'a2',
        product: {
            name: 'Nike Air Jordan 1 Retro High OG Limited',
            image: 'https://images.unsplash.com/photo-1600269452121-4f2416e55c28?w=400',
            originalPrice: 8990000,
        },
        currentBid: 5200000,
        startingBid: 3000000,
        bidCount: 32,
        endTime: new Date(Date.now() + 45 * 60 * 1000).toISOString(), // 45 mins
        highestBidder: { name: 'Trần T***', avatar: 'https://ui-avatars.com/api/?name=TT' },
        status: 'ending',
    },
    {
        id: 'a3',
        product: {
            name: 'MacBook Pro 14" M3 Pro Chip',
            image: 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=400',
            originalPrice: 52990000,
        },
        currentBid: 42000000,
        startingBid: 35000000,
        bidCount: 28,
        endTime: new Date(Date.now() + 5 * 60 * 60 * 1000).toISOString(), // 5 hours
        highestBidder: { name: 'Lê H***', avatar: 'https://ui-avatars.com/api/?name=LH' },
        status: 'active',
    },
];

export default function AuctionPage() {
    const [auctions, setAuctions] = useState(AUCTION_ITEMS);
    const [selectedAuction, setSelectedAuction] = useState<AuctionItem | null>(null);
    const [bidAmount, setBidAmount] = useState('');
    const [timeLeft, setTimeLeft] = useState<Record<string, string>>({});
    const [notification, setNotification] = useState<string | null>(null);

    // Update countdown timers
    useEffect(() => {
        const timer = setInterval(() => {
            const newTimeLeft: Record<string, string> = {};
            auctions.forEach(auction => {
                const end = new Date(auction.endTime).getTime();
                const now = Date.now();
                const diff = end - now;

                if (diff <= 0) {
                    newTimeLeft[auction.id] = 'Đã kết thúc';
                } else {
                    const hours = Math.floor(diff / (1000 * 60 * 60));
                    const mins = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
                    const secs = Math.floor((diff % (1000 * 60)) / 1000);
                    newTimeLeft[auction.id] = `${hours.toString().padStart(2, '0')}:${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
                }
            });
            setTimeLeft(newTimeLeft);
        }, 1000);

        return () => clearInterval(timer);
    }, [auctions]);

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    const placeBid = (auctionId: string) => {
        const auction = auctions.find(a => a.id === auctionId);
        if (!auction) return;

        const amount = parseInt(bidAmount.replace(/\D/g, ''));
        const minBid = auction.currentBid + 100000;

        if (!amount || amount < minBid) {
            setNotification(`❌ Giá đặt tối thiểu là ₫${formatPrice(minBid)}`);
            setTimeout(() => setNotification(null), 3000);
            return;
        }

        setAuctions(prev => prev.map(a =>
            a.id === auctionId
                ? {
                    ...a,
                    currentBid: amount,
                    bidCount: a.bidCount + 1,
                    highestBidder: { name: 'Bạn', avatar: 'https://ui-avatars.com/api/?name=Me' }
                }
                : a
        ));

        setNotification('🎉 Đặt giá thành công! Bạn đang dẫn đầu!');
        setTimeout(() => setNotification(null), 3000);
        setBidAmount('');
        setSelectedAuction(null);
    };

    return (
        <div className="min-h-screen bg-[#1a1a2e]">
            {notification && (
                <div className="fixed top-20 left-1/2 -translate-x-1/2 bg-white rounded-lg shadow-lg px-6 py-3 z-50 animate-fade-in">
                    {notification}
                </div>
            )}

            {/* Bid Modal */}
            {selectedAuction && (
                <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4">
                    <div className="bg-white rounded-2xl max-w-md w-full p-6 animate-fade-in-up">
                        <div className="flex items-center gap-4 mb-4">
                            <div className="w-20 h-20 rounded-lg overflow-hidden relative">
                                <Image src={selectedAuction.product.image} alt="" fill className="object-cover" unoptimized />
                            </div>
                            <div>
                                <h3 className="font-bold line-clamp-2">{selectedAuction.product.name}</h3>
                                <p className="text-[#ee4d2d] font-bold">
                                    Giá hiện tại: ₫{formatPrice(selectedAuction.currentBid)}
                                </p>
                            </div>
                        </div>

                        <div className="mb-4">
                            <label className="block text-sm text-gray-600 mb-2">
                                Giá đặt của bạn (tối thiểu ₫{formatPrice(selectedAuction.currentBid + 100000)})
                            </label>
                            <input
                                type="text"
                                value={bidAmount}
                                onChange={e => setBidAmount(e.target.value)}
                                placeholder="Nhập số tiền..."
                                className="w-full border-2 border-gray-200 rounded-lg px-4 py-3 text-lg font-bold focus:border-[#ee4d2d] outline-none"
                            />
                        </div>

                        <div className="flex gap-3">
                            <button
                                onClick={() => setSelectedAuction(null)}
                                className="flex-1 py-3 border-2 border-gray-200 rounded-lg font-medium"
                            >
                                Hủy
                            </button>
                            <button
                                onClick={() => placeBid(selectedAuction.id)}
                                className="flex-1 py-3 bg-[#ee4d2d] text-white rounded-lg font-medium"
                            >
                                🔨 Đặt giá
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {/* Header */}
            <div className="bg-gradient-to-r from-purple-600 to-pink-500 p-6">
                <div className="container mx-auto">
                    <h1 className="text-3xl font-bold text-white flex items-center gap-3">
                        🔨 Đấu Giá Shopee
                    </h1>
                    <p className="text-white/80 mt-2">Săn deal độc quyền với giá không tưởng!</p>
                </div>
            </div>

            <div className="container mx-auto px-4 py-8">
                {/* Live Auctions */}
                <div className="mb-8">
                    <h2 className="text-xl font-bold text-white mb-4 flex items-center gap-2">
                        <span className="w-3 h-3 bg-red-500 rounded-full animate-pulse" />
                        Đang diễn ra
                    </h2>

                    <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {auctions.map((auction, index) => (
                            <div
                                key={auction.id}
                                className="bg-white rounded-2xl overflow-hidden shadow-lg animate-fade-in-up"
                                style={{ animationDelay: `${index * 100}ms` }}
                            >
                                {/* Image */}
                                <div className="relative aspect-square">
                                    <Image
                                        src={auction.product.image}
                                        alt={auction.product.name}
                                        fill
                                        className="object-cover"
                                        unoptimized
                                    />
                                    <div className={`absolute top-4 left-4 px-3 py-1 rounded-full text-xs font-bold text-white ${auction.status === 'ending' ? 'bg-red-500 animate-pulse' : 'bg-purple-600'
                                        }`}>
                                        {auction.status === 'ending' ? '🔥 Sắp kết thúc' : '⚡ Đang đấu giá'}
                                    </div>
                                    <div className="absolute top-4 right-4 bg-black/70 text-white px-3 py-1 rounded-full text-sm font-mono">
                                        ⏰ {timeLeft[auction.id] || '--:--:--'}
                                    </div>
                                </div>

                                {/* Info */}
                                <div className="p-4">
                                    <h3 className="font-bold line-clamp-2 mb-2">{auction.product.name}</h3>

                                    <div className="flex justify-between items-center mb-3">
                                        <div>
                                            <p className="text-xs text-gray-500">Giá khởi điểm</p>
                                            <p className="text-sm line-through text-gray-400">
                                                ₫{formatPrice(auction.product.originalPrice)}
                                            </p>
                                        </div>
                                        <div className="text-right">
                                            <p className="text-xs text-gray-500">Giá hiện tại</p>
                                            <p className="text-xl font-bold text-[#ee4d2d]">
                                                ₫{formatPrice(auction.currentBid)}
                                            </p>
                                        </div>
                                    </div>

                                    <div className="flex items-center justify-between mb-4">
                                        <div className="flex items-center gap-2">
                                            <div className="w-6 h-6 rounded-full overflow-hidden relative">
                                                <img src={auction.highestBidder.avatar} alt="" className="w-full h-full object-cover" />
                                            </div>
                                            <span className="text-sm text-gray-600">{auction.highestBidder.name}</span>
                                        </div>
                                        <span className="text-sm text-gray-500">{auction.bidCount} lượt đặt</span>
                                    </div>

                                    <button
                                        onClick={() => setSelectedAuction(auction)}
                                        className="w-full py-3 bg-gradient-to-r from-purple-600 to-pink-500 text-white rounded-xl font-bold hover:opacity-90 transition-all"
                                    >
                                        🔨 Đặt giá ngay
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Rules */}
                <div className="bg-white/10 backdrop-blur-sm rounded-2xl p-6">
                    <h3 className="text-white font-bold mb-4">📋 Quy tắc đấu giá</h3>
                    <ul className="text-white/80 text-sm space-y-2">
                        <li>• Mỗi lần đặt giá phải cao hơn giá hiện tại tối thiểu ₫100.000</li>
                        <li>• Người đặt giá cao nhất khi hết thời gian sẽ thắng</li>
                        <li>• Thanh toán trong vòng 24h sau khi thắng đấu giá</li>
                        <li>• Không hoàn tiền cọc nếu từ chối thanh toán</li>
                    </ul>
                </div>
            </div>
        </div>
    );
}
