'use client';

import React, { useState } from 'react';

export default function TwoFactorAuthPage() {
    const [is2FAEnabled, setIs2FAEnabled] = useState(false);
    const [step, setStep] = useState(1);
    const [verificationCode, setVerificationCode] = useState('');
    const [backupCodes, setBackupCodes] = useState<string[]>([]);
    const [notification, setNotification] = useState<string | null>(null);

    const generatedSecret = 'JBSWY3DPEHPK3PXP'; // Mock secret
    const qrCodeUrl = `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=otpauth://totp/Shopee:nguyenvana@email.com?secret=${generatedSecret}%26issuer=Shopee`;

    const enable2FA = () => {
        if (verificationCode.length !== 6) {
            setNotification('❌ Mã xác thực phải có 6 số');
            setTimeout(() => setNotification(null), 2000);
            return;
        }

        // Generate backup codes
        const codes = Array.from({ length: 8 }, () =>
            Math.random().toString(36).substring(2, 6).toUpperCase() + '-' +
            Math.random().toString(36).substring(2, 6).toUpperCase()
        );
        setBackupCodes(codes);
        setIs2FAEnabled(true);
        setStep(3);
        setNotification('✓ Đã bật xác thực 2 lớp thành công!');
        setTimeout(() => setNotification(null), 3000);
    };

    const disable2FA = () => {
        setIs2FAEnabled(false);
        setStep(1);
        setBackupCodes([]);
        setNotification('✓ Đã tắt xác thực 2 lớp');
        setTimeout(() => setNotification(null), 2000);
    };

    return (
        <div className="animate-fade-in">
            {notification && <div className="toast toast-success">{notification}</div>}

            <div className="bg-white dark:bg-gray-800 rounded-sm shadow-sm p-6 mb-4">
                <div className="flex items-center justify-between mb-4">
                    <div>
                        <h1 className="text-lg font-medium dark:text-white">🔐 Xác Thực 2 Lớp (2FA)</h1>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                            Bảo vệ tài khoản của bạn với lớp bảo mật bổ sung
                        </p>
                    </div>
                    <div className={`px-3 py-1 rounded-full text-sm ${is2FAEnabled
                            ? 'bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-300'
                            : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-400'
                        }`}>
                        {is2FAEnabled ? '✓ Đã bật' : 'Chưa bật'}
                    </div>
                </div>

                {!is2FAEnabled ? (
                    <div className="space-y-6">
                        {step === 1 && (
                            <div className="animate-fade-in">
                                <h3 className="font-medium mb-4 dark:text-white">Bước 1: Tải ứng dụng xác thực</h3>
                                <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
                                    Tải một trong các ứng dụng xác thực sau trên điện thoại của bạn:
                                </p>
                                <div className="flex gap-4 mb-6">
                                    <div className="flex-1 p-4 border dark:border-gray-700 rounded-lg text-center">
                                        <div className="text-3xl mb-2">🔐</div>
                                        <div className="font-medium dark:text-white">Google Authenticator</div>
                                        <div className="text-sm text-gray-500 dark:text-gray-400">iOS & Android</div>
                                    </div>
                                    <div className="flex-1 p-4 border dark:border-gray-700 rounded-lg text-center">
                                        <div className="text-3xl mb-2">🔒</div>
                                        <div className="font-medium dark:text-white">Microsoft Authenticator</div>
                                        <div className="text-sm text-gray-500 dark:text-gray-400">iOS & Android</div>
                                    </div>
                                    <div className="flex-1 p-4 border dark:border-gray-700 rounded-lg text-center">
                                        <div className="text-3xl mb-2">🛡️</div>
                                        <div className="font-medium dark:text-white">Authy</div>
                                        <div className="text-sm text-gray-500 dark:text-gray-400">iOS & Android</div>
                                    </div>
                                </div>
                                <button
                                    onClick={() => setStep(2)}
                                    className="px-6 py-3 bg-[#ee4d2d] text-white rounded-sm hover:opacity-90"
                                >
                                    Tiếp tục →
                                </button>
                            </div>
                        )}

                        {step === 2 && (
                            <div className="animate-fade-in">
                                <h3 className="font-medium mb-4 dark:text-white">Bước 2: Quét mã QR</h3>
                                <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
                                    Mở ứng dụng xác thực và quét mã QR bên dưới:
                                </p>

                                <div className="flex gap-8 items-start">
                                    <div className="bg-white p-4 rounded-lg shadow-sm">
                                        <img src={qrCodeUrl} alt="QR Code" className="w-48 h-48" />
                                    </div>
                                    <div className="flex-1">
                                        <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">
                                            Hoặc nhập mã thủ công:
                                        </p>
                                        <div className="bg-gray-100 dark:bg-gray-700 p-3 rounded-sm font-mono text-sm mb-4 dark:text-white">
                                            {generatedSecret}
                                        </div>

                                        <label className="block text-sm font-medium mb-2 dark:text-white">
                                            Nhập mã 6 số từ ứng dụng xác thực:
                                        </label>
                                        <div className="flex gap-2 mb-4">
                                            {[...Array(6)].map((_, i) => (
                                                <input
                                                    key={i}
                                                    type="text"
                                                    maxLength={1}
                                                    value={verificationCode[i] || ''}
                                                    onChange={(e) => {
                                                        const newCode = verificationCode.split('');
                                                        newCode[i] = e.target.value;
                                                        setVerificationCode(newCode.join(''));
                                                        if (e.target.value && e.target.nextElementSibling) {
                                                            (e.target.nextElementSibling as HTMLInputElement).focus();
                                                        }
                                                    }}
                                                    className="w-12 h-12 border dark:border-gray-600 text-center text-lg font-medium rounded-sm outline-none focus:border-[#ee4d2d] dark:bg-gray-700 dark:text-white"
                                                />
                                            ))}
                                        </div>

                                        <div className="flex gap-3">
                                            <button
                                                onClick={() => setStep(1)}
                                                className="px-4 py-2 border dark:border-gray-600 rounded-sm hover:bg-gray-50 dark:hover:bg-gray-700 dark:text-white"
                                            >
                                                ← Quay lại
                                            </button>
                                            <button
                                                onClick={enable2FA}
                                                className="px-6 py-2 bg-[#ee4d2d] text-white rounded-sm hover:opacity-90"
                                            >
                                                Xác nhận
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        )}
                    </div>
                ) : step === 3 ? (
                    <div className="animate-fade-in">
                        <h3 className="font-medium mb-4 text-green-600">✓ Xác thực 2 lớp đã được bật!</h3>

                        <div className="bg-yellow-50 dark:bg-yellow-900 border border-yellow-200 dark:border-yellow-700 rounded-lg p-4 mb-6">
                            <h4 className="font-medium text-yellow-700 dark:text-yellow-300 mb-2">
                                ⚠️ Lưu mã khôi phục
                            </h4>
                            <p className="text-sm text-yellow-600 dark:text-yellow-400 mb-4">
                                Lưu giữ các mã này ở nơi an toàn. Bạn sẽ cần chúng nếu mất quyền truy cập vào ứng dụng xác thực.
                            </p>
                            <div className="grid grid-cols-4 gap-2 mb-4">
                                {backupCodes.map((code, i) => (
                                    <div key={i} className="bg-white dark:bg-gray-800 p-2 text-center font-mono text-sm rounded dark:text-white">
                                        {code}
                                    </div>
                                ))}
                            </div>
                            <button
                                onClick={() => {
                                    navigator.clipboard.writeText(backupCodes.join('\n'));
                                    setNotification('✓ Đã sao chép mã khôi phục');
                                    setTimeout(() => setNotification(null), 2000);
                                }}
                                className="text-sm text-[#ee4d2d] hover:underline"
                            >
                                📋 Sao chép tất cả
                            </button>
                        </div>

                        <button
                            onClick={() => setStep(4)}
                            className="px-6 py-2 bg-[#ee4d2d] text-white rounded-sm hover:opacity-90"
                        >
                            Hoàn tất
                        </button>
                    </div>
                ) : (
                    <div className="animate-fade-in">
                        <div className="p-4 bg-green-50 dark:bg-green-900 rounded-lg mb-4">
                            <p className="text-green-700 dark:text-green-300">
                                ✓ Tài khoản của bạn đang được bảo vệ bởi xác thực 2 lớp
                            </p>
                        </div>
                        <button
                            onClick={disable2FA}
                            className="px-4 py-2 border border-red-500 text-red-500 rounded-sm hover:bg-red-50 dark:hover:bg-red-900"
                        >
                            Tắt xác thực 2 lớp
                        </button>
                    </div>
                )}
            </div>

            {/* Security Tips */}
            <div className="bg-white dark:bg-gray-800 rounded-sm shadow-sm p-6">
                <h3 className="font-medium mb-4 dark:text-white">💡 Mẹo bảo mật</h3>
                <ul className="space-y-2 text-sm text-gray-600 dark:text-gray-400">
                    <li>• Không chia sẻ mã 2FA với bất kỳ ai</li>
                    <li>• Lưu trữ mã khôi phục ở nơi an toàn</li>
                    <li>• Sử dụng mật khẩu mạnh và duy nhất</li>
                    <li>• Đăng xuất khi sử dụng thiết bị công cộng</li>
                    <li>• Thường xuyên kiểm tra hoạt động tài khoản</li>
                </ul>
            </div>
        </div>
    );
}
