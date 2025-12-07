package memory

import (
	"sync"
	"time"

	"github.com/titan-commerce/backend/tracking-service/internal/domain"
)

// TrackingRepository is an in-memory implementation
type TrackingRepository struct {
	trackings map[string]*domain.OrderTracking
	mu        sync.RWMutex
}

// NewTrackingRepository creates a new in-memory repository with sample data
func NewTrackingRepository() *TrackingRepository {
	repo := &TrackingRepository{
		trackings: make(map[string]*domain.OrderTracking),
	}
	repo.seedData()
	return repo
}

func (r *TrackingRepository) seedData() {
	// Sample order tracking
	tracking := &domain.OrderTracking{
		OrderID:           "SPX2024120712345",
		Status:            domain.StatusOutForDelivery,
		EstimatedDelivery: "Hôm nay, 14:00 - 18:00",
		Carrier:           "Shopee Express",
		TrackingNumber:    "SPXVN123456789",
		CurrentLocation: &domain.Location{
			Lat:     10.7769,
			Lng:     106.7009,
			Address: "Quận 1, TP. Hồ Chí Minh",
		},
		Driver: &domain.Driver{
			ID:        "d1",
			Name:      "Nguyễn Văn Shipper",
			Phone:     "0912***456",
			Vehicle:   "Honda Vision",
			PlateNo:   "59A1-12345",
			AvatarURL: "https://ui-avatars.com/api/?name=NVS",
		},
		Steps: []domain.TrackingStep{
			{Status: "Đang giao", Description: "Shipper đang trên đường giao hàng đến bạn", Time: "10:30", Location: "Quận 1, TP.HCM", Completed: true},
			{Status: "Đến kho phát", Description: "Đơn hàng đã đến bưu cục phát", Time: "08:15", Location: "Quận 1, TP.HCM", Completed: true},
			{Status: "Đang vận chuyển", Description: "Đơn hàng đang trên đường vận chuyển", Time: "06:00", Location: "Bình Dương", Completed: true},
			{Status: "Rời kho phân loại", Description: "Đơn hàng đã rời kho phân loại", Time: "Hôm qua, 22:00", Location: "Kho Long An", Completed: true},
			{Status: "Đến kho phân loại", Description: "Đơn hàng đã đến kho phân loại", Time: "Hôm qua, 18:00", Location: "Kho Long An", Completed: true},
			{Status: "Đã lấy hàng", Description: "Shipper đã lấy hàng từ người bán", Time: "Hôm qua, 15:30", Location: "Quận 7, TP.HCM", Completed: true},
			{Status: "Đơn hàng đã xác nhận", Description: "Người bán đã xác nhận đơn hàng", Time: "Hôm qua, 14:00", Location: "", Completed: true},
		},
	}
	r.trackings["SPX2024120712345"] = tracking

	// Add another sample
	tracking2 := &domain.OrderTracking{
		OrderID:           "SPX2024120754321",
		Status:            domain.StatusInTransit,
		EstimatedDelivery: "Ngày mai, 09:00 - 12:00",
		Carrier:           "Shopee Express",
		TrackingNumber:    "SPXVN987654321",
		CurrentLocation: &domain.Location{
			Lat:     21.0285,
			Lng:     105.8542,
			Address: "Từ Liêm, Hà Nội",
		},
		Steps: []domain.TrackingStep{
			{Status: "Đang vận chuyển", Description: "Đơn hàng đang trên đường vận chuyển", Time: time.Now().Format("15:04"), Location: "Hà Nội", Completed: true},
			{Status: "Đã lấy hàng", Description: "Shipper đã lấy hàng từ người bán", Time: "Hôm qua, 16:00", Location: "Quận 1, TP.HCM", Completed: true},
		},
	}
	r.trackings["SPX2024120754321"] = tracking2
}

func (r *TrackingRepository) GetTrackingByOrderID(ctx interface{}, orderID string) (*domain.OrderTracking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if t, ok := r.trackings[orderID]; ok {
		return t, nil
	}
	return nil, nil
}

func (r *TrackingRepository) SaveTracking(ctx interface{}, tracking *domain.OrderTracking) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.trackings[tracking.OrderID] = tracking
	return nil
}

func (r *TrackingRepository) UpdateDriverLocation(ctx interface{}, orderID string, lat, lng float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if t, ok := r.trackings[orderID]; ok {
		if t.CurrentLocation == nil {
			t.CurrentLocation = &domain.Location{}
		}
		t.CurrentLocation.Lat = lat
		t.CurrentLocation.Lng = lng
	}
	return nil
}
