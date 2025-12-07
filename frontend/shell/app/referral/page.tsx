'use client';

import React, { useState } from 'react';
import { coinsService } from '@/services/coinsService';

interface Referral {
    id: string;
    name: string;
    avatar: string;
    joinedAt: string;
    status: 'pending' | 'completed';
    reward: number;
}

const MOCK_REFERRALS: Referral[] = [
    { id: 'r1', name: 'Nguyễn Văn B', avatar: 'https://ui-avatars.com/api/?name=NB&background=random', joinedAt: '2024-12-05', status: 'completed', reward: 500 },
    { id: 'r2', name: 'Trần Thị C', avatar: 'https://ui-avatars.com/api/?name=TC&background=random', joinedAt: '2024-12-04', status: 'completed', reward: 500 },
    { id: 'r3', name: 'Lê Văn D', avatar: 'https://ui-avatars.com/api/?name=LD&background=random', joinedAt: '2024-12-06', status: 'pending', reward: 0 },
];

export default function ReferralPage() {
    const [referrals] = useState(MOCK_REFERRALS);
    const [copied, setCopied] = useState(false);

    const referralCode = 'SHOPEE-NGUYENVANA-2024';
    const referralLink = `https://shopee.vn/register?ref=${referralCode}`;

    const totalEarned = referrals.filter(r => r.status === 'completed').reduce((sum, r) => sum + r.reward, 0);
    const pendingRewards = referrals.filter(r => r.status === 'pending').length * 500;

    const copyLink = () => {
        navigator.clipboard.writeText(referralLink);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    const shareToSocial = (platform: string) => {
        const text = encodeURIComponent('🎁 Đăng ký Shopee qua link của mình để nhận ngay voucher 50K! ');
        const url = encodeURIComponent(referralLink);

        const urls: Record<string, string> = {
            facebook: `https://www.facebook.com/sharer/sharer.php?u=${url}&quote=${text}`,
            twitter: `https://twitter.com/intent/tweet?text=${text}&url=${url}`,
            telegram: `https://t.me/share/url?url=${url}&text=${text}`,
            zalo: `https://zalo.me/share?url=${url}`,
        };

        window.open(urls[platform], '_blank', 'width=600,height=400');
    };

    return (
        <div className="min-h-screen bg-gradient-to-b from-[#ee4d2d] to-[#ff6633]">
            <div className="container mx-auto px-4 py-8">
                {/* Header */}
                <div className="text-center text-white mb-8">
                    <h1 className="text-3xl font-bold mb-2">🎁 Mời Bạn Bè</h1>
                    <p className="opacity-90">Mỗi người bạn giới thiệu thành công = 500 xu</p>
                </div>

                {/* Stats */}
                <div className="grid grid-cols-3 gap-4 mb-8">
                    <div className="bg-white/20 backdrop-blur-sm rounded-xl p-4 text-center text-white">
                        <div className="text-3xl font-bold">{referrals.length}</div>
                        <div className="text-sm opacity-80">Đã mời</div>
                    </div>
                    <div className="bg-white/20 backdrop-blur-sm rounded-xl p-4 text-center text-white">
                        <div className="text-3xl font-bold">🪙 {totalEarned}</div>
                        <div className="text-sm opacity-80">Đã nhận</div>
                    </div>
                    <div className="bg-white/20 backdrop-blur-sm rounded-xl p-4 text-center text-white">
                        <div className="text-3xl font-bold">🪙 {pendingRewards}</div>
                        <div className="text-sm opacity-80">Đang chờ</div>
                    </div>
                </div>

                {/* Referral Link Card */}
                <div className="bg-white rounded-2xl shadow-lg p-6 mb-6 animate-fade-in-up">
                    <h3 className="font-bold mb-4">🔗 Link giới thiệu của bạn</h3>

                    <div className="flex gap-2 mb-4">
                        <input
                            type="text"
                            value={referralLink}
                            readOnly
                            className="flex-1 bg-gray-100 px-4 py-3 rounded-lg text-sm text-gray-600"
                        />
                        <button
                            onClick={copyLink}
                            className={`px-6 py-3 rounded-lg font-medium transition-all ${copied ? 'bg-green-500 text-white' : 'bg-[#ee4d2d] text-white'
                                }`}
                        >
                            {copied ? '✓ Đã sao chép' : '📋 Copy'}
                        </button>
                    </div>

                    <div className="flex gap-3 justify-center">
                        <button
                            onClick={() => shareToSocial('facebook')}
                            className="w-12 h-12 bg-blue-600 text-white rounded-full flex items-center justify-center text-xl hover:opacity-80"
                        >
                            f
                        </button>
                        <button
                            onClick={() => shareToSocial('twitter')}
                            className="w-12 h-12 bg-sky-500 text-white rounded-full flex items-center justify-center text-xl hover:opacity-80"
                        >
                            𝕏
                        </button>
                        <button
                            onClick={() => shareToSocial('telegram')}
                            className="w-12 h-12 bg-blue-400 text-white rounded-full flex items-center justify-center text-xl hover:opacity-80"
                        >
                            ✈️
                        </button>
                        <button
                            onClick={() => shareToSocial('zalo')}
                            className="w-12 h-12 bg-blue-500 text-white rounded-full flex items-center justify-center text-xl hover:opacity-80"
                        >
                            Z
                        </button>
                    </div>
                </div>

                {/* How it works */}
                <div className="bg-white rounded-2xl shadow-lg p-6 mb-6 animate-fade-in-up" style={{ animationDelay: '100ms' }}>
                    <h3 className="font-bold mb-4">💡 Cách hoạt động</h3>
                    <div className="flex justify-between items-center">
                        <div className="text-center flex-1">
                            <div className="w-12 h-12 bg-[#fef6f5] rounded-full flex items-center justify-center mx-auto text-2xl mb-2">
                                1️⃣
                            </div>
                            <p className="text-xs">Chia sẻ link</p>
                        </div>
                        <div className="text-gray-300">→</div>
                        <div className="text-center flex-1">
                            <div className="w-12 h-12 bg-[#fef6f5] rounded-full flex items-center justify-center mx-auto text-2xl mb-2">
                                2️⃣
                            </div>
                            <p className="text-xs">Bạn bè đăng ký</p>
                        </div>
                        <div className="text-gray-300">→</div>
                        <div className="text-center flex-1">
                            <div className="w-12 h-12 bg-[#fef6f5] rounded-full flex items-center justify-center mx-auto text-2xl mb-2">
                                3️⃣
                            </div>
                            <p className="text-xs">Đặt đơn đầu</p>
                        </div>
                        <div className="text-gray-300">→</div>
                        <div className="text-center flex-1">
                            <div className="w-12 h-12 bg-[#fef6f5] rounded-full flex items-center justify-center mx-auto text-2xl mb-2">
                                🎁
                            </div>
                            <p className="text-xs">Nhận 500 xu!</p>
                        </div>
                    </div>
                </div>

                {/* Referral List */}
                <div className="bg-white rounded-2xl shadow-lg overflow-hidden animate-fade-in-up" style={{ animationDelay: '200ms' }}>
                    <div className="p-4 border-b">
                        <h3 className="font-bold">👥 Bạn bè đã mời ({referrals.length})</h3>
                    </div>
                    {referrals.length === 0 ? (
                        <div className="p-8 text-center text-gray-500">
                            <div className="text-4xl mb-2">🙋‍♂️</div>
                            <p>Chưa có ai. Mời bạn bè ngay!</p>
                        </div>
                    ) : (
                        <div className="divide-y">
                            {referrals.map(referral => (
                                <div key={referral.id} className="p-4 flex items-center gap-3">
                                    <div className="w-10 h-10 rounded-full overflow-hidden relative">
                                        <img src={referral.avatar} alt="" className="w-full h-full object-cover" />
                                    </div>
                                    <div className="flex-1">
                                        <div className="font-medium">{referral.name}</div>
                                        <div className="text-xs text-gray-400">
                                            Tham gia {new Date(referral.joinedAt).toLocaleDateString('vi-VN')}
                                        </div>
                                    </div>
                                    <div className={`px-3 py-1 rounded-full text-sm ${referral.status === 'completed'
                                            ? 'bg-green-100 text-green-700'
                                            : 'bg-yellow-100 text-yellow-700'
                                        }`}>
                                        {referral.status === 'completed' ? `+${referral.reward} xu` : 'Chờ đặt đơn'}
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </div>

                {/* Terms */}
                <div className="mt-6 text-center text-white/70 text-xs">
                    <p>Điều khoản: Bạn bè cần hoàn thành đơn hàng đầu tiên trong vòng 30 ngày.</p>
                    <p>Mỗi tài khoản có thể mời tối đa 50 người/tháng.</p>
                </div>
            </div>
        </div>
    );
}
