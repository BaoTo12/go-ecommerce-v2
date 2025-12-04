# Videocall Service

1-on-1 customer support video calls using WebRTC.

## Purpose
Provides video calling infrastructure for customer support and seller consultations.

## Technology Stack
- **WebRTC**: Peer-to-peer video/audio
- **Signaling**: Redis Pub/Sub
- **STUN/TURN**: For NAT traversal
- **Database**: PostgreSQL (call metadata)
- **API**: gRPC + WebSocket

## Key Features
- ✅ 1-on-1 video calls
- ✅ WebRTC peer-to-peer connections
- ✅ STUN/TURN server configuration
- ✅ Call signaling via Redis
- ✅ Call status tracking
- ✅ Call quality metrics
- ✅ Call recording (optional)
- ✅ Call history
- ✅ Multiple call types (support, sales, consult)

## Architecture

### WebRTC Flow
1. Caller initiates call
2. Service creates room and generates ICE servers
3. WebRTC signaling via Redis Pub/Sub
4. Peer connection established
5. Call quality monitored
6. Call ended and metrics saved

### Signaling Events
- `offer`: SDP offer from caller
- `answer`: SDP answer from callee
- `ice-candidate`: ICE candidates exchange
- `hangup`: End call signal

## Quick Start

```bash
export SERVICE_NAME=videocall-service
export CELL_ID=cell-001
export REDIS_HOST=localhost
export TURN_SERVER=turn:localhost:3478
go run cmd/server/main.go
```

## API Overview

### Commands
- `InitiateCall`: Start new call
- `RingCall`: Mark call as ringing
- `AnswerCall`: Accept call
- `RejectCall`: Reject call
- `EndCall`: Hang up
- `SendSignal`: Send WebRTC signaling data
- `UpdateCallQuality`: Report quality metrics

### Queries
- `GetCall`: Get call details
- `GetUserCalls`: Get call history
- `SubscribeToSignals`: Listen for signaling events

## Integration

### Events Published
- `CallInitiated`: New call started
- `CallAnswered`: Call accepted
- `CallEnded`: Call finished
- `CallMissed`: Call not answered
- `CallRejected`: Call declined

### Events Consumed
- `UserOffline`: End active calls
- `SupportTicketCreated`: Offer video support

## WebRTC Configuration

### ICE Servers
```json
{
  "iceServers": [
    {"urls": "stun:stun.l.google.com:19302"},
    {"urls": "turn:turn.example.com:3478", "username": "user", "credential": "pass"}
  ]
}
```

### Quality Metrics Tracked
- Average latency (ms)
- Packet loss (%)
- Average bitrate (kbps)
- Video resolution
- Connection stability
