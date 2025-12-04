package domain

import (
	"time"
)

type CartItem struct {
	ProductID   string
	ProductName string
	Quantity    int
	UnitPrice   float64
	Subtotal    float64
	AddedAt     time.Time
}

type Cart struct {
	UserID    string
	Items     []CartItem
	Total     float64
	UpdatedAt time.Time
}

func NewCart(userID string) *Cart {
	return &Cart{
		UserID:    userID,
		Items:     make([]CartItem, 0),
		Total:     0.0,
		UpdatedAt: time.Now(),
	}
}

func (c *Cart) AddItem(productID, productName string, quantity int, unitPrice float64) {
	// Check if item already exists
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items[i].Quantity += quantity
			c.Items[i].Subtotal = float64(c.Items[i].Quantity) * item.UnitPrice
			c.recalculateTotal()
			c.UpdatedAt = time.Now()
			return
		}
	}

	// Add new item
	subtotal := float64(quantity) * unitPrice
	c.Items = append(c.Items, CartItem{
		ProductID:   productID,
		ProductName: productName,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		Subtotal:    subtotal,
		AddedAt:     time.Now(),
	})
	c.recalculateTotal()
	c.UpdatedAt = time.Now()
}

func (c *Cart) RemoveItem(productID string) {
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			break
		}
	}
	c.recalculateTotal()
	c.UpdatedAt = time.Now()
}

func (c *Cart) UpdateQuantity(productID string, quantity int) {
	if quantity <= 0 {
		c.RemoveItem(productID)
		return
	}

	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items[i].Quantity = quantity
			c.Items[i].Subtotal = float64(quantity) * item.UnitPrice
			break
		}
	}
	c.recalculateTotal()
	c.UpdatedAt = time.Now()
}

func (c *Cart) Clear() {
	c.Items = make([]CartItem, 0)
	c.Total = 0.0
	c.UpdatedAt = time.Now()
}

func (c *Cart) recalculateTotal() {
	total := 0.0
	for _, item := range c.Items {
		total += item.Subtotal
	}
	c.Total = total
}
