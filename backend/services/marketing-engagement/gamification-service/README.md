# Gamification Service

Points, badges, and rewards system for user engagement.

## Purpose
Increases user engagement through gamification mechanics (points, levels, badges, rewards).

## Technology Stack
- **Database**: PostgreSQL (user points, badges)
- **API**: gRPC

## Key Features
- ✅ Points earning system
- ✅ Points redemption for rewards
- ✅ User levels based on lifetime points
- ✅ Achievement badges
- ✅ Reward catalog
- ✅ Transaction history
- ✅ Leaderboards

## Points Earning
- Purchase: 1 point per $1 spent
- Review: 50 points
- Referral: 500 points
- Daily login: 10 points
- Social share: 25 points

## Badges
- First Purchase
- Loyal Customer (10+ orders)
- Super Reviewer (50+ reviews)
- Shopaholic (100+ orders)
- VIP Member (Level 10+)

## API
- `EarnPoints`: Award points to user
- `RedeemPoints`: Redeem reward
- `GetUserPoints`: Get balance and level
