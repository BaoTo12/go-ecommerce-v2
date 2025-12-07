// Authentication Service - Handles user login, register, and session management

export interface AuthUser {
    id: string;
    email: string;
    name: string;
    phone: string;
    avatar: string;
    isVerified: boolean;
}

export interface LoginCredentials {
    email: string;
    password: string;
}

export interface RegisterData {
    name: string;
    email: string;
    phone: string;
    password: string;
}

export interface AuthState {
    user: AuthUser | null;
    isAuthenticated: boolean;
    token: string | null;
}

// Mock users database
const USERS: Record<string, { user: AuthUser; password: string }> = {
    'nguyenvana@email.com': {
        user: {
            id: 'u1',
            email: 'nguyenvana@email.com',
            name: 'Nguyễn Văn A',
            phone: '0901234567',
            avatar: 'https://ui-avatars.com/api/?name=Nguyen+Van+A&background=ee4d2d&color=fff&size=200',
            isVerified: true,
        },
        password: '123456',
    },
};

// Auth state management
let currentAuth: AuthState = {
    user: null,
    isAuthenticated: false,
    token: null,
};

// Initialize from localStorage
if (typeof window !== 'undefined') {
    const saved = localStorage.getItem('auth');
    if (saved) {
        try {
            currentAuth = JSON.parse(saved);
        } catch (e) {
            console.error('Failed to parse auth state:', e);
        }
    }
}

const saveAuthState = () => {
    if (typeof window !== 'undefined') {
        localStorage.setItem('auth', JSON.stringify(currentAuth));
    }
};

export const authService = {
    // Get current auth state
    getAuthState: (): AuthState => {
        return currentAuth;
    },

    // Check if user is authenticated
    isAuthenticated: (): boolean => {
        return currentAuth.isAuthenticated && !!currentAuth.token;
    },

    // Get current user
    getCurrentUser: (): AuthUser | null => {
        return currentAuth.user;
    },

    // Login
    login: async (credentials: LoginCredentials): Promise<{ success: boolean; error?: string }> => {
        await new Promise(resolve => setTimeout(resolve, 500)); // Simulate API call

        const userRecord = USERS[credentials.email];
        if (!userRecord) {
            return { success: false, error: 'Email không tồn tại' };
        }
        if (userRecord.password !== credentials.password) {
            return { success: false, error: 'Mật khẩu không đúng' };
        }

        currentAuth = {
            user: userRecord.user,
            isAuthenticated: true,
            token: `token_${Date.now()}_${Math.random().toString(36)}`,
        };
        saveAuthState();

        return { success: true };
    },

    // Register
    register: async (data: RegisterData): Promise<{ success: boolean; error?: string }> => {
        await new Promise(resolve => setTimeout(resolve, 500));

        if (USERS[data.email]) {
            return { success: false, error: 'Email đã được sử dụng' };
        }

        const newUser: AuthUser = {
            id: `u${Date.now()}`,
            email: data.email,
            name: data.name,
            phone: data.phone,
            avatar: `https://ui-avatars.com/api/?name=${encodeURIComponent(data.name)}&background=ee4d2d&color=fff&size=200`,
            isVerified: false,
        };

        USERS[data.email] = { user: newUser, password: data.password };

        currentAuth = {
            user: newUser,
            isAuthenticated: true,
            token: `token_${Date.now()}_${Math.random().toString(36)}`,
        };
        saveAuthState();

        return { success: true };
    },

    // Send OTP for verification
    sendOTP: async (phone: string): Promise<{ success: boolean; otp?: string }> => {
        await new Promise(resolve => setTimeout(resolve, 300));
        // In production, this would send a real SMS
        const otp = Math.floor(100000 + Math.random() * 900000).toString();
        console.log(`[Mock] OTP for ${phone}: ${otp}`);
        return { success: true, otp }; // Return OTP for demo purposes
    },

    // Verify OTP
    verifyOTP: async (phone: string, otp: string): Promise<{ success: boolean }> => {
        await new Promise(resolve => setTimeout(resolve, 300));
        // In production, this would verify against stored OTP
        if (otp.length === 6) {
            if (currentAuth.user) {
                currentAuth.user.isVerified = true;
                saveAuthState();
            }
            return { success: true };
        }
        return { success: false };
    },

    // Logout
    logout: async (): Promise<void> => {
        await new Promise(resolve => setTimeout(resolve, 100));
        currentAuth = {
            user: null,
            isAuthenticated: false,
            token: null,
        };
        saveAuthState();
    },

    // Reset password
    resetPassword: async (email: string): Promise<{ success: boolean }> => {
        await new Promise(resolve => setTimeout(resolve, 300));
        if (USERS[email]) {
            console.log(`[Mock] Password reset email sent to ${email}`);
            return { success: true };
        }
        return { success: false };
    },

    // Update profile
    updateProfile: async (data: Partial<AuthUser>): Promise<{ success: boolean }> => {
        await new Promise(resolve => setTimeout(resolve, 200));
        if (currentAuth.user) {
            currentAuth.user = { ...currentAuth.user, ...data };
            saveAuthState();
            return { success: true };
        }
        return { success: false };
    },
};

export default authService;
