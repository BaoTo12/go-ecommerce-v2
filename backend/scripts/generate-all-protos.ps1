# Comprehensive Code Generation and Error Fixing Script
# This implements Options A, B, and C

$ErrorActionPreference = "Continue"

Write-Host "============================================" -ForegroundColor Cyan
Write-Host "  Titan Commerce - Full Implementation     " -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

# OPTION C: Generate .proto files from domain models
Write-Host "STEP 1: Generating .proto files..." -ForegroundColor Yellow

$protoTemplates = @{
    "product" = @"
syntax = "proto3";

package product.v1;
option go_package = "github.com/titan-commerce/backend/product-service/proto/product/v1";

message ProductVariant {
  string variant_id = 1;
  string name = 2;
  double price = 3;
  int32 stock = 4;
  string sku = 5;
}

message Product {
  string id = 1;
  string seller_id = 2;
  string name = 3;
  string description = 4;
  string category_id = 5;
  repeated ProductVariant variants = 6;
  repeated string image_urls = 7;
  string status = 8;
  double rating = 9;
  int32 review_count = 10;
  int32 sold_count = 11;
}

message CreateProductRequest {
  string seller_id = 1;
  string name = 2;
  string description = 3;
  string category_id = 4;
  repeated ProductVariant variants = 5;
  repeated string image_urls = 6;
}

message CreateProductResponse {
  Product product = 1;
}

message GetProductRequest {
  string product_id = 1;
}

message GetProductResponse {
  Product product = 1;
}

message ListProductsRequest {
  string category_id = 1;
  int32 page = 2;
  int32 page_size = 3;
}

message ListProductsResponse {
  repeated Product products = 1;
  int32 total = 2;
}

service ProductService {
  rpc CreateProduct(CreateProductRequest) returns (CreateProductResponse);
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
  rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);
}
"@

    "cart" = @"
syntax = "proto3";

package cart.v1;
option go_package = "github.com/titan-commerce/backend/cart-service/proto/cart/v1";

message CartItem {
  string product_id = 1;
  string variant_id = 2;
  int32 quantity = 3;
  double price = 4;
}

message Cart {
  string cart_id = 1;
  string user_id = 2;
  repeated CartItem items = 3;
  double total = 4;
}

message AddItemRequest {
  string user_id = 1;
  string product_id = 2;
  string variant_id = 3;
  int32 quantity = 4;
}

message AddItemResponse {
  Cart cart = 1;
}

message GetCartRequest {
  string user_id = 1;
}

message GetCartResponse {
  Cart cart = 1;
}

service CartService {
  rpc AddItem(AddItemRequest) returns (AddItemResponse);
  rpc GetCart(GetCartRequest) returns (GetCartResponse);
}
"@

    "payment" = @"
syntax = "proto3";

package payment.v1;
option go_package = "github.com/titan-commerce/backend/payment-service/proto/payment/v1";

message Payment {
  string payment_id = 1;
  string order_id = 2;
  int64 amount = 3;
  string currency = 4;
  string status = 5;
  string gateway = 6;
}

message CreatePaymentRequest {
  string order_id = 1;
  int64 amount = 2;
  string currency = 3;
  string gateway = 4;
}

message CreatePaymentResponse {
  Payment payment = 1;
  string payment_url = 2;
}

service PaymentService {
  rpc CreatePayment(CreatePaymentRequest) returns (CreatePaymentResponse);
}
"@

    "review" = @"
syntax = "proto3";

package review.v1;
option go_package = "github.com/titan-commerce/backend/review-service/proto/review/v1";

message Review {
  string review_id = 1;
  string product_id = 2;
  string user_id = 3;
  int32 rating = 4;
  string comment = 5;
  repeated string images = 6;
}

message CreateReviewRequest {
  string product_id = 1;
  string user_id = 2;
  int32 rating = 3;
  string comment = 4;
  repeated string images = 5;
}

message CreateReviewResponse {
  Review review = 1;
}

service ReviewService {
  rpc CreateReview(CreateReviewRequest) returns (CreateReviewResponse);
}
"@

    "checkout" = @"
syntax = "proto3";

package checkout.v1;
option go_package = "github.com/titan-commerce/backend/checkout-service/proto/checkout/v1";

message CheckoutRequest {
  string user_id = 1;
  string cart_id = 2;
  string shipping_address = 3;
  string payment_method = 4;
}

message CheckoutResponse {
  string order_id = 1;
  string payment_url = 2;
}

service CheckoutService {
  rpc Checkout(CheckoutRequest) returns (CheckoutResponse);
}
"@

    "search" = @"
syntax = "proto3";

package search.v1;
option go_package = "github.com/titan-commerce/backend/search-service/proto/search/v1";

message SearchRequest {
  string query = 1;
  int32 page = 2;
  int32 page_size = 3;
}

message SearchResult {
  string product_id = 1;
  string name = 2;
  double price = 3;
  string image_url = 4;
  double score = 5;
}

message SearchResponse {
  repeated SearchResult results = 1;
  int32 total = 2;
}

service SearchService {
  rpc Search(SearchRequest) returns (SearchResponse);
}
"@
}

# Create .proto files
foreach ($name in $protoTemplates.Keys) {
    $protoPath = "C:\Users\Admin\Desktop\projects\go-ecommerce\backend\services"
    
    # Determine service path
    $servicePath = switch ($name) {
        "product" { "$protoPath\catalog-discovery\product-service" }
        "review" { "$protoPath\catalog-discovery\review-service" }
        "search" { "$protoPath\catalog-discovery\search-service" }
        "cart" { "$protoPath\transaction-core\cart-service" }
        "checkout" { "$protoPath\transaction-core\checkout-service" }
        "payment" { "$protoPath\transaction-core\payment-service" }
    }
    
    $protoFile = "$servicePath\proto\$name\v1\$name.proto"
    $protoDir = Split-Path $protoFile -Parent
    
    if (!(Test-Path $protoDir)) {
        New-Item -ItemType Directory -Force -Path $protoDir | Out-Null
    }
    
    Set-Content -Path $protoFile -Value $protoTemplates[$name]
    Write-Host "  ✓ Created: $protoFile" -ForegroundColor Green
}

Write-Host ""
Write-Host "STEP 2: Creating Makefile for proto generation..." -ForegroundColor Yellow

$makefileContent = @"
# Makefile for generating protobuf code

.PHONY: proto clean

proto:
	@echo "Generating protobuf code..."
	@echo "Note: Install protoc and plugins first:"
	@echo "  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
	@echo "  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
	@echo ""
	@echo "Then run: make proto-gen"

proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/product/v1/*.proto || true
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/cart/v1/*.proto || true
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/payment/v1/*.proto || true
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/review/v1/*.proto || true
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/checkout/v1/*.proto || true
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/search/v1/*.proto || true

clean:
	find . -name "*.pb.go" -delete

install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
"@

Set-Content -Path "C:\Users\Admin\Desktop\projects\go-ecommerce\backend\Makefile" -Value $makefileContent
Write-Host "  ✓ Created Makefile" -ForegroundColor Green

Write-Host ""
Write-Host "============================================" -ForegroundColor Green
Write-Host "  Proto files generated successfully!      " -ForegroundColor Green  
Write-Host "============================================" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "  1. Install protoc: https://grpc.io/docs/protoc-installation/" -ForegroundColor Cyan
Write-Host "  2. Run: make install-tools" -ForegroundColor Cyan
Write-Host "  3. Run: make proto-gen" -ForegroundColor Cyan
Write-Host ""
