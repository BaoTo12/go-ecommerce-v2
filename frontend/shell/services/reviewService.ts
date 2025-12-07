// Review Service - Product reviews and ratings

import { Product } from './productService';

export interface Review {
    id: string;
    productId: string;
    userId: string;
    userName: string;
    userAvatar: string;
    rating: number;
    content: string;
    images: string[];
    variant?: string;
    likes: number;
    isLiked: boolean;
    sellerReply?: {
        content: string;
        timestamp: string;
    };
    timestamp: string;
}

export interface ReviewStats {
    totalReviews: number;
    averageRating: number;
    ratingCounts: Record<number, number>;
    withPhotos: number;
    withContent: number;
}

// Mock reviews data
const REVIEWS: Record<string, Review[]> = {
    'p1': [
        {
            id: 'r1',
            productId: 'p1',
            userId: 'u2',
            userName: 'Trần Văn B',
            userAvatar: 'https://ui-avatars.com/api/?name=Tran+Van+B&background=random',
            rating: 5,
            content: 'Sản phẩm chính hãng, đóng gói cẩn thận. Giao hàng nhanh chóng. Rất hài lòng!',
            images: ['https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=200'],
            variant: 'Titan Xanh, 256GB',
            likes: 24,
            isLiked: false,
            sellerReply: {
                content: 'Cảm ơn bạn đã tin tưởng và ủng hộ shop ạ! Chúc bạn sử dụng sản phẩm vui vẻ! ❤️',
                timestamp: '2024-12-05T10:00:00',
            },
            timestamp: '2024-12-04T15:30:00',
        },
        {
            id: 'r2',
            productId: 'p1',
            userId: 'u3',
            userName: 'Lê Thị C',
            userAvatar: 'https://ui-avatars.com/api/?name=Le+Thi+C&background=random',
            rating: 4,
            content: 'Máy đẹp, dùng mượt. Chỉ tiếc là giá hơi cao.',
            images: [],
            variant: 'Titan Đen, 512GB',
            likes: 12,
            isLiked: true,
            timestamp: '2024-12-03T09:15:00',
        },
        {
            id: 'r3',
            productId: 'p1',
            userId: 'u4',
            userName: 'Nguyễn Hoàng D',
            userAvatar: 'https://ui-avatars.com/api/?name=Nguyen+Hoang+D&background=random',
            rating: 5,
            content: 'Đã mua nhiều lần ở shop, lần nào cũng hài lòng. Camera chụp ảnh cực đẹp!',
            images: [
                'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=200',
                'https://images.unsplash.com/photo-1592750475338-74b7b21085ab?w=200',
            ],
            variant: 'Titan Trắng, 256GB',
            likes: 45,
            isLiked: false,
            timestamp: '2024-12-01T18:45:00',
        },
    ],
};

export const reviewService = {
    // Get reviews for a product
    getProductReviews: async (productId: string): Promise<Review[]> => {
        await new Promise(resolve => setTimeout(resolve, 100));
        return REVIEWS[productId] || [];
    },

    // Get review stats for a product
    getReviewStats: async (productId: string): Promise<ReviewStats> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        const reviews = REVIEWS[productId] || [];

        const ratingCounts: Record<number, number> = { 1: 0, 2: 0, 3: 0, 4: 0, 5: 0 };
        reviews.forEach(r => {
            ratingCounts[r.rating] = (ratingCounts[r.rating] || 0) + 1;
        });

        return {
            totalReviews: reviews.length,
            averageRating: reviews.length > 0
                ? reviews.reduce((sum, r) => sum + r.rating, 0) / reviews.length
                : 0,
            ratingCounts,
            withPhotos: reviews.filter(r => r.images.length > 0).length,
            withContent: reviews.filter(r => r.content.length > 0).length,
        };
    },

    // Add a review
    addReview: async (review: Omit<Review, 'id' | 'timestamp' | 'likes' | 'isLiked'>): Promise<Review> => {
        await new Promise(resolve => setTimeout(resolve, 200));

        const newReview: Review = {
            ...review,
            id: `r_${Date.now()}`,
            timestamp: new Date().toISOString(),
            likes: 0,
            isLiked: false,
        };

        if (!REVIEWS[review.productId]) {
            REVIEWS[review.productId] = [];
        }
        REVIEWS[review.productId].unshift(newReview);

        return newReview;
    },

    // Like a review
    toggleLike: async (reviewId: string, productId: string): Promise<boolean> => {
        await new Promise(resolve => setTimeout(resolve, 50));

        const reviews = REVIEWS[productId];
        if (!reviews) return false;

        const review = reviews.find(r => r.id === reviewId);
        if (!review) return false;

        review.isLiked = !review.isLiked;
        review.likes += review.isLiked ? 1 : -1;
        return true;
    },

    // Get user's reviews
    getUserReviews: async (userId: string): Promise<Review[]> => {
        await new Promise(resolve => setTimeout(resolve, 100));

        const allReviews: Review[] = [];
        Object.values(REVIEWS).forEach(productReviews => {
            allReviews.push(...productReviews.filter(r => r.userId === userId));
        });
        return allReviews;
    },
};

export default reviewService;
