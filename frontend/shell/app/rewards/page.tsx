'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { coinsService, CoinState } from '@/services/coinsService';

export default function ShopeeCoinsPage() {
    const [coinState, setCoinState] = useState<CoinState | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [isSpinning, setIsSpinning] = useState(false);
    const [spinResult, setSpinResult] = useState<{ reward: number; type: string } | null>(null);
    const [checkInResult, setCheckInResult] = useState<{ reward: number; streak: number } | null>(null);
    const [notification, setNotification] = useState<string | null>(null);

    const loadCoins = async () => {
        const state = await coinsService.getState();
        setCoinState(state);
        setIsLoading(false);
    };

    useEffect(() => {
        loadCoins();
    }, []);

    const handleCheckIn = async () => {
        const result = await coinsService.dailyCheckIn();
        if (result.success) {
            setCheckInResult({ reward: result.reward, streak: result.streak });
            loadCoins();
        } else {
            setNotification(result.error || 'Đã điểm danh rồi!');
            setTimeout(() => setNotification(null), 2000);
        }
    };

    const handleSpin = async () => {
        setIsSpinning(true);
        setSpinResult(null);
        const result = await coinsService.spinWheel();
        setSpinResult(result);
        setIsSpinning(false);
        loadCoins();
    };

    const formatNumber = (num: number) => new Intl.NumberFormat('vi-VN').format(num);

    if (isLoading || !coinState) {
        return (
            <div className="min-h-screen bg-gradient-to-b from-[#ee4d2d] to-[#ff8f70] flex items-center justify-center">
                <div className="loading-spinner" style={{ borderTopColor: 'white' }} />
            </div>
        );
    }

    const hasCheckedIn = coinsService.hasCheckedInToday();

    return (
        <div className="min-h-screen bg-gradient-to-b from-[#ee4d2d] to-[#ff8f70]">
            {/* Toast */}
            {notification && <div className="toast toast-error">{notification}</div>}

            {/* Check-in Success Modal */}
            {checkInResult && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 animate-fade-in">
                    <div className="bg-white rounded-lg p-6 max-w-sm mx-4 text-center animate-fade-in-up">
                        <div className="text-6xl mb-4 animate-bounce">🎉</div>
                        <h3 className="text-xl font-bold text-[#ee4d2d]">+{checkInResult.reward} Xu</h3>
                        <p className="text-gray-600 mt-2">Điểm danh ngày {checkInResult.streak} liên tiếp!</p>
                        <button
                            onClick={() => setCheckInResult(null)}
                            className="mt-4 px-6 py-2 bg-[#ee4d2d] text-white rounded-full hover:opacity-90"
                        >
                            Tuyệt vời!
                        </button>
                    </div>
                </div>
            )}

            {/* Spin Result Modal */}
            {spinResult && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 animate-fade-in">
                    <div className="bg-white rounded-lg p-6 max-w-sm mx-4 text-center animate-fade-in-up">
                        <div className="text-6xl mb-4">
                            {spinResult.type === 'coins' ? '🪙' : spinResult.type === 'voucher' ? '🎟️' : '😅'}
                        </div>
                        <h3 className="text-xl font-bold">
                            {spinResult.type === 'coins' && `+${spinResult.reward} Xu!`}
                            {spinResult.type === 'voucher' && 'Voucher ₫10.000!'}
                            {spinResult.type === 'nothing' && 'Chúc bạn may mắn lần sau!'}
                        </h3>
                        <button
                            onClick={() => setSpinResult(null)}
                            className="mt-4 px-6 py-2 bg-[#ee4d2d] text-white rounded-full hover:opacity-90"
                        >
                            OK
                        </button>
                    </div>
                </div>
            )}

            <div className="container mx-auto px-4 py-6">
                {/* Header */}
                <div className="text-white text-center mb-6">
                    <h1 className="text-2xl font-bold flex items-center justify-center gap-2">
                        <span className="text-3xl">🪙</span> Shopee Xu
                    </h1>
                </div>

                {/* Balance Card */}
                <div className="bg-white rounded-2xl p-6 shadow-lg mb-6 animate-fade-in-up">
                    <div className="text-center">
                        <div className="text-sm text-gray-500">Số dư hiện tại</div>
                        <div className="text-4xl font-bold text-[#ee4d2d] flex items-center justify-center gap-2 my-2">
                            <span className="text-3xl">🪙</span>
                            {formatNumber(coinState.balance)}
                        </div>
                        <div className="text-xs text-gray-400">
                            Tổng xu đã tích lũy: {formatNumber(coinState.lifetime)}
                        </div>
                        {coinState.expiring && (
                            <div className="text-xs text-orange-500 mt-2">
                                ⚠️ {formatNumber(coinState.expiring.amount)} xu sắp hết hạn ({coinState.expiring.date})
                            </div>
                        )}
                    </div>

                    <div className="mt-4 pt-4 border-t flex justify-around">
                        <Link href="/rewards/history" className="text-center">
                            <div className="text-sm text-gray-500">Lịch sử</div>
                            <div className="text-[#ee4d2d]">📋</div>
                        </Link>
                        <Link href="/rewards/earn" className="text-center">
                            <div className="text-sm text-gray-500">Kiếm xu</div>
                            <div className="text-[#ee4d2d]">💰</div>
                        </Link>
                        <Link href="/rewards/spend" className="text-center">
                            <div className="text-sm text-gray-500">Đổi xu</div>
                            <div className="text-[#ee4d2d]">🎁</div>
                        </Link>
                    </div>
                </div>

                {/* Daily Check-in */}
                <div className="bg-white rounded-2xl p-6 shadow-lg mb-6 animate-fade-in-up" style={{ animationDelay: '100ms' }}>
                    <h3 className="font-bold flex items-center gap-2 mb-4">
                        📅 Điểm Danh Hàng Ngày
                    </h3>

                    <div className="flex justify-between mb-4">
                        {[1, 2, 3, 4, 5, 6, 7].map(day => {
                            const isCompleted = day <= coinState.dailyCheckIn.streak;
                            const isToday = day === coinState.dailyCheckIn.streak + 1;
                            return (
                                <div key={day} className="text-center">
                                    <div className={`w-10 h-10 rounded-full flex items-center justify-center text-sm font-medium ${isCompleted
                                            ? 'bg-[#ee4d2d] text-white'
                                            : isToday
                                                ? 'border-2 border-[#ee4d2d] text-[#ee4d2d] animate-pulse'
                                                : 'bg-gray-100 text-gray-400'
                                        }`}>
                                        {isCompleted ? '✓' : `+${5 + (day - 1) * 2}`}
                                    </div>
                                    <div className="text-xs text-gray-400 mt-1">Ngày {day}</div>
                                </div>
                            );
                        })}
                    </div>

                    <button
                        onClick={handleCheckIn}
                        disabled={hasCheckedIn}
                        className={`w-full py-3 rounded-full font-medium transition-all ${hasCheckedIn
                                ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                                : 'bg-gradient-to-r from-[#ee4d2d] to-[#ff8f70] text-white hover:opacity-90'
                            }`}
                    >
                        {hasCheckedIn ? '✓ Đã điểm danh hôm nay' : `Điểm danh nhận +${coinState.dailyCheckIn.todayReward} xu`}
                    </button>
                </div>

                {/* Spin Wheel */}
                <div className="bg-white rounded-2xl p-6 shadow-lg mb-6 animate-fade-in-up" style={{ animationDelay: '200ms' }}>
                    <h3 className="font-bold flex items-center gap-2 mb-4">
                        🎰 Vòng Quay May Mắn
                    </h3>

                    <div className="text-center">
                        <div className={`w-40 h-40 mx-auto rounded-full bg-gradient-to-r from-yellow-400 via-red-500 to-pink-500 flex items-center justify-center text-white text-4xl ${isSpinning ? 'animate-spin' : ''
                            }`} style={{ animationDuration: '0.5s' }}>
                            🎡
                        </div>

                        <button
                            onClick={handleSpin}
                            disabled={isSpinning}
                            className={`mt-4 px-8 py-3 rounded-full font-medium transition-all ${isSpinning
                                    ? 'bg-gray-100 text-gray-400 cursor-wait'
                                    : 'bg-gradient-to-r from-[#ee4d2d] to-[#ff8f70] text-white hover:opacity-90'
                                }`}
                        >
                            {isSpinning ? 'Đang quay...' : 'Quay ngay (1 lượt/ngày)'}
                        </button>
                    </div>

                    <div className="mt-4 pt-4 border-t grid grid-cols-3 gap-2 text-center text-xs">
                        <div className="p-2 bg-yellow-50 rounded">
                            <div>🪙 10-100 xu</div>
                            <div className="text-gray-400">30%</div>
                        </div>
                        <div className="p-2 bg-pink-50 rounded">
                            <div>🎟️ Voucher</div>
                            <div className="text-gray-400">10%</div>
                        </div>
                        <div className="p-2 bg-gray-50 rounded">
                            <div>😅 Chúc may mắn</div>
                            <div className="text-gray-400">60%</div>
                        </div>
                    </div>
                </div>

                {/* Ways to Earn */}
                <div className="bg-white rounded-2xl p-6 shadow-lg animate-fade-in-up" style={{ animationDelay: '300ms' }}>
                    <h3 className="font-bold flex items-center gap-2 mb-4">
                        💡 Cách Kiếm Xu
                    </h3>

                    <div className="space-y-3">
                        {[
                            { icon: '🛒', title: 'Mua hàng', desc: '1 xu cho mỗi ₫10.000', coins: '~100 xu/đơn' },
                            { icon: '⭐', title: 'Đánh giá sản phẩm', desc: 'Đánh giá kèm ảnh/video', coins: '+50 xu' },
                            { icon: '📅', title: 'Điểm danh', desc: 'Điểm danh hàng ngày', coins: '+5-50 xu' },
                            { icon: '🎮', title: 'Chơi game', desc: 'Tham gia mini game', coins: '+10-200 xu' },
                            { icon: '👥', title: 'Mời bạn bè', desc: 'Bạn bè đăng ký thành công', coins: '+500 xu' },
                        ].map((item, i) => (
                            <div key={i} className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                                <div className="text-2xl">{item.icon}</div>
                                <div className="flex-1">
                                    <div className="font-medium text-sm">{item.title}</div>
                                    <div className="text-xs text-gray-500">{item.desc}</div>
                                </div>
                                <div className="text-[#ee4d2d] text-sm font-medium">{item.coins}</div>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Transaction History */}
                <div className="bg-white rounded-2xl p-6 shadow-lg mt-6 animate-fade-in-up" style={{ animationDelay: '400ms' }}>
                    <h3 className="font-bold flex items-center justify-between mb-4">
                        <span>📋 Lịch Sử Giao Dịch</span>
                        <Link href="/rewards/history" className="text-sm text-[#ee4d2d]">Xem tất cả →</Link>
                    </h3>

                    <div className="divide-y">
                        {coinState.transactions.slice(0, 5).map((tx, i) => (
                            <div key={tx.id} className="py-3 flex items-center justify-between">
                                <div>
                                    <div className="text-sm">{tx.description}</div>
                                    <div className="text-xs text-gray-400">
                                        {new Date(tx.timestamp).toLocaleDateString('vi-VN')}
                                    </div>
                                </div>
                                <div className={`font-medium ${tx.amount >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                                    {tx.amount >= 0 ? '+' : ''}{formatNumber(tx.amount)}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}
