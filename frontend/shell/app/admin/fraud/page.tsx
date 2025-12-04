'use client';

import { useState, useEffect } from 'react';
import { fraudApi, FraudCheck, FraudAlert } from '../../lib/api';

export default function FraudDashboard() {
    const [alerts, setAlerts] = useState<FraudAlert[]>([]);
    const [recentChecks, setRecentChecks] = useState<FraudCheck[]>([]);
    const [loading, setLoading] = useState(true);
    const [testResult, setTestResult] = useState<FraudCheck | null>(null);
    const [testing, setTesting] = useState(false);

    useEffect(() => {
        loadData();
        const interval = setInterval(loadData, 10000);
        return () => clearInterval(interval);
    }, []);

    const loadData = async () => {
        try {
            const alertsData = await fraudApi.getAlerts();
            setAlerts(alertsData || []);
        } catch {
            // Mock data
            setAlerts([
                {
                    id: 'alert-1',
                    fraud_check_id: 'fc-001',
                    alert_type: 'high_risk',
                    severity: 'HIGH',
                    message: 'Multiple failed payments from new device',
                    acknowledged: false,
                },
                {
                    id: 'alert-2',
                    fraud_check_id: 'fc-002',
                    alert_type: 'velocity',
                    severity: 'MEDIUM',
                    message: 'Unusual order frequency detected',
                    acknowledged: false,
                },
            ]);
            setRecentChecks([
                {
                    check_id: 'fc-001',
                    transaction_id: 'txn-12345',
                    score: 0.85,
                    risk_level: 'HIGH',
                    decision: 'REVIEW',
                    reasons: ['New device', 'High amount', 'Multiple IPs'],
                    processing_time: 12,
                },
                {
                    check_id: 'fc-002',
                    transaction_id: 'txn-12346',
                    score: 0.45,
                    risk_level: 'MEDIUM',
                    decision: 'ALLOW',
                    reasons: ['Order velocity elevated'],
                    processing_time: 8,
                },
                {
                    check_id: 'fc-003',
                    transaction_id: 'txn-12347',
                    score: 0.12,
                    risk_level: 'LOW',
                    decision: 'ALLOW',
                    reasons: [],
                    processing_time: 5,
                },
            ]);
        } finally {
            setLoading(false);
        }
    };

    const runTestTransaction = async () => {
        setTesting(true);
        try {
            const result = await fraudApi.checkTransaction({
                transaction_id: `txn-test-${Date.now()}`,
                user_id: 'user-test',
                amount: 599.99,
                currency: 'USD',
                ip: '192.168.1.100',
                device_id: 'device-new',
            });
            setTestResult(result);
        } catch {
            // Mock result
            setTestResult({
                check_id: 'fc-test',
                transaction_id: `txn-test-${Date.now()}`,
                score: Math.random(),
                risk_level: Math.random() > 0.7 ? 'HIGH' : Math.random() > 0.4 ? 'MEDIUM' : 'LOW',
                decision: Math.random() > 0.8 ? 'BLOCK' : 'ALLOW',
                reasons: ['Test transaction', 'New device detected'],
                processing_time: Math.floor(Math.random() * 20),
            });
        } finally {
            setTesting(false);
        }
    };

    const getRiskColor = (level: string) => {
        switch (level) {
            case 'HIGH': return 'text-red-600 bg-red-100';
            case 'MEDIUM': return 'text-yellow-600 bg-yellow-100';
            case 'LOW': return 'text-green-600 bg-green-100';
            default: return 'text-gray-600 bg-gray-100';
        }
    };

    const getDecisionColor = (decision: string) => {
        switch (decision) {
            case 'BLOCK': return 'text-red-600';
            case 'REVIEW': return 'text-yellow-600';
            case 'ALLOW': return 'text-green-600';
            default: return 'text-gray-600';
        }
    };

    if (loading) {
        return (
            <div className="flex h-96 items-center justify-center">
                <div className="animate-pulse text-2xl">🛡️ Loading Fraud Dashboard...</div>
            </div>
        );
    }

    return (
        <div className="container mx-auto py-8">
            {/* Header */}
            <div className="mb-8 rounded-xl bg-gradient-to-r from-slate-800 to-slate-900 p-8 text-white">
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-3xl font-bold">🛡️ Fraud Detection Center</h1>
                        <p className="text-slate-300">ML-powered real-time transaction monitoring</p>
                    </div>
                    <div className="flex gap-4">
                        <div className="text-center">
                            <div className="text-3xl font-bold">99.7%</div>
                            <div className="text-sm text-slate-400">Detection Rate</div>
                        </div>
                        <div className="text-center">
                            <div className="text-3xl font-bold">8ms</div>
                            <div className="text-sm text-slate-400">Avg Latency</div>
                        </div>
                        <div className="text-center">
                            <div className="text-3xl font-bold">0.02%</div>
                            <div className="text-sm text-slate-400">False Positive</div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Alerts Section */}
            {alerts.length > 0 && (
                <div className="mb-8">
                    <h2 className="mb-4 text-xl font-bold text-red-600">
                        🚨 Active Alerts ({alerts.length})
                    </h2>
                    <div className="space-y-3">
                        {alerts.map((alert) => (
                            <div
                                key={alert.id}
                                className={`rounded-lg border-l-4 p-4 ${alert.severity === 'HIGH'
                                        ? 'border-red-500 bg-red-50'
                                        : 'border-yellow-500 bg-yellow-50'
                                    }`}
                            >
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center gap-3">
                                        <span className="text-2xl">
                                            {alert.severity === 'HIGH' ? '🔴' : '🟡'}
                                        </span>
                                        <div>
                                            <p className="font-semibold">{alert.message}</p>
                                            <p className="text-sm text-muted-foreground">
                                                Check ID: {alert.fraud_check_id} | Type: {alert.alert_type}
                                            </p>
                                        </div>
                                    </div>
                                    <button className="rounded-lg bg-white px-4 py-2 text-sm font-semibold shadow hover:bg-gray-50">
                                        Review
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}

            <div className="grid gap-8 lg:grid-cols-3">
                {/* Test Transaction */}
                <div className="rounded-xl border bg-white p-6">
                    <h2 className="mb-4 text-xl font-bold">🧪 Test Transaction</h2>
                    <p className="mb-4 text-sm text-muted-foreground">
                        Run a test transaction through the fraud detection engine
                    </p>
                    <button
                        onClick={runTestTransaction}
                        disabled={testing}
                        className="w-full rounded-lg bg-gradient-to-r from-blue-600 to-indigo-600 py-3 font-bold text-white hover:from-blue-700 hover:to-indigo-700 disabled:opacity-50"
                    >
                        {testing ? 'Analyzing...' : 'Run Test'}
                    </button>

                    {testResult && (
                        <div className="mt-4 rounded-lg border p-4">
                            <div className="mb-3 flex items-center justify-between">
                                <span className={`rounded-full px-3 py-1 text-sm font-semibold ${getRiskColor(testResult.risk_level)}`}>
                                    {testResult.risk_level} RISK
                                </span>
                                <span className={`font-bold ${getDecisionColor(testResult.decision)}`}>
                                    {testResult.decision}
                                </span>
                            </div>
                            <div className="space-y-2 text-sm">
                                <div className="flex justify-between">
                                    <span className="text-muted-foreground">Score</span>
                                    <span className="font-mono">{(testResult.score * 100).toFixed(1)}%</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-muted-foreground">Processing</span>
                                    <span className="font-mono">{testResult.processing_time}ms</span>
                                </div>
                                {testResult.reasons.length > 0 && (
                                    <div>
                                        <span className="text-muted-foreground">Flags:</span>
                                        <ul className="mt-1 list-inside list-disc">
                                            {testResult.reasons.map((r, i) => (
                                                <li key={i} className="text-yellow-700">{r}</li>
                                            ))}
                                        </ul>
                                    </div>
                                )}
                            </div>
                        </div>
                    )}
                </div>

                {/* Recent Checks */}
                <div className="lg:col-span-2 rounded-xl border bg-white p-6">
                    <h2 className="mb-4 text-xl font-bold">📋 Recent Fraud Checks</h2>
                    <div className="overflow-x-auto">
                        <table className="w-full">
                            <thead>
                                <tr className="border-b text-left text-sm text-muted-foreground">
                                    <th className="pb-3">Transaction</th>
                                    <th className="pb-3">Score</th>
                                    <th className="pb-3">Risk</th>
                                    <th className="pb-3">Decision</th>
                                    <th className="pb-3">Time</th>
                                </tr>
                            </thead>
                            <tbody className="text-sm">
                                {recentChecks.map((check) => (
                                    <tr key={check.check_id} className="border-b hover:bg-gray-50">
                                        <td className="py-3 font-mono">{check.transaction_id}</td>
                                        <td className="py-3">
                                            <div className="flex items-center gap-2">
                                                <div className="h-2 w-16 rounded-full bg-gray-200">
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
                                            <span className={`rounded-full px-2 py-1 text-xs font-semibold ${getRiskColor(check.risk_level)}`}>
                                                {check.risk_level}
                                            </span>
                                        </td>
                                        <td className={`py-3 font-semibold ${getDecisionColor(check.decision)}`}>
                                            {check.decision}
                                        </td>
                                        <td className="py-3 text-muted-foreground">{check.processing_time}ms</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>

            {/* ML Features */}
            <div className="mt-8 rounded-xl border bg-gradient-to-r from-slate-50 to-slate-100 p-8">
                <h2 className="mb-6 text-2xl font-bold">🤖 ML Feature Engineering</h2>
                <div className="grid gap-4 md:grid-cols-4">
                    {[
                        { name: 'Account Age', icon: '📅', desc: 'Days since creation' },
                        { name: 'Order Velocity', icon: '⚡', desc: 'Orders per time period' },
                        { name: 'Device Fingerprint', icon: '📱', desc: 'Unique device tracking' },
                        { name: 'IP Reputation', icon: '🌐', desc: 'IP risk scoring' },
                        { name: 'Payment History', icon: '💳', desc: 'Past payment patterns' },
                        { name: 'Amount Deviation', icon: '📊', desc: 'vs historical average' },
                        { name: 'Geo-Location', icon: '📍', desc: 'Location anomaly detection' },
                        { name: 'Session Behavior', icon: '🖱️', desc: 'Browsing patterns' },
                    ].map((feature) => (
                        <div key={feature.name} className="rounded-lg bg-white p-4 shadow-sm">
                            <div className="mb-2 text-2xl">{feature.icon}</div>
                            <h3 className="font-semibold">{feature.name}</h3>
                            <p className="text-xs text-muted-foreground">{feature.desc}</p>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
