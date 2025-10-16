package model

import "time"

// Trade represents a forex trade record
type Trade struct {
	ID           int       `json:"id" db:"id"`
	TradeTime    time.Time `json:"trade_time" db:"trade_time"`
	LotSize      float64   `json:"lot_size" db:"lot_size"`
	PurchaseRate float64   `json:"purchase_rate" db:"purchase_rate"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CreateTradeRequest represents the request body for creating a trade
type CreateTradeRequest struct {
	TradeTime    time.Time `json:"trade_time" binding:"required"`
	LotSize      float64   `json:"lot_size" binding:"required,gt=0"`
	PurchaseRate float64   `json:"purchase_rate" binding:"required,gt=0"`
}

// UpdateTradeRequest represents the request body for updating a trade
type UpdateTradeRequest struct {
	TradeTime    *time.Time `json:"trade_time"`
	LotSize      *float64   `json:"lot_size" binding:"omitempty,gt=0"`
	PurchaseRate *float64   `json:"purchase_rate" binding:"omitempty,gt=0"`
}
