# Fraud Service

Real-time fraud detection and risk analysis for transactions.

## Purpose
Analyzes transactions in real-time to detect fraudulent activities using ML models and rule-based detection.

## Technology Stack
- **Database**: ClickHouse (fraud event storage, pattern analysis)
- **ML**: Anomaly detection, behavior analysis
- **API**: gRPC

## Key Features
- ✅ Real-time transaction risk scoring
- ✅ Velocity checks (too many orders)
- ✅ Device fingerprinting
- ✅ IP/location analysis
- ✅ User behavior profiling
- ✅ Chargeback tracking
- ✅ Account age validation
- ✅ Trust score calculation
- ✅ Rule-based + ML hybrid detection

## Fraud Indicators
- Velocity: Unusual order frequency
- Location: Mismatched shipping/billing
- Device: Multiple accounts, same device
- Behavior: Unusual purchase patterns
- Payment: Card testing, failed payments
- Account: New account, high-value order

## Risk Levels
- **Low** (0-0.3): Auto-approve
- **Medium** (0.3-0.6): Manual review
- **High** (0.6-0.8): Manual review required
- **Critical** (0.8-1.0): Auto-reject

## API
- `AnalyzeTransaction`: Check order for fraud
- `GetFraudCheck`: Retrieve fraud analysis
- `GetUserRiskProfile`: Get user trust score
