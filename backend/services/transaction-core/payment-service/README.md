# Payment Service

Multi-gateway payment processing with split payments and escrow support.

## Features

- ✅ Multi-gateway support (Stripe, PayPal, Adyen)
- ✅ Split payments for multi-vendor orders
- ✅ Escrow management  
- ✅ Idempotency for retry safety
- ✅ PCI DSS compliance patterns
- ✅ Saga participant for distributed transactions

## Quick Start

```bash
export SERVICE_NAME=payment-service
export CELL_ID=cell-001
export DATABASE_URL=postgresql://user:pass@localhost:5432/payments
go run cmd/server/main.go
```

## API

See `proto/payment/v1/payment.proto` for API definition.

## Status

🚧 **Under Development** - Skeleton structure created
