'use client';

import React, { useState, useEffect, useRef } from 'react';

interface Stream {
    id: string;
    host: string;
    title: string;
    viewers: number;
    thumbnail: string;
    products: number;
}

interface ChatMessage {
    id: string;
    user: string;
    message: string;
    isGift?: boolean;
}

export default function LivestreamPage() {
    const [streams, setStreams] = useState<Stream[]>([]);
    const [selectedStream, setSelectedStream] = useState<Stream | null>(null);
    const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);
    const [newMessage, setNewMessage] = useState('');
    const chatRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        setStreams([
            { id: '1', host: 'TechStore Official', title: '🔥 Sale iPhone 15 - Giảm 5 Triệu!', viewers: 15234, thumbnail: '📱', products: 12 },
            { id: '2', host: 'Fashion Queen', title: '👗 Thời Trang Mùa Đông - Freeship', viewers: 8721, thumbnail: '👗', products: 25 },
            { id: '3', host: 'BeautyLive', title: '💄 Review Mỹ Phẩm Hàn - Tặng Voucher', viewers: 12543, thumbnail: '💄', products: 18 },
            { id: '4', host: 'HomeDecor.vn', title: '🏠 Đồ Gia Dụng Thông Minh', viewers: 5432, thumbnail: '🏠', products: 30 },
            { id: '5', host: 'SneakerHub', title: '👟 Unbox Giày Limited Edition', viewers: 21456, thumbnail: '👟', products: 8 },
            { id: '6', host: 'GadgetWorld', title: '🎧 Deal Phụ Kiện Công Nghệ', viewers: 9876, thumbnail: '🎧', products: 45 },
        ]);
    }, []);

    useEffect(() => {
        if (!selectedStream) return;

        const messages = [
            '🔥 Sản phẩm này có freeship không shop?',
            'Giá tốt quá!',
            '❤️❤️❤️',
            'Cho xem kỹ hơn được không ạ?',
            'Đã đặt hàng!',
            '⭐⭐⭐⭐⭐',
            'Sale cho em voucher với shop ơi',
            'Size này còn màu khác không?',
        ];
        const users = ['Minh Anh', 'Hoàng Long', 'Thu Hà', 'Đức Minh', 'Lan Anh', 'Việt Hà'];

        const interval = setInterval(() => {
            setChatMessages(prev => [
                ...prev.slice(-30),
                {
                    id: crypto.randomUUID(),
                    user: users[Math.floor(Math.random() * users.length)],
                    message: messages[Math.floor(Math.random() * messages.length)],
                    isGift: Math.random() > 0.9,
                },
            ]);
        }, 1500);

        return () => clearInterval(interval);
    }, [selectedStream]);

    useEffect(() => {
        if (chatRef.current) {
            chatRef.current.scrollTop = chatRef.current.scrollHeight;
        }
    }, [chatMessages]);

    const sendMessage = () => {
        if (!newMessage.trim()) return;
        setChatMessages(prev => [...prev, { id: crypto.randomUUID(), user: 'Bạn', message: newMessage }]);
        setNewMessage('');
    };

    if (selectedStream) {
        return (
            <div className="min-h-screen bg-black">
                <div className="container mx-auto">
                    <div className="grid lg:grid-cols-3 gap-0">
                        {/* Video Area */}
                        <div className="lg:col-span-2 relative">
                            <button
                                onClick={() => setSelectedStream(null)}
                                className="absolute top-4 left-4 z-10 bg-black/50 text-white px-3 py-1 rounded text-sm hover:bg-black/70"
                            >
                                ← Quay lại
                            </button>

                            <div className="aspect-video bg-gradient-to-br from-gray-900 to-gray-800 flex items-center justify-center relative">
                                <span className="text-9xl">{selectedStream.thumbnail}</span>

                                {/* Live Badge */}
                                <div className="absolute top-4 right-4 bg-[#EE4D2D] text-white px-3 py-1 rounded flex items-center gap-2 text-sm font-bold">
                                    <span className="w-2 h-2 bg-white rounded-full animate-pulse" />
                                    LIVE
                                </div>

                                {/* Viewers */}
                                <div className="absolute bottom-4 left-4 bg-black/60 text-white px-3 py-1 rounded text-sm">
                                    👁️ {selectedStream.viewers.toLocaleString()} đang xem
                                </div>

                                {/* Products Overlay */}
                                <div className="absolute bottom-4 right-4 flex gap-2">
                                    {[1, 2, 3].map(i => (
                                        <div key={i} className="w-16 h-16 bg-white rounded shadow-lg flex items-center justify-center text-2xl cursor-pointer hover:scale-110 transition-transform">
                                            📦
                                        </div>
                                    ))}
                                </div>
                            </div>

                            {/* Stream Info */}
                            <div className="bg-gray-900 p-4">
                                <div className="flex items-center gap-3">
                                    <div className="w-10 h-10 bg-[#EE4D2D] rounded-full flex items-center justify-center text-white font-bold">
                                        {selectedStream.host[0]}
                                    </div>
                                    <div className="flex-1">
                                        <h2 className="text-white font-bold">{selectedStream.title}</h2>
                                        <p className="text-gray-400 text-sm">{selectedStream.host}</p>
                                    </div>
                                    <button className="bg-[#EE4D2D] text-white px-6 py-2 rounded font-bold hover:bg-[#D73211]">
                                        + Theo dõi
                                    </button>
                                </div>
                            </div>
                        </div>

                        {/* Chat Area */}
                        <div className="bg-gray-900 flex flex-col h-[calc(100vh-64px)]">
                            <div className="p-3 border-b border-gray-700">
                                <h3 className="text-white font-bold">💬 Trò chuyện trực tiếp</h3>
                            </div>

                            <div ref={chatRef} className="flex-1 overflow-y-auto p-3 space-y-2">
                                {chatMessages.map(msg => (
                                    <div
                                        key={msg.id}
                                        className={`rounded p-2 text-sm ${msg.isGift
                                                ? 'bg-gradient-to-r from-yellow-500/20 to-orange-500/20 border border-yellow-500/50'
                                                : msg.user === 'Bạn'
                                                    ? 'bg-[#EE4D2D]/20'
                                                    : 'bg-gray-800'
                                            }`}
                                    >
                                        {msg.isGift && <span className="text-yellow-400">🎁 </span>}
                                        <span className="text-[#EE4D2D] font-semibold">{msg.user}: </span>
                                        <span className="text-white">{msg.message}</span>
                                    </div>
                                ))}
                            </div>

                            <div className="p-3 border-t border-gray-700">
                                <div className="flex gap-2">
                                    <input
                                        type="text"
                                        value={newMessage}
                                        onChange={e => setNewMessage(e.target.value)}
                                        onKeyPress={e => e.key === 'Enter' && sendMessage()}
                                        placeholder="Nhập tin nhắn..."
                                        className="flex-1 bg-gray-800 text-white rounded px-3 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-[#EE4D2D]"
                                    />
                                    <button
                                        onClick={sendMessage}
                                        className="bg-[#EE4D2D] text-white px-4 py-2 rounded font-semibold hover:bg-[#D73211]"
                                    >
                                        Gửi
                                    </button>
                                </div>
                                <div className="flex gap-2 mt-2">
                                    <button className="bg-yellow-500/20 text-yellow-400 px-3 py-1 rounded text-xs hover:bg-yellow-500/30">
                                        🎁 Tặng quà
                                    </button>
                                    <button className="bg-red-500/20 text-red-400 px-3 py-1 rounded text-xs hover:bg-red-500/30">
                                        ❤️ Thích
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-[#F5F5F5]">
            {/* Header */}
            <div className="bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] py-4">
                <div className="container mx-auto px-4">
                    <div className="flex items-center gap-3 text-white">
                        <span className="w-3 h-3 bg-white rounded-full animate-pulse" />
                        <h1 className="text-2xl font-bold">Shopee Live</h1>
                        <span className="text-sm opacity-80">| {streams.length} đang phát trực tiếp</span>
                    </div>
                </div>
            </div>

            {/* Streams Grid */}
            <div className="container mx-auto px-4 py-6">
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-3">
                    {streams.map(stream => (
                        <div
                            key={stream.id}
                            onClick={() => setSelectedStream(stream)}
                            className="bg-white rounded overflow-hidden cursor-pointer hover:shadow-xl transition-shadow group"
                        >
                            <div className="relative aspect-[3/4] bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center overflow-hidden">
                                <span className="text-6xl group-hover:scale-110 transition-transform">{stream.thumbnail}</span>

                                {/* Live Badge */}
                                <div className="absolute top-2 left-2 bg-[#EE4D2D] text-white px-2 py-0.5 rounded text-xs font-bold flex items-center gap-1">
                                    <span className="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
                                    LIVE
                                </div>

                                {/* Viewer Count */}
                                <div className="absolute bottom-2 left-2 bg-black/60 text-white px-2 py-0.5 rounded text-xs">
                                    👁️ {stream.viewers.toLocaleString()}
                                </div>

                                {/* Products Count */}
                                <div className="absolute bottom-2 right-2 bg-[#EE4D2D] text-white px-2 py-0.5 rounded text-xs">
                                    🛒 {stream.products}
                                </div>
                            </div>

                            <div className="p-2">
                                <h3 className="text-sm font-semibold line-clamp-2 h-10">{stream.title}</h3>
                                <p className="text-xs text-gray-500 mt-1">{stream.host}</p>
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {/* Features */}
            <div className="container mx-auto px-4 py-8">
                <div className="bg-white rounded p-6">
                    <h2 className="text-lg font-bold text-center mb-6 text-[#EE4D2D]">Tính năng Shopee Live</h2>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
                        <div className="p-4">
                            <div className="text-4xl mb-2">🛍️</div>
                            <h3 className="font-semibold text-sm">Mua Ngay</h3>
                            <p className="text-xs text-gray-500">Thêm vào giỏ khi xem</p>
                        </div>
                        <div className="p-4">
                            <div className="text-4xl mb-2">💬</div>
                            <h3 className="font-semibold text-sm">Chat Trực Tiếp</h3>
                            <p className="text-xs text-gray-500">Hỏi đáp ngay lập tức</p>
                        </div>
                        <div className="p-4">
                            <div className="text-4xl mb-2">🎁</div>
                            <h3 className="font-semibold text-sm">Quà Tặng</h3>
                            <p className="text-xs text-gray-500">Voucher độc quyền</p>
                        </div>
                        <div className="p-4">
                            <div className="text-4xl mb-2">⚡</div>
                            <h3 className="font-semibold text-sm">Flash Deal</h3>
                            <p className="text-xs text-gray-500">Giá sốc live only</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
