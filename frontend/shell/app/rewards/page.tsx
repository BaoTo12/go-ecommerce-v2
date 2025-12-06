'use client';

import React, { useState, useEffect, useRef } from 'react';

export default function RewardsPage() {
    const [coins, setCoins] = useState(1250);
    const [isSpinning, setIsSpinning] = useState(false);
    const [spinResult, setSpinResult] = useState<string | null>(null);
    const [rotation, setRotation] = useState(0);
    const [streak, setStreak] = useState(5);
    const [checkedIn, setCheckedIn] = useState(false);
    const [showConfetti, setShowConfetti] = useState(false);
    const wheelRef = useRef<HTMLDivElement>(null);

    const prizes = [
        { name: '1000 Xu', value: 1000, color: 'from-yellow-400 to-amber-500', icon: '🪙' },
        { name: '5000 Xu', value: 5000, color: 'from-yellow-500 to-orange-500', icon: '💰' },
        { name: 'Voucher 50K', value: 0, color: 'from-purple-400 to-pink-500', icon: '🎟️' },
        { name: '100 Xu', value: 100, color: 'from-gray-400 to-gray-500', icon: '🥉' },
        { name: '10000 Xu', value: 10000, color: 'from-yellow-300 to-yellow-500', icon: '👑' },
        { name: '500 Xu', value: 500, color: 'from-blue-400 to-cyan-500', icon: '🥈' },
        { name: 'Freeship', value: 0, color: 'from-green-400 to-emerald-500', icon: '🚚' },
        { name: '200 Xu', value: 200, color: 'from-rose-400 to-red-500', icon: '🎁' },
    ];

    const missions = [
        { id: 1, name: 'Mua hàng đầu tiên', reward: 500, progress: 0, total: 1, icon: '🛒' },
        { id: 2, name: 'Đánh giá 3 sản phẩm', reward: 300, progress: 2, total: 3, icon: '⭐' },
        { id: 3, name: 'Chia sẻ lên Facebook', reward: 200, progress: 1, total: 1, icon: '📱', completed: true },
        { id: 4, name: 'Mời 2 bạn bè', reward: 1000, progress: 1, total: 2, icon: '👥' },
        { id: 5, name: 'Xem 5 livestream', reward: 250, progress: 3, total: 5, icon: '🔴' },
    ];

    const spinWheel = () => {
        if (isSpinning || coins < 50) return;

        setIsSpinning(true);
        setCoins(prev => prev - 50);
        setSpinResult(null);

        const prizeIndex = Math.floor(Math.random() * prizes.length);
        const extraRotation = 360 * 5 + (prizeIndex * (360 / prizes.length)) + (360 / prizes.length / 2);

        setRotation(prev => prev + extraRotation);

        setTimeout(() => {
            setIsSpinning(false);
            setSpinResult(prizes[prizeIndex].name);
            setCoins(prev => prev + prizes[prizeIndex].value);
            setShowConfetti(true);
            setTimeout(() => setShowConfetti(false), 3000);
        }, 4000);
    };

    const handleCheckIn = () => {
        if (checkedIn) return;
        setCheckedIn(true);
        setStreak(prev => prev + 1);
        setCoins(prev => prev + streak * 10 + 50);
        setShowConfetti(true);
        setTimeout(() => setShowConfetti(false), 2000);
    };

    return (
        <div className="min-h-screen bg-gradient-to-b from-purple-900 via-purple-800 to-indigo-900 animate-fade-in">
            {/* Confetti Effect */}
            {showConfetti && (
                <div className="fixed inset-0 pointer-events-none z-50 overflow-hidden">
                    {[...Array(50)].map((_, i) => (
                        <div
                            key={i}
                            className="absolute w-3 h-3 animate-bounce"
                            style={{
                                left: `${Math.random() * 100}%`,
                                top: '-20px',
                                backgroundColor: ['#EE4D2D', '#FFD700', '#FF69B4', '#00CED1', '#7CFC00'][Math.floor(Math.random() * 5)],
                                borderRadius: Math.random() > 0.5 ? '50%' : '0',
                                animation: `fall ${2 + Math.random() * 2}s linear forwards`,
                                animationDelay: `${Math.random() * 0.5}s`,
                            }}
                        />
                    ))}
                </div>
            )}

            <style jsx>{`
        @keyframes fall {
          to {
            transform: translateY(100vh) rotate(720deg);
            opacity: 0;
          }
        }
      `}</style>

            {/* Header */}
            <div className="relative overflow-hidden py-8">
                <div className="absolute inset-0">
                    {[...Array(30)].map((_, i) => (
                        <div
                            key={i}
                            className="absolute text-2xl animate-float opacity-30"
                            style={{
                                left: `${Math.random() * 100}%`,
                                top: `${Math.random() * 100}%`,
                                animationDelay: `${Math.random() * 3}s`,
                            }}
                        >
                            {['🪙', '💰', '🎁', '⭐', '🎮'][Math.floor(Math.random() * 5)]}
                        </div>
                    ))}
                </div>

                <div className="container mx-auto px-4 relative z-10">
                    <div className="text-center text-white">
                        <h1 className="text-4xl md:text-5xl font-black mb-2">
                            <span className="inline-block animate-bounce">🎮</span> Shopee Rewards
                        </h1>
                        <p className="text-purple-200">Chơi game, nhận xu, đổi quà!</p>

                        {/* Coin display */}
                        <div className="inline-flex items-center gap-3 bg-white/10 backdrop-blur-lg px-6 py-3 rounded-full mt-4 border border-white/20">
                            <span className="text-3xl animate-bounce">🪙</span>
                            <span className="text-3xl font-black text-yellow-300">{coins.toLocaleString()}</span>
                            <span className="text-purple-200">Shopee Xu</span>
                        </div>
                    </div>
                </div>
            </div>

            <div className="container mx-auto px-4 pb-12">
                <div className="grid lg:grid-cols-2 gap-8">
                    {/* Lucky Wheel */}
                    <div className="bg-white/10 backdrop-blur-lg rounded-2xl p-6 border border-white/20 animate-slide-up">
                        <h2 className="text-2xl font-bold text-white text-center mb-6 flex items-center justify-center gap-2">
                            <span className="animate-wiggle">🎰</span> Vòng Quay May Mắn
                        </h2>

                        <div className="relative mx-auto" style={{ width: 300, height: 300 }}>
                            {/* Wheel */}
                            <div
                                ref={wheelRef}
                                className="w-full h-full rounded-full relative overflow-hidden shadow-2xl"
                                style={{
                                    transform: `rotate(${rotation}deg)`,
                                    transition: isSpinning ? 'transform 4s cubic-bezier(0.17, 0.67, 0.12, 0.99)' : 'none',
                                }}
                            >
                                {prizes.map((prize, index) => (
                                    <div
                                        key={prize.name}
                                        className={`absolute w-1/2 h-1/2 origin-bottom-right bg-gradient-to-br ${prize.color}`}
                                        style={{
                                            transform: `rotate(${index * 45}deg) skewY(-45deg)`,
                                            top: 0,
                                            right: '50%',
                                        }}
                                    >
                                        <div
                                            className="absolute text-white font-bold text-xs"
                                            style={{
                                                transform: 'skewY(45deg) rotate(22.5deg)',
                                                top: '40%',
                                                left: '10%',
                                            }}
                                        >
                                            <span className="text-lg">{prize.icon}</span>
                                        </div>
                                    </div>
                                ))}
                            </div>

                            {/* Center button */}
                            <div className="absolute inset-0 flex items-center justify-center">
                                <button
                                    onClick={spinWheel}
                                    disabled={isSpinning || coins < 50}
                                    className={`w-20 h-20 rounded-full text-white font-bold shadow-xl transition-all z-10 ${isSpinning
                                            ? 'bg-gray-500 cursor-not-allowed'
                                            : coins < 50
                                                ? 'bg-gray-500 cursor-not-allowed'
                                                : 'bg-gradient-to-br from-[#EE4D2D] to-[#FF6633] hover:scale-110 animate-pulse-glow'
                                        }`}
                                >
                                    {isSpinning ? '...' : 'QUAY'}
                                </button>
                            </div>

                            {/* Pointer */}
                            <div className="absolute -top-2 left-1/2 -translate-x-1/2 z-20">
                                <div className="w-0 h-0 border-l-[15px] border-r-[15px] border-t-[30px] border-l-transparent border-r-transparent border-t-white drop-shadow-lg" />
                            </div>
                        </div>

                        <p className="text-center text-purple-200 mt-4">
                            Chi phí: <span className="text-yellow-300 font-bold">50 Xu</span> / lượt
                        </p>

                        {/* Result */}
                        {spinResult && (
                            <div className="mt-4 p-4 bg-gradient-to-r from-yellow-400 to-orange-500 rounded-xl text-center animate-bounce-in">
                                <p className="text-white font-bold text-lg">🎉 Chúc mừng!</p>
                                <p className="text-white text-2xl font-black">{spinResult}</p>
                            </div>
                        )}
                    </div>

                    {/* Daily Check-in */}
                    <div className="bg-white/10 backdrop-blur-lg rounded-2xl p-6 border border-white/20 animate-slide-up" style={{ animationDelay: '100ms' }}>
                        <h2 className="text-2xl font-bold text-white text-center mb-6 flex items-center justify-center gap-2">
                            <span className="animate-heartbeat">📅</span> Điểm Danh Hàng Ngày
                        </h2>

                        {/* Streak display */}
                        <div className="text-center mb-6">
                            <div className="inline-flex items-center gap-3 bg-gradient-to-r from-orange-500 to-red-500 px-6 py-3 rounded-full">
                                <span className="text-3xl">🔥</span>
                                <div className="text-white">
                                    <p className="text-sm opacity-80">Chuỗi điểm danh</p>
                                    <p className="text-2xl font-black">{streak} ngày</p>
                                </div>
                            </div>
                        </div>

                        {/* Calendar */}
                        <div className="grid grid-cols-7 gap-2 mb-6">
                            {[...Array(7)].map((_, i) => {
                                const isCompleted = i < streak % 7;
                                const isToday = i === streak % 7;
                                return (
                                    <div
                                        key={i}
                                        className={`aspect-square rounded-lg flex flex-col items-center justify-center text-white transition-all ${isCompleted
                                                ? 'bg-gradient-to-br from-green-400 to-emerald-500 scale-100'
                                                : isToday
                                                    ? 'bg-gradient-to-br from-yellow-400 to-orange-500 animate-pulse scale-110 ring-2 ring-yellow-300'
                                                    : 'bg-white/10'
                                            }`}
                                    >
                                        <span className="text-lg">{isCompleted ? '✓' : isToday ? '🎁' : '?'}</span>
                                        <span className="text-[10px] opacity-80">Ngày {i + 1}</span>
                                    </div>
                                );
                            })}
                        </div>

                        {/* Check-in button */}
                        <button
                            onClick={handleCheckIn}
                            disabled={checkedIn}
                            className={`w-full py-4 rounded-xl font-bold text-lg transition-all ${checkedIn
                                    ? 'bg-gray-500 text-white cursor-not-allowed'
                                    : 'bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] text-white hover:scale-105 hover-shine ripple'
                                }`}
                        >
                            {checkedIn ? '✓ Đã điểm danh hôm nay' : `Điểm danh (+${streak * 10 + 50} Xu)`}
                        </button>

                        {/* Bonus info */}
                        <div className="mt-4 p-4 bg-white/5 rounded-xl">
                            <p className="text-purple-200 text-sm text-center">
                                💡 Điểm danh liên tục {7 - (streak % 7)} ngày nữa để nhận
                                <span className="text-yellow-300 font-bold"> x2 Xu!</span>
                            </p>
                        </div>
                    </div>
                </div>

                {/* Missions */}
                <div className="mt-8 bg-white/10 backdrop-blur-lg rounded-2xl p-6 border border-white/20 animate-slide-up" style={{ animationDelay: '200ms' }}>
                    <h2 className="text-2xl font-bold text-white mb-6 flex items-center gap-2">
                        <span className="animate-wiggle">🎯</span> Nhiệm Vụ Hàng Ngày
                    </h2>

                    <div className="space-y-4">
                        {missions.map((mission, index) => (
                            <div
                                key={mission.id}
                                className={`p-4 rounded-xl transition-all animate-slide-up ${mission.completed
                                        ? 'bg-green-500/20 border border-green-500/30'
                                        : 'bg-white/5 hover:bg-white/10 border border-white/10'
                                    }`}
                                style={{ animationDelay: `${300 + index * 50}ms` }}
                            >
                                <div className="flex items-center gap-4">
                                    <span className="text-3xl">{mission.icon}</span>
                                    <div className="flex-1">
                                        <div className="flex items-center justify-between mb-1">
                                            <h3 className="text-white font-semibold">{mission.name}</h3>
                                            <span className="text-yellow-300 font-bold flex items-center gap-1">
                                                +{mission.reward} 🪙
                                            </span>
                                        </div>

                                        {/* Progress bar */}
                                        <div className="h-2 bg-white/10 rounded-full overflow-hidden">
                                            <div
                                                className={`h-full rounded-full transition-all duration-500 ${mission.completed
                                                        ? 'bg-gradient-to-r from-green-400 to-emerald-500'
                                                        : 'bg-gradient-to-r from-[#EE4D2D] to-[#FF6633]'
                                                    }`}
                                                style={{ width: `${(mission.progress / mission.total) * 100}%` }}
                                            />
                                        </div>
                                        <p className="text-purple-300 text-xs mt-1">
                                            {mission.completed ? '✓ Hoàn thành' : `${mission.progress}/${mission.total}`}
                                        </p>
                                    </div>

                                    {mission.completed && (
                                        <button className="px-4 py-2 bg-gradient-to-r from-green-400 to-emerald-500 text-white rounded-full text-sm font-bold hover:scale-105 transition-transform">
                                            Nhận Xu
                                        </button>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}
