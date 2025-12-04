module github.com/titan-commerce/backend/product-service

go 1.23

require (
	github.com/google/uuid v1.5.0
	github.com/titan-commerce/backend/pkg v0.0.0
	go.mongodb.org/mongo-driver v1.13.1
	google.golang.org/grpc v1.60.1
)

replace github.com/titan-commerce/backend/pkg => ../../pkg
