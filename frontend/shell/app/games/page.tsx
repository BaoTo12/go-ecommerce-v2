'use client';

import React, { useState, useEffect } from 'react';
import { coinsService } from '@/services/coinsService';

interface GameResult {
    won: boolean;
    coins: number;
}

export default function MiniGamesPage() {
    const [activeGame, setActiveGame] = useState<string | null>(null);
    const [coins, setCoins] = useState(0);
    const [gameResult, setGameResult] = useState<GameResult | null>(null);

    // Memory Game State
    const [memoryCards, setMemoryCards] = useState<number[]>([]);
    const [flippedCards, setFlippedCards] = useState<number[]>([]);
    const [matchedPairs, setMatchedPairs] = useState<number[]>([]);

    // Coin Flip State
    const [coinFlipResult, setCoinFlipResult] = useState<'heads' | 'tails' | null>(null);
    const [isFlipping, setIsFlipping] = useState(false);

    // Lucky Number State
    const [luckyNumber, setLuckyNumber] = useState<number | null>(null);
    const [selectedNumber, setSelectedNumber] = useState<number | null>(null);
    const [isRevealing, setIsRevealing] = useState(false);

    useEffect(() => {
        coinsService.getState().then(state => setCoins(state.balance));
    }, []);

    // Memory Game Logic
    const initMemoryGame = () => {
        const pairs = [1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6];
        setMemoryCards(pairs.sort(() => Math.random() - 0.5));
        setFlippedCards([]);
        setMatchedPairs([]);
        setActiveGame('memory');
        setGameResult(null);
    };

    const flipCard = (index: number) => {
        if (flippedCards.length >= 2 || flippedCards.includes(index) || matchedPairs.includes(memoryCards[index])) {
            return;
        }

        const newFlipped = [...flippedCards, index];
        setFlippedCards(newFlipped);

        if (newFlipped.length === 2) {
            const [first, second] = newFlipped;
            if (memoryCards[first] === memoryCards[second]) {
                const newMatched = [...matchedPairs, memoryCards[first]];
                setMatchedPairs(newMatched);
                setFlippedCards([]);

                if (newMatched.length === 6) {
                    const reward = 30;
                    coinsService.earnCoins(reward, 'Thắng game Lật Thẻ');
                    setGameResult({ won: true, coins: reward });
                    setCoins(c => c + reward);
                }
            } else {
                setTimeout(() => setFlippedCards([]), 1000);
            }
        }
    };

    // Coin Flip Game Logic
    const playCoinFlip = (bet: 'heads' | 'tails') => {
        if (isFlipping) return;
        setIsFlipping(true);
        setGameResult(null);

        setTimeout(() => {
            const result = Math.random() < 0.5 ? 'heads' : 'tails';
            setCoinFlipResult(result);
            setIsFlipping(false);

            if (result === bet) {
                const reward = 20;
                coinsService.earnCoins(reward, 'Thắng game Tung Xu');
                setGameResult({ won: true, coins: reward });
                setCoins(c => c + reward);
            } else {
                setGameResult({ won: false, coins: 0 });
            }
        }, 1500);
    };

    // Lucky Number Game Logic
    const playLuckyNumber = (num: number) => {
        if (isRevealing) return;
        setSelectedNumber(num);
        setIsRevealing(true);
        setGameResult(null);

        setTimeout(() => {
            const lucky = Math.floor(Math.random() * 9) + 1;
            setLuckyNumber(lucky);
            setIsRevealing(false);

            if (lucky === num) {
                const reward = 100;
                coinsService.earnCoins(reward, 'Thắng game Số May Mắn');
                setGameResult({ won: true, coins: reward });
                setCoins(c => c + reward);
            } else {
                setGameResult({ won: false, coins: 0 });
            }
        }, 2000);
    };

    const cardEmojis = ['🎁', '⭐', '🎈', '🎉', '💎', '🔥'];

    return (
        <div className="min-h-screen bg-gradient-to-b from-purple-600 to-blue-600">
            <div className="container mx-auto px-4 py-6">
                {/* Header */}
                <div className="text-center mb-8">
                    <h1 className="text-3xl font-bold text-white mb-2">🎮 Mini Games</h1>
                    <p className="text-white/80">Chơi game, kiếm xu mỗi ngày!</p>
                    <div className="inline-flex items-center gap-2 bg-white/20 rounded-full px-4 py-2 mt-4">
                        <span className="text-2xl">🪙</span>
                        <span className="text-white font-bold text-lg">{coins}</span>
                    </div>
                </div>

                {/* Game Result Modal */}
                {gameResult && (
                    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50" onClick={() => setGameResult(null)}>
                        <div className="bg-white rounded-2xl p-8 text-center animate-bounce-in" onClick={e => e.stopPropagation()}>
                            <div className="text-6xl mb-4">{gameResult.won ? '🎉' : '😅'}</div>
                            <h3 className="text-xl font-bold">
                                {gameResult.won ? `Thắng! +${gameResult.coins} xu` : 'Chúc may mắn lần sau!'}
                            </h3>
                            <button
                                onClick={() => setGameResult(null)}
                                className="mt-4 px-6 py-2 bg-[#ee4d2d] text-white rounded-full"
                            >
                                OK
                            </button>
                        </div>
                    </div>
                )}

                {/* Game Selection or Active Game */}
                {!activeGame ? (
                    <div className="grid md:grid-cols-3 gap-6">
                        {/* Memory Game */}
                        <button
                            onClick={initMemoryGame}
                            className="bg-white rounded-2xl p-6 text-center shadow-lg hover:shadow-xl transition-all hover:-translate-y-1"
                        >
                            <div className="text-5xl mb-4">🃏</div>
                            <h3 className="font-bold text-lg mb-2">Lật Thẻ</h3>
                            <p className="text-sm text-gray-500 mb-4">Tìm các cặp thẻ giống nhau</p>
                            <div className="text-[#ee4d2d] font-medium">🪙 +30 xu</div>
                        </button>

                        {/* Coin Flip */}
                        <button
                            onClick={() => setActiveGame('coinflip')}
                            className="bg-white rounded-2xl p-6 text-center shadow-lg hover:shadow-xl transition-all hover:-translate-y-1"
                        >
                            <div className="text-5xl mb-4">🪙</div>
                            <h3 className="font-bold text-lg mb-2">Tung Xu</h3>
                            <p className="text-sm text-gray-500 mb-4">Đoán mặt sấp hay ngửa</p>
                            <div className="text-[#ee4d2d] font-medium">🪙 +20 xu</div>
                        </button>

                        {/* Lucky Number */}
                        <button
                            onClick={() => setActiveGame('luckynumber')}
                            className="bg-white rounded-2xl p-6 text-center shadow-lg hover:shadow-xl transition-all hover:-translate-y-1"
                        >
                            <div className="text-5xl mb-4">🔢</div>
                            <h3 className="font-bold text-lg mb-2">Số May Mắn</h3>
                            <p className="text-sm text-gray-500 mb-4">Chọn đúng số trúng thưởng</p>
                            <div className="text-[#ee4d2d] font-medium">🪙 +100 xu</div>
                        </button>
                    </div>
                ) : (
                    <div className="bg-white rounded-2xl p-6 shadow-lg max-w-lg mx-auto">
                        <button
                            onClick={() => {
                                setActiveGame(null);
                                setGameResult(null);
                            }}
                            className="text-gray-500 mb-4 hover:text-gray-700"
                        >
                            ← Quay lại
                        </button>

                        {/* Memory Game */}
                        {activeGame === 'memory' && (
                            <div className="animate-fade-in">
                                <h3 className="text-lg font-bold mb-4 text-center">🃏 Lật Thẻ</h3>
                                <div className="grid grid-cols-4 gap-2">
                                    {memoryCards.map((card, index) => {
                                        const isFlipped = flippedCards.includes(index);
                                        const isMatched = matchedPairs.includes(card);
                                        return (
                                            <button
                                                key={index}
                                                onClick={() => flipCard(index)}
                                                className={`aspect-square rounded-lg text-3xl transition-all ${isFlipped || isMatched
                                                        ? 'bg-[#ee4d2d] text-white rotate-0'
                                                        : 'bg-gray-200 hover:bg-gray-300 rotate-180'
                                                    } ${isMatched ? 'opacity-50' : ''}`}
                                                disabled={isMatched}
                                            >
                                                {(isFlipped || isMatched) && cardEmojis[card - 1]}
                                            </button>
                                        );
                                    })}
                                </div>
                                <p className="text-center text-sm text-gray-500 mt-4">
                                    Đã ghép: {matchedPairs.length}/6
                                </p>
                            </div>
                        )}

                        {/* Coin Flip Game */}
                        {activeGame === 'coinflip' && (
                            <div className="animate-fade-in text-center">
                                <h3 className="text-lg font-bold mb-4">🪙 Tung Xu</h3>
                                <div className={`text-7xl my-8 ${isFlipping ? 'animate-spin' : ''}`}>
                                    {coinFlipResult === 'heads' ? '😊' : coinFlipResult === 'tails' ? '😎' : '🪙'}
                                </div>
                                {coinFlipResult && !isFlipping && (
                                    <p className="mb-4 font-medium">
                                        Kết quả: {coinFlipResult === 'heads' ? 'MẶT NGỬA' : 'MẶT SẤP'}
                                    </p>
                                )}
                                <div className="flex gap-4 justify-center">
                                    <button
                                        onClick={() => playCoinFlip('heads')}
                                        disabled={isFlipping}
                                        className="px-6 py-3 bg-yellow-500 text-white rounded-full font-medium hover:bg-yellow-600 disabled:opacity-50"
                                    >
                                        😊 Mặt Ngửa
                                    </button>
                                    <button
                                        onClick={() => playCoinFlip('tails')}
                                        disabled={isFlipping}
                                        className="px-6 py-3 bg-blue-500 text-white rounded-full font-medium hover:bg-blue-600 disabled:opacity-50"
                                    >
                                        😎 Mặt Sấp
                                    </button>
                                </div>
                            </div>
                        )}

                        {/* Lucky Number Game */}
                        {activeGame === 'luckynumber' && (
                            <div className="animate-fade-in text-center">
                                <h3 className="text-lg font-bold mb-4">🔢 Số May Mắn</h3>
                                {luckyNumber && !isRevealing && (
                                    <div className="text-6xl font-bold text-[#ee4d2d] mb-4">
                                        {luckyNumber}
                                    </div>
                                )}
                                {isRevealing && (
                                    <div className="text-6xl font-bold text-gray-300 mb-4 animate-pulse">
                                        ?
                                    </div>
                                )}
                                <p className="text-sm text-gray-500 mb-4">Chọn một số từ 1-9</p>
                                <div className="grid grid-cols-3 gap-2 max-w-xs mx-auto">
                                    {[1, 2, 3, 4, 5, 6, 7, 8, 9].map(num => (
                                        <button
                                            key={num}
                                            onClick={() => playLuckyNumber(num)}
                                            disabled={isRevealing}
                                            className={`aspect-square rounded-lg text-2xl font-bold transition-all ${selectedNumber === num
                                                    ? 'bg-[#ee4d2d] text-white'
                                                    : 'bg-gray-100 hover:bg-gray-200'
                                                } disabled:opacity-50`}
                                        >
                                            {num}
                                        </button>
                                    ))}
                                </div>
                            </div>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
}
