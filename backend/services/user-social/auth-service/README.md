# Auth Service 🔐

Authentication and authorization with JWT, MFA, and OAuth2.

## Features

- 🔑 JWT access tokens (15 min expiry) + refresh tokens (30 days)
- 📱 MFA support (TOTP, SMS)
- 🌐 OAuth2/OIDC integration (Google, Facebook, Apple)
- 🚫 Redis for token blacklist (logout handling)
- 🛡️ Rate limiting on login attempts (prevent brute force)
- 🔒 Password hashing with bcrypt
- 🔄 Token refresh mechanism
- 📊 Session management across devices

## API

```protobuf
service AuthService {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse);
  rpc Logout(LogoutRequest) returns (LogoutResponse);
  rpc EnableMFA(EnableMFARequest) returns (EnableMFAResponse);
}
```

## Token Structure

```go
// Access Token (JWT)
{
  "sub": "user-123",      // User ID
  "email": "user@example.com",
  "roles": ["customer"],
  "cell_id": "cell-042",  // Cell assignment
  "exp": 1234567890       // 15 minutes from now
}

// Refresh Token (stored in DB + Redis)
{
  "token_id": "refresh-xyz",
  "user_id": "user-123",
  "device": "iPhone 15",
  "expires_at": "2025-02-04T00:00:00Z"  // 30 days
}
```

## Status

🚧 **Under Development** - Skeleton structure created
