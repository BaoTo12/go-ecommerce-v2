'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { authService } from '@/services/authService';

export default function LoginPage() {
    const router = useRouter();
    const [email, setEmail] = useState('nguyenvana@email.com');
    const [password, setPassword] = useState('123456');
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [showPassword, setShowPassword] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setIsLoading(true);

        try {
            const result = await authService.login({ email, password });
            if (result.success) {
                router.push('/');
            } else {
                setError(result.error || 'Đăng nhập thất bại');
            }
        } catch (err) {
            setError('Có lỗi xảy ra. Vui lòng thử lại.');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-[#ee4d2d] flex items-center justify-center p-4">
            <div className="w-full max-w-md animate-fade-in">
                {/* Logo */}
                <div className="text-center mb-8">
                    <Link href="/" className="inline-block">
                        <h1 className="text-white text-4xl font-bold">Shopee</h1>
                    </Link>
                    <p className="text-white/80 mt-2">Nền tảng thương mại điện tử hàng đầu</p>
                </div>

                {/* Login Form */}
                <div className="bg-white rounded-sm shadow-lg p-8 animate-fade-in-up">
                    <h2 className="text-xl font-medium mb-6">Đăng Nhập</h2>

                    {error && (
                        <div className="bg-red-50 border border-red-200 text-red-600 text-sm p-3 rounded-sm mb-4 animate-shake">
                            {error}
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div>
                            <input
                                type="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                placeholder="Email/Số điện thoại"
                                className="w-full border px-4 py-3 text-sm outline-none focus:border-[#ee4d2d] rounded-sm transition-all"
                                required
                            />
                        </div>

                        <div className="relative">
                            <input
                                type={showPassword ? 'text' : 'password'}
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                placeholder="Mật khẩu"
                                className="w-full border px-4 py-3 text-sm outline-none focus:border-[#ee4d2d] rounded-sm transition-all pr-12"
                                required
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                            >
                                {showPassword ? '🙈' : '👁️'}
                            </button>
                        </div>

                        <button
                            type="submit"
                            disabled={isLoading}
                            className={`w-full py-3 bg-[#ee4d2d] text-white font-medium rounded-sm hover:opacity-90 transition-all ${isLoading ? 'opacity-70 cursor-wait' : ''
                                }`}
                        >
                            {isLoading ? (
                                <span className="flex items-center justify-center gap-2">
                                    <span className="loading-spinner" /> Đang đăng nhập...
                                </span>
                            ) : (
                                'ĐĂNG NHẬP'
                            )}
                        </button>
                    </form>

                    <div className="flex items-center justify-between mt-4 text-sm">
                        <Link href="/auth/forgot-password" className="text-[#0055aa] hover:text-[#ee4d2d]">
                            Quên mật khẩu
                        </Link>
                        <Link href="/auth/sms-login" className="text-[#0055aa] hover:text-[#ee4d2d]">
                            Đăng nhập với SMS
                        </Link>
                    </div>

                    <div className="flex items-center my-6">
                        <div className="flex-1 border-t border-gray-200" />
                        <span className="px-4 text-gray-400 text-sm">HOẶC</span>
                        <div className="flex-1 border-t border-gray-200" />
                    </div>

                    {/* Social Login */}
                    <div className="grid grid-cols-2 gap-3">
                        <button className="flex items-center justify-center gap-2 py-3 border rounded-sm hover:bg-gray-50 transition-all">
                            <svg className="w-5 h-5" viewBox="0 0 24 24" fill="#1877f2">
                                <path d="M12 2C6.477 2 2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.879V14.89h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.89h-2.33v6.989C18.343 21.129 22 16.99 22 12c0-5.523-4.477-10-10-10z" />
                            </svg>
                            <span className="text-sm">Facebook</span>
                        </button>
                        <button className="flex items-center justify-center gap-2 py-3 border rounded-sm hover:bg-gray-50 transition-all">
                            <svg className="w-5 h-5" viewBox="0 0 24 24">
                                <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
                                <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
                                <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" />
                                <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
                            </svg>
                            <span className="text-sm">Google</span>
                        </button>
                    </div>

                    <div className="text-center mt-6 text-sm text-gray-500">
                        Bạn mới biết đến Shopee?{' '}
                        <Link href="/auth/register" className="text-[#ee4d2d] hover:underline font-medium">
                            Đăng ký
                        </Link>
                    </div>
                </div>

                {/* Demo credentials */}
                <div className="mt-4 p-3 bg-white/10 rounded-sm text-white text-sm text-center">
                    <p>🔑 Demo: nguyenvana@email.com / 123456</p>
                </div>
            </div>
        </div>
    );
}
