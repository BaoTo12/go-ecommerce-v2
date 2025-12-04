export default function LivePage() {
    const liveStreams = Array.from({ length: 12 }, (_, i) => ({
        id: `live-${i + 1}`,
        seller: `Seller ${i + 1}`,
        title: `Live Shopping Stream ${i + 1}`,
        viewers: Math.floor(Math.random() * 50000),
        isLive: Math.random() > 0.3,
    }))

    return (
        <div className="container mx-auto py-8">
            <div className="mb-6 flex items-center gap-3">
                <h1 className="text-3xl font-bold">🔴 Live Streaming</h1>
                <span className="rounded-full bg-red-500 px-3 py-1 text-sm font-semibold text-white">
                    LIVE
                </span>
            </div>

            <p className="mb-8 text-muted-foreground">
                Watch live streams, shop directly from sellers, and get exclusive deals!
            </p>

            {/* Live Stream Grid */}
            <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
                {liveStreams.map((stream) => (
                    <div
                        key={stream.id}
                        className="group cursor-pointer overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-xl"
                    >
                        {/* Stream Thumbnail */}
                        <div className="relative aspect-video bg-gradient-to-br from-purple-500 to-pink-500">
                            <div className="flex h-full items-center justify-center text-4xl">
                                📹
                            </div>

                            {stream.isLive && (
                                <div className="absolute left-3 top-3 flex items-center gap-2 rounded-full bg-red-500 px-3 py-1">
                                    <span className="h-2 w-2 animate-pulse rounded-full bg-white"></span>
                                    <span className="text-sm font-bold text-white">LIVE</span>
                                </div>
                            )}

                            <div className="absolute bottom-3 right-3 rounded-full bg-black/60 px-3 py-1 text-sm text-white">
                                👁️ {stream.viewers.toLocaleString()}
                            </div>
                        </div>

                        {/* Stream Info */}
                        <div className="p-4">
                            <h3 className="mb-2 font-semibold">{stream.title}</h3>
                            <p className="text-sm text-muted-foreground">{stream.seller}</p>
                        </div>
                    </div>
                ))}
            </div>

            {/* Info Banner */}
            <div className="mt-12 rounded-lg border bg-gradient-to-r from-purple-50 to-pink-50 p-8 dark:from-purple-950/20 dark:to-pink-950/20">
                <h2 className="mb-4 text-2xl font-bold">Start Your Own Live Stream!</h2>
                <p className="mb-4 text-muted-foreground">
                    Sell directly to thousands of customers through live video. Showcase products,
                    answer questions in real-time, and trigger flash sales during your stream.
                </p>
                <button className="rounded-lg bg-gradient-to-r from-purple-500 to-pink-500 px-6 py-3 font-semibold text-white hover:opacity-90">
                    🎥 Start Streaming
                </button>
            </div>
        </div>
    )
}
