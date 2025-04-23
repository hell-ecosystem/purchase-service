package domain

import "time"

type Purchase struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Name           string     `json:"name"`
	Category       string     `json:"category"`
	EstimatedPrice *float64   `json:"estimated_price,omitempty"`
	ActualPrice    *float64   `json:"actual_price,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	Deadline       *time.Time `json:"deadline,omitempty"`
	IsActive       bool       `json:"is_active"`
	IsPurchased    bool       `json:"is_purchased"`
}
