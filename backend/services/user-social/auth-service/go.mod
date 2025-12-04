module github.com/titan-commerce/backend/auth-service

go 1.23

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.5.0
	github.com/redis/go-redis/v9 v9.4.0
	github.com/titan-commerce/backend/pkg v0.0.0
	google.golang.org/grpc v1.60.1
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
)

replace github.com/titan-commerce/backend/pkg => ../../pkg
