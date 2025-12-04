# Livestream Service 🔴

**SHOPEE LIVE** - Live streaming shopping like TikTok Shop.

## Features

- 🎥 RTMP video ingestion from seller mobile app
- 🔄 Multi-bitrate transcoding (1080p, 720p, 480p, 360p) via FFmpeg
- 📺 HLS packaging for adaptive bitrate streaming
- ☁️ CDN integration (CloudFlare Stream / AWS MediaLive)
- 💬 Live chat overlay during stream (WebSocket)
- 📌 Pinned products during stream
- ⚡ Flash sale triggers during live
- 📊 Analytics: peak viewers, total views, purchases during stream

## Architecture

```
Seller Mobile App (OBS) → RTMP → Livestream Service → FFmpeg Transcoding
                                          ↓
                                    HLS Segments → S3 → CDN
                                          ↓
                                   Viewers (HLS Player)
```

## Tech Stack

- **RTMP Server**: nginx-rtmp or custom Go RTMP server
- **Transcoding**: FFmpeg
- **Packaging**: HLS segmenter
- **Storage**: S3/MinIO for HLS segments
- **CDN**: CloudFlare / AWS CloudFront
- **Chat**: WebSocket for live chat
- **Analytics**: Redis for viewer count, ClickHouse for historical data

## Complexity

⚠️ **HIGH** - This is the most complex service in the platform.

Requires:
- Video encoding pipeline
- Real-time streaming protocols
- CDN integration
- WebSocket for chat
- Complex state management

## Status

🚧 **Under Development** - Skeleton structure created
📝 **Implementation Priority**: After core transaction services
