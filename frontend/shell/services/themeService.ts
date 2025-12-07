// Theme Service - Dark mode support with system preference detection

export type Theme = 'light' | 'dark' | 'system';

interface ThemeState {
    theme: Theme;
    isDark: boolean;
}

// Get from localStorage or system preference
const getStoredTheme = (): ThemeState => {
    if (typeof window === 'undefined') {
        return { theme: 'system', isDark: false };
    }

    const saved = localStorage.getItem('theme') as Theme | null;
    const theme = saved || 'system';

    let isDark = false;
    if (theme === 'dark') {
        isDark = true;
    } else if (theme === 'system') {
        isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    }

    return { theme, isDark };
};

let state = getStoredTheme();

const applyTheme = (isDark: boolean) => {
    if (typeof document === 'undefined') return;

    if (isDark) {
        document.documentElement.classList.add('dark');
    } else {
        document.documentElement.classList.remove('dark');
    }
};

// Listen to system preference changes
if (typeof window !== 'undefined') {
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
        if (state.theme === 'system') {
            state.isDark = e.matches;
            applyTheme(state.isDark);
        }
    });

    // Apply initial theme
    applyTheme(state.isDark);
}

export const themeService = {
    // Get current theme
    getTheme: (): ThemeState => state,

    // Set theme
    setTheme: (theme: Theme): void => {
        state.theme = theme;

        if (theme === 'dark') {
            state.isDark = true;
        } else if (theme === 'light') {
            state.isDark = false;
        } else {
            state.isDark = typeof window !== 'undefined'
                ? window.matchMedia('(prefers-color-scheme: dark)').matches
                : false;
        }

        if (typeof window !== 'undefined') {
            localStorage.setItem('theme', theme);
        }

        applyTheme(state.isDark);
    },

    // Toggle dark mode
    toggleDark: (): boolean => {
        const newTheme = state.isDark ? 'light' : 'dark';
        themeService.setTheme(newTheme);
        return state.isDark;
    },

    // Check if dark mode
    isDarkMode: (): boolean => state.isDark,
};

export default themeService;
