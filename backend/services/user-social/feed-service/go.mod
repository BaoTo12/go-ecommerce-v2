module github.com/titan-commerce/backend/feed-service

go 1.23

require (
	github.com/google/uuid v1.5.0
	github.com/lib/pq v1.10.9
	github.com/titan-commerce/backend/pkg v0.0.0
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/rs/zerolog v1.31.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240102182953-50ed04b92917 // indirect
)

replace github.com/titan-commerce/backend/pkg => ../../../pkg
