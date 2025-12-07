'use client';

import React, { useState, useEffect, useRef } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useRouter } from 'next/navigation';
import { searchService, SearchSuggestion } from '@/services/searchService';

interface SearchAutocompleteProps {
    onSearch?: (query: string) => void;
    className?: string;
}

export default function SearchAutocomplete({ onSearch, className = '' }: SearchAutocompleteProps) {
    const router = useRouter();
    const [query, setQuery] = useState('');
    const [suggestions, setSuggestions] = useState<SearchSuggestion[]>([]);
    const [isOpen, setIsOpen] = useState(false);
    const [selectedIndex, setSelectedIndex] = useState(-1);
    const inputRef = useRef<HTMLInputElement>(null);
    const containerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const fetchSuggestions = async () => {
            const results = await searchService.getSuggestions(query);
            setSuggestions(results);
        };

        const debounce = setTimeout(fetchSuggestions, 200);
        return () => clearTimeout(debounce);
    }, [query]);

    useEffect(() => {
        const handleClickOutside = (e: MouseEvent) => {
            if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
                setIsOpen(false);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleSearch = (searchQuery: string) => {
        if (!searchQuery.trim()) return;
        searchService.addToHistory(searchQuery);
        setIsOpen(false);
        if (onSearch) {
            onSearch(searchQuery);
        } else {
            router.push(`/products?search=${encodeURIComponent(searchQuery)}`);
        }
    };

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'ArrowDown') {
            e.preventDefault();
            setSelectedIndex(prev => Math.min(prev + 1, suggestions.length - 1));
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            setSelectedIndex(prev => Math.max(prev - 1, -1));
        } else if (e.key === 'Enter') {
            e.preventDefault();
            if (selectedIndex >= 0 && suggestions[selectedIndex]) {
                const suggestion = suggestions[selectedIndex];
                if (suggestion.type === 'product' && suggestion.productId) {
                    router.push(`/products/${suggestion.productId}`);
                } else {
                    handleSearch(suggestion.text);
                }
            } else {
                handleSearch(query);
            }
        } else if (e.key === 'Escape') {
            setIsOpen(false);
        }
    };

    const handleRemoveHistory = (e: React.MouseEvent, text: string) => {
        e.stopPropagation();
        searchService.removeFromHistory(text);
        searchService.getSuggestions(query).then(setSuggestions);
    };

    const formatPrice = (price: number) => new Intl.NumberFormat('vi-VN').format(price);

    return (
        <div ref={containerRef} className={`relative ${className}`}>
            <div className="flex bg-white rounded-sm overflow-hidden">
                <input
                    ref={inputRef}
                    type="text"
                    value={query}
                    onChange={(e) => {
                        setQuery(e.target.value);
                        setSelectedIndex(-1);
                    }}
                    onFocus={() => setIsOpen(true)}
                    onKeyDown={handleKeyDown}
                    placeholder="Tìm kiếm sản phẩm, thương hiệu..."
                    className="flex-1 px-4 py-[10px] text-[14px] text-gray-700 dark:bg-gray-800 dark:text-white outline-none"
                />
                <button
                    onClick={() => handleSearch(query)}
                    className="px-5 bg-[#fb5533] hover:opacity-90 flex items-center justify-center transition-opacity"
                >
                    <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                </button>
            </div>

            {/* Suggestions Dropdown */}
            {isOpen && suggestions.length > 0 && (
                <div className="absolute top-full left-0 right-0 bg-white dark:bg-gray-800 mt-1 rounded-sm shadow-lg z-50 overflow-hidden animate-fade-in">
                    {suggestions.map((suggestion, index) => (
                        <div key={`${suggestion.type}-${suggestion.text}-${index}`}>
                            {suggestion.type === 'product' ? (
                                <Link
                                    href={`/products/${suggestion.productId}`}
                                    className={`flex items-center gap-3 px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors ${index === selectedIndex ? 'bg-gray-50 dark:bg-gray-700' : ''
                                        }`}
                                    onClick={() => setIsOpen(false)}
                                >
                                    {suggestion.image && (
                                        <div className="relative w-10 h-10 bg-gray-100 rounded overflow-hidden flex-shrink-0">
                                            <Image
                                                src={suggestion.image}
                                                alt=""
                                                fill
                                                className="object-cover"
                                                unoptimized
                                            />
                                        </div>
                                    )}
                                    <div className="flex-1 min-w-0">
                                        <div className="text-sm truncate dark:text-white">{suggestion.text}</div>
                                        {suggestion.price && (
                                            <div className="text-xs text-[#ee4d2d]">₫{formatPrice(suggestion.price)}</div>
                                        )}
                                    </div>
                                </Link>
                            ) : (
                                <button
                                    onClick={() => handleSearch(suggestion.text)}
                                    className={`w-full flex items-center gap-3 px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors text-left ${index === selectedIndex ? 'bg-gray-50 dark:bg-gray-700' : ''
                                        }`}
                                >
                                    <span className="text-gray-400">
                                        {suggestion.type === 'history' && '🕐'}
                                        {suggestion.type === 'trending' && '🔥'}
                                        {suggestion.type === 'category' && '📁'}
                                    </span>
                                    <span className="flex-1 text-sm dark:text-white">{suggestion.text}</span>
                                    {suggestion.type === 'history' && (
                                        <button
                                            onClick={(e) => handleRemoveHistory(e, suggestion.text)}
                                            className="text-gray-400 hover:text-gray-600 text-xs"
                                        >
                                            ✕
                                        </button>
                                    )}
                                    {suggestion.type === 'trending' && (
                                        <span className="text-xs text-orange-500">Hot</span>
                                    )}
                                </button>
                            )}
                        </div>
                    ))}

                    {query && (
                        <div className="border-t px-4 py-2">
                            <button
                                onClick={() => handleSearch(query)}
                                className="text-sm text-[#ee4d2d] hover:underline"
                            >
                                Tìm kiếm &quot;{query}&quot; →
                            </button>
                        </div>
                    )}
                </div>
            )}
        </div>
    );
}
