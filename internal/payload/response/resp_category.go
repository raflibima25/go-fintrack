package response

import "time"

type CategoryResponse struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Color           string    `json:"color"`
	IconColor       string    `json:"icon_color"`
	UsageCount      int64     `json:"usage_count"`
	UsagePercentage float64   `json:"usage_percentage"`
	UserID          uint      `json:"user_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at,omitempty"`
}

type CategoryListResponse struct {
	Categories []CategoryResponse `json:"categories"`
}
