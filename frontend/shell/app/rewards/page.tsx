export default function RewardsPage() {
    const userCoins = 1250

    return (
        <div className="container mx-auto py-8">
            <div className="mb-8 text-center">
                <h1 className="mb-4 text-4xl font-bold">🎮 Rewards & Games</h1>
                <div className="inline-flex items-center gap-3 rounded-full bg-gradient-to-r from-yellow-400 to-orange-500 px-8 py-4 text-white">
                    <span className="text-2xl">💰</span>
                    <div>
                        <div className="text-sm opacity-90">Your Shopee Coins</div>
                        <div className="text-3xl font-bold">{userCoins.toLocaleString()}</div>
                    </div>
                </div>
            </div>

            {/* Games Grid */}
            <div className="mb-12 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
                {/* Shake Game */}
                <div className="group cursor-pointer overflow-hidden rounded-lg border bg-gradient-to-br from-purple-500 to-pink-500 p-6 text-white transition-transform hover:scale-105">
                    <div className="mb-4 text-6xl">📱</div>
                    <h3 className="mb-2 text-2xl font-bold">Shake-Shake</h3>
                    <p className="mb-4 opacity-90">Shake your phone to win coins!</p>
                    <button className="w-full rounded-lg bg-white/20 px-4 py-2 font-semibold backdrop-blur hover:bg-white/30">
                        Play Now
                    </button>
                </div>

                {/* Lucky Draw */}
                <div className="group cursor-pointer overflow-hidden rounded-lg border bg-gradient-to-br from-blue-500 to-cyan-500 p-6 text-white transition-transform hover:scale-105">
                    <div className="mb-4 text-6xl">🎰</div>
                    <h3 className="mb-2 text-2xl font-bold">Lucky Draw</h3>
                    <p className="mb-4 opacity-90">Spin the wheel for prizes!</p>
                    <button className="w-full rounded-lg bg-white/20 px-4 py-2 font-semibold backdrop-blur hover:bg-white/30">
                        Spin (50 coins)
                    </button>
                </div>

                {/* Daily Check-in */}
                <div className="group cursor-pointer overflow-hidden rounded-lg border bg-gradient-to-br from-green-500 to-emerald-500 p-6 text-white transition-transform hover:scale-105">
                    <div className="mb-4 text-6xl">📅</div>
                    <h3 className="mb-2 text-2xl font-bold">Daily Check-in</h3>
                    <p className="mb-4 opacity-90">Login daily for streak bonuses!</p>
                    <button className="w-full rounded-lg bg-white/20 px-4 py-2 font-semibold backdrop-blur hover:bg-white/30">
                        Check In
                    </button>
                </div>
            </div>

            {/* Missions */}
            <div className="mb-12">
                <h2 className="mb-6 text-2xl font-bold">🎯 Active Missions</h2>
                <div className="space-y-4">
                    {[
                        { title: 'Make your first purchase', reward: 100, progress: 0, goal: 1 },
                        { title: 'Buy 3 items this week', reward: 500, progress: 1, goal: 3 },
                        { title: 'Follow 5 sellers', reward: 200, progress: 3, goal: 5 },
                        { title: 'Share a product', reward: 50, progress: 0, goal: 1 },
                    ].map((mission, i) => (
                        <div key={i} className="rounded-lg border bg-card p-4">
                            <div className="mb-2 flex items-center justify-between">
                                <h3 className="font-semibold">{mission.title}</h3>
                                <span className="font-bold text-yellow-500">+{mission.reward} 💰</span>
                            </div>
                            <div className="mb-2">
                                <div className="h-2 overflow-hidden rounded-full bg-muted">
                                    <div
                                        className="h-full bg-gradient-to-r from-yellow-400 to-orange-500"
                                        style={{ width: `${(mission.progress / mission.goal) * 100}%` }}
                                    />
                                </div>
                            </div>
                            <div className="text-sm text-muted-foreground">
                                {mission.progress}/{mission.goal} completed
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {/* Coin Economy Info */}
            <div className="rounded-lg border bg-gradient-to-r from-yellow-50 to-orange-50 p-6 dark:from-yellow-950/20 dark:to-orange-950/20">
                <h2 className="mb-4 text-xl font-bold">💡 How to Use Your Coins</h2>
                <div className="grid gap-4 md:grid-cols-3">
                    <div>
                        <div className="mb-2 text-3xl">💸</div>
                        <h3 className="mb-1 font-semibold">Get Discounts</h3>
                        <p className="text-sm text-muted-foreground">100 coins = $1 off</p>
                    </div>
                    <div>
                        <div className="mb-2 text-3xl">🎁</div>
                        <h3 className="mb-1 font-semibold">Redeem Prizes</h3>
                        <p className="text-sm text-muted-foreground">Use coins for lucky draws</p>
                    </div>
                    <div>
                        <div className="mb-2 text-3xl">🔓</div>
                        <h3 className="mb-1 font-semibold">Unlock Deals</h3>
                        <p className="text-sm text-muted-foreground">Access exclusive offers</p>
                    </div>
                </div>
            </div>
        </div>
    )
}
