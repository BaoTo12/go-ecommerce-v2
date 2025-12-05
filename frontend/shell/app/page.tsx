import React from 'react';
import Link from 'next/link';

export default function HomePage() {
    const categories = [
        { icon: '📱', name: 'Điện Thoại' },
        { icon: '💻', name: 'Laptop' },
        { icon: '👗', name: 'Thời Trang' },
        { icon: '💄', name: 'Làm Đẹp' },
        { icon: '🏠', name: 'Nhà Cửa' },
        { icon: '🎮', name: 'Gaming' },
        { icon: '👟', name: 'Giày Dép' },
        { icon: '⌚', name: 'Đồng Hồ' },
        { icon: '📚', name: 'Sách' },
        { icon: '🧸', name: 'Đồ Chơi' },
    ];

    const flashSaleProducts = [
        { id: 1, name: 'iPhone 15 Pro', price: 29990000, originalPrice: 34990000, discount: 14, sold: 87, image: '📱' },
        { id: 2, name: 'AirPods Pro 2', price: 4990000, originalPrice: 6990000, discount: 29, sold: 156, image: '🎧' },
        { id: 3, name: 'MacBook Air', price: 24990000, originalPrice: 27990000, discount: 11, sold: 45, image: '💻' },
    ];

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    return (
        <div className="min-h-screen bg-[#F5F5F5]">
            {/* Hero Banner */}
            <section className="bg-gradient-to-r from-[#EE4D2D] to-[#FF7337] py-8">
                <div className="container mx-auto px-4">
                    <div className="grid md:grid-cols-2 gap-6 items-center">
                        <div className="text-white">
                            <h1 className="text-4xl md:text-5xl font-bold mb-4">
                                Mua Sắm Online
                                <br />
                                <span className="text-yellow-300">Giá Siêu Rẻ!</span>
                            </h1>
                            <p className="text-lg opacity-90 mb-6">
                                Flash Sale mỗi ngày - Freeship mọi đơn - Đổi trả miễn phí
                            </p>
                            <div className="flex gap-3">
                                <Link
                                    href="/deals/flash-sale"
                                    className="bg-white text-[#EE4D2D] px-6 py-3 rounded font-bold hover:bg-yellow-100 transition-colors"
                                >
                                    ⚡ Flash Sale
                                </Link>
                                <Link
                                    href="/live"
                                    className="bg-white/20 text-white px-6 py-3 rounded font-bold hover:bg-white/30 transition-colors border border-white/50"
                                >
                                    🔴 Xem Live
                                </Link>
                            </div>
                        </div>
                        <div className="hidden md:flex justify-center">
                            <div className="text-9xl animate-bounce">🛒</div>
                        </div>
                    </div>
                </div>
            </section>

            {/* Categories */}
            <section className="container mx-auto px-4 py-6">
                <div className="bg-white rounded">
                    <div className="p-4 border-b">
                        <h2 className="font-bold text-[#EE4D2D] uppercase">Danh mục</h2>
                    </div>
                    <div className="grid grid-cols-5 md:grid-cols-10 gap-2 p-4">
                        {categories.map(cat => (
                            <Link
                                key={cat.name}
                                href="/products"
                                className="flex flex-col items-center p-2 hover:bg-gray-50 rounded transition-colors text-center group"
                            >
                                <span className="text-3xl mb-2 group-hover:scale-110 transition-transform">{cat.icon}</span>
                                <span className="text-xs text-gray-700">{cat.name}</span>
                            </Link>
                        ))}
                    </div>
                </div>
            </section>

            {/* Flash Sale Preview */}
            <section className="container mx-auto px-4 py-4">
                <div className="bg-white rounded">
                    <div className="p-4 border-b flex items-center justify-between">
                        <div className="flex items-center gap-3">
                            <span className="text-[#EE4D2D] text-xl font-bold uppercase">⚡ Flash Sale</span>
                            <div className="flex gap-1">
                                <span className="bg-[#333] text-white px-2 py-1 rounded text-sm font-mono">02</span>
                                <span className="text-[#333]">:</span>
                                <span className="bg-[#333] text-white px-2 py-1 rounded text-sm font-mono">45</span>
                                <span className="text-[#333]">:</span>
                                <span className="bg-[#333] text-white px-2 py-1 rounded text-sm font-mono">30</span>
                            </div>
                        </div>
                        <Link href="/deals/flash-sale" className="text-[#EE4D2D] text-sm hover:underline">
                            Xem tất cả &gt;
                        </Link>
                    </div>
                    <div className="grid grid-cols-3 md:grid-cols-6 gap-2 p-4">
                        {flashSaleProducts.map(product => (
                            <div key={product.id} className="group cursor-pointer">
                                <div className="relative aspect-square bg-gray-100 rounded flex items-center justify-center text-5xl">
                                    {product.image}
                                    <span className="absolute top-0 right-0 bg-[#FFEB3B] text-[#EE4D2D] text-xs font-bold px-1">
                                        -{product.discount}%
                                    </span>
                                </div>
                                <div className="mt-2">
                                    <div className="text-[#EE4D2D] font-bold">₫{formatPrice(product.price)}</div>
                                    <div className="h-3 bg-[#FFE0DB] rounded-full overflow-hidden">
                                        <div
                                            className="h-full bg-gradient-to-r from-[#EE4D2D] to-[#FF6633] rounded-full relative"
                                            style={{ width: `${product.sold}%` }}
                                        >
                                            <div className="absolute inset-0 flex items-center justify-center text-[8px] text-white font-bold">
                                                ĐÃ BÁN {product.sold}
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </section>

            {/* Quick Features */}
            <section className="container mx-auto px-4 py-4">
                <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                    <Link
                        href="/rewards"
                        className="bg-gradient-to-br from-yellow-400 to-orange-500 rounded p-4 text-white hover:opacity-90 transition-opacity"
                    >
                        <div className="text-3xl mb-2">🎮</div>
                        <h3 className="font-bold">Shopee Xu</h3>
                        <p className="text-xs opacity-90">Chơi game nhận xu</p>
                    </Link>
                    <Link
                        href="/deals/coupons"
                        className="bg-gradient-to-br from-purple-500 to-pink-500 rounded p-4 text-white hover:opacity-90 transition-opacity"
                    >
                        <div className="text-3xl mb-2">🎟️</div>
                        <h3 className="font-bold">Mã Giảm Giá</h3>
                        <p className="text-xs opacity-90">Voucher hot</p>
                    </Link>
                    <Link
                        href="/live"
                        className="bg-gradient-to-br from-red-500 to-pink-600 rounded p-4 text-white hover:opacity-90 transition-opacity"
                    >
                        <div className="text-3xl mb-2">🔴</div>
                        <h3 className="font-bold">Shopee Live</h3>
                        <p className="text-xs opacity-90">Xem & mua sắm</p>
                    </Link>
                    <Link
                        href="/deals/flash-sale"
                        className="bg-gradient-to-br from-[#EE4D2D] to-[#FF6633] rounded p-4 text-white hover:opacity-90 transition-opacity"
                    >
                        <div className="text-3xl mb-2">⚡</div>
                        <h3 className="font-bold">Flash Sale</h3>
                        <p className="text-xs opacity-90">Giảm đến 90%</p>
                    </Link>
                </div>
            </section>

            {/* Platform Stats */}
            <section className="container mx-auto px-4 py-6">
                <div className="bg-gradient-to-r from-slate-800 to-slate-900 rounded p-6 text-white">
                    <h2 className="text-center font-bold text-xl mb-6">🚀 Hyperscale Platform</h2>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
                        <div>
                            <div className="text-3xl font-bold text-yellow-400">50M+</div>
                            <div className="text-sm opacity-75">Người dùng hàng ngày</div>
                        </div>
                        <div>
                            <div className="text-3xl font-bold text-green-400">1.2M</div>
                            <div className="text-sm opacity-75">Đơn hàng/ngày</div>
                        </div>
                        <div>
                            <div className="text-3xl font-bold text-blue-400">50ms</div>
                            <div className="text-sm opacity-75">P99 Latency</div>
                        </div>
                        <div>
                            <div className="text-3xl font-bold text-purple-400">99.99%</div>
                            <div className="text-sm opacity-75">Uptime SLA</div>
                        </div>
                    </div>
                </div>
            </section>

            {/* Admin Quick Access */}
            <section className="container mx-auto px-4 py-4">
                <div className="bg-white rounded p-4">
                    <h2 className="font-bold text-[#EE4D2D] uppercase mb-4">🔧 Admin Tools</h2>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <Link
                            href="/admin/analytics"
                            className="flex items-center gap-3 p-4 bg-blue-50 rounded hover:bg-blue-100 transition-colors"
                        >
                            <span className="text-3xl">📊</span>
                            <div>
                                <h3 className="font-bold">Analytics Dashboard</h3>
                                <p className="text-xs text-gray-500">Số liệu thời gian thực</p>
                            </div>
                        </Link>
                        <Link
                            href="/admin/fraud"
                            className="flex items-center gap-3 p-4 bg-slate-100 rounded hover:bg-slate-200 transition-colors"
                        >
                            <span className="text-3xl">🛡️</span>
                            <div>
                                <h3 className="font-bold">Fraud Detection</h3>
                                <p className="text-xs text-gray-500">Giám sát ML 99.7%</p>
                            </div>
                        </Link>
                        <Link
                            href="/admin/pricing"
                            className="flex items-center gap-3 p-4 bg-emerald-50 rounded hover:bg-emerald-100 transition-colors"
                        >
                            <span className="text-3xl">💹</span>
                            <div>
                                <h3 className="font-bold">Dynamic Pricing</h3>
                                <p className="text-xs text-gray-500">Tối ưu giá AI</p>
                            </div>
                        </Link>
                    </div>
                </div>
            </section>
        </div>
    );
}
