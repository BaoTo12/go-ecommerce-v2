export default function ProductsPage() {
    const categories = [
        'All', 'Electronics', 'Fashion', 'Home & Living',
        'Beauty', 'Sports', 'Books', 'Toys', 'Food'
    ]

    const products = Array.from({ length: 24 }, (_, i) => ({
        id: `prod-${i + 1}`,
        name: `Product ${i + 1}`,
        price: Math.floor(Math.random() * 500) + 10,
        rating: (Math.random() * 2 + 3).toFixed(1),
        sold: Math.floor(Math.random() * 10000),
        discount: Math.random() > 0.5 ? Math.floor(Math.random() * 50) + 10 : 0,
    }))

    return (
        <div className="container mx-auto py-8">
            <h1 className="mb-6 text-3xl font-bold">Browse Products</h1>

            {/* Category Filter */}
            <div className="mb-8 flex gap-2 overflow-x-auto pb-2">
                {categories.map((cat) => (
                    <button
                        key={cat}
                        className="whitespace-nowrap rounded-full border px-4 py-2 text-sm hover:bg-accent"
                    >
                        {cat}
                    </button>
                ))}
            </div>

            {/* Product Grid */}
            <div className="grid grid-cols-2 gap-4 md:grid-cols-4 lg:grid-cols-6">
                {products.map((product) => (
                    <div
                        key={product.id}
                        className="group cursor-pointer rounded-lg border bg-card transition-shadow hover:shadow-lg"
                    >
                        {/* Product Image */}
                        <div className="relative aspect-square overflow-hidden rounded-t-lg bg-muted">
                            <div className="flex h-full items-center justify-center text-muted-foreground">
                                📦 {product.id}
                            </div>
                            {product.discount > 0 && (
                                <div className="absolute right-2 top-2 rounded bg-red-500 px-2 py-1 text-xs font-bold text-white">
                                    -{product.discount}%
                                </div>
                            )}
                        </div>

                        {/* Product Info */}
                        <div className="p-3">
                            <h3 className="mb-1 line-clamp-2 text-sm">{product.name}</h3>
                            <div className="mb-2 flex items-baseline gap-2">
                                <span className="text-lg font-bold text-primary">
                                    ${Math.floor(product.price * (1 - product.discount / 100))}
                                </span>
                                {product.discount > 0 && (
                                    <span className="text-xs text-muted-foreground line-through">
                                        ${product.price}
                                    </span>
                                )}
                            </div>
                            <div className="flex items-center gap-2 text-xs text-muted-foreground">
                                <span>⭐ {product.rating}</span>
                                <span>|</span>
                                <span>{product.sold.toLocaleString()} sold</span>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    )
}
