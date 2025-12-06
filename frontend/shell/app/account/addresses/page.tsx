'use client';

import React, { useState, useEffect } from 'react';
import { userService, Address } from '@/services/userService';

export default function AddressesPage() {
    const [addresses, setAddresses] = useState<Address[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [showForm, setShowForm] = useState(false);
    const [notification, setNotification] = useState<string | null>(null);
    const [newAddress, setNewAddress] = useState({
        name: '',
        phone: '',
        address: '',
        district: '',
        city: '',
        type: 'home' as 'home' | 'office',
        isDefault: false,
    });

    useEffect(() => {
        userService.getAddresses().then(data => {
            setAddresses(data);
            setIsLoading(false);
        });
    }, []);

    const handleAddAddress = async () => {
        try {
            const added = await userService.addAddress(newAddress);
            setAddresses([...addresses, added]);
            setShowForm(false);
            setNewAddress({
                name: '',
                phone: '',
                address: '',
                district: '',
                city: '',
                type: 'home',
                isDefault: false,
            });
            setNotification('✓ Đã thêm địa chỉ mới');
            setTimeout(() => setNotification(null), 2500);
        } catch (error) {
            console.error('Failed to add address:', error);
        }
    };

    if (isLoading) {
        return (
            <div className="bg-white rounded-sm shadow-sm p-6 animate-pulse">
                <div className="h-6 bg-gray-200 rounded w-1/4 mb-4" />
                <div className="space-y-4">
                    {[1, 2].map(i => (
                        <div key={i} className="h-24 bg-gray-200 rounded" />
                    ))}
                </div>
            </div>
        );
    }

    return (
        <div className="animate-fade-in">
            {/* Toast */}
            {notification && <div className="toast toast-success">{notification}</div>}

            {/* Header */}
            <div className="bg-white rounded-sm shadow-sm p-4 mb-4 flex items-center justify-between">
                <h1 className="text-lg font-medium">Địa Chỉ Của Tôi</h1>
                <button
                    onClick={() => setShowForm(true)}
                    className="px-4 py-2 bg-[#ee4d2d] text-white text-sm hover:opacity-90 transition-all hover-shrink flex items-center gap-2"
                >
                    <span>+</span> Thêm Địa Chỉ Mới
                </button>
            </div>

            {/* Add Address Form */}
            {showForm && (
                <div className="bg-white rounded-sm shadow-sm p-6 mb-4 animate-fade-in-down">
                    <h2 className="font-medium mb-4">Địa Chỉ Mới</h2>
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm text-gray-500 mb-1">Họ và tên</label>
                            <input
                                type="text"
                                value={newAddress.name}
                                onChange={(e) => setNewAddress({ ...newAddress, name: e.target.value })}
                                className="w-full border px-3 py-2 text-sm outline-none focus:border-[#ee4d2d] rounded-sm"
                                placeholder="Nhập họ và tên"
                            />
                        </div>
                        <div>
                            <label className="block text-sm text-gray-500 mb-1">Số điện thoại</label>
                            <input
                                type="tel"
                                value={newAddress.phone}
                                onChange={(e) => setNewAddress({ ...newAddress, phone: e.target.value })}
                                className="w-full border px-3 py-2 text-sm outline-none focus:border-[#ee4d2d] rounded-sm"
                                placeholder="Nhập số điện thoại"
                            />
                        </div>
                        <div className="col-span-2">
                            <label className="block text-sm text-gray-500 mb-1">Địa chỉ cụ thể</label>
                            <input
                                type="text"
                                value={newAddress.address}
                                onChange={(e) => setNewAddress({ ...newAddress, address: e.target.value })}
                                className="w-full border px-3 py-2 text-sm outline-none focus:border-[#ee4d2d] rounded-sm"
                                placeholder="Số nhà, tên đường"
                            />
                        </div>
                        <div>
                            <label className="block text-sm text-gray-500 mb-1">Quận/Huyện</label>
                            <input
                                type="text"
                                value={newAddress.district}
                                onChange={(e) => setNewAddress({ ...newAddress, district: e.target.value })}
                                className="w-full border px-3 py-2 text-sm outline-none focus:border-[#ee4d2d] rounded-sm"
                                placeholder="Nhập quận/huyện"
                            />
                        </div>
                        <div>
                            <label className="block text-sm text-gray-500 mb-1">Tỉnh/Thành phố</label>
                            <input
                                type="text"
                                value={newAddress.city}
                                onChange={(e) => setNewAddress({ ...newAddress, city: e.target.value })}
                                className="w-full border px-3 py-2 text-sm outline-none focus:border-[#ee4d2d] rounded-sm"
                                placeholder="Nhập tỉnh/thành phố"
                            />
                        </div>
                        <div className="col-span-2 flex items-center gap-6">
                            <label className="flex items-center gap-2 cursor-pointer">
                                <input
                                    type="radio"
                                    name="type"
                                    checked={newAddress.type === 'home'}
                                    onChange={() => setNewAddress({ ...newAddress, type: 'home' })}
                                    className="accent-[#ee4d2d]"
                                />
                                <span className="text-sm">🏠 Nhà riêng</span>
                            </label>
                            <label className="flex items-center gap-2 cursor-pointer">
                                <input
                                    type="radio"
                                    name="type"
                                    checked={newAddress.type === 'office'}
                                    onChange={() => setNewAddress({ ...newAddress, type: 'office' })}
                                    className="accent-[#ee4d2d]"
                                />
                                <span className="text-sm">🏢 Văn phòng</span>
                            </label>
                        </div>
                        <div className="col-span-2">
                            <label className="flex items-center gap-2 cursor-pointer">
                                <input
                                    type="checkbox"
                                    checked={newAddress.isDefault}
                                    onChange={(e) => setNewAddress({ ...newAddress, isDefault: e.target.checked })}
                                    className="accent-[#ee4d2d]"
                                />
                                <span className="text-sm">Đặt làm địa chỉ mặc định</span>
                            </label>
                        </div>
                        <div className="col-span-2 flex justify-end gap-2 pt-4">
                            <button
                                onClick={() => setShowForm(false)}
                                className="px-6 py-2 border text-sm hover:bg-gray-50 transition-all"
                            >
                                Hủy
                            </button>
                            <button
                                onClick={handleAddAddress}
                                className="px-6 py-2 bg-[#ee4d2d] text-white text-sm hover:opacity-90 transition-all"
                            >
                                Hoàn thành
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {/* Address List */}
            <div className="bg-white rounded-sm shadow-sm">
                <div className="p-4 border-b">
                    <h2 className="text-sm font-medium">Địa Chỉ</h2>
                </div>

                {addresses.length === 0 ? (
                    <div className="p-12 text-center">
                        <div className="text-5xl mb-4">📍</div>
                        <p className="text-gray-500">Bạn chưa có địa chỉ nào</p>
                    </div>
                ) : (
                    <div className="divide-y">
                        {addresses.map((addr, index) => (
                            <div
                                key={addr.id}
                                className="p-4 flex items-start gap-4 animate-fade-in-up hover:bg-gray-50 transition-colors"
                                style={{ animationDelay: `${index * 50}ms` }}
                            >
                                <div className="flex-1">
                                    <div className="flex items-center gap-2 mb-1">
                                        <span className="font-medium">{addr.name}</span>
                                        <span className="text-gray-300">|</span>
                                        <span className="text-gray-500">{addr.phone}</span>
                                        {addr.isDefault && (
                                            <span className="text-xs text-[#ee4d2d] border border-[#ee4d2d] px-2 py-0.5">
                                                Mặc định
                                            </span>
                                        )}
                                        <span className={`text-xs px-2 py-0.5 rounded ${addr.type === 'home' ? 'bg-green-50 text-green-600' : 'bg-blue-50 text-blue-600'
                                            }`}>
                                            {addr.type === 'home' ? '🏠 Nhà riêng' : '🏢 Văn phòng'}
                                        </span>
                                    </div>
                                    <p className="text-sm text-gray-500">
                                        {addr.address}, {addr.district}, {addr.city}
                                    </p>
                                </div>
                                <div className="flex flex-col gap-2">
                                    <div className="flex gap-2">
                                        <button className="text-sm text-[#4080ee] hover:underline">Cập nhật</button>
                                        {!addr.isDefault && (
                                            <button className="text-sm text-[#4080ee] hover:underline">Xóa</button>
                                        )}
                                    </div>
                                    {!addr.isDefault && (
                                        <button className="text-xs border px-2 py-1 hover:bg-gray-100 transition-all">
                                            Thiết lập mặc định
                                        </button>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
