// Chat Service - Real-time chat with sellers

export interface ChatMessage {
    id: string;
    senderId: string;
    senderName: string;
    senderAvatar: string;
    content: string;
    timestamp: string;
    type: 'text' | 'image' | 'product';
    productId?: string;
    read: boolean;
}

export interface ChatConversation {
    id: string;
    participantId: string;
    participantName: string;
    participantAvatar: string;
    isShop: boolean;
    lastMessage: string;
    lastMessageTime: string;
    unreadCount: number;
    messages: ChatMessage[];
}

// Mock conversations
const CONVERSATIONS: ChatConversation[] = [
    {
        id: 'conv1',
        participantId: 'shop1',
        participantName: 'Apple Store Official',
        participantAvatar: 'https://ui-avatars.com/api/?name=Apple&background=000&color=fff',
        isShop: true,
        lastMessage: 'Dạ, sản phẩm còn hàng ạ. Anh/chị đặt ngay nhé!',
        lastMessageTime: '2024-12-06T10:30:00',
        unreadCount: 1,
        messages: [
            {
                id: 'm1',
                senderId: 'u1',
                senderName: 'Bạn',
                senderAvatar: 'https://ui-avatars.com/api/?name=You&background=ee4d2d&color=fff',
                content: 'Chào shop, sản phẩm iPhone 15 Pro Max còn hàng không ạ?',
                timestamp: '2024-12-06T10:25:00',
                type: 'text',
                read: true,
            },
            {
                id: 'm2',
                senderId: 'shop1',
                senderName: 'Apple Store Official',
                senderAvatar: 'https://ui-avatars.com/api/?name=Apple&background=000&color=fff',
                content: 'Dạ, sản phẩm còn hàng ạ. Anh/chị đặt ngay nhé!',
                timestamp: '2024-12-06T10:30:00',
                type: 'text',
                read: false,
            },
        ],
    },
    {
        id: 'conv2',
        participantId: 'shop5',
        participantName: 'Dior Beauty Official',
        participantAvatar: 'https://ui-avatars.com/api/?name=Dior&background=9c27b0&color=fff',
        isShop: true,
        lastMessage: 'Cảm ơn bạn đã ủng hộ shop! ❤️',
        lastMessageTime: '2024-12-05T16:20:00',
        unreadCount: 0,
        messages: [
            {
                id: 'm1',
                senderId: 'u1',
                senderName: 'Bạn',
                senderAvatar: 'https://ui-avatars.com/api/?name=You&background=ee4d2d&color=fff',
                content: 'Shop ơi, son này có màu gì đẹp nhất ạ?',
                timestamp: '2024-12-05T16:00:00',
                type: 'text',
                read: true,
            },
            {
                id: 'm2',
                senderId: 'shop5',
                senderName: 'Dior Beauty Official',
                senderAvatar: 'https://ui-avatars.com/api/?name=Dior&background=9c27b0&color=fff',
                content: 'Dạ chào bạn! Màu 001 Pink là best seller của shop, tone hồng tự nhiên phù hợp với mọi tone da ạ!',
                timestamp: '2024-12-05T16:10:00',
                type: 'text',
                read: true,
            },
            {
                id: 'm3',
                senderId: 'u1',
                senderName: 'Bạn',
                senderAvatar: 'https://ui-avatars.com/api/?name=You&background=ee4d2d&color=fff',
                content: 'Ok shop, mình đặt 2 cây màu này nhé!',
                timestamp: '2024-12-05T16:15:00',
                type: 'text',
                read: true,
            },
            {
                id: 'm4',
                senderId: 'shop5',
                senderName: 'Dior Beauty Official',
                senderAvatar: 'https://ui-avatars.com/api/?name=Dior&background=9c27b0&color=fff',
                content: 'Cảm ơn bạn đã ủng hộ shop! ❤️',
                timestamp: '2024-12-05T16:20:00',
                type: 'text',
                read: true,
            },
        ],
    },
    {
        id: 'conv3',
        participantId: 'shop4',
        participantName: 'Nike Official Store',
        participantAvatar: 'https://ui-avatars.com/api/?name=Nike&background=000&color=fff',
        isShop: true,
        lastMessage: 'Đơn hàng của bạn đang được vận chuyển!',
        lastMessageTime: '2024-12-04T09:00:00',
        unreadCount: 0,
        messages: [],
    },
];

let conversations = [...CONVERSATIONS];

export const chatService = {
    // Get all conversations
    getConversations: async (): Promise<ChatConversation[]> => {
        await new Promise(resolve => setTimeout(resolve, 100));
        return conversations.sort((a, b) =>
            new Date(b.lastMessageTime).getTime() - new Date(a.lastMessageTime).getTime()
        );
    },

    // Get conversation by ID
    getConversation: async (conversationId: string): Promise<ChatConversation | null> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        return conversations.find(c => c.id === conversationId) || null;
    },

    // Get or create conversation with a shop
    getOrCreateConversation: async (shopId: string, shopName: string, shopAvatar: string): Promise<ChatConversation> => {
        await new Promise(resolve => setTimeout(resolve, 50));

        let conv = conversations.find(c => c.participantId === shopId);
        if (!conv) {
            conv = {
                id: `conv_${Date.now()}`,
                participantId: shopId,
                participantName: shopName,
                participantAvatar: shopAvatar,
                isShop: true,
                lastMessage: '',
                lastMessageTime: new Date().toISOString(),
                unreadCount: 0,
                messages: [],
            };
            conversations.push(conv);
        }
        return conv;
    },

    // Send message
    sendMessage: async (conversationId: string, content: string, type: 'text' | 'image' | 'product' = 'text'): Promise<ChatMessage> => {
        await new Promise(resolve => setTimeout(resolve, 100));

        const conv = conversations.find(c => c.id === conversationId);
        if (!conv) throw new Error('Conversation not found');

        const message: ChatMessage = {
            id: `msg_${Date.now()}`,
            senderId: 'u1',
            senderName: 'Bạn',
            senderAvatar: 'https://ui-avatars.com/api/?name=You&background=ee4d2d&color=fff',
            content,
            timestamp: new Date().toISOString(),
            type,
            read: true,
        };

        conv.messages.push(message);
        conv.lastMessage = content;
        conv.lastMessageTime = message.timestamp;

        // Simulate shop reply after 1-2 seconds
        setTimeout(() => {
            const replies = [
                'Dạ, cảm ơn bạn! Shop sẽ hỗ trợ ngay ạ.',
                'Dạ vâng, bạn có cần hỗ trợ gì thêm không ạ?',
                'Shop ghi nhận yêu cầu của bạn rồi ạ!',
                'Dạ, bạn chờ shop một chút nhé!',
            ];
            const replyMessage: ChatMessage = {
                id: `msg_${Date.now()}`,
                senderId: conv.participantId,
                senderName: conv.participantName,
                senderAvatar: conv.participantAvatar,
                content: replies[Math.floor(Math.random() * replies.length)],
                timestamp: new Date().toISOString(),
                type: 'text',
                read: false,
            };
            conv.messages.push(replyMessage);
            conv.lastMessage = replyMessage.content;
            conv.lastMessageTime = replyMessage.timestamp;
            conv.unreadCount++;
        }, 1000 + Math.random() * 1000);

        return message;
    },

    // Mark conversation as read
    markAsRead: async (conversationId: string): Promise<void> => {
        await new Promise(resolve => setTimeout(resolve, 50));
        const conv = conversations.find(c => c.id === conversationId);
        if (conv) {
            conv.unreadCount = 0;
            conv.messages.forEach(m => m.read = true);
        }
    },

    // Get total unread count
    getUnreadCount: (): number => {
        return conversations.reduce((sum, c) => sum + c.unreadCount, 0);
    },
};

export default chatService;
