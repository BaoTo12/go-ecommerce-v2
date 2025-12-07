'use client';

import React, { useState, useEffect, useRef } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { chatService, ChatConversation, ChatMessage } from '@/services/chatService';

export default function ChatPage() {
    const [conversations, setConversations] = useState<ChatConversation[]>([]);
    const [selectedConv, setSelectedConv] = useState<ChatConversation | null>(null);
    const [message, setMessage] = useState('');
    const [isLoading, setIsLoading] = useState(true);
    const [isSending, setIsSending] = useState(false);
    const messagesEndRef = useRef<HTMLDivElement>(null);

    const loadConversations = async () => {
        const data = await chatService.getConversations();
        setConversations(data);
        setIsLoading(false);
    };

    useEffect(() => {
        loadConversations();
        // Poll for new messages
        const interval = setInterval(loadConversations, 3000);
        return () => clearInterval(interval);
    }, []);

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [selectedConv?.messages]);

    const selectConversation = async (conv: ChatConversation) => {
        setSelectedConv(conv);
        await chatService.markAsRead(conv.id);
        loadConversations();
    };

    const sendMessage = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!message.trim() || !selectedConv) return;

        setIsSending(true);
        await chatService.sendMessage(selectedConv.id, message);
        setMessage('');

        // Refresh to get updated messages
        const updated = await chatService.getConversation(selectedConv.id);
        if (updated) setSelectedConv(updated);
        loadConversations();
        setIsSending(false);
    };

    const formatTime = (timestamp: string) => {
        const date = new Date(timestamp);
        const now = new Date();
        const isToday = date.toDateString() === now.toDateString();

        if (isToday) {
            return date.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' });
        }
        return date.toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit' });
    };

    if (isLoading) {
        return (
            <div className="min-h-screen bg-[#f5f5f5] flex items-center justify-center">
                <div className="loading-spinner" />
            </div>
        );
    }

    return (
        <div className="min-h-[calc(100vh-200px)] bg-white rounded-sm shadow-sm overflow-hidden animate-fade-in">
            <div className="grid grid-cols-3 h-[600px]">
                {/* Conversation List */}
                <div className="col-span-1 border-r overflow-y-auto">
                    <div className="p-4 border-b bg-gray-50">
                        <h2 className="font-medium">💬 Tin Nhắn</h2>
                    </div>

                    {conversations.length === 0 ? (
                        <div className="p-8 text-center text-gray-500">
                            <div className="text-4xl mb-2">💬</div>
                            <p>Chưa có cuộc trò chuyện nào</p>
                        </div>
                    ) : (
                        <div className="divide-y">
                            {conversations.map(conv => (
                                <button
                                    key={conv.id}
                                    onClick={() => selectConversation(conv)}
                                    className={`w-full p-4 flex items-start gap-3 hover:bg-gray-50 transition-colors text-left ${selectedConv?.id === conv.id ? 'bg-[#fef6f5]' : ''
                                        }`}
                                >
                                    <div className="relative">
                                        <div className="w-12 h-12 rounded-full overflow-hidden">
                                            <Image
                                                src={conv.participantAvatar}
                                                alt={conv.participantName}
                                                width={48}
                                                height={48}
                                                className="object-cover"
                                                unoptimized
                                            />
                                        </div>
                                        <span className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 border-2 border-white rounded-full" />
                                    </div>
                                    <div className="flex-1 min-w-0">
                                        <div className="flex items-center justify-between">
                                            <span className="font-medium text-sm truncate">{conv.participantName}</span>
                                            <span className="text-xs text-gray-400">{formatTime(conv.lastMessageTime)}</span>
                                        </div>
                                        <p className="text-xs text-gray-500 truncate">{conv.lastMessage || 'Bắt đầu cuộc trò chuyện'}</p>
                                    </div>
                                    {conv.unreadCount > 0 && (
                                        <span className="w-5 h-5 bg-[#ee4d2d] text-white text-xs rounded-full flex items-center justify-center">
                                            {conv.unreadCount}
                                        </span>
                                    )}
                                </button>
                            ))}
                        </div>
                    )}
                </div>

                {/* Chat Area */}
                <div className="col-span-2 flex flex-col">
                    {selectedConv ? (
                        <>
                            {/* Chat Header */}
                            <div className="p-4 border-b flex items-center gap-3">
                                <div className="w-10 h-10 rounded-full overflow-hidden">
                                    <Image
                                        src={selectedConv.participantAvatar}
                                        alt={selectedConv.participantName}
                                        width={40}
                                        height={40}
                                        className="object-cover"
                                        unoptimized
                                    />
                                </div>
                                <div>
                                    <h3 className="font-medium">{selectedConv.participantName}</h3>
                                    <p className="text-xs text-green-600 flex items-center gap-1">
                                        <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
                                        Đang hoạt động
                                    </p>
                                </div>
                                <div className="ml-auto flex gap-2">
                                    <button className="p-2 hover:bg-gray-100 rounded-full">📞</button>
                                    <button className="p-2 hover:bg-gray-100 rounded-full">📹</button>
                                    <button className="p-2 hover:bg-gray-100 rounded-full">⋮</button>
                                </div>
                            </div>

                            {/* Messages */}
                            <div className="flex-1 overflow-y-auto p-4 space-y-4 bg-gray-50">
                                {selectedConv.messages.map((msg, i) => {
                                    const isMe = msg.senderId === 'u1';
                                    return (
                                        <div
                                            key={msg.id}
                                            className={`flex items-end gap-2 animate-fade-in-up ${isMe ? 'flex-row-reverse' : ''}`}
                                            style={{ animationDelay: `${i * 30}ms` }}
                                        >
                                            {!isMe && (
                                                <div className="w-8 h-8 rounded-full overflow-hidden flex-shrink-0">
                                                    <Image
                                                        src={msg.senderAvatar}
                                                        alt={msg.senderName}
                                                        width={32}
                                                        height={32}
                                                        className="object-cover"
                                                        unoptimized
                                                    />
                                                </div>
                                            )}
                                            <div className={`max-w-[70%] ${isMe ? 'text-right' : ''}`}>
                                                <div className={`inline-block px-4 py-2 rounded-2xl text-sm ${isMe
                                                        ? 'bg-[#ee4d2d] text-white rounded-br-sm'
                                                        : 'bg-white shadow-sm rounded-bl-sm'
                                                    }`}>
                                                    {msg.content}
                                                </div>
                                                <p className="text-xs text-gray-400 mt-1">
                                                    {formatTime(msg.timestamp)}
                                                    {isMe && msg.read && ' ✓✓'}
                                                </p>
                                            </div>
                                        </div>
                                    );
                                })}
                                <div ref={messagesEndRef} />
                            </div>

                            {/* Message Input */}
                            <form onSubmit={sendMessage} className="p-4 border-t bg-white">
                                <div className="flex items-center gap-2">
                                    <button type="button" className="p-2 hover:bg-gray-100 rounded-full text-gray-500">
                                        😊
                                    </button>
                                    <button type="button" className="p-2 hover:bg-gray-100 rounded-full text-gray-500">
                                        📷
                                    </button>
                                    <input
                                        type="text"
                                        value={message}
                                        onChange={(e) => setMessage(e.target.value)}
                                        placeholder="Nhập tin nhắn..."
                                        className="flex-1 px-4 py-2 bg-gray-100 rounded-full text-sm outline-none focus:ring-2 focus:ring-[#ee4d2d]/20"
                                    />
                                    <button
                                        type="submit"
                                        disabled={!message.trim() || isSending}
                                        className={`p-2 rounded-full transition-colors ${message.trim()
                                                ? 'bg-[#ee4d2d] text-white hover:opacity-90'
                                                : 'bg-gray-100 text-gray-400'
                                            }`}
                                    >
                                        {isSending ? '⏳' : '➤'}
                                    </button>
                                </div>
                            </form>
                        </>
                    ) : (
                        <div className="flex-1 flex items-center justify-center text-gray-500">
                            <div className="text-center">
                                <div className="text-6xl mb-4">💬</div>
                                <p>Chọn một cuộc trò chuyện để bắt đầu</p>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
