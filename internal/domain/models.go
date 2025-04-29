package domain

import "time"

type PurchaseUpdate struct {
	Name           *string    `json:"name,omitempty"`
	CategoryID     *string    `json:"category_id,omitempty"`
	EstimatedPrice *float64   `json:"estimated_price,omitempty"`
	Deadline       *time.Time `json:"deadline,omitempty"`
}

type PurchaseStats struct {
	Total          int     `json:"total"`
	Active         int     `json:"active"`
	Purchased      int     `json:"purchased"`
	EstimatedTotal float64 `json:"estimated_total"`
	ActualTotal    float64 `json:"actual_total"`
}

type Reminder struct {
	RemindAt time.Time `json:"remind_at"`
	Message  string    `json:"message,omitempty"`
}
