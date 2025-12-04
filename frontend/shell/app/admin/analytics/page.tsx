'use client';

import { useState, useEffect } from 'react';
import { analyticsApi, DashboardMetrics, SalesReport, ConversionFunnel } from '../../lib/api';

export default function AnalyticsDashboard() {
    const [metrics, setMetrics] = useState<DashboardMetrics | null>(null);
    const [salesReport, setSalesReport] = useState<SalesReport | null>(null);
    const [funnel, setFunnel] = useState<ConversionFunnel | null>(null);
    const [loading, setLoading] = useState(true);
    const [period, setPeriod] = useState('daily');

    useEffect(() => {
        loadData();
        const interval = setInterval(loadData, 30000);
        return () => clearInterval(interval);
    }, [period]);

    const loadData = async () => {
        try {
            const [metricsData, salesData, funnelData] = await Promise.all([
                analyticsApi.getDashboard(),
                analyticsApi.getSalesReport(period),
                analyticsApi.getConversionFunnel('purchase'),
            ]);
            setMetrics(metricsData);
            setSalesReport(salesData);
            setFunnel(funnelData);
        } catch {
            // Mock data
            setMetrics({
                active_users: 12453,
                orders_in_progress: 847,
                current_revenue: 156789.50,
                today_orders: 3421,
                today_revenue: 234567.89,
                today_page_views: 89432,
                orders_change: 12.5,
                revenue_change: 8.3,
            });
            setSalesReport({
                period: period,
                total_orders: 15234,
                total_revenue: 1234567.89,
                avg_order_value: 81.02,
            });
            setFunnel({
                name: 'Purchase Funnel',
                steps: [
                    { name: 'Page View', users: 100000, conversion: 100 },
                    { name: 'Product View', users: 45000, conversion: 45 },
                    { name: 'Add to Cart', users: 18000, conversion: 40 },
                    { name: 'Checkout', users: 9000, conversion: 50 },
                    { name: 'Purchase', users: 7200, conversion: 80 },
                ],
                overall_rate: 7.2,
            });
        } finally {
            setLoading(false);
        }
    };

    if (loading) {
        return (
            <div className="flex h-96 items-center justify-center">
                <div className="animate-pulse text-2xl">📊 Loading Analytics...</div>
            </div>
        );
    }

    return (
        <div className="container mx-auto py-8">
            {/* Header */}
            <div className="mb-8 flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold">📊 Analytics Dashboard</h1>
                    <p className="text-muted-foreground">Real-time business intelligence</p>
                </div>
                <div className="flex gap-2">
                    {['daily', 'weekly', 'monthly'].map((p) => (
                        <button
                            key={p}
                            onClick={() => setPeriod(p)}
                            className={`rounded-lg px-4 py-2 capitalize transition-all ${period === p
                                    ? 'bg-blue-600 text-white'
                                    : 'bg-gray-100 hover:bg-gray-200'
                                }`}
                        >
                            {p}
                        </button>
                    ))}
                </div>
            </div>

            {/* Real-time Stats */}
            <div className="mb-8 grid gap-4 md:grid-cols-4">
                <StatCard
                    icon="👥"
                    label="Active Users"
                    value={metrics?.active_users.toLocaleString() || '0'}
                    color="bg-blue-500"
                    pulse
                />
                <StatCard
                    icon="🛒"
                    label="Orders in Progress"
                    value={metrics?.orders_in_progress.toLocaleString() || '0'}
                    color="bg-green-500"
                    pulse
                />
                <StatCard
                    icon="💰"
                    label="Live Revenue"
                    value={`$${metrics?.current_revenue.toLocaleString() || '0'}`}
                    color="bg-yellow-500"
                    pulse
                />
                <StatCard
                    icon="📈"
                    label="Today's Orders"
                    value={metrics?.today_orders.toLocaleString() || '0'}
                    change={metrics?.orders_change}
                    color="bg-purple-500"
                />
            </div>

            {/* Main Charts */}
            <div className="grid gap-8 lg:grid-cols-2">
                {/* Revenue Overview */}
                <div className="rounded-xl border bg-white p-6">
                    <h2 className="mb-4 text-xl font-bold">💵 Revenue Overview</h2>
                    <div className="space-y-4">
                        <div className="flex items-center justify-between">
                            <span>Today's Revenue</span>
                            <span className="text-2xl font-bold text-green-600">
                                ${metrics?.today_revenue.toLocaleString()}
                            </span>
                        </div>
                        <div className="flex items-center justify-between">
                            <span>Average Order Value</span>
                            <span className="text-xl font-semibold">
                                ${salesReport?.avg_order_value.toFixed(2)}
                            </span>
                        </div>
                        <div className="flex items-center justify-between">
                            <span>Total Orders ({period})</span>
                            <span className="text-xl font-semibold">
                                {salesReport?.total_orders.toLocaleString()}
                            </span>
                        </div>

                        {/* Mock Chart */}
                        <div className="mt-4 h-48 rounded-lg bg-gradient-to-t from-blue-500/20 to-transparent flex items-end justify-around p-4">
                            {[65, 82, 45, 91, 78, 95, 88].map((h, i) => (
                                <div
                                    key={i}
                                    className="w-8 bg-gradient-to-t from-blue-600 to-blue-400 rounded-t transition-all hover:opacity-80"
                                    style={{ height: `${h}%` }}
                                />
                            ))}
                        </div>
                        <div className="flex justify-around text-xs text-muted-foreground">
                            {['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'].map((d) => (
                                <span key={d}>{d}</span>
                            ))}
                        </div>
                    </div>
                </div>

                {/* Conversion Funnel */}
                <div className="rounded-xl border bg-white p-6">
                    <h2 className="mb-4 text-xl font-bold">🔄 Conversion Funnel</h2>
                    <div className="space-y-4">
                        {funnel?.steps.map((step, index) => {
                            const width = (step.users / (funnel.steps[0]?.users || 1)) * 100;
                            return (
                                <div key={step.name}>
                                    <div className="mb-1 flex justify-between text-sm">
                                        <span className="font-medium">{step.name}</span>
                                        <span className="text-muted-foreground">
                                            {step.users.toLocaleString()} ({step.conversion}%)
                                        </span>
                                    </div>
                                    <div className="h-8 rounded bg-gray-100">
                                        <div
                                            className="h-full rounded bg-gradient-to-r from-indigo-500 to-purple-500 transition-all flex items-center justify-end pr-2"
                                            style={{ width: `${width}%` }}
                                        >
                                            {width > 30 && (
                                                <span className="text-xs text-white font-semibold">
                                                    {step.conversion}%
                                                </span>
                                            )}
                                        </div>
                                    </div>
                                    {index < funnel.steps.length - 1 && (
                                        <div className="ml-4 text-xs text-gray-400">↓</div>
                                    )}
                                </div>
                            );
                        })}
                        <div className="mt-4 rounded-lg bg-green-100 p-4 text-center">
                            <span className="text-sm text-green-700">Overall Conversion Rate</span>
                            <div className="text-3xl font-bold text-green-600">
                                {funnel?.overall_rate}%
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Additional Metrics */}
            <div className="mt-8 grid gap-4 md:grid-cols-3">
                <div className="rounded-xl border bg-white p-6">
                    <h3 className="mb-4 text-lg font-bold">🌍 Traffic by Region</h3>
                    <div className="space-y-3">
                        {[
                            { region: 'Southeast Asia', percent: 45 },
                            { region: 'East Asia', percent: 28 },
                            { region: 'South Asia', percent: 15 },
                            { region: 'Others', percent: 12 },
                        ].map((r) => (
                            <div key={r.region}>
                                <div className="flex justify-between text-sm">
                                    <span>{r.region}</span>
                                    <span>{r.percent}%</span>
                                </div>
                                <div className="h-2 rounded-full bg-gray-200">
                                    <div
                                        className="h-full rounded-full bg-gradient-to-r from-teal-500 to-emerald-500"
                                        style={{ width: `${r.percent}%` }}
                                    />
                                </div>
                            </div>
                        ))}
                    </div>
                </div>

                <div className="rounded-xl border bg-white p-6">
                    <h3 className="mb-4 text-lg font-bold">📱 Device Breakdown</h3>
                    <div className="flex items-center justify-center gap-4">
                        <div className="text-center">
                            <div className="text-4xl">📱</div>
                            <div className="text-2xl font-bold">68%</div>
                            <div className="text-sm text-muted-foreground">Mobile</div>
                        </div>
                        <div className="text-center">
                            <div className="text-4xl">💻</div>
                            <div className="text-2xl font-bold">24%</div>
                            <div className="text-sm text-muted-foreground">Desktop</div>
                        </div>
                        <div className="text-center">
                            <div className="text-4xl">📟</div>
                            <div className="text-2xl font-bold">8%</div>
                            <div className="text-sm text-muted-foreground">Tablet</div>
                        </div>
                    </div>
                </div>

                <div className="rounded-xl border bg-white p-6">
                    <h3 className="mb-4 text-lg font-bold">🏆 Top Categories</h3>
                    <div className="space-y-3">
                        {[
                            { name: 'Electronics', revenue: 45234 },
                            { name: 'Fashion', revenue: 38756 },
                            { name: 'Home', revenue: 28123 },
                            { name: 'Beauty', revenue: 19876 },
                        ].map((c, i) => (
                            <div key={c.name} className="flex items-center gap-3">
                                <span className="text-xl">
                                    {['🥇', '🥈', '🥉', '4️⃣'][i]}
                                </span>
                                <span className="flex-1">{c.name}</span>
                                <span className="font-semibold">${c.revenue.toLocaleString()}</span>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}

function StatCard({
    icon,
    label,
    value,
    change,
    color,
    pulse,
}: {
    icon: string;
    label: string;
    value: string;
    change?: number;
    color: string;
    pulse?: boolean;
}) {
    return (
        <div className="rounded-xl border bg-white p-6 transition-shadow hover:shadow-lg">
            <div className="flex items-center gap-4">
                <div className={`rounded-lg ${color} p-3 text-2xl text-white ${pulse ? 'animate-pulse' : ''}`}>
                    {icon}
                </div>
                <div>
                    <p className="text-sm text-muted-foreground">{label}</p>
                    <p className="text-2xl font-bold">{value}</p>
                    {change !== undefined && (
                        <p className={`text-sm ${change >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                            {change >= 0 ? '↑' : '↓'} {Math.abs(change)}%
                        </p>
                    )}
                </div>
            </div>
        </div>
    );
}
