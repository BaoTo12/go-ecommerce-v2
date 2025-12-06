'use client';

import React, { useState, useEffect } from 'react';
import Image from 'next/image';
import { userService, User } from '@/services/userService';

export default function ProfilePage() {
    const [user, setUser] = useState<User | null>(null);
    const [isEditing, setIsEditing] = useState(false);
    const [isSaving, setIsSaving] = useState(false);
    const [notification, setNotification] = useState<string | null>(null);
    const [formData, setFormData] = useState({
        name: '',
        email: '',
        phone: '',
    });

    useEffect(() => {
        userService.getCurrentUser().then(data => {
            setUser(data);
            setFormData({
                name: data.name,
                email: data.email,
                phone: data.phone,
            });
        });
    }, []);

    const handleSave = async () => {
        setIsSaving(true);
        try {
            const updated = await userService.updateProfile(formData);
            setUser(updated);
            setIsEditing(false);
            setNotification('✓ Cập nhật thông tin thành công');
            setTimeout(() => setNotification(null), 2500);
        } catch (error) {
            console.error('Failed to update profile:', error);
        } finally {
            setIsSaving(false);
        }
    };

    if (!user) {
        return (
            <div className="bg-white rounded-sm shadow-sm p-6">
                <div className="animate-pulse space-y-4">
                    <div className="h-6 bg-gray-200 rounded w-1/4" />
                    <div className="h-4 bg-gray-200 rounded w-1/2" />
                    <div className="h-32 bg-gray-200 rounded" />
                </div>
            </div>
        );
    }

    return (
        <div className="animate-fade-in">
            {/* Toast */}
            {notification && <div className="toast toast-success">{notification}</div>}

            {/* Header */}
            <div className="bg-white rounded-sm shadow-sm p-6 mb-4">
                <h1 className="text-xl font-medium">Hồ Sơ Của Tôi</h1>
                <p className="text-sm text-gray-500 mt-1">Quản lý thông tin hồ sơ để bảo mật tài khoản</p>
            </div>

            {/* Profile Form */}
            <div className="bg-white rounded-sm shadow-sm p-6">
                <div className="grid md:grid-cols-3 gap-8">
                    {/* Form */}
                    <div className="md:col-span-2 space-y-6">
                        <div className="flex items-center gap-4">
                            <label className="w-32 text-sm text-gray-500 text-right">Tên đăng nhập</label>
                            <span className="text-sm">{user.email.split('@')[0]}</span>
                        </div>

                        <div className="flex items-center gap-4">
                            <label className="w-32 text-sm text-gray-500 text-right">Tên</label>
                            {isEditing ? (
                                <input
                                    type="text"
                                    value={formData.name}
                                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                    className="flex-1 max-w-sm border px-3 py-2 text-sm outline-none focus:border-[#ee4d2d] rounded-sm"
                                />
                            ) : (
                                <span className="text-sm">{user.name}</span>
                            )}
                        </div>

                        <div className="flex items-center gap-4">
                            <label className="w-32 text-sm text-gray-500 text-right">Email</label>
                            {isEditing ? (
                                <input
                                    type="email"
                                    value={formData.email}
                                    onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                                    className="flex-1 max-w-sm border px-3 py-2 text-sm outline-none focus:border-[#ee4d2d] rounded-sm"
                                />
                            ) : (
                                <div className="flex items-center gap-2">
                                    <span className="text-sm">{user.email}</span>
                                    <span className="text-xs text-green-600 bg-green-50 px-2 py-0.5 rounded">Đã xác thực</span>
                                </div>
                            )}
                        </div>

                        <div className="flex items-center gap-4">
                            <label className="w-32 text-sm text-gray-500 text-right">Số điện thoại</label>
                            {isEditing ? (
                                <input
                                    type="tel"
                                    value={formData.phone}
                                    onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                                    className="flex-1 max-w-sm border px-3 py-2 text-sm outline-none focus:border-[#ee4d2d] rounded-sm"
                                />
                            ) : (
                                <div className="flex items-center gap-2">
                                    <span className="text-sm">{user.phone.replace(/(\d{4})(\d{3})(\d{3})/, '$1 $2 $3')}</span>
                                    <span className="text-xs text-green-600 bg-green-50 px-2 py-0.5 rounded">Đã xác thực</span>
                                </div>
                            )}
                        </div>

                        <div className="flex items-center gap-4">
                            <label className="w-32 text-sm text-gray-500 text-right">Ngày tham gia</label>
                            <span className="text-sm">{new Date(user.joinDate).toLocaleDateString('vi-VN', { year: 'numeric', month: 'long', day: 'numeric' })}</span>
                        </div>

                        <div className="flex items-center gap-4 pt-4">
                            <div className="w-32" />
                            {isEditing ? (
                                <div className="flex gap-2">
                                    <button
                                        onClick={handleSave}
                                        disabled={isSaving}
                                        className="px-6 py-2 bg-[#ee4d2d] text-white text-sm hover:opacity-90 transition-all hover-shrink"
                                    >
                                        {isSaving ? 'Đang lưu...' : 'Lưu'}
                                    </button>
                                    <button
                                        onClick={() => setIsEditing(false)}
                                        className="px-6 py-2 border text-sm hover:bg-gray-50 transition-all"
                                    >
                                        Hủy
                                    </button>
                                </div>
                            ) : (
                                <button
                                    onClick={() => setIsEditing(true)}
                                    className="px-6 py-2 border border-[#ee4d2d] text-[#ee4d2d] text-sm hover:bg-[#fef6f5] transition-all hover-shrink"
                                >
                                    ✏️ Chỉnh Sửa
                                </button>
                            )}
                        </div>
                    </div>

                    {/* Avatar */}
                    <div className="flex flex-col items-center pt-4 border-l pl-8">
                        <div className="relative w-24 h-24 rounded-full overflow-hidden border-4 border-[#ee4d2d] mb-4 hover:scale-105 transition-transform">
                            <Image
                                src={user.avatar}
                                alt={user.name}
                                fill
                                className="object-cover"
                                unoptimized
                            />
                        </div>
                        <button className="px-4 py-2 border text-sm hover:bg-gray-50 transition-all">
                            Chọn Ảnh
                        </button>
                        <p className="text-xs text-gray-400 mt-2 text-center">
                            Dung lượng file tối đa 1 MB<br />
                            Định dạng: JPEG, PNG
                        </p>
                    </div>
                </div>
            </div>

            {/* Security Section */}
            <div className="bg-white rounded-sm shadow-sm p-6 mt-4">
                <h2 className="font-medium mb-4">Bảo Mật</h2>
                <div className="space-y-4">
                    <div className="flex items-center justify-between py-3 border-b">
                        <div>
                            <p className="text-sm font-medium">Đổi mật khẩu</p>
                            <p className="text-xs text-gray-500">Thay đổi mật khẩu định kỳ để bảo mật tài khoản</p>
                        </div>
                        <button className="text-sm text-[#ee4d2d] hover:underline">Thay đổi</button>
                    </div>
                    <div className="flex items-center justify-between py-3 border-b">
                        <div>
                            <p className="text-sm font-medium">Xác minh 2 bước</p>
                            <p className="text-xs text-gray-500">Bảo vệ tài khoản với xác minh 2 bước</p>
                        </div>
                        <span className="text-xs text-green-600 bg-green-50 px-2 py-1 rounded">Đã bật</span>
                    </div>
                    <div className="flex items-center justify-between py-3">
                        <div>
                            <p className="text-sm font-medium">Liên kết tài khoản</p>
                            <p className="text-xs text-gray-500">Google, Facebook</p>
                        </div>
                        <button className="text-sm text-[#ee4d2d] hover:underline">Quản lý</button>
                    </div>
                </div>
            </div>
        </div>
    );
}
