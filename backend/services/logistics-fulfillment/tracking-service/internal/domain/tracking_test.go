package domain_test

import (
	"testing"
	"time"

	"github.com/titan-commerce/backend/tracking-service/internal/domain"
)

func TestNewOrderTracking(t *testing.T) {
	tracking := domain.NewOrderTracking("ORDER123", "Shopee Express")

	if tracking.OrderID != "ORDER123" {
		t.Errorf("Expected OrderID ORDER123, got %s", tracking.OrderID)
	}
	if tracking.Carrier != "Shopee Express" {
		t.Errorf("Expected Carrier Shopee Express, got %s", tracking.Carrier)
	}
	if tracking.Status != domain.StatusProcessing {
		t.Errorf("Expected Status processing, got %s", tracking.Status)
	}
	if tracking.TrackingNumber == "" {
		t.Error("Expected non-empty TrackingNumber")
	}
}

func TestOrderTracking_UpdateDriverLocation(t *testing.T) {
	tracking := domain.NewOrderTracking("ORDER123", "SPX")

	if tracking.CurrentLocation != nil {
		t.Error("Expected CurrentLocation to be nil initially")
	}

	tracking.UpdateDriverLocation(10.7769, 106.7009, "Quận 1, TP.HCM")

	if tracking.CurrentLocation == nil {
		t.Fatal("Expected CurrentLocation to be set")
	}
	if tracking.CurrentLocation.Lat != 10.7769 {
		t.Errorf("Expected Lat 10.7769, got %f", tracking.CurrentLocation.Lat)
	}
	if tracking.CurrentLocation.Lng != 106.7009 {
		t.Errorf("Expected Lng 106.7009, got %f", tracking.CurrentLocation.Lng)
	}
	if tracking.CurrentLocation.Address != "Quận 1, TP.HCM" {
		t.Errorf("Expected Address 'Quận 1, TP.HCM', got %s", tracking.CurrentLocation.Address)
	}
}

func TestOrderTracking_AddStep(t *testing.T) {
	tracking := domain.NewOrderTracking("ORDER123", "SPX")

	if len(tracking.Steps) != 0 {
		t.Error("Expected no steps initially")
	}

	tracking.AddStep("Đã lấy hàng", "Shipper đã lấy hàng", "Quận 7, TP.HCM")

	if len(tracking.Steps) != 1 {
		t.Fatalf("Expected 1 step, got %d", len(tracking.Steps))
	}

	step := tracking.Steps[0]
	if step.Status != "Đã lấy hàng" {
		t.Errorf("Expected Status 'Đã lấy hàng', got %s", step.Status)
	}
	if step.Description != "Shipper đã lấy hàng" {
		t.Errorf("Expected Description 'Shipper đã lấy hàng', got %s", step.Description)
	}
	if step.Location != "Quận 7, TP.HCM" {
		t.Errorf("Expected Location 'Quận 7, TP.HCM', got %s", step.Location)
	}
	if !step.Completed {
		t.Error("Expected Completed to be true")
	}
	if step.Time == "" {
		t.Error("Expected Time to be set")
	}
}

func TestOrderTracking_AddStep_PrependOrder(t *testing.T) {
	tracking := domain.NewOrderTracking("ORDER123", "SPX")

	tracking.AddStep("Step 1", "First step", "Location 1")
	time.Sleep(10 * time.Millisecond) // Small delay to ensure different times
	tracking.AddStep("Step 2", "Second step", "Location 2")

	if len(tracking.Steps) != 2 {
		t.Fatalf("Expected 2 steps, got %d", len(tracking.Steps))
	}

	// Newest step should be first
	if tracking.Steps[0].Status != "Step 2" {
		t.Errorf("Expected newest step first, got %s", tracking.Steps[0].Status)
	}
	if tracking.Steps[1].Status != "Step 1" {
		t.Errorf("Expected oldest step last, got %s", tracking.Steps[1].Status)
	}
}

func TestTrackingStatus_Values(t *testing.T) {
	tests := []struct {
		status domain.TrackingStatus
		value  string
	}{
		{domain.StatusProcessing, "processing"},
		{domain.StatusShipped, "shipped"},
		{domain.StatusInTransit, "in_transit"},
		{domain.StatusOutForDelivery, "out_for_delivery"},
		{domain.StatusDelivered, "delivered"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.value {
			t.Errorf("Expected %s, got %s", tt.value, tt.status)
		}
	}
}

func TestDriver_Fields(t *testing.T) {
	driver := &domain.Driver{
		ID:        "d1",
		Name:      "Nguyễn Shipper",
		Phone:     "0912345678",
		Vehicle:   "Honda Vision",
		PlateNo:   "59A1-12345",
		AvatarURL: "https://example.com/avatar.jpg",
	}

	if driver.ID != "d1" {
		t.Errorf("Expected ID d1, got %s", driver.ID)
	}
	if driver.Name != "Nguyễn Shipper" {
		t.Errorf("Expected Name 'Nguyễn Shipper', got %s", driver.Name)
	}
}

func TestLocation_Fields(t *testing.T) {
	loc := &domain.Location{
		Lat:     10.7769,
		Lng:     106.7009,
		Address: "Test Address",
	}

	if loc.Lat != 10.7769 {
		t.Errorf("Expected Lat 10.7769, got %f", loc.Lat)
	}
	if loc.Lng != 106.7009 {
		t.Errorf("Expected Lng 106.7009, got %f", loc.Lng)
	}
	if loc.Address != "Test Address" {
		t.Errorf("Expected Address 'Test Address', got %s", loc.Address)
	}
}

func BenchmarkAddStep(b *testing.B) {
	tracking := domain.NewOrderTracking("ORDER123", "SPX")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tracking.AddStep("Status", "Description", "Location")
	}
}
