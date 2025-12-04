export default function HomePage() {
    return (
        <div className="container mx-auto py-12">
            <section className="mb-16 text-center">
                <h1 className="mb-4 text-5xl font-bold tracking-tight">
                    Welcome to Titan Commerce
                </h1>
                <p className="mb-8 text-xl text-muted-foreground">
                    Hyperscale e-commerce platform with 50M DAU capacity
                </p>
                <div className="flex justify-center space-x-4">
                    <a
                        href="/products"
                        className="rounded-lg bg-primary px-6 py-3 font-semibold text-primary-foreground hover:bg-primary/90"
                    >
                        Browse Products
                    </a>
                    <a
                        href="/live"
                        className="rounded-lg border bg-background px-6 py-3 font-semibold hover:bg-accent"
                    >
                        Watch Live Streams
                    </a>
                </div>
            </section>

            <section className="mb-16">
                <h2 className="mb-8 text-3xl font-bold">Featured Categories</h2>
                <div className="grid grid-cols-1 gap-6 md:grid-cols-3 lg:grid-cols-4">
                    {['Electronics', 'Fashion', 'Home & Living', 'Beauty', 'Sports', 'Books', 'Toys', 'Food'].map((category) => (
                        <div
                            key={category}
                            className="rounded-lg border p-6 text-center hover:shadow-lg transition-shadow cursor-pointer"
                        >
                            <h3 className="font-semibold text-lg">{category}</h3>
                        </div>
                    ))}
                </div>
            </section>

            <section className="mb-16">
                <h2 className="mb-8 text-3xl font-bold">Modern Features</h2>
                <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
                    <FeatureCard
                        title="🔴 Live Shopping"
                        description="Watch live streams and shop directly from sellers"
                    />
                    <FeatureCard
                        title="⚡ Flash Sales"
                        description="Mega deals with 1M concurrent shoppers"
                    />
                    <FeatureCard
                        title="🎮 Gamification"
                        description="Earn coins, play games, get rewards"
                    />
                    <FeatureCard
                        title="🤝 Social Commerce"
                        description="Follow sellers, share products, discover together"
                    />
                </div>
            </section>

            <section>
                <div className="rounded-lg bg-primary/10 p-8 text-center">
                    <h2 className="mb-2 text-2xl font-bold">Powered by Cell-Based Architecture</h2>
                    <p className="text-muted-foreground">
                        500 isolated cells serving 50M users with 99.99% uptime
                    </p>
                </div>
            </section>
        </div>
    )
}

function FeatureCard({ title, description }: { title: string; description: string }) {
    return (
        <div className="rounded-lg border p-6 hover:shadow-lg transition-shadow">
            <h3 className="mb-2 text-lg font-semibold">{title}</h3>
            <p className="text-sm text-muted-foreground">{description}</p>
        </div>
    )
}
