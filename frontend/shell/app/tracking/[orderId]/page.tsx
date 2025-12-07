'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';

interface TrackingStep {
    status: string;
    description: string;
    time: string;
    location?: string;
    completed: boolean;
}

interface OrderTracking {
    orderId: string;
    status: 'processing' | 'shipped' | 'in_transit' | 'out_for_delivery' | 'delivered';
    estimatedDelivery: string;
    carrier: string;
    trackingNumber: string;
    currentLocation: { lat: number; lng: number; address: string };
    steps: TrackingStep[];
}

const MOCK_ORDER: OrderTracking = {
    orderId: 'SPX2024120712345',
    status: 'out_for_delivery',
    estimatedDelivery: 'Hôm nay, 14:00 - 18:00',
    carrier: 'Shopee Express',
    trackingNumber: 'SPXVN123456789',
    currentLocation: {
        lat: 10.7769,
        lng: 106.7009,
        address: 'Quận 1, TP. Hồ Chí Minh',
    },
    steps: [
        { status: 'Đang giao', description: 'Shipper đang trên đường giao hàng đến bạn', time: '10:30', location: 'Quận 1, TP.HCM', completed: true },
        { status: 'Đến kho phát', description: 'Đơn hàng đã đến bưu cục phát', time: '08:15', location: 'Quận 1, TP.HCM', completed: true },
        { status: 'Đang vận chuyển', description: 'Đơn hàng đang trên đường vận chuyển', time: '06:00', location: 'Bình Dương', completed: true },
        { status: 'Rời kho phân loại', description: 'Đơn hàng đã rời kho phân loại', time: 'Hôm qua, 22:00', location: 'Kho Long An', completed: true },
        { status: 'Đến kho phân loại', description: 'Đơn hàng đã đến kho phân loại', time: 'Hôm qua, 18:00', location: 'Kho Long An', completed: true },
        { status: 'Đã lấy hàng', description: 'Shipper đã lấy hàng từ người bán', time: 'Hôm qua, 15:30', location: 'Quận 7, TP.HCM', completed: true },
        { status: 'Đơn hàng đã xác nhận', description: 'Người bán đã xác nhận đơn hàng', time: 'Hôm qua, 14:00', completed: true },
    ],
};

export default function LiveTrackingPage() {
    const [order] = useState(MOCK_ORDER);
    const [driverLocation, setDriverLocation] = useState({ lat: 10.775, lng: 106.698 });
    const [eta, setEta] = useState(25);

    // Simulate driver movement
    useEffect(() => {
        const interval = setInterval(() => {
            setDriverLocation(prev => ({
                lat: prev.lat + (Math.random() - 0.5) * 0.002,
                lng: prev.lng + (Math.random() - 0.5) * 0.002,
            }));
            setEta(prev => Math.max(1, prev - 0.5));
        }, 3000);

        return () => clearInterval(interval);
    }, []);

    const statusColors: Record<string, string> = {
        processing: 'bg-yellow-500',
        shipped: 'bg-blue-500',
        in_transit: 'bg-purple-500',
        out_for_delivery: 'bg-orange-500',
        delivered: 'bg-green-500',
    };

    const statusLabels: Record<string, string> = {
        processing: 'Đang xử lý',
        shipped: 'Đã giao cho ĐVVC',
        in_transit: 'Đang vận chuyển',
        out_for_delivery: 'Đang giao hàng',
        delivered: 'Đã giao',
    };

    return (
        <div className="min-h-screen bg-gray-100 dark:bg-gray-900">
            {/* Map Placeholder */}
            <div className="h-[300px] bg-gradient-to-b from-blue-400 to-blue-600 relative overflow-hidden">
                {/* Simulated Map */}
                <div className="absolute inset-0 bg-[url('data:image/svg+xml,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22100%22%20height%3D%22100%22%3E%3Crect%20fill%3D%22%23e0e0e0%22%20width%3D%22100%22%20height%3D%22100%22%2F%3E%3Cpath%20fill%3D%22%23fff%22%20d%3D%22M0%2020h100M0%2040h100M0%2060h100M0%2080h100M20%200v100M40%200v100M60%200v100M80%200v100%22%20stroke%3D%22%23ccc%22%20stroke-width%3D%220.5%22%2F%3E%3C%2Fsvg%3E')] opacity-30" />

                {/* Driver Marker */}
                <div
                    className="absolute w-12 h-12 -translate-x-1/2 -translate-y-1/2 animate-bounce"
                    style={{
                        top: `${50 + (driverLocation.lat - 10.775) * 500}%`,
                        left: `${50 + (driverLocation.lng - 106.7) * 500}%`,
                    }}
                >
                    <div className="w-12 h-12 bg-[#ee4d2d] rounded-full flex items-center justify-center text-white text-2xl shadow-lg">
                        🛵
                    </div>
                </div>

                {/* Destination Marker */}
                <div className="absolute bottom-20 right-20">
                    <div className="w-10 h-10 bg-green-500 rounded-full flex items-center justify-center text-white text-xl shadow-lg">
                        📍
                    </div>
                </div>

                {/* ETA Card */}
                <div className="absolute top-4 left-4 bg-white rounded-lg shadow-lg p-3">
                    <div className="text-sm text-gray-500">Thời gian dự kiến</div>
                    <div className="text-2xl font-bold text-[#ee4d2d]">{Math.round(eta)} phút</div>
                </div>

                {/* Driver Info Card */}
                <div className="absolute bottom-4 left-4 right-4 bg-white rounded-lg shadow-lg p-4">
                    <div className="flex items-center gap-3">
                        <div className="w-12 h-12 bg-gray-200 rounded-full flex items-center justify-center text-2xl">
                            🧑
                        </div>
                        <div className="flex-1">
                            <div className="font-bold">Nguyễn Văn Shipper</div>
                            <div className="text-sm text-gray-500">59A1-12345 • Honda Vision</div>
                        </div>
                        <div className="flex gap-2">
                            <button className="w-10 h-10 bg-green-500 text-white rounded-full flex items-center justify-center">
                                📞
                            </button>
                            <button className="w-10 h-10 bg-blue-500 text-white rounded-full flex items-center justify-center">
                                💬
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            {/* Order Details */}
            <div className="container mx-auto px-4 py-6">
                {/* Status Bar */}
                <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-4 mb-4">
                    <div className="flex items-center justify-between mb-3">
                        <div className="flex items-center gap-2">
                            <span className={`w-3 h-3 rounded-full ${statusColors[order.status]} animate-pulse`} />
                            <span className="font-bold dark:text-white">{statusLabels[order.status]}</span>
                        </div>
                        <span className="text-sm text-gray-500">{order.trackingNumber}</span>
                    </div>

                    {/* Progress Bar */}
                    <div className="flex items-center gap-1">
                        {['processing', 'shipped', 'in_transit', 'out_for_delivery', 'delivered'].map((step, i, arr) => {
                            const isCompleted = arr.indexOf(order.status) >= i;
                            return (
                                <React.Fragment key={step}>
                                    <div className={`w-6 h-6 rounded-full flex items-center justify-center text-xs ${isCompleted ? 'bg-green-500 text-white' : 'bg-gray-200 text-gray-400'
                                        }`}>
                                        {isCompleted ? '✓' : i + 1}
                                    </div>
                                    {i < arr.length - 1 && (
                                        <div className={`flex-1 h-1 ${isCompleted ? 'bg-green-500' : 'bg-gray-200'}`} />
                                    )}
                                </React.Fragment>
                            );
                        })}
                    </div>

                    <div className="mt-3 text-sm text-gray-600 dark:text-gray-400">
                        📦 Dự kiến giao: <span className="font-medium text-[#ee4d2d]">{order.estimatedDelivery}</span>
                    </div>
                </div>

                {/* Timeline */}
                <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-4">
                    <h3 className="font-bold mb-4 dark:text-white">📋 Chi tiết vận chuyển</h3>

                    <div className="space-y-0">
                        {order.steps.map((step, index) => (
                            <div key={index} className="flex gap-3">
                                <div className="flex flex-col items-center">
                                    <div className={`w-3 h-3 rounded-full ${index === 0 ? 'bg-[#ee4d2d] animate-pulse' : step.completed ? 'bg-green-500' : 'bg-gray-300'
                                        }`} />
                                    {index < order.steps.length - 1 && (
                                        <div className={`w-0.5 h-16 ${step.completed ? 'bg-green-500' : 'bg-gray-300'}`} />
                                    )}
                                </div>
                                <div className="flex-1 pb-4">
                                    <div className={`font-medium ${index === 0 ? 'text-[#ee4d2d]' : 'dark:text-white'}`}>
                                        {step.status}
                                    </div>
                                    <div className="text-sm text-gray-500 dark:text-gray-400">{step.description}</div>
                                    <div className="text-xs text-gray-400 mt-1">
                                        {step.time} {step.location && `• ${step.location}`}
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Actions */}
                <div className="mt-4 grid grid-cols-2 gap-3">
                    <button className="py-3 bg-gray-200 dark:bg-gray-700 rounded-lg font-medium dark:text-white">
                        📍 Báo sai địa chỉ
                    </button>
                    <button className="py-3 bg-[#ee4d2d] text-white rounded-lg font-medium">
                        📦 Xác nhận đã nhận
                    </button>
                </div>
            </div>
        </div>
    );
}
