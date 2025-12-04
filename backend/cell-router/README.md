# Cell Router

Routes users to the appropriate cell using consistent hashing.

## Features

- ✅ Consistent hashing for user-to-cell routing
- ✅ Health checking all 500 cells
- ✅ Automatic failover on cell failure
- ✅ HTTP API for routing queries
- ✅ Real-time cell status monitoring

## Algorithm

```go
cellID = Hash(userID) % 500 + 1
```

This ensures:
- Same user always routes to same cell
- Uniform distribution across cells
- Stateless routing (no lookup table needed)

## Quick Start

```bash
export SERVICE_NAME=cell-router
export LOG_LEVEL=info
export HTTP_PORT=8080
go run main.go
```

## API

### Route User to Cell
```bash
curl "http://localhost:8080/route?user_id=user-123"
# Response: {"user_id":"user-123","cell_id":42,"endpoint":"cell-042.svc.cluster.local:9000"}
```

### List All Cells
```bash
curl "http://localhost:8080/cells"
# Returns status of all 500 cells
```

### Health Check
```bash
curl "http://localhost:8080/health"
# Response: {"status":"healthy"}
```

## Architecture

```
User Request → Cell Router → Hash(userID) → Cell #42
                           ↓
                    Health Check (every 5s)
                           ↓
                    If unhealthy → Failover to Cell #43
```

## Deployment

In production, deploy as Kubernetes service:
- 3+ replicas for high availability
- Service mesh (Istio) for advanced routing
- Prometheus metrics for monitoring
