# Chat Service

Real-time messaging between buyers and sellers with WebSocket support.

## Purpose
Provides instant messaging capabilities for buyer-seller communication, customer support, and group chats.

## Technology Stack
- **Database**: ScyllaDB (for time-series message storage and high write throughput)
- **Real-time**: WebSocket for bi-directional streaming
- **API**: gRPC + WebSocket

## Key Features
- ✅ One-on-one buyer-seller chat
- ✅ Real-time message delivery via WebSocket
- ✅ Message read receipts and delivery status
- ✅ Typing indicators
- ✅ Online/offline presence tracking
- ✅ Rich media messages (images, files, products, orders)
- ✅ Message editing and deletion
- ✅ Conversation threading
- ✅ Unread message counts
- ✅ Message history pagination

## Domain Model

### Message
- Text, image, file, product, or order messages
- Delivery and read status tracking
- Edit and delete capabilities
- Metadata for rich content

### Conversation
- One-on-one or group conversations
- Participant management
- Last message tracking
- Unread counts per user
- Mute functionality

### Presence
- Online/offline status
- Typing indicators
- Last seen timestamp

## Quick Start

```bash
export SERVICE_NAME=chat-service
export CELL_ID=cell-001
export SCYLLA_HOSTS=localhost
go run cmd/server/main.go
```

## API Overview

### Commands
- `CreateConversation`: Start new conversation
- `SendMessage`: Send message to conversation
- `EditMessage`: Edit sent message
- `DeleteMessage`: Delete message
- `MarkMessagesAsRead`: Mark messages as read
- `SetTypingIndicator`: Show/hide typing status
- `SetOnlineStatus`: Update user presence

### Queries
- `GetConversation`: Get conversation details
- `GetUserConversations`: List user's conversations
- `GetMessages`: Retrieve message history
- `GetTypingIndicators`: Get who's typing
- `GetOnlineStatus`: Check user presence

### Streaming
- `StreamMessages`: WebSocket stream for real-time updates

## Integration

### Events Published
- `MessageSent`: New message sent
- `MessageRead`: Message marked as read
- `ConversationCreated`: New conversation started
- `UserOnline`: User came online
- `UserOffline`: User went offline

### Events Consumed
- `OrderCreated`: Link order in chat
- `ProductUpdated`: Update product cards in chat
- `UserBlocked`: Prevent messaging

## Database Schema

See `migrations/001_init.cql` for complete ScyllaDB schema.

## WebSocket Protocol

WebSocket endpoint: `ws://<host>/ws/chat`

Connection requires authentication token in query params:
```
ws://localhost:8080/ws/chat?token=<jwt_token>&user_id=<user_id>
```

Real-time events streamed:
- New messages
- Typing indicators
- Presence updates
- Message status changes
