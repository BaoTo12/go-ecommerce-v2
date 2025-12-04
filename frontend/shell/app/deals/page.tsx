'use client'

import { useState, useEffect } from 'react'

export default function DealsPage() {
    const [timeLeft, setTimeLeft] = useState({ hours: 12, minutes: 34, seconds: 56 })

    useEffect(() => {
        const timer = setInterval(() => {
            setTimeLeft(prev => {
                let { hours, minutes, seconds } = prev
                if (seconds > 0) seconds--
                else if (minutes > 0) { minutes--; seconds = 59 }
                else if (hours > 0) { hours--; minutes = 59; seconds = 59 }
                return { hours, minutes, seconds }
            })
        }, 1000)
        return () => clearInterval(timer)
    }, [])

    const flashSales = Array.from({ length: 8 }, (_, i) => ({
        id: `flash-${i + 1}`,
        name: `Flash Deal ${i + 1}`,
        original: Math.floor(Math.random() * 200) + 50,
        price: Math.floor(Math.random() * 50) + 10,
        stock: Math.floor(Math.random() * 100),
        total: 100,
    }))

    return (
        <div className="container mx-auto py-8">
            <div className="mb-8 text-center">
                <h1 className="mb-4 text-4xl font-bold">⚡ Flash Sales</h1>
                <p className="mb-6 text-xl text-muted-foreground">
                    Lightning deals ending soon!
                </p>

                {/* Countdown Timer */}
                <div className="inline-flex gap-4 rounded-lg bg-gradient-to-r from-orange-500 to-red-500 p-6 text-white">
                    <div className="text-center">
                        <div className="text-4xl font-bold">{String(timeLeft.hours).padStart(2, '0')}</div>
                        <div className="text-sm">Hours</div>
                    </div>
                    <div className="text-4xl font-bold">:</div>
                    <div className="text-center">
                        <div className="text-4xl font-bold">{String(timeLeft.minutes).padStart(2, '0')}</div>
                        <div className="text-sm">Minutes</div>
                    </div>
                    <div className="text-4xl font-bold">:</div>
                    <div className="text-center">
                        <div className="text-4xl font-bold">{String(timeLeft.seconds).padStart(2, '0')}</div>
                        <div className="text-sm">Seconds</div>
                    </div>
                </div>
            </div>

            {/* Flash Sale Grid */}
            <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
                {flashSales.map((deal) => {
                    const progress = ((deal.total - deal.stock) / deal.total) * 100
                    const discount = Math.round(((deal.original - deal.price) / deal.original) * 100)

                    return (
                        <div
                            key={deal.id}
                            className="group cursor-pointer rounded-lg border bg-card p-4 transition-all hover:shadow-xl hover:scale-105"
                        >
                            {/* Product Image */}
                            <div className="relative mb-3 aspect-square overflow-hidden rounded-lg bg-gradient-to-br from-yellow-400 to-orange-500">
                                <div className="flex h-full items-center justify-center text-4xl">
                                    ⚡
                                </div>
                                <div className="absolute right-2 top-2 rounded-full bg-red-500 px-2 py-1 text-xs font-bold text-white">
                                    -{discount}%
                                </div>
                            </div>

                            {/* Product Info */}
                            <h3 className="mb-2 line-clamp-2 text-sm font-semibold">{deal.name}</h3>

                            <div className="mb-2 flex items-baseline gap-2">
                                <span className="text-xl font-bold text-red-500">${deal.price}</span>
                                <span className="text-sm text-muted-foreground line-through">
                                    ${deal.original}
                                </span>
                            </div>

                            {/* Stock Progress */}
                            <div className="mb-2">
                                <div className="mb-1 flex justify-between text-xs">
                                    <span>Stock</span>
                                    <span className="font-semibold">{deal.stock}/{deal.total}</span>
                                </div>
                                <div className="h-2 overflow-hidden rounded-full bg-muted">
                                    <div
                                        className="h-full bg-gradient-to-r from-orange-500 to-red-500 transition-all"
                                        style={{ width: `${progress}%` }}
                                    />
                                </div>
                            </div>

                            <button className="w-full rounded-lg bg-gradient-to-r from-orange-500 to-red-500 py-2 text-sm font-semibold text-white hover:opacity-90">
                                Buy Now ⚡
                            </button>
                        </div>
                    )
                })}
            </div>

            {/* Flash Sale Info */}
            <div className="mt-12 rounded-lg border bg-muted/50 p-6">
                <h2 className="mb-3 text-xl font-bold">🔥 How Flash Sales Work</h2>
                <ul className="space-y-2 text-sm text-muted-foreground">
                    <li>✨ New flash sales every hour</li>
                    <li>⚡ Limited stock - first come, first served</li>
                    <li>🎯 Up to 90% discount on selected items</li>
                    <li>🚀 Handles 1M concurrent users (The "11.11 Problem")</li>
                </ul>
            </div>
        </div>
    )
}
