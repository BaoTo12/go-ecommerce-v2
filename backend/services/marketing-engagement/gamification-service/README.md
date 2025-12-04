# Gamification Service 🎮

**SHOPEE COINS** - Gamification like Shopee with coins, games, and rewards.

## Features

- 💰 **Shopee Coins Wallet** (separate from real money)
- 📱 **Shake-Shake Game** (shake phone to win coins via gyroscope)
- 📅 **Daily Check-in** with streak tracking (day 1: 10 coins, day 7: 100 coins)
- 🎰 **Lucky Draw / Spin Wheel**
- 🎯 **Missions & Challenges** ("Buy 3 items this week → 100 coins")
- 💳 **Coin Redemption** for discounts (100 coins = $1 off)
- 🏆 **Leaderboards** (top coin earners)
- 🏅 **Achievement Badges**

## Coin Economy

```
Earn Coins:
- Daily check-in: 10-100 coins (based on streak)
- Shake game: 1-50 coins (random)
- Lucky draw: 5-500 coins
- Complete missions: 50-1000 coins
- Purchase rewards: 1 coin per $1 spent

Spend Coins:
- 100 coins = $1 discount
- Enter lucky draws (50 coins per entry)
- Unlock exclusive deals
```

## Shake-Shake Game

```go
// Mobile app sends gyroscope data
type ShakeRequest struct {
    UserID      string
    Intensity   float64  // Shake intensity (0-100)
    DeviceData  string   // Prevent cheating
}

// Server validates and awards coins
coinsWon := calculateCoinsFromIntensity(intensity)  // 1-50 coins
```

## Status

🚧 **Under Development** - Skeleton structure created
