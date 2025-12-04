'use client';

import { useState, useEffect, useRef } from 'react';
import { gamificationApi } from '../../lib/api';
import { useGamificationStore } from '../../lib/store';

const PRIZES = [
    { id: '1', name: '100 Coins', type: 'coins', value: 100, color: '#FFD700', icon: '🪙' },
    { id: '2', name: '500 Coins', type: 'coins', value: 500, color: '#FFA500', icon: '💰' },
    { id: '3', name: '1000 Coins', type: 'coins', value: 1000, color: '#FF4500', icon: '🏆' },
    { id: '4', name: 'Try Again', type: 'nothing', value: 0, color: '#808080', icon: '😢' },
    { id: '5', name: '50 Coins', type: 'coins', value: 50, color: '#98FB98', icon: '🌟' },
    { id: '6', name: '200 Coins', type: 'coins', value: 200, color: '#87CEEB', icon: '💎' },
    { id: '7', name: 'Voucher $5', type: 'voucher', value: 5, color: '#DDA0DD', icon: '🎟️' },
    { id: '8', name: 'Mystery Box', type: 'mystery', value: 0, color: '#FF69B4', icon: '🎁' },
];

export default function GamificationPage() {
    const { balance, streak, setBalance, addCoins, setStreak, setLastCheckIn } = useGamificationStore();
    const [spinning, setSpinning] = useState(false);
    const [rotation, setRotation] = useState(0);
    const [result, setResult] = useState<typeof PRIZES[0] | null>(null);
    const [checkingIn, setCheckingIn] = useState(false);
    const [missions, setMissions] = useState<any[]>([]);
    const [userMissions, setUserMissions] = useState<any[]>([]);
    const wheelRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        try {
            const walletData = await gamificationApi.getBalance('user-123');
            setBalance(walletData.balance);

            const missionData = await gamificationApi.getMissions('user-123');
            setMissions(missionData.missions || []);
            setUserMissions(missionData.user_progress || []);
        } catch {
            // Use mock data
            setMissions([
                { id: 'm1', name: 'First Purchase', description: 'Complete your first order', target: 1, reward: 100 },
                { id: 'm2', name: 'Review Products', description: 'Write 3 product reviews', target: 3, reward: 150 },
                { id: 'm3', name: 'Invite Friends', description: 'Invite 2 friends to join', target: 2, reward: 500 },
            ]);
            setUserMissions([
                { mission_id: 'm1', progress: 1, completed: true, claimed_at: null },
                { mission_id: 'm2', progress: 2, completed: false },
                { mission_id: 'm3', progress: 0, completed: false },
            ]);
        }
    };

    const handleSpin = async () => {
        if (spinning || balance < 50) return;

        setSpinning(true);
        setResult(null);

        // Calculate random prize and rotation
        const prizeIndex = Math.floor(Math.random() * PRIZES.length);
        const segmentAngle = 360 / PRIZES.length;
        const targetRotation = 360 * 5 + (360 - prizeIndex * segmentAngle - segmentAngle / 2);

        setRotation(prev => prev + targetRotation);

        try {
            const spinResult = await gamificationApi.spinLuckyDraw('user-123', 50);

            setTimeout(() => {
                const prize = PRIZES[prizeIndex];
                setResult(prize);
                setSpinning(false);
                if (prize.value > 0) {
                    addCoins(prize.value - 50); // Net gain
                } else {
                    addCoins(-50); // Lost spin cost
                }
            }, 4000);
        } catch {
            setTimeout(() => {
                const prize = PRIZES[prizeIndex];
                setResult(prize);
                setSpinning(false);
                if (prize.value > 0 && prize.type === 'coins') {
                    addCoins(prize.value - 50);
                } else {
                    addCoins(-50);
                }
            }, 4000);
        }
    };

    const handleCheckIn = async () => {
        setCheckingIn(true);
        try {
            const result = await gamificationApi.dailyCheckIn('user-123');
            addCoins(result.reward);
            setStreak(result.streak);
            setLastCheckIn(new Date().toISOString());
            alert(`🎉 +${result.reward} coins! Streak: ${result.streak} days`);
        } catch {
            // Mock success
            const reward = Math.min(streak + 1, 10);
            addCoins(reward);
            setStreak(streak + 1);
            setLastCheckIn(new Date().toISOString());
            alert(`🎉 +${reward} coins! Streak: ${streak + 1} days`);
        } finally {
            setCheckingIn(false);
        }
    };

    return (
        <div className="container mx-auto py-8">
            {/* Header */}
            <div className="mb-8 rounded-xl bg-gradient-to-r from-yellow-500 via-orange-500 to-red-500 p-8 text-white">
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-4xl font-bold">🎮 Gamification Center</h1>
                        <p className="text-lg opacity-90">Earn coins, spin the wheel, complete missions!</p>
                    </div>
                    <div className="text-right">
                        <div className="text-sm opacity-75">Your Balance</div>
                        <div className="text-4xl font-bold">🪙 {balance.toLocaleString()}</div>
                    </div>
                </div>
            </div>

            <div className="grid gap-8 lg:grid-cols-2">
                {/* Lucky Draw Wheel */}
                <div className="rounded-xl border bg-white p-8">
                    <h2 className="mb-6 text-2xl font-bold text-center">🎡 Lucky Draw</h2>

                    <div className="relative mx-auto w-[320px] h-[320px]">
                        {/* Pointer */}
                        <div className="absolute top-0 left-1/2 -translate-x-1/2 -translate-y-2 z-20">
                            <div className="w-0 h-0 border-l-[15px] border-r-[15px] border-t-[25px] border-l-transparent border-r-transparent border-t-red-600" />
                        </div>

                        {/* Wheel */}
                        <div
                            ref={wheelRef}
                            className="w-full h-full rounded-full border-8 border-yellow-500 shadow-2xl overflow-hidden transition-transform"
                            style={{
                                transform: `rotate(${rotation}deg)`,
                                transitionDuration: spinning ? '4s' : '0s',
                                transitionTimingFunction: 'cubic-bezier(0.17, 0.67, 0.12, 0.99)',
                            }}
                        >
                            <div className="relative w-full h-full">
                                {PRIZES.map((prize, index) => {
                                    const angle = (360 / PRIZES.length) * index;
                                    const skewY = 90 - 360 / PRIZES.length;

                                    return (
                                        <div
                                            key={prize.id}
                                            className="absolute w-1/2 h-1/2 origin-bottom-right overflow-hidden"
                                            style={{
                                                transform: `rotate(${angle}deg) skewY(-${skewY}deg)`,
                                                backgroundColor: prize.color,
                                            }}
                                        >
                                            <div
                                                className="absolute bottom-4 left-4 text-2xl font-bold text-white drop-shadow-lg"
                                                style={{
                                                    transform: `skewY(${skewY}deg) rotate(${180 / PRIZES.length}deg)`,
                                                }}
                                            >
                                                {prize.icon}
                                            </div>
                                        </div>
                                    );
                                })}
                            </div>
                        </div>

                        {/* Center Button */}
                        <button
                            onClick={handleSpin}
                            disabled={spinning || balance < 50}
                            className={`absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-20 h-20 rounded-full font-bold text-white shadow-lg transition-all z-10 ${spinning || balance < 50
                                    ? 'bg-gray-400 cursor-not-allowed'
                                    : 'bg-gradient-to-br from-red-500 to-orange-600 hover:scale-110 active:scale-95'
                                }`}
                        >
                            {spinning ? '🎰' : 'SPIN'}
                        </button>
                    </div>

                    <div className="mt-6 text-center">
                        <p className="text-muted-foreground">Cost: 50 coins per spin</p>
                        {result && (
                            <div className={`mt-4 p-4 rounded-lg ${result.value > 0 ? 'bg-green-100 text-green-800' : 'bg-gray-100'}`}>
                                <span className="text-3xl">{result.icon}</span>
                                <p className="text-lg font-bold">{result.name}</p>
                                {result.value > 0 && result.type === 'coins' && (
                                    <p className="text-sm">Net: +{result.value - 50} coins</p>
                                )}
                            </div>
                        )}
                    </div>
                </div>

                {/* Daily Check-in & Missions */}
                <div className="space-y-8">
                    {/* Daily Check-in */}
                    <div className="rounded-xl border bg-white p-6">
                        <h2 className="mb-4 text-xl font-bold">📅 Daily Check-in</h2>
                        <div className="flex items-center justify-between">
                            <div className="flex gap-2">
                                {[1, 2, 3, 4, 5, 6, 7].map((day) => (
                                    <div
                                        key={day}
                                        className={`w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold ${day <= streak
                                                ? 'bg-green-500 text-white'
                                                : day === streak + 1
                                                    ? 'bg-yellow-100 border-2 border-yellow-500'
                                                    : 'bg-gray-100'
                                            }`}
                                    >
                                        {day <= streak ? '✓' : day}
                                    </div>
                                ))}
                            </div>
                            <button
                                onClick={handleCheckIn}
                                disabled={checkingIn}
                                className="rounded-lg bg-gradient-to-r from-green-500 to-emerald-600 px-6 py-3 font-bold text-white hover:from-green-600 hover:to-emerald-700 disabled:opacity-50"
                            >
                                {checkingIn ? '...' : 'Check In'}
                            </button>
                        </div>
                        <p className="mt-3 text-sm text-muted-foreground">
                            🔥 Current streak: {streak} days | Reward: {Math.min(streak + 1, 10)} coins
                        </p>
                    </div>

                    {/* Missions */}
                    <div className="rounded-xl border bg-white p-6">
                        <h2 className="mb-4 text-xl font-bold">🎯 Missions</h2>
                        <div className="space-y-4">
                            {missions.map((mission) => {
                                const userProgress = userMissions.find((um) => um.mission_id === mission.id);
                                const progress = userProgress?.progress || 0;
                                const completed = userProgress?.completed || false;
                                const claimed = userProgress?.claimed_at != null;
                                const percent = Math.min((progress / mission.target) * 100, 100);

                                return (
                                    <div
                                        key={mission.id}
                                        className={`rounded-lg border p-4 ${completed ? 'bg-green-50' : 'bg-white'}`}
                                    >
                                        <div className="flex items-start justify-between">
                                            <div>
                                                <h3 className="font-semibold">{mission.name}</h3>
                                                <p className="text-sm text-muted-foreground">{mission.description}</p>
                                            </div>
                                            <div className="text-right">
                                                <span className="text-lg font-bold text-yellow-600">
                                                    +{mission.reward} 🪙
                                                </span>
                                            </div>
                                        </div>
                                        <div className="mt-3">
                                            <div className="flex items-center justify-between text-sm">
                                                <span>{progress}/{mission.target}</span>
                                                <span>{Math.round(percent)}%</span>
                                            </div>
                                            <div className="mt-1 h-2 rounded-full bg-gray-200">
                                                <div
                                                    className={`h-full rounded-full transition-all ${completed ? 'bg-green-500' : 'bg-blue-500'
                                                        }`}
                                                    style={{ width: `${percent}%` }}
                                                />
                                            </div>
                                        </div>
                                        {completed && !claimed && (
                                            <button className="mt-3 w-full rounded-lg bg-green-500 py-2 font-bold text-white hover:bg-green-600">
                                                Claim Reward
                                            </button>
                                        )}
                                        {claimed && (
                                            <div className="mt-3 text-center text-green-600 font-semibold">
                                                ✓ Claimed
                                            </div>
                                        )}
                                    </div>
                                );
                            })}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
