'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { authService } from '@/services/authService';

export default function RegisterPage() {
    const router = useRouter();
    const [step, setStep] = useState(1); // 1: form, 2: OTP verification
    const [formData, setFormData] = useState({
        name: '',
        email: '',
        phone: '',
        password: '',
        confirmPassword: '',
    });
    const [otp, setOtp] = useState('');
    const [generatedOtp, setGeneratedOtp] = useState('');
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [showPassword, setShowPassword] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (formData.password !== formData.confirmPassword) {
            setError('Mật khẩu xác nhận không khớp');
            return;
        }

        if (formData.password.length < 6) {
            setError('Mật khẩu phải có ít nhất 6 ký tự');
            return;
        }

        setIsLoading(true);

        try {
            const result = await authService.register({
                name: formData.name,
                email: formData.email,
                phone: formData.phone,
                password: formData.password,
            });

            if (result.success) {
                // Send OTP for verification
                const otpResult = await authService.sendOTP(formData.phone);
                if (otpResult.success && otpResult.otp) {
                    setGeneratedOtp(otpResult.otp);
                    setStep(2);
                }
            } else {
                setError(result.error || 'Đăng ký thất bại');
            }
        } catch (err) {
            setError('Có lỗi xảy ra. Vui lòng thử lại.');
        } finally {
            setIsLoading(false);
        }
    };

    const handleVerifyOTP = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setIsLoading(true);

        try {
            const result = await authService.verifyOTP(formData.phone, otp);
            if (result.success) {
                router.push('/');
            } else {
                setError('Mã OTP không đúng');
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
                    <p className="text-white/80 mt-2">Đăng ký tài khoản mới</p>
                </div>

                {/* Register Form */}
                <div className="bg-white rounded-sm shadow-lg p-8 animate-fade-in-up">
                    {step === 1 ? (
                        <>
                            <h2 className="text-xl font-medium mb-6">Đăng Ký</h2>

                            {error && (
                                <div className="bg-red-50 border border-red-200 text-red-600 text-sm p-3 rounded-sm mb-4 animate-shake">
                                    {error}
                                </div>
                            )}

                            <form onSubmit={handleSubmit} className="space-y-4">
                                <div>
                                    <input
                                        type="text"
                                        value={formData.name}
                                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                        placeholder="Họ và tên"
                                        className="w-full border px-4 py-3 text-sm outline-none focus:border-[#ee4d2d] rounded-sm transition-all"
                                        required
                                    />
                                </div>

                                <div>
                                    <input
                                        type="email"
                                        value={formData.email}
                                        onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                                        placeholder="Email"
                                        className="w-full border px-4 py-3 text-sm outline-none focus:border-[#ee4d2d] rounded-sm transition-all"
                                        required
                                    />
                                </div>

                                <div>
                                    <input
                                        type="tel"
                                        value={formData.phone}
                                        onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                                        placeholder="Số điện thoại"
                                        className="w-full border px-4 py-3 text-sm outline-none focus:border-[#ee4d2d] rounded-sm transition-all"
                                        required
                                    />
                                </div>

                                <div className="relative">
                                    <input
                                        type={showPassword ? 'text' : 'password'}
                                        value={formData.password}
                                        onChange={(e) => setFormData({ ...formData, password: e.target.value })}
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

                                <div>
                                    <input
                                        type="password"
                                        value={formData.confirmPassword}
                                        onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
                                        placeholder="Xác nhận mật khẩu"
                                        className="w-full border px-4 py-3 text-sm outline-none focus:border-[#ee4d2d] rounded-sm transition-all"
                                        required
                                    />
                                </div>

                                <button
                                    type="submit"
                                    disabled={isLoading}
                                    className={`w-full py-3 bg-[#ee4d2d] text-white font-medium rounded-sm hover:opacity-90 transition-all ${isLoading ? 'opacity-70 cursor-wait' : ''
                                        }`}
                                >
                                    {isLoading ? (
                                        <span className="flex items-center justify-center gap-2">
                                            <span className="loading-spinner" /> Đang xử lý...
                                        </span>
                                    ) : (
                                        'ĐĂNG KÝ'
                                    )}
                                </button>
                            </form>

                            <p className="text-xs text-gray-500 mt-4 text-center">
                                Bằng việc đăng kí, bạn đã đồng ý với Shopee về{' '}
                                <Link href="#" className="text-[#ee4d2d]">Điều khoản dịch vụ</Link>
                                {' & '}
                                <Link href="#" className="text-[#ee4d2d]">Chính sách bảo mật</Link>
                            </p>

                            <div className="flex items-center my-6">
                                <div className="flex-1 border-t border-gray-200" />
                                <span className="px-4 text-gray-400 text-sm">HOẶC</span>
                                <div className="flex-1 border-t border-gray-200" />
                            </div>

                            {/* Social Login */}
                            <div className="grid grid-cols-2 gap-3">
                                <button className="flex items-center justify-center gap-2 py-3 border rounded-sm hover:bg-gray-50 transition-all">
                                    <span className="text-blue-600">📘</span>
                                    <span className="text-sm">Facebook</span>
                                </button>
                                <button className="flex items-center justify-center gap-2 py-3 border rounded-sm hover:bg-gray-50 transition-all">
                                    <span>🔵</span>
                                    <span className="text-sm">Google</span>
                                </button>
                            </div>

                            <div className="text-center mt-6 text-sm text-gray-500">
                                Bạn đã có tài khoản?{' '}
                                <Link href="/auth/login" className="text-[#ee4d2d] hover:underline font-medium">
                                    Đăng nhập
                                </Link>
                            </div>
                        </>
                    ) : (
                        <>
                            <h2 className="text-xl font-medium mb-2">Xác Minh Số Điện Thoại</h2>
                            <p className="text-sm text-gray-500 mb-6">
                                Mã xác minh đã được gửi đến {formData.phone}
                            </p>

                            {/* Demo OTP display */}
                            <div className="bg-green-50 border border-green-200 text-green-700 text-sm p-3 rounded-sm mb-4">
                                🔑 Demo OTP: {generatedOtp}
                            </div>

                            {error && (
                                <div className="bg-red-50 border border-red-200 text-red-600 text-sm p-3 rounded-sm mb-4 animate-shake">
                                    {error}
                                </div>
                            )}

                            <form onSubmit={handleVerifyOTP} className="space-y-4">
                                <div className="flex gap-2 justify-center">
                                    {[...Array(6)].map((_, i) => (
                                        <input
                                            key={i}
                                            type="text"
                                            maxLength={1}
                                            value={otp[i] || ''}
                                            onChange={(e) => {
                                                const newOtp = otp.split('');
                                                newOtp[i] = e.target.value;
                                                setOtp(newOtp.join(''));
                                                if (e.target.value && e.target.nextElementSibling) {
                                                    (e.target.nextElementSibling as HTMLInputElement).focus();
                                                }
                                            }}
                                            className="w-12 h-12 border text-center text-lg font-medium outline-none focus:border-[#ee4d2d] rounded-sm"
                                        />
                                    ))}
                                </div>

                                <button
                                    type="submit"
                                    disabled={isLoading || otp.length !== 6}
                                    className={`w-full py-3 bg-[#ee4d2d] text-white font-medium rounded-sm hover:opacity-90 transition-all ${isLoading || otp.length !== 6 ? 'opacity-70 cursor-not-allowed' : ''
                                        }`}
                                >
                                    {isLoading ? 'Đang xác minh...' : 'XÁC NHẬN'}
                                </button>

                                <button
                                    type="button"
                                    onClick={() => setStep(1)}
                                    className="w-full py-3 border text-gray-600 rounded-sm hover:bg-gray-50 transition-all"
                                >
                                    Quay lại
                                </button>
                            </form>

                            <p className="text-center text-sm text-gray-500 mt-4">
                                Không nhận được mã?{' '}
                                <button className="text-[#ee4d2d] hover:underline">Gửi lại</button>
                            </p>
                        </>
                    )}
                </div>
            </div>
        </div>
    );
}
