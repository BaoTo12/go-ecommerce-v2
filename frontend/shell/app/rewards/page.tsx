'use client';

import React, { useState, useEffect, useRef } from 'react';

interface Prize {
    id: string;
    name: string;
    value: number;
    type: string;
    color: string;
}

interface Mission {
    id: string;
    name: string;
    description: string;
    target: number;
    reward: number;
    icon: string;
}

const PRIZES: Prize[] = [
    { id: '1', name: '100 Xu', value: 100, type: 'coins', color: '#FFD700' },
    { id: '2', name: '50 Xu', value: 50, type: 'coins', color: '#FFA500' },
    { id: '3', name: '200 Xu', value: 200, type: 'coins', color: '#FF6B6B' },
    { id: '4', name: 'Thử lại', value: 0, type: 'nothing', color: '#CCCCCC' },
    { id: '5', name: '500 Xu', value: 500, type: 'coins', color: '#4ECDC4' },
    { id: '6', name: '20 Xu', value: 20, type: 'coins', color: '#95E1D3' },
    { id: '7', name: 'Voucher 10K', value: 10000, type: 'voucher', color: '#F38181' },
    { id: '8', name: '1000 Xu', value: 1000, type: 'coins', color: '#AA96DA' },
];

export default function GamificationPage() {
    const [balance, setBalance] = useState(1250);
    const [spinning, setSpinning] = useState(false);
    const [rotation, setRotation] = useState(0);
    const [result, setResult] = useState<Prize | null>(null);
    const [streak, setStreak] = useState(5);
    const [checkedIn, setCheckedIn] = useState(false);
    const wheelRef = useRef<HTMLDivElement>(null);

    const missions: Mission[] = [
        { id: 'm1', name: 'Mua hàng đầu tiên', description: 'Hoàn thành đơn hàng đầu tiên', target: 1, reward: 100, icon: '🛒' },
        { id: 'm2', name: 'Đánh giá sản phẩm', description: 'Viết 3 đánh giá sản phẩm', target: 3, reward: 150, icon: '⭐' },
        { id: 'm3', name: 'Mời bạn bè', description: 'Mời 2 bạn tham gia Shopee', target: 2, reward: 500, icon: '👥' },
        { id: 'm4', name: 'Xem Shopee Live', description: 'Xem 5 phút livestream', target: 5, reward: 50, icon: '📺' },
    ];

    const [missionProgress] = useState<Record<string, number>>({
        m1: 1,
        m2: 2,
        m3: 0,
        m4: 5,
    });

    const handleSpin = () => {
        if (spinning || balance < 50) return;

        setSpinning(true);
        setResult(null);
        setBalance(prev => prev - 50);

        const prizeIndex = Math.floor(Math.random() * PRIZES.length);
        const segmentAngle = 360 / PRIZES.length;
        const targetRotation = 360 * 5 + (360 - prizeIndex * segmentAngle - segmentAngle / 2);

        setRotation(prev => prev + targetRotation);

        setTimeout(() => {
            const prize = PRIZES[prizeIndex];
            setResult(prize);
            setSpinning(false);
            if (prize.value > 0 && prize.type === 'coins') {
                setBalance(prev => prev + prize.value);
            }
        }, 4000);
    };

    const handleCheckIn = () => {
        if (checkedIn) return;
        const reward = Math.min(streak + 1, 10);
        setBalance(prev => prev + reward);
        setStreak(prev => prev + 1);
        setCheckedIn(true);
    };

    return (
        <div className="min-h-screen bg-[#F5F5F5]">
            {/* Header */}
            <div className="bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] py-6">
                <div className="container mx-auto px-4">
                    <div className="flex items-center justify-between text-white">
                        <div>
                            <h1 className="text-2xl font-bold">🎮 Shopee Rewards</h1>
                            <p className="text-sm opacity-90">Chơi game, nhận xu, đổi quà!</p>
                        </div>
                        <div className="text-right">
                            <div className="text-sm opacity-75">Số dư Xu</div>
                            <div className="text-3xl font-bold flex items-center gap-1">
                                <span className="text-yellow-300">🪙</span> {balance.toLocaleString()}
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div className="container mx-auto px-4 py-6">
                <div className="grid gap-6 lg:grid-cols-2">
                    {/* Lucky Wheel Section */}
                    <div className="bg-white rounded-sm p-6">
                        <h2 className="text-lg font-bold text-center mb-6 text-[#EE4D2D]">
                            🎡 Vòng Quay May Mắn
                        </h2>

                        {/* Wheel Container */}
                        <div className="relative w-72 h-72 mx-auto mb-6">
                            {/* Pointer */}
                            <div className="absolute top-0 left-1/2 -translate-x-1/2 -translate-y-2 z-20">
                                <div className="w-0 h-0 border-l-[12px] border-r-[12px] border-t-[20px] border-l-transparent border-r-transparent border-t-[#EE4D2D]" />
                            </div>

                            {/* Wheel */}
                            <div
                                ref={wheelRef}
                                className="w-full h-full rounded-full border-8 border-[#EE4D2D] shadow-xl overflow-hidden relative"
                                style={{
                                    transform: `rotate(${rotation}deg)`,
                                    transition: spinning ? 'transform 4s cubic-bezier(0.17, 0.67, 0.12, 0.99)' : 'none',
                                }}
                            >
                                {PRIZES.map((prize, index) => {
                                    const angle = (360 / PRIZES.length) * index;
                                    return (
                                        <div
                                            key={prize.id}
                                            className="absolute w-1/2 h-1/2 origin-bottom-right"
                                            style={{
                                                transform: `rotate(${angle}deg) skewY(-${90 - 360 / PRIZES.length}deg)`,
                                                background: prize.color,
                                            }}
                                        />
                                    );
                                })}
                                {/* Center circle */}
                                <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-16 h-16 bg-white rounded-full shadow-lg flex items-center justify-center z-10">
                                    <span className="text-2xl">🎰</span>
                                </div>
                            </div>
                        </div>

                        {/* Spin Button */}
                        <button
                            onClick={handleSpin}
                            disabled={spinning || balance < 50}
                            className={`w-full py-3 rounded font-bold text-white transition-all ${spinning || balance < 50
                                    ? 'bg-gray-400 cursor-not-allowed'
                                    : 'bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] hover:opacity-90 active:scale-[0.98]'
                                }`}
                        >
                            {spinning ? '⏳ Đang quay...' : '🎲 Quay (50 Xu)'}
                        </button>

                        {/* Result */}
                        {result && (
                            <div className={`mt-4 p-4 rounded text-center ${result.value > 0 ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'
                                }`}>
                                <p className="text-lg font-bold">
                                    {result.value > 0 ? '🎉 Chúc mừng!' : '😢 Tiếc quá!'}
                                </p>
                                <p>{result.name}</p>
                            </div>
                        )}
                    </div>

                    {/* Check-in & Missions */}
                    <div className="space-y-6">
                        {/* Daily Check-in */}
                        <div className="bg-white rounded-sm p-6">
                            <h2 className="text-lg font-bold mb-4 text-[#EE4D2D]">📅 Điểm Danh Hàng Ngày</h2>

                            {/* Streak Calendar */}
                            <div className="flex justify-between mb-4">
                                {[1, 2, 3, 4, 5, 6, 7].map(day => (
                                    <div
                                        key={day}
                                        className={`w-10 h-10 rounded-full flex flex-col items-center justify-center text-xs ${day < streak
                                                ? 'bg-[#EE4D2D] text-white'
                                                : day === streak
                                                    ? 'bg-yellow-400 text-white border-2 border-yellow-500'
                                                    : 'bg-gray-100 text-gray-500'
                                            }`}
                                    >
                                        <span className="font-bold">{day}</span>
                                        {day < streak && <span>✓</span>}
                                    </div>
                                ))}
                            </div>

                            <div className="flex items-center justify-between">
                                <div>
                                    <p className="text-sm text-gray-600">
                                        🔥 Chuỗi: <span className="font-bold text-[#EE4D2D]">{streak} ngày</span>
                                    </p>
                                    <p className="text-xs text-gray-400">
                                        Phần thưởng: {Math.min(streak + 1, 10)} xu
                                    </p>
                                </div>
                                <button
                                    onClick={handleCheckIn}
                                    disabled={checkedIn}
                                    className={`px-6 py-2 rounded font-semibold ${checkedIn
                                            ? 'bg-gray-200 text-gray-500 cursor-not-allowed'
                                            : 'bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] text-white hover:opacity-90'
                                        }`}
                                >
                                    {checkedIn ? '✓ Đã điểm danh' : 'Điểm danh'}
                                </button>
                            </div>
                        </div>

                        {/* Missions */}
                        <div className="bg-white rounded-sm p-6">
                            <h2 className="text-lg font-bold mb-4 text-[#EE4D2D]">🎯 Nhiệm Vụ</h2>

                            <div className="space-y-3">
                                {missions.map(mission => {
                                    const progress = missionProgress[mission.id] || 0;
                                    const completed = progress >= mission.target;
                                    const percent = Math.min((progress / mission.target) * 100, 100);

                                    return (
                                        <div
                                            key={mission.id}
                                            className={`p-3 rounded border ${completed ? 'bg-green-50 border-green-200' : 'border-gray-200'
                                                }`}
                                        >
                                            <div className="flex items-start justify-between">
                                                <div className="flex items-start gap-3">
                                                    <span className="text-2xl">{mission.icon}</span>
                                                    <div>
                                                        <h3 className="font-semibold text-sm">{mission.name}</h3>
                                                        <p className="text-xs text-gray-500">{mission.description}</p>
                                                    </div>
                                                </div>
                                                <span className="text-[#EE4D2D] font-bold text-sm">+{mission.reward} 🪙</span>
                                            </div>

                                            <div className="mt-2">
                                                <div className="flex justify-between text-xs text-gray-500 mb-1">
                                                    <span>{progress}/{mission.target}</span>
                                                    <span>{Math.round(percent)}%</span>
                                                </div>
                                                <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
                                                    <div
                                                        className={`h-full rounded-full transition-all ${completed ? 'bg-green-500' : 'bg-[#EE4D2D]'
                                                            }`}
                                                        style={{ width: `${percent}%` }}
                                                    />
                                                </div>
                                            </div>

                                            {completed && (
                                                <button className="mt-2 w-full py-1.5 bg-green-500 text-white text-sm rounded font-semibold hover:bg-green-600">
                                                    Nhận thưởng
                                                </button>
                                            )}
                                        </div>
                                    );
                                })}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
