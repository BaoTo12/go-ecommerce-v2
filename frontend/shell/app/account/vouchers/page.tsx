'use client';

import React, { useState, useEffect } from 'react';
import { voucherService, Voucher } from '@/services/voucherService';

export default function VouchersPage() {
    const [vouchers, setVouchers] = useState<Voucher[]>([]);
    const [userVouchers, setUserVouchers] = useState<Voucher[]>([]);
    const [activeTab, setActiveTab] = useState<'all' | 'collected'>('collected');
    const [isLoading, setIsLoading] = useState(true);
    const [notification, setNotification] = useState<string | null>(null);

    const loadVouchers = async () => {
        const [all, user] = await Promise.all([
            voucherService.getAllVouchers(),
            voucherService.getUserVouchers(),
        ]);
        setVouchers(all);
        setUserVouchers(user);
        setIsLoading(false);
    };

    useEffect(() => {
        loadVouchers();
    }, []);

    const collectVoucher = async (voucherId: string) => {
        const success = await voucherService.collectVoucher(voucherId);
        if (success) {
            setNotification('🎉 Đã lưu voucher thành công!');
            loadVouchers();
        } else {
            setNotification('❌ Không thể lưu voucher');
        }
        setTimeout(() => setNotification(null), 2000);
    };

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    const displayVouchers = activeTab === 'collected' ? userVouchers : vouchers;

    const getVoucherColor = (type: Voucher['type']) => {
        switch (type) {
            case 'freeship': return 'from-cyan-500 to-blue-500';
            case 'percentage': return 'from-orange-500 to-red-500';
            default: return 'from-[#ee4d2d] to-[#f63]';
        }
    };

    const getVoucherIcon = (type: Voucher['type']) => {
        switch (type) {
            case 'freeship': return '🚚';
            case 'percentage': return '🏷️';
            default: return '🎁';
        }
    };

    if (isLoading) {
        return (
            <div className="flex items-center justify-center h-40">
                <div className="loading-spinner" />
            </div>
        );
    }

    return (
        <div className="animate-fade-in">
            {/* Toast */}
            {notification && <div className="toast toast-success">{notification}</div>}

            {/* Header */}
            <div className="bg-white rounded-sm shadow-sm p-4 mb-4">
                <h1 className="text-lg font-medium">Kho Voucher</h1>
                <p className="text-sm text-gray-500">Bạn có {userVouchers.length} voucher</p>
            </div>

            {/* Tabs */}
            <div className="bg-white rounded-sm shadow-sm mb-4">
                <div className="flex">
                    <button
                        onClick={() => setActiveTab('collected')}
                        className={`flex-1 py-3 text-sm font-medium border-b-2 transition-colors ${activeTab === 'collected'
                                ? 'text-[#ee4d2d] border-[#ee4d2d]'
                                : 'text-gray-500 border-transparent'
                            }`}
                    >
                        Voucher Của Tôi ({userVouchers.length})
                    </button>
                    <button
                        onClick={() => setActiveTab('all')}
                        className={`flex-1 py-3 text-sm font-medium border-b-2 transition-colors ${activeTab === 'all'
                                ? 'text-[#ee4d2d] border-[#ee4d2d]'
                                : 'text-gray-500 border-transparent'
                            }`}
                    >
                        Khám Phá ({vouchers.filter(v => !v.isCollected).length})
                    </button>
                </div>
            </div>

            {/* Voucher List */}
            {displayVouchers.length === 0 ? (
                <div className="bg-white rounded-sm shadow-sm p-12 text-center">
                    <div className="text-5xl mb-4">🎟️</div>
                    <p className="text-gray-500">Chưa có voucher nào</p>
                </div>
            ) : (
                <div className="space-y-3">
                    {displayVouchers.map((voucher, index) => {
                        const isExpired = new Date(voucher.expiresAt) < new Date();
                        const daysLeft = Math.ceil((new Date(voucher.expiresAt).getTime() - Date.now()) / (1000 * 60 * 60 * 24));

                        return (
                            <div
                                key={voucher.id}
                                className={`bg-white rounded-sm shadow-sm overflow-hidden flex animate-fade-in-up ${isExpired ? 'opacity-60' : ''
                                    }`}
                                style={{ animationDelay: `${index * 50}ms` }}
                            >
                                {/* Left side */}
                                <div className={`w-28 bg-gradient-to-br ${getVoucherColor(voucher.type)} text-white p-4 flex flex-col items-center justify-center relative`}>
                                    <div className="text-3xl mb-1">{getVoucherIcon(voucher.type)}</div>
                                    {voucher.type === 'freeship' ? (
                                        <div className="text-center">
                                            <div className="text-xs">Miễn phí</div>
                                            <div className="text-sm font-bold">Vận chuyển</div>
                                        </div>
                                    ) : voucher.type === 'percentage' ? (
                                        <div className="text-center">
                                            <div className="text-2xl font-bold">{voucher.value}%</div>
                                            <div className="text-xs">Giảm</div>
                                        </div>
                                    ) : (
                                        <div className="text-center">
                                            <div className="text-lg font-bold">₫{formatPrice(voucher.value)}</div>
                                            <div className="text-xs">Giảm</div>
                                        </div>
                                    )}
                                    {/* Decorative dots */}
                                    <div className="absolute right-0 top-0 bottom-0 flex flex-col justify-around">
                                        {[...Array(5)].map((_, i) => (
                                            <div key={i} className="w-3 h-3 bg-white rounded-full -mr-1.5" />
                                        ))}
                                    </div>
                                </div>

                                {/* Right side */}
                                <div className="flex-1 p-4 flex items-center justify-between">
                                    <div>
                                        <h3 className="font-medium text-sm">{voucher.description}</h3>
                                        {voucher.shopName && (
                                            <p className="text-xs text-gray-500 mt-1">🏪 {voucher.shopName}</p>
                                        )}
                                        {voucher.minOrder > 0 && (
                                            <p className="text-xs text-gray-400 mt-1">
                                                Đơn tối thiểu ₫{formatPrice(voucher.minOrder)}
                                            </p>
                                        )}
                                        <div className="flex items-center gap-3 mt-2">
                                            <span className="text-xs text-gray-400">
                                                HSD: {new Date(voucher.expiresAt).toLocaleDateString('vi-VN')}
                                            </span>
                                            {!isExpired && daysLeft <= 7 && (
                                                <span className="text-xs text-orange-500">
                                                    ⏰ Còn {daysLeft} ngày
                                                </span>
                                            )}
                                        </div>
                                    </div>

                                    <div className="flex-shrink-0">
                                        {voucher.isCollected ? (
                                            <button className="px-4 py-2 border border-gray-300 text-gray-400 text-sm rounded-sm cursor-default">
                                                Đã lưu
                                            </button>
                                        ) : isExpired ? (
                                            <button className="px-4 py-2 border border-gray-300 text-gray-400 text-sm rounded-sm cursor-not-allowed">
                                                Hết hạn
                                            </button>
                                        ) : (
                                            <button
                                                onClick={() => collectVoucher(voucher.id)}
                                                className="px-4 py-2 bg-[#ee4d2d] text-white text-sm rounded-sm hover:opacity-90 transition-all"
                                            >
                                                Lưu
                                            </button>
                                        )}
                                    </div>
                                </div>
                            </div>
                        );
                    })}
                </div>
            )}

            {/* Voucher Code Input */}
            <div className="bg-white rounded-sm shadow-sm p-4 mt-4">
                <h3 className="text-sm font-medium mb-3">Nhập Mã Voucher</h3>
                <div className="flex gap-2">
                    <input
                        type="text"
                        placeholder="Nhập mã voucher"
                        className="flex-1 border px-4 py-2 text-sm outline-none focus:border-[#ee4d2d] rounded-sm uppercase"
                    />
                    <button className="px-6 py-2 bg-[#ee4d2d] text-white text-sm hover:opacity-90 rounded-sm">
                        Áp dụng
                    </button>
                </div>
            </div>
        </div>
    );
}
