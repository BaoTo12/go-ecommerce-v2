'use client';

import React, { useState, useEffect } from 'react';

interface FraudAlert {
    id: string;
    severity: 'high' | 'medium' | 'low';
    message: string;
    transactionId: string;
    amount: number;
    time: string;
}

interface FraudCheck {
    id: string;
    transactionId: string;
    score: number;
    riskLevel: string;
    decision: string;
    reasons: string[];
    processingTime: number;
}

export default function FraudDashboard() {
    const [alerts, setAlerts] = useState<FraudAlert[]>([]);
    const [recentChecks, setRecentChecks] = useState<FraudCheck[]>([]);
    const [testResult, setTestResult] = useState<FraudCheck | null>(null);
    const [testing, setTesting] = useState(false);

    useEffect(() => {
        setAlerts([
            { id: 'a1', severity: 'high', message: 'Nhiều giao dịch thất bại từ thiết bị mới', transactionId: 'TXN-12345', amount: 15000000, time: '2 phút trước' },
            { id: 'a2', severity: 'medium', message: 'Tần suất đặt hàng bất thường', transactionId: 'TXN-12346', amount: 8500000, time: '5 phút trước' },
        ]);

        setRecentChecks([
            { id: 'fc1', transactionId: 'TXN-12345', score: 0.85, riskLevel: 'HIGH', decision: 'REVIEW', reasons: ['Thiết bị mới', 'Số tiền lớn'], processingTime: 12 },
            { id: 'fc2', transactionId: 'TXN-12346', score: 0.45, riskLevel: 'MEDIUM', decision: 'ALLOW', reasons: ['Tần suất cao'], processingTime: 8 },
            { id: 'fc3', transactionId: 'TXN-12347', score: 0.12, riskLevel: 'LOW', decision: 'ALLOW', reasons: [], processingTime: 5 },
        ]);
    }, []);

    const runTest = () => {
        setTesting(true);
        setTimeout(() => {
            setTestResult({
                id: 'test',
                transactionId: `TXN-TEST-${Date.now()}`,
                score: Math.random(),
                riskLevel: Math.random() > 0.7 ? 'HIGH' : Math.random() > 0.4 ? 'MEDIUM' : 'LOW',
                decision: Math.random() > 0.8 ? 'BLOCK' : 'ALLOW',
                reasons: ['Giao dịch test', 'Thiết bị mới'],
                processingTime: Math.floor(Math.random() * 20) + 5,
            });
            setTesting(false);
        }, 1500);
    };

    const getRiskColor = (level: string) => {
        switch (level) {
            case 'HIGH': return 'bg-red-100 text-red-700';
            case 'MEDIUM': return 'bg-yellow-100 text-yellow-700';
            case 'LOW': return 'bg-green-100 text-green-700';
            default: return 'bg-gray-100 text-gray-700';
        }
    };

    return (
        <div className="min-h-screen bg-[#F5F5F5]">
            {/* Header */}
            <div className="bg-gradient-to-r from-slate-800 to-slate-900 text-white px-6 py-6">
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-2xl font-bold">🛡️ Fraud Detection Center</h1>
                        <p className="text-slate-300 text-sm">Giám sát giao dịch thời gian thực với ML</p>
                    </div>
                    <div className="flex gap-6 text-center">
                        <div>
                            <div className="text-2xl font-bold">99.7%</div>
                            <div className="text-xs text-slate-400">Tỷ lệ phát hiện</div>
                        </div>
                        <div>
                            <div className="text-2xl font-bold">8ms</div>
                            <div className="text-xs text-slate-400">Độ trễ TB</div>
                        </div>
                        <div>
                            <div className="text-2xl font-bold">0.02%</div>
                            <div className="text-xs text-slate-400">False Positive</div>
                        </div>
                    </div>
                </div>
            </div>

            <div className="p-6">
                {/* Alerts */}
                {alerts.length > 0 && (
                    <div className="mb-6">
                        <h2 className="font-bold text-lg mb-3 text-red-600">🚨 Cảnh báo ({alerts.length})</h2>
                        <div className="space-y-3">
                            {alerts.map(alert => (
                                <div
                                    key={alert.id}
                                    className={`p-4 rounded border-l-4 ${alert.severity === 'high' ? 'bg-red-50 border-red-500' : 'bg-yellow-50 border-yellow-500'
                                        }`}
                                >
                                    <div className="flex items-center justify-between">
                                        <div>
                                            <span className={`inline-block w-2 h-2 rounded-full mr-2 ${alert.severity === 'high' ? 'bg-red-500' : 'bg-yellow-500'
                                                }`} />
                                            <span className="font-semibold">{alert.message}</span>
                                            <div className="text-sm text-gray-500 mt-1">
                                                {alert.transactionId} | {alert.amount.toLocaleString()}₫ | {alert.time}
                                            </div>
                                        </div>
                                        <button className="bg-white px-4 py-2 rounded text-sm font-semibold hover:bg-gray-50 border">
                                            Xem xét
                                        </button>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}

                <div className="grid lg:grid-cols-3 gap-6">
                    {/* Test Transaction */}
                    <div className="bg-white rounded p-6">
                        <h2 className="font-bold text-lg mb-4">🧪 Kiểm tra giao dịch</h2>
                        <p className="text-sm text-gray-500 mb-4">Chạy giao dịch thử qua hệ thống phát hiện gian lận</p>

                        <button
                            onClick={runTest}
                            disabled={testing}
                            className={`w-full py-3 rounded font-bold text-white ${testing ? 'bg-gray-400' : 'bg-gradient-to-r from-blue-600 to-indigo-600 hover:opacity-90'
                                }`}
                        >
                            {testing ? '⏳ Đang phân tích...' : '🔍 Chạy kiểm tra'}
                        </button>

                        {testResult && (
                            <div className="mt-4 p-4 border rounded">
                                <div className="flex items-center justify-between mb-3">
                                    <span className={`px-3 py-1 rounded text-sm font-bold ${getRiskColor(testResult.riskLevel)}`}>
                                        {testResult.riskLevel}
                                    </span>
                                    <span className={`font-bold ${testResult.decision === 'BLOCK' ? 'text-red-600' : 'text-green-600'
                                        }`}>
                                        {testResult.decision}
                                    </span>
                                </div>
                                <div className="space-y-2 text-sm">
                                    <div className="flex justify-between">
                                        <span className="text-gray-500">Điểm rủi ro</span>
                                        <span className="font-mono">{(testResult.score * 100).toFixed(1)}%</span>
                                    </div>
                                    <div className="flex justify-between">
                                        <span className="text-gray-500">Thời gian xử lý</span>
                                        <span className="font-mono">{testResult.processingTime}ms</span>
                                    </div>
                                    {testResult.reasons.length > 0 && (
                                        <div>
                                            <span className="text-gray-500">Cảnh báo:</span>
                                            <ul className="list-disc list-inside text-yellow-600">
                                                {testResult.reasons.map((r, i) => <li key={i}>{r}</li>)}
                                            </ul>
                                        </div>
                                    )}
                                </div>
                            </div>
                        )}
                    </div>

                    {/* Recent Checks */}
                    <div className="lg:col-span-2 bg-white rounded p-6">
                        <h2 className="font-bold text-lg mb-4">📋 Kiểm tra gần đây</h2>
                        <div className="overflow-x-auto">
                            <table className="w-full">
                                <thead>
                                    <tr className="border-b text-left text-sm text-gray-500">
                                        <th className="pb-3">Giao dịch</th>
                                        <th className="pb-3">Điểm</th>
                                        <th className="pb-3">Rủi ro</th>
                                        <th className="pb-3">Quyết định</th>
                                        <th className="pb-3">Thời gian</th>
                                    </tr>
                                </thead>
                                <tbody className="text-sm">
                                    {recentChecks.map(check => (
                                        <tr key={check.id} className="border-b hover:bg-gray-50">
                                            <td className="py-3 font-mono">{check.transactionId}</td>
                                            <td className="py-3">
                                                <div className="flex items-center gap-2">
                                                    <div className="w-16 h-2 bg-gray-200 rounded-full overflow-hidden">
                                                        <div
                                                            className={`h-full rounded-full ${check.score > 0.7 ? 'bg-red-500' : check.score > 0.4 ? 'bg-yellow-500' : 'bg-green-500'
                                                                }`}
                                                            style={{ width: `${check.score * 100}%` }}
                                                        />
                                                    </div>
                                                    <span>{(check.score * 100).toFixed(0)}%</span>
                                                </div>
                                            </td>
                                            <td className="py-3">
                                                <span className={`px-2 py-1 rounded text-xs font-bold ${getRiskColor(check.riskLevel)}`}>
                                                    {check.riskLevel}
                                                </span>
                                            </td>
                                            <td className={`py-3 font-semibold ${check.decision === 'BLOCK' ? 'text-red-600' : check.decision === 'REVIEW' ? 'text-yellow-600' : 'text-green-600'
                                                }`}>
                                                {check.decision}
                                            </td>
                                            <td className="py-3 text-gray-500">{check.processingTime}ms</td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>

                {/* ML Features */}
                <div className="mt-6 bg-white rounded p-6">
                    <h2 className="font-bold text-lg mb-4">🤖 Đặc trưng ML</h2>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        {[
                            { icon: '📅', name: 'Tuổi tài khoản' },
                            { icon: '⚡', name: 'Tần suất giao dịch' },
                            { icon: '📱', name: 'Device Fingerprint' },
                            { icon: '🌐', name: 'IP Reputation' },
                            { icon: '💳', name: 'Lịch sử thanh toán' },
                            { icon: '📊', name: 'Độ lệch số tiền' },
                            { icon: '📍', name: 'Vị trí địa lý' },
                            { icon: '🖱️', name: 'Hành vi phiên' },
                        ].map(f => (
                            <div key={f.name} className="p-3 bg-gray-50 rounded text-center">
                                <div className="text-2xl mb-1">{f.icon}</div>
                                <div className="text-sm font-semibold">{f.name}</div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}
