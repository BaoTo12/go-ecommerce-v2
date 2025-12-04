'use client';

import { useState, useEffect, useRef } from 'react';
import { livestreamApi, Livestream } from '../../lib/api';

export default function LivestreamPage() {
    const [streams, setStreams] = useState<Livestream[]>([]);
    const [selectedStream, setSelectedStream] = useState<Livestream | null>(null);
    const [loading, setLoading] = useState(true);
    const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);
    const [newMessage, setNewMessage] = useState('');
    const chatRef = useRef<HTMLDivElement>(null);

    interface ChatMessage {
        id: string;
        user: string;
        message: string;
        timestamp: Date;
        isGift?: boolean;
    }

    useEffect(() => {
        loadStreams();
        const interval = setInterval(loadStreams, 10000);
        return () => clearInterval(interval);
    }, []);

    useEffect(() => {
        if (selectedStream) {
            // Simulate live chat messages
            const chatInterval = setInterval(() => {
                const randomMessages = [
                    '🔥 Great product!',
                    'How much is shipping?',
                    '❤️❤️❤️',
                    'Can you show the back?',
                    'Just ordered!',
                    '⭐⭐⭐⭐⭐',
                    'Love it!',
                    'What colors available?',
                    '🎉 Amazing deal!',
                ];
                const randomUsers = ['Alice', 'Bob', 'Charlie', 'Diana', 'Eve', 'Frank'];
                setChatMessages((prev) => [
                    ...prev.slice(-50),
                    {
                        id: crypto.randomUUID(),
                        user: randomUsers[Math.floor(Math.random() * randomUsers.length)],
                        message: randomMessages[Math.floor(Math.random() * randomMessages.length)],
                        timestamp: new Date(),
                        isGift: Math.random() > 0.9,
                    },
                ]);
            }, 2000);
            return () => clearInterval(chatInterval);
        }
    }, [selectedStream]);

    useEffect(() => {
        if (chatRef.current) {
            chatRef.current.scrollTop = chatRef.current.scrollHeight;
        }
    }, [chatMessages]);

    const loadStreams = async () => {
        try {
            const data = await livestreamApi.getStreams();
            setStreams(data || []);
        } catch {
            // Mock data
            setStreams([
                {
                    id: 'stream-001',
                    host_id: 'seller-001',
                    title: '🔥 Electronics Flash Sale - iPhone Deals!',
                    status: 'live',
                    viewer_count: 12453,
                    product_ids: ['prod-001', 'prod-002'],
                },
                {
                    id: 'stream-002',
                    host_id: 'seller-002',
                    title: '👗 Fashion Week - New Arrivals',
                    status: 'live',
                    viewer_count: 8721,
                    product_ids: ['prod-003', 'prod-004'],
                },
                {
                    id: 'stream-003',
                    host_id: 'seller-003',
                    title: '🏠 Home & Living Must-Haves',
                    status: 'live',
                    viewer_count: 5432,
                    product_ids: ['prod-005'],
                },
                {
                    id: 'stream-004',
                    host_id: 'seller-004',
                    title: '💄 Beauty Tutorial + Giveaway',
                    status: 'live',
                    viewer_count: 15678,
                    product_ids: ['prod-006', 'prod-007'],
                },
            ]);
        } finally {
            setLoading(false);
        }
    };

    const sendMessage = () => {
        if (!newMessage.trim()) return;
        setChatMessages((prev) => [
            ...prev,
            {
                id: crypto.randomUUID(),
                user: 'You',
                message: newMessage,
                timestamp: new Date(),
            },
        ]);
        setNewMessage('');
    };

    if (loading) {
        return (
            <div className="flex h-96 items-center justify-center">
                <div className="animate-pulse text-2xl">🔴 Loading Live Streams...</div>
            </div>
        );
    }

    if (selectedStream) {
        return (
            <div className="container mx-auto py-4">
                <button
                    onClick={() => setSelectedStream(null)}
                    className="mb-4 flex items-center gap-2 text-blue-600 hover:underline"
                >
                    ← Back to streams
                </button>

                <div className="grid gap-4 lg:grid-cols-3">
                    {/* Video Player */}
                    <div className="lg:col-span-2">
                        <div className="relative aspect-video rounded-xl bg-gradient-to-br from-slate-800 to-slate-900 overflow-hidden">
                            {/* Mock video placeholder */}
                            <div className="absolute inset-0 flex items-center justify-center">
                                <div className="text-center text-white">
                                    <div className="text-8xl mb-4">📺</div>
                                    <p className="text-xl">{selectedStream.title}</p>
                                </div>
                            </div>

                            {/* Live badge */}
                            <div className="absolute top-4 left-4 flex items-center gap-2 rounded-full bg-red-600 px-3 py-1 text-sm font-bold text-white">
                                <span className="h-2 w-2 animate-pulse rounded-full bg-white"></span>
                                LIVE
                            </div>

                            {/* Viewer count */}
                            <div className="absolute top-4 right-4 rounded-full bg-black/50 px-3 py-1 text-sm text-white">
                                👁️ {selectedStream.viewer_count.toLocaleString()}
                            </div>

                            {/* Products overlay */}
                            <div className="absolute bottom-4 left-4 right-4 flex gap-2 overflow-x-auto">
                                {[1, 2, 3].map((i) => (
                                    <div
                                        key={i}
                                        className="flex-shrink-0 rounded-lg bg-white p-2 shadow-lg cursor-pointer hover:scale-105 transition-transform"
                                    >
                                        <div className="w-16 h-16 bg-gray-100 rounded flex items-center justify-center text-2xl">
                                            📱
                                        </div>
                                        <p className="mt-1 text-xs font-semibold text-red-600">$99.99</p>
                                    </div>
                                ))}
                            </div>
                        </div>

                        {/* Stream info */}
                        <div className="mt-4 rounded-xl border bg-white p-4">
                            <div className="flex items-center gap-4">
                                <div className="h-12 w-12 rounded-full bg-gradient-to-br from-pink-500 to-purple-500 flex items-center justify-center text-xl text-white">
                                    👤
                                </div>
                                <div>
                                    <h2 className="font-bold">{selectedStream.title}</h2>
                                    <p className="text-sm text-muted-foreground">Host ID: {selectedStream.host_id}</p>
                                </div>
                                <button className="ml-auto rounded-lg bg-red-600 px-6 py-2 font-semibold text-white hover:bg-red-700">
                                    Follow
                                </button>
                            </div>
                        </div>
                    </div>

                    {/* Chat */}
                    <div className="rounded-xl border bg-white flex flex-col h-[600px]">
                        <div className="border-b p-4 font-bold">💬 Live Chat</div>

                        <div ref={chatRef} className="flex-1 overflow-y-auto p-4 space-y-2">
                            {chatMessages.map((msg) => (
                                <div
                                    key={msg.id}
                                    className={`rounded-lg p-2 ${msg.isGift
                                            ? 'bg-gradient-to-r from-yellow-100 to-orange-100 border border-yellow-300'
                                            : msg.user === 'You'
                                                ? 'bg-blue-100'
                                                : 'bg-gray-100'
                                        }`}
                                >
                                    {msg.isGift && <span className="text-lg">🎁 </span>}
                                    <span className="font-semibold text-blue-600">{msg.user}: </span>
                                    <span>{msg.message}</span>
                                </div>
                            ))}
                        </div>

                        <div className="border-t p-4">
                            <div className="flex gap-2">
                                <input
                                    type="text"
                                    value={newMessage}
                                    onChange={(e) => setNewMessage(e.target.value)}
                                    onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
                                    placeholder="Type a message..."
                                    className="flex-1 rounded-lg border px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                />
                                <button
                                    onClick={sendMessage}
                                    className="rounded-lg bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
                                >
                                    Send
                                </button>
                            </div>
                            <div className="mt-2 flex gap-2">
                                <button className="rounded-full bg-yellow-100 px-3 py-1 text-sm hover:bg-yellow-200">
                                    🎁 Gift
                                </button>
                                <button className="rounded-full bg-red-100 px-3 py-1 text-sm hover:bg-red-200">
                                    ❤️ Like
                                </button>
                                <button className="rounded-full bg-blue-100 px-3 py-1 text-sm hover:bg-blue-200">
                                    📢 Share
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="container mx-auto py-8">
            {/* Header */}
            <div className="mb-8 rounded-xl bg-gradient-to-r from-red-600 to-pink-600 p-8 text-white">
                <div className="flex items-center gap-4">
                    <div className="h-4 w-4 animate-pulse rounded-full bg-white"></div>
                    <h1 className="text-4xl font-bold">Live Shopping</h1>
                </div>
                <p className="mt-2 text-lg opacity-90">
                    Watch live streams and shop directly from your favorite sellers
                </p>
                <div className="mt-4 flex gap-6 text-sm">
                    <span>🔴 {streams.length} streams live</span>
                    <span>👁️ {streams.reduce((s, st) => s + st.viewer_count, 0).toLocaleString()} watching</span>
                </div>
            </div>

            {/* Stream Grid */}
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                {streams.map((stream) => (
                    <div
                        key={stream.id}
                        onClick={() => setSelectedStream(stream)}
                        className="group cursor-pointer rounded-xl border bg-white overflow-hidden transition-transform hover:scale-[1.02] hover:shadow-xl"
                    >
                        {/* Thumbnail */}
                        <div className="relative aspect-video bg-gradient-to-br from-slate-700 to-slate-900">
                            <div className="absolute inset-0 flex items-center justify-center text-6xl">
                                🎥
                            </div>

                            {/* Live badge */}
                            <div className="absolute top-2 left-2 flex items-center gap-1 rounded bg-red-600 px-2 py-0.5 text-xs font-bold text-white">
                                <span className="h-1.5 w-1.5 animate-pulse rounded-full bg-white"></span>
                                LIVE
                            </div>

                            {/* Viewers */}
                            <div className="absolute bottom-2 right-2 rounded bg-black/60 px-2 py-0.5 text-xs text-white">
                                👁️ {stream.viewer_count.toLocaleString()}
                            </div>

                            {/* Overlay on hover */}
                            <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                                <span className="rounded-full bg-white/90 p-4 text-2xl">▶️</span>
                            </div>
                        </div>

                        {/* Info */}
                        <div className="p-4">
                            <h3 className="font-semibold line-clamp-2">{stream.title}</h3>
                            <p className="mt-1 text-sm text-muted-foreground">
                                {stream.product_ids.length} products featured
                            </p>
                        </div>
                    </div>
                ))}
            </div>

            {/* Features */}
            <div className="mt-12 grid gap-6 md:grid-cols-3">
                {[
                    { icon: '🛍️', title: 'Shop While Watching', desc: 'Add products to cart without leaving the stream' },
                    { icon: '💬', title: 'Real-time Chat', desc: 'Interact with hosts and other viewers instantly' },
                    { icon: '🎁', title: 'Exclusive Deals', desc: 'Stream-only discounts and flash sales' },
                ].map((f) => (
                    <div key={f.title} className="rounded-xl border bg-white p-6 text-center">
                        <div className="text-4xl mb-2">{f.icon}</div>
                        <h3 className="font-bold">{f.title}</h3>
                        <p className="text-sm text-muted-foreground">{f.desc}</p>
                    </div>
                ))}
            </div>
        </div>
    );
}
