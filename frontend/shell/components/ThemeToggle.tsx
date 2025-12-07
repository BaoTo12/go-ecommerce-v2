'use client';

import React, { useState, useEffect } from 'react';
import { themeService, Theme } from '@/services/themeService';

export default function ThemeToggle() {
    const [isDark, setIsDark] = useState(false);
    const [mounted, setMounted] = useState(false);

    useEffect(() => {
        setMounted(true);
        setIsDark(themeService.isDarkMode());
    }, []);

    const toggleTheme = () => {
        const newIsDark = themeService.toggleDark();
        setIsDark(newIsDark);
    };

    if (!mounted) {
        return null; // Avoid hydration mismatch
    }

    return (
        <button
            onClick={toggleTheme}
            className="p-2 rounded-full hover:bg-white/10 transition-colors"
            title={isDark ? 'Chế độ sáng' : 'Chế độ tối'}
        >
            {isDark ? (
                <svg className="w-5 h-5 text-yellow-400" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M12 3a1 1 0 011 1v1a1 1 0 11-2 0V4a1 1 0 011-1zm0 15a1 1 0 011 1v1a1 1 0 11-2 0v-1a1 1 0 011-1zm9-6a1 1 0 100-2h-1a1 1 0 100 2h1zM5 12a1 1 0 100-2H4a1 1 0 100 2h1zm14.364-5.364a1 1 0 00-1.414-1.414l-.707.707a1 1 0 101.414 1.414l.707-.707zM6.757 18.243a1 1 0 00-1.414-1.414l-.707.707a1 1 0 001.414 1.414l.707-.707zm12.02 0l.707.707a1 1 0 001.414-1.414l-.707-.707a1 1 0 00-1.414 1.414zM6.757 5.757a1 1 0 00-1.414 0L4.636 6.464a1 1 0 101.414 1.414l.707-.707a1 1 0 000-1.414zM12 8a4 4 0 100 8 4 4 0 000-8z" />
                </svg>
            ) : (
                <svg className="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                </svg>
            )}
        </button>
    );
}
