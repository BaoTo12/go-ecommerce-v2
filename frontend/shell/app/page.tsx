import Link from 'next/link';

export default function HomePage() {
    return (
        <div className="container mx-auto py-12 px-4">
            {/* Hero Section */}
            <section className="mb-16 text-center">
                <h1 className="mb-4 text-5xl font-bold tracking-tight bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                    Welcome to Titan Commerce
                </h1>
                <p className="mb-8 text-xl text-muted-foreground">
                    Hyperscale e-commerce platform with 50M DAU capacity
                </p>
                <div className="flex justify-center gap-4 flex-wrap">
                    <Link
                        href="/products"
                        className="rounded-lg bg-gradient-to-r from-blue-600 to-purple-600 px-6 py-3 font-semibold text-white hover:from-blue-700 hover:to-purple-700 transition-all"
                    >
                        🛍️ Browse Products
                    </Link>
                    <Link
                        href="/live"
                        className="rounded-lg border-2 border-red-500 bg-white px-6 py-3 font-semibold text-red-500 hover:bg-red-50 transition-all"
                    >
                        🔴 Watch Live
                    </Link>
                    <Link
                        href="/deals/flash-sale"
                        className="rounded-lg bg-gradient-to-r from-red-600 to-orange-500 px-6 py-3 font-semibold text-white hover:from-red-700 hover:to-orange-600 transition-all"
                    >
                        ⚡ Flash Sales
                    </Link>
                </div>
            </section>

            {/* Live Stats Banner */}
            <section className="mb-16 rounded-xl bg-gradient-to-r from-slate-900 to-slate-800 p-6 text-white">
                <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
                    <div className="text-center">
                        <div className="text-3xl font-bold text-yellow-400">50M+</div>
                        <div className="text-sm opacity-75">Daily Active Users</div>
                    </div>
                    <div className="text-center">
                        <div className="text-3xl font-bold text-green-400">1.2M</div>
                        <div className="text-sm opacity-75">Orders Today</div>
                    </div>
                    <div className="text-center">
                        <div className="text-3xl font-bold text-blue-400">500</div>
                        <div className="text-sm opacity-75">Cell Shards</div>
                    </div>
                    <div className="text-center">
                        <div className="text-3xl font-bold text-purple-400">99.99%</div>
                        <div className="text-sm opacity-75">Uptime SLA</div>
                    </div>
                </div>
            </section>

            {/* Feature Categories */}
            <section className="mb-16">
                <h2 className="mb-8 text-3xl font-bold text-center">Platform Features</h2>
                <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                    <FeatureCard
                        href="/deals/flash-sale"
                        icon="⚡"
                        title="Flash Sale"
                        description="PoW-protected mega deals handling 1M concurrent shoppers. Anti-bot with cryptographic challenges."
                        gradient="from-red-500 to-orange-500"
                    />
                    <FeatureCard
                        href="/rewards"
                        icon="🎮"
                        title="Gamification"
                        description="Earn Shopee Coins, daily check-ins, lucky draw, and missions. Boost engagement through play."
                        gradient="from-yellow-500 to-orange-500"
                    />
                    <FeatureCard
                        href="/live"
                        icon="🔴"
                        title="Live Shopping"
                        description="Watch sellers stream live, interact in real-time chat, and shop products directly."
                        gradient="from-pink-500 to-red-500"
                    />
                    <FeatureCard
                        href="/deals/coupons"
                        icon="🎟️"
                        title="Coupons"
                        description="Claim exclusive discount coupons and vouchers. Stack savings on every purchase."
                        gradient="from-purple-500 to-pink-500"
                    />
                    <FeatureCard
                        href="/admin/fraud"
                        icon="🛡️"
                        title="Fraud Detection"
                        description="ML-powered real-time fraud scoring with 99.7% detection rate and 8ms latency."
                        gradient="from-slate-600 to-slate-800"
                    />
                    <FeatureCard
                        href="/admin/analytics"
                        icon="📊"
                        title="Analytics"
                        description="Real-time dashboards, conversion funnels, cohort analysis, and ClickHouse-powered insights."
                        gradient="from-blue-500 to-indigo-500"
                    />
                </div>
            </section>

            {/* Quick Actions */}
            <section className="mb-16">
                <h2 className="mb-8 text-3xl font-bold text-center">Quick Actions</h2>
                <div className="grid gap-4 md:grid-cols-4">
                    <Link
                        href="/rewards"
                        className="group rounded-xl border bg-white p-6 text-center hover:shadow-lg hover:border-yellow-400 transition-all"
                    >
                        <div className="text-4xl mb-2 group-hover:scale-110 transition-transform">📅</div>
                        <h3 className="font-semibold">Daily Check-in</h3>
                        <p className="text-sm text-muted-foreground">Earn coins every day</p>
                    </Link>
                    <Link
                        href="/rewards"
                        className="group rounded-xl border bg-white p-6 text-center hover:shadow-lg hover:border-green-400 transition-all"
                    >
                        <div className="text-4xl mb-2 group-hover:scale-110 transition-transform">🎡</div>
                        <h3 className="font-semibold">Lucky Draw</h3>
                        <p className="text-sm text-muted-foreground">Spin to win prizes</p>
                    </Link>
                    <Link
                        href="/deals/flash-sale"
                        className="group rounded-xl border bg-white p-6 text-center hover:shadow-lg hover:border-red-400 transition-all"
                    >
                        <div className="text-4xl mb-2 group-hover:scale-110 transition-transform">⚡</div>
                        <h3 className="font-semibold">Flash Deals</h3>
                        <p className="text-sm text-muted-foreground">Up to 90% off</p>
                    </Link>
                    <Link
                        href="/deals/coupons"
                        className="group rounded-xl border bg-white p-6 text-center hover:shadow-lg hover:border-purple-400 transition-all"
                    >
                        <div className="text-4xl mb-2 group-hover:scale-110 transition-transform">🎟️</div>
                        <h3 className="font-semibold">Claim Coupons</h3>
                        <p className="text-sm text-muted-foreground">Extra savings</p>
                    </Link>
                </div>
            </section>

            {/* Categories */}
            <section className="mb-16">
                <h2 className="mb-8 text-3xl font-bold">Featured Categories</h2>
                <div className="grid grid-cols-2 gap-4 md:grid-cols-4 lg:grid-cols-8">
                    {['📱 Electronics', '👗 Fashion', '🏠 Home', '💄 Beauty', '🏃 Sports', '📚 Books', '🧸 Toys', '🍜 Food'].map((category) => {
                        const [icon, name] = category.split(' ');
                        return (
                            <div
                                key={name}
                                className="rounded-lg border bg-white p-4 text-center hover:shadow-lg transition-shadow cursor-pointer"
                            >
                                <div className="text-3xl mb-2">{icon}</div>
                                <h3 className="font-semibold text-sm">{name}</h3>
                            </div>
                        );
                    })}
                </div>
            </section>

            {/* Admin Tools */}
            <section className="mb-16">
                <h2 className="mb-8 text-3xl font-bold">Admin Tools</h2>
                <div className="grid gap-6 md:grid-cols-3">
                    <Link
                        href="/admin/analytics"
                        className="rounded-xl border bg-gradient-to-br from-blue-50 to-indigo-50 p-6 hover:shadow-lg transition-all"
                    >
                        <div className="text-4xl mb-4">📊</div>
                        <h3 className="text-xl font-bold">Analytics Dashboard</h3>
                        <p className="text-sm text-muted-foreground mt-2">
                            Real-time metrics, conversion funnels, and business intelligence
                        </p>
                    </Link>
                    <Link
                        href="/admin/fraud"
                        className="rounded-xl border bg-gradient-to-br from-slate-50 to-gray-100 p-6 hover:shadow-lg transition-all"
                    >
                        <div className="text-4xl mb-4">🛡️</div>
                        <h3 className="text-xl font-bold">Fraud Detection</h3>
                        <p className="text-sm text-muted-foreground mt-2">
                            ML-powered transaction monitoring and alert management
                        </p>
                    </Link>
                    <Link
                        href="/admin/pricing"
                        className="rounded-xl border bg-gradient-to-br from-emerald-50 to-teal-50 p-6 hover:shadow-lg transition-all"
                    >
                        <div className="text-4xl mb-4">💹</div>
                        <h3 className="text-xl font-bold">Dynamic Pricing</h3>
                        <p className="text-sm text-muted-foreground mt-2">
                            AI-powered price optimization with competitor tracking
                        </p>
                    </Link>
                </div>
            </section>

            {/* Architecture Banner */}
            <section>
                <div className="rounded-xl bg-gradient-to-r from-indigo-600 to-purple-600 p-8 text-center text-white">
                    <h2 className="mb-2 text-2xl font-bold">🏗️ Powered by Cell-Based Architecture</h2>
                    <p className="opacity-90 mb-4">
                        500 isolated cells serving 50M users with complete fault isolation and 99.99% uptime
                    </p>
                    <div className="flex justify-center gap-8 text-sm flex-wrap">
                        <span>⚡ 50ms P99 Latency</span>
                        <span>🔄 1M TPS</span>
                        <span>📦 Eventual Consistency</span>
                        <span>🌐 Multi-Region</span>
                    </div>
                </div>
            </section>
        </div>
    );
}

function FeatureCard({
    href,
    icon,
    title,
    description,
    gradient,
}: {
    href: string;
    icon: string;
    title: string;
    description: string;
    gradient: string;
}) {
    return (
        <Link
            href={href}
            className="group rounded-xl border bg-white overflow-hidden hover:shadow-xl transition-all"
        >
            <div className={`h-2 bg-gradient-to-r ${gradient}`} />
            <div className="p-6">
                <div className="text-4xl mb-3 group-hover:scale-110 transition-transform inline-block">
                    {icon}
                </div>
                <h3 className="text-lg font-bold mb-2">{title}</h3>
                <p className="text-sm text-muted-foreground">{description}</p>
            </div>
        </Link>
    );
}
