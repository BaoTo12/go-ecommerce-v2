# Livestream Service

Live video shopping infrastructure with real-time streaming.

## Purpose
Powers TikTok-style live shopping streams where sellers showcase products in real-time to buyers.

## Technology Stack
- **Streaming**: RTMP ingestion, HLS/DASH playback
- **Storage**: Redis (viewer tracking, real-time data), S3 (recordings)
- **Database**: PostgreSQL (stream metadata)
- **API**: gRPC + WebSocket

## Key Features
- ✅ Live video streaming (RTMP → HLS/DASH)
- ✅ Real-time viewer count tracking
- ✅ Live chat comments
- ✅ Featured product showcase
- ✅ In-stream purchases
- ✅ Like and share functionality
- ✅ Stream scheduling
- ✅ Automatic recording to S3
- ✅ Stream analytics (views, watch time, revenue)
- ✅ Peak viewer tracking
- ✅ Viewer retention metrics

## Architecture

### Stream Flow
1. Seller streams via RTMP to media server
2. Media server transcodes to HLS/DASH
3. CDN distributes to viewers
4. Service tracks viewers and interactions
5. Stream recorded to S3

### Real-time Features
- Viewer join/leave tracking
- Live comment stream
- Product feature notifications
- Sale notifications

## Quick Start

```bash
export SERVICE_NAME=livestream-service
export CELL_ID=cell-001
export REDIS_HOST=localhost
export S3_BUCKET=livestreams
go run cmd/server/main.go
```

## API Overview

### Commands
- `CreateStream`: Schedule new livestream
- `StartStream`: Go live
- `EndStream`: End stream
- `JoinStream`: Viewer joins
- `LeaveStream`: Viewer leaves
- `PostComment`: Post live comment
- `LikeStream`: Like the stream
- `AddFeaturedProduct`: Showcase product
- `RecordSale`: Track sale

### Queries
- `GetStream`: Get stream details
- `GetLiveStreams`: List all live streams
- `GetStreamComments`: Get live comments
- `GetStreamAnalytics`: Get performance metrics

## Integration

### Events Published
- `StreamStarted`: Stream went live
- `StreamEnded`: Stream finished
- `ProductFeatured`: Product showcased
- `StreamSaleMade`: Sale during stream
- `ViewerMilestone`: 100/1000/10000 viewers

### Events Consumed
- `ProductCreated`: Allow featuring new products
- `OrderCreated`: Track stream-attributed sales

## Stream Keys

Each stream gets unique RTMP credentials:
```
rtmp://ingest.domain.com/live/{stream_key}
```

Playback URL:
```
https://cdn.domain.com/hls/{stream_id}/playlist.m3u8
```

## Analytics Tracked
- Total/unique viewers
- Peak concurrent viewers
- Average watch time
- Viewer retention rate
- Total revenue
- Product click-through rate
- Conversion rate
