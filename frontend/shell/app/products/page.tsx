'use client';

import React, { useState, useEffect, Suspense } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useSearchParams } from 'next/navigation';
import { productService, Product } from '@/services/productService';

function ProductsContent() {
    const searchParams = useSearchParams();
    const categoryParam = searchParams.get('category');
    const searchParam = searchParams.get('search');

    const [products, setProducts] = useState<Product[]>([]);
    const [selectedCategory, setSelectedCategory] = useState(categoryParam || 'Tất cả');
    const [sortBy, setSortBy] = useState('popular');
    const [searchQuery, setSearchQuery] = useState(searchParam || '');
    const [priceRange, setPriceRange] = useState({ min: '', max: '' });
    const [notification, setNotification] = useState<string | null>(null);
    const [isLoaded, setIsLoaded] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [addingToCart, setAddingToCart] = useState<string | null>(null);
    const [totalProducts, setTotalProducts] = useState(0);

    const categories = ['Tất cả', 'Điện thoại', 'Laptop', 'Thời trang', 'Làm đẹp', 'Nhà cửa', 'Giày dép', 'Túi ví', 'Đồng hồ'];

    // Load products from service
    useEffect(() => {
        const loadProducts = async () => {
            setIsLoading(true);
            try {
                const { products: data, total } = await productService.getProducts({
                    category: selectedCategory,
                    search: searchQuery,
                    sort: sortBy,
                });
                setProducts(data);
                setTotalProducts(total);
                setIsLoaded(true);
            } catch (error) {
                console.error('Failed to load products:', error);
            } finally {
                setIsLoading(false);
            }
        };
        loadProducts();
    }, [selectedCategory, sortBy, searchQuery]);

    // Update category when URL param changes
    useEffect(() => {
        if (categoryParam) {
            setSelectedCategory(categoryParam);
        }
    }, [categoryParam]);

    // Update search when URL param changes
    useEffect(() => {
        if (searchParam) {
            setSearchQuery(searchParam);
        }
    }, [searchParam]);

    const addToCart = (productId: string, productName: string, e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setAddingToCart(productId);

        setTimeout(() => {
            setAddingToCart(null);
            setNotification(`✓ Đã thêm "${productName.substring(0, 25)}..." vào giỏ hàng`);
            setTimeout(() => setNotification(null), 2500);
        }, 500);
    };

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    return (
        <div className={`min-h-screen bg-[#f5f5f5] ${isLoaded ? 'animate-fade-in' : 'opacity-0'}`}>
            {/* Toast */}
            {notification && (
                <div className="toast toast-success">{notification}</div>
            )}

            <div className="container mx-auto px-4 py-4">
                <div className="flex gap-4">
                    {/* Sidebar Filters */}
                    <aside className="w-[190px] flex-shrink-0 hidden lg:block animate-fade-in-left">
                        <div className="bg-white rounded-sm shadow-sm p-4 sticky top-[140px]">
                            <h3 className="font-bold text-sm mb-3 flex items-center gap-2">
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h7" />
                                </svg>
                                Bộ Lọc Tìm Kiếm
                            </h3>

                            {/* Categories */}
                            <div className="border-b pb-4 mb-4">
                                <h4 className="text-sm font-medium mb-2">Theo Danh Mục</h4>
                                <div className="space-y-1">
                                    {categories.map((cat, index) => (
                                        <button
                                            key={cat}
                                            onClick={() => setSelectedCategory(cat)}
                                            className={`block w-full text-left text-sm py-1.5 px-2 rounded transition-all duration-200 ${selectedCategory === cat
                                                    ? 'text-[#ee4d2d] bg-[#fef6f5] font-medium'
                                                    : 'text-gray-600 hover:text-[#ee4d2d] hover:bg-gray-50'
                                                }`}
                                            style={{ animationDelay: `${index * 50}ms` }}
                                        >
                                            {cat}
                                            {selectedCategory === cat && (
                                                <span className="ml-1 text-xs">({totalProducts})</span>
                                            )}
                                        </button>
                                    ))}
                                </div>
                            </div>

                            {/* Price Range */}
                            <div className="border-b pb-4 mb-4">
                                <h4 className="text-sm font-medium mb-2">Khoảng Giá</h4>
                                <div className="flex items-center gap-2">
                                    <input
                                        type="text"
                                        placeholder="₫ TỪ"
                                        value={priceRange.min}
                                        onChange={(e) => setPriceRange({ ...priceRange, min: e.target.value })}
                                        className="w-full border px-2 py-1.5 text-sm outline-none focus:border-[#ee4d2d] rounded-sm transition-all"
                                    />
                                    <span className="text-gray-400">-</span>
                                    <input
                                        type="text"
                                        placeholder="₫ ĐẾN"
                                        value={priceRange.max}
                                        onChange={(e) => setPriceRange({ ...priceRange, max: e.target.value })}
                                        className="w-full border px-2 py-1.5 text-sm outline-none focus:border-[#ee4d2d] rounded-sm transition-all"
                                    />
                                </div>
                                <button className="w-full mt-2 py-1.5 bg-[#ee4d2d] text-white text-sm hover:opacity-90 transition-all hover-shrink ripple">
                                    ÁP DỤNG
                                </button>
                            </div>

                            {/* Rating Filter */}
                            <div>
                                <h4 className="text-sm font-medium mb-2">Đánh Giá</h4>
                                {[5, 4, 3].map(stars => (
                                    <button key={stars} className="flex items-center gap-1 py-1.5 text-sm text-gray-600 hover:text-[#ee4d2d] transition-colors w-full">
                                        {[...Array(5)].map((_, i) => (
                                            <span key={i} className={`transition-transform hover:scale-110 ${i < stars ? 'star-rating' : 'text-gray-300'}`}>★</span>
                                        ))}
                                        <span className="ml-1">trở lên</span>
                                    </button>
                                ))}
                            </div>
                        </div>
                    </aside>

                    {/* Main Content */}
                    <main className="flex-1">
                        {/* Breadcrumb */}
                        <div className="bg-white rounded-sm shadow-sm p-3 mb-3 animate-fade-in-down">
                            <div className="flex items-center gap-2 text-sm">
                                <Link href="/" className="text-gray-500 hover:text-[#ee4d2d]">Trang chủ</Link>
                                <span className="text-gray-400">›</span>
                                <span className="text-gray-700">{selectedCategory === 'Tất cả' ? 'Tất cả sản phẩm' : selectedCategory}</span>
                                <span className="ml-auto text-gray-500">{totalProducts} sản phẩm</span>
                            </div>
                        </div>

                        {/* Search Bar */}
                        <div className="bg-white rounded-sm shadow-sm p-3 mb-3 animate-fade-in-down">
                            <input
                                type="text"
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                placeholder="Tìm sản phẩm trong danh mục này..."
                                className="w-full border px-3 py-2 text-sm outline-none focus:border-[#ee4d2d] rounded-sm transition-all"
                            />
                        </div>

                        {/* Sort Bar */}
                        <div className="bg-[#ededed] rounded-sm p-3 mb-3 flex items-center gap-2 animate-fade-in">
                            <span className="text-gray-500 text-sm">Sắp xếp theo</span>
                            {[
                                { value: 'popular', label: 'Phổ Biến' },
                                { value: 'newest', label: 'Mới Nhất' },
                                { value: 'best-seller', label: 'Bán Chạy' },
                            ].map(opt => (
                                <button
                                    key={opt.value}
                                    onClick={() => setSortBy(opt.value)}
                                    className={`px-4 py-1.5 text-sm rounded-sm transition-all duration-200 hover-shrink ${sortBy === opt.value
                                            ? 'bg-[#ee4d2d] text-white'
                                            : 'bg-white text-gray-700 hover:bg-gray-100'
                                        }`}
                                >
                                    {opt.label}
                                </button>
                            ))}
                            <select
                                value={sortBy.startsWith('price') ? sortBy : ''}
                                onChange={(e) => e.target.value && setSortBy(e.target.value)}
                                className="ml-auto px-3 py-1.5 text-sm border bg-white outline-none cursor-pointer transition-all hover:border-[#ee4d2d]"
                            >
                                <option value="">Giá</option>
                                <option value="price-asc">Giá: Thấp đến Cao</option>
                                <option value="price-desc">Giá: Cao đến Thấp</option>
                            </select>
                        </div>

                        {/* Loading State */}
                        {isLoading ? (
                            <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-[10px]">
                                {[...Array(10)].map((_, i) => (
                                    <div key={i} className="bg-white rounded-sm overflow-hidden animate-pulse">
                                        <div className="aspect-square bg-gray-200" />
                                        <div className="p-3 space-y-2">
                                            <div className="h-4 bg-gray-200 rounded" />
                                            <div className="h-4 bg-gray-200 rounded w-2/3" />
                                            <div className="h-5 bg-gray-200 rounded w-1/2" />
                                        </div>
                                    </div>
                                ))}
                            </div>
                        ) : products.length === 0 ? (
                            <div className="bg-white rounded-sm p-12 text-center animate-fade-in">
                                <div className="text-5xl mb-4 animate-float">🔍</div>
                                <p className="text-gray-500 mb-2">Không tìm thấy sản phẩm phù hợp</p>
                                <p className="text-sm text-gray-400">Thử tìm kiếm với từ khóa khác hoặc chọn danh mục khác</p>
                            </div>
                        ) : (
                            <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-[10px]">
                                {products.map((product, index) => (
                                    <Link
                                        key={product.id}
                                        href={`/products/${product.id}`}
                                        className="product-card group animate-fade-in-up"
                                        style={{ animationDelay: `${index * 50}ms` }}
                                    >
                                        <div className="relative aspect-square bg-gray-100 overflow-hidden">
                                            <Image
                                                src={product.thumbnail}
                                                alt={product.name}
                                                fill
                                                className={`object-cover product-image ${addingToCart === product.id ? 'animate-bounce-in' : ''}`}
                                                unoptimized
                                            />

                                            {/* Discount badge */}
                                            {product.discount > 0 && (
                                                <div className="discount-badge">-{product.discount}%</div>
                                            )}

                                            {/* Top left badges */}
                                            <div className="absolute top-0 left-0 flex flex-col gap-0.5">
                                                {product.isOfficial && (
                                                    <div className="bg-[#ee4d2d] text-white text-[10px] px-1">Mall</div>
                                                )}
                                                {product.isFavorite && (
                                                    <div className="bg-[#ee4d2d] text-white text-[10px] px-1">Yêu thích</div>
                                                )}
                                            </div>

                                            {/* Add to cart */}
                                            <button
                                                onClick={(e) => addToCart(product.id, product.name, e)}
                                                disabled={addingToCart === product.id}
                                                className={`absolute bottom-2 right-2 w-8 h-8 bg-[#ee4d2d] text-white rounded-sm flex items-center justify-center 
                                    opacity-0 group-hover:opacity-100 transition-all duration-300 text-sm 
                                    hover:bg-[#d73211] hover:scale-110 transform translate-y-2 group-hover:translate-y-0
                                    ${addingToCart === product.id ? 'animate-spin' : ''}`}
                                            >
                                                {addingToCart === product.id ? '⏳' : '+'}
                                            </button>
                                        </div>

                                        <div className="p-2">
                                            <h3 className="text-xs line-clamp-2 h-8 mb-1 group-hover:text-[#ee4d2d] transition-colors">{product.name}</h3>

                                            {/* Free ship badge */}
                                            {product.freeShip && (
                                                <div className="inline-flex items-center text-[10px] text-[#00bfa5] border border-[#00bfa5] px-1 mb-1">
                                                    <svg className="w-3 h-3 mr-0.5" fill="currentColor" viewBox="0 0 24 24">
                                                        <path d="M20 8h-3V4H3c-1.1 0-2 .9-2 2v11h2c0 1.66 1.34 3 3 3s3-1.34 3-3h6c0 1.66 1.34 3 3 3s3-1.34 3-3h2v-5l-3-4z" />
                                                    </svg>
                                                    Miễn phí
                                                </div>
                                            )}

                                            <div className="flex items-end justify-between">
                                                <div>
                                                    <span className="price-current text-sm font-medium">₫{formatPrice(product.price)}</span>
                                                    {product.originalPrice > product.price && (
                                                        <span className="price-original block">₫{formatPrice(product.originalPrice)}</span>
                                                    )}
                                                </div>
                                            </div>

                                            <div className="flex items-center justify-between mt-1 text-[11px] text-gray-500">
                                                <span className="flex items-center gap-0.5">
                                                    <span className="star-rating">★</span> {product.rating}
                                                </span>
                                                <span>Đã bán {product.soldDisplay}</span>
                                            </div>
                                            <div className="text-[11px] text-gray-400 mt-0.5">{product.location}</div>
                                        </div>
                                    </Link>
                                ))}
                            </div>
                        )}
                    </main>
                </div>
            </div>
        </div>
    );
}

export default function ProductsPage() {
    return (
        <Suspense fallback={
            <div className="min-h-screen bg-[#f5f5f5] flex items-center justify-center">
                <div className="loading-spinner" />
            </div>
        }>
            <ProductsContent />
        </Suspense>
    );
}
